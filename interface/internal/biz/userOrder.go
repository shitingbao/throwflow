package biz

import (
	"context"
	v1 "interface/api/service/weixin/v1"
)

type UserOrderRepo interface {
	Create(context.Context, uint64, uint64, uint64, float64, string) (*v1.CreateUserOrdersReply, error)
	Upgrade(context.Context, uint64, uint64, float64, string) (*v1.UpgradeUserOrdersReply, error)
}
