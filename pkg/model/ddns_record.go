package model

type DDNSRecordType int

const (
	A DDNSRecordType = 1 + iota
	AAAA
)

type DDNSRecord struct {
	IPAddress string
	Type      DDNSRecordType
}
