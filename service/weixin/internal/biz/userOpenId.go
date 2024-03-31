package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"weixin/internal/domain"
)

var (
	WeixinUserOpenidNotFound = errors.NotFound("WEIXIN_USER_OPENID_NOT_FOUND", "微信用户openid不存在")
)

type UserOpenIdRepo interface {
	Get(context.Context, uint64, string, string) (*domain.UserOpenId, error)
	List(context.Context, uint64) ([]*domain.UserOpenId, error)
	Update(context.Context, *domain.UserOpenId) (*domain.UserOpenId, error)
	Save(context.Context, *domain.UserOpenId) (*domain.UserOpenId, error)
}
