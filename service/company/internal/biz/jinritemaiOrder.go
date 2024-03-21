package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyJinritemaiOrderListError = errors.InternalServer("COMPANY_JINRITEMAI_ORDER_LIST_ERROR", "达人订单列表获取失败")
)

type JinritemaiOrderRepo interface {
	ListCommissionRate(context.Context, string) (*v1.ListCommissionRateJinritemaiOrdersReply, error)
	GetIsTop(context.Context, uint64) (*v1.GetIsTopJinritemaiOrdersReply, error)
}
