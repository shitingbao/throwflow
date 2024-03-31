package merchant

import (
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
	"weixin/internal/pkg/gongmall"
	"weixin/internal/pkg/tool"
)

func DoSinglePayment(serviceId uint64, appkey, appSecret, name, mobile, identity, bankAccount, amount, datetime, requestId string) (*DoSinglePaymentResponse, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce := tool.GetRandCode(timestamp)

	doSinglePaymentRequest := make(map[string]string)
	doSinglePaymentRequest["appKey"] = appkey
	doSinglePaymentRequest["nonce"] = nonce
	doSinglePaymentRequest["timestamp"] = timestamp
	doSinglePaymentRequest["sign"] = gongmall.Sign("amount="+amount+"&appKey="+appkey+"&bankAccount="+bankAccount+"&dateTime="+datetime+"&identity="+identity+"&mobile="+mobile+"&name="+name+"&nonce="+nonce+"&requestId="+requestId+"&salaryType="+gongmall.SalaryTypeBank+"&serviceId="+strconv.FormatUint(serviceId, 10)+"&timestamp="+timestamp, appSecret)
	doSinglePaymentRequest["serviceId"] = strconv.FormatUint(serviceId, 10)
	doSinglePaymentRequest["requestId"] = requestId
	doSinglePaymentRequest["mobile"] = mobile
	doSinglePaymentRequest["name"] = name
	doSinglePaymentRequest["amount"] = amount
	doSinglePaymentRequest["identity"] = identity
	doSinglePaymentRequest["bankAccount"] = bankAccount
	doSinglePaymentRequest["dateTime"] = datetime
	doSinglePaymentRequest["salaryType"] = gongmall.SalaryTypeBank

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", gongmall.ApplicationForm).
		SetHeader("Authorization", gongmall.HeadAuthorization).
		SetFormData(doSinglePaymentRequest).
		Post(gongmall.BaseDomain + "/api/merchant/doSinglePayment")

	if err != nil {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/doSinglePayment", "请求失败："+err.Error(), "")
	}

	if resp.StatusCode() != 200 {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/doSinglePayment", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var doSinglePaymentResponse DoSinglePaymentResponse

	if err := doSinglePaymentResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &doSinglePaymentResponse, nil
}
