package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userBankRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserBankRepo(data *Data, logger log.Logger) biz.UserBankRepo {
	return &userBankRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ubr *userBankRepo) List(ctx context.Context, userId, pageNum, pageSize uint64) (*v1.ListUserBanksReply, error) {
	list, err := ubr.data.weixinuc.ListUserBanks(ctx, &v1.ListUserBanksRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		UserId:   userId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (ubr *userBankRepo) Save(ctx context.Context, userId uint64, bankCode string) (*v1.CreateUserBanksReply, error) {
	userBank, err := ubr.data.weixinuc.CreateUserBanks(ctx, &v1.CreateUserBanksRequest{
		UserId:   userId,
		BankCode: bankCode,
	})

	if err != nil {
		return nil, err
	}

	return userBank, err
}

func (ubr *userBankRepo) Delete(ctx context.Context, userId uint64, bankCode string) (*v1.DeleteUserBanksReply, error) {
	userBank, err := ubr.data.weixinuc.DeleteUserBanks(ctx, &v1.DeleteUserBanksRequest{
		UserId:   userId,
		BankCode: bankCode,
	})

	if err != nil {
		return nil, err
	}

	return userBank, err
}
