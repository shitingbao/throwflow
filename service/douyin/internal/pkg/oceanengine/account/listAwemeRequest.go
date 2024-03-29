package account

import (
	"github.com/google/go-querystring/query"
)

type ListAwemeRequest struct {
	AdvertiserId uint64 `url:"advertiser_id"`
	Page         uint32 `url:"page"`
	PageSize     uint32 `url:"page_size"`
}

func (lar ListAwemeRequest) String() string {
	v, _ := query.Values(lar)

	return v.Encode()
}
