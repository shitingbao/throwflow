package product

import (
	"douyin/internal/pkg/csj"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

func List(appId, appSecret, reqId string, page, cosRatioMin uint64, firstCids, secondCids, thirdCids []uint64) (*ListResponse, error) {
	timestamp := time.Now().Unix()

	listDataRequest := ListDataRequest{
		Page:        page,
		PageSize:    csj.PageSize20,
		FirstCids:   firstCids,
		SecondCids:  secondCids,
		ThirdCids:   thirdCids,
		SearchType:  csj.SearchType3,
		OrderType:   csj.OrderType1,
		CosRatioMin: cosRatioMin,
	}

	sign := csj.Sign(timestamp, appId, appSecret, listDataRequest.String(), reqId)

	commonRequest := csj.CommonRequest{
		AppId:     appId,
		Timestamp: timestamp,
		Sign:      sign,
		ReqId:     reqId,
		Data:      listDataRequest.String(),
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", csj.ApplicationJson).
		SetBody(commonRequest.String()).
		Post(csj.BaseDomain + "/product/search")

	if err != nil {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/product/search", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/product/search", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var listResponse ListResponse

	if err := listResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listResponse, nil
}

func Link(appId, appSecret, productUrl, externalInfo, reqId string) (*LinkResponse, error) {
	timestamp := time.Now().Unix()

	linkDataRequest := LinkDataRequest{
		ProductUrl:   productUrl,
		ExternalInfo: externalInfo,
		ShareType:    []int{3},
	}

	sign := csj.Sign(timestamp, appId, appSecret, linkDataRequest.String(), reqId)

	commonRequest := csj.CommonRequest{
		AppId:     appId,
		Timestamp: timestamp,
		Sign:      sign,
		ReqId:     reqId,
		Data:      linkDataRequest.String(),
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", csj.ApplicationJson).
		SetBody(commonRequest.String()).
		Post(csj.BaseDomain + "/product/link")

	if err != nil {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/product/link", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/product/link", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var linkResponse LinkResponse

	if err := linkResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &linkResponse, nil
}

func Detail(appId, appSecret, reqId string, productIds []uint64) (*DetailResponse, error) {
	timestamp := time.Now().Unix()

	detailDataRequest := DetailDataRequest{
		ProductIds: productIds,
	}

	sign := csj.Sign(timestamp, appId, appSecret, detailDataRequest.String(), reqId)

	commonRequest := csj.CommonRequest{
		AppId:     appId,
		Timestamp: timestamp,
		Sign:      sign,
		ReqId:     reqId,
		Data:      detailDataRequest.String(),
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", csj.ApplicationJson).
		SetBody(commonRequest.String()).
		Post(csj.BaseDomain + "/product/detail")

	if err != nil {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/product/detail", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/product/detail", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}
	fmt.Println(resp.String())
	var detailResponse DetailResponse

	if err := detailResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &detailResponse, nil
}

func Category(appId, appSecret, reqId string, parentId uint64) (*CategoryResponse, error) {
	timestamp := time.Now().Unix()

	categoryDataRequest := CategoryDataRequest{
		ParentId: parentId,
	}

	sign := csj.Sign(timestamp, appId, appSecret, categoryDataRequest.String(), reqId)

	commonRequest := csj.CommonRequest{
		AppId:     appId,
		Timestamp: timestamp,
		Sign:      sign,
		ReqId:     reqId,
		Data:      categoryDataRequest.String(),
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", csj.ApplicationJson).
		SetBody(commonRequest.String()).
		Post(csj.BaseDomain + "/product/category")

	if err != nil {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/product/category", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/product/category", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}
	fmt.Println(resp.String())
	var categoryResponse CategoryResponse

	if err := categoryResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &categoryResponse, nil
}
