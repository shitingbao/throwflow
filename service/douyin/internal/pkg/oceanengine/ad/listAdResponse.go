package ad

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type ProductInfo struct {
	Id                  uint64  `json:"id"`
	Name                string  `json:"name"`
	DiscountPrice       float64 `json:"discount_price"`
	Img                 string  `json:"img"`
	MarketPrice         float64 `json:"market_price"`
	DiscountLowerPrice  float64 `json:"discount_lower_price"`
	DiscountHigherPrice float64 `json:"discount_higher_price"`
	ChannelId           uint64  `json:"channel_id"`
	ChannelType         string  `json:"channel_type"`
}

type RoomInfo struct {
	RoomTitle    string `json:"room_title"`
	RoomStatus   string `json:"room_status"`
	AnchorId     uint64 `json:"anchor_id"`
	AnchorName   string `json:"anchor_name"`
	AnchorAvatar string `json:"anchor_avatar"`
}

type AwemeInfo struct {
	AwemeId     uint64 `json:"aweme_id"`
	AwemeName   string `json:"aweme_name"`
	AwemeShowId string `json:"aweme_show_id"`
	AwemeAvatar string `json:"aweme_avatar"`
}

type DeliverySetting struct {
	DeepExternalAction string  `json:"deep_external_action"`
	DeepBidType        string  `json:"deep_bid_type"`
	RoiGoal            float64 `json:"roi_goal"`
	SmartBidType       string  `json:"smart_bid_type"`
	ExternalAction     string  `json:"external_action"`
	Budget             float64 `json:"budget"`
	ReviveBudget       float64 `json:"revive_budget"`
	BudgetMode         string  `json:"budget_mode"`
	CpaBid             float64 `json:"cpa_bid"`
	StartTime          string  `json:"start_time"`
	EndTime            string  `json:"end_time"`
	ProductNewOpen     bool    `json:"product_new_open"`
	QcpxMode           string  `json:"qcpx_mode"`
	AllowQcpx          bool    `json:"allow_qcpx"`
}

type ListAds struct {
	AdId           uint64 `json:"ad_id"`
	CampaignId     uint64 `json:"campaign_id"`
	CampaignScene  string `json:"campaign_scene"`
	MarketingGoal  string `json:"marketing_goal"`
	MarketingScene string `json:"marketing_scene"`
	PromotionWay   string `json:"promotion_way"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	OptStatus      string `json:"opt_status"`
	AdCreateTime   string `json:"ad_create_time"`
	AdModifyTime   string `json:"ad_modify_time"`
	LabAdType      string `json:"lab_ad_type"`

	ProductInfo     []*ProductInfo  `json:"product_info"`
	AwemeInfo       []*AwemeInfo    `json:"aweme_info"`
	DeliverySetting DeliverySetting `json:"delivery_setting"`
}

type ListAdResponseData struct {
	List     []*ListAds               `json:"list"`
	FailList []uint64                 `json:"fail_list"`
	PageInfo oceanengine.PageResponse `json:"page_info"`
}

type ListAdResponse struct {
	oceanengine.CommonResponse
	Data ListAdResponseData `json:"data"`
}

func (lar *ListAdResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lar); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/ad/get/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if lar.Code != 0 {
			return oceanengine.NewOceanengineError(lar.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/ad/get/", lar.Message, oceanengine.ResponseDescription[lar.Code], lar.RequestId, response)
		}
	}

	return nil
}
