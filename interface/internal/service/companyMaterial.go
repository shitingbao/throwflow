package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) ListCompanyMaterials(ctx context.Context, in *emptypb.Empty) (*v1.ListCompanyMaterialsReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "product")

	if err != nil {
		return nil, err
	}

	companyMaterials, err := is.cmuc.ListCompanyMaterials(ctx, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	return &v1.ListCompanyMaterialsReply{
		Code: 200,
		Data: &v1.ListCompanyMaterialsReply_Data{
			List: companyMaterials.Data.List,
		},
	}, nil
}
