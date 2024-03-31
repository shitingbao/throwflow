package merchant

import (
	"encoding/json"
	"weixin/internal/pkg/gongmall"
)

type DoSinglePaymentDataResponse struct {
	RequestId   string `json:"requestId"`
	AppmentTime string `json:"appmentTime"`
}

type DoSinglePaymentResponse struct {
	gongmall.CommonResponse
	Data DoSinglePaymentDataResponse `json:"data"`
}

func (dspr *DoSinglePaymentResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), dspr); err != nil {
		return gongmall.NewGongmallError("CO_MSG_900100", gongmall.BaseDomain+"/api/merchant/doSinglePayment", "解析json失败："+err.Error(), response)
	} else {
		if !dspr.Success {
			return gongmall.NewGongmallError(dspr.ErrorCode, gongmall.BaseDomain+"/api/merchant/doSinglePayment", dspr.ErrorMsg, response)
		}
	}

	return nil
}
