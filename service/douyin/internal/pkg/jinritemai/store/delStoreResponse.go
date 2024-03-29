package store

import (
	"douyin/internal/pkg/jinritemai"
	"encoding/json"
)

type DelStoreDataResponse struct {
	ErrorCode   uint64 `json:"error_code"`
	Description string `json:"description"`
	Success     bool   `json:"success"`
}

type DelStoreResponse struct {
	jinritemai.CommonResponse
	Data DelStoreDataResponse `json:"data"`
}

func (dsr *DelStoreResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), dsr); err != nil {
		return jinritemai.NewJinritemaiError(90000, jinritemai.BaseDomain+"/alliance/kol/store/remove", "解析json失败："+err.Error(), response)
	} else {
		if dsr.CommonResponse.Extra.ErrorCode != 0 {
			return jinritemai.NewJinritemaiError(dsr.CommonResponse.Extra.ErrorCode, jinritemai.BaseDomain+"/alliance/kol/store/remove", dsr.CommonResponse.Extra.Description, response)
		}
	}

	return nil
}
