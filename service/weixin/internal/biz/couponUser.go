package biz

import (
	"context"
	"weixin/internal/domain"
)

type CouponUserRepo interface {
	List(context.Context) ([]*domain.CouponUser, error)
}
