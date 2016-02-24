package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
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

	cmd_show         = kingpin.Command("show", "List details about an instance.")
	show_instance    = cmd_show.Arg("instance", "Target instance id.").Required().String()
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
    }
}
