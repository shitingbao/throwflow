package store

import (
	"douyin/internal/pkg/jinritemai"
	"encoding/json"
)

type Store struct {
	ProductId         int64   `json:"product_id"`          // 商品ID
	PromotionId       int64   `json:"promotion_id"`        // 推广ID
	Title             string  `json:"title"`               // 商品标题
	Cover             string  `json:"cover"`               // 商品主图
	PromotionType     int64   `json:"promotion_type"`      // 推广类型，0 非团长，1 团长
	Price             int64   `json:"price"`               // 商品售价（单位为分）
	CosType           int64   `json:"cos_type"`            // 佣金类型，0 未定义（异常） 1 专属团长 2 普通佣金 3 定向佣金，4 提报活动 5 招募佣金
	CosRatio          float64 `json:"cos_ratio"`           // 佣金率（10表示佣金率为10%）
	ColonelActivityId int64   `json:"colonel_activity_id"` // 团长活动id
	HideStatus        bool    `json:"hide_status"`         // 隐藏状态（true代表隐藏，false代表显示）
}

type ListOrderDataResponse struct {
	ErrorCode   uint64  `json:"error_code"`
	Description string  `json:"description"`
	Total       int64   `json:"total"`
	Results     []Store `json:"results"`
}

type ListStoreResponse struct {
	jinritemai.CommonResponse
	Data ListOrderDataResponse `json:"data"`
}

func (lsr *ListStoreResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lsr); err != nil {
		return jinritemai.NewJinritemaiError(90000, jinritemai.BaseDomain+"/alliance/kol/store/list", "解析json失败："+err.Error(), response)
	} else {
		if lsr.CommonResponse.Extra.ErrorCode != 0 {
			return jinritemai.NewJinritemaiError(lsr.CommonResponse.Extra.ErrorCode, jinritemai.BaseDomain+"/alliance/kol/store/list", lsr.CommonResponse.Extra.Description, response)
		}
	}

	return nil
}
