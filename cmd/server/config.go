package main

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Zone      string    `yaml:"zone"`
	ListenTCP bool      `yaml:"listen_tcp"`
	ListenUDP bool      `yaml:"listen_udp"`
	API       APIConfig `yaml:"api"`
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
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
