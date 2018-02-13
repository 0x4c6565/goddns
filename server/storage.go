package main

import "ddns/lib"

type Storage interface {
	Get(host string, recordType lib.DDNSRecordType) string
	Update(host string, ipAddress string, recordType lib.DDNSRecordType) error
}
