package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) GetMiniUserBalances(ctx context.Context, in *empty.Empty) (*v1.GetMiniUserBalancesReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userBalance, err := is.ubuc.GetMiniUserBalances(ctx, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetMiniUserBalancesReply{
		Code: 200,
		Data: &v1.GetMiniUserBalancesReply_Data{
			EstimatedCommissionBalance: userBalance.Data.EstimatedCommissionBalance,
			RealCommissionBalance:      userBalance.Data.RealCommissionBalance,
			EstimatedCostBalance:       userBalance.Data.EstimatedCostBalance,
			RealCostBalance:            userBalance.Data.RealCostBalance,
		},
	}, nil
}

func (is *InterfaceService) ListMinUserBalances(ctx context.Context, in *v1.ListMinUserBalancesRequest) (*v1.ListMinUserBalancesReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userBalances, err := is.ubuc.ListMinUserBalances(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId, in.OperationType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMinUserBalancesReply_UserBalance, 0)

	for _, userBalance := range userBalances.Data.List {
		list = append(list, &v1.ListMinUserBalancesReply_UserBalance{
			Amount:           userBalance.Amount,
			BalanceType:      userBalance.BalanceType,
			OperationType:    userBalance.OperationType,
			OperationContent: userBalance.OperationContent,
			BalanceStatus:    userBalance.BalanceStatus,
			CreateTime:       userBalance.CreateTime,
		})
	}

	return &v1.ListMinUserBalancesReply{
		Code: 200,
		Data: &v1.ListMinUserBalancesReply_Data{
			PageNum:   userBalances.Data.PageNum,
			PageSize:  userBalances.Data.PageSize,
			Total:     userBalances.Data.Total,
			TotalPage: userBalances.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) CreateMiniUserBalances(ctx context.Context, in *v1.CreateMiniUserBalancesRequest) (*v1.CreateMiniUserBalancesReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.ubuc.CreateMiniUserBalances(ctx, userInfo.Data.UserId, in.BalanceType, in.Amount, in.BankCode); err != nil {
		return nil, err
	}

	return &v1.CreateMiniUserBalancesReply{
		Code: 200,
		Data: &v1.CreateMiniUserBalancesReply_Data{},
	}, nil
}
