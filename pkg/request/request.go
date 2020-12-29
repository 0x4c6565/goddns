package request

import "github.com/0x4c6565/goddns/pkg/model"

type DDNSUpdateRequest struct {
	AuthKey string           `json:"authKey,omitempty"`
	Record  model.DDNSRecord `json:"record,omitempty"`
}
