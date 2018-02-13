package main

import "github.com/Lee303/goddns/lib"

func main() {
	config, err := LoadConfig("ddnsserver.yml")
	check(err)

	storage, err := NewFlatFileStorage("records.json")
	check(err)

	sanitizedZone := lib.SanitizeHost(config.Zone)
	api := NewAPI(config.API.Port, config.API.AuthKey, sanitizedZone, storage)
	server := NewServer(sanitizedZone, config.Port, storage)

	go api.Start()
	go server.Start()
	select {}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
