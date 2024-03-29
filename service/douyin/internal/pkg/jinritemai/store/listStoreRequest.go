package store

import (
	"douyin/internal/pkg/tool"
	"github.com/google/go-querystring/query"
)

type ListStoreRequest struct {
	OpenId string `url:"open_id"`
}

type ListStoreBodyDataRequest struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

func (lobdr ListStoreBodyDataRequest) String() string {
	data, _ := tool.Marshal(lobdr)

	return string(data)
}

func (lsr ListStoreRequest) String() string {
	v, _ := query.Values(lsr)

	return v.Encode()
}
