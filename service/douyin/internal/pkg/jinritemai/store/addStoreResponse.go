package store

import (
	"douyin/internal/pkg/jinritemai"
	"encoding/json"
)

type AddStoreInfo struct {
	ProductId int64  `json:"product_id"`
	ErrorCode int64  `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type AddStoreDataResponse struct {
	ErrorCode   uint64         `json:"error_code"`
	Description string         `json:"description"`
	Results     []AddStoreInfo `json:"results"`
}

type AddStoreResponse struct {
	jinritemai.CommonResponse
	Data AddStoreDataResponse `json:"data"`
}

func (asr *AddStoreResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), asr); err != nil {
		return jinritemai.NewJinritemaiError(90000, jinritemai.BaseDomain+"/alliance/store/add/", "解析json失败："+err.Error(), response)
	} else {
		if asr.CommonResponse.Extra.ErrorCode != 0 {
			return jinritemai.NewJinritemaiError(asr.CommonResponse.Extra.ErrorCode, jinritemai.BaseDomain+"/alliance/store/add/", asr.CommonResponse.Extra.Description, response)
		}
	}

	return nil
}
