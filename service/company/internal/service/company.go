package service

import (
	v1 "company/api/company/v1"
	"company/internal/biz"
	"company/internal/pkg/tool"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"time"
)

func (cs *CompanyService) ListCompanys(ctx context.Context, in *v1.ListCompanysRequest) (*v1.ListCompanysReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	companys, err := cs.couc.ListCompanys(ctx, in.PageNum, in.PageSize, in.IndustryId, in.Keyword, in.Status, uint8(in.CompanyType))

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanysReply_Company, 0)

	for _, company := range companys.List {
		contactInformationList := make([]*v1.ListCompanysReply_ContactInformation, 0)

		for _, contactInformation := range company.Clue.ContactInformations {
			contactInformationList = append(contactInformationList, &v1.ListCompanysReply_ContactInformation{
				ContactUsername: contactInformation.ContactUsername,
				ContactPosition: contactInformation.ContactPosition,
				ContactPhone:    contactInformation.ContactPhone,
				ContactWeixin:   contactInformation.ContactWeixin,
			})
		}

		var adminName string
		var adminPhone string

		for _, companyUser := range company.CompanyUser {
			if companyUser.Role == 1 {
				adminName = companyUser.Username
				adminPhone = companyUser.Phone
				break
			}
		}

		list = append(list, &v1.ListCompanysReply_Company{
			Id:                   company.Id,
			CompanyName:          company.Clue.CompanyName,
			IndustryId:           company.Clue.IndustryId,
			IndustryName:         company.Clue.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      uint32(company.Clue.CompanyType),
			QianchuanUse:         uint32(company.Clue.QianchuanUse),
			Sale:                 company.Clue.Sale,
			Seller:               company.Clue.Seller,
			Facilitator:          company.Clue.Facilitator,
			CompanyType:          uint32(company.CompanyType),
			CompanyTypeName:      company.CompanyTypeName,
			Status:               uint32(company.Status),
			StartTime:            tool.TimeToString("2006-01-02", company.StartTime),
			EndTime:              tool.TimeToString("2006-01-02", company.EndTime),
			AdminName:            adminName,
			AdminPhone:           adminPhone,
			ClueId:               company.ClueId,
			MenuId:               company.MenuId,
			Accounts:             company.Accounts,
			QianchuanAdvertisers: company.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.MiniQrCodeUrl,
			IsTermwork:           uint32(company.IsTermwork),
			Address:              company.Clue.Address,
			AreaName:             company.Clue.AreaName,
			AreaCode:             company.Clue.AreaCode,
		})
	}

	totalPage := uint64(math.Ceil(float64(companys.Total) / float64(companys.PageSize)))

	return &v1.ListCompanysReply{
		Code: 200,
		Data: &v1.ListCompanysReply_Data{
			PageNum:   companys.PageNum,
			PageSize:  companys.PageSize,
			Total:     companys.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) ListSelectCompanys(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCompanysReply, error) {
	selects, err := cs.couc.ListSelectCompanys(ctx)

	if err != nil {
		return nil, err
	}

	status := make([]*v1.ListSelectCompanysReply_Status, 0)
	companyType := make([]*v1.ListSelectCompanysReply_CompanyType, 0)

	for _, lstatus := range selects.Status {
		status = append(status, &v1.ListSelectCompanysReply_Status{
			Key:   lstatus.Key,
			Value: lstatus.Value,
		})
	}

	for _, lcompanyType := range selects.CompanyType {
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

func (cs *CompanyService) StatisticsCompanys(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsCompanysReply, error) {
	statistics, err := cs.couc.StatisticsCompanys(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsCompanysReply_Statistic, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsCompanysReply_Statistic{
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

func (cs *CompanyService) GetCompanys(ctx context.Context, in *v1.GetCompanysRequest) (*v1.GetCompanysReply, error) {
	company, err := cs.couc.GetCompanys(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.GetCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.GetCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	var adminName string
	var adminPhone string

	for _, companyUser := range company.CompanyUser {
		if companyUser.Role == 1 {
			adminName = companyUser.Username
			adminPhone = companyUser.Phone
			break
		}
	}

	return &v1.GetCompanysReply{
		Code: 200,
		Data: &v1.GetCompanysReply_Data{
			Id:                   company.Id,
			CompanyName:          company.Clue.CompanyName,
			IndustryId:           company.Clue.IndustryId,
			IndustryName:         company.Clue.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      uint32(company.Clue.CompanyType),
			QianchuanUse:         uint32(company.Clue.QianchuanUse),
			Sale:                 company.Clue.Sale,
			Seller:               company.Clue.Seller,
			Facilitator:          company.Clue.Facilitator,
			CompanyType:          uint32(company.CompanyType),
			CompanyTypeName:      company.CompanyTypeName,
			Status:               uint32(company.Status),
			StartTime:            tool.TimeToString("2006-01-02", company.StartTime),
			EndTime:              tool.TimeToString("2006-01-02", company.EndTime),
			AdminName:            adminName,
			AdminPhone:           adminPhone,
			ClueId:               company.ClueId,
			MenuId:               company.MenuId,
			Accounts:             company.Accounts,
			QianchuanAdvertisers: company.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.MiniQrCodeUrl,
			IsTermwork:           uint32(company.IsTermwork),
			Address:              company.Clue.Address,
			AreaName:             company.Clue.AreaName,
			AreaCode:             company.Clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) CreateCompanys(ctx context.Context, in *v1.CreateCompanysRequest) (*v1.CreateCompanysReply, error) {
	company, err := cs.couc.CreateCompanys(ctx, in.CompanyName, in.ContactInformation, in.Source, in.Seller, in.Facilitator, in.AdminName, in.AdminPhone, in.Address, in.IndustryId, in.UserId, in.ClueId, in.AreaCode, uint8(in.CompanyType), uint8(in.QianchuanUse), uint8(in.Status))

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.CreateCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.CreateCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	var adminName string
	var adminPhone string

	for _, companyUser := range company.CompanyUser {
		if companyUser.Role == 1 {
			adminName = companyUser.Username
			adminPhone = companyUser.Phone
			break
		}
	}

	return &v1.CreateCompanysReply{
		Code: 200,
		Data: &v1.CreateCompanysReply_Data{
			Id:                   company.Id,
			CompanyName:          company.Clue.CompanyName,
			IndustryId:           company.Clue.IndustryId,
			IndustryName:         company.Clue.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      uint32(company.Clue.CompanyType),
			QianchuanUse:         uint32(company.Clue.QianchuanUse),
			Sale:                 company.Clue.Sale,
			Seller:               company.Clue.Seller,
			Facilitator:          company.Clue.Facilitator,
			CompanyType:          uint32(company.CompanyType),
			CompanyTypeName:      company.CompanyTypeName,
			Status:               uint32(company.Status),
			StartTime:            tool.TimeToString("2006-01-02", company.StartTime),
			EndTime:              tool.TimeToString("2006-01-02", company.EndTime),
			AdminName:            adminName,
			AdminPhone:           adminPhone,
			ClueId:               company.ClueId,
			MenuId:               company.MenuId,
			Accounts:             company.Accounts,
			QianchuanAdvertisers: company.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.MiniQrCodeUrl,
			IsTermwork:           uint32(company.IsTermwork),
			Address:              company.Clue.Address,
			AreaName:             company.Clue.AreaName,
			AreaCode:             company.Clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) UpdateCompanys(ctx context.Context, in *v1.UpdateCompanysRequest) (*v1.UpdateCompanysReply, error) {
	company, err := cs.couc.UpdateCompanys(ctx, in.Id, in.CompanyName, in.ContactInformation, in.Seller, in.Facilitator, in.AdminName, in.AdminPhone, in.Address, in.IndustryId, in.ClueId, in.AreaCode, uint8(in.CompanyType), uint8(in.QianchuanUse), uint8(in.Status))

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	var adminName string
	var adminPhone string

	for _, companyUser := range company.CompanyUser {
		if companyUser.Role == 1 {
			adminName = companyUser.Username
			adminPhone = companyUser.Phone
			break
		}
	}

	return &v1.UpdateCompanysReply{
		Code: 200,
		Data: &v1.UpdateCompanysReply_Data{
			Id:                   company.Id,
			CompanyName:          company.Clue.CompanyName,
			IndustryId:           company.Clue.IndustryId,
			IndustryName:         company.Clue.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      uint32(company.Clue.CompanyType),
			QianchuanUse:         uint32(company.Clue.QianchuanUse),
			Sale:                 company.Clue.Sale,
			Seller:               company.Clue.Seller,
			Facilitator:          company.Clue.Facilitator,
			CompanyType:          uint32(company.CompanyType),
			CompanyTypeName:      company.CompanyTypeName,
			Status:               uint32(company.Status),
			StartTime:            tool.TimeToString("2006-01-02", company.StartTime),
			EndTime:              tool.TimeToString("2006-01-02", company.EndTime),
			AdminName:            adminName,
			AdminPhone:           adminPhone,
			ClueId:               company.ClueId,
			MenuId:               company.MenuId,
			Accounts:             company.Accounts,
			QianchuanAdvertisers: company.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.MiniQrCodeUrl,
			IsTermwork:           uint32(company.IsTermwork),
			Address:              company.Clue.Address,
			AreaName:             company.Clue.AreaName,
			AreaCode:             company.Clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) UpdateStatusCompanys(ctx context.Context, in *v1.UpdateStatusCompanysRequest) (*v1.UpdateStatusCompanysReply, error) {
	company, err := cs.couc.UpdateStatusCompanys(ctx, in.Id, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateStatusCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateStatusCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	var adminName string
	var adminPhone string

	for _, companyUser := range company.CompanyUser {
		if companyUser.Role == 1 {
			adminName = companyUser.Username
			adminPhone = companyUser.Phone
			break
		}
	}

	return &v1.UpdateStatusCompanysReply{
		Code: 200,
		Data: &v1.UpdateStatusCompanysReply_Data{
			Id:                   company.Id,
			CompanyName:          company.Clue.CompanyName,
			IndustryId:           company.Clue.IndustryId,
			IndustryName:         company.Clue.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      uint32(company.Clue.CompanyType),
			QianchuanUse:         uint32(company.Clue.QianchuanUse),
			Sale:                 company.Clue.Sale,
			Seller:               company.Clue.Seller,
			Facilitator:          company.Clue.Facilitator,
			CompanyType:          uint32(company.CompanyType),
			CompanyTypeName:      company.CompanyTypeName,
			Status:               uint32(company.Status),
			StartTime:            tool.TimeToString("2006-01-02", company.StartTime),
			EndTime:              tool.TimeToString("2006-01-02", company.EndTime),
			AdminName:            adminName,
			AdminPhone:           adminPhone,
			ClueId:               company.ClueId,
			MenuId:               company.MenuId,
			Accounts:             company.Accounts,
			QianchuanAdvertisers: company.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.MiniQrCodeUrl,
			IsTermwork:           uint32(company.IsTermwork),
			Address:              company.Clue.Address,
			AreaName:             company.Clue.AreaName,
			AreaCode:             company.Clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) UpdateRoleCompanys(ctx context.Context, in *v1.UpdateRoleCompanysRequest) (*v1.UpdateRoleCompanysReply, error) {
	ids, err := cs.verifyMenu(ctx, in.MenuIds)

	if err != nil {
		return nil, err
	}

	if in.Accounts < 0 {
		return nil, biz.CompanyValidatorError
	}

	if in.QianchuanAdvertisers < 0 {
		return nil, biz.CompanyValidatorError
	}

	startTime, err := tool.StringToTime("2006-01-02", in.StartTime)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndTime)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	if !startTime.Equal(endTime) {
		if startTime.After(endTime) {
			return nil, biz.CompanyValidatorError
		}
	}

	company, err := cs.couc.UpdateRoleCompanys(ctx, in.Id, in.UserId, ids, startTime, endTime, in.Accounts, in.QianchuanAdvertisers, uint8(in.CompanyType), uint8(in.IsTermwork))

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateRoleCompanysReply_ContactInformation, 0)

	for _, contactInformation := range company.Clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateRoleCompanysReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	var adminName string
	var adminPhone string

	for _, companyUser := range company.CompanyUser {
		if companyUser.Role == 1 {
			adminName = companyUser.Username
			adminPhone = companyUser.Phone
			break
		}
	}

	return &v1.UpdateRoleCompanysReply{
		Code: 200,
		Data: &v1.UpdateRoleCompanysReply_Data{
			Id:                   company.Id,
			CompanyName:          company.Clue.CompanyName,
			IndustryId:           company.Clue.IndustryId,
			IndustryName:         company.Clue.IndustryName,
			ContactInformations:  contactInformationList,
			ClueCompanyType:      uint32(company.Clue.CompanyType),
			QianchuanUse:         uint32(company.Clue.QianchuanUse),
			Sale:                 company.Clue.Sale,
			Seller:               company.Clue.Seller,
			Facilitator:          company.Clue.Facilitator,
			CompanyType:          uint32(company.CompanyType),
			CompanyTypeName:      company.CompanyTypeName,
			Status:               uint32(company.Status),
			StartTime:            tool.TimeToString("2006-01-02", company.StartTime),
			EndTime:              tool.TimeToString("2006-01-02", company.EndTime),
			AdminName:            adminName,
			AdminPhone:           adminPhone,
			ClueId:               company.ClueId,
			MenuId:               company.MenuId,
			Accounts:             company.Accounts,
			QianchuanAdvertisers: company.QianchuanAdvertisers,
			MiniQrCodeUrl:        company.MiniQrCodeUrl,
			IsTermwork:           uint32(company.IsTermwork),
			Address:              company.Clue.Address,
			AreaName:             company.Clue.AreaName,
			AreaCode:             company.Clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) DeleteCompanys(ctx context.Context, in *v1.DeleteCompanysRequest) (*v1.DeleteCompanysReply, error) {
	err := cs.couc.DeleteCompanys(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteCompanysReply{
		Code: 200,
		Data: &v1.DeleteCompanysReply_Data{},
	}, nil
}

func (cs *CompanyService) SyncUpdateStatusCompanys(ctx context.Context, in *emptypb.Empty) (*v1.SyncUpdateStatusCompanysReply, error) {
	ctx = context.Background()

	if err := cs.couc.SyncUpdateStatusCompanys(ctx); err != nil {
		return nil, err
	}

	return &v1.SyncUpdateStatusCompanysReply{
		Code: 200,
		Data: &v1.SyncUpdateStatusCompanysReply_Data{},
	}, nil
}
