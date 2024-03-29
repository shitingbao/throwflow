package biz

import (
	"context"
	"douyin/internal/domain"
)

type QianchuanProductRepo interface {
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.QianchuanProduct) error
}
