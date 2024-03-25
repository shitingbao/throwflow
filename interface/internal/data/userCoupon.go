package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userCouponRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserCouponRepo(data *Data, logger log.Logger) biz.UserCouponRepo {
	return &userCouponRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ucr *userCouponRepo) GetBind(ctx context.Context, userId, organizationId uint64) (*v1.GetUserCouponsReply, error) {
	userCoupon, err := ucr.data.weixinuc.GetUserCoupons(ctx, &v1.GetUserCouponsRequest{
		UserId:         userId,
		OrganizationId: organizationId,
	})

	if err != nil {
		return nil, err
	}

	return userCoupon, err
}

func (ucr *userCouponRepo) List(ctx context.Context, userId, organizationId, pageNum, pageSize uint64) (*v1.ListUserCouponsReply, error) {
	list, err := ucr.data.weixinuc.ListUserCoupons(ctx, &v1.ListUserCouponsRequest{
		PageNum:        pageNum,
		PageSize:       pageSize,
		UserId:         userId,
		OrganizationId: organizationId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (ucr *userCouponRepo) Bind(ctx context.Context, userId uint64, phone string) (*v1.BindUserCouponsReply, error) {
	userCoupon, err := ucr.data.weixinuc.BindUserCoupons(ctx, &v1.BindUserCouponsRequest{
		UserId: userId,
		Phone:  phone,
	})

	if err != nil {
		return nil, err
	}

	return userCoupon, err
}

func (ucr *userCouponRepo) Activate(ctx context.Context, userId, parentUserId, organizationId uint64) (*v1.ActivateUserCouponsReply, error) {
	userCoupon, err := ucr.data.weixinuc.ActivateUserCoupons(ctx, &v1.ActivateUserCouponsRequest{
		UserId:         userId,
		ParentUserId:   parentUserId,
		OrganizationId: organizationId,
	})

	if err != nil {
		return nil, err
	}

	return userCoupon, err
}
