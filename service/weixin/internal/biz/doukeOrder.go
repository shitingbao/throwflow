package biz

import (
	"context"
	v1 "weixin/api/service/douyin/v1"
)

type DoukeOrderRepo interface {
	ListUserId(context.Context) (*v1.ListUserIdDoukeOrdersReply, error)
	StatisticsByDay(context.Context, uint64, string) (*v1.StatisticsDoukeOrderByDaysReply, error)
}
