package js

import (
	"douyin/internal/pkg/openDouyin"
	"encoding/json"
)

type GetTicketDataResponse struct {
	ExpiresIn   uint64 `json:"expires_in"`
	Ticket      string `json:"ticket"`
	Description string `json:"description"`
	ErrorCode   uint64 `json:"error_code"`
}

type GetTicketResponse struct {
	openDouyin.CommonResponse
	Data GetTicketDataResponse `json:"data"`
}

func (gtr *GetTicketResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gtr); err != nil {
		return openDouyin.NewOpenDouyinError(90000, openDouyin.BaseDomain+"/js/getticket/", "解析json失败："+err.Error(), response)
	} else {
		if gtr.CommonResponse.Extra.ErrorCode != 0 {
			return openDouyin.NewOpenDouyinError(gtr.CommonResponse.Extra.ErrorCode, openDouyin.BaseDomain+"/js/getticket/", gtr.CommonResponse.Extra.Description, response)
		}
	}

	return nil
}
