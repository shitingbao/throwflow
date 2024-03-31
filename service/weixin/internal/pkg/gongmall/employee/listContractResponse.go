package employee

import (
	"encoding/json"
	"weixin/internal/pkg/gongmall"
)

type ListContractDataResponse struct {
	TemplateId          string `json:"templateId"`
	TemplateName        string `json:"templateName"`
	PlatformCompanyName string `json:"platformCompanyName"`
	ContractAddr        string `json:"contractAddr"`
}

type ListContractResponse struct {
	gongmall.CommonResponse
	Data []ListContractDataResponse `json:"data"`
}

func (lcr *ListContractResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lcr); err != nil {
		return gongmall.NewGongmallError("CO_MSG_900100", gongmall.BaseDomain+"/api/merchant/contract/getList", "解析json失败："+err.Error(), response)
	} else {
		if !lcr.Success {
			return gongmall.NewGongmallError(lcr.ErrorCode, gongmall.BaseDomain+"/api/merchant/contract/getList", lcr.ErrorMsg, response)
		}
	}

	return nil
}
