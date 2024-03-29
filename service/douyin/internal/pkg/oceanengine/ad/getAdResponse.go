package ad

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type GetAdResponseData struct {
	AdId            uint64          `json:"ad_id"`
	CampaignId      uint64          `json:"campaign_id"`
	PromotionWay    string          `json:"promotion_way"`
	MarketingGoal   string          `json:"marketing_goal"`
	CampaignScene   string          `json:"campaign_scene"`
	MarketingScene  string          `json:"marketing_scene"`
	Name            string          `json:"name"`
	Status          string          `json:"status"`
	OptStatus       string          `json:"opt_status"`
	AdCreateTime    string          `json:"ad_create_time"`
	AdModifyTime    string          `json:"ad_modify_time"`
	LabAdType       string          `json:"lab_ad_type"`
	ProductInfo     []*ProductInfo  `json:"product_info"`
	RoomInfo        []*RoomInfo     `json:"room_info"`
	AwemeInfo       []*AwemeInfo    `json:"aweme_info"`
	DeliverySetting DeliverySetting `json:"delivery_setting"`
}

type GetAdResponse struct {
	oceanengine.CommonResponse
	Data GetAdResponseData `json:"data"`
}

func (gar *GetAdResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gar); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/ad/detail/get/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if gar.Code != 0 {
			return oceanengine.NewOceanengineError(gar.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/ad/detail/get/", gar.Message, oceanengine.ResponseDescription[gar.Code], gar.RequestId, response)
		}
	}

	return nil
}
