package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"time"
	"weixin/internal/domain"
)

var (
	WeixinUserBalanceLogNotFound               = errors.NotFound("WEIXIN_USER_BALANCE_LOG_NOT_FOUND", "微信用户余额提现不存在")
	WeixinUserBalanceLogListError              = errors.InternalServer("WEIXIN_USER_BALANCE_LOG_LIST_ERROR", "微信用户余额明细列表获取失败")
	WeixinUserBalanceLogAsyncNotificationError = errors.InternalServer("WEIXIN_USER_BALANCE_LOG_ASYNC_NOTIFICATION_ERROR", "微信用户余额提现异步消息同步失败")
)

type UserBalanceLogRepo interface {
	NextId(context.Context) (uint64, error)

	GetByOutTradeNo(context.Context, string) (*domain.UserBalanceLog, error)
	List(context.Context, int, int, uint64, uint8, string, []string, []string) ([]*domain.UserBalanceLog, error)
	Count(context.Context, uint64, uint8, []string, []string) (int64, error)
	Statistics(context.Context, uint64, uint8, []string, []string) (*domain.UserBalanceLog, error)
	Save(context.Context, *domain.UserBalanceLog) (*domain.UserBalanceLog, error)
	Update(context.Context, *domain.UserBalanceLog) (*domain.UserBalanceLog, error)
	DeleteByDay(context.Context, uint8, uint32, []string) error

	SaveCacheString(context.Context, string, string, time.Duration) (bool, error)

	DeleteCache(context.Context, string) error
}
