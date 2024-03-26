package biz

import (
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CommonShortCodeCreateError = errors.InternalServer("COMMON_SHORT_CODE_CREATE_ERROR", "短码创建失败")
)

type ShortCodeLogRepo interface {
	Get(context.Context, string) (*domain.ShortCodeLog, error)
	Save(context.Context, *domain.ShortCodeLog) (*domain.ShortCodeLog, error)
}
