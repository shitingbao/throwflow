package domain

import (
	"context"
	"time"
)

type UserCouponCreateLog struct {
	Id             uint64
	UserId         uint64
	OrganizationId uint64
	Level          uint8
	Num            uint32
	IsHandle       uint8
	CreateTime     time.Time
	UpdateTime     time.Time
}

func NewUserCouponCreateLog(ctx context.Context, userId, organizationId uint64, num uint32, level, isHandle uint8) *UserCouponCreateLog {
	return &UserCouponCreateLog{
		UserId:         userId,
		OrganizationId: organizationId,
		Level:          level,
		Num:            num,
		IsHandle:       isHandle,
	}
}

func (uccl *UserCouponCreateLog) SetUserId(ctx context.Context, userId uint64) {
	uccl.UserId = userId
}

func (uccl *UserCouponCreateLog) SetOrganizationId(ctx context.Context, organizationId uint64) {
	uccl.OrganizationId = organizationId
}

func (uccl *UserCouponCreateLog) SetLevel(ctx context.Context, level uint8) {
	uccl.Level = level
}

func (uccl *UserCouponCreateLog) SetNum(ctx context.Context, num uint32) {
	uccl.Num = num
}

func (uccl *UserCouponCreateLog) SetIsHandle(ctx context.Context, isHandle uint8) {
	uccl.IsHandle = isHandle
}

func (uccl *UserCouponCreateLog) SetUpdateTime(ctx context.Context) {
	uccl.UpdateTime = time.Now()
}

func (uccl *UserCouponCreateLog) SetCreateTime(ctx context.Context) {
	uccl.CreateTime = time.Now()
}
