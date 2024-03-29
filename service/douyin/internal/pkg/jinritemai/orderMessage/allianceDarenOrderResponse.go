package orderMessage

import (
	"douyin/internal/pkg/douke"
	"douyin/internal/pkg/jinritemai"
	"douyin/internal/pkg/jinritemai/order"
	"encoding/json"
)

type AllianceDarenOrderDataResponse struct {
	OrderId                        string        `json:"order_id"`                   // 订单ID
	ProductId                      string        `json:"product_id"`                 // 商品id
	ProductName                    string        `json:"product_name"`               // 商品名称
	ProductImg                     string        `json:"product_img"`                // 商品图片URL
	AuthorAccount                  string        `json:"author_account"`             // 作者账号昵称(抖音/火山作者)
	AuthorOpenId                   string        `json:"author_openid"`              // 作者抖音open_id
	ShopName                       string        `json:"shop_name"`                  // 商家名称
	TotalPayAmount                 float64       `json:"total_pay_amount"`           // 订单支付金额，单位分
	CommissionRate                 float64       `json:"commission_rate"`            // 达人佣金率，此处保存为真实数据x1万之后，如真实是0.35，这里是3500
	FlowPoint                      string        `json:"flow_point"`                 // 订单状态(PAY_SUCC:支付完成 REFUND:退款 SETTLE:结算 CONFIRM: 确认收货)
	App                            string        `json:"app"`                        // App名称（抖音，火山）
	UpdateTime                     string        `json:"update_time"`                // 更新时间 [联盟侧订单更新时间]
	PaySuccessTime                 string        `json:"pay_success_time"`           // 付款时间
	SettleTime                     string        `json:"settle_time"`                // 结算时间，结算前为空字符串
	PayGoodsAmount                 int64         `json:"pay_goods_amount"`           // 预估参与结算金额
	SettledGoodsAmount             int64         `json:"settled_goods_amount"`       // 实际参与结算金额
	EstimatedCommission            int64         `json:"estimated_commission"`       // 达人预估佣金收入，单位分
	RealCommission                 int64         `json:"real_commission"`            // 达人实际佣金收入，单位分
	Extra                          string        `json:"extra"`                      // 其他
	ItemNum                        int64         `json:"item_num"`                   // 商品数目
	ShopId                         int64         `json:"shop_id"`                    // 店铺ID
	RefundTime                     string        `json:"refund_time"`                // 退款订单退款时间
	PidInfo                        order.PidInfo `json:"pid_info"`                   // 分销订单相关参数
	EstimatedTotalCommission       int64         `json:"estimated_total_commission"` // 总佣金（预估），对应百应订单明细中的总佣金
	EstimatedTechServiceFee        int64         `json:"estimated_tech_service_fee"` // 预估平台技术服务费
	PickSourceClientKey            string        `json:"pick_source_client_key"`     // 选品App client_key
	PickExtra                      string        `json:"pick_extra"`                 // 选品来源自定义参数
	AuthorShortId                  string        `json:"author_short_id"`            // 达人抖音号/火山号
	MediaType                      string        `json:"media_type"`                 // 带货体裁。shop_list：橱窗；video：视频；live：直播；others：其他(如图文、微头条、问答、西瓜长视频等)
	IsSteppedPlan                  bool          `json:"is_stepped_plan"`
	PlatformSubsidy                int64         `json:"platform_subsidy"`
	AuthorSubsidy                  int64         `json:"author_subsidy"`
	ProductActivityId              string        `json:"product_activity_id"`
	AppId                          int64         `json:"app_id"`
	SettleUserSteppedCommission    int64         `json:"settle_user_stepped_commission"`
	SettleInstSteppedCommission    int64         `json:"settle_inst_stepped_commission"`
	PaySubsidy                     int64         `json:"pay_subsidy"`
	MediaId                        int64         `json:"media_id"`
	AuthorBuyinId                  string        `json:"author_buyin_id"`
	ConfirmTime                    string        `json:"confirm_time"`
	EstimatedInstSteppedCommission int64         `json:"estimated_inst_stepped_commission"`
	EstimatedUserSteppedCommission int64         `json:"estimated_user_stepped_commission"`
}

type AllianceDarenOrderResponse struct {
	jinritemai.MessageResponse
}

func (dador *AllianceDarenOrderResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), dador); err != nil {
		return douke.NewDoukeError(90100, douke.BaseDomain+"/buyin/kolProductShare", "解析json失败："+err.Error(), response)
	}

	return nil
}
