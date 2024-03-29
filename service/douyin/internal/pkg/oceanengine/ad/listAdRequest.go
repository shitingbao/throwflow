package ad

import (
	"douyin/internal/pkg/tool"
	"github.com/google/go-querystring/query"
)

type FilteringNotAdModifyTime struct {
	CampaignScene  []string `json:"campaign_scene"`
	MarketingGoal  string   `json:"marketing_goal"`
	Status         string   `json:"status"`
	MarketingScene string   `json:"marketing_scene"`
}

type Filtering struct {
	CampaignScene  []string `json:"campaign_scene"`
	MarketingGoal  string   `json:"marketing_goal"`
	Status         string   `json:"status"`
	MarketingScene string   `json:"marketing_scene"`
	AdModifyTime   string   `json:"ad_modify_time"`
}

type ListAdRequest struct {
	AdvertiserId     uint64 `url:"advertiser_id"`
	RequestAwemeInfo uint8  `url:"request_aweme_info"`
	Filtering        string `url:"filtering"`
	Page             uint32 `url:"page"`
	PageSize         uint32 `url:"page_size"`
}

func (f Filtering) String() string {
	data, _ := tool.Marshal(f)

	return string(data)
}

func (f FilteringNotAdModifyTime) String() string {
	data, _ := tool.Marshal(f)

	return string(data)
}

func (lar ListAdRequest) String() string {
	v, _ := query.Values(lar)

	return v.Encode()
}
