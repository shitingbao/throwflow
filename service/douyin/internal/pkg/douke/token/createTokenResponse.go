package token

import (
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type CreateTokenDataResponse struct {
	AccessToken     string `json:"access_token"`
	ExpiresIn       uint64 `json:"expires_in"`
	RefreshToken    string `json:"refresh_token"`
	Scope           string `json:"scope"`
	ShopId          uint64 `json:"shop_id"`
	ShopName        string `json:"shop_name"`
	AuthorityId     string `json:"authority_id"`
	AuthSubjectType string `json:"auth_subject_type"`
	EncryptOperator string `json:"encrypt_operator"`
	OperatorName    string `json:"operator_name"`
}

type CreateTokenResponse struct {
	douke.CommonResponse
	Data CreateTokenDataResponse `json:"data"`
}

func (ctr *CreateTokenResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), ctr); err != nil {
		return douke.NewDoukeError(90100, douke.BaseDomain+"/token/create", "解析json失败："+err.Error(), response)
	} else {
		if ctr.CommonResponse.Code != 10000 {
			return douke.NewDoukeError(ctr.CommonResponse.Code, douke.BaseDomain+"/token/create", ctr.CommonResponse.Msg, response)
		}
	}

	return nil
}
