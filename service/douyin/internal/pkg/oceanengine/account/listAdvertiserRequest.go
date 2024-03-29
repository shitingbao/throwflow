package account

import (
	"github.com/google/go-querystring/query"
)

type ListAdvertiserRequest struct {
	AccessToken string `url:"access_token"`
	AppId       string `url:"app_id"`
	Secret      string `url:"secret"`
}

func (lar ListAdvertiserRequest) String() string {
	v, _ := query.Values(lar)

	return v.Encode()
}
