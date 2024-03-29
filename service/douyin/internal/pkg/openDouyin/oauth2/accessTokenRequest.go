package oauth2

import (
	"douyin/internal/pkg/tool"
)

type AccessTokenRequest struct {
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	ClientKey    string `json:"client_key"`
}

func (atr AccessTokenRequest) String() string {
	data, _ := tool.Marshal(atr)

	return string(data)
}
