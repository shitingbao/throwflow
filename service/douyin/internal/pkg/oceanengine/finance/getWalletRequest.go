package finance

import (
	"github.com/google/go-querystring/query"
)

type GetWalletRequest struct {
	AdvertiserId uint64 `url:"advertiser_id"`
}

func (gwr GetWalletRequest) String() string {
	v, _ := query.Values(gwr)

	return v.Encode()
}
