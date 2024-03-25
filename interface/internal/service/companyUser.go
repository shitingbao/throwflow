package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) GetCompanyUser(ctx context.Context, in *emptypb.Empty) (*v1.GetCompanyUserReply, error) {
	companyUser, err := is.verifyLogin(ctx, false, false, "")

	if err != nil {
		return nil, err
	}

	userCompanys := make([]*v1.GetCompanyUserReply_Company, 0)

	for _, userCompany := range companyUser.Data.UserCompany {
		userCompanys = append(userCompanys, &v1.GetCompanyUserReply_Company{
			CompanyId:   userCompany.CompanyId,
			CompanyName: userCompany.CompanyName,
		})
	}

	var isMaterial uint32 = 0

	if companyUser.Data.CurrentCompanyId > 0 {
		if companyUserMenus, err := is.cuuc.ListCompanyUserMenu(ctx, companyUser.Data.CurrentCompanyId); err == nil {
			for _, companyUserMenu := range companyUserMenus.Data.List {
				if companyUserMenu.PermissionCode == "material" {
					isMaterial = 1
					break
				}

				for _, childCompanyUserMenu := range companyUserMenu.ChildList {
					if childCompanyUserMenu.PermissionCode == "material" {
						isMaterial = 1
						break
					}
				}
			}
		}
	}

	return &v1.GetCompanyUserReply{
		Code: 200,
		Data: &v1.GetCompanyUserReply_Data{
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
			IsMaterial:           isMaterial,
			Reason:               companyUser.Data.Reason,
			CurrentCompanyId:     companyUser.Data.CurrentCompanyId,
			UserCompany:          userCompanys,
		},
	}, nil
}

func (is *InterfaceService) ChangeCompanyUserCompany(ctx context.Context, in *v1.ChangeCompanyUserCompanyRequest) (*v1.ChangeCompanyUserCompanyReply, error) {
	companyUser, err := is.verifyLogin(ctx, false, false, "")

	if err != nil {
		return nil, err
	}

	isExit := true

	for _, userCompany := range companyUser.Data.UserCompany {
		if userCompany.CompanyId == in.CompanyId {
			isExit = false
		}
	}

	if isExit {
		return nil, errors.InternalServer("INTERFACE_CHANGE_COMPANY_USER_COMPANY_FAILED", "所选企业不存在")
	}

	if err := is.cuuc.ChangeCompanyUserCompany(ctx, in.CompanyId); err != nil {
		return nil, err
	}

	return &v1.ChangeCompanyUserCompanyReply{
		Code: 200,
		Data: &v1.ChangeCompanyUserCompanyReply_Data{},
	}, nil
}

func (is *InterfaceService) ListCompanyUserMenu(ctx context.Context, in *emptypb.Empty) (*v1.ListCompanyUserMenuReply, error) {
	companyUser, err := is.verifyLogin(ctx, false, false, "")

	if err != nil {
		return nil, err
	}

	companyUserMenus, err := is.cuuc.ListCompanyUserMenu(ctx, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserMenuReply_Menus, 0)

	for _, companyUserMenu := range companyUserMenus.Data.List {
		childList := make([]*v1.ListCompanyUserMenuReply_ChildMenus, 0)

		for _, companyUserChildMenu := range companyUserMenu.ChildList {
			childList = append(childList, &v1.ListCompanyUserMenuReply_ChildMenus{
				MenuName:     companyUserChildMenu.MenuName,
				MenuType:     companyUserChildMenu.MenuType,
				Filename:     companyUserChildMenu.Filename,
				Filepath:     companyUserChildMenu.Filepath,
				IsPermission: companyUserChildMenu.IsPermission,
			})
		}

		list = append(list, &v1.ListCompanyUserMenuReply_Menus{
			MenuName:  companyUserMenu.MenuName,
			ChildList: childList,
			MenuType:  companyUserMenu.MenuType,
			Filename:  companyUserMenu.Filename,
			Filepath:  companyUserMenu.Filepath,
		})
	}

	return &v1.ListCompanyUserMenuReply{
		Code: 200,
		Data: &v1.ListCompanyUserMenuReply_Data{
			List: list,
		},
	}, nil
}

func (is *InterfaceService) ListCompanyUserQianchuanAdvertisers(ctx context.Context, in *v1.ListCompanyUserQianchuanAdvertisersRequest) (*v1.ListCompanyUserQianchuanAdvertisersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "account")

	if err != nil {
		return nil, err
	}

	qianchuanAdvertisers, err := is.cuuc.ListCompanyUserQianchuanAdvertisers(ctx, in.PageNum, in.PageSize, companyUser.Data.CurrentCompanyId, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserQianchuanAdvertisersReply_QianchuanAdvertisers, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
		list = append(list, &v1.ListCompanyUserQianchuanAdvertisersReply_QianchuanAdvertisers{
			AdvertiserId:     qianchuanAdvertiser.AdvertiserId,
			AdvertiserName:   qianchuanAdvertiser.AdvertiserName,
			CompanyName:      qianchuanAdvertiser.CompanyName,
			Status:           qianchuanAdvertiser.Status,
			OtherCompanyId:   qianchuanAdvertiser.OtherCompanyId,
			OtherCompanyName: qianchuanAdvertiser.OtherCompanyName,
			UpdateTime:       qianchuanAdvertiser.UpdateTime,
		})
	}

	return &v1.ListCompanyUserQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.ListCompanyUserQianchuanAdvertisersReply_Data{
			PageNum:   qianchuanAdvertisers.Data.PageNum,
			PageSize:  qianchuanAdvertisers.Data.PageSize,
			Total:     qianchuanAdvertisers.Data.Total,
			TotalPage: qianchuanAdvertisers.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) UpdateStatusCompanyUserQianchuanAdvertisers(ctx context.Context, in *v1.UpdateStatusCompanyUserQianchuanAdvertisersRequest) (*v1.UpdateStatusCompanyUserQianchuanAdvertisersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "account")

	if err != nil {
		return nil, err
	}

	qianchuanAdvertiser, err := is.cuuc.UpdateStatusCompanyUserQianchuanAdvertisers(ctx, companyUser.Data.CurrentCompanyId, in.AdvertiserId, in.Status)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusCompanyUserQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.UpdateStatusCompanyUserQianchuanAdvertisersReply_Data{
			AdvertiserId:   qianchuanAdvertiser.Data.AdvertiserId,
			AdvertiserName: qianchuanAdvertiser.Data.AdvertiserName,
			CompanyName:    qianchuanAdvertiser.Data.CompanyName,
			Status:         qianchuanAdvertiser.Data.Status,
		},
	}, nil
}

func (is *InterfaceService) StatisticsQianchuanAdvertisers(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsQianchuanAdvertisersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "account")

	if err != nil {
		return nil, err
	}

	statistics, err := is.cuuc.StatisticsQianchuanAdvertisers(ctx, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsQianchuanAdvertisersReply_Statistics, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsQianchuanAdvertisersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.StatisticsQianchuanAdvertisersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (is *InterfaceService) GetUrlCompanyUserOceanengine(ctx context.Context, in *v1.GetUrlCompanyUserOceanengineRequest) (*v1.GetUrlCompanyUserOceanengineReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "account")

	if err != nil {
		return nil, err
	}

	oceanengineConfig, err := is.cuuc.GetUrlCompanyUserOceanengineAccounts(ctx, in.OceanengineType, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	return &v1.GetUrlCompanyUserOceanengineReply{
		Code: 200,
		Data: &v1.GetUrlCompanyUserOceanengineReply_Data{
			RedirectUrl: oceanengineConfig.Data.RedirectUrl,
		},
	}, nil
}

func (is *InterfaceService) ListCompanyUsers(ctx context.Context, in *v1.ListCompanyUsersRequest) (*v1.ListCompanyUsersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	companyUsers, err := is.cuuc.ListCompanyUsers(ctx, in.PageNum, in.PageSize, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUsersReply_CompanyUsers, 0)

	for _, companyUser := range companyUsers.Data.List {
		list = append(list, &v1.ListCompanyUsersReply_CompanyUsers{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      companyUser.Role,
			Status:    companyUser.Status,
			RoleName:  companyUser.RoleName,
		})
	}

	return &v1.ListCompanyUsersReply{
		Code: 200,
		Data: &v1.ListCompanyUsersReply_Data{
			PageNum:   companyUsers.Data.PageNum,
			PageSize:  companyUsers.Data.PageSize,
			Total:     companyUsers.Data.Total,
			TotalPage: companyUsers.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) StatisticsCompanyUsers(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsCompanyUsersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	statistics, err := is.cuuc.StatisticsCompanyUsers(ctx, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsCompanyUsersReply_Statistics, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsCompanyUsersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsCompanyUsersReply{
		Code: 200,
		Data: &v1.StatisticsCompanyUsersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (is *InterfaceService) CreateCompanyUsers(ctx context.Context, in *v1.CreateCompanyUsersRequest) (*v1.CreateCompanyUsersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if companyUser.Data.IsWhite == 0 {
		if companyUser.Data.Role == 3 {
			return nil, errors.InternalServer("INTERFACE_COMPANY_USER_CREATE_ERROR", "暂无权限")
		}
	}

	createCompanyUser, err := is.cuuc.CreateCompanyUsers(ctx, companyUser.Data.CurrentCompanyId, in.Username, in.Job, in.Phone, in.Role)

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyUsersReply{
		Code: 200,
		Data: &v1.CreateCompanyUsersReply_Data{
			Id:        createCompanyUser.Data.Id,
			CompanyId: createCompanyUser.Data.CompanyId,
			Username:  createCompanyUser.Data.Username,
			Job:       createCompanyUser.Data.Job,
			Phone:     createCompanyUser.Data.Phone,
			Role:      createCompanyUser.Data.Role,
			Status:    createCompanyUser.Data.Status,
			RoleName:  createCompanyUser.Data.RoleName,
		},
	}, nil
}

func (is *InterfaceService) UpdateCompanyUsers(ctx context.Context, in *v1.UpdateCompanyUsersRequest) (*v1.UpdateCompanyUsersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	if companyUser.Data.IsWhite == 0 {
		if companyUser.Data.Role == 3 {
			return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", "暂无权限")
		}
	}

	updateCompanyUser, err := is.cuuc.UpdateCompanyUsers(ctx, in.Id, companyUser.Data.CurrentCompanyId, in.Username, in.Job, in.Phone, in.Role)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateCompanyUsersReply_Data{
			Id:        updateCompanyUser.Data.Id,
			CompanyId: updateCompanyUser.Data.CompanyId,
			Username:  updateCompanyUser.Data.Username,
			Job:       updateCompanyUser.Data.Job,
			Phone:     updateCompanyUser.Data.Phone,
			Role:      updateCompanyUser.Data.Role,
			Status:    updateCompanyUser.Data.Status,
			RoleName:  updateCompanyUser.Data.RoleName,
		},
	}, nil
}

func (is *InterfaceService) UpdateStatusCompanyUsers(ctx context.Context, in *v1.UpdateStatusCompanyUsersRequest) (*v1.UpdateStatusCompanyUsersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	if companyUser.Data.IsWhite == 0 {
		if companyUser.Data.Role == 3 {
			return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", "暂无权限")
		}
	}

	updateStatusCompanyUser, err := is.cuuc.UpdateStatusCompanyUsers(ctx, in.Id, companyUser.Data.CurrentCompanyId, in.Status)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateStatusCompanyUsersReply_Data{
			Id:        updateStatusCompanyUser.Data.Id,
			CompanyId: updateStatusCompanyUser.Data.CompanyId,
			Username:  updateStatusCompanyUser.Data.Username,
			Job:       updateStatusCompanyUser.Data.Job,
			Phone:     updateStatusCompanyUser.Data.Phone,
			Role:      updateStatusCompanyUser.Data.Role,
			Status:    updateStatusCompanyUser.Data.Status,
			RoleName:  updateStatusCompanyUser.Data.RoleName,
		},
	}, nil
}

func (is *InterfaceService) ListQianchuanAdvertisersCompanyUsers(ctx context.Context, in *v1.ListQianchuanAdvertisersCompanyUsersRequest) (*v1.ListQianchuanAdvertisersCompanyUsersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	if companyUser.Data.IsWhite == 0 {
		if companyUser.Data.Role == 3 {
			return nil, errors.InternalServer("INTERFACE_COMPANY_USER_LIST_QIANCHUAN_ADVERTISERS_ERROR", "暂无权限")
		}
	}

	qianchuanAdvertisers, err := is.cuuc.ListQianchuanAdvertisersCompanyUsers(ctx, in.Id, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdvertisersCompanyUsersReply_Advertisers, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
		list = append(list, &v1.ListQianchuanAdvertisersCompanyUsersReply_Advertisers{
			AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
			AdvertiserName: qianchuanAdvertiser.AdvertiserName,
			CompanyName:    qianchuanAdvertiser.CompanyName,
			Status:         qianchuanAdvertiser.Status,
			IsSelect:       qianchuanAdvertiser.IsSelect,
		})
	}

	return &v1.ListQianchuanAdvertisersCompanyUsersReply{
		Code: 200,
		Data: &v1.ListQianchuanAdvertisersCompanyUsersReply_Data{
			List: list,
		},
	}, nil
}

func (is *InterfaceService) UpdateRoleCompanyUsers(ctx context.Context, in *v1.UpdateRoleCompanyUsersRequest) (*v1.UpdateRoleCompanyUsersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	if companyUser.Data.IsWhite == 0 {
		if companyUser.Data.Role == 3 {
			return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", "暂无权限")
		}
	}

	updateRoleCompanyUser, err := is.cuuc.UpdateRoleCompanyUsers(ctx, in.Id, companyUser.Data.CurrentCompanyId, in.RoleIds)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateRoleCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateRoleCompanyUsersReply_Data{
			Id:        updateRoleCompanyUser.Data.Id,
			CompanyId: updateRoleCompanyUser.Data.CompanyId,
			Username:  updateRoleCompanyUser.Data.Username,
			Job:       updateRoleCompanyUser.Data.Job,
			Phone:     updateRoleCompanyUser.Data.Phone,
			Role:      updateRoleCompanyUser.Data.Role,
			Status:    updateRoleCompanyUser.Data.Status,
			RoleName:  updateRoleCompanyUser.Data.RoleName,
		},
	}, nil
}

func (is *InterfaceService) DeleteCompanyUsers(ctx context.Context, in *v1.DeleteCompanyUsersRequest) (*v1.DeleteCompanyUsersReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "")

	if err != nil {
		return nil, err
	}

	if companyUser.Data.IsWhite == 0 {
		if companyUser.Data.Role == 3 {
			return nil, errors.InternalServer("INTERFACE_COMPANY_USER_DELETE_ERROR", "暂无权限")
		}
	}

	if err := is.cuuc.DeleteCompanyUsers(ctx, in.Id, companyUser.Data.Id, companyUser.Data.CurrentCompanyId); err != nil {
		return nil, err
	}

	return &v1.DeleteCompanyUsersReply{
		Code: 200,
		Data: &v1.DeleteCompanyUsersReply_Data{},
	}, nil
}
