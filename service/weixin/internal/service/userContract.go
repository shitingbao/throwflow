package service

import (
	"context"
	v1 "weixin/api/weixin/v1"
)

func (ws *WeixinService) GetContractUserContracts(ctx context.Context, in *v1.GetContractUserContractsRequest) (*v1.GetContractUserContractsReply, error) {
	contract, err := ws.uconuc.GetContractUserContracts(ctx, in.UserId, uint8(in.ContractType))

	if err != nil {
		return nil, err
	}

	return &v1.GetContractUserContractsReply{
		Code: 200,
		Data: &v1.GetContractUserContractsReply_Data{
			ContractUrl: contract.ContractAddr,
		},
	}, nil
}

func (ws *WeixinService) GetUserContracts(ctx context.Context, in *v1.GetUserContractsRequest) (*v1.GetUserContractsReply, error) {
	userContract, err := ws.uconuc.GetUserContracts(ctx, in.UserId, uint8(in.ContractType))

	if err != nil {
		return nil, err
	}

	return &v1.GetUserContractsReply{
		Code: 200,
		Data: &v1.GetUserContractsReply_Data{
			ContractStatus: uint32(userContract.ContractStatus),
		},
	}, nil
}

func (ws *WeixinService) CreateUserContracts(ctx context.Context, in *v1.CreateUserContractsRequest) (*v1.CreateUserContractsReply, error) {
	userContract, err := ws.uconuc.CreateUserContracts(ctx, in.UserId, uint8(in.ContractType), in.Name, in.Phone, in.IdentityCard)

	if err != nil {
		return nil, err
	}

	return &v1.CreateUserContractsReply{
		Code: 200,
		Data: &v1.CreateUserContractsReply_Data{
			ContractId:     userContract.ContractId,
			ContractStatus: uint32(userContract.ContractStatus),
		},
	}, nil
}

func (ws *WeixinService) ConfirmUserContracts(ctx context.Context, in *v1.ConfirmUserContractsRequest) (*v1.ConfirmUserContractsReply, error) {
	if err := ws.uconuc.ConfirmUserContracts(ctx, in.UserId, in.ContractId, in.Phone, in.Code); err != nil {
		return nil, err
	}

	return &v1.ConfirmUserContractsReply{
		Code: 200,
		Data: &v1.ConfirmUserContractsReply_Data{},
	}, nil
}

func (ws *WeixinService) AsyncNotificationUserContracts(ctx context.Context, in *v1.AsyncNotificationUserContractsRequest) (*v1.AsyncNotificationUserContractsReply, error) {
	if err := ws.uconuc.AsyncNotificationUserContracts(ctx, in.Content); err != nil {
		return nil, err
	}

	return &v1.AsyncNotificationUserContractsReply{
		Code: 200,
		Data: &v1.AsyncNotificationUserContractsReply_Data{},
	}, nil
}
