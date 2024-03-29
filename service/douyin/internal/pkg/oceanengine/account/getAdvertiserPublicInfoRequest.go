package account

import (
	"github.com/google/go-querystring/query"
)

type GetAdvertiserPublicInfoRequest struct {
	AdvertiserIds string `url:"advertiser_ids"`
}

func (gapir GetAdvertiserPublicInfoRequest) String() string {
	v, _ := query.Values(gapir)

	return v.Encode()
}
