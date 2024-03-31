package biz

import (
	"context"
	"weixin/internal/domain"
)

type AwemesAdvertiserWeixinAuthRepo interface {
	Upsert(context.Context, *domain.AwemesAdvertiserWeixinAuth) error
}
