package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserCouponRepo interface {
	GetBind(context.Context, uint64, uint64) (*v1.GetUserCouponsReply, error)
	List(context.Context, uint64, uint64, uint64, uint64) (*v1.ListUserCouponsReply, error)
	Bind(context.Context, uint64, string) (*v1.BindUserCouponsReply, error)
	Activate(context.Context, uint64, uint64, uint64) (*v1.ActivateUserCouponsReply, error)
}

type UserCouponUsecase struct {
	repo UserCouponRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserCouponUsecase(repo UserCouponRepo, conf *conf.Data, logger log.Logger) *UserCouponUsecase {
	return &UserCouponUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (ucuc *UserCouponUsecase) GetBindMinUserCoupons(ctx context.Context, userId, organizationId uint64) (*v1.GetUserCouponsReply, error) {
	bindUserCoupon, err := ucuc.repo.GetBind(ctx, userId, organizationId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_BIND_USER_COUPON_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return bindUserCoupon, nil
}

func (ucuc *UserCouponUsecase) ListMinUserCoupons(ctx context.Context, pageNum, pageSize, userId, organizationId uint64) (*v1.ListUserCouponsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(ucuc.conf.Database.PageSize)
	}

	list, err := ucuc.repo.List(ctx, userId, organizationId, pageNum, pageSize)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_COUPON_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (ucuc *UserCouponUsecase) BindMiniUserCoupons(ctx context.Context, userId uint64, phone string) (*v1.BindUserCouponsReply, error) {
	userCoupon, err := ucuc.repo.Bind(ctx, userId, phone)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_BIND_USER_COUPON_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userCoupon, nil
}

func (ucuc *UserCouponUsecase) ActivateMiniUserCoupons(ctx context.Context, userId, parentUserId, organizationId uint64) (*v1.ActivateUserCouponsReply, error) {
	userCoupon, err := ucuc.repo.Activate(ctx, userId, parentUserId, organizationId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_ACTIVATE_USER_COUPON_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userCoupon, nil
}
