package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"math"
	"time"
	v1 "weixin/api/weixin/v1"
)

func (ws *WeixinService) GetUserCoupons(ctx context.Context, in *v1.GetUserCouponsRequest) (*v1.GetUserCouponsReply, error) {
	userCoupon, err := ws.ucouc.GetUserCoupons(ctx, in.UserId, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	return &v1.GetUserCouponsReply{
		Code: 200,
		Data: &v1.GetUserCouponsReply_Data{
			CouponCode: userCoupon.CouponCode,
			Level:      uint32(userCoupon.Level),
		},
	}, nil
}

func (ws *WeixinService) ListUserCoupons(ctx context.Context, in *v1.ListUserCouponsRequest) (*v1.ListUserCouponsReply, error) {
	userCoupons, err := ws.ucouc.ListUserCoupons(ctx, in.PageNum, in.PageSize, in.UserId, in.OrganizationId, uint8(in.Level))

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserCouponsReply_UserCoupon, 0)

	for _, userCoupon := range userCoupons.List {
		list = append(list, &v1.ListUserCouponsReply_UserCoupon{
			CouponCode: userCoupon.CouponCode,
			Content:    userCoupon.Content,
		})
	}

	totalPage := uint64(math.Ceil(float64(userCoupons.Total) / float64(userCoupons.PageSize)))

	return &v1.ListUserCouponsReply{
		Code: 200,
		Data: &v1.ListUserCouponsReply_Data{
			PageNum:   userCoupons.PageNum,
			PageSize:  userCoupons.PageSize,
			Total:     userCoupons.Total,
			TotalPage: totalPage,
			TotalUsed: userCoupons.TotalUsed,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) BindUserCoupons(ctx context.Context, in *v1.BindUserCouponsRequest) (*v1.BindUserCouponsReply, error) {
	if err := ws.ucouc.BindUserCoupons(ctx, in.UserId, uint8(in.Level), in.Phone); err != nil {
		return nil, err
	}

	return &v1.BindUserCouponsReply{
		Code: 200,
		Data: &v1.BindUserCouponsReply_Data{},
	}, nil
}

func (ws *WeixinService) ActivateUserCoupons(ctx context.Context, in *v1.ActivateUserCouponsRequest) (*v1.ActivateUserCouponsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ws.ucouc.ActivateUserCoupons(ctx, in.UserId, in.ParentUserId, in.OrganizationId); err != nil {
		return nil, err
	}

	return &v1.ActivateUserCouponsReply{
		Code: 200,
		Data: &v1.ActivateUserCouponsReply_Data{},
	}, nil
}

func (ws *WeixinService) CreateUserCoupons(ctx context.Context, in *v1.CreateUserCouponsRequest) (*v1.CreateUserCouponsReply, error) {
	if err := ws.ucouc.CreateUserCoupons(ctx, in.UserId, in.OrganizationId, in.Num, uint8(in.Level)); err != nil {
		return nil, err
	}

	return &v1.CreateUserCouponsReply{
		Code: 200,
		Data: &v1.CreateUserCouponsReply_Data{},
	}, nil
}

func (ws *WeixinService) SyncUserCoupons(ctx context.Context, in *empty.Empty) (*v1.SyncUserCouponsReply, error) {
	ws.log.Infof("用户券码生成，回收数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.ucouc.SyncUserCoupons(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("用户券码生成，回收数据, 结束时间 %s \n", time.Now())

	return &v1.SyncUserCouponsReply{
		Code: 200,
		Data: &v1.SyncUserCouponsReply_Data{},
	}, nil
}
