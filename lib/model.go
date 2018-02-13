package lib

type DDNSRecordBody struct {
	AuthKey    string         `json:"authKey,omitempty"`
	IPAddress  string         `json:"ipAddress,omitempty"`
	RecordType DDNSRecordType `json:"type,omitempty"`
}
