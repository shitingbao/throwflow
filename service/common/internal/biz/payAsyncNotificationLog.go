package biz

import (
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CommonPayAsyncNotificationLogCreateError = errors.InternalServer("COMMON_PAY_ASYNC_NOTIFICATION_LOG_CREATE_ERROR", "支付异步通知日志创建失败")
)

type PayAsyncNotificationLogRepo interface {
	Save(context.Context, *domain.PayAsyncNotificationLog) (*domain.PayAsyncNotificationLog, error)
}
