package service

import (
	v1 "company/api/company/v1"
	"company/internal/pkg/tool"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
)

func (cs *CompanyService) GetCompanyOrganizations(ctx context.Context, in *v1.GetCompanyOrganizationsRequest) (*v1.GetCompanyOrganizationsReply, error) {
	companyOrganization, err := cs.coruc.GetCompanyOrganizations(ctx, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.GetCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.GetCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn,
		})
	}

	organizationCommission := &v1.GetCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
		OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
	}

	organizationColonelCommission := &v1.GetCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
		PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
		IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
		IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
		AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
		CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
		PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
		PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
		IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
		IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
		AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
		OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
		IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
		IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
		AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
		CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
		PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
		PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
		IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
		IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
		AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
	}

	organizationCourses := make([]*v1.GetCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.OrganizationCourses {
		courseModules := make([]*v1.GetCompanyOrganizationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.GetCompanyOrganizationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.GetCompanyOrganizationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.GetCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.GetCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.GetCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.GetCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (cs *CompanyService) GetCompanyOrganizationByOrganizationCodes(ctx context.Context, in *v1.GetCompanyOrganizationByOrganizationCodesRequest) (*v1.GetCompanyOrganizationByOrganizationCodesReply, error) {
	companyOrganization, err := cs.coruc.GetCompanyOrganizationByOrganizationCodes(ctx, in.OrganizationCode)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.GetCompanyOrganizationByOrganizationCodesReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.GetCompanyOrganizationByOrganizationCodesReply_OrganizationMcn{
			OrganizationMcn: organizationMcn,
		})
	}

	organizationCommission := &v1.GetCompanyOrganizationByOrganizationCodesReply_OrganizationCommission{
		CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
		OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
	}

	organizationColonelCommission := &v1.GetCompanyOrganizationByOrganizationCodesReply_OrganizationColonelCommission{
		ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
		PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
		IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
		IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
		AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
		CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
		PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
		PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
		IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
		IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
		AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
		OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
		IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
		IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
		AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
		CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
		PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
		PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
		IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
		IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
		AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
	}

	organizationCourses := make([]*v1.GetCompanyOrganizationByOrganizationCodesReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.OrganizationCourses {
		courseModules := make([]*v1.GetCompanyOrganizationByOrganizationCodesReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.GetCompanyOrganizationByOrganizationCodesReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.GetCompanyOrganizationByOrganizationCodesReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.GetCompanyOrganizationByOrganizationCodesReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.GetCompanyOrganizationByOrganizationCodesReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.GetCompanyOrganizationByOrganizationCodesReply{
		Code: 200,
		Data: &v1.GetCompanyOrganizationByOrganizationCodesReply_Data{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyOrganizations(ctx context.Context, in *v1.ListCompanyOrganizationsRequest) (*v1.ListCompanyOrganizationsReply, error) {
	companyOrganizations, err := cs.coruc.ListCompanyOrganizations(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyOrganizationsReply_CompanyOrganization, 0)

	for _, companyOrganization := range companyOrganizations.List {
		organizationMcns := make([]*v1.ListCompanyOrganizationsReply_OrganizationMcn, 0)

		for _, organizationMcn := range companyOrganization.OrganizationMcns {
			organizationMcns = append(organizationMcns, &v1.ListCompanyOrganizationsReply_OrganizationMcn{
				OrganizationMcn: organizationMcn,
			})
		}

		organizationCommission := &v1.ListCompanyOrganizationsReply_OrganizationCommission{
			CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
			OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
		}

		organizationColonelCommission := &v1.ListCompanyOrganizationsReply_OrganizationColonelCommission{
			ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
			PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
			PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
			IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
			IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
			AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
			CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
			PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
			PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
			IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
			IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
			AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
			OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
			IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
			IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
			AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
			CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
			PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
			PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
			IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
			IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
			AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
		}

		organizationCourses := make([]*v1.ListCompanyOrganizationsReply_OrganizationCourse, 0)

		for _, organizationCourse := range companyOrganization.OrganizationCourses {
			courseModules := make([]*v1.ListCompanyOrganizationsReply_CourseModule, 0)

			for _, courseModule := range organizationCourse.CourseModules {
				courseModules = append(courseModules, &v1.ListCompanyOrganizationsReply_CourseModule{
					CourseModuleName:    courseModule.CourseModuleName,
					CourseModuleContent: courseModule.CourseModuleContent,
				})
			}

			organizationCourses = append(organizationCourses, &v1.ListCompanyOrganizationsReply_OrganizationCourse{
				CourseName:          organizationCourse.CourseName,
				CourseSubName:       organizationCourse.CourseSubName,
				CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
				CourseDuration:      organizationCourse.CourseDuration,
				CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
				CourseLevel:         uint32(organizationCourse.CourseLevel),
				CourseModules:       courseModules,
			})
		}

		organizationUsers := make([]*v1.ListCompanyOrganizationsReply_OrganizationUser, 0)

		for _, organizationUser := range companyOrganization.OrganizationUsers {
			organizationUsers = append(organizationUsers, &v1.ListCompanyOrganizationsReply_OrganizationUser{
				UserId:   organizationUser.UserId,
				Username: organizationUser.Username,
				Phone:    organizationUser.Phone,
			})
		}

		list = append(list, &v1.ListCompanyOrganizationsReply_CompanyOrganization{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		})
	}

	totalPage := uint64(math.Ceil(float64(companyOrganizations.Total) / float64(companyOrganizations.PageSize)))

	return &v1.ListCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.ListCompanyOrganizationsReply_Data{
			PageNum:   companyOrganizations.PageNum,
			PageSize:  companyOrganizations.PageSize,
			Total:     companyOrganizations.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) ListSelectCompanyOrganizations(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCompanyOrganizationsReply, error) {
	selects, err := cs.coruc.ListSelectCompanyOrganizations(ctx)

	if err != nil {
		return nil, err
	}

	courseLevel := make([]*v1.ListSelectCompanyOrganizationsReply_CourseLevel, 0)

	for _, lcourseLevel := range selects.CourseLevel {
		courseLevel = append(courseLevel, &v1.ListSelectCompanyOrganizationsReply_CourseLevel{
			Key:   lcourseLevel.Key,
			Value: lcourseLevel.Value,
		})
	}

	organizationMcn := make([]*v1.ListSelectCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, lorganizationMcn := range selects.OrganizationMcn {
		organizationMcn = append(organizationMcn, &v1.ListSelectCompanyOrganizationsReply_OrganizationMcn{
			Key:   lorganizationMcn.Key,
			Value: lorganizationMcn.Value,
		})
	}

	return &v1.ListSelectCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.ListSelectCompanyOrganizationsReply_Data{
			CourseLevel:     courseLevel,
			OrganizationMcn: organizationMcn,
		},
	}, nil
}

func (cs *CompanyService) CreateCompanyOrganizations(ctx context.Context, in *v1.CreateCompanyOrganizationsRequest) (*v1.CreateCompanyOrganizationsReply, error) {
	companyOrganization, err := cs.coruc.CreateCompanyOrganizations(ctx, in.OrganizationName, in.OrganizationLogo, in.OrganizationMcn, in.CompanyName, in.BankCode, in.BankDeposit)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.CreateCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.CreateCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn,
		})
	}

	organizationCommission := &v1.CreateCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
		OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
	}

	organizationColonelCommission := &v1.CreateCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
		PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
		IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
		IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
		AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
		CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
		PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
		PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
		IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
		IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
		AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
		OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
		IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
		IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
		AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
		CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
		PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
		PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
		IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
		IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
		AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
	}

	organizationCourses := make([]*v1.CreateCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.OrganizationCourses {
		courseModules := make([]*v1.CreateCompanyOrganizationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.CreateCompanyOrganizationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.CreateCompanyOrganizationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.CreateCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.CreateCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.CreateCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.CreateCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (cs *CompanyService) UpdateCompanyOrganizations(ctx context.Context, in *v1.UpdateCompanyOrganizationsRequest) (*v1.UpdateCompanyOrganizationsReply, error) {
	companyOrganization, err := cs.coruc.UpdateCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationName, in.OrganizationLogo, in.OrganizationMcn, in.CompanyName, in.BankCode, in.BankDeposit)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn,
		})
	}

	organizationCommission := &v1.UpdateCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
		OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
	}

	organizationColonelCommission := &v1.UpdateCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
		PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
		IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
		IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
		AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
		CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
		PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
		PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
		IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
		IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
		AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
		OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
		IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
		IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
		AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
		CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
		PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
		PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
		IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
		IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
		AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
	}

	organizationCourses := make([]*v1.UpdateCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.OrganizationCourses {
		courseModules := make([]*v1.UpdateCompanyOrganizationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.UpdateCompanyOrganizationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.UpdateCompanyOrganizationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (cs *CompanyService) UpdateTeamCompanyOrganizations(ctx context.Context, in *v1.UpdateTeamCompanyOrganizationsRequest) (*v1.UpdateTeamCompanyOrganizationsReply, error) {
	companyOrganization, err := cs.coruc.UpdateTeamCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationUser)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateTeamCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateTeamCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn,
		})
	}

	organizationCommission := &v1.UpdateTeamCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
		OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
	}

	organizationColonelCommission := &v1.UpdateTeamCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
		PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
		IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
		IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
		AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
		CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
		PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
		PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
		IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
		IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
		AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
		OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
		IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
		IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
		AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
		CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
		PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
		PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
		IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
		IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
		AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
	}

	organizationCourses := make([]*v1.UpdateTeamCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.OrganizationCourses {
		courseModules := make([]*v1.UpdateTeamCompanyOrganizationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.UpdateTeamCompanyOrganizationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.UpdateTeamCompanyOrganizationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateTeamCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateTeamCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateTeamCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateTeamCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (cs *CompanyService) UpdateCommissionCompanyOrganizations(ctx context.Context, in *v1.UpdateCommissionCompanyOrganizationsRequest) (*v1.UpdateCommissionCompanyOrganizationsReply, error) {
	companyOrganization, err := cs.coruc.UpdateCommissionCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationCommission)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateCommissionCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn,
		})
	}

	organizationCommission := &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
		OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
	}

	organizationColonelCommission := &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
		PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
		IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
		IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
		AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
		CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
		PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
		PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
		IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
		IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
		AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
		OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
		IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
		IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
		AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
		CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
		PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
		PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
		IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
		IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
		AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
	}

	organizationCourses := make([]*v1.UpdateCommissionCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.OrganizationCourses {
		courseModules := make([]*v1.UpdateCommissionCompanyOrganizationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.UpdateCommissionCompanyOrganizationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateCommissionCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateCommissionCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateCommissionCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (cs *CompanyService) UpdateColonelCommissionCompanyOrganizations(ctx context.Context, in *v1.UpdateColonelCommissionCompanyOrganizationsRequest) (*v1.UpdateColonelCommissionCompanyOrganizationsReply, error) {
	companyOrganization, err := cs.coruc.UpdateColonelCommissionCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationColonelCommission)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn,
		})
	}

	organizationCommission := &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
		OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
	}

	organizationColonelCommission := &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
		PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
		IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
		IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
		AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
		CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
		PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
		PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
		IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
		IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
		AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
		OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
		IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
		IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
		AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
		CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
		PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
		PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
		IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
		IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
		AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
	}

	organizationCourses := make([]*v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.OrganizationCourses {
		courseModules := make([]*v1.UpdateColonelCommissionCompanyOrganizationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.UpdateColonelCommissionCompanyOrganizationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateColonelCommissionCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateColonelCommissionCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (cs *CompanyService) UpdateCourseCompanyOrganizations(ctx context.Context, in *v1.UpdateCourseCompanyOrganizationsRequest) (*v1.UpdateCourseCompanyOrganizationsReply, error) {
	companyOrganization, err := cs.coruc.UpdateCourseCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationCourse)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateCourseCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateCourseCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn,
		})
	}

	organizationCommission := &v1.UpdateCourseCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: tool.Decimal(float64(companyOrganization.OrganizationCommissions.CostOrderCommissionRatio), 2),
		OrderCommissionRatio:     tool.Decimal(float64(companyOrganization.OrganizationCommissions.OrderCommissionRatio), 2),
	}

	organizationColonelCommission := &v1.UpdateCourseCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.ZeroCourseRatio), 2),
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterZeroCourseCommissionRule), 2),
		PrimaryAdvancedTutorZeroCourseCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorZeroCourseCommissionRule), 2),
		IntermediateAdvancedPresenterZeroCourseCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterZeroCourseCommissionRule), 2),
		IntermediateAdvancedTutorZeroCourseCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorZeroCourseCommissionRule), 2),
		AdvancedPresenterZeroCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterZeroCourseCommissionRule), 2),
		CourseRatio: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CourseRatio), 2),
		PrimaryAdvancedPresenterCourseCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCourseCommissionRule), 2),
		PrimaryAdvancedTutorCourseCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCourseCommissionRule), 2),
		IntermediateAdvancedPresenterCourseCommissionRule:    tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCourseCommissionRule), 2),
		IntermediateAdvancedTutorCourseCommissionRule:        tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCourseCommissionRule), 2),
		AdvancedPresenterCourseCommissionRule:                tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCourseCommissionRule), 2),
		OrderRatio:                                           tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.OrderRatio), 2),
		IntermediateAdvancedPresenterOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterOrderCommissionRule), 2),
		IntermediateAdvancedTutorOrderCommissionRule:         tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorOrderCommissionRule), 2),
		AdvancedPresenterOrderCommissionRule:                 tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterOrderCommissionRule), 2),
		CostOrderRatio:                                       tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.CostOrderRatio), 2),
		PrimaryAdvancedPresenterCostOrderCommissionRule:      tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedPresenterCostOrderCommissionRule), 2),
		PrimaryAdvancedTutorCostOrderCommissionRule:          tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.PrimaryAdvancedTutorCostOrderCommissionRule), 2),
		IntermediateAdvancedPresenterCostOrderCommissionRule: tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedPresenterCostOrderCommissionRule), 2),
		IntermediateAdvancedTutorCostOrderCommissionRule:     tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.IntermediateAdvancedTutorCostOrderCommissionRule), 2),
		AdvancedPresenterCostOrderCommissionRule:             tool.Decimal(float64(companyOrganization.OrganizationColonelCommissions.AdvancedPresenterCostOrderCommissionRule), 2),
	}

	organizationCourses := make([]*v1.UpdateCourseCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.OrganizationCourses {
		courseModules := make([]*v1.UpdateCourseCompanyOrganizationsReply_CourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, &v1.UpdateCourseCompanyOrganizationsReply_CourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, &v1.UpdateCourseCompanyOrganizationsReply_OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         tool.Decimal(float64(organizationCourse.CoursePrice), 2),
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: tool.Decimal(float64(organizationCourse.CourseOriginalPrice), 2),
			CourseLevel:         uint32(organizationCourse.CourseLevel),
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateCourseCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateCourseCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateCourseCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateCourseCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Id,
			OrganizationName:              companyOrganization.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.CompanyName,
			BankCode:                      companyOrganization.BankCode,
			BankDeposit:                   companyOrganization.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (cs *CompanyService) SyncUpdateQrCodeCompanyOrganizations(ctx context.Context, in *emptypb.Empty) (*v1.SyncUpdateQrCodeCompanyOrganizationsReply, error) {
	ctx = context.Background()

	if err := cs.coruc.SyncUpdateQrCodeCompanyOrganizations(ctx); err != nil {
		return nil, err
	}

	return &v1.SyncUpdateQrCodeCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.SyncUpdateQrCodeCompanyOrganizationsReply_Data{},
	}, nil
}
