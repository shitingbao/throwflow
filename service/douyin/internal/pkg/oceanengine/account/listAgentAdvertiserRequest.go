package account

import (
	"github.com/google/go-querystring/query"
)

type ListAgentAdvertiserRequest struct {
	AdvertiserId uint64 `url:"advertiser_id"`
	Page         uint32 `url:"page"`
	PageSize     uint32 `url:"page_size"`
}

func (laar ListAgentAdvertiserRequest) String() string {
	v, _ := query.Values(laar)

	return v.Encode()
}
