package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"os/exec"
	"os/user"
	"path"
	"text/tabwriter"
)

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
	total := 0
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

func (a Aws) showInstance(instanceId string) {
	inst := a.getInstance(instanceId)
	if inst != nil {
		fmt.Printf("%#v\n", inst)
	} else {
		fmt.Printf("Couldn't find %s\n", instanceId)
	}
}

func (a Aws) sshInstance(username, instanceId string) {
	inst := a.getInstance(instanceId)
	if inst != nil {
		userHost := fmt.Sprintf("%s@%s", username, *inst.PrivateIpAddress)

		me, err := user.Current()
		if err != nil {
			fmt.Println("Can't determine username details:", err)
		}
		keyFile := fmt.Sprintf(".ssh/%s.pem", *inst.KeyName)
		keyFile = path.Join(me.HomeDir, keyFile)

		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			fmt.Println("Can't find private key:", keyFile)
			return
		}

		cmd := exec.Command("ssh", "-i", keyFile, userHost)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("Couldn't find %s\n", instanceId)
	}
}

func (a Aws) getInstance(instanceId string) *ec2.Instance {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					&instanceId,
				},
			},
		},
	}

	resp, _ := a.svc.DescribeInstances(params)
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			return inst
		}
	}
	return nil
}