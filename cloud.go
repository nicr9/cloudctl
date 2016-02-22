package main

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
)

type Cloud interface {
	listInstances()
}

type Aws struct {
	config Config
}

func NewAws(config Config) Aws {
	creds := aws.Creds(config.AccessKey, config.SecretKey, "")
	svc := ec2.New(&aws.Config{
		Credentials: creds,
		Region:      config.Region,
	})
}

func (a Aws) listInstances() {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}

	resp, _ := svc.DescribeInstances(params)
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println(inst.InstanceID)
		}
	}
}

func NewCloud(config Config) Cloud {
	return Aws{config}
}
