package main

import (
	"log"
	"os"

	"flag"

	"github.com/0x4c6565/goddns/pkg/helper"
)

func main() {
	log.SetOutput(os.Stdout)

	var configPath string
	flag.StringVar(&configPath, "config", "ddnsserver.yml", "Path to config file")
	flag.Parse()

	config, err := LoadConfig(configPath)
	check(err)

	storage, err := NewFlatFileStorage("records.json")
	check(err)

	sanitizedZone := helper.SanitizeHost(config.Zone)
	api := NewAPI(config.API.Port, config.API.AuthKey, sanitizedZone, storage)

	go api.Start()

	if config.ListenUDP {
		udpServer := NewServer(sanitizedZone, storage, "udp")
		go udpServer.Start()
	}

	if config.ListenTCP {
		tcpServer := NewServer(sanitizedZone, storage, "tcp")
		go tcpServer.Start()
	}
	select {}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
