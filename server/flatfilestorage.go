package main

import (
	"ddns/lib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type FlatFileDDNSRecord struct {
	A    string
	AAAA string
}

type FlatFileStorage struct {
	path    string
	mtx     sync.Mutex
	records map[string]*FlatFileDDNSRecord
}

func NewFlatFileStorage(path string) (FlatFileStorage, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return FlatFileStorage{}, err
	}

	storage := FlatFileStorage{path: path}
	json.Unmarshal(file, &storage.records)

	return storage, nil
}

func (s FlatFileStorage) Get(host string, recordType lib.DDNSRecordType) string {
	record, exists := s.records[host]
	if exists {
		switch recordType {
		case lib.A:
			return record.A
		case lib.AAAA:
			return record.AAAA
		}
	}
	return ""
}

func (s FlatFileStorage) Update(host string, ipAddress string, recordType lib.DDNSRecordType) error {
	if s.records == nil {
		s.records = map[string]*FlatFileDDNSRecord{}
	}

	if s.records[host] == nil {
		s.records[host] = &FlatFileDDNSRecord{}
	}

	switch recordType {
	case lib.A:
		s.records[host].A = ipAddress
		break
	case lib.AAAA:
		s.records[host].AAAA = ipAddress
		break
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	json, err := json.Marshal(s.records)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for storage: %s", err)
	}

	err = ioutil.WriteFile(s.path, json, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file for storage: %s", err)
	}

	return nil
}
