package campaign

import (
	"douyin/internal/pkg/tool"
	"github.com/google/go-querystring/query"
)

type Filter struct {
	MarketingGoal  string `json:"marketing_goal"`
	MarketingScene string `json:"marketing_scene"`
	Status         string `json:"status"`
}

type ListCampaignRequest struct {
	AdvertiserId uint64 `url:"advertiser_id"`
	Filter       string `url:"filter"`
	Page         uint32 `url:"page"`
	PageSize     uint32 `url:"page_size"`
}

func (f Filter) String() string {
	data, _ := tool.Marshal(f)

	return string(data)
}

func (lcr ListCampaignRequest) String() string {
	v, _ := query.Values(lcr)

	return v.Encode()
}
