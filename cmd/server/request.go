package main

type DDNSUpdateRequest struct {
	AuthKey string     `json:"authKey,omitempty"`
	Record  DDNSRecord `json:"record,omitempty"`
}
