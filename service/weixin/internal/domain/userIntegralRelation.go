package domain

import (
	"context"
	"time"
)

type UserIntegralRelation struct {
	Id                 uint64
	UserId             uint64
	OrganizationId     uint64
	OrganizationUserId uint64
	Level              uint8
	CreateTime         time.Time
	UpdateTime         time.Time
}

func NewUserIntegralRelation(ctx context.Context, userId, organizationId, organizationUserId uint64, level uint8) *UserIntegralRelation {
	return &UserIntegralRelation{
		UserId:             userId,
		OrganizationId:     organizationId,
		OrganizationUserId: organizationUserId,
		Level:              level,
	}
}

func (uir *UserIntegralRelation) SetUserId(ctx context.Context, userId uint64) {
	uir.UserId = userId
}

func (uir *UserIntegralRelation) SetOrganizationId(ctx context.Context, organizationId uint64) {
	uir.OrganizationId = organizationId
}

func (uir *UserIntegralRelation) SetOrganizationUserId(ctx context.Context, organizationUserId uint64) {
	uir.OrganizationUserId = organizationUserId
}

func (uir *UserIntegralRelation) SetLevel(ctx context.Context, level uint8) {
	uir.Level = level
}

func (uir *UserIntegralRelation) SetUpdateTime(ctx context.Context) {
	uir.UpdateTime = time.Now()
}

func (uir *UserIntegralRelation) SetCreateTime(ctx context.Context) {
	uir.CreateTime = time.Now()
}
