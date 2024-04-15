package service

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"time"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/pkg/tool"
)

func (ws *WeixinService) GetUserBalances(ctx context.Context, in *v1.GetUserBalancesRequest) (*v1.GetUserBalancesReply, error) {
	userBalance, err := ws.ubuc.GetUserBalances(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetUserBalancesReply{
		Code: 200,
		Data: &v1.GetUserBalancesReply_Data{
			EstimatedCommissionBalance:     fmt.Sprintf("%.2f", tool.Decimal(float64(userBalance.EstimatedCommissionBalance), 2)),
			SettleCommissionBalance:        fmt.Sprintf("%.2f", tool.Decimal(float64(userBalance.SettleCommissionBalance), 2)),
			CashableCommissionBalance:      fmt.Sprintf("%.2f", tool.Decimal(float64(userBalance.CashableCommissionBalance), 2)),
			EstimatedCommissionCostBalance: fmt.Sprintf("%.2f", tool.Decimal(float64(userBalance.EstimatedCommissionCostBalance), 2)),
			SettleCommissionCostBalance:    fmt.Sprintf("%.2f", tool.Decimal(float64(userBalance.SettleCommissionCostBalance), 2)),
			CashableCommissionCostBalance:  fmt.Sprintf("%.2f", tool.Decimal(float64(userBalance.CashableCommissionCostBalance), 2)),
		},
	}, nil
}

func (ws *WeixinService) ListUserBalances(ctx context.Context, in *v1.ListUserBalancesRequest) (*v1.ListUserBalancesReply, error) {
	userBalances, err := ws.ubuc.ListUserBalances(ctx, in.PageNum, in.PageSize, in.UserId, uint8(in.OperationType), in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserBalancesReply_UserBalance, 0)

	for _, userBalance := range userBalances.List {
		list = append(list, &v1.ListUserBalancesReply_UserBalance{
			Amount:             tool.Decimal(float64(userBalance.CommissionAmount), 2),
			NickName:           userBalance.ChildNickName,
			Phone:              tool.FormatPhone(userBalance.ChildPhone),
			CommissionTypeName: userBalance.CommissionTypeName,
			OperationType:      uint32(userBalance.OperationType),
			OperationContent:   userBalance.OperationContent,
			CreateTime:         tool.TimeToString("2006/01/02 15:04", userBalance.CreateTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(userBalances.Total) / float64(userBalances.PageSize)))

	return &v1.ListUserBalancesReply{
		Code: 200,
		Data: &v1.ListUserBalancesReply_Data{
			PageNum:   userBalances.PageNum,
			PageSize:  userBalances.PageSize,
			Total:     userBalances.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) CreateUserBalances(ctx context.Context, in *v1.CreateUserBalancesRequest) (*v1.CreateUserBalancesReply, error) {
	if err := ws.ubuc.CreateUserBalances(ctx, in.UserId, uint8(in.BalanceType), in.BankCode, in.Amount); err != nil {
		return nil, err
	}

	return &v1.CreateUserBalancesReply{
		Code: 200,
		Data: &v1.CreateUserBalancesReply_Data{},
	}, nil
}

func (ws *WeixinService) AsyncNotificationUserBalances(ctx context.Context, in *v1.AsyncNotificationUserBalancesRequest) (*v1.AsyncNotificationUserBalancesReply, error) {
	if err := ws.ubuc.AsyncNotificationUserBalances(ctx, in.Content); err != nil {
		return nil, err
	}

	return &v1.AsyncNotificationUserBalancesReply{
		Code: 200,
		Data: &v1.AsyncNotificationUserBalancesReply_Data{},
	}, nil
}

func (ws *WeixinService) SyncUserBalances(ctx context.Context, in *emptypb.Empty) (*v1.SyncUserBalancesReply, error) {
	ws.log.Infof("同步微信用户提现数据到工猫, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.ubuc.SyncUserBalances(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("同步微信用户提现数据到工猫, 结束时间 %s \n", time.Now())

	return &v1.SyncUserBalancesReply{
		Code: 200,
		Data: &v1.SyncUserBalancesReply_Data{},
	}, nil
}
