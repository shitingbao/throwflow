package biz

import (
	"context"
	"weixin/internal/domain"
)

type TuUserRepo interface {
	List(context.Context) ([]*domain.TuUser, error)
}
