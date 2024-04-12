package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) GetBindMinUserCoupons(ctx context.Context, in *v1.GetBindMinUserCouponsRequest) (*v1.GetBindMinUserCouponsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	bindUserCoupon, err := is.ucuc.GetBindMinUserCoupons(ctx, userInfo.Data.UserId, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	return &v1.GetBindMinUserCouponsReply{
		Code: 200,
		Data: &v1.GetBindMinUserCouponsReply_Data{
			CouponCode: bindUserCoupon.Data.CouponCode,
			Level:      bindUserCoupon.Data.Level,
		},
	}, nil
}

func (is *InterfaceService) ListMinUserCoupons(ctx context.Context, in *v1.ListMinUserCouponsRequest) (*v1.ListMinUserCouponsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userCoupons, err := is.ucuc.ListMinUserCoupons(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId, in.OrganizationId, in.Level)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMinUserCouponsReply_UserCoupon, 0)

	for _, userCoupon := range userCoupons.Data.List {
		list = append(list, &v1.ListMinUserCouponsReply_UserCoupon{
			CouponCode: userCoupon.CouponCode,
			Content:    userCoupon.Content,
		})
	}

	return &v1.ListMinUserCouponsReply{
		Code: 200,
		Data: &v1.ListMinUserCouponsReply_Data{
			PageNum:   userCoupons.Data.PageNum,
			PageSize:  userCoupons.Data.PageSize,
			Total:     userCoupons.Data.Total,
			TotalPage: userCoupons.Data.TotalPage,
			TotalUsed: userCoupons.Data.TotalUsed,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) BindMiniUserCoupons(ctx context.Context, in *v1.BindMiniUserCouponsRequest) (*v1.BindMiniUserCouponsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.ucuc.BindMiniUserCoupons(ctx, userInfo.Data.UserId, in.Level, in.Phone); err != nil {
		return nil, err
	}

	return &v1.BindMiniUserCouponsReply{
		Code: 200,
		Data: &v1.BindMiniUserCouponsReply_Data{},
	}, nil
}

func (is *InterfaceService) ActivateMiniUserCoupons(ctx context.Context, in *v1.ActivateMiniUserCouponsRequest) (*v1.ActivateMiniUserCouponsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.ucuc.ActivateMiniUserCoupons(ctx, userInfo.Data.UserId, in.ParentUserId, in.OrganizationId); err != nil {
		return nil, err
	}

	return &v1.ActivateMiniUserCouponsReply{
		Code: 200,
		Data: &v1.ActivateMiniUserCouponsReply_Data{},
	}, nil
}
