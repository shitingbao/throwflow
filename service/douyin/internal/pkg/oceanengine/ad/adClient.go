package ad

import (
	"douyin/internal/pkg/oceanengine"
	"github.com/go-resty/resty/v2"
	"strconv"
	"unicode/utf8"
)

func ListAd(advertiserId uint64, accessToken, marketingGoal, modifyTime string, page uint32) (*ListAdResponse, error) {
	var listAdRequest ListAdRequest

	campaignScene := make([]string, 0)
	campaignScene = append(campaignScene, "DAILY_SALE")
	campaignScene = append(campaignScene, "NEW_CUSTOMER_TRANSFORMATION")
	campaignScene = append(campaignScene, "LIVE_HEAT")
	campaignScene = append(campaignScene, "PLANT_GRASS")
	campaignScene = append(campaignScene, "PRODUCT_HEAT")

	if l := utf8.RuneCountInString(modifyTime); l > 0 {
		filtering := Filtering{
			CampaignScene:  campaignScene,
			MarketingGoal:  marketingGoal,
			Status:         "ALL_INCLUDE_DELETED",
			MarketingScene: "ALL",
			AdModifyTime:   modifyTime,
		}

		listAdRequest = ListAdRequest{
			AdvertiserId:     advertiserId,
			RequestAwemeInfo: 1,
			Filtering:        filtering.String(),
			Page:             page,
			PageSize:         oceanengine.PageSize200,
		}
	} else {
		filtering := FilteringNotAdModifyTime{
			CampaignScene:  campaignScene,
			MarketingGoal:  marketingGoal,
			Status:         "ALL_INCLUDE_DELETED",
			MarketingScene: "ALL",
		}

		listAdRequest = ListAdRequest{
			AdvertiserId:     advertiserId,
			RequestAwemeInfo: 1,
			Filtering:        filtering.String(),
			Page:             page,
			PageSize:         oceanengine.PageSize200,
		}
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listAdRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/v1.0/qianchuan/ad/get/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/ad/get/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/ad/get/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var listAdResponse ListAdResponse

	if err := listAdResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listAdResponse, nil
}

func GetAd(advertiserId, adId uint64, accessToken string) (*GetAdResponse, error) {
	getAdRequest := GetAdRequest{
		AdvertiserId: advertiserId,
		AdId:         adId,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(getAdRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/v1.0/qianchuan/ad/detail/get/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/ad/detail/get/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/ad/detail/get/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var getAdResponse GetAdResponse

	if err := getAdResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getAdResponse, nil
}
