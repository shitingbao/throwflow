package biz

import (
	"context"
	"douyin/internal/domain"
)

type JinritemaiApiLogRepo interface {
	Save(context.Context, *domain.JinritemaiApiLog) error
}
