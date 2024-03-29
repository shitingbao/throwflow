package product

import (
	"douyin/internal/pkg/oceanengine"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func ListAvailableProduct(advertiserId uint64, accessToken, marketingScene string, page uint32) (*ListAvailableProductResponse, error) {
	listAvailableProductRequest := ListAvailableProductRequest{
		AdvertiserId:   advertiserId,
		MarketingScene: marketingScene,
		Page:           page,
		PageSize:       oceanengine.PageSize100,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listAvailableProductRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/v1.0/qianchuan/product/available/get/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/product/available/get/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/product/available/get/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var listAvailableProductResponse ListAvailableProductResponse

	if err := listAvailableProductResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listAvailableProductResponse, nil
}
