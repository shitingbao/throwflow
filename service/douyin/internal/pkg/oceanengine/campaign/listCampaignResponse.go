package campaign

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type ListCampaigns struct {
	Id             uint64  `json:"id"`
	Name           string  `json:"name"`
	Budget         float64 `json:"budget"`
	BudgetMode     string  `json:"budget_mode"`
	MarketingGoal  string  `json:"marketing_goal"`
	MarketingScene string  `json:"marketing_scene"`
	Status         string  `json:"status"`
	CreateDate     string  `json:"create_date"`
}

type ListCampaignData struct {
	List     []*ListCampaigns         `json:"list"`
	PageInfo oceanengine.PageResponse `json:"page_info"`
}

type ListCampaignResponse struct {
	oceanengine.CommonResponse
	Data ListCampaignData `json:"data"`
}

func (lcr *ListCampaignResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lcr); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/campaign_list/get/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if lcr.Code != 0 {
			return oceanengine.NewOceanengineError(lcr.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/campaign_list/get/", lcr.Message, oceanengine.ResponseDescription[lcr.Code], lcr.RequestId, response)
		}
	}

	return nil
}
