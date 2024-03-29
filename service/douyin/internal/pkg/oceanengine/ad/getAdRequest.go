package ad

import (
	"github.com/google/go-querystring/query"
)

type GetAdRequest struct {
	AdvertiserId uint64 `url:"advertiser_id"`
	AdId         uint64 `url:"ad_id"`
}

func (gar GetAdRequest) String() string {
	v, _ := query.Values(gar)

	return v.Encode()
}
