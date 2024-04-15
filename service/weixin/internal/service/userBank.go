package service

import (
	"context"
	"math"
	v1 "weixin/api/weixin/v1"
)

func (ws *WeixinService) ListUserBanks(ctx context.Context, in *v1.ListUserBanksRequest) (*v1.ListUserBanksReply, error) {
	userBanks, err := ws.ubauc.ListUserBanks(ctx, in.PageNum, in.PageSize, in.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserBanksReply_UserBank, 0)

	for _, userBank := range userBanks.List {
		list = append(list, &v1.ListUserBanksReply_UserBank{
			BankCode: userBank.BankCode,
			BankName: userBank.BankName,
		})
	}

	totalPage := uint64(math.Ceil(float64(userBanks.Total) / float64(userBanks.PageSize)))

	return &v1.ListUserBanksReply{
		Code: 200,
		Data: &v1.ListUserBanksReply_Data{
			PageNum:   userBanks.PageNum,
			PageSize:  userBanks.PageSize,
			Total:     userBanks.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) CreateUserBanks(ctx context.Context, in *v1.CreateUserBanksRequest) (*v1.CreateUserBanksReply, error) {
	userBank, err := ws.ubauc.CreateUserBanks(ctx, in.UserId, in.BankCode)

	if err != nil {
		return nil, err
	}

	return &v1.CreateUserBanksReply{
		Code: 200,
		Data: &v1.CreateUserBanksReply_Data{
			BankCode: userBank.BankCode,
			BankName: userBank.BankName,
		},
	}, nil
}

func (ws *WeixinService) DeleteUserBanks(ctx context.Context, in *v1.DeleteUserBanksRequest) (*v1.DeleteUserBanksReply, error) {
	err := ws.ubauc.DeleteUserBanks(ctx, in.UserId, in.BankCode)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteUserBanksReply{
		Code: 200,
		Data: &v1.DeleteUserBanksReply_Data{},
	}, nil
}

func (ws *WeixinService) DecryptDatas(ctx context.Context, in *v1.DecryptDatasRequest) (*v1.DecryptDatasReply, error) {
	content, err := ws.ubauc.DecryptDatas(ctx, in.OrganizationId, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.DecryptDatasReply{
		Code: 200,
		Data: &v1.DecryptDatasReply_Data{
			Content: content,
		},
	}, nil
}
