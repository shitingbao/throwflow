package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/pkg/tool"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"strings"
)

func (as *AdminService) ListCompanys(ctx context.Context, in *v1.ListCompanysRequest) (*v1.ListCompanysReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:company:list"); err != nil {
		return nil, err
	}

	companys, err := as.couc.ListCompanys(ctx, in.PageNum, in.IndustryId, in.Keyword, in.Status, in.CompanyType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanysReply_Companys, 0)

	for _, company := range companys.Data.List {
		contactInformationList := make([]*v1.ListCompanysReply_ContactInformation, 0)

		for _, contactInformation := range company.ContactInformations {
			contactInformationList = append(contactInformationList, &v1.ListCompanysReply_ContactInformation{
				ContactUsername: contactInformation.ContactUsername,
				ContactPosition: contactInformation.ContactPosition,
				ContactPhone:    contactInformation.ContactPhone,
				ContactWeixin:   contactInformation.ContactWeixin,
			})
		}

		list = append(list, &v1.ListCompanysReply_Companys{
			Id:                   company.Id,
			CompanyName:          company.CompanyName,
			IndustryId:           company.IndustryId,
			IndustryName:         company.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      company.ClueCompanyType,
			QianchuanUse:         company.QianchuanUse,
			Sale:                 company.Sale,
			Seller:               company.Seller,
			Facilitator:          company.Facilitator,
			CompanyType:          company.CompanyType,
			CompanyTypeName:      company.CompanyTypeName,
			Status:               company.Status,
			StartTime:            company.StartTime,
			EndTime:              company.EndTime,
			AdminName:            company.AdminName,
			AdminPhone:           company.AdminPhone,
			ClueId:               company.ClueId,
			MenuId:               company.MenuId,
			Accounts:             company.Accounts,
			QianchuanAdvertisers: company.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.MiniQrCodeUrl,
			IsTermwork:           company.IsTermwork,
			Address:              company.Address,
			AreaName:             company.AreaName,
			AreaCode:             company.AreaCode,
		})
	}

	totalPage := uint64(math.Ceil(float64(companys.Data.Total) / float64(companys.Data.PageSize)))

	return &v1.ListCompanysReply{
		Code: 200,
		Data: &v1.ListCompanysReply_Data{
			PageNum:   companys.Data.PageNum,
			PageSize:  companys.Data.PageSize,
			Total:     companys.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (as *AdminService) ListSelectCompanys(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCompanysReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	selects, err := as.couc.ListSelectCompanys(ctx)

	if err != nil {
		return nil, err
	}

	status := make([]*v1.ListSelectCompanysReply_Status, 0)
	companyType := make([]*v1.ListSelectCompanysReply_CompanyType, 0)

	for _, lstatus := range selects.Data.Status {
		status = append(status, &v1.ListSelectCompanysReply_Status{
			Key:   lstatus.Key,
			Value: lstatus.Value,
		})
	}

	for _, lcompanyType := range selects.Data.CompanyType {
		companyType = append(companyType, &v1.ListSelectCompanysReply_CompanyType{
			Key:   lcompanyType.Key,
			Value: lcompanyType.Value,
		})
	}

	return &v1.ListSelectCompanysReply{
		Code: 200,
		Data: &v1.ListSelectCompanysReply_Data{
			Status:      status,
			CompanyType: companyType,
		},
	}, nil
}

func (as *AdminService) StatisticsCompanys(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsCompanysReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	statistics, err := as.couc.StatisticsCompanys(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsCompanysReply_Statistics, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsCompanysReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsCompanysReply{
		Code: 200,
		Data: &v1.StatisticsCompanysReply_Data{
			Statistics: list,
		},
	}, nil
}

func (as *AdminService) CreateCompanys(ctx context.Context, in *v1.CreateCompanysRequest) (*v1.CreateCompanysReply, error) {
	user, err := as.verifyPermission(ctx, "admin:company:create")

	if err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	sIndustryId := strings.Split(in.IndustryId, ",")
	sIndustryId = tool.RemoveEmptyString(sIndustryId)

	company, err := as.couc.CreateCompanys(ctx, in.CompanyName, in.ContactInformation, in.Seller, in.Facilitator, in.AdminName, in.AdminPhone, in.Address, strings.Join(sIndustryId, ","), in.CompanyType, in.QianchuanUse, user.Id, in.ClueId, in.AreaCode)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.CreateCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Data.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.CreateCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	return &v1.CreateCompanysReply{
		Code: 200,
		Data: &v1.CreateCompanysReply_Data{
			Id:                   company.Data.Id,
			CompanyName:          company.Data.CompanyName,
			IndustryId:           company.Data.IndustryId,
			IndustryName:         company.Data.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      company.Data.ClueCompanyType,
			QianchuanUse:         company.Data.QianchuanUse,
			Sale:                 company.Data.Sale,
			Seller:               company.Data.Seller,
			Facilitator:          company.Data.Facilitator,
			CompanyType:          company.Data.CompanyType,
			CompanyTypeName:      company.Data.CompanyTypeName,
			Status:               company.Data.Status,
			StartTime:            company.Data.StartTime,
			EndTime:              company.Data.EndTime,
			AdminName:            company.Data.AdminName,
			AdminPhone:           company.Data.AdminPhone,
			ClueId:               company.Data.ClueId,
			MenuId:               company.Data.MenuId,
			Accounts:             company.Data.Accounts,
			QianchuanAdvertisers: company.Data.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.Data.MiniQrCodeUrl,
			IsTermwork:           company.Data.IsTermwork,
			Address:              company.Data.Address,
			AreaName:             company.Data.AreaName,
			AreaCode:             company.Data.AreaCode,
		},
	}, nil
}

func (as *AdminService) UpdateCompanys(ctx context.Context, in *v1.UpdateCompanysRequest) (*v1.UpdateCompanysReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:company:update"); err != nil {
		return nil, err
	}

	sIndustryId := strings.Split(in.IndustryId, ",")
	sIndustryId = tool.RemoveEmptyString(sIndustryId)

	company, err := as.couc.UpdateCompanys(ctx, in.Id, in.CompanyName, in.ContactInformation, in.Seller, in.Facilitator, in.AdminName, in.AdminPhone, in.Address, strings.Join(sIndustryId, ","), in.CompanyType, in.QianchuanUse, in.ClueId, in.AreaCode)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Data.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	return &v1.UpdateCompanysReply{
		Code: 200,
		Data: &v1.UpdateCompanysReply_Data{
			Id:                   company.Data.Id,
			CompanyName:          company.Data.CompanyName,
			IndustryId:           company.Data.IndustryId,
			IndustryName:         company.Data.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      company.Data.ClueCompanyType,
			QianchuanUse:         company.Data.QianchuanUse,
			Sale:                 company.Data.Sale,
			Seller:               company.Data.Seller,
			Facilitator:          company.Data.Facilitator,
			CompanyType:          company.Data.CompanyType,
			CompanyTypeName:      company.Data.CompanyTypeName,
			Status:               company.Data.Status,
			StartTime:            company.Data.StartTime,
			EndTime:              company.Data.EndTime,
			AdminName:            company.Data.AdminName,
			AdminPhone:           company.Data.AdminPhone,
			ClueId:               company.Data.ClueId,
			MenuId:               company.Data.MenuId,
			Accounts:             company.Data.Accounts,
			QianchuanAdvertisers: company.Data.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.Data.MiniQrCodeUrl,
			IsTermwork:           company.Data.IsTermwork,
			Address:              company.Data.Address,
			AreaName:             company.Data.AreaName,
			AreaCode:             company.Data.AreaCode,
		},
	}, nil
}

func (as *AdminService) UpdateStatusCompanys(ctx context.Context, in *v1.UpdateStatusCompanysRequest) (*v1.UpdateStatusCompanysReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:company:updateStatus"); err != nil {
		return nil, err
	}

	ctx = context.Background()

	company, err := as.couc.UpdateStatusCompanys(ctx, in.Id, in.Status)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateStatusCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Data.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateStatusCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	return &v1.UpdateStatusCompanysReply{
		Code: 200,
		Data: &v1.UpdateStatusCompanysReply_Data{
			Id:                   company.Data.Id,
			CompanyName:          company.Data.CompanyName,
			IndustryId:           company.Data.IndustryId,
			IndustryName:         company.Data.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      company.Data.ClueCompanyType,
			QianchuanUse:         company.Data.QianchuanUse,
			Sale:                 company.Data.Sale,
			Seller:               company.Data.Seller,
			Facilitator:          company.Data.Facilitator,
			CompanyType:          company.Data.CompanyType,
			CompanyTypeName:      company.Data.CompanyTypeName,
			Status:               company.Data.Status,
			StartTime:            company.Data.StartTime,
			EndTime:              company.Data.EndTime,
			AdminName:            company.Data.AdminName,
			AdminPhone:           company.Data.AdminPhone,
			ClueId:               company.Data.ClueId,
			MenuId:               company.Data.MenuId,
			Accounts:             company.Data.Accounts,
			QianchuanAdvertisers: company.Data.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.Data.MiniQrCodeUrl,
			IsTermwork:           company.Data.IsTermwork,
			Address:              company.Data.Address,
			AreaName:             company.Data.AreaName,
			AreaCode:             company.Data.AreaCode,
		},
	}, nil
}

func (as *AdminService) UpdateRoleCompanys(ctx context.Context, in *v1.UpdateRoleCompanysRequest) (*v1.UpdateRoleCompanysReply, error) {
	user, err := as.verifyPermission(ctx, "admin:company:updateRole")

	if err != nil {
		return nil, err
	}

	if in.Accounts < 0 {
		return nil, biz.AdminValidatorError
	}

	if in.QianchuanAdvertisers < 0 {
		return nil, biz.AdminValidatorError
	}

	startTime, err := tool.StringToTime("2006-01-02", in.StartTime)

	if err != nil {
		return nil, biz.AdminValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndTime)

	if err != nil {
		return nil, biz.AdminValidatorError
	}

	if startTime.After(endTime) {
		return nil, biz.AdminValidatorError
	}

	company, err := as.couc.UpdateRoleCompanys(ctx, in.Id, user.Id, in.MenuIds, in.StartTime, in.EndTime, in.Accounts, in.QianchuanAdvertisers, in.CompanyType, in.IsTermwork)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateRoleCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Data.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateRoleCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	return &v1.UpdateRoleCompanysReply{
		Code: 200,
		Data: &v1.UpdateRoleCompanysReply_Data{
			Id:                   company.Data.Id,
			CompanyName:          company.Data.CompanyName,
			IndustryId:           company.Data.IndustryId,
			IndustryName:         company.Data.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      company.Data.ClueCompanyType,
			QianchuanUse:         company.Data.QianchuanUse,
			Sale:                 company.Data.Sale,
			Seller:               company.Data.Seller,
			Facilitator:          company.Data.Facilitator,
			CompanyType:          company.Data.CompanyType,
			CompanyTypeName:      company.Data.CompanyTypeName,
			Status:               company.Data.Status,
			StartTime:            company.Data.StartTime,
			EndTime:              company.Data.EndTime,
			AdminName:            company.Data.AdminName,
			AdminPhone:           company.Data.AdminPhone,
			ClueId:               company.Data.ClueId,
			MenuId:               company.Data.MenuId,
			Accounts:             company.Data.Accounts,
			QianchuanAdvertisers: company.Data.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.Data.MiniQrCodeUrl,
			IsTermwork:           company.Data.IsTermwork,
			Address:              company.Data.Address,
			AreaName:             company.Data.AreaName,
			AreaCode:             company.Data.AreaCode,
		},
	}, nil
}

func (as *AdminService) DeleteCompanys(ctx context.Context, in *v1.DeleteCompanysRequest) (*v1.DeleteCompanysReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:company:delete"); err != nil {
		return nil, err
	}

	if _, err := as.couc.DeleteCompanys(ctx, in.Id); err != nil {
		return nil, err
	}

	return &v1.DeleteCompanysReply{
		Code: 200,
		Data: &v1.DeleteCompanysReply_Data{},
	}, nil
}
