package order

import (
	"encoding/json"
)

type ListDataRequest struct {
	Size      int    `json:"size"`
	Cursor    string `json:"cursor"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	OrderType int    `json:"order_type"`
	TimeType  string `json:"time_type"`
}

func (ldr ListDataRequest) String() string {
	data, _ := json.Marshal(ldr)

	return string(data)
}
