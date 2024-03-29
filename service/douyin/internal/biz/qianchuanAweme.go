package biz

import (
	"context"
	"douyin/internal/domain"
)

type QianchuanAwemeRepo interface {
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.QianchuanAweme) error
}
