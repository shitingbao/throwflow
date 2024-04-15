package order

import (
	"douyin/internal/pkg/csj"
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type LinkDataOrderResponse struct {
	OrderId                 string  `json:"order_id"`                   // 订单id
	AppId                   string  `json:"app_id"`                     // 应⽤id
	ProductId               string  `json:"product_id"`                 // 商品id
	ProductName             string  `json:"product_name"`               // 商品名称
	AuthorAccount           string  `json:"author_account"`             // 直播间达⼈昵称，仅直播间分销订单才会有该字段
	AdsAttribution          string  `json:"ads_attribution"`            // 结算⽅式, intersect为跨播，directIn为直推
	ProductImg              string  `json:"product_img"`                // 商品图⽚url
	TotalPayAmount          int     `json:"total_pay_amount"`           // 总付款⾦额，单位分
	PaySuccessTime          string  `json:"pay_success_time"`           // ⽀付成功时间
	RefundTime              string  `json:"refund_time"`                // 退款时间
	PayGoodsAmount          int     `json:"pay_goods_amount"`           // 预估结算⾦额，单位分。如果有⽀付优惠， pay_goods_amount会略⼤于total_pay_amount
	EstimatedCommission     float32 `json:"estimated_commission"`       // 预估佣⾦收⼊，单位分
	AdsRealCommission       float32 `json:"ads_real_commission"`        // 实际可结算金额，单位分
	SplitRate               float32 `json:"split_rate"`                 // 推⼴费率，单位万分之⼀，10代表推⼴费率为0.10%
	AfterSalesStatus        int     `json:"after_sales_status"`         // 售后状态，1-空，2-产⽣退款
	FlowPoint               string  `json:"flow_point"`                 // PAY_SUCC:⽀付完成，REFUND:退款，SETTLE:结算。此状态代表商家确定会结算佣⾦，CONFIRM: 确认收货
	ExternalInfo            string  `json:"external_info"`              // 开发者在转链时⾃⼰上传的external_info
	SettleTime              string  `json:"settle_time"`                // 结算时间
	ConfirmTime             string  `json:"confirm_time"`               // 确认收货时间
	MediaTypeName           string  `json:"media_type_name"`            // 分销类型名称：Live-直播间，ProductDetail-商品详情，ActivityMaterial-活动页
	UpdateTime              string  `json:"update_time"`                // 订单状态更新时间
	EstimatedTechServiceFee int     `json:"estimated_tech_service_fee"` // 预估技术服务费，为 pay_goods_amount*split_rate*0.1。
}

type LinkDataResponse struct {
	Cursor string                  `json:"cursor"`
	Orders []LinkDataOrderResponse `json:"orders"`
}

type ListResponse struct {
	csj.CommonResponse
	Data LinkDataResponse `json:"data"`
}

func (lr *ListResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lr); err != nil {
		return csj.NewCsjError(90100, douke.BaseDomain+"/order/search", "解析json失败："+err.Error(), response)
	} else {
		if lr.CommonResponse.Code != 0 {
			return csj.NewCsjError(lr.CommonResponse.Code, douke.BaseDomain+"/order/search", lr.CommonResponse.Desc, response)
		}
	}

	return nil
}
