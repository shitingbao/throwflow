package order

import (
	"douyin/internal/pkg/tool"
	"github.com/google/go-querystring/query"
)

type ListOrderRequest struct {
	OpenId string `url:"open_id"`
}

type ListOrderBodyDataRequest struct {
	Size      int64  `json:"size"`
	Cursor    string `json:"cursor,omitempty"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	TimeType  string `json:"time_type"`
}

func (lobdr ListOrderBodyDataRequest) String() string {
	data, _ := tool.Marshal(lobdr)

	return string(data)
}

func (lor ListOrderRequest) String() string {
	v, _ := query.Values(lor)

	return v.Encode()
}
