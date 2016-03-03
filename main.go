package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"strings"
)

var (
	default_config = `platform: aws
region: us-west-1
`
	cloud_name       = kingpin.Flag("cloud", "Name of public/private cloud").Short('c').Default("default").String()
	cmd_init         = kingpin.Command("init", "Create a configuration file for a new cloud.")
	cmd_config_print = kingpin.Command("config-print", "Print the cloudctl configuration for this cloud.")
	cmd_config_edit  = kingpin.Command("config-edit", "Edit the cloudctl configuration for this cloud with $EDITOR.")
	cmd_ls           = kingpin.Command("ls", "List the instances in this cloud.")

	cmd_show      = kingpin.Command("show", "List details about an instance.")
	show_instance = cmd_show.Arg("instance", "Target instance id.").Required().String()

	cmd_ssh       = kingpin.Command("ssh", "Sign into instance over ssh.")
	ssh_user_host = cmd_ssh.Arg("user_host", "username and instance id, e.g., centos@i-12345678. username will default to $USER").Required().String()

	cmd_rm       = kingpin.Command("rm", "Remove one or more instances.")
	rm_instances = cmd_rm.Arg("instance", "instance id").Required().Strings()
)

func main() {
	// Initialise
	ensureConfigDir()

	// Handle cli args
	switch kingpin.Parse() {
	case cmd_init.FullCommand():
		config_path := ConfigPath(*cloud_name)
		if _, err := os.Stat(config_path); os.IsNotExist(err) {
			WriteConfig(config_path, default_config)
		}
	case cmd_config_print.FullCommand():
		config := GetConfig(*cloud_name)
		fmt.Printf("%+v\n", config)
	case cmd_config_edit.FullCommand():
		config_path := ConfigPath(*cloud_name)

		cmd := exec.Command("vim", config_path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	case cmd_ls.FullCommand():
		config := GetConfig(*cloud_name)
		cloud, err := NewCloud(config)
		if err != nil {
			fmt.Println("Couldn't create cloud interface.")
			fmt.Println(err)
		}
		cloud.listInstances()
	case cmd_show.FullCommand():
		config := GetConfig(*cloud_name)
		cloud, err := NewCloud(config)
		if err != nil {
			fmt.Println("Couldn't create cloud interface.")
			fmt.Println(err)
		}
		cloud.showInstance(*show_instance)
	case cmd_ssh.FullCommand():
		user_host := strings.Split(*ssh_user_host, "@")
		var user, host string
		if len(user_host) == 1 {
			user = os.Getenv("USER")
			host = user_host[0]
		} else {
			user, host = user_host[0], user_host[1]
		}

		config := GetConfig(*cloud_name)
		cloud, err := NewCloud(config)
		if err != nil {
			fmt.Println("Couldn't create cloud interface.")
			fmt.Println(err)
		}
		cloud.sshInstance(user, host)
	case cmd_rm.FullCommand():
		config := GetConfig(*cloud_name)
		cloud, err := NewCloud(config)
		if err != nil {
			fmt.Println("Couldn't create cloud interface.")
			fmt.Println(err)
		}
		cloud.removeInstances(*rm_instances)
	}
}
