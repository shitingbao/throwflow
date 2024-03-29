package store

import (
	"douyin/internal/pkg/tool"
	"github.com/google/go-querystring/query"
)

type AddStoreRequest struct {
	OpenId string `url:"open_id"`
}

type AddStoreBodyDataProductsRequest struct {
	ProductId   int64  `json:"product_id"`
	ActivityUrl string `json:"activity_url"`
}

type AddStoreBodyDataRequest struct {
	Products       []*AddStoreBodyDataProductsRequest `json:"products"`
	NeedHide       bool                               `json:"need_hide"`
	PickExtra      string                             `json:"pick_extra"`
	KeepPicksource bool                               `json:"keep_picksource"`
}

func (asbdr AddStoreBodyDataRequest) String() string {
	data, _ := tool.Marshal(asbdr)

	return string(data)
}

func (asr AddStoreRequest) String() string {
	v, _ := query.Values(asr)

	return v.Encode()
}
