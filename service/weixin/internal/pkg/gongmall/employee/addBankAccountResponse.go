package employee

import (
	"encoding/json"
	"weixin/internal/pkg/gongmall"
)

type AddBankAccountDataResponse struct {
	BankAccountNo string `json:"bankAccountNo"`
	BankName      string `json:"bankName"`
}

type AddBankAccountResponse struct {
	gongmall.CommonResponse
	Data AddBankAccountDataResponse `json:"data"`
}

func (abar *AddBankAccountResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), abar); err != nil {
		return gongmall.NewGongmallError("CO_MSG_900100", gongmall.BaseDomain+"/api/merchant/employee/addBankAccount", "解析json失败："+err.Error(), response)
	} else {
		if !abar.Success {
			return gongmall.NewGongmallError(abar.ErrorCode, gongmall.BaseDomain+"/api/merchant/employee/addBankAccount", abar.ErrorMsg, response)
		}
	}

	return nil
}
