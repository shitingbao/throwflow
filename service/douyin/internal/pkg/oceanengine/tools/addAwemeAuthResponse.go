package tools

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type AddAwemeAuthData struct {
	AuthSuccess bool `json:"auth_success"`
}

type AddAwemeAuthResponse struct {
	oceanengine.CommonResponse
	Data AddAwemeAuthData `json:"data"`
}

func (aaar *AddAwemeAuthResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), aaar); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/tools/aweme_auth/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if aaar.Code != 0 {
			return oceanengine.NewOceanengineError(aaar.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/tools/aweme_auth/", aaar.Message, oceanengine.ResponseDescription[aaar.Code], aaar.RequestId, response)
		}
	}

	return nil
}
