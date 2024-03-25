package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userBalanceRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserBalanceRepo(data *Data, logger log.Logger) biz.UserBalanceRepo {
	return &userBalanceRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ubr *userBalanceRepo) Get(ctx context.Context, userId uint64) (*v1.GetUserBalancesReply, error) {
	userBalance, err := ubr.data.weixinuc.GetUserBalances(ctx, &v1.GetUserBalancesRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return userBalance, err
}

func (ubr *userBalanceRepo) List(ctx context.Context, userId, pageNum, pageSize uint64, operationType uint32) (*v1.ListUserBalancesReply, error) {
	list, err := ubr.data.weixinuc.ListUserBalances(ctx, &v1.ListUserBalancesRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		UserId:        userId,
		OperationType: operationType,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (ubr *userBalanceRepo) Save(ctx context.Context, userId uint64, balanceType uint32, amount float64, bankCode string) (*v1.CreateUserBalancesReply, error) {
	userBalance, err := ubr.data.weixinuc.CreateUserBalances(ctx, &v1.CreateUserBalancesRequest{
		UserId:      userId,
		BankCode:    bankCode,
		Amount:      amount,
		BalanceType: balanceType,
	})

	if err != nil {
		return nil, err
	}

	return userBalance, err
}
