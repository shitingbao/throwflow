package product

import (
	"douyin/internal/pkg/csj"
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type Score struct {
	Text  string `json:"text"`
	Score string `json:"score"`
	Level uint64 `json:"level"`
}

type ShopTotalScore struct {
	ShopScore      Score `json:"shop_score"`      //商家体验分
	ProductScore   Score `json:"product_score"`   //商品体验分
	LogisticsScore Score `json:"logistics_score"` //物流体验分
	ServiceScore   Score `json:"service_score"`   //商家服务分
}

type DailyStatistic struct {
	Date     string `json:"date"`      // 日期
	OrderNum uint64 `json:"order_num"` // 商品总订单量
	ViewNum  uint64 `json:"view_num"`  // 商品总浏览量
	KolNum   uint64 `json:"kol_num"`   // 推⼴总达⼈数
}

type DetailDataProductResponse struct {
	ProductId       uint64           `json:"product_id"`       // 商品id
	Title           string           `json:"title"`            // 商品名称
	IsKolProduct    bool             `json:"is_kol_product"`   // 是否有达⼈特殊佣⾦
	Price           uint64           `json:"price"`            // 商品价格，单位分
	CosRatio        uint64           `json:"cos_ratio"`        // 分佣⽐例，百分⽐乘以 100，⽐如1% 返回 1*100 = 100
	CosFee          uint64           `json:"cos_fee"`          // 佣⾦⾦额，单位分
	FirstCid        uint64           `json:"first_cid"`        // 商品⼀级类⽬
	SecondCid       uint64           `json:"second_cid"`       // 商品⼆级类⽬
	ThirdCid        uint64           `json:"third_cid"`        // 商品三级类⽬
	InStock         bool             `json:"in_stock"`         // 有⽆库存
	Sales           uint64           `json:"sales"`            // 商品历史销量
	Cover           string           `json:"cover"`            // 商品主图
	Imgs            []string         `json:"imgs"`             // 商品轮播图
	DetailUrl       string           `json:"detail_url"`       // 商品链接
	ShopId          uint64           `json:"shop_id"`          // 商铺id
	ShopName        string           `json:"shop_name"`        // 商铺名称
	CommentScore    float64          `json:"comment_score"`    // 商品评分
	CommentNum      uint64           `json:"comment_num"`      // 商品评价数⽬
	OrderNum        uint64           `json:"order_num"`        // 近30天商品总订单量
	ViewNum         uint64           `json:"view_num"`         // 近30天商品总浏览量
	KolNum          uint64           `json:"kol_num"`          // 近30天推⼴总达⼈数
	DailyStatistics []DailyStatistic `json:"daily_statistics"` // 近30天推⼴达⼈数、浏览量、和订单量明细
	LogisticsInfo   string           `json:"logistics_info"`   // 商品物流说明
	HasSxt          bool             `json:"has_sxt"`          // 是否具有短视频随⼼推资
	ShopTotalScore  ShopTotalScore   `json:"shop_total_score"` // 是否具有短视频随⼼推资
	CouponPrice     uint64           `json:"coupon_price"`     // 券后价格，单位分（0或者没传则为原价）
}

type DetailDataResponse struct {
	Total    uint64                      `json:"total"` // 商品总数
	Products []DetailDataProductResponse `json:"products"`
}

type DetailResponse struct {
	csj.CommonResponse
	Data DetailDataResponse `json:"data"`
}

func (dr *DetailResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), dr); err != nil {
		return csj.NewCsjError(90100, douke.BaseDomain+"/product/detail", "解析json失败："+err.Error(), response)
	} else {
		if dr.CommonResponse.Code != 0 {
			return csj.NewCsjError(dr.CommonResponse.Code, douke.BaseDomain+"/product/detail", dr.CommonResponse.Desc, response)
		}
	}

	return nil
}
