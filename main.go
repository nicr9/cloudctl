package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	default_config = `platform: aws
`
	cloud_name       = kingpin.Flag("cloud", "Name of public/private cloud").Short('c').Default("default").String()
	cmd_init         = kingpin.Command("init", "Create a configuration file for a new cloud.")
	cmd_print_config = kingpin.Command("print-config", "Print the cloudctl configuration for this cloud.")
)

func main() {
	// Initialise
	ensureConfigDir()

	// Handle cli args
	switch kingpin.Parse() {
	case "init":
		WriteConfig(*cloud_name, default_config)
	case "print-config":
		config := GetConfig(*cloud_name)
		fmt.Printf("%+v\n", config)
	}
}
