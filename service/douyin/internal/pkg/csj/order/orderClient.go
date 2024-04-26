package order

import (
	"douyin/internal/pkg/csj"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

func List(startTime, endTime int, appId, appSecret, cursor, reqId string) (*ListResponse, error) {
	timestamp := time.Now().Unix()

	listDataRequest := ListDataRequest{
		Size:      csj.PageSize50,
		Cursor:    cursor,
		StartTime: startTime,
		EndTime:   endTime,
		OrderType: csj.ListOrderOrderType1,
		TimeType:  csj.ListOrderTimeTypeUpdate,
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
		Post(csj.BaseDomain + "/order/search")

	if err != nil {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/order/search", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/order/search", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}
	
	var listResponse ListResponse

	if err := listResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listResponse, nil
}
