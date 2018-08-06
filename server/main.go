package main

import (
	"log"
	"os"

	"github.com/Lee303/goddns/lib"
	"flag"
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

	sanitizedZone := lib.SanitizeHost(config.Zone)
	api := NewAPI(config.API.Port, config.API.AuthKey, sanitizedZone, storage)
	udpServer := NewServer(sanitizedZone, storage, "udp")
	tcpServer := NewServer(sanitizedZone, storage, "tcp")

	go api.Start()
	go udpServer.Start()
	go tcpServer.Start()
	select {}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
