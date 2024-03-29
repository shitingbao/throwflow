package orderMessage

import (
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type DoudianAllianceDistributorOrderPidInfoResponse struct {
	ExternalInfo  string `json:"external_info"`   // 自定义参数
	MediaTypeName string `json:"media_type_name"` // 分销类型，Live：直播
	Pid           string `json:"pid"`             // PID
}

type DoudianAllianceDistributorOrderProductTagsResponse struct {
	HasSubsidyTag     bool `json:"has_subsidy_tag"`     // 是否超值购商品
	HasSupermarketTag bool `json:"has_supermarket_tag"` // 是否抖音超市(次日达)商品
}

type DoudianAllianceDistributorOrderDataResponse struct {
	AuthorAccount          string                                             `json:"author_account"`           // 作者账号昵称(抖音/火山作者)，仅直播间分销订单显示
	MediaType              string                                             `json:"media_type"`               // 下单体裁：shop_list 橱窗，video 视频，live 直播，others 其他
	AdsActivityId          int64                                              `json:"ads_activity_id"`          // 分销活动Id，1000-超值购 1001-秒杀（活动页推广）
	ProductImg             string                                             `json:"product_img"`              // 商品图片URL
	UpdateTime             string                                             `json:"update_time"`              // 更新时间 [联盟侧订单更新时间]
	PaySuccessTime         string                                             `json:"pay_success_time"`         // 付款时间
	AdsRealCommission      int64                                              `json:"ads_real_commission"`      // 渠道实际推广费
	AdsEstimatedCommission int64                                              `json:"ads_estimated_commission"` // 渠道预估推广费
	ProductId              string                                             `json:"product_id"`               // 商品id
	TotalPayAmount         int64                                              `json:"total_pay_amount"`         // 订单支付金额，单位分
	FlowPoint              string                                             `json:"flow_point"`               // 订单状态(PAY_SUCC:支付完成 REFUND:退款 SETTLE:结算, COMFIRM: 确认收货)
	SettleTime             string                                             `json:"settle_time"`              // 结算时间，结算前为空字符串
	SettledGoodsAmount     int64                                              `json:"settled_goods_amount"`     // 实际参与结算金额
	PidInfo                DoudianAllianceDistributorOrderPidInfoResponse     `json:"pid_info"`                 // PID的信息
	ItemNum                int64                                              `json:"item_num"`                 // 商品数目
	AuthorBuyinId          string                                             `json:"author_buyin_id"`          // 达人百应ID（仅直播间分销订单显示）
	ShopId                 int64                                              `json:"shop_id"`                  // 店铺ID
	PayGoodsAmount         int64                                              `json:"pay_goods_amount"`         // 预估参与结算金额
	AdsDistributorId       int64                                              `json:"ads_distributor_id"`       // 抖客百应ID
	AdsPromotionTate       int64                                              `json:"ads_promotion_rate"`       // 渠道推广费率
	OrderId                string                                             `json:"order_id"`                 // 订单号
	ProductName            string                                             `json:"product_name"`             // 商品名称
	DistributionType       string                                             `json:"distribution_type"`        // 分销类型，直播：Live
	AuthorUid              int64                                              `json:"author_uid"`               // 作者uid（仅直播间分销订单显示）
	ShopName               string                                             `json:"shop_name"`                // 商家名称
	ProductActivityId      string                                             `json:"product_activity_id"`      // 商品参与的活动id，0: 未参加活动 1: 超值购（活动页单品推广）
	MaterialId             string                                             `json:"material_id"`              // 活动页物料ID
	RefundTime             string                                             `json:"refund_time"`              // 退款时间
	ConfirmTime            string                                             `json:"confirm_time"`             // 收货时间
	ProductTags            DoudianAllianceDistributorOrderProductTagsResponse `json:"product_tags"`             // 商品分销订单-商品标签信息
	BuyerAppId             string                                             `json:"buyer_app_id"`             // 购买端appId
	DistributorRightType   string                                             `json:"distributor_right_type"`   // 权益类型，下单用券类型，空或者Unknown为无权益，ColonelExclusiveCoupon(抖客团长券)
}

type DoudianAllianceDistributorOrderResponse struct {
	douke.MessageResponse
	Data string `json:"data"`
}

func (dador *DoudianAllianceDistributorOrderResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), dador); err != nil {
		return douke.NewDoukeError(90100, douke.BaseDomain+"/buyin/kolProductShare", "解析json失败："+err.Error(), response)
	}

	return nil
}
