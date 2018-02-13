package main

import "github.com/Lee303/goddns/lib"

func main() {
	config, err := LoadConfig("ddnsserver.yml")
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
		panic(err)
	}
}
