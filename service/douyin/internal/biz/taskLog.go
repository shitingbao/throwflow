package biz

import (
	"context"
	"douyin/internal/domain"
)

type TaskLogRepo interface {
	Save(context.Context, *domain.TaskLog) error
}
