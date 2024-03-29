package oauth2

import (
	"douyin/internal/pkg/openDouyin"
	"encoding/json"
)

type ClientTokenDataResponse struct {
	AccessToken string `json:"access_token"`
	Captcha     string `json:"captcha"`
	DescUrl     string `json:"desc_url"`
	Description string `json:"description"`
	ErrorCode   uint64 `json:"error_code"`
	ExpiresIn   uint64 `json:"expires_in"`
	LogId       string `json:"log_id"`
}

type ClientTokenResponse struct {
	Message string                  `json:"message"`
	Data    ClientTokenDataResponse `json:"data"`
}

func (ctr *ClientTokenResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), ctr); err != nil {
		return openDouyin.NewOpenDouyinError(90000, openDouyin.BaseDomain+"/oauth/client_token/", "解析json失败："+err.Error(), response)
	} else {
		if ctr.Message != "success" {
			return openDouyin.NewOpenDouyinError(ctr.Data.ErrorCode, openDouyin.BaseDomain+"/oauth/client_token/", openDouyin.ResponseDescription[ctr.Data.ErrorCode], response)
		}
	}

	return nil
}
