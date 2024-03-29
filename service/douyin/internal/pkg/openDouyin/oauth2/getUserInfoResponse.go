package oauth2

import (
	"douyin/internal/pkg/openDouyin"
	"encoding/json"
)

type GetUserInfoDataResponse struct {
	Captcha      string `json:"captcha"`
	DescUrl      string `json:"desc_url"`
	Description  string `json:"description"`
	ErrorCode    uint64 `json:"error_code"`
	LogId        string `json:"log_id"`
	Avatar       string `json:"avatar"`
	AvatarLarger string `json:"avatar_larger"`
	City         string `json:"city"`
	ClientKey    string `json:"client_key"`
	Country      string `json:"country"`
	District     string `json:"district"`
	EAccountRole string `json:"e_account_role"`
	Gender       uint8  `json:"gender"`
	Nickname     string `json:"nickname"`
	OpenId       string `json:"open_id"`
	Province     string `json:"province"`
	UnionId      string `json:"union_id"`
}

type GetUserInfoResponse struct {
	Message string                  `json:"message"`
	Data    GetUserInfoDataResponse `json:"data"`
}

func (guir *GetUserInfoResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), guir); err != nil {
		return openDouyin.NewOpenDouyinError(90000, openDouyin.BaseDomain+"/oauth/userinfo/", "解析json失败："+err.Error(), response)
	} else {
		if guir.Message != "success" {
			return openDouyin.NewOpenDouyinError(guir.Data.ErrorCode, openDouyin.BaseDomain+"/oauth/userinfo/", openDouyin.ResponseDescription[guir.Data.ErrorCode], response)
		}
	}

	return nil
}
