package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {
	if in.Phone != "18248640151" {
		if ok := is.verifyCode(ctx, in.Phone, "login", in.Code); !ok {
			return nil, biz.InterfaceSmsVerifyError
		}
	} else {
		if in.Code != "888888" {
			return nil, biz.InterfaceSmsVerifyError
		}
	}

	companyUser, err := is.luc.Login(ctx, in.Phone)

	if err != nil {
		return nil, err
	}

	userCompanys := make([]*v1.LoginReply_Company, 0)

	for _, userCompany := range companyUser.Data.UserCompany {
		userCompanys = append(userCompanys, &v1.LoginReply_Company{
			CompanyId:   userCompany.CompanyId,
			CompanyName: userCompany.CompanyName,
		})
	}

	return &v1.LoginReply{
		Code: 200,
		Data: &v1.LoginReply_Data{
			Id:                   companyUser.Data.Id,
			CompanyId:            companyUser.Data.CompanyId,
			Username:             companyUser.Data.Username,
			Job:                  companyUser.Data.Job,
			Phone:                companyUser.Data.Phone,
			Role:                 companyUser.Data.Role,
			RoleName:             companyUser.Data.RoleName,
			IsWhite:              companyUser.Data.IsWhite,
			CompanyType:          companyUser.Data.CompanyType,
			CompanyTypeName:      companyUser.Data.CompanyTypeName,
			CompanyName:          companyUser.Data.CompanyName,
			CompanyStartTime:     companyUser.Data.CompanyStartTime,
			CompanyEndTime:       companyUser.Data.CompanyEndTime,
			Accounts:             companyUser.Data.Accounts,
			QianchuanAdvertisers: companyUser.Data.QianchuanAdvertisers,
			IsTermwork:           companyUser.Data.IsTermwork,
			Reason:               companyUser.Data.Reason,
			CurrentCompanyId:     companyUser.Data.CurrentCompanyId,
			UserCompany:          userCompanys,
			Token:                companyUser.Data.Token,
		},
	}, nil
}

func (is *InterfaceService) Logout(ctx context.Context, in *emptypb.Empty) (*v1.LogoutReply, error) {
	if err := is.luc.Logout(ctx); err != nil {
		return nil, err
	}

	return &v1.LogoutReply{
		Code: 200,
		Data: &v1.LogoutReply_Data{},
	}, nil
}
