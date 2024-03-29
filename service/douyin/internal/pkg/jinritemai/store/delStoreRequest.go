package store

import (
	"douyin/internal/pkg/tool"
	"github.com/google/go-querystring/query"
)

type DelStoreRequest struct {
	OpenId string `url:"open_id"`
}

type DelStoreBodyDataProductsRequest struct {
	ProductId   int64 `json:"product_id"`
	PromotionId int64 `json:"promotion_id"`
}

type DelStoreBodyDataRequest struct {
	Products []*DelStoreBodyDataProductsRequest `json:"products"`
}

func (dsbdr DelStoreBodyDataRequest) String() string {
	data, _ := tool.Marshal(dsbdr)

	return string(data)
}

func (dsr DelStoreRequest) String() string {
	v, _ := query.Values(dsr)

	return v.Encode()
}
