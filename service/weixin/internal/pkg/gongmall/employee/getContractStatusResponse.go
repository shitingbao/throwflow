package employee

import (
	"encoding/json"
	"weixin/internal/pkg/gongmall"
)

type GetContractStatusDataResponse struct {
	ProcessStatus uint8  `json:"processStatus"`
	FileUrl       string `json:"fileUrl"`
}

type GetContractStatusResponse struct {
	gongmall.CommonResponse
	Data GetContractStatusDataResponse `json:"data"`
}

func (gcsr *GetContractStatusResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gcsr); err != nil {
		return gongmall.NewGongmallError("CO_MSG_900100", gongmall.BaseDomain+"/api/merchant/employee/getContractStatus", "解析json失败："+err.Error(), response)
	} else {
		if !gcsr.Success {
			return gongmall.NewGongmallError(gcsr.ErrorCode, gongmall.BaseDomain+"/api/merchant/employee/getContractStatus", gcsr.ErrorMsg, response)
		}
	}

	return nil
}
