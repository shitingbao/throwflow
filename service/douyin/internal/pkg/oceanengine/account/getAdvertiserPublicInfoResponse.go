package account

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type GetAdvertiserPublicInfoData struct {
	Id                 uint64 `json:"id"`
	Name               string `json:"name"`
	Company            string `json:"company"`
	FirstIndustryName  string `json:"first_industry_name"`
	SecondIndustryName string `json:"second_industry_name"`
}

type GetAdvertiserPublicInfoResponse struct {
	oceanengine.CommonResponse
	Data []*GetAdvertiserPublicInfoData `json:"data"`
}

func (gapir *GetAdvertiserPublicInfoResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gapir); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/2/advertiser/public_info/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if gapir.Code != 0 {
			return oceanengine.NewOceanengineError(gapir.Code, oceanengine.BaseDomain+"/open_api/2/advertiser/public_info/", gapir.Message, oceanengine.ResponseDescription[gapir.Code], gapir.RequestId, response)
		}
	}

	return nil
}
