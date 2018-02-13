package lib

type DDNSRecordType int

const (
	A DDNSRecordType = 1 + iota
	AAAA
)
