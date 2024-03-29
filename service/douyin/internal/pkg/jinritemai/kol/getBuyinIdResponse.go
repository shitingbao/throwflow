package kol

import (
	"douyin/internal/pkg/jinritemai"
	"encoding/json"
)

type GetBuyinIdDataResponse struct {
	ErrorCode   uint64 `json:"error_code"`
	Description string `json:"description"`
	BuyinId     string `json:"buyin_id"`
}

type GetBuyinIdResponse struct {
	jinritemai.CommonResponse
	Data GetBuyinIdDataResponse `json:"data"`
}

func (gbir *GetBuyinIdResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gbir); err != nil {
		return jinritemai.NewJinritemaiError(90000, jinritemai.BaseDomain+"/alliance/kol/buyin_id/", "解析json失败："+err.Error(), response)
	} else {
		if gbir.CommonResponse.Extra.ErrorCode != 0 {
			return jinritemai.NewJinritemaiError(gbir.CommonResponse.Extra.ErrorCode, jinritemai.BaseDomain+"/alliance/kol/buyin_id/", gbir.CommonResponse.Extra.Description, response)
		}
	}

	return nil
}
