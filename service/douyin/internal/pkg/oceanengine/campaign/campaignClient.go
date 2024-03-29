package campaign

import (
	"douyin/internal/pkg/oceanengine"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func ListCampaign(advertiserId uint64, accessToken, marketingGoal, marketingScene string, page uint32) (*ListCampaignResponse, error) {
	filter := Filter{
		MarketingGoal:  marketingGoal,
		MarketingScene: marketingScene,
		Status:         "ALL",
	}

	listCampaignRequest := ListCampaignRequest{
		AdvertiserId: advertiserId,
		Filter:       filter.String(),
		Page:         page,
		PageSize:     oceanengine.PageSize1000,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listCampaignRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/v1.0/qianchuan/campaign_list/get/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/campaign_list/get/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/campaign_list/get/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var listCampaignResponse ListCampaignResponse

	if err := listCampaignResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listCampaignResponse, nil
}
