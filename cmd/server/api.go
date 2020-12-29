package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/0x4c6565/goddns/pkg/helper"
	"github.com/gorilla/mux"
)

type API struct {
	port    int
	authKey string
	zone    string
	storage Storage
	router  *mux.Router
}

// NewAPI returns a new instance of the API struct
func NewAPI(port int, authKey string, zone string, storage Storage) *API {
	return &API{port: port, authKey: authKey, zone: zone, storage: storage}
}

func (a *API) UpdateDDNSRecordEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var req DDNSUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if a.authKey != req.AuthKey {
		w.WriteHeader(403)
		return
	}

	sanitizedHost := helper.SanitizeHost(params["host"])

	ip := net.ParseIP(req.Record.IPAddress)
	if ip == nil || !strings.HasSuffix(sanitizedHost, a.zone) {
		w.WriteHeader(400)
		return
	}

	if (req.Record.Type == A && ip.To4() == nil) || (req.Record.Type == AAAA && ip.To16() == nil) {
		w.WriteHeader(400)
		return
	}

	err = a.storage.Update(sanitizedHost, req.Record)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(200)
}

// Start adds routes and starts the API listener
func (a *API) Start() {
	a.router = mux.NewRouter()
	a.router.HandleFunc("/ddns/{host}", a.UpdateDDNSRecordEndpoint).Methods("PUT")

	log.Printf("Starting API at %d\n", a.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", a.port), a.router)
	if err != nil {
		log.Fatalf("Failed to start API listener: %s\n ", err)
	}
}
