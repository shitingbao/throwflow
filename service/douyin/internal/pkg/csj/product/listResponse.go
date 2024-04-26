package product

import (
	"douyin/internal/pkg/csj"
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type ListDataProductResponse struct {
	ProductId    uint64 `json:"product_id"`     // 商品id
	Title        string `json:"title"`          // 商品名称
	IsKolProduct bool   `json:"is_kol_product"` // 是否有达⼈特殊佣⾦
	Price        uint64 `json:"price"`          // 商品价格，单位分
	InStock      bool   `json:"in_stock"`       // 有⽆库存
	FirstCid     uint64 `json:"first_cid"`      // 商品⼀级类⽬
	SecondCid    uint64 `json:"second_cid"`     // 商品⼆级类⽬
	ThirdCid     uint64 `json:"third_cid"`      // 商品三级类⽬
	Sales        uint64 `json:"sales"`          // 商品历史销量
	Cover        string `json:"cover"`          // 商品主图
	DetailUrl    string `json:"detail_url"`     // 商品链接
	ShopId       uint64 `json:"shop_id"`        // 商铺id
	ShopName     string `json:"shop_name"`      // 商铺名称
	CouponPrice  uint64 `json:"coupon_price"`   // 券后价格，单位分（0或者没传则为原价）
	CosRatio     uint64 `json:"cos_ratio"`      // 分佣⽐例，百分⽐乘以 100，⽐如1% 返回 1*100 = 100
	CosFee       uint64 `json:"cos_fee"`        // 佣⾦⾦额，单位分
	ActivityId   uint64 `json:"activity_id"`    // 获取活动商品。1: 返回超值购商品；0:返回全量商品
	Ext          string `json:"ext"`            // ⼀个加密字段，需要在转链接⼝当中回传
	PostFree     bool   `json:"post_free"`      // 是否包邮
	LimitMinSale bool   `json:"limit_min_sale"` // 是否多件起购
}

type ListDataResponse struct {
	Total    uint64                    `json:"total"` // 满⾜搜索条件的商品总数，如果超过2000，则只展⽰前2000
	Products []ListDataProductResponse `json:"products"`
}

type ListResponse struct {
	csj.CommonResponse
	Data ListDataResponse `json:"data"`
}

func (lr *ListResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lr); err != nil {
		return csj.NewCsjError(90100, douke.BaseDomain+"/product/search", "解析json失败："+err.Error(), response)
	} else {
		if lr.CommonResponse.Code != 0 {
			return csj.NewCsjError(lr.CommonResponse.Code, douke.BaseDomain+"/product/search", lr.CommonResponse.Desc, response)
		}
	}

	return nil
}
