package employee

import (
	"encoding/json"
	"weixin/internal/pkg/gongmall"
)

type SignContractResponse struct {
	gongmall.CommonResponse
}

func (scr *SignContractResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), scr); err != nil {
		return gongmall.NewGongmallError("CO_MSG_900100", gongmall.BaseDomain+"/api/merchant/employee/signContract", "解析json失败："+err.Error(), response)
	} else {
		if !scr.Success {
			return gongmall.NewGongmallError(scr.ErrorCode, gongmall.BaseDomain+"/api/merchant/employee/signContract", scr.ErrorMsg, response)
		}
	}

	return nil
}
