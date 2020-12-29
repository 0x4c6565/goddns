package main

type Storage interface {
	Get(host string, recordType DDNSRecordType) (DDNSRecord, bool)
	Update(host string, record DDNSRecord) error
}
