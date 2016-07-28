package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/cloudwatch"

	"github.com/sthulb/lambda-signer/aws"
	"github.com/sthulb/lambda-signer/cert"
	"github.com/sthulb/lambda-signer/lambda"
)

var (
	// config file location
	configFile string
)

func init() {
	flagSet := flag.NewFlagSet("lambda", flag.ContinueOnError)
	flagSet.StringVar(&configFile, "config-file", lambda.DefaultConfigFile, "Location of config file")
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}

func main() {
	log.Printf("Using log file: %s", configFile)

	cloudwatch.HandleFunc(func(evt *cloudwatch.Event, ctx *apex.Context) error {
		// set log prefix to the ID
		log.SetPrefix(fmt.Sprintf("%s ", evt.ID))

		if evt.Source != "aws.autoscaling" {
			msg := fmt.Sprintf("Not an expected AutoScaling event, got %q instead.", evt.Source)
			log.Println(msg)
			return errors.New(msg)
		}

		log.Printf("Loading configuration file: %v", configFile)

		configLoader, err := lambda.LoadConfig(configFile)
		if err != nil {
			log.Printf("Config error: %v", err)
			return err
		}

		u, err := url.Parse(configLoader.ConfigURL)
		if err != nil {
			log.Printf("Unable to parse %q: %v", configLoader.ConfigURL, err)
			return err
		}

		configR, err := aws.GetObject(&aws.ObjectConfig{
			Bucket: u.Host,
			Key:    u.Path,
			Region: evt.Region,
		})

		config := &lambda.Config{}
		if err := config.ReadFromReader(configR); err != nil {
			return err
		}

		var details cloudwatch.AutoScalingGroupDetail
		if err := json.Unmarshal(evt.Detail, &details); err != nil {
			log.Printf("Unable to unmarshal detail body: %v", err)
			return err
		}

		ip, err := aws.InstanceIP(details.EC2InstanceID, evt.Region)
		if err != nil {
			log.Printf("Unable to find instance IP %q: %v", details.EC2InstanceID, err)
			return err
		}

		log.Printf("Found IP address for instance %q: %v", details.EC2InstanceID, ip)

		log.Printf("Pulling CA certificate from '%s/ca.pem'", config.Bucket)

		body, err := aws.GetObject(&aws.ObjectConfig{
			Bucket: config.Bucket,
			Key:    "ca.pem",
			Region: evt.Region,
		})
		if err != nil {
			log.Printf("Unable to get object from %s/%s: %v", config.Bucket, "ca.pem", err)
			return err
		}

		caBytes, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		log.Printf("Generating new certificate for %q", details.EC2InstanceID)
		keyPair, err := cert.GenerateX509KeyPair(&cert.Options{
			Hosts:        []string{ip},
			Org:          config.EnvironmentName,
			RawCAKeyPair: caBytes,
			Bits:         2048,
		})
		if err != nil {
			return err
		}

		keyPairBytes, _ := ioutil.ReadAll(keyPair)
		if err != nil {
			log.Printf("Unable to read cert: %v", err)
			return err
		}

		keyPairRS := bytes.NewReader(keyPairBytes)

		log.Printf("Storing certificate for %q as: %s/%s.pem", details.EC2InstanceID, config.Bucket, details.EC2InstanceID)
		if err := aws.PutObject(&aws.ObjectConfig{
			Body:     keyPairRS,
			Bucket:   config.Bucket,
			KMSKeyID: config.KMSKeyID,
			Key:      fmt.Sprintf("%s.pem", details.EC2InstanceID),
			Region:   evt.Region,
		}); err != nil {
			return err
		}

		return nil
	})
}
