package account

import (
	"douyin/internal/pkg/oceanengine"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"
)

func ListAdvertiser(appId, appSecret, accessToken string) (*ListAdvertiserResponse, error) {
	listAdvertiserRequest := ListAdvertiserRequest{
		AccessToken: accessToken,
		AppId:       appId,
		Secret:      appSecret,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listAdvertiserRequest.String()).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/oauth2/advertiser/get/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/oauth2/advertiser/get/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/oauth2/advertiser/get/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var ListAdvertiserResponse ListAdvertiserResponse

	if err := ListAdvertiserResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &ListAdvertiserResponse, nil
}

func ListShopAdvertiser(advertiserId uint64, accessToken string, page uint32) (*ListShopAdvertiserResponse, error) {
	listShopAdvertiserRequest := ListShopAdvertiserRequest{
		ShopId:   advertiserId,
		Page:     page,
		PageSize: oceanengine.PageSize100,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listShopAdvertiserRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/v1.0/qianchuan/shop/advertiser/list/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/shop/advertiser/list/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/shop/advertiser/list/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var listShopAdvertiserResponse ListShopAdvertiserResponse

	if err := listShopAdvertiserResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listShopAdvertiserResponse, nil
}

func ListAgentAdvertiser(advertiserId uint64, accessToken string, page uint32) (*ListAgentAdvertiserResponse, error) {
	listAgentAdvertiserRequest := ListAgentAdvertiserRequest{
		AdvertiserId: advertiserId,
		Page:         page,
		PageSize:     oceanengine.PageSize1000,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listAgentAdvertiserRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/2/agent/advertiser/select/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/2/agent/advertiser/select/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/2/agent/advertiser/select/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var listAgentAdvertiserResponse ListAgentAdvertiserResponse

	if err := listAgentAdvertiserResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listAgentAdvertiserResponse, nil
}

func GetAdvertiserPublicInfo(advertiserIds []string, accessToken string) (*GetAdvertiserPublicInfoResponse, error) {
	getAdvertiserPublicInfoRequest := GetAdvertiserPublicInfoRequest{
		AdvertiserIds: "[" + strings.Join(advertiserIds, ",") + "]",
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(getAdvertiserPublicInfoRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/2/advertiser/public_info/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/2/advertiser/public_info/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/2/advertiser/public_info/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var getAdvertiserPublicInfoResponse GetAdvertiserPublicInfoResponse

	if err := getAdvertiserPublicInfoResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getAdvertiserPublicInfoResponse, nil
}

func ListAweme(advertiserId uint64, accessToken string, page uint32) (*ListAwemeResponse, error) {
	listAgentAdvertiserRequest := ListAwemeRequest{
		AdvertiserId: advertiserId,
		Page:         page,
		PageSize:     oceanengine.PageSize100,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listAgentAdvertiserRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/v1.0/qianchuan/aweme/authorized/get/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/aweme/authorized/get/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/aweme/authorized/get/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}
	
	var listAwemeResponse ListAwemeResponse

	if err := listAwemeResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listAwemeResponse, nil
}
