package service

import (
	v1 "company/api/company/v1"
	"context"
)

func (cs *CompanyService) GetCompanySets(ctx context.Context, in *v1.GetCompanySetsRequest) (*v1.GetCompanySetsReply, error) {
	companySet, err := cs.csuc.GetCompanySets(ctx, in.CompanyId, in.Day, in.SetKey)

	if err != nil {
		return nil, err
	}

	return &v1.GetCompanySetsReply{
		Code: 200,
		Data: &v1.GetCompanySetsReply_Data{
			CompanyId: companySet.CompanyId,
			SetKey:    companySet.SetKey,
			SetValue:  companySet.SetValue,
		},
	}, nil
}

func (cs *CompanyService) UpdateCompanySets(ctx context.Context, in *v1.UpdateCompanySetsRequest) (*v1.UpdateCompanySetsReply, error) {
	if err := cs.csuc.UpdateCompanySets(ctx, in.CompanyId, in.SetKey, in.SetValue); err != nil {
		return nil, err
	}

	return &v1.UpdateCompanySetsReply{
		Code: 200,
		Data: &v1.UpdateCompanySetsReply_Data{},
	}, nil
}
