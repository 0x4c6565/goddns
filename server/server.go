package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Lee303/goddns/lib"
	"github.com/miekg/dns"
)

type Server struct {
	storage   Storage
	Zone      string
	Port      int
	dnsServer *dns.Server
}

func NewServer(zone string, port int, storage Storage) *Server {
	return &Server{Zone: zone, Port: port, storage: storage}
}

func (s *Server) parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			ip := s.storage.Get(q.Name, lib.A)
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeAAAA:
			ip := s.storage.Get(q.Name, lib.AAAA)
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func (s *Server) handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		s.parseQuery(m)
	}

	w.WriteMsg(m)
}

func (s *Server) Start() {
	dns.HandleFunc(s.Zone, s.handleDnsRequest)

	s.dnsServer = &dns.Server{Addr: ":" + strconv.Itoa(s.Port), Net: "udp"}
	log.Printf("Starting server at %d\n", s.Port)
	err := s.dnsServer.ListenAndServe()
	if err != nil {
		panic(fmt.Errorf("Failed to start server: %s\n ", err))
	}
}
