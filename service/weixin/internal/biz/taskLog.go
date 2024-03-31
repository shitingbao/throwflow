package biz

import (
	"context"
	"weixin/internal/domain"
)

type TaskLogRepo interface {
	Save(context.Context, *domain.TaskLog) error
}
