package token

import (
	"douyin/internal/pkg/douke"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

func (rtr RefreshTokenRequest) String() string {
	return douke.Marshal(rtr)
}
