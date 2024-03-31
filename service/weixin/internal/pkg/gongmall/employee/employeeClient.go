package employee

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
	"weixin/internal/pkg/gongmall"
	"weixin/internal/pkg/tool"
)

func ListContract(serviceId uint64, appkey, appSecret string) (*ListContractResponse, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce := tool.GetRandCode(timestamp)

	listContractRequest := make(map[string]string)
	listContractRequest["appKey"] = appkey
	listContractRequest["nonce"] = nonce
	listContractRequest["timestamp"] = timestamp
	listContractRequest["sign"] = gongmall.Sign("appKey="+appkey+"&nonce="+nonce+"&serviceId="+strconv.FormatUint(serviceId, 10)+"&timestamp="+timestamp, appSecret)
	listContractRequest["serviceId"] = strconv.FormatUint(serviceId, 10)

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", gongmall.ApplicationForm).
		SetFormData(listContractRequest).
		Post(gongmall.BaseDomain + "/api/merchant/contract/getList")

	if err != nil {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/contract/getList", "请求失败："+err.Error(), "")
	}

	if resp.StatusCode() != 200 {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/contract/getList", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var listContractResponse ListContractResponse

	if err := listContractResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listContractResponse, nil
}

func SyncInfo(serviceId, templateId uint64, appkey, appSecret, name, mobile, certificateType, identity string) (*SyncInfoResponse, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce := tool.GetRandCode(timestamp)

	sserviceId := strconv.FormatUint(serviceId, 10)
	stemplateId := strconv.FormatUint(templateId, 10)

	syncInfoRequest := make(map[string]string)
	syncInfoRequest["appKey"] = appkey
	syncInfoRequest["nonce"] = nonce
	syncInfoRequest["timestamp"] = timestamp
	syncInfoRequest["sign"] = gongmall.Sign("appKey="+appkey+"&certificateType="+certificateType+"&identity="+identity+"&mobile="+mobile+"&name="+name+"&nonce="+nonce+"&serviceId="+sserviceId+"&templateId="+stemplateId+"&timestamp="+timestamp, appSecret)
	syncInfoRequest["serviceId"] = sserviceId
	syncInfoRequest["name"] = name
	syncInfoRequest["mobile"] = mobile
	syncInfoRequest["certificateType"] = certificateType
	syncInfoRequest["identity"] = identity
	syncInfoRequest["templateId"] = stemplateId

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", gongmall.ApplicationForm).
		SetHeader("Authorization", gongmall.HeadAuthorization).
		SetFormData(syncInfoRequest).
		Post(gongmall.BaseDomain + "/api/merchant/employee/syncInfo")

	if err != nil {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/employee/syncInfo", "请求失败："+err.Error(), "")
	}

	if resp.StatusCode() != 200 {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/employee/syncInfo", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var syncInfoResponse SyncInfoResponse
	fmt.Println("####################################")
	fmt.Println(resp.String())
	fmt.Println("####################################")
	if err := syncInfoResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &syncInfoResponse, nil
}

func SignContract(serviceId, contractId uint64, appkey, appSecret, mobile, captcha string) (*SyncInfoResponse, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce := tool.GetRandCode(timestamp)

	sserviceId := strconv.FormatUint(serviceId, 10)
	scontractId := strconv.FormatUint(contractId, 10)

	signContractRequest := make(map[string]string)
	signContractRequest["appKey"] = appkey
	signContractRequest["nonce"] = nonce
	signContractRequest["timestamp"] = timestamp
	signContractRequest["sign"] = gongmall.Sign("appKey="+appkey+"&captcha="+captcha+"&contractId="+scontractId+"&mobile="+mobile+"&nonce="+nonce+"&serviceId="+sserviceId+"&timestamp="+timestamp, appSecret)
	signContractRequest["serviceId"] = sserviceId
	signContractRequest["mobile"] = mobile
	signContractRequest["contractId"] = scontractId
	signContractRequest["captcha"] = captcha

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", gongmall.ApplicationForm).
		SetHeader("Authorization", gongmall.HeadAuthorization).
		SetFormData(signContractRequest).
		Post(gongmall.BaseDomain + "/api/merchant/employee/signContract")

	if err != nil {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/employee/signContract", "请求失败："+err.Error(), "")
	}

	if resp.StatusCode() != 200 {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/employee/signContract", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var syncInfoResponse SyncInfoResponse

	if err := syncInfoResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &syncInfoResponse, nil
}

func GetContractStatus(serviceId, templateId uint64, appkey, appSecret, identity string) (*GetContractStatusResponse, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce := tool.GetRandCode(timestamp)

	sserviceId := strconv.FormatUint(serviceId, 10)
	stemplateId := strconv.FormatUint(templateId, 10)

	signContractRequest := make(map[string]string)
	signContractRequest["appKey"] = appkey
	signContractRequest["nonce"] = nonce
	signContractRequest["timestamp"] = timestamp
	signContractRequest["sign"] = gongmall.Sign("appKey="+appkey+"&identity="+identity+"&nonce="+nonce+"&serviceId="+sserviceId+"&templateId="+stemplateId+"&timestamp="+timestamp, appSecret)
	signContractRequest["serviceId"] = sserviceId
	signContractRequest["identity"] = identity
	signContractRequest["templateId"] = stemplateId

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", gongmall.ApplicationForm).
		SetHeader("Authorization", gongmall.HeadAuthorization).
		SetFormData(signContractRequest).
		Post(gongmall.BaseDomain + "/api/merchant/employee/getContractStatus")

	if err != nil {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/employee/getContractStatus", "请求失败："+err.Error(), "")
	}

	if resp.StatusCode() != 200 {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/employee/getContractStatus", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var getContractStatusResponse GetContractStatusResponse

	if err := getContractStatusResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getContractStatusResponse, nil
}

func AddBankAccount(appkey, appSecret, name, identity, bankAccountNo string) (*AddBankAccountResponse, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce := tool.GetRandCode(timestamp)

	signContractRequest := make(map[string]string)
	signContractRequest["appKey"] = appkey
	signContractRequest["nonce"] = nonce
	signContractRequest["timestamp"] = timestamp
	signContractRequest["sign"] = gongmall.Sign("appKey="+appkey+"&bankAccountNo="+bankAccountNo+"&identity="+identity+"&name="+name+"&nonce="+nonce+"&timestamp="+timestamp, appSecret)
	signContractRequest["name"] = name
	signContractRequest["identity"] = identity
	signContractRequest["bankAccountNo"] = bankAccountNo

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", gongmall.ApplicationForm).
		SetHeader("Authorization", gongmall.HeadAuthorization).
		SetFormData(signContractRequest).
		Post(gongmall.BaseDomain + "/api/merchant/employee/addBankAccount")

	if err != nil {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/employee/addBankAccount", "请求失败："+err.Error(), "")
	}

	if resp.StatusCode() != 200 {
		return nil, gongmall.NewGongmallError("CO_MSG_900101", gongmall.BaseDomain+"/api/merchant/employee/addBankAccount", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var addBankAccountResponse AddBankAccountResponse

	if err := addBankAccountResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &addBankAccountResponse, nil
}
