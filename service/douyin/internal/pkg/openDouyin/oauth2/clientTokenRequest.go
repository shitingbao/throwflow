package oauth2

import (
	"douyin/internal/pkg/tool"
)

type ClientTokenRequest struct {
	ClientKey    string `json:"client_key"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

func (ctr ClientTokenRequest) String() string {
	data, _ := tool.Marshal(ctr)

	return string(data)
}
