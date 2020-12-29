package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	server "github.com/0x4c6565/goddns/cmd/server"
)

type Poller struct {
	IPHost        string
	Host          string
	APIURL        string
	APIAuthKey    string
	RecordType    server.DDNSRecordType
	lastIPAddress string
	Interval      time.Duration
}

func main() {
	config, err := LoadConfig("ddnsclient.yml")
	check(err)

	interval, err := time.ParseDuration(config.Interval)
	check(err)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	go func() {
		_ = <-signals
		log.Print("Terminating")
		os.Exit(0)
	}()

	if !config.IPv4 && !config.IPv6 {
		panic("at least one protocol should be provided")
	}

	if config.IPv4 {
		ipv4Poller := &Poller{
			IPHost:     config.IPHost,
			Host:       config.Host,
			APIURL:     config.API.URL,
			APIAuthKey: config.API.AuthKey,
			RecordType: server.A,
			Interval:   interval,
		}

		go ipv4Poller.Start()
	}

	if config.IPv6 {
		ipv6Poller := &Poller{
			IPHost:     config.IPHost,
			Host:       config.Host,
			APIURL:     config.API.URL,
			APIAuthKey: config.API.AuthKey,
			RecordType: server.AAAA,
			Interval:   interval,
		}

		go ipv6Poller.Start()
	}

	select {}
}

func (p *Poller) Start() {
	first := true
	for {

		if first != true {
			time.Sleep(p.Interval)
		}
		first = false

		newIPAddress, err := p.getIPAddress()
		if err != nil {
			log.Printf("failed to retrieve ip address: %s", err)
			continue
		}

		if newIPAddress != p.lastIPAddress {

			log.Printf("updating DDNS host %s with new IP address %s", p.Host, newIPAddress)

			err = p.updateDDNSHost(newIPAddress)
			if err != nil {
				log.Printf("failed to update DDNS host: %s", err)
				continue
			}

			log.Print("successfully updated DDNS host")

			p.lastIPAddress = newIPAddress
		}
	}
}

func (p *Poller) updateDDNSHost(ipAddress string) error {

	payloadJSON, err := json.Marshal(server.DDNSUpdateRequest{
		AuthKey: p.APIAuthKey,
		Record: server.DDNSRecord{
			IPAddress: ipAddress,
			Type:      p.RecordType,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %s", err)
	}

	client := http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/ddns/%s", p.APIURL, p.Host), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unsuccessful request: statuscode=%d", resp.StatusCode)
	}

	return nil
}

func (p *Poller) getIPAddress() (string, error) {
	var ipAddress string
	var err error
	switch p.RecordType {
	case server.A:
		ipAddr, err := net.ResolveIPAddr("ip4", p.IPHost)
		if err != nil {
			return "", err
		}

		ipAddress = ipAddr.IP.String()
		break
	case server.AAAA:
		ipAddr, err := net.ResolveIPAddr("ip6", p.IPHost)
		if err != nil {
			return "", err
		}
		ipAddress = fmt.Sprintf("[%s]", ipAddr.IP.String())
		break
	default:
		return "", errors.New("invalid DDNS record type")
	}

	client := http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest("GET", "http://"+ipAddress, nil)
	if err != nil {
		return "", err
	}
	req.Host = p.IPHost
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)

	}

	return strings.TrimSpace(string(bodyBytes)), nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
