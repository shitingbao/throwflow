package product

import (
	"github.com/google/go-querystring/query"
)

type ListAvailableProductRequest struct {
	AdvertiserId   uint64 `url:"advertiser_id"`
	MarketingScene string `url:"marketing_scene"`
	Page           uint32 `url:"page"`
	PageSize       uint32 `url:"page_size"`
}

func (lapr ListAvailableProductRequest) String() string {
	v, _ := query.Values(lapr)

	return v.Encode()
}
