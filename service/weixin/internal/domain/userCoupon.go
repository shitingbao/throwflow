package domain

import (
	"context"
	"time"
	"weixin/internal/pkg/tool"
)

type UserCoupon struct {
	Id               uint64
	UserId           uint64
	CouponCode       string
	Level            uint8
	Phone            string
	UserCouponStatus uint8
	OrganizationId   uint64
	Content          string
	CreateTime       time.Time
	UpdateTime       time.Time
}

type UserCouponList struct {
	PageNum   uint64
	PageSize  uint64
	Total     uint64
	TotalUsed uint64
	List      []*UserCoupon
}

func NewUserCoupon(ctx context.Context, userId, organizationId uint64, level, userCouponStatus uint8, couponCode, phone string) *UserCoupon {
	return &UserCoupon{
		UserId:           userId,
		CouponCode:       couponCode,
		Level:            level,
		Phone:            phone,
		UserCouponStatus: userCouponStatus,
		OrganizationId:   organizationId,
	}
}

func (uc *UserCoupon) SetUserId(ctx context.Context, userId uint64) {
	uc.UserId = userId
}

func (uc *UserCoupon) SetCouponCode(ctx context.Context, couponCode string) {
	uc.CouponCode = couponCode
}

func (uc *UserCoupon) SetLevel(ctx context.Context, level uint8) {
	uc.Level = level
}

func (uc *UserCoupon) SetPhone(ctx context.Context, phone string) {
	uc.Phone = phone
}

func (uc *UserCoupon) SetUserCouponStatus(ctx context.Context, userCouponStatus uint8) {
	uc.UserCouponStatus = userCouponStatus
}

func (uc *UserCoupon) SetOrganizationId(ctx context.Context, organizationId uint64) {
	uc.OrganizationId = organizationId
}

func (uc *UserCoupon) SetContent(ctx context.Context) {
	if uc.UserCouponStatus == 1 {
		uc.Content = "未激活"
	} else if uc.UserCouponStatus == 2 {
		uc.Content = tool.FormatPhone(uc.Phone) + "已派发"
	} else if uc.UserCouponStatus == 3 {
		uc.Content = tool.FormatPhone(uc.Phone) + "已激活"
	}
}

func (uc *UserCoupon) SetUpdateTime(ctx context.Context) {
	uc.UpdateTime = time.Now()
}

func (uc *UserCoupon) SetCreateTime(ctx context.Context) {
	uc.CreateTime = time.Now()
}
