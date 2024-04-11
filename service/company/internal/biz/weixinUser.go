package biz

import (
	v1 "company/api/service/weixin/v1"
	"context"
)

type WeixinUserRepo interface {
	ListByIds(context.Context, string, string, string) (*v1.ListByIdsReply, error)
}
