package account

import (
	"github.com/google/go-querystring/query"
	"weixin/internal/pkg/tool"
)

type GetPhoneNumberQueryRequest struct {
	AccessToken string `url:"access_token"`
}

type GetPhoneNumberRequest struct {
	Code string `json:"code"`
}

func (gpnqr GetPhoneNumberQueryRequest) String() string {
	v, _ := query.Values(gpnqr)

	return v.Encode()
}

func (gpnr GetPhoneNumberRequest) String() string {
	data, _ := tool.Marshal(gpnr)

	return string(data)
}
