package biz

import (
	"context"
	v1 "weixin/api/service/common/v1"
)

type PayRepo interface {
	Pay(context.Context, uint64, float64, string, string, string, string, string) (*v1.PayReply, error)
	PayAsyncNotification(context.Context, string) (*v1.PayAsyncNotificationReply, error)
}
