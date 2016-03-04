package main

import (
	"errors"
	"fmt"
)

type Cloud interface {
	listInstances()
	showInstance(instanceId string)
	sshInstance(username, instanceId string)
	removeInstances(instances []string)
}

func NewCloud(config Config) (Cloud, error) {
	switch config.Provider {
	case "aws":
		return NewAws(config), nil
	case "digitalocean":
		return NewDigitalOcean(config), nil
	default:
		msg := fmt.Sprintf("Unrecognised provider: %s", config.Provider)
		return nil, errors.New(msg)
	}
}
