package aws

import (
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ObjectConfig contains data for objects
type ObjectConfig struct {
	Body     io.ReadSeeker
	Bucket   string
	KMSKeyID string
	Key      string
	Region   string
}

// GetObject returns an object from S3
func GetObject(cfg *ObjectConfig) (io.Reader, error) {
	svc := s3.New(session.New(), aws.NewConfig().WithRegion(cfg.Region))
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(cfg.Key),
	})

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// PutObject puts an object into S3
func PutObject(cfg *ObjectConfig) error {
	svc := s3.New(session.New(), aws.NewConfig().WithRegion(cfg.Region))
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(cfg.Bucket),
		Key:                  aws.String(cfg.Key),
		Body:                 cfg.Body,
		ServerSideEncryption: aws.String("aws:kms"),
		SSEKMSKeyId:          aws.String(cfg.KMSKeyID),
	})

	return err
}

// InstanceIP returns the IP address of an instance
func InstanceIP(ID string, region string) (string, error) {
	svc := ec2.New(session.New(), aws.NewConfig().WithRegion(region))
	resp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(ID)},
	})

	if err != nil {
		return "", nil
	}

	var ip string

	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			ip = *i.PrivateIpAddress
		}
	}

	if len(ip) == 0 {
		return "", errors.New("No IP found")
	}

	return ip, nil
}
