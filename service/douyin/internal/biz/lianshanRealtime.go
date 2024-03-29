package biz

import (
	"context"
	"douyin/internal/domain"
)

type LianshanRealtimeRepo interface {
	List(context.Context, string, string, string, string) ([]*domain.LianshanRealtime, error)
}
