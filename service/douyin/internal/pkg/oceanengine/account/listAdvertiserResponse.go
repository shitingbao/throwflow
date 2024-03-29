package account

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type ListAdvertisers struct {
	AdvertiserId   uint64 `json:"advertiser_id"`
	AdvertiserName string `json:"advertiser_name"`
	IsValid        bool   `json:"is_valid"`
	AccountRole    string `json:"account_role"`
}

type ListAdvertiserData struct {
	List []*ListAdvertisers `json:"list"`
}

type ListAdvertiserResponse struct {
	oceanengine.CommonResponse
	Data ListAdvertiserData `json:"data"`
}

func (lar *ListAdvertiserResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lar); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/oauth2/advertiser/get/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if lar.Code != 0 {
			return oceanengine.NewOceanengineError(lar.Code, oceanengine.BaseDomain+"/open_api/oauth2/advertiser/get/", lar.Message, oceanengine.ResponseDescription[lar.Code], lar.RequestId, response)
		}
	}

	return nil
}
