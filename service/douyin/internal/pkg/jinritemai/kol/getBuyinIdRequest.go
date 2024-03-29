package kol

import (
	"github.com/google/go-querystring/query"
)

type GetBuyinIdRequest struct {
	OpenId string `url:"open_id"`
}

func (gbir GetBuyinIdRequest) String() string {
	v, _ := query.Values(gbir)

	return v.Encode()
}
