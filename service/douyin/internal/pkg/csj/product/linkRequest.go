package product

import (
	"encoding/json"
)

type LinkDataRequest struct {
	ProductUrl   string `json:"product_url"`
	ExternalInfo string `json:"external_info"`
	ShareType    []int  `json:"share_type"`
}

func (ldr LinkDataRequest) String() string {
	data, _ := json.Marshal(ldr)

	return string(data)
}
