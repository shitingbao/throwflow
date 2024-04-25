package biz

import (
	"context"
	v1 "douyin/api/service/weixin/v1"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	DouyinWeixinUserCommissionOrderCreateError     = errors.NotFound("DOUYIN_WEIXIN_USER_COMMISSION_ORDER_CREATE_ERROR", "微信用户电商分佣创建失败")
	DouyinWeixinUserCommissionCostOrderCreateError = errors.NotFound("DOUYIN_WEIXIN_USER_COMMISSION_COST_ORDER_CREATE_ERROR", "微信用户成本购分佣创建失败")
)

type WeixinUserCommissionRepo interface {
	CreateOrder(context.Context, float64, float64, string, string, string, string, string) (*v1.CreateOrderUserCommissionsReply, error)
	CreateCostOrder(context.Context, uint64, float64, float64, string, string, string, string) (*v1.CreateCostOrderUserCommissionsReply, error)
}
