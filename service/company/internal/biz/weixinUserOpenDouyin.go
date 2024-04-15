package biz

import (
	v1 "company/api/service/weixin/v1"
	"context"
)

type WeixinUserOpenDouyinRepo interface {
	ListByClientKeyAndOpenIds(context.Context, uint64, uint64, string, string) (*v1.ListByClientKeyAndOpenIdsReply, error)
}
