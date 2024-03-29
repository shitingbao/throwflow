package token

import (
	"douyin/internal/pkg/douke"
)

type CreateTokenRequest struct {
	AuthId          string `json:"auth_id"`
	AuthSubjectType string `json:"auth_subject_type"`
	Code            string `json:"code"`
	GrantType       string `json:"grant_type"`
	ShopId          string `json:"shop_id"`
	TestShop        string `json:"test_shop"`
}

func (ctr CreateTokenRequest) String() string {
	return douke.Marshal(ctr)
}
