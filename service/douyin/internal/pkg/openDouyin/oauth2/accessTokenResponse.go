package oauth2

import (
	"douyin/internal/pkg/openDouyin"
	"encoding/json"
)

type AccessTokenDataResponse struct {
	AccessToken      string `json:"access_token"`
	Captcha          string `json:"captcha"`
	DescUrl          string `json:"desc_url"`
	Description      string `json:"description"`
	ErrorCode        uint64 `json:"error_code"`
	ExpiresIn        uint64 `json:"expires_in"`
	LogId            string `json:"log_id"`
	OpenId           string `json:"open_id"`
	RefreshExpiresIn uint64 `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
	RefreshToken     string `json:"refresh_token"`
}

type AccessTokenResponse struct {
	Message string                  `json:"message"`
	Data    AccessTokenDataResponse `json:"data"`
}

func (atr *AccessTokenResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), atr); err != nil {
		return openDouyin.NewOpenDouyinError(90000, openDouyin.BaseDomain+"/oauth/access_token/", "解析json失败："+err.Error(), response)
	} else {
		if atr.Message != "success" {
			return openDouyin.NewOpenDouyinError(atr.Data.ErrorCode, openDouyin.BaseDomain+"/oauth/access_token/", openDouyin.ResponseDescription[atr.Data.ErrorCode], response)
		}
	}

	return nil
}
