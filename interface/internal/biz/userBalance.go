package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserBalanceRepo interface {
	Get(context.Context, uint64) (*v1.GetUserBalancesReply, error)
	List(context.Context, uint64, uint64, uint64, uint32) (*v1.ListUserBalancesReply, error)
	Save(context.Context, uint64, uint32, float64, string) (*v1.CreateUserBalancesReply, error)
}

type UserBalanceUsecase struct {
	repo UserBalanceRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserBalanceUsecase(repo UserBalanceRepo, conf *conf.Data, logger log.Logger) *UserBalanceUsecase {
	return &UserBalanceUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (ubuc *UserBalanceUsecase) GetMiniUserBalances(ctx context.Context, userId uint64) (*v1.GetUserBalancesReply, error) {
	userBalance, err := ubuc.repo.Get(ctx, userId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_USER_BALANCE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userBalance, nil
}

func (ubuc *UserBalanceUsecase) ListMinUserBalances(ctx context.Context, pageNum, pageSize, userId uint64, balanceType uint32) (*v1.ListUserBalancesReply, error) {
	if pageSize == 0 {
		pageSize = uint64(ubuc.conf.Database.PageSize)
	}

	list, err := ubuc.repo.List(ctx, userId, pageNum, pageSize, balanceType)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_BALANCE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (ubuc *UserBalanceUsecase) CreateMiniUserBalances(ctx context.Context, userId uint64, balanceType uint32, amount float64, bankCode string) (*v1.CreateUserBalancesReply, error) {
	userBalance, err := ubuc.repo.Save(ctx, userId, balanceType, amount, bankCode)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_BALANCE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userBalance, nil
}
