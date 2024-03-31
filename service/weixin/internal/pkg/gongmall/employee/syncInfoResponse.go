package employee

import (
	"encoding/json"
	"weixin/internal/pkg/gongmall"
)

type SyncInfoDataResponse struct {
	ContractId    string `json:"contractId"`
	ProcessStatus uint8  `json:"processStatus"`
	BankName      string `json:"bankName"`
}

type SyncInfoResponse struct {
	gongmall.CommonResponse
	Data SyncInfoDataResponse `json:"data"`
}

func (sir *SyncInfoResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), sir); err != nil {
		return gongmall.NewGongmallError("CO_MSG_900100", gongmall.BaseDomain+"/api/merchant/employee/syncInfo", "解析json失败："+err.Error(), response)
	} else {
		if !sir.Success {
			return gongmall.NewGongmallError(sir.ErrorCode, gongmall.BaseDomain+"/api/merchant/employee/syncInfo", sir.ErrorMsg, response)
		}
	}

	return nil
}
