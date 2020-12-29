package main

import "github.com/0x4c6565/goddns/pkg/model"

type Storage interface {
	Get(host string, recordType model.DDNSRecordType) (model.DDNSRecord, bool)
	Update(host string, record model.DDNSRecord) error
}
