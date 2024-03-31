package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"weixin/internal/domain"
)

var (
	WeixinUserCouponCreateLogNotFound    = errors.InternalServer("WEIXIN_USER_COUPON_CREATE_LOG_NOT_FOUND", "微信用户券码新增记录不存在")
	WeixinUserCouponCreateLogUpdateError = errors.InternalServer("WEIXIN_USER_COUPON_CREATE_LOG_UPDATE_ERROR", "微信用户券码新增记录更新失败")
)

type UserCouponCreateLogRepo interface {
	Get(context.Context, string) (*domain.UserCouponCreateLog, error)
	List(context.Context, string) ([]*domain.UserCouponCreateLog, error)
	Save(context.Context, *domain.UserCouponCreateLog) (*domain.UserCouponCreateLog, error)
	Update(context.Context, *domain.UserCouponCreateLog) (*domain.UserCouponCreateLog, error)
	Delete(context.Context, *domain.UserCouponCreateLog) error
}
