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
			EstimatedCommissionBalance:     userBalance.Data.EstimatedCommissionBalance,
			SettleCommissionBalance:        userBalance.Data.SettleCommissionBalance,
			CashableCommissionBalance:      userBalance.Data.CashableCommissionBalance,
			EstimatedCommissionCostBalance: userBalance.Data.EstimatedCommissionCostBalance,
			SettleCommissionCostBalance:    userBalance.Data.SettleCommissionCostBalance,
			CashableCommissionCostBalance:  userBalance.Data.CashableCommissionCostBalance,
		},
	}, nil
}

func (is *InterfaceService) ListMinUserBalances(ctx context.Context, in *v1.ListMinUserBalancesRequest) (*v1.ListMinUserBalancesReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userBalances, err := is.ubuc.ListMinUserBalances(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId, in.OperationType, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMinUserBalancesReply_UserBalance, 0)

	for _, userBalance := range userBalances.Data.List {
		list = append(list, &v1.ListMinUserBalancesReply_UserBalance{
			Amount:             userBalance.Amount,
			NickName:           userBalance.NickName,
			Phone:              userBalance.Phone,
			CommissionTypeName: userBalance.CommissionTypeName,
			OperationType:      userBalance.OperationType,
			OperationContent:   userBalance.OperationContent,
			CreateTime:         userBalance.CreateTime,
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
