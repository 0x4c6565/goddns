package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Config struct {
	Zone string    `yaml:"zone"`
	Port int       `yaml:"port"`
	API  APIConfig `yaml:"api"`
}

type APIConfig struct {
	Port    int    `yaml:"port"`
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
