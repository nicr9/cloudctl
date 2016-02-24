package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"text/tabwriter"
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
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "InstanceID\tPublic IP\tPrivateIP")
	fmt.Fprintln(w, "---\t---\t---")
    total :=0
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
            // Replace public ip with "-" if instance doesn't have one
            publicIp := "-"
            if inst.PublicIpAddress != nil {
                publicIp = *inst.PublicIpAddress
            }

            // Replace private ip with "-" if instance doesn't have one
            privateIp := "-"
            if inst.PrivateIpAddress != nil {
                privateIp = *inst.PrivateIpAddress
            }

			fmt.Fprintf(
				w,
				"%s\t%s\t%s\n",
				*inst.InstanceId,
				publicIp,
				privateIp,
			)
            total++
		}
	}
	w.Flush()
    fmt.Printf("---\nFound %d instances.\n", total)
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
