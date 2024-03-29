package biz

import (
	"context"
	"douyin/internal/domain"
)

type OpenDouyinApiLogRepo interface {
	Save(context.Context, *domain.OpenDouyinApiLog) error
}
