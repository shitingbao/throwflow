package command

import (
	"encoding/json"
)

type ParseDataRequest struct {
	Command string `json:"command"` // 直播间/直播间商品/商品的抖口令
}

func (pdr ParseDataRequest) String() string {
	data, _ := json.Marshal(pdr)

	return string(data)
}
