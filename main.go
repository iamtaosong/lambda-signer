package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/cloudwatch"

	"github.com/sthulb/lambda-signer/aws"
	"github.com/sthulb/lambda-signer/cert"
	"github.com/sthulb/lambda-signer/lambda"
)

func main() {
	cloudwatch.HandleFunc(func(evt *cloudwatch.Event, ctx *apex.Context) error {
		// set log prefix to the ID
		log.SetPrefix(fmt.Sprintf("%s ", evt.ID))

		if evt.Source != "aws.autoscaling" {
			msg := fmt.Sprintf("Not an expected AutoScaling event, got %q instead.", evt.Source)
			log.Println(msg)
			return errors.New(msg)
		}

		log.Printf("Loading configuration file: %s/config.json", ctx.FunctionName)

		configR, err := aws.GetObject(&aws.ObjectConfig{
			Bucket: ctx.FunctionName,
			Key:    "config.json",
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
		if len(ip) == 0 || err != nil {
			log.Printf("Unable to find instance IP %q: %v", details.EC2InstanceID, err)
			return err
		}

		log.Printf("Found IP address for instance %q: %v", details.EC2InstanceID, ip)

		log.Printf("Pulling CA certificate from '%s/ca.pem'", ctx.FunctionName)

		body, err := aws.GetObject(&aws.ObjectConfig{
			Bucket: ctx.FunctionName,
			Key:    "ca.pem",
			Region: evt.Region,
		})
		if err != nil {
			log.Printf("Unable to get object from %s/%s: %v", ctx.FunctionName, "ca.pem", err)
			return err
		}

		caBytes, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		log.Printf("Generating new certificate for %q", details.EC2InstanceID)
		keyPair, err := cert.GenerateX509KeyPair(&cert.Options{
			Hosts:        []string{ip},
			Org:          ctx.FunctionName,
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

		log.Printf("Storing certificate for %q as: %s/%s.pem", details.EC2InstanceID, ctx.FunctionName, details.EC2InstanceID)
		if err := aws.PutObject(&aws.ObjectConfig{
			Body:     keyPairRS,
			Bucket:   ctx.FunctionName,
			KMSKeyID: config.KMSKeyID,
			Key:      fmt.Sprintf("%s.pem", details.EC2InstanceID),
			Region:   evt.Region,
		}); err != nil {
			return err
		}

		return nil
	})
}
