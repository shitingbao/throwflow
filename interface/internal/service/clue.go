package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) ListSelectClues(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCluesReply, error) {
	selects, err := is.cuc.ListSelectClues(ctx)

	if err != nil {
		return nil, err
	}

	status := make([]*v1.ListSelectCluesReply_Status, 0)
	companyType := make([]*v1.ListSelectCluesReply_CompanyType, 0)
	qianchuanUse := make([]*v1.ListSelectCluesReply_QianchuanUse, 0)

	for _, lstatus := range selects.Data.Status {
		status = append(status, &v1.ListSelectCluesReply_Status{
			Key:   lstatus.Key,
			Value: lstatus.Value,
		})
	}

	for _, lcompanyType := range selects.Data.CompanyType {
		companyType = append(companyType, &v1.ListSelectCluesReply_CompanyType{
			Key:   lcompanyType.Key,
			Value: lcompanyType.Value,
		})
	}

	for _, lqianchuanUse := range selects.Data.QianchuanUse {
		qianchuanUse = append(qianchuanUse, &v1.ListSelectCluesReply_QianchuanUse{
			Key:   lqianchuanUse.Key,
			Value: lqianchuanUse.Value,
		})
	}

	return &v1.ListSelectCluesReply{
		Code: 200,
		Data: &v1.ListSelectCluesReply_Data{
			Status:       status,
			CompanyType:  companyType,
			QianchuanUse: qianchuanUse,
		},
	}, nil
}

func (is *InterfaceService) UpdateClues(ctx context.Context, in *v1.UpdateCluesRequest) (*v1.UpdateCluesReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	if _, err := is.cuc.UpdateClues(ctx, companyUser.Data.CurrentCompanyId, in.CompanyName); err != nil {
		return nil, err
	}

	return &v1.UpdateCluesReply{
		Code: 200,
		Data: &v1.UpdateCluesReply_Data{},
	}, nil
}
