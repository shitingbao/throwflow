package biz

import (
	"context"
	"weixin/internal/domain"
)

type UserIntegralRelationRepo interface {
	Get(context.Context, uint64, uint64) (*domain.UserIntegralRelation, error)
	GetByUserId(context.Context, uint64, uint64, uint64) (*domain.UserIntegralRelation, error)
	GetSuperior(context.Context, uint64, uint8, []*domain.UserIntegralRelation) *domain.UserIntegralRelation
	GetChildNum(context.Context, uint64, *uint64, []*domain.UserIntegralRelation)
	List(context.Context, uint64) ([]*domain.UserIntegralRelation, error)
	ListChildId(context.Context, uint64, *[]uint64, []*domain.UserIntegralRelation)
	Save(context.Context, *domain.UserIntegralRelation) (*domain.UserIntegralRelation, error)
	Update(context.Context, *domain.UserIntegralRelation) (*domain.UserIntegralRelation, error)
}
