package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

// Default config
var (
	me             user.User
	config_dir     string
	default_config = `platform: aws
`
	cloud_name       = kingpin.Flag("cloud", "Name of public/private cloud").Short('c').Default("default").String()
	cmd_init         = kingpin.Command("init", "Create a configuration file for a new cloud.")
	cmd_print_config = kingpin.Command("print-config", "Print the cloudctl configuration for this cloud.")
)

type Config struct {
	Platform  string "platform"
	AccessKey string "access_key,omitempty"
	SecretKey string "secret_key,omitempty"
}

func ConfigPath(cloud_name string) string {
	config_name := fmt.Sprintf("%s.yaml", cloud_name)
	config_path := path.Join(config_dir, config_name)

	return config_path
}

func GetConfig(cloud_name string) Config {
	config_path := ConfigPath(cloud_name)

	result := Config{}
	config, err := ioutil.ReadFile(config_path)
	if err != nil {
		msg := fmt.Sprintf("Couldn't open cloud config (try running cloudctl init)")
		log.Fatal(msg)
	}

	err = yaml.Unmarshal([]byte(config), &result)
	if err != nil {
		msg := fmt.Sprintf("Couldn't parse config: %s", err)
		log.Fatal(msg)
	}

	// If not defined in config, grab AWS keys from ENV
	if result.Platform == "aws" && result.AccessKey == "" {
		result.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if result.Platform == "aws" && result.SecretKey == "" {
		result.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	return result
}

func WriteConfig(cloud_name, config_str string) {
	config_path := ConfigPath(cloud_name)
	ioutil.WriteFile(config_path, []byte(config_str), 0600)
}

func ensureConfigDir() {
	_ = os.Mkdir(config_dir, 0700)
}

func main() {
	// Initialise
	me, err := user.Current()
	if err != nil {
		fmt.Println("Can't determine user details:", err)
	}
	config_dir = path.Join(me.HomeDir, ".cloudctl")

	// Handle cli args
	switch kingpin.Parse() {
	case "init":
		ensureConfigDir()
		WriteConfig(*cloud_name, default_config)
	case "print-config":
		config := GetConfig(*cloud_name)
		fmt.Printf("%+v\n", config)
	}
}
