package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserContractRepo interface {
	GetContract(context.Context, uint64, uint32) (*v1.GetContractUserContractsReply, error)
	Get(context.Context, uint64, uint32) (*v1.GetUserContractsReply, error)
	Save(context.Context, uint64, uint32, string, string, string) (*v1.CreateUserContractsReply, error)
	Confirm(context.Context, uint64, uint64, string, string) (*v1.ConfirmUserContractsReply, error)
}

type UserContractUsecase struct {
	repo UserContractRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserContractUsecase(repo UserContractRepo, conf *conf.Data, logger log.Logger) *UserContractUsecase {
	return &UserContractUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (ucuc *UserContractUsecase) GetContractMinUserContracts(ctx context.Context, userId uint64, contractType uint32) (*v1.GetContractUserContractsReply, error) {
	contract, err := ucuc.repo.GetContract(ctx, userId, contractType)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_CONTRACT_USER_CONTRACT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return contract, nil
}

func (ucuc *UserContractUsecase) GetMinUserContracts(ctx context.Context, userId uint64, contractType uint32) (*v1.GetUserContractsReply, error) {
	userContract, err := ucuc.repo.Get(ctx, userId, contractType)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_USER_CONTRACT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userContract, nil
}

func (ucuc *UserContractUsecase) CreateMiniUserContracts(ctx context.Context, userId uint64, contractType uint32, name, phone, identityCard string) (*v1.CreateUserContractsReply, error) {
	userContract, err := ucuc.repo.Save(ctx, userId, contractType, name, phone, identityCard)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_CONTRACT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userContract, nil
}

func (ucuc *UserContractUsecase) ConfirmMinUserContracts(ctx context.Context, userId, contractId uint64, phone, code string) (*v1.ConfirmUserContractsReply, error) {
	userContract, err := ucuc.repo.Confirm(ctx, userId, contractId, phone, code)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CONFIRM_USER_CONTRACT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userContract, nil
}
