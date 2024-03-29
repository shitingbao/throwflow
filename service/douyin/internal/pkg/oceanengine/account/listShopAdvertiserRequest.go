package account

import (
	"github.com/google/go-querystring/query"
)

type ListShopAdvertiserRequest struct {
	ShopId   uint64 `url:"shop_id"`
	Page     uint32 `url:"page"`
	PageSize uint32 `url:"page_size"`
}

func (lsar ListShopAdvertiserRequest) String() string {
	v, _ := query.Values(lsar)

	return v.Encode()
}
