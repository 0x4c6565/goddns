package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Lee303/goddns/lib"

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

func (a *API) UpdateDDNSRecordEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	var record lib.DDNSRecordBody
	_ = json.NewDecoder(req.Body).Decode(&record)

	if a.authKey != record.AuthKey {
		w.WriteHeader(403)
		return
	}

	sanitizedHost := lib.SanitizeHost(params["host"])

	ip := net.ParseIP(record.IPAddress)
	if ip == nil || !strings.HasSuffix(sanitizedHost, a.zone) {
		w.WriteHeader(400)
		return
	}

	if (record.RecordType == lib.A && ip.To4() == nil) || (record.RecordType == lib.AAAA && ip.To16() == nil) {
		w.WriteHeader(400)
		return
	}

	err := a.storage.Update(sanitizedHost, record.IPAddress, record.RecordType)
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
		panic(fmt.Errorf("Failed to start API listener: %s\n ", err))
	}
}
