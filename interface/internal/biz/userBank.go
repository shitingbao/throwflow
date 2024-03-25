package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserBankRepo interface {
	List(context.Context, uint64, uint64, uint64) (*v1.ListUserBanksReply, error)
	Save(context.Context, uint64, string) (*v1.CreateUserBanksReply, error)
	Delete(context.Context, uint64, string) (*v1.DeleteUserBanksReply, error)
}

type UserBankUsecase struct {
	repo UserBankRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserBankUsecase(repo UserBankRepo, conf *conf.Data, logger log.Logger) *UserBankUsecase {
	return &UserBankUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (ubuc *UserBankUsecase) ListMinUserBanks(ctx context.Context, pageNum, pageSize, userId uint64) (*v1.ListUserBanksReply, error) {
	if pageSize == 0 {
		pageSize = uint64(ubuc.conf.Database.PageSize)
	}

	list, err := ubuc.repo.List(ctx, userId, pageNum, pageSize)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_BANK_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (ubuc *UserBankUsecase) CreateMiniUserBanks(ctx context.Context, userId uint64, bankCode string) (*v1.CreateUserBanksReply, error) {
	userBank, err := ubuc.repo.Save(ctx, userId, bankCode)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_BANK_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userBank, nil
}

func (ubuc *UserBankUsecase) DeleteMiniUserBanks(ctx context.Context, userId uint64, bankCode string) (*v1.DeleteUserBanksReply, error) {
	userBank, err := ubuc.repo.Delete(ctx, userId, bankCode)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_DELETE_USER_BANK_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userBank, nil
}
