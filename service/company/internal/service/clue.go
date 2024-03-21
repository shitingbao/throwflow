package service

import (
	v1 "company/api/company/v1"
	"company/internal/biz"
	"company/internal/pkg/tool"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
)

func (cs *CompanyService) ListClues(ctx context.Context, in *v1.ListCluesRequest) (*v1.ListCluesReply, error) {
	clues, err := cs.cuc.ListClues(ctx, in.PageNum, in.PageSize, in.IndustryId, in.Keyword, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCluesReply_Clue, 0)

	for _, clue := range clues.List {
		contactInformationList := make([]*v1.ListCluesReply_ContactInformation, 0)
		operationLogList := make([]*v1.ListCluesReply_OperationLog, 0)

		for _, contactInformation := range clue.ContactInformations {
			contactInformationList = append(contactInformationList, &v1.ListCluesReply_ContactInformation{
				ContactUsername: contactInformation.ContactUsername,
				ContactPosition: contactInformation.ContactPosition,
				ContactPhone:    contactInformation.ContactPhone,
				ContactWeixin:   contactInformation.ContactWeixin,
			})
		}

		for _, operationLog := range clue.OperationLogs {
			operationLogList = append(operationLogList, &v1.ListCluesReply_OperationLog{
				UserId:     operationLog.UserId,
				UserName:   operationLog.UserName,
				Content:    operationLog.Content,
				CreateTime: operationLog.CreateTime,
			})
		}

		list = append(list, &v1.ListCluesReply_Clue{
			Id:                  clue.Id,
			CompanyName:         clue.CompanyName,
			IndustryId:          clue.IndustryId,
			IndustryName:        clue.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         uint32(clue.CompanyType),
			QianchuanUse:        uint32(clue.QianchuanUse),
			Sale:                clue.Sale,
			Seller:              clue.Seller,
			Facilitator:         clue.Facilitator,
			Source:              clue.Source,
			Status:              clue.StatusName,
			OperationLogs:       operationLogList,
			Address:             clue.Address,
			AreaName:            clue.AreaName,
			AreaCode:            clue.AreaCode,
			IsAffiliates:        uint32(clue.IsAffiliates),
			AdminName:           clue.AdminName,
			AdminPhone:          clue.AdminPhone,
		})
	}

	totalPage := uint64(math.Ceil(float64(clues.Total) / float64(clues.PageSize)))

	return &v1.ListCluesReply{
		Code: 200,
		Data: &v1.ListCluesReply_Data{
			PageNum:   clues.PageNum,
			PageSize:  clues.PageSize,
			Total:     clues.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) ListSelectClues(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCluesReply, error) {
	selects, err := cs.cuc.ListSelectClues(ctx)

	if err != nil {
		return nil, err
	}

	status := make([]*v1.ListSelectCluesReply_Status, 0)
	companyType := make([]*v1.ListSelectCluesReply_CompanyType, 0)
	qianchuanUse := make([]*v1.ListSelectCluesReply_QianchuanUse, 0)

	for _, lstatus := range selects.Status {
		status = append(status, &v1.ListSelectCluesReply_Status{
			Key:   lstatus.Key,
			Value: lstatus.Value,
		})
	}

	for _, lcompanyType := range selects.CompanyType {
		companyType = append(companyType, &v1.ListSelectCluesReply_CompanyType{
			Key:   lcompanyType.Key,
			Value: lcompanyType.Value,
		})
	}

	for _, lqianchuanUse := range selects.QianchuanUse {
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

func (cs *CompanyService) StatisticsClues(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsCluesReply, error) {
	statistics, err := cs.cuc.StatisticsClues(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsCluesReply_Statistic, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsCluesReply_Statistic{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsCluesReply{
		Code: 200,
		Data: &v1.StatisticsCluesReply_Data{
			Statistics: list,
		},
	}, nil
}

func (cs *CompanyService) CreateClues(ctx context.Context, in *v1.CreateCluesRequest) (*v1.CreateCluesReply, error) {
	clue, err := cs.cuc.CreateClues(ctx, in.CompanyName, in.ContactInformation, in.Source, in.Seller, in.Facilitator, in.Address, in.IndustryId, in.UserId, in.AreaCode, uint8(in.CompanyType), uint8(in.QianchuanUse), uint8(in.Status))

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.CreateCluesReply_ContactInformation, 0)
	operationLogList := make([]*v1.CreateCluesReply_OperationLog, 0)

	for _, contactInformation := range clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.CreateCluesReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	for _, operationLog := range clue.OperationLogs {
		operationLogList = append(operationLogList, &v1.CreateCluesReply_OperationLog{
			UserId:     operationLog.UserId,
			UserName:   operationLog.UserName,
			Content:    operationLog.Content,
			CreateTime: operationLog.CreateTime,
		})
	}

	return &v1.CreateCluesReply{
		Code: 200,
		Data: &v1.CreateCluesReply_Data{
			Id:                  clue.Id,
			CompanyName:         clue.CompanyName,
			IndustryId:          clue.IndustryId,
			IndustryName:        clue.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         uint32(clue.CompanyType),
			QianchuanUse:        uint32(clue.QianchuanUse),
			Sale:                clue.Sale,
			Seller:              clue.Seller,
			Facilitator:         clue.Facilitator,
			Source:              clue.Source,
			Status:              clue.StatusName,
			OperationLogs:       operationLogList,
			Address:             clue.Address,
			AreaName:            clue.AreaName,
			AreaCode:            clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) UpdateClues(ctx context.Context, in *v1.UpdateCluesRequest) (*v1.UpdateCluesReply, error) {
	clue, err := cs.cuc.UpdateClues(ctx, in.Id, in.UserId, in.CompanyName, in.ContactInformation, in.Seller, in.Facilitator, in.Address, in.IndustryId, in.AreaCode, uint8(in.CompanyType), uint8(in.QianchuanUse), uint8(in.Status))

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateCluesReply_ContactInformation, 0)
	operationLogList := make([]*v1.UpdateCluesReply_OperationLog, 0)

	for _, contactInformation := range clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateCluesReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	for _, operationLog := range clue.OperationLogs {
		operationLogList = append(operationLogList, &v1.UpdateCluesReply_OperationLog{
			UserId:     operationLog.UserId,
			UserName:   operationLog.UserName,
			Content:    operationLog.Content,
			CreateTime: operationLog.CreateTime,
		})
	}

	return &v1.UpdateCluesReply{
		Code: 200,
		Data: &v1.UpdateCluesReply_Data{
			Id:                  clue.Id,
			CompanyName:         clue.CompanyName,
			IndustryId:          clue.IndustryId,
			IndustryName:        clue.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         uint32(clue.CompanyType),
			QianchuanUse:        uint32(clue.QianchuanUse),
			Sale:                clue.Sale,
			Seller:              clue.Seller,
			Facilitator:         clue.Facilitator,
			Source:              clue.Source,
			Status:              clue.StatusName,
			OperationLogs:       operationLogList,
			Address:             clue.Address,
			AreaName:            clue.AreaName,
			AreaCode:            clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) UpdateCompanyNameClues(ctx context.Context, in *v1.UpdateCompanyNameCluesRequest) (*v1.UpdateCompanyNameCluesReply, error) {
	clue, err := cs.cuc.UpdateCompanyNameClues(ctx, in.CompanyId, in.CompanyName)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateCompanyNameCluesReply_ContactInformation, 0)
	operationLogList := make([]*v1.UpdateCompanyNameCluesReply_OperationLog, 0)

	for _, contactInformation := range clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateCompanyNameCluesReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	for _, operationLog := range clue.OperationLogs {
		operationLogList = append(operationLogList, &v1.UpdateCompanyNameCluesReply_OperationLog{
			UserId:     operationLog.UserId,
			UserName:   operationLog.UserName,
			Content:    operationLog.Content,
			CreateTime: operationLog.CreateTime,
		})
	}

	return &v1.UpdateCompanyNameCluesReply{
		Code: 200,
		Data: &v1.UpdateCompanyNameCluesReply_Data{
			Id:                  clue.Id,
			CompanyName:         clue.CompanyName,
			IndustryId:          clue.IndustryId,
			IndustryName:        clue.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         uint32(clue.CompanyType),
			QianchuanUse:        uint32(clue.QianchuanUse),
			Sale:                clue.Sale,
			Seller:              clue.Seller,
			Facilitator:         clue.Facilitator,
			Source:              clue.Source,
			Status:              clue.StatusName,
			OperationLogs:       operationLogList,
			Address:             clue.Address,
			AreaName:            clue.AreaName,
			AreaCode:            clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) UpdateOperationLogClues(ctx context.Context, in *v1.UpdateOperationLogCluesRequest) (*v1.UpdateOperationLogCluesReply, error) {
	operationTime, err := tool.StringToTime("2006-01-02 15:04", in.OperationTime)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	clue, err := cs.cuc.UpdateOperationLogClues(ctx, in.Id, in.UserId, in.Content, operationTime)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateOperationLogCluesReply_ContactInformation, 0)
	operationLogList := make([]*v1.UpdateOperationLogCluesReply_OperationLog, 0)

	for _, contactInformation := range clue.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateOperationLogCluesReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	for _, operationLog := range clue.OperationLogs {
		operationLogList = append(operationLogList, &v1.UpdateOperationLogCluesReply_OperationLog{
			UserId:     operationLog.UserId,
			UserName:   operationLog.UserName,
			Content:    operationLog.Content,
			CreateTime: operationLog.CreateTime,
		})
	}

	return &v1.UpdateOperationLogCluesReply{
		Code: 200,
		Data: &v1.UpdateOperationLogCluesReply_Data{
			Id:                  clue.Id,
			CompanyName:         clue.CompanyName,
			IndustryId:          clue.IndustryId,
			IndustryName:        clue.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         uint32(clue.CompanyType),
			QianchuanUse:        uint32(clue.QianchuanUse),
			Sale:                clue.Sale,
			Seller:              clue.Seller,
			Facilitator:         clue.Facilitator,
			Source:              clue.Source,
			Status:              clue.StatusName,
			OperationLogs:       operationLogList,
			Address:             clue.Address,
			AreaName:            clue.AreaName,
			AreaCode:            clue.AreaCode,
		},
	}, nil
}

func (cs *CompanyService) DeleteClues(ctx context.Context, in *v1.DeleteCluesRequest) (*v1.DeleteCluesReply, error) {
	err := cs.cuc.DeleteClues(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteCluesReply{
		Code: 200,
		Data: &v1.DeleteCluesReply_Data{},
	}, nil
}
