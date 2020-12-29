package main

import (
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
	mtx     *sync.Mutex
	records map[string]*FlatFileDDNSRecord
}

func NewFlatFileStorage(path string) (FlatFileStorage, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return FlatFileStorage{}, err
	}

	storage := FlatFileStorage{path: path, mtx: &sync.Mutex{}}
	json.Unmarshal(file, &storage.records)

	return storage, nil
}

func (s FlatFileStorage) Get(host string, recordType DDNSRecordType) (DDNSRecord, bool) {
	record, exists := s.records[host]
	if exists {
		switch recordType {
		case A:
			return DDNSRecord{
				IPAddress: record.A,
				Type:      A,
			}, true
		case AAAA:
			return DDNSRecord{
				IPAddress: record.AAAA,
				Type:      AAAA,
			}, true
		}
	}

	return DDNSRecord{}, false
}

func (s FlatFileStorage) Update(host string, record DDNSRecord) error {
	if s.records == nil {
		s.records = map[string]*FlatFileDDNSRecord{}
	}

	if s.records[host] == nil {
		s.records[host] = &FlatFileDDNSRecord{}
	}

	switch record.Type {
	case A:
		s.records[host].A = record.IPAddress
		break
	case AAAA:
		s.records[host].AAAA = record.IPAddress
		break
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	json, err := json.Marshal(s.records)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for storage: %w", err)
	}

	err = ioutil.WriteFile(s.path, json, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file for storage: %w", err)
	}

	return nil
}
