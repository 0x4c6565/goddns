package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/0x4c6565/goddns/pkg/model"
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

func (s FlatFileStorage) Get(host string, recordType model.DDNSRecordType) (model.DDNSRecord, bool) {
	record, exists := s.records[host]
	if exists {
		switch recordType {
		case model.A:
			return model.DDNSRecord{
				IPAddress: record.A,
				Type:      model.A,
			}, true
		case model.AAAA:
			return model.DDNSRecord{
				IPAddress: record.AAAA,
				Type:      model.AAAA,
			}, true
		}
	}

	return model.DDNSRecord{}, false
}

func (s FlatFileStorage) Update(host string, record model.DDNSRecord) error {
	if s.records == nil {
		s.records = map[string]*FlatFileDDNSRecord{}
	}

	if s.records[host] == nil {
		s.records[host] = &FlatFileDDNSRecord{}
	}

	switch record.Type {
	case model.A:
		s.records[host].A = record.IPAddress
		break
	case model.AAAA:
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
