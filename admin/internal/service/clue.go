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

func (as *AdminService) ListClues(ctx context.Context, in *v1.ListCluesRequest) (*v1.ListCluesReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:clue:list"); err != nil {
		return nil, err
	}

	clues, err := as.cluc.ListClues(ctx, in.PageNum, in.IndustryId, in.Keyword, in.Status)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCluesReply_Clues, 0)

	for _, clue := range clues.Data.List {
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
				UserName:   operationLog.UserName,
				Content:    operationLog.Content,
				CreateTime: operationLog.CreateTime,
			})
		}

		list = append(list, &v1.ListCluesReply_Clues{
			Id:                  clue.Id,
			CompanyName:         clue.CompanyName,
			IndustryId:          clue.IndustryId,
			IndustryName:        clue.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         clue.CompanyType,
			QianchuanUse:        clue.QianchuanUse,
			Sale:                clue.Sale,
			Seller:              clue.Seller,
			Facilitator:         clue.Facilitator,
			Source:              clue.Source,
			Status:              clue.Status,
			OperationLogs:       operationLogList,
			Address:             clue.Address,
			AreaName:            clue.AreaName,
			AreaCode:            clue.AreaCode,
			IsAffiliates:        clue.IsAffiliates,
			AdminName:           clue.AdminName,
			AdminPhone:          clue.AdminPhone,
		})
	}

	totalPage := uint64(math.Ceil(float64(clues.Data.Total) / float64(clues.Data.PageSize)))

	return &v1.ListCluesReply{
		Code: 200,
		Data: &v1.ListCluesReply_Data{
			PageNum:   clues.Data.PageNum,
			PageSize:  clues.Data.PageSize,
			Total:     clues.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (as *AdminService) ListSelectClues(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCluesReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	selects, err := as.cluc.ListSelectClues(ctx)

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

func (as *AdminService) StatisticsClues(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsCluesReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	statistics, err := as.cluc.StatisticsClues(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsCluesReply_Statistics, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsCluesReply_Statistics{
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

func (as *AdminService) CreateClues(ctx context.Context, in *v1.CreateCluesRequest) (*v1.CreateCluesReply, error) {
	user, err := as.verifyPermission(ctx, "admin:clue:create")

	if err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	sIndustryId := strings.Split(in.IndustryId, ",")
	sIndustryId = tool.RemoveEmptyString(sIndustryId)

	clue, err := as.cluc.CreateClues(ctx, in.CompanyName, in.ContactInformation, in.Seller, in.Facilitator, in.Address, strings.Join(sIndustryId, ","), in.CompanyType, in.QianchuanUse, user.Id, in.AreaCode)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.CreateCluesReply_ContactInformation, 0)
	operationLogList := make([]*v1.CreateCluesReply_OperationLog, 0)

	for _, contactInformation := range clue.Data.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.CreateCluesReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	for _, operationLog := range clue.Data.OperationLogs {
		operationLogList = append(operationLogList, &v1.CreateCluesReply_OperationLog{
			UserName:   operationLog.UserName,
			Content:    operationLog.Content,
			CreateTime: operationLog.CreateTime,
		})
	}

	return &v1.CreateCluesReply{
		Code: 200,
		Data: &v1.CreateCluesReply_Data{
			Id:                  clue.Data.Id,
			CompanyName:         clue.Data.CompanyName,
			IndustryId:          clue.Data.IndustryId,
			IndustryName:        clue.Data.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         clue.Data.CompanyType,
			QianchuanUse:        clue.Data.QianchuanUse,
			Sale:                clue.Data.Sale,
			Seller:              clue.Data.Seller,
			Facilitator:         clue.Data.Facilitator,
			Source:              clue.Data.Source,
			Status:              clue.Data.Status,
			OperationLogs:       operationLogList,
			Address:             clue.Data.Address,
			AreaName:            clue.Data.AreaName,
			AreaCode:            clue.Data.AreaCode,
		},
	}, nil
}

func (as *AdminService) UpdateClues(ctx context.Context, in *v1.UpdateCluesRequest) (*v1.UpdateCluesReply, error) {
	user, err := as.verifyPermission(ctx, "admin:clue:update")

	if err != nil {
		return nil, err
	}

	sIndustryId := strings.Split(in.IndustryId, ",")
	sIndustryId = tool.RemoveEmptyString(sIndustryId)

	clue, err := as.cluc.UpdateClues(ctx, in.Id, user.Id, in.CompanyName, in.ContactInformation, in.Seller, in.Facilitator, in.Address, strings.Join(sIndustryId, ","), in.CompanyType, in.QianchuanUse, in.AreaCode)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateCluesReply_ContactInformation, 0)
	operationLogList := make([]*v1.UpdateCluesReply_OperationLog, 0)

	for _, contactInformation := range clue.Data.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateCluesReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	for _, operationLog := range clue.Data.OperationLogs {
		operationLogList = append(operationLogList, &v1.UpdateCluesReply_OperationLog{
			UserName:   operationLog.UserName,
			Content:    operationLog.Content,
			CreateTime: operationLog.CreateTime,
		})
	}

	return &v1.UpdateCluesReply{
		Code: 200,
		Data: &v1.UpdateCluesReply_Data{
			Id:                  clue.Data.Id,
			CompanyName:         clue.Data.CompanyName,
			IndustryId:          clue.Data.IndustryId,
			IndustryName:        clue.Data.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         clue.Data.CompanyType,
			QianchuanUse:        clue.Data.QianchuanUse,
			Sale:                clue.Data.Sale,
			Seller:              clue.Data.Seller,
			Facilitator:         clue.Data.Facilitator,
			Source:              clue.Data.Source,
			Status:              clue.Data.Status,
			OperationLogs:       operationLogList,
			Address:             clue.Data.Address,
			AreaName:            clue.Data.AreaName,
			AreaCode:            clue.Data.AreaCode,
		},
	}, nil
}

func (as *AdminService) UpdateOperationLogClues(ctx context.Context, in *v1.UpdateOperationLogCluesRequest) (*v1.UpdateOperationLogCluesReply, error) {
	user, err := as.verifyPermission(ctx, "admin:clue:updateOperationLog")

	if err != nil {
		return nil, err
	}

	operationTime, err := tool.StringToTime("2006-01-02 15:04", in.OperationTime)

	if err != nil {
		return nil, biz.AdminValidatorError
	}

	clue, err := as.cluc.UpdateOperationLogClues(ctx, in.Id, user.Id, in.Content, operationTime)

	if err != nil {
		return nil, err
	}

	contactInformationList := make([]*v1.UpdateOperationLogCluesReply_ContactInformation, 0)
	operationLogList := make([]*v1.UpdateOperationLogCluesReply_OperationLog, 0)

	for _, contactInformation := range clue.Data.ContactInformations {
		contactInformationList = append(contactInformationList, &v1.UpdateOperationLogCluesReply_ContactInformation{
			ContactUsername: contactInformation.ContactUsername,
			ContactPosition: contactInformation.ContactPosition,
			ContactPhone:    contactInformation.ContactPhone,
			ContactWeixin:   contactInformation.ContactWeixin,
		})
	}

	for _, operationLog := range clue.Data.OperationLogs {
		operationLogList = append(operationLogList, &v1.UpdateOperationLogCluesReply_OperationLog{
			UserName:   operationLog.UserName,
			Content:    operationLog.Content,
			CreateTime: operationLog.CreateTime,
		})
	}

	return &v1.UpdateOperationLogCluesReply{
		Code: 200,
		Data: &v1.UpdateOperationLogCluesReply_Data{
			Id:                  clue.Data.Id,
			CompanyName:         clue.Data.CompanyName,
			IndustryId:          clue.Data.IndustryId,
			IndustryName:        clue.Data.IndustryName,
			ContactInformations: contactInformationList,
			CompanyType:         clue.Data.CompanyType,
			QianchuanUse:        clue.Data.QianchuanUse,
			Sale:                clue.Data.Sale,
			Seller:              clue.Data.Seller,
			Facilitator:         clue.Data.Facilitator,
			Source:              clue.Data.Source,
			Status:              clue.Data.Status,
			OperationLogs:       operationLogList,
			Address:             clue.Data.Address,
			AreaName:            clue.Data.AreaName,
			AreaCode:            clue.Data.AreaCode,
		},
	}, nil
}

func (as *AdminService) DeleteClues(ctx context.Context, in *v1.DeleteCluesRequest) (*v1.DeleteCluesReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:clue:delete"); err != nil {
		return nil, err
	}

	if _, err := as.cluc.DeleteClues(ctx, in.Id); err != nil {
		return nil, err
	}

	return &v1.DeleteCluesReply{
		Code: 200,
		Data: &v1.DeleteCluesReply_Data{},
	}, nil
}
