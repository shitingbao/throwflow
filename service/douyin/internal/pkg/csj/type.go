package csj

import (
	"douyin/internal/pkg/tool"
	"fmt"
	"strconv"
)

type CommonRequest struct {
	AppId     string `json:"app_id"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
	ReqId     string `json:"req_id"`
	Data      string `json:"data"`
}

func (cr CommonRequest) String() string {
	data, _ := tool.Marshal(cr)

	return string(data)
}

func Sign(timestamp int64, appId, appSecret, data, reqId string) string {
	paramPattern := "app_id=" + appId + "&data=" + data + "&req_id=" + reqId + "&timestamp=" + strconv.FormatInt(timestamp, 10)

	signPattern := paramPattern + appSecret
	fmt.Println(signPattern)
	return tool.GetMd5(signPattern)
}

type CommonResponse struct {
	Code int64  `json:"code"`
	Desc string `json:"desc"`
}
