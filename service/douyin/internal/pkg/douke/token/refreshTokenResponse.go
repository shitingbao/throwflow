package token

import (
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type RefreshTokenDataResponse struct {
	AccessToken     string `json:"access_token"`
	ExpiresIn       uint64 `json:"expires_in"`
	RefreshToken    string `json:"refresh_token"`
	Scope           string `json:"scope"`
	ShopId          uint64 `json:"shop_id"`
	ShopName        string `json:"shop_name"`
	AuthorityId     string `json:"authority_id"`
	AuthSubjectType string `json:"auth_subject_type"`
}

type RefreshTokenResponse struct {
	douke.CommonResponse
	Data RefreshTokenDataResponse `json:"data"`
}

func (rtr *RefreshTokenResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), rtr); err != nil {
		return douke.NewDoukeError(90100, douke.BaseDomain+"/token/refresh", "解析json失败："+err.Error(), response)
	} else {
		if rtr.CommonResponse.Code != 10000 {
			return douke.NewDoukeError(rtr.CommonResponse.Code, douke.BaseDomain+"/token/refresh", rtr.CommonResponse.Msg, response)
		}
	}

	return nil
}
