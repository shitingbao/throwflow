package kol

import (
	"douyin/internal/pkg/jinritemai"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func GetBuyinId(accessToken, openId string) (*GetBuyinIdResponse, error) {
	getBuyinIdRequest := GetBuyinIdRequest{
		OpenId: openId,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(getBuyinIdRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", jinritemai.ApplicationJson).
		Post(jinritemai.BaseDomain + "/alliance/kol/buyin_id/")

	if err != nil {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/kol/buyin_id/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/kol/buyin_id/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var getBuyinIdResponse GetBuyinIdResponse

	if err := getBuyinIdResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getBuyinIdResponse, nil
}
