package oauth2

import (
	"github.com/google/go-querystring/query"
)

type GetAccessTokenRequest struct {
	GrantType string `url:"grant_type"`
	Appid     string `url:"appid"`
	Secret    string `url:"secret"`
}

func (gatr GetAccessTokenRequest) String() string {
	v, _ := query.Values(gatr)

	return v.Encode()
}
