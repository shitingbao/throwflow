package domain

import (
	"context"
	"time"
)

type CompanyUserRole struct {
	UserId       uint64
	AdvertiserId uint64
	CompanyId    uint64
	RoleType     uint8
	Day          uint32
	CreateTime   time.Time
	UpdateTime   time.Time
}

type CompanyUserRoleAndUsername struct {
	UserId       uint64
	AdvertiserId uint64
	CampaignId   uint64
	CompanyId    uint64
	Username     string
}

func NewCompanyUserRole(ctx context.Context, userId, advertiserId, companyId uint64, day uint32, roleType uint8) *CompanyUserRole {
	return &CompanyUserRole{
		UserId:       userId,
		AdvertiserId: advertiserId,
		CompanyId:    companyId,
		RoleType:     roleType,
		Day:          day,
	}
}

func (cur *CompanyUserRole) SetUserId(ctx context.Context, userId uint64) {
	cur.UserId = userId
}

func (cur *CompanyUserRole) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	cur.AdvertiserId = advertiserId
}

func (cur *CompanyUserRole) SetCompanyId(ctx context.Context, companyId uint64) {
	cur.CompanyId = companyId
}

func (cur *CompanyUserRole) SetDay(ctx context.Context, day uint32) {
	cur.Day = day
}

func (cur *CompanyUserRole) SetRoleType(ctx context.Context, roleType uint8) {
	cur.RoleType = roleType
}

func (cur *CompanyUserRole) SetUpdateTime(ctx context.Context) {
	cur.UpdateTime = time.Now()
}

func (cur *CompanyUserRole) SetCreateTime(ctx context.Context) {
	cur.CreateTime = time.Now()
}
