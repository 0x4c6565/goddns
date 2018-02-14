package main

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Host     string    `yaml:"host"`
	IPHost   string    `yaml:"ip_host"`
	Interval string    `yaml:"interval"`
	IPv4     bool      `yaml:"ipv4"`
	IPv6     bool      `yaml:"ipv6"`
	API      APIConfig `yaml:"api"`
}

type APIConfig struct {
	URL     string `yaml:"url"`
	AuthKey string `yaml:"auth_key"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal([]byte(file), &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %s\n", err)
	}

	return config, nil
}
