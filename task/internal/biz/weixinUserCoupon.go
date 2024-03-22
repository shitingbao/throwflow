package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/weixin/v1"
	"task/internal/pkg/tool"
)

type WeixinUserCouponRepo interface {
	Sync(ctx context.Context) (*v1.SyncUserCouponsReply, error)
}

type WeixinUserCouponUsecase struct {
	repo WeixinUserCouponRepo
	log  *log.Helper
}

func NewWeixinUserCouponUsecase(repo WeixinUserCouponRepo, logger log.Logger) *WeixinUserCouponUsecase {
	return &WeixinUserCouponUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (wucuc *WeixinUserCouponUsecase) SyncUserCoupons(ctx context.Context) (*v1.SyncUserCouponsReply, error) {
	userCoupon, err := wucuc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_USER_COUPON_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return userCoupon, nil
}
