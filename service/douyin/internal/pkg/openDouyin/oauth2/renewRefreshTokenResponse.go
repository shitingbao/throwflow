package oauth2

import (
	"douyin/internal/pkg/openDouyin"
	"encoding/json"
)

type RenewRefreshTokenDataResponse struct {
	Captcha      string `json:"captcha"`
	DescUrl      string `json:"desc_url"`
	Description  string `json:"description"`
	ErrorCode    uint64 `json:"error_code"`
	LogId        string `json:"log_id"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint64 `json:"expires_in"`
}

type RenewRefreshTokenResponse struct {
	Message string                        `json:"message"`
	Data    RenewRefreshTokenDataResponse `json:"data"`
}

func (rrtr *RenewRefreshTokenResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), rrtr); err != nil {
		return openDouyin.NewOpenDouyinError(90000, openDouyin.BaseDomain+"/oauth/renew_refresh_token/", "解析json失败："+err.Error(), response)
	} else {
		if rrtr.Message != "success" {
			return openDouyin.NewOpenDouyinError(rrtr.Data.ErrorCode, openDouyin.BaseDomain+"/oauth/renew_refresh_token/", openDouyin.ResponseDescription[rrtr.Data.ErrorCode], response)
		}
	}

	return nil
}
