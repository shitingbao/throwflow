package command

import (
	"douyin/internal/pkg/csj"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

func Parse(appId, appSecret, command, reqId string) (*ParseResponse, error) {
	timestamp := time.Now().Unix()

	parseDataRequest := ParseDataRequest{
		Command: command,
	}

	sign := csj.Sign(timestamp, appId, appSecret, parseDataRequest.String(), reqId)

	commonRequest := csj.CommonRequest{
		AppId:     appId,
		Timestamp: timestamp,
		Sign:      sign,
		ReqId:     reqId,
		Data:      parseDataRequest.String(),
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", csj.ApplicationJson).
		SetBody(commonRequest.String()).
		Post(csj.BaseDomain + "/command_parse")

	if err != nil {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/command_parse", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, csj.NewCsjError(90101, csj.BaseDomain+"/command_parse", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var parseResponse ParseResponse

	if err := parseResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &parseResponse, nil
}
