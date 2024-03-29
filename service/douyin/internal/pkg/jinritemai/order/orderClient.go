package order

import (
	"douyin/internal/pkg/jinritemai"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func ListOrder(accessToken, openId, statTime, endTime, cursor string) (*ListOrderResponse, error) {
	listOrderRequest := ListOrderRequest{
		OpenId: openId,
	}
	
	listOrderBodyDataRequest := ListOrderBodyDataRequest{
		Size:      jinritemai.PageSize20,
		Cursor:    cursor,
		StartTime: statTime,
		EndTime:   endTime,
		TimeType:  "update",
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listOrderRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", jinritemai.ApplicationJson).
		SetBody(listOrderBodyDataRequest.String()).
		Post(jinritemai.BaseDomain + "/alliance/kol/orders/")

	if err != nil {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/kol/orders/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, jinritemai.NewJinritemaiError(90001, jinritemai.BaseDomain+"/alliance/kol/orders/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var listOrderResponse ListOrderResponse

	if err := listOrderResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listOrderResponse, nil
}
