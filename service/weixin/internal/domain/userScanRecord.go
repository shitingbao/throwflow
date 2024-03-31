package domain

import (
	"context"
	"time"
)

type UserScanRecord struct {
	Id                 uint64
	UserId             uint64
	OrganizationId     uint64
	OrganizationUserId uint64
	AwemeUserId        uint64
	CreateTime         time.Time
	UpdateTime         time.Time
}

func NewUserScanRecord(ctx context.Context, userId, organizationId, organizationUserId, awemeUserId uint64) *UserScanRecord {
	return &UserScanRecord{
		UserId:             userId,
		OrganizationId:     organizationId,
		OrganizationUserId: organizationUserId,
		AwemeUserId:        awemeUserId,
	}
}

func (usr *UserScanRecord) SetUserId(ctx context.Context, userId uint64) {
	usr.UserId = userId
}

func (usr *UserScanRecord) SetOrganizationId(ctx context.Context, organizationId uint64) {
	usr.OrganizationId = organizationId
}

func (usr *UserScanRecord) SetOrganizationUserId(ctx context.Context, organizationUserId uint64) {
	usr.OrganizationUserId = organizationUserId
}

func (usr *UserScanRecord) SetAwemeUserId(ctx context.Context, awemeUserId uint64) {
	usr.AwemeUserId = awemeUserId
}

func (usr *UserScanRecord) SetUpdateTime(ctx context.Context) {
	usr.UpdateTime = time.Now()
}

func (usr *UserScanRecord) SetCreateTime(ctx context.Context) {
	usr.CreateTime = time.Now()
}
