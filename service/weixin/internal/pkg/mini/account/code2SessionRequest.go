package account

import (
	"github.com/google/go-querystring/query"
)

type Code2SessionRequest struct {
	Appid     string `url:"appid"`
	Secret    string `url:"secret"`
	JsCode    string `url:"js_code"`
	GrantType string `url:"grant_type"`
}

func (csr Code2SessionRequest) String() string {
	v, _ := query.Values(csr)

	return v.Encode()
}
