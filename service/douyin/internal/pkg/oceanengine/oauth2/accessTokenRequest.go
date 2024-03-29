package oauth2

import (
	"douyin/internal/pkg/tool"
)

type AccessTokenRequest struct {
	AppId     string `json:"app_id"`
	Secret    string `json:"secret"`
	GrantType string `json:"grant_type"`
	AuthCode  string `json:"auth_code"`
}

func (atr AccessTokenRequest) String() string {
	data, _ := tool.Marshal(atr)

	return string(data)
}
