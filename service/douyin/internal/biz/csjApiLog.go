package biz

import (
	"context"
	"douyin/internal/domain"
)

type CsjApiLogRepo interface {
	Save(context.Context, *domain.CsjApiLog) error
}
