package oauth2

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type RefreshTokenResponse struct {
	oceanengine.CommonResponse
	Data oceanengine.TokenResponse `json:"data"`
}

func (rtr *RefreshTokenResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), rtr); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/oauth2/refresh_token/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if rtr.Code != 0 {
			return oceanengine.NewOceanengineError(rtr.Code, oceanengine.BaseDomain+"/open_api/oauth2/refresh_token/", rtr.Message, oceanengine.ResponseDescription[rtr.Code], rtr.RequestId, response)
		}
	}

	return nil
}
