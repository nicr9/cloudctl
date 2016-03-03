package main

import (
	"errors"
	"fmt"
)

type Cloud interface {
	listInstances()
	showInstance(instanceId string)
	sshInstance(username, instanceId string)
}

func NewCloud(config Config) (Cloud, error) {
	switch config.Platform {
	case "aws":
		return NewAws(config), nil
	case "digitalocean":
		return NewDigitalOcean(config), nil
	default:
		msg := fmt.Sprintf("Unrecognised platform: %s", config.Platform)
		return nil, errors.New(msg)
	}
}
