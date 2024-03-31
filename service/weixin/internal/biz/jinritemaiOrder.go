package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/douyin/v1"
)

var (
	WeixinJinritemaiOrderGetError = errors.InternalServer("WEIXIN_JINRITEMAI_ORDER_GET_ERROR", "达人订单数据获取失败")
)

type JinritemaiOrderRepo interface {
	ListByPickExtra(context.Context) (*v1.ListJinritemaiOrderByPickExtrasReply, error)
	Statistics(context.Context, uint64, string, string) (*v1.StatisticsJinritemaiOrdersReply, error)
	StatisticsByDay(context.Context, string, string, string) (*v1.StatisticsJinritemaiOrderByDaysReply, error)
	StatisticsByPayTimeDay(context.Context, string, string, string, string) (*v1.StatisticsJinritemaiOrderByPayTimeAndDaysReply, error)
	StatisticsByClientKeyAndOpenId(context.Context, string, string, string, string) (*v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsReply, error)
}
