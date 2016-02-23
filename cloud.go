package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Cloud interface {
	listInstances()
}

type Aws struct {
	config Config
	svc    *ec2.EC2
}

func NewAws(config Config) Aws {
	creds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
	region := &config.Region
	sess := session.New(&aws.Config{
		Credentials: creds,
		Region:      region,
	})
	svc := ec2.New(sess)

	return Aws{config, svc}
}

func (a Aws) listInstances() {
	params := &ec2.DescribeInstancesInput{}

	resp, _ := a.svc.DescribeInstances(params)
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println(*inst.InstanceId)
		}
	}
}

func NewCloud(config Config) (Cloud, error) {
	switch config.Platform {
	case "aws":
		return NewAws(config), nil
	default:
		msg := fmt.Sprintf("Unrecognised platform: %s", config.Platform)
		return nil, errors.New(msg)
	}
}
