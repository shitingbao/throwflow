package oauth2

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type AccessTokenResponse struct {
	oceanengine.CommonResponse
	Data oceanengine.TokenResponse `json:"data"`
}

func (atr *AccessTokenResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), atr); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/oauth2/access_token/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if atr.Code != 0 {
			return oceanengine.NewOceanengineError(atr.Code, oceanengine.BaseDomain+"/open_api/oauth2/access_token/", atr.Message, oceanengine.ResponseDescription[atr.Code], atr.RequestId, response)
		}
	}

	return nil
}
