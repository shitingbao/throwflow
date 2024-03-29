package oauth2

import (
	"douyin/internal/pkg/tool"
)

type RefreshTokenRequest struct {
	AppId        string `json:"app_id"`
	Secret       string `json:"secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

func (rtr RefreshTokenRequest) String() string {
	data, _ := tool.Marshal(rtr)

	return string(data)
}
