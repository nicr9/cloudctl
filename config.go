package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

var config_dir string

type Config struct {
	Platform  string "platform"
	AccessKey string "access_key,omitempty"
	SecretKey string "secret_key,omitempty"
	Region    string "region"
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

func WriteConfig(config_path, config_str string) {
	ioutil.WriteFile(config_path, []byte(config_str), 0600)
}

func ensureConfigDir() {
	me, err := user.Current()
	if err != nil {
		fmt.Println("Can't determine user details:", err)
	}
	config_dir = path.Join(me.HomeDir, ".cloudctl")
	_ = os.Mkdir(config_dir, 0700)
}
