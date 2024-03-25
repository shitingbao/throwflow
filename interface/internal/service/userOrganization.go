package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) GetMiniUserOrganizationRelations(ctx context.Context, in *empty.Empty) (*v1.GetMiniUserOrganizationRelationsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userOrganization, err := is.uouc.GetMiniUserOrganizationRelations(ctx, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	organizationCourses := make([]*v1.GetMiniUserOrganizationRelationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range userOrganization.Data.OrganizationCourses {
		courseModules := make([]*v1.GetMiniUserOrganizationRelationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.GetMiniUserOrganizationRelationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.GetMiniUserOrganizationRelationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         organizationCourse.CoursePrice,
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
			CourseLevel:         organizationCourse.CourseLevel,
			CourseModules:       courseModules,
		})
	}

	return &v1.GetMiniUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.GetMiniUserOrganizationRelationsReply_Data{
			OrganizationId:            userOrganization.Data.OrganizationId,
			CompanyId:                 userOrganization.Data.CompanyId,
			OrganizationName:          userOrganization.Data.OrganizationName,
			OrganizationLogoUrl:       userOrganization.Data.OrganizationLogoUrl,
			OrganizationCourses:       organizationCourses,
			CompanyName:               userOrganization.Data.CompanyName,
			BankCode:                  userOrganization.Data.BankCode,
			BankDeposit:               userOrganization.Data.BankDeposit,
			ActivationTime:            userOrganization.Data.ActivationTime,
			LevelName:                 userOrganization.Data.LevelName,
			Level:                     userOrganization.Data.Level,
			OrganizationUserQrCodeUrl: userOrganization.Data.OrganizationUserQrCodeUrl,
			ParentUserId:              userOrganization.Data.ParentUserId,
			ParentNickName:            userOrganization.Data.ParentNickName,
			ParentAvatarUrl:           userOrganization.Data.ParentAvatarUrl,
			Total:                     10000 + userOrganization.Data.Total,
		},
	}, nil
}

func (is *InterfaceService) GetBindMiniUserOrganizationRelations(ctx context.Context, in *v1.GetBindMiniUserOrganizationRelationsRequest) (*v1.GetBindMiniUserOrganizationRelationsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	bindUserOrganization, err := is.uouc.GetBindMiniUserOrganizations(ctx, userInfo.Data.UserId, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	mcns := make([]*v1.GetBindMiniUserOrganizationRelationsReply_Mcn, 0)

	for _, mcn := range bindUserOrganization.Data.Mcn {
		mcns = append(mcns, &v1.GetBindMiniUserOrganizationRelationsReply_Mcn{
			Name:          mcn.Name,
			BindStartTime: mcn.BindStartTime,
			BindEndTime:   mcn.BindEndTime,
		})
	}

	return &v1.GetBindMiniUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.GetBindMiniUserOrganizationRelationsReply_Data{
			OrganizationId: bindUserOrganization.Data.OrganizationId,
			CompanyId:      bindUserOrganization.Data.CompanyId,
			ParentNickName: bindUserOrganization.Data.ParentNickName,
			TutorId:        bindUserOrganization.Data.TutorId,
			TutorNickName:  bindUserOrganization.Data.TutorNickName,
			CreateTime:     bindUserOrganization.Data.CreateTime,
			Mcn:            mcns,
		},
	}, nil
}

func (is *InterfaceService) ListMinUserOrders(ctx context.Context, in *v1.ListMinUserOrdersRequest) (*v1.ListMinUserOrdersReply, error) {
	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	userOrders, err := is.uouc.ListMinUserOrders(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMinUserOrdersReply_Materials, 0)

	for _, userOrder := range userOrders.Data.List {
		list = append(list, &v1.ListMinUserOrdersReply_Materials{
			NickName:  userOrder.NickName,
			AvatarUrl: userOrder.AvatarUrl,
			PayTime:   userOrder.PayTime,
		})
	}

	return &v1.ListMinUserOrdersReply{
		Code: 200,
		Data: &v1.ListMinUserOrdersReply_Data{
			PageNum:   userOrders.Data.PageNum,
			PageSize:  userOrders.Data.PageSize,
			Total:     userOrders.Data.Total,
			TotalPage: userOrders.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListParentMiniUserOrganizationRelations(ctx context.Context, in *v1.ListParentMiniUserOrganizationRelationsRequest) (*v1.ListParentMiniUserOrganizationRelationsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	parentUserOrganizationRelations, err := is.uouc.ListParentMiniUserOrganizationRelations(ctx, userInfo.Data.UserId, in.OrganizationId, in.RelationType)

	if err != nil {
		return nil, err
	}

	parentUsers := make([]*v1.ListParentMiniUserOrganizationRelationsReply_ParentUser, 0)

	for _, parentUserOrganizationRelation := range parentUserOrganizationRelations.Data.ParentUser {
		parentUsers = append(parentUsers, &v1.ListParentMiniUserOrganizationRelationsReply_ParentUser{
			ParentUserId:    parentUserOrganizationRelation.ParentUserId,
			ParentNickName:  parentUserOrganizationRelation.ParentNickName,
			ParentAvatarUrl: parentUserOrganizationRelation.ParentAvatarUrl,
			ParentUserType:  parentUserOrganizationRelation.ParentUserType,
		})
	}

	return &v1.ListParentMiniUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.ListParentMiniUserOrganizationRelationsReply_Data{
			ParentUser: parentUsers,
		},
	}, nil
}

func (is *InterfaceService) ListMiniUserOrganizationRelations(ctx context.Context, in *v1.ListMiniUserOrganizationRelationsRequest) (*v1.ListMiniUserOrganizationRelationsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userCommissions, err := is.uouc.ListMiniUserOrganizations(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId, in.OrganizationId, in.IsDirect, in.Month, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMiniUserOrganizationRelationsReply_UserCommission, 0)

	for _, userCommission := range userCommissions.Data.List {
		list = append(list, &v1.ListMiniUserOrganizationRelationsReply_UserCommission{
			NickName:                userCommission.NickName,
			AvatarUrl:               userCommission.AvatarUrl,
			Phone:                   userCommission.Phone,
			ActivationTime:          userCommission.ActivationTime,
			RelationName:            userCommission.RelationName,
			TotalPayAmount:          userCommission.TotalPayAmount,
			CommissionPool:          userCommission.CommissionPool,
			EstimatedUserCommission: userCommission.EstimatedUserCommission,
			CommissionRatio:         userCommission.CommissionRatio,
			RealUserCommission:      userCommission.RealUserCommission,
			UserCommissionTypeName:  userCommission.UserCommissionTypeName,
		})
	}

	return &v1.ListMiniUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.ListMiniUserOrganizationRelationsReply_Data{
			PageNum:   userCommissions.Data.PageNum,
			PageSize:  userCommissions.Data.PageSize,
			Total:     userCommissions.Data.Total,
			TotalPage: userCommissions.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) CreateMiniUserOrders(ctx context.Context, in *v1.CreateMiniUserOrdersRequest) (*v1.CreateMiniUserOrdersReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userOrder, err := is.uouc.CreateMiniUserOrders(ctx, userInfo.Data.UserId, in.ParentUserId, in.OrganizationId, in.PayAmount)

	if err != nil {
		return nil, err
	}

	return &v1.CreateMiniUserOrdersReply{
		Code: 200,
		Data: &v1.CreateMiniUserOrdersReply_Data{
			TimeStamp:  userOrder.Data.TimeStamp,
			NonceStr:   userOrder.Data.NonceStr,
			Package:    userOrder.Data.Package,
			SignType:   userOrder.Data.SignType,
			PaySign:    userOrder.Data.PaySign,
			OutTradeNo: userOrder.Data.OutTradeNo,
			PayAmount:  userOrder.Data.PayAmount,
			LevelName:  userOrder.Data.LevelName,
		},
	}, nil
}

func (is *InterfaceService) UpgradeMiniUserOrders(ctx context.Context, in *v1.UpgradeMiniUserOrdersRequest) (*v1.UpgradeMiniUserOrdersReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userOrder, err := is.uouc.UpgradeMiniUserOrders(ctx, userInfo.Data.UserId, in.OrganizationId, in.PayAmount)

	if err != nil {
		return nil, err
	}

	return &v1.UpgradeMiniUserOrdersReply{
		Code: 200,
		Data: &v1.UpgradeMiniUserOrdersReply_Data{
			TimeStamp:  userOrder.Data.TimeStamp,
			NonceStr:   userOrder.Data.NonceStr,
			Package:    userOrder.Data.Package,
			SignType:   userOrder.Data.SignType,
			PaySign:    userOrder.Data.PaySign,
			OutTradeNo: userOrder.Data.OutTradeNo,
			PayAmount:  userOrder.Data.PayAmount,
			LevelName:  userOrder.Data.LevelName,
		},
	}, nil
}

func (is *InterfaceService) StatisticsMiniUserOrganizationRelations(ctx context.Context, in *v1.StatisticsMiniUserOrganizationRelationsRequest) (*v1.StatisticsMiniUserOrganizationRelationsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	statistics, err := is.uouc.StatisticsMiniUserOrganizations(ctx, userInfo.Data.UserId, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsMiniUserOrganizationRelationsReply_Statistic, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsMiniUserOrganizationRelationsReply_Statistic{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsMiniUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.StatisticsMiniUserOrganizationRelationsReply_Data{
			Statistics: list,
		},
	}, nil
}
