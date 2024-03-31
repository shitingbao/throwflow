package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户表
type CouponUser struct {
	Phone string `gorm:"column:phone;type:varchar(50);not null;comment:手机号"`
	Num   uint64 `gorm:"column:num;type:bigint(20) UNSIGNED;not null;default:0;comment:券码数"`
}

func (CouponUser) TableName() string {
	return "coupon_user"
}

type couponUserRepo struct {
	data *Data
	log  *log.Helper
}

func (cu *CouponUser) ToDomain(ctx context.Context) *domain.CouponUser {
	couponUser := &domain.CouponUser{
		Phone: cu.Phone,
		Num:   cu.Num,
	}

	return couponUser
}

func NewCouponUserRepo(data *Data, logger log.Logger) biz.CouponUserRepo {
	return &couponUserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cur *couponUserRepo) List(ctx context.Context) ([]*domain.CouponUser, error) {
	var couponUsers []CouponUser
	list := make([]*domain.CouponUser, 0)

	if result := cur.data.db.WithContext(ctx).
		Find(&couponUsers); result.Error != nil {
		return nil, result.Error
	}

	for _, couponUser := range couponUsers {
		list = append(list, couponUser.ToDomain(ctx))
	}

	return list, nil
}
