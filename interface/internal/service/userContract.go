package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) GetContractMinUserContracts(ctx context.Context, in *v1.GetContractMinUserContractsRequest) (*v1.GetContractMinUserContractsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	contract, err := is.ucouc.GetContractMinUserContracts(ctx, userInfo.Data.UserId, in.ContractType)

	if err != nil {
		return nil, err
	}

	return &v1.GetContractMinUserContractsReply{
		Code: 200,
		Data: &v1.GetContractMinUserContractsReply_Data{
			ContractUrl: contract.Data.ContractUrl,
		},
	}, nil
}

func (is *InterfaceService) GetMinUserContracts(ctx context.Context, in *v1.GetMinUserContractsRequest) (*v1.GetMinUserContractsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userContract, err := is.ucouc.GetMinUserContracts(ctx, userInfo.Data.UserId, in.ContractType)

	if err != nil {
		return nil, err
	}

	return &v1.GetMinUserContractsReply{
		Code: 200,
		Data: &v1.GetMinUserContractsReply_Data{
			ContractStatus: userContract.Data.ContractStatus,
		},
	}, nil
}

func (is *InterfaceService) CreateMiniUserContracts(ctx context.Context, in *v1.CreateMiniUserContractsRequest) (*v1.CreateMiniUserContractsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userContract, err := is.ucouc.CreateMiniUserContracts(ctx, userInfo.Data.UserId, in.ContractType, in.Name, in.Phone, in.IdentityCard)

	if err != nil {
		return nil, err
	}

	return &v1.CreateMiniUserContractsReply{
		Code: 200,
		Data: &v1.CreateMiniUserContractsReply_Data{
			ContractId:     userContract.Data.ContractId,
			ContractStatus: userContract.Data.ContractStatus,
		},
	}, nil
}

func (is *InterfaceService) ConfirmMinUserContracts(ctx context.Context, in *v1.ConfirmMinUserContractsRequest) (*v1.ConfirmMinUserContractsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.ucouc.ConfirmMinUserContracts(ctx, userInfo.Data.UserId, in.ContractId, in.Phone, in.Code); err != nil {
		return nil, err
	}

	return &v1.ConfirmMinUserContractsReply{
		Code: 200,
		Data: &v1.ConfirmMinUserContractsReply_Data{},
	}, nil
}
