package account

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type ListShopAdvertiserData struct {
	List     []*uint64                `json:"list"`
	PageInfo oceanengine.PageResponse `json:"page_info"`
}

type ListShopAdvertiserResponse struct {
	oceanengine.CommonResponse
	Data ListShopAdvertiserData `json:"data"`
}

func (lsar *ListShopAdvertiserResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lsar); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/shop/advertiser/list/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if lsar.Code != 0 {
			return oceanengine.NewOceanengineError(lsar.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/shop/advertiser/list/", lsar.Message, oceanengine.ResponseDescription[lsar.Code], lsar.RequestId, response)
		}
	}

	return nil
}
