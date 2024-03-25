package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) GetCompanySets(ctx context.Context, in *v1.GetCompanySetsRequest) (*v1.GetCompanySetsReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "product")

	if err != nil {
		return nil, err
	}

	companySet, err := is.csuc.GetCompanySets(ctx, companyUser.Data.CurrentCompanyId, in.SetKey)

	if err != nil {
		return nil, err
	}

	return &v1.GetCompanySetsReply{
		Code: 200,
		Data: &v1.GetCompanySetsReply_Data{
			SetValue: companySet.Data.SetValue,
		},
	}, nil
}

func (is *InterfaceService) UpdateCompanySets(ctx context.Context, in *v1.UpdateCompanySetsRequest) (*v1.UpdateCompanySetsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	if _, err := is.csuc.UpdateCompanySets(ctx, companyUser.Data.CurrentCompanyId, in.SetKey, in.SetValue); err != nil {
		return nil, err
	}

	return &v1.UpdateCompanySetsReply{
		Code: 200,
		Data: &v1.UpdateCompanySetsReply_Data{},
	}, nil
}
