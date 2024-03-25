package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) ListMinUserBanks(ctx context.Context, in *v1.ListMinUserBanksRequest) (*v1.ListMinUserBanksReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userBanks, err := is.ubauc.ListMinUserBanks(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMinUserBanksReply_UserBank, 0)

	for _, userBank := range userBanks.Data.List {
		list = append(list, &v1.ListMinUserBanksReply_UserBank{
			BankCode: userBank.BankCode,
			BankName: userBank.BankName,
		})
	}

	return &v1.ListMinUserBanksReply{
		Code: 200,
		Data: &v1.ListMinUserBanksReply_Data{
			PageNum:   userBanks.Data.PageNum,
			PageSize:  userBanks.Data.PageSize,
			Total:     userBanks.Data.Total,
			TotalPage: userBanks.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) CreateMiniUserBanks(ctx context.Context, in *v1.CreateMiniUserBanksRequest) (*v1.CreateMiniUserBanksReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userBank, err := is.ubauc.CreateMiniUserBanks(ctx, userInfo.Data.UserId, in.BankCode)

	if err != nil {
		return nil, err
	}

	return &v1.CreateMiniUserBanksReply{
		Code: 200,
		Data: &v1.CreateMiniUserBanksReply_Data{
			BankCode: userBank.Data.BankCode,
			BankName: userBank.Data.BankName,
		},
	}, nil
}

func (is *InterfaceService) DeleteMiniUserBanks(ctx context.Context, in *v1.DeleteMiniUserBanksRequest) (*v1.DeleteMiniUserBanksReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.ubauc.DeleteMiniUserBanks(ctx, userInfo.Data.UserId, in.BankCode); err != nil {
		return nil, err
	}

	return &v1.DeleteMiniUserBanksReply{
		Code: 200,
		Data: &v1.DeleteMiniUserBanksReply_Data{},
	}, nil
}
