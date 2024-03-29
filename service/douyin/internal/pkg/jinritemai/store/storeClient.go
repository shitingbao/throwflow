package store

import (
	"douyin/internal/pkg/jinritemai"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func ListStore(accessToken, openId string, page int64) (*ListStoreResponse, error) {
	listStoreRequest := ListStoreRequest{
		OpenId: openId,
	}

	listStoreBodyDataRequest := ListStoreBodyDataRequest{
		PageSize: jinritemai.PageSize20,
		Page:     page,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listStoreRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", jinritemai.ApplicationJson).
		SetBody(listStoreBodyDataRequest.String()).
		Post(jinritemai.BaseDomain + "/alliance/kol/store/list")

	if err != nil {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/kol/store/list", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/kol/store/list", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var listStoreResponse ListStoreResponse

	if err := listStoreResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listStoreResponse, nil
}

func AddStore(accessToken, openId, extra string, products []*AddStoreBodyDataProductsRequest) (*AddStoreResponse, error) {
	addStoreRequest := AddStoreRequest{
		OpenId: openId,
	}

	addStoreBodyDataRequest := AddStoreBodyDataRequest{
		Products:       products,
		NeedHide:       false,
		PickExtra:      extra,
		KeepPicksource: true,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(addStoreRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", jinritemai.ApplicationJson).
		SetBody(addStoreBodyDataRequest.String()).
		Post(jinritemai.BaseDomain + "/alliance/store/add/")

	if err != nil {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/store/add/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/store/add/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var addStoreResponse AddStoreResponse

	if err := addStoreResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &addStoreResponse, nil
}

func DelStore(accessToken, openId string, products []*DelStoreBodyDataProductsRequest) (*DelStoreResponse, error) {
	delStoreRequest := DelStoreRequest{
		OpenId: openId,
	}

	delStoreBodyDataRequest := DelStoreBodyDataRequest{
		Products: products,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(delStoreRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", jinritemai.ApplicationJson).
		SetBody(delStoreBodyDataRequest.String()).
		Post(jinritemai.BaseDomain + "/alliance/kol/store/remove")

	if err != nil {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/kol/store/remove", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/kol/store/remove", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var delStoreResponse DelStoreResponse

	if err := delStoreResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &delStoreResponse, nil
}
