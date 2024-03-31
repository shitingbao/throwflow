package domain

import (
	"context"
)

type CouponUser struct {
	Phone string
	Num   uint64
}

func NewCouponUser(ctx context.Context, num uint64, phone string) *CouponUser {
	return &CouponUser{
		Phone: phone,
		Num:   num,
	}
}

func (cu *CouponUser) SetPhone(ctx context.Context, phone string) {
	cu.Phone = phone
}

func (cu *CouponUser) SetNum(ctx context.Context, num uint64) {
	cu.Num = num
}
