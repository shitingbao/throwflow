package product

import (
	"encoding/json"
)

type CategoryDataRequest struct {
	ParentId uint64 `json:"parent_id"`
}

func (cdr CategoryDataRequest) String() string {
	data, _ := json.Marshal(cdr)

	return string(data)
}
