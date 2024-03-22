package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/weixin/v1"
	"task/internal/biz"
)

type weixinUserCouponRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserCouponRepo(data *Data, logger log.Logger) biz.WeixinUserCouponRepo {
	return &weixinUserCouponRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (wucr *weixinUserCouponRepo) Sync(ctx context.Context) (*v1.SyncUserCouponsReply, error) {
	userCoupon, err := wucr.data.weixinuc.SyncUserCoupons(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return userCoupon, err
}
