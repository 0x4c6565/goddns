package main

import "github.com/Lee303/goddns/lib"

type Storage interface {
	Get(host string, recordType lib.DDNSRecordType) string
	Update(host string, ipAddress string, recordType lib.DDNSRecordType) error
}
