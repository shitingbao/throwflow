package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userContractRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserContractRepo(data *Data, logger log.Logger) biz.UserContractRepo {
	return &userContractRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ucr *userContractRepo) GetContract(ctx context.Context, userId uint64, contractType uint32) (*v1.GetContractUserContractsReply, error) {
	contract, err := ucr.data.weixinuc.GetContractUserContracts(ctx, &v1.GetContractUserContractsRequest{
		UserId:       userId,
		ContractType: contractType,
	})

	if err != nil {
		return nil, err
	}

	return contract, err
}

func (ucr *userContractRepo) Get(ctx context.Context, userId uint64, contractType uint32) (*v1.GetUserContractsReply, error) {
	userContract, err := ucr.data.weixinuc.GetUserContracts(ctx, &v1.GetUserContractsRequest{
		UserId:       userId,
		ContractType: contractType,
	})

	if err != nil {
		return nil, err
	}

	return userContract, err
}

func (ucr *userContractRepo) Save(ctx context.Context, userId uint64, contractType uint32, name, phone, identityCard string) (*v1.CreateUserContractsReply, error) {
	userContract, err := ucr.data.weixinuc.CreateUserContracts(ctx, &v1.CreateUserContractsRequest{
		UserId:       userId,
		Name:         name,
		Phone:        phone,
		IdentityCard: identityCard,
		ContractType: contractType,
	})

	if err != nil {
		return nil, err
	}

	return userContract, err
}

func (ucr *userContractRepo) Confirm(ctx context.Context, userId, contractId uint64, phone, code string) (*v1.ConfirmUserContractsReply, error) {
	userContract, err := ucr.data.weixinuc.ConfirmUserContracts(ctx, &v1.ConfirmUserContractsRequest{
		UserId:     userId,
		ContractId: contractId,
		Phone:      phone,
		Code:       code,
	})

	if err != nil {
		return nil, err
	}

	return userContract, err
}
