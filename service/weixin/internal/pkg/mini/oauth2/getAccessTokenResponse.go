package oauth2

import (
	"encoding/json"
	"weixin/internal/pkg/mini"
)

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint64 `json:"expires_in"`
}

func (gatr *GetAccessTokenResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gatr); err != nil {
		return mini.NewMiniError("获取 token 失败")
	}

	return nil
}
