package biz

import (
	"context"
	"douyin/internal/domain"
)

type OceanengineApiLogRepo interface {
	Save(context.Context, *domain.OceanengineApiLog) error
}
