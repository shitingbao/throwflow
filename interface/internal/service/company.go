package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) GetMiniQrCodeCompanys(ctx context.Context, in *emptypb.Empty) (*v1.GetMiniQrCodeCompanysReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	company, err := is.couc.GetMiniQrCodeCompanys(ctx, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	return &v1.GetMiniQrCodeCompanysReply{
		Code: 200,
		Data: &v1.GetMiniQrCodeCompanysReply_Data{
			MiniQrCodeUrl: company.Data.MiniQrCodeUrl,
		},
	}, nil
}
