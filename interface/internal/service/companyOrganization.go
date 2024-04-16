package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) ListCompanyOrganizations(ctx context.Context, in *v1.ListCompanyOrganizationsRequest) (*v1.ListCompanyOrganizationsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "companyOrganization"); err != nil {
		return nil, err
	}

	companyOrganizations, err := is.coruc.ListCompanyOrganizations(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyOrganizationsReply_CompanyOrganization, 0)

	for _, companyOrganization := range companyOrganizations.Data.List {
		organizationMcns := make([]*v1.ListCompanyOrganizationsReply_OrganizationMcn, 0)

		for _, organizationMcn := range companyOrganization.OrganizationMcns {
			organizationMcns = append(organizationMcns, &v1.ListCompanyOrganizationsReply_OrganizationMcn{
				OrganizationMcn: organizationMcn.OrganizationMcn,
			})
		}

		organizationCommission := &v1.ListCompanyOrganizationsReply_OrganizationCommission{
			CostOrderCommissionRatio: companyOrganization.OrganizationCommission.CostOrderCommissionRatio,
			OrderCommissionRatio:     companyOrganization.OrganizationCommission.OrderCommissionRatio,
		}

		organizationColonelCommission := &v1.ListCompanyOrganizationsReply_OrganizationColonelCommission{
			ZeroCourseRatio: companyOrganization.OrganizationColonelCommission.ZeroCourseRatio,
			ZeroAdvancedPresenterZeroCourseCommissionRule:         companyOrganization.OrganizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule,
			ZeroAdvancedTutorZeroCourseCommissionRule:             companyOrganization.OrganizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule,
			PrimaryAdvancedPresenterZeroCourseCommissionRule:      companyOrganization.OrganizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule,
			PrimaryAdvancedTutorZeroCourseCommissionRule:          companyOrganization.OrganizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule,
			IntermediateAdvancedPresenterZeroCourseCommissionRule: companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule,
			IntermediateAdvancedTutorZeroCourseCommissionRule:     companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule,
			AdvancedPresenterZeroCourseCommissionRule:             companyOrganization.OrganizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule,
			PrimaryCourseRatio: companyOrganization.OrganizationColonelCommission.PrimaryCourseRatio,
			PrimaryAdvancedPresenterPrimaryCourseCommissionRule:           companyOrganization.OrganizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule,
			PrimaryAdvancedTutorPrimaryCourseCommissionRule:               companyOrganization.OrganizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule,
			IntermediateAdvancedPresenterPrimaryCourseCommissionRule:      companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule,
			IntermediateAdvancedTutorPrimaryCourseCommissionRule:          companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule,
			AdvancedPresenterPrimaryCourseCommissionRule:                  companyOrganization.OrganizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule,
			IntermediateCourseRatio:                                       companyOrganization.OrganizationColonelCommission.IntermediateCourseRatio,
			PrimaryAdvancedPresenterIntermediateCourseCommissionRule:      companyOrganization.OrganizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule,
			PrimaryAdvancedTutorIntermediateCourseCommissionRule:          companyOrganization.OrganizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule,
			IntermediateAdvancedPresenterIntermediateCourseCommissionRule: companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule,
			IntermediateAdvancedTutorIntermediateCourseCommissionRule:     companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule,
			AdvancedPresenterIntermediateCourseCommissionRule:             companyOrganization.OrganizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule,
			AdvancedCourseRatio: companyOrganization.OrganizationColonelCommission.AdvancedCourseRatio,
			PrimaryAdvancedPresenterAdvancedCourseCommissionRule:      companyOrganization.OrganizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule,
			PrimaryAdvancedTutorAdvancedCourseCommissionRule:          companyOrganization.OrganizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule,
			IntermediateAdvancedPresenterAdvancedCourseCommissionRule: companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule,
			IntermediateAdvancedTutorAdvancedCourseCommissionRule:     companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule,
			AdvancedPresenterAdvancedCourseCommissionRule:             companyOrganization.OrganizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule,
			OrderRatio: companyOrganization.OrganizationColonelCommission.OrderRatio,
			PrimaryAdvancedPresenterOrderCommissionRule:          companyOrganization.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule,
			PrimaryAdvancedTutorOrderCommissionRule:              companyOrganization.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule,
			IntermediateAdvancedPresenterOrderCommissionRule:     companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule,
			IntermediateAdvancedTutorOrderCommissionRule:         companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule,
			AdvancedPresenterOrderCommissionRule:                 companyOrganization.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule,
			CostOrderRatio:                                       companyOrganization.OrganizationColonelCommission.CostOrderRatio,
			ZeroAdvancedPresenterCostOrderCommissionRule:         companyOrganization.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule,
			ZeroAdvancedTutorCostOrderCommissionRule:             companyOrganization.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule,
			PrimaryAdvancedPresenterCostOrderCommissionRule:      companyOrganization.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule,
			PrimaryAdvancedTutorCostOrderCommissionRule:          companyOrganization.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule,
			IntermediateAdvancedPresenterCostOrderCommissionRule: companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule,
			IntermediateAdvancedTutorCostOrderCommissionRule:     companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule,
			AdvancedPresenterCostOrderCommissionRule:             companyOrganization.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule,
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
				CoursePrice:         float64(organizationCourse.CoursePrice),
				CourseDuration:      organizationCourse.CourseDuration,
				CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
				CourseLevel:         organizationCourse.CourseLevel,
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
			OrganizationId:                companyOrganization.OrganizationId,
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

	return &v1.ListCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.ListCompanyOrganizationsReply_Data{
			PageNum:   companyOrganizations.Data.PageNum,
			PageSize:  companyOrganizations.Data.PageSize,
			Total:     companyOrganizations.Data.Total,
			TotalPage: companyOrganizations.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListSelectCompanyOrganizations(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCompanyOrganizationsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "companyOrganization"); err != nil {
		return nil, err
	}

	selects, err := is.coruc.ListSelectCompanyOrganizations(ctx)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	courseLevel := make([]*v1.ListSelectCompanyOrganizationsReply_CourseLevel, 0)

	for _, lcourseLevel := range selects.Data.CourseLevel {
		courseLevel = append(courseLevel, &v1.ListSelectCompanyOrganizationsReply_CourseLevel{
			Key:   lcourseLevel.Key,
			Value: lcourseLevel.Value,
		})
	}

	organizationMcn := make([]*v1.ListSelectCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, lorganizationMcn := range selects.Data.OrganizationMcn {
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

func (is *InterfaceService) CreateCompanyOrganizations(ctx context.Context, in *v1.CreateCompanyOrganizationsRequest) (*v1.CreateCompanyOrganizationsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.verifyLogin(ctx, true, false, "companyOrganization"); err != nil {
		return nil, err
	}

	companyOrganization, err := is.coruc.CreateCompanyOrganizations(ctx, in.OrganizationName, in.OrganizationLogo, in.OrganizationMcn, in.CompanyName, in.BankCode, in.BankDeposit)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.CreateCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.Data.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.CreateCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn.OrganizationMcn,
		})
	}

	organizationCommission := &v1.CreateCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: companyOrganization.Data.OrganizationCommission.CostOrderCommissionRatio,
		OrderCommissionRatio:     companyOrganization.Data.OrganizationCommission.OrderCommissionRatio,
	}

	organizationColonelCommission := &v1.CreateCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: companyOrganization.Data.OrganizationColonelCommission.ZeroCourseRatio,
		ZeroAdvancedPresenterZeroCourseCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule,
		ZeroAdvancedTutorZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule,
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule,
		PrimaryAdvancedTutorZeroCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule,
		IntermediateAdvancedPresenterZeroCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule,
		IntermediateAdvancedTutorZeroCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule,
		AdvancedPresenterZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule,
		PrimaryCourseRatio: companyOrganization.Data.OrganizationColonelCommission.PrimaryCourseRatio,
		PrimaryAdvancedPresenterPrimaryCourseCommissionRule:           companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule,
		PrimaryAdvancedTutorPrimaryCourseCommissionRule:               companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule,
		IntermediateAdvancedPresenterPrimaryCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateAdvancedTutorPrimaryCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule,
		AdvancedPresenterPrimaryCourseCommissionRule:                  companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateCourseRatio:                                       companyOrganization.Data.OrganizationColonelCommission.IntermediateCourseRatio,
		PrimaryAdvancedPresenterIntermediateCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule,
		PrimaryAdvancedTutorIntermediateCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule,
		IntermediateAdvancedPresenterIntermediateCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule,
		IntermediateAdvancedTutorIntermediateCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule,
		AdvancedPresenterIntermediateCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule,
		AdvancedCourseRatio: companyOrganization.Data.OrganizationColonelCommission.AdvancedCourseRatio,
		PrimaryAdvancedPresenterAdvancedCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule,
		PrimaryAdvancedTutorAdvancedCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule,
		IntermediateAdvancedPresenterAdvancedCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule,
		IntermediateAdvancedTutorAdvancedCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule,
		AdvancedPresenterAdvancedCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule,
		OrderRatio: companyOrganization.Data.OrganizationColonelCommission.OrderRatio,
		PrimaryAdvancedPresenterOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule,
		PrimaryAdvancedTutorOrderCommissionRule:              companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule,
		IntermediateAdvancedPresenterOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule,
		IntermediateAdvancedTutorOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule,
		AdvancedPresenterOrderCommissionRule:                 companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule,
		CostOrderRatio:                                       companyOrganization.Data.OrganizationColonelCommission.CostOrderRatio,
		ZeroAdvancedPresenterCostOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule,
		ZeroAdvancedTutorCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule,
		PrimaryAdvancedPresenterCostOrderCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule,
		PrimaryAdvancedTutorCostOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule,
		IntermediateAdvancedPresenterCostOrderCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule,
		IntermediateAdvancedTutorCostOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule,
		AdvancedPresenterCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule,
	}

	organizationCourses := make([]*v1.CreateCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.Data.OrganizationCourses {
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
			CoursePrice:         organizationCourse.CoursePrice,
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
			CourseLevel:         organizationCourse.CourseLevel,
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.CreateCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.Data.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.CreateCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.CreateCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.CreateCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Data.OrganizationId,
			OrganizationName:              companyOrganization.Data.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.Data.CompanyName,
			BankCode:                      companyOrganization.Data.BankCode,
			BankDeposit:                   companyOrganization.Data.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.Data.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.Data.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.Data.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.Data.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (is *InterfaceService) UpdateCompanyOrganizations(ctx context.Context, in *v1.UpdateCompanyOrganizationsRequest) (*v1.UpdateCompanyOrganizationsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "companyOrganization"); err != nil {
		return nil, err
	}

	companyOrganization, err := is.coruc.UpdateCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationName, in.OrganizationLogo, in.OrganizationMcn, in.CompanyName, in.BankCode, in.BankDeposit)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.Data.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn.OrganizationMcn,
		})
	}

	organizationCommission := &v1.UpdateCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: companyOrganization.Data.OrganizationCommission.CostOrderCommissionRatio,
		OrderCommissionRatio:     companyOrganization.Data.OrganizationCommission.OrderCommissionRatio,
	}

	organizationColonelCommission := &v1.UpdateCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: companyOrganization.Data.OrganizationColonelCommission.ZeroCourseRatio,
		ZeroAdvancedPresenterZeroCourseCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule,
		ZeroAdvancedTutorZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule,
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule,
		PrimaryAdvancedTutorZeroCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule,
		IntermediateAdvancedPresenterZeroCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule,
		IntermediateAdvancedTutorZeroCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule,
		AdvancedPresenterZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule,
		PrimaryCourseRatio: companyOrganization.Data.OrganizationColonelCommission.PrimaryCourseRatio,
		PrimaryAdvancedPresenterPrimaryCourseCommissionRule:           companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule,
		PrimaryAdvancedTutorPrimaryCourseCommissionRule:               companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule,
		IntermediateAdvancedPresenterPrimaryCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateAdvancedTutorPrimaryCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule,
		AdvancedPresenterPrimaryCourseCommissionRule:                  companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateCourseRatio:                                       companyOrganization.Data.OrganizationColonelCommission.IntermediateCourseRatio,
		PrimaryAdvancedPresenterIntermediateCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule,
		PrimaryAdvancedTutorIntermediateCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule,
		IntermediateAdvancedPresenterIntermediateCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule,
		IntermediateAdvancedTutorIntermediateCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule,
		AdvancedPresenterIntermediateCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule,
		AdvancedCourseRatio: companyOrganization.Data.OrganizationColonelCommission.AdvancedCourseRatio,
		PrimaryAdvancedPresenterAdvancedCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule,
		PrimaryAdvancedTutorAdvancedCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule,
		IntermediateAdvancedPresenterAdvancedCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule,
		IntermediateAdvancedTutorAdvancedCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule,
		AdvancedPresenterAdvancedCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule,
		OrderRatio: companyOrganization.Data.OrganizationColonelCommission.OrderRatio,
		PrimaryAdvancedPresenterOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule,
		PrimaryAdvancedTutorOrderCommissionRule:              companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule,
		IntermediateAdvancedPresenterOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule,
		IntermediateAdvancedTutorOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule,
		AdvancedPresenterOrderCommissionRule:                 companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule,
		CostOrderRatio:                                       companyOrganization.Data.OrganizationColonelCommission.CostOrderRatio,
		ZeroAdvancedPresenterCostOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule,
		ZeroAdvancedTutorCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule,
		PrimaryAdvancedPresenterCostOrderCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule,
		PrimaryAdvancedTutorCostOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule,
		IntermediateAdvancedPresenterCostOrderCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule,
		IntermediateAdvancedTutorCostOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule,
		AdvancedPresenterCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule,
	}

	organizationCourses := make([]*v1.UpdateCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.Data.OrganizationCourses {
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
			CoursePrice:         organizationCourse.CoursePrice,
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
			CourseLevel:         organizationCourse.CourseLevel,
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.Data.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Data.OrganizationId,
			OrganizationName:              companyOrganization.Data.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.Data.CompanyName,
			BankCode:                      companyOrganization.Data.BankCode,
			BankDeposit:                   companyOrganization.Data.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.Data.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.Data.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.Data.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.Data.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (is *InterfaceService) UpdateTeamCompanyOrganizations(ctx context.Context, in *v1.UpdateTeamCompanyOrganizationsRequest) (*v1.UpdateTeamCompanyOrganizationsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "companyOrganization"); err != nil {
		return nil, err
	}

	companyOrganization, err := is.coruc.UpdateTeamCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationUser)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateTeamCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.Data.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateTeamCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn.OrganizationMcn,
		})
	}

	organizationCommission := &v1.UpdateTeamCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: companyOrganization.Data.OrganizationCommission.CostOrderCommissionRatio,
		OrderCommissionRatio:     companyOrganization.Data.OrganizationCommission.OrderCommissionRatio,
	}

	organizationColonelCommission := &v1.UpdateTeamCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: companyOrganization.Data.OrganizationColonelCommission.ZeroCourseRatio,
		ZeroAdvancedPresenterZeroCourseCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule,
		ZeroAdvancedTutorZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule,
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule,
		PrimaryAdvancedTutorZeroCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule,
		IntermediateAdvancedPresenterZeroCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule,
		IntermediateAdvancedTutorZeroCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule,
		AdvancedPresenterZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule,
		PrimaryCourseRatio: companyOrganization.Data.OrganizationColonelCommission.PrimaryCourseRatio,
		PrimaryAdvancedPresenterPrimaryCourseCommissionRule:           companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule,
		PrimaryAdvancedTutorPrimaryCourseCommissionRule:               companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule,
		IntermediateAdvancedPresenterPrimaryCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateAdvancedTutorPrimaryCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule,
		AdvancedPresenterPrimaryCourseCommissionRule:                  companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateCourseRatio:                                       companyOrganization.Data.OrganizationColonelCommission.IntermediateCourseRatio,
		PrimaryAdvancedPresenterIntermediateCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule,
		PrimaryAdvancedTutorIntermediateCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule,
		IntermediateAdvancedPresenterIntermediateCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule,
		IntermediateAdvancedTutorIntermediateCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule,
		AdvancedPresenterIntermediateCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule,
		AdvancedCourseRatio: companyOrganization.Data.OrganizationColonelCommission.AdvancedCourseRatio,
		PrimaryAdvancedPresenterAdvancedCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule,
		PrimaryAdvancedTutorAdvancedCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule,
		IntermediateAdvancedPresenterAdvancedCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule,
		IntermediateAdvancedTutorAdvancedCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule,
		AdvancedPresenterAdvancedCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule,
		OrderRatio: companyOrganization.Data.OrganizationColonelCommission.OrderRatio,
		PrimaryAdvancedPresenterOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule,
		PrimaryAdvancedTutorOrderCommissionRule:              companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule,
		IntermediateAdvancedPresenterOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule,
		IntermediateAdvancedTutorOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule,
		AdvancedPresenterOrderCommissionRule:                 companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule,
		CostOrderRatio:                                       companyOrganization.Data.OrganizationColonelCommission.CostOrderRatio,
		ZeroAdvancedPresenterCostOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule,
		ZeroAdvancedTutorCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule,
		PrimaryAdvancedPresenterCostOrderCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule,
		PrimaryAdvancedTutorCostOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule,
		IntermediateAdvancedPresenterCostOrderCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule,
		IntermediateAdvancedTutorCostOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule,
		AdvancedPresenterCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule,
	}

	organizationCourses := make([]*v1.UpdateTeamCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.Data.OrganizationCourses {
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
			CoursePrice:         organizationCourse.CoursePrice,
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
			CourseLevel:         organizationCourse.CourseLevel,
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateTeamCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.Data.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateTeamCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateTeamCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateTeamCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Data.OrganizationId,
			OrganizationName:              companyOrganization.Data.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.Data.CompanyName,
			BankCode:                      companyOrganization.Data.BankCode,
			BankDeposit:                   companyOrganization.Data.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.Data.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.Data.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.Data.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.Data.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (is *InterfaceService) UpdateCommissionCompanyOrganizations(ctx context.Context, in *v1.UpdateCommissionCompanyOrganizationsRequest) (*v1.UpdateCommissionCompanyOrganizationsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "companyOrganization"); err != nil {
		return nil, err
	}

	companyOrganization, err := is.coruc.UpdateCommissionCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationCommission)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateCommissionCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.Data.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn.OrganizationMcn,
		})
	}

	organizationCommission := &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: companyOrganization.Data.OrganizationCommission.CostOrderCommissionRatio,
		OrderCommissionRatio:     companyOrganization.Data.OrganizationCommission.OrderCommissionRatio,
	}

	organizationColonelCommission := &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: companyOrganization.Data.OrganizationColonelCommission.ZeroCourseRatio,
		ZeroAdvancedPresenterZeroCourseCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule,
		ZeroAdvancedTutorZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule,
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule,
		PrimaryAdvancedTutorZeroCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule,
		IntermediateAdvancedPresenterZeroCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule,
		IntermediateAdvancedTutorZeroCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule,
		AdvancedPresenterZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule,
		PrimaryCourseRatio: companyOrganization.Data.OrganizationColonelCommission.PrimaryCourseRatio,
		PrimaryAdvancedPresenterPrimaryCourseCommissionRule:           companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule,
		PrimaryAdvancedTutorPrimaryCourseCommissionRule:               companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule,
		IntermediateAdvancedPresenterPrimaryCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateAdvancedTutorPrimaryCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule,
		AdvancedPresenterPrimaryCourseCommissionRule:                  companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateCourseRatio:                                       companyOrganization.Data.OrganizationColonelCommission.IntermediateCourseRatio,
		PrimaryAdvancedPresenterIntermediateCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule,
		PrimaryAdvancedTutorIntermediateCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule,
		IntermediateAdvancedPresenterIntermediateCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule,
		IntermediateAdvancedTutorIntermediateCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule,
		AdvancedPresenterIntermediateCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule,
		AdvancedCourseRatio: companyOrganization.Data.OrganizationColonelCommission.AdvancedCourseRatio,
		PrimaryAdvancedPresenterAdvancedCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule,
		PrimaryAdvancedTutorAdvancedCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule,
		IntermediateAdvancedPresenterAdvancedCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule,
		IntermediateAdvancedTutorAdvancedCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule,
		AdvancedPresenterAdvancedCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule,
		OrderRatio: companyOrganization.Data.OrganizationColonelCommission.OrderRatio,
		PrimaryAdvancedPresenterOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule,
		PrimaryAdvancedTutorOrderCommissionRule:              companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule,
		IntermediateAdvancedPresenterOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule,
		IntermediateAdvancedTutorOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule,
		AdvancedPresenterOrderCommissionRule:                 companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule,
		CostOrderRatio:                                       companyOrganization.Data.OrganizationColonelCommission.CostOrderRatio,
		ZeroAdvancedPresenterCostOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule,
		ZeroAdvancedTutorCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule,
		PrimaryAdvancedPresenterCostOrderCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule,
		PrimaryAdvancedTutorCostOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule,
		IntermediateAdvancedPresenterCostOrderCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule,
		IntermediateAdvancedTutorCostOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule,
		AdvancedPresenterCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule,
	}

	organizationCourses := make([]*v1.UpdateCommissionCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.Data.OrganizationCourses {
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
			CoursePrice:         organizationCourse.CoursePrice,
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
			CourseLevel:         organizationCourse.CourseLevel,
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateCommissionCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.Data.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateCommissionCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateCommissionCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateCommissionCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Data.OrganizationId,
			OrganizationName:              companyOrganization.Data.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.Data.CompanyName,
			BankCode:                      companyOrganization.Data.BankCode,
			BankDeposit:                   companyOrganization.Data.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.Data.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.Data.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.Data.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.Data.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (is *InterfaceService) UpdateColonelCommissionCompanyOrganizations(ctx context.Context, in *v1.UpdateColonelCommissionCompanyOrganizationsRequest) (*v1.UpdateColonelCommissionCompanyOrganizationsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	companyOrganization, err := is.coruc.UpdateColonelCommissionCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationColonelCommission)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.Data.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn.OrganizationMcn,
		})
	}

	organizationCommission := &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: companyOrganization.Data.OrganizationCommission.CostOrderCommissionRatio,
		OrderCommissionRatio:     companyOrganization.Data.OrganizationCommission.OrderCommissionRatio,
	}

	organizationColonelCommission := &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: companyOrganization.Data.OrganizationColonelCommission.ZeroCourseRatio,
		ZeroAdvancedPresenterZeroCourseCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule,
		ZeroAdvancedTutorZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule,
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule,
		PrimaryAdvancedTutorZeroCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule,
		IntermediateAdvancedPresenterZeroCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule,
		IntermediateAdvancedTutorZeroCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule,
		AdvancedPresenterZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule,
		PrimaryCourseRatio: companyOrganization.Data.OrganizationColonelCommission.PrimaryCourseRatio,
		PrimaryAdvancedPresenterPrimaryCourseCommissionRule:           companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule,
		PrimaryAdvancedTutorPrimaryCourseCommissionRule:               companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule,
		IntermediateAdvancedPresenterPrimaryCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateAdvancedTutorPrimaryCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule,
		AdvancedPresenterPrimaryCourseCommissionRule:                  companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateCourseRatio:                                       companyOrganization.Data.OrganizationColonelCommission.IntermediateCourseRatio,
		PrimaryAdvancedPresenterIntermediateCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule,
		PrimaryAdvancedTutorIntermediateCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule,
		IntermediateAdvancedPresenterIntermediateCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule,
		IntermediateAdvancedTutorIntermediateCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule,
		AdvancedPresenterIntermediateCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule,
		AdvancedCourseRatio: companyOrganization.Data.OrganizationColonelCommission.AdvancedCourseRatio,
		PrimaryAdvancedPresenterAdvancedCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule,
		PrimaryAdvancedTutorAdvancedCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule,
		IntermediateAdvancedPresenterAdvancedCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule,
		IntermediateAdvancedTutorAdvancedCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule,
		AdvancedPresenterAdvancedCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule,
		OrderRatio: companyOrganization.Data.OrganizationColonelCommission.OrderRatio,
		PrimaryAdvancedPresenterOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule,
		PrimaryAdvancedTutorOrderCommissionRule:              companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule,
		IntermediateAdvancedPresenterOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule,
		IntermediateAdvancedTutorOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule,
		AdvancedPresenterOrderCommissionRule:                 companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule,
		CostOrderRatio:                                       companyOrganization.Data.OrganizationColonelCommission.CostOrderRatio,
		ZeroAdvancedPresenterCostOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule,
		ZeroAdvancedTutorCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule,
		PrimaryAdvancedPresenterCostOrderCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule,
		PrimaryAdvancedTutorCostOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule,
		IntermediateAdvancedPresenterCostOrderCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule,
		IntermediateAdvancedTutorCostOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule,
		AdvancedPresenterCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule,
	}

	organizationCourses := make([]*v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.Data.OrganizationCourses {
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
			CoursePrice:         organizationCourse.CoursePrice,
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
			CourseLevel:         organizationCourse.CourseLevel,
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.Data.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateColonelCommissionCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateColonelCommissionCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateColonelCommissionCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Data.OrganizationId,
			OrganizationName:              companyOrganization.Data.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.Data.CompanyName,
			BankCode:                      companyOrganization.Data.BankCode,
			BankDeposit:                   companyOrganization.Data.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.Data.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.Data.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.Data.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.Data.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}

func (is *InterfaceService) UpdateCourseCompanyOrganizations(ctx context.Context, in *v1.UpdateCourseCompanyOrganizationsRequest) (*v1.UpdateCourseCompanyOrganizationsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	companyOrganization, err := is.coruc.UpdateCourseCompanyOrganizations(ctx, in.OrganizationId, in.OrganizationCourse)

	if err != nil {
		return nil, err
	}

	organizationMcns := make([]*v1.UpdateCourseCompanyOrganizationsReply_OrganizationMcn, 0)

	for _, organizationMcn := range companyOrganization.Data.OrganizationMcns {
		organizationMcns = append(organizationMcns, &v1.UpdateCourseCompanyOrganizationsReply_OrganizationMcn{
			OrganizationMcn: organizationMcn.OrganizationMcn,
		})
	}

	organizationCommission := &v1.UpdateCourseCompanyOrganizationsReply_OrganizationCommission{
		CostOrderCommissionRatio: companyOrganization.Data.OrganizationCommission.CostOrderCommissionRatio,
		OrderCommissionRatio:     companyOrganization.Data.OrganizationCommission.OrderCommissionRatio,
	}

	organizationColonelCommission := &v1.UpdateCourseCompanyOrganizationsReply_OrganizationColonelCommission{
		ZeroCourseRatio: companyOrganization.Data.OrganizationColonelCommission.ZeroCourseRatio,
		ZeroAdvancedPresenterZeroCourseCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule,
		ZeroAdvancedTutorZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule,
		PrimaryAdvancedPresenterZeroCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule,
		PrimaryAdvancedTutorZeroCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule,
		IntermediateAdvancedPresenterZeroCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule,
		IntermediateAdvancedTutorZeroCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule,
		AdvancedPresenterZeroCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule,
		PrimaryCourseRatio: companyOrganization.Data.OrganizationColonelCommission.PrimaryCourseRatio,
		PrimaryAdvancedPresenterPrimaryCourseCommissionRule:           companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule,
		PrimaryAdvancedTutorPrimaryCourseCommissionRule:               companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule,
		IntermediateAdvancedPresenterPrimaryCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateAdvancedTutorPrimaryCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule,
		AdvancedPresenterPrimaryCourseCommissionRule:                  companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule,
		IntermediateCourseRatio:                                       companyOrganization.Data.OrganizationColonelCommission.IntermediateCourseRatio,
		PrimaryAdvancedPresenterIntermediateCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule,
		PrimaryAdvancedTutorIntermediateCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule,
		IntermediateAdvancedPresenterIntermediateCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule,
		IntermediateAdvancedTutorIntermediateCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule,
		AdvancedPresenterIntermediateCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule,
		AdvancedCourseRatio: companyOrganization.Data.OrganizationColonelCommission.AdvancedCourseRatio,
		PrimaryAdvancedPresenterAdvancedCourseCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule,
		PrimaryAdvancedTutorAdvancedCourseCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule,
		IntermediateAdvancedPresenterAdvancedCourseCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule,
		IntermediateAdvancedTutorAdvancedCourseCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule,
		AdvancedPresenterAdvancedCourseCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule,
		OrderRatio: companyOrganization.Data.OrganizationColonelCommission.OrderRatio,
		PrimaryAdvancedPresenterOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule,
		PrimaryAdvancedTutorOrderCommissionRule:              companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule,
		IntermediateAdvancedPresenterOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule,
		IntermediateAdvancedTutorOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule,
		AdvancedPresenterOrderCommissionRule:                 companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule,
		CostOrderRatio:                                       companyOrganization.Data.OrganizationColonelCommission.CostOrderRatio,
		ZeroAdvancedPresenterCostOrderCommissionRule:         companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule,
		ZeroAdvancedTutorCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule,
		PrimaryAdvancedPresenterCostOrderCommissionRule:      companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule,
		PrimaryAdvancedTutorCostOrderCommissionRule:          companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule,
		IntermediateAdvancedPresenterCostOrderCommissionRule: companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule,
		IntermediateAdvancedTutorCostOrderCommissionRule:     companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule,
		AdvancedPresenterCostOrderCommissionRule:             companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule,
	}

	organizationCourses := make([]*v1.UpdateCourseCompanyOrganizationsReply_OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.Data.OrganizationCourses {
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
			CoursePrice:         organizationCourse.CoursePrice,
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
			CourseLevel:         organizationCourse.CourseLevel,
			CourseModules:       courseModules,
		})
	}

	organizationUsers := make([]*v1.UpdateCourseCompanyOrganizationsReply_OrganizationUser, 0)

	for _, organizationUser := range companyOrganization.Data.OrganizationUsers {
		organizationUsers = append(organizationUsers, &v1.UpdateCourseCompanyOrganizationsReply_OrganizationUser{
			UserId:   organizationUser.UserId,
			Username: organizationUser.Username,
			Phone:    organizationUser.Phone,
		})
	}

	return &v1.UpdateCourseCompanyOrganizationsReply{
		Code: 200,
		Data: &v1.UpdateCourseCompanyOrganizationsReply_Data{
			OrganizationId:                companyOrganization.Data.OrganizationId,
			OrganizationName:              companyOrganization.Data.OrganizationName,
			OrganizationMcns:              organizationMcns,
			CompanyName:                   companyOrganization.Data.CompanyName,
			BankCode:                      companyOrganization.Data.BankCode,
			BankDeposit:                   companyOrganization.Data.BankDeposit,
			OrganizationLogoUrl:           companyOrganization.Data.OrganizationLogoUrl,
			OrganizationCode:              companyOrganization.Data.OrganizationCode,
			OrganizationQrCodeUrl:         companyOrganization.Data.OrganizationQrCodeUrl,
			OrganizationShortUrl:          companyOrganization.Data.OrganizationShortUrl,
			OrganizationCommission:        organizationCommission,
			OrganizationColonelCommission: organizationColonelCommission,
			OrganizationCourses:           organizationCourses,
			OrganizationUsers:             organizationUsers,
		},
	}, nil
}
