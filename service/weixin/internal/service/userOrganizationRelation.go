package service

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/pkg/tool"
)

func (ws *WeixinService) GetUserOrganizationRelations(ctx context.Context, in *v1.GetUserOrganizationRelationsRequest) (*v1.GetUserOrganizationRelationsReply, error) {
	userOrganizationRelation, err := ws.uoruc.GetUserOrganizationRelations(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	organizationCourses := make([]*v1.GetUserOrganizationRelationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range userOrganizationRelation.OrganizationCourses {
		courseModules := make([]*v1.GetUserOrganizationRelationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.GetUserOrganizationRelationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.GetUserOrganizationRelationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         fmt.Sprintf("%.2f", tool.Decimal(organizationCourse.CoursePrice, 2)),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: fmt.Sprintf("%.2f", tool.Decimal(organizationCourse.CourseOriginalPrice, 2)),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	activationTime := ""

	if !userOrganizationRelation.ActivationTime.IsZero() {
		activationTime = tool.TimeToString("2006/01/02 15:04", userOrganizationRelation.ActivationTime)
	}

	return &v1.GetUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.GetUserOrganizationRelationsReply_Data{
			OrganizationId:            userOrganizationRelation.OrganizationId,
			OrganizationName:          userOrganizationRelation.OrganizationName,
			OrganizationLogoUrl:       userOrganizationRelation.OrganizationLogoUrl,
			OrganizationCourses:       organizationCourses,
			CompanyName:               userOrganizationRelation.CompanyName,
			BankCode:                  userOrganizationRelation.BankCode,
			BankDeposit:               userOrganizationRelation.BankDeposit,
			ActivationTime:            activationTime,
			LevelName:                 userOrganizationRelation.LevelName,
			Level:                     uint32(userOrganizationRelation.Level),
			OrganizationUserQrCodeUrl: userOrganizationRelation.OrganizationUserQrCodeUrl,
			ParentUserId:              userOrganizationRelation.ParentUserId,
			ParentNickName:            userOrganizationRelation.ParentNickName,
			ParentAvatarUrl:           userOrganizationRelation.ParentAvatarUrl,
			Total:                     userOrganizationRelation.Total + 1,
		},
	}, nil
}

func (ws *WeixinService) GetBindUserOrganizationRelations(ctx context.Context, in *v1.GetBindUserOrganizationRelationsRequest) (*v1.GetBindUserOrganizationRelationsReply, error) {
	bindUserOrganizationRelation, err := ws.uoruc.GetBindUserOrganizationRelations(ctx, in.UserId, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	mcns := make([]*v1.GetBindUserOrganizationRelationsReply_Mcn, 0)

	for _, mcn := range bindUserOrganizationRelation.Mcn {
		mcns = append(mcns, &v1.GetBindUserOrganizationRelationsReply_Mcn{
			Name:          mcn.Name,
			BindStartTime: mcn.BindStartTime,
			BindEndTime:   mcn.BindEndTime,
		})
	}

	return &v1.GetBindUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.GetBindUserOrganizationRelationsReply_Data{
			OrganizationId: bindUserOrganizationRelation.OrganizationId,
			ParentNickName: bindUserOrganizationRelation.ParentNickName,
			TutorId:        bindUserOrganizationRelation.TutorId,
			TutorNickName:  bindUserOrganizationRelation.TutorNickName,
			CreateTime:     tool.TimeToString("2006/01/02 15:04", bindUserOrganizationRelation.CreateTime),
			Mcn:            mcns,
		},
	}, nil
}

func (ws *WeixinService) ListParentUserOrganizationRelations(ctx context.Context, in *v1.ListParentUserOrganizationRelationsRequest) (*v1.ListParentUserOrganizationRelationsReply, error) {
	parentUserOrganizationRelations, err := ws.uoruc.ListParentUserOrganizationRelations(ctx, in.UserId, in.OrganizationId, in.RelationType)

	if err != nil {
		return nil, err
	}

	parentUsers := make([]*v1.ListParentUserOrganizationRelationsReply_ParentUser, 0)

	for _, parentUserOrganizationRelation := range parentUserOrganizationRelations {
		parentUsers = append(parentUsers, &v1.ListParentUserOrganizationRelationsReply_ParentUser{
			ParentUserId:    parentUserOrganizationRelation.ParentUserId,
			ParentNickName:  parentUserOrganizationRelation.ParentNickName,
			ParentAvatarUrl: parentUserOrganizationRelation.ParentAvatarUrl,
			ParentUserType:  parentUserOrganizationRelation.ParentUserType,
		})
	}

	return &v1.ListParentUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.ListParentUserOrganizationRelationsReply_Data{
			ParentUser: parentUsers,
		},
	}, nil
}

func (ws *WeixinService) UpdateLevelUserOrganizationRelations(ctx context.Context, in *v1.UpdateLevelUserOrganizationRelationsRequest) (*v1.UpdateLevelUserOrganizationRelationsReply, error) {
	if err := ws.uoruc.UpdateLevelUserOrganizationRelations(ctx, in.UserId, in.OrganizationId, uint8(in.Level)); err != nil {
		return nil, err
	}

	return &v1.UpdateLevelUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.UpdateLevelUserOrganizationRelationsReply_Data{},
	}, nil
}

func (ws *WeixinService) SyncUserOrganizationRelations(ctx context.Context, in *emptypb.Empty) (*v1.SyncUserOrganizationRelationsReply, error) {
	ws.log.Infof("同步微信用户账单机构关系数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.uoruc.SyncUserOrganizationRelations(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("同步微信用户账单机构关系数据, 结束时间 %s \n", time.Now())

	return &v1.SyncUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.SyncUserOrganizationRelationsReply_Data{},
	}, nil
}

func (ws *WeixinService) UpdateTutorUserOrganizationRelations(ctx context.Context, in *emptypb.Empty) (*v1.UpdateTutorUserOrganizationRelationsReply, error) {
	ws.log.Infof("同步微信用户账单机构关系数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.uoruc.UpdateTutorUserOrganizationRelations(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("同步微信用户账单机构关系数据, 结束时间 %s \n", time.Now())

	return &v1.UpdateTutorUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.UpdateTutorUserOrganizationRelationsReply_Data{},
	}, nil
}
