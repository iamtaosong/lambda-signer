package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
		if evt.Source != "aws.autoscaling" {
			log.Printf("Not an autoscaling event")
			return errors.New("Not an autoscaling event")
		}

		config, err := lambda.LoadConfig(configFile)
		if err != nil {
			log.Printf("Config error: %v", err)
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

		certGen := cert.NewX509CertGenerator()
		cert, err := certGen.GenerateCert(&cert.Options{
			Hosts:        []string{ip},
			Org:          "kube",
			RawCAKeyPair: caBytes,
			Bits:         2048,
		})
		if err != nil {
			return err
		}

		certBytes, _ := ioutil.ReadAll(cert)
		if err != nil {
			log.Printf("Unable to read cert: %v", err)
			return err
		}

		certRS := bytes.NewReader(certBytes)

		if err := aws.PutObject(&aws.ObjectConfig{
			Body:     certRS,
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
