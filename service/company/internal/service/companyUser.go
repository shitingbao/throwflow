package service

import (
	v1 "company/api/company/v1"
	"company/internal/biz"
	"company/internal/pkg/tool"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"strconv"
	"strings"
	"time"
)

func (cs *CompanyService) ListCompanyUsers(ctx context.Context, in *v1.ListCompanyUsersRequest) (*v1.ListCompanyUsersReply, error) {
	companyUsers, err := cs.cuuc.ListCompanyUsers(ctx, in.CompanyId, in.PageNum, in.PageSize, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUsersReply_CompanyUser, 0)

	for _, companyUser := range companyUsers.List {
		list = append(list, &v1.ListCompanyUsersReply_CompanyUser{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      uint32(companyUser.Role),
			Status:    uint32(companyUser.Status),
			IsWhite:   uint32(companyUser.IsWhite),
			RoleName:  companyUser.RoleName,
		})
	}

	totalPage := uint64(math.Ceil(float64(companyUsers.Total) / float64(companyUsers.PageSize)))

	return &v1.ListCompanyUsersReply{
		Code: 200,
		Data: &v1.ListCompanyUsersReply_Data{
			PageNum:   companyUsers.PageNum,
			PageSize:  companyUsers.PageSize,
			Total:     companyUsers.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUsersByPhone(ctx context.Context, in *v1.ListCompanyUsersByPhoneRequest) (*v1.ListCompanyUsersByPhoneReply, error) {
	companyUsers, err := cs.cuuc.ListCompanyUsersByPhone(ctx, in.Phone)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUsersByPhoneReply_CompanyUser, 0)

	for _, companyUser := range companyUsers {
		list = append(list, &v1.ListCompanyUsersByPhoneReply_CompanyUser{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      uint32(companyUser.Role),
			Status:    uint32(companyUser.Status),
			IsWhite:   uint32(companyUser.IsWhite),
			RoleName:  companyUser.RoleName,
		})
	}

	return &v1.ListCompanyUsersByPhoneReply{
		Code: 200,
		Data: &v1.ListCompanyUsersByPhoneReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListSelectCompanyUsers(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCompanyUsersReply, error) {
	selects, err := cs.cuuc.ListSelectCompanyUsers(ctx)

	if err != nil {
		return nil, err
	}

	role := make([]*v1.ListSelectCompanyUsersReply_Role, 0)

	for _, lrole := range selects.Role {
		role = append(role, &v1.ListSelectCompanyUsersReply_Role{
			Key:   lrole.Key,
			Value: lrole.Value,
		})
	}

	return &v1.ListSelectCompanyUsersReply{
		Code: 200,
		Data: &v1.ListSelectCompanyUsersReply_Data{
			Role: role,
		},
	}, nil
}

func (cs *CompanyService) StatisticsCompanyUsers(ctx context.Context, in *v1.StatisticsCompanyUsersRequest) (*v1.StatisticsCompanyUsersReply, error) {
	selects, err := cs.cuuc.StatisticsCompanyUsers(ctx, in.CompanyId)

	if err != nil {
		return nil, err
	}

	statistics := make([]*v1.StatisticsCompanyUsersReply_Statistic, 0)

	for _, sstatistics := range selects.Statistics {
		statistics = append(statistics, &v1.StatisticsCompanyUsersReply_Statistic{
			Key:   sstatistics.Key,
			Value: sstatistics.Value,
		})
	}

	return &v1.StatisticsCompanyUsersReply{
		Code: 200,
		Data: &v1.StatisticsCompanyUsersReply_Data{
			Statistics: statistics,
		},
	}, nil
}

func (cs *CompanyService) ListQianchuanAdvertisersCompanyUsers(ctx context.Context, in *v1.ListQianchuanAdvertisersCompanyUsersRequest) (*v1.ListQianchuanAdvertisersCompanyUsersReply, error) {
	qianchuanAdvertisers, err := cs.cuuc.ListQianchuanAdvertisersCompanyUsers(ctx, in.Id, in.CompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdvertisersCompanyUsersReply_Advertiser, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		list = append(list, &v1.ListQianchuanAdvertisersCompanyUsersReply_Advertiser{
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

func (cs *CompanyService) GetCompanyUsers(ctx context.Context, in *v1.GetCompanyUsersRequest) (*v1.GetCompanyUsersReply, error) {
	companyUser, err := cs.cuuc.GetCompanyUsers(ctx, in.Id, in.CompanyId)

	if err != nil {
		return nil, err
	}

	return &v1.GetCompanyUsersReply{
		Code: 200,
		Data: &v1.GetCompanyUsersReply_Data{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      uint32(companyUser.Role),
			Status:    uint32(companyUser.Status),
			IsWhite:   uint32(companyUser.IsWhite),
			RoleName:  companyUser.RoleName,
		},
	}, nil
}

func (cs *CompanyService) CreateCompanyUsers(ctx context.Context, in *v1.CreateCompanyUsersRequest) (*v1.CreateCompanyUsersReply, error) {
	companyUser, err := cs.cuuc.CreateCompanyUsers(ctx, in.CompanyId, in.Username, in.Job, in.Phone, uint8(in.Role))

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyUsersReply{
		Code: 200,
		Data: &v1.CreateCompanyUsersReply_Data{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      uint32(companyUser.Role),
			Status:    uint32(companyUser.Status),
			IsWhite:   uint32(companyUser.IsWhite),
			RoleName:  companyUser.RoleName,
		},
	}, nil
}

func (cs *CompanyService) UpdateCompanyUsers(ctx context.Context, in *v1.UpdateCompanyUsersRequest) (*v1.UpdateCompanyUsersReply, error) {
	companyUser, err := cs.cuuc.UpdateCompanyUsers(ctx, in.Id, in.CompanyId, in.Username, in.Job, in.Phone, uint8(in.Role))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateCompanyUsersReply_Data{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      uint32(companyUser.Role),
			Status:    uint32(companyUser.Status),
			IsWhite:   uint32(companyUser.IsWhite),
			RoleName:  companyUser.RoleName,
		},
	}, nil
}

func (cs *CompanyService) UpdateStatusCompanyUsers(ctx context.Context, in *v1.UpdateStatusCompanyUsersRequest) (*v1.UpdateStatusCompanyUsersReply, error) {
	companyUser, err := cs.cuuc.UpdateStatusCompanyUsers(ctx, in.Id, in.CompanyId, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateStatusCompanyUsersReply_Data{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      uint32(companyUser.Role),
			Status:    uint32(companyUser.Status),
			IsWhite:   uint32(companyUser.IsWhite),
			RoleName:  companyUser.RoleName,
		},
	}, nil
}

func (cs *CompanyService) UpdateRoleCompanyUsers(ctx context.Context, in *v1.UpdateRoleCompanyUsersRequest) (*v1.UpdateRoleCompanyUsersReply, error) {
	companyUser, err := cs.cuuc.UpdateRoleCompanyUsers(ctx, in.Id, in.CompanyId, in.RoleIds)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateRoleCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateRoleCompanyUsersReply_Data{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      uint32(companyUser.Role),
			Status:    uint32(companyUser.Status),
			IsWhite:   uint32(companyUser.IsWhite),
			RoleName:  companyUser.RoleName,
		},
	}, nil
}

func (cs *CompanyService) UpdateWhiteCompanyUsers(ctx context.Context, in *v1.UpdateWhiteCompanyUsersRequest) (*v1.UpdateWhiteCompanyUsersReply, error) {
	companyUser, err := cs.cuuc.UpdateWhiteCompanyUsers(ctx, in.Id, in.CompanyId, uint8(in.IsWhite))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateWhiteCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateWhiteCompanyUsersReply_Data{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      uint32(companyUser.Role),
			Status:    uint32(companyUser.Status),
			IsWhite:   uint32(companyUser.IsWhite),
			RoleName:  companyUser.RoleName,
		},
	}, nil
}

func (cs *CompanyService) DeleteCompanyUsers(ctx context.Context, in *v1.DeleteCompanyUsersRequest) (*v1.DeleteCompanyUsersReply, error) {
	if err := cs.cuuc.DeleteCompanyUsers(ctx, in.Id, in.CompanyId); err != nil {
		return nil, err
	}

	return &v1.DeleteCompanyUsersReply{
		Code: 200,
		Data: &v1.DeleteCompanyUsersReply_Data{},
	}, nil
}

func (cs *CompanyService) DeleteRoleCompanyUsers(ctx context.Context, in *v1.DeleteRoleCompanyUsersRequest) (*v1.DeleteRoleCompanyUsersReply, error) {
	if err := cs.cuuc.DeleteRoleCompanyUsers(ctx, in.CompanyId, in.AdvertiserId); err != nil {
		return nil, err
	}

	return &v1.DeleteRoleCompanyUsersReply{
		Code: 200,
		Data: &v1.DeleteRoleCompanyUsersReply_Data{},
	}, nil
}

func (cs *CompanyService) LoginCompanyUser(ctx context.Context, in *v1.LoginCompanyUserRequest) (*v1.LoginCompanyUserReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	loginCompanyUser, err := cs.cuuc.LoginCompanyUser(ctx, in.Phone)

	if err != nil {
		return nil, err
	}

	userCompanys := make([]*v1.LoginCompanyUserReply_Company, 0)

	for _, userCompany := range loginCompanyUser.UserCompany {
		userCompanys = append(userCompanys, &v1.LoginCompanyUserReply_Company{
			CompanyId:   userCompany.Id,
			CompanyName: userCompany.CompanyName,
		})
	}

	return &v1.LoginCompanyUserReply{
		Code: 200,
		Data: &v1.LoginCompanyUserReply_Data{
			Id:                   loginCompanyUser.Id,
			CompanyId:            loginCompanyUser.CompanyId,
			Username:             loginCompanyUser.Username,
			Job:                  loginCompanyUser.Job,
			Phone:                loginCompanyUser.Phone,
			Role:                 uint32(loginCompanyUser.Role),
			RoleName:             loginCompanyUser.RoleName,
			IsWhite:              uint32(loginCompanyUser.IsWhite),
			CompanyType:          uint32(loginCompanyUser.CompanyType),
			CompanyTypeName:      loginCompanyUser.CompanyTypeName,
			CompanyName:          loginCompanyUser.CompanyName,
			CompanyStartTime:     tool.TimeToString("2006-01-02", loginCompanyUser.CompanyStartTime),
			CompanyEndTime:       tool.TimeToString("2006-01-02", loginCompanyUser.CompanyEndTime),
			Accounts:             loginCompanyUser.Accounts,
			QianchuanAdvertisers: loginCompanyUser.QianchuanAdvertisers,
			IsTermwork:           uint32(loginCompanyUser.IsTermwork),
			Token:                loginCompanyUser.Token,
			Reason:               loginCompanyUser.Reason,
			CurrentCompanyId:     loginCompanyUser.CurrentCompanyId,
			UserCompany:          userCompanys,
		},
	}, nil
}

func (cs *CompanyService) GetCompanyUser(ctx context.Context, in *v1.GetCompanyUserRequest) (*v1.GetCompanyUserReply, error) {
	loginCompanyUser, err := cs.cuuc.GetCompanyUser(ctx, in.Token)

	if err != nil {
		return nil, err
	}

	userCompanys := make([]*v1.GetCompanyUserReply_Company, 0)

	for _, userCompany := range loginCompanyUser.UserCompany {
		userCompanys = append(userCompanys, &v1.GetCompanyUserReply_Company{
			CompanyId:   userCompany.Id,
			CompanyName: userCompany.CompanyName,
		})
	}

	return &v1.GetCompanyUserReply{
		Code: 200,
		Data: &v1.GetCompanyUserReply_Data{
			Id:                   loginCompanyUser.Id,
			CompanyId:            loginCompanyUser.CompanyId,
			Username:             loginCompanyUser.Username,
			Job:                  loginCompanyUser.Job,
			Phone:                loginCompanyUser.Phone,
			Role:                 uint32(loginCompanyUser.Role),
			RoleName:             loginCompanyUser.RoleName,
			IsWhite:              uint32(loginCompanyUser.IsWhite),
			CompanyType:          uint32(loginCompanyUser.CompanyType),
			CompanyTypeName:      loginCompanyUser.CompanyTypeName,
			CompanyName:          loginCompanyUser.CompanyName,
			CompanyStartTime:     tool.TimeToString("2006-01-02", loginCompanyUser.CompanyStartTime),
			CompanyEndTime:       tool.TimeToString("2006-01-02", loginCompanyUser.CompanyEndTime),
			Accounts:             loginCompanyUser.Accounts,
			QianchuanAdvertisers: loginCompanyUser.QianchuanAdvertisers,
			IsTermwork:           uint32(loginCompanyUser.IsTermwork),
			Reason:               loginCompanyUser.Reason,
			CurrentCompanyId:     loginCompanyUser.CurrentCompanyId,
			UserCompany:          userCompanys,
		},
	}, nil
}

func (cs *CompanyService) LogoutCompanyUser(ctx context.Context, in *v1.LogoutCompanyUserRequest) (*v1.LogoutCompanyUserReply, error) {
	if err := cs.cuuc.LogoutCompanyUser(ctx, in.Token); err != nil {
		return nil, err
	}

	return &v1.LogoutCompanyUserReply{
		Code: 200,
		Data: &v1.LogoutCompanyUserReply_Data{},
	}, nil
}

func (cs *CompanyService) ChangeCompanyUserCompany(ctx context.Context, in *v1.ChangeCompanyUserCompanyRequest) (*v1.ChangeCompanyUserCompanyReply, error) {
	if err := cs.cuuc.ChangeCompanyUserCompany(ctx, in.Token, in.CompanyId); err != nil {
		return nil, err
	}

	return &v1.ChangeCompanyUserCompanyReply{
		Code: 200,
		Data: &v1.ChangeCompanyUserCompanyReply_Data{},
	}, nil
}

func (cs *CompanyService) ListCompanyUserMenu(ctx context.Context, in *v1.ListCompanyUserMenuRequest) (*v1.ListCompanyUserMenuReply, error) {
	menus, err := cs.muc.ListMenus(ctx)

	if err != nil {
		return nil, err
	}

	company, err := cs.couc.GetCompanys(ctx, in.CompanyId)

	if err != nil {
		return nil, err
	}

	companyMenus := make([]string, 0)

	if len(company.MenuId) > 0 {
		companyMenus = strings.Split(company.MenuId, ",")
	}

	list := make([]*v1.ListCompanyUserMenuReply_Menu, 0)

	for _, menu := range menus {
		if menu.Status == 1 {
			childList := make([]*v1.ListCompanyUserMenuReply_ChildMenu, 0)

			for _, childMenu := range menu.ChildList {
				if childMenu.Status == 1 {
					lchildMenu := &v1.ListCompanyUserMenuReply_ChildMenu{
						MenuName:       childMenu.MenuName,
						MenuType:       childMenu.MenuType,
						Filename:       childMenu.Filename,
						Filepath:       childMenu.Filepath,
						PermissionCode: childMenu.PermissionCode,
						IsPermission:   "0",
					}

					for _, companyMenu := range companyMenus {
						if iCompanyMenu, err := strconv.ParseUint(companyMenu, 10, 64); err == nil {
							if iCompanyMenu == childMenu.Id {
								lchildMenu.IsPermission = "1"

								break
							}
						}
					}

					childList = append(childList, lchildMenu)
				}
			}

			list = append(list, &v1.ListCompanyUserMenuReply_Menu{
				MenuName:       menu.MenuName,
				ChildList:      childList,
				MenuType:       menu.MenuType,
				Filename:       menu.Filename,
				Filepath:       menu.Filepath,
				PermissionCode: menu.PermissionCode,
			})
		}
	}

	return &v1.ListCompanyUserMenuReply{
		Code: 200,
		Data: &v1.ListCompanyUserMenuReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserQianchuanAdvertisers(ctx context.Context, in *v1.ListCompanyUserQianchuanAdvertisersRequest) (*v1.ListCompanyUserQianchuanAdvertisersReply, error) {
	qianchuanAdvertisers, err := cs.cuuc.ListCompanyUserQianchuanAdvertisers(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), in.Day, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserQianchuanAdvertisersReply_Advertiser, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		list = append(list, &v1.ListCompanyUserQianchuanAdvertisersReply_Advertiser{
			AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
			CompanyId:      qianchuanAdvertiser.CompanyId,
			AccountId:      qianchuanAdvertiser.AccountId,
			AdvertiserName: qianchuanAdvertiser.AdvertiserName,
			CompanyName:    qianchuanAdvertiser.CompanyName,
		})
	}

	return &v1.ListCompanyUserQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.ListCompanyUserQianchuanAdvertisersReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserQianchuanReportAdvertisers(ctx context.Context, in *v1.ListCompanyUserQianchuanReportAdvertisersRequest) (*v1.ListCompanyUserQianchuanReportAdvertisersReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	qianchuanReportAdvertisers, err := cs.cuuc.ListCompanyUserQianchuanReportAdvertisers(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), in.Day)

	if err != nil {
		return nil, err
	}

	payOrderAmountList := make([]*v1.ListCompanyUserQianchuanReportAdvertisersReply_PayOrderAmount, 0)
	statCostList := make([]*v1.ListCompanyUserQianchuanReportAdvertisersReply_StatCost, 0)

	for _, qianchuanReportAdvertiser := range qianchuanReportAdvertisers {
		if qianchuanReportAdvertiser.Key == "payOrderAmount" {
			for _, l := range qianchuanReportAdvertiser.List {
				payOrderAmountList = append(payOrderAmountList, &v1.ListCompanyUserQianchuanReportAdvertisersReply_PayOrderAmount{
					Time:           l.Time,
					Value:          fmt.Sprintf("%.2f", l.Value),
					YesterdayValue: fmt.Sprintf("%.2f", l.YesterdayValue),
				})
			}
		} else if qianchuanReportAdvertiser.Key == "statCost" {
			for _, l := range qianchuanReportAdvertiser.List {
				statCostList = append(statCostList, &v1.ListCompanyUserQianchuanReportAdvertisersReply_StatCost{
					Time:           l.Time,
					Value:          fmt.Sprintf("%.2f", l.Value),
					YesterdayValue: fmt.Sprintf("%.2f", l.YesterdayValue),
				})
			}
		}
	}

	return &v1.ListCompanyUserQianchuanReportAdvertisersReply{
		Code: 200,
		Data: &v1.ListCompanyUserQianchuanReportAdvertisersReply_Data{
			PayOrderAmounts: payOrderAmountList,
			StatCosts:       statCostList,
		},
	}, nil
}

func (cs *CompanyService) StatisticsCompanyUserDashboardQianchuanAdvertisers(ctx context.Context, in *v1.StatisticsCompanyUserDashboardQianchuanAdvertisersRequest) (*v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dashboardQianchuanAdvertisers, err := cs.cuuc.StatisticsCompanyUserDashboardQianchuanAdvertisers(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), in.Day)

	if err != nil {
		return nil, err
	}

	statistics := make([]*v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic, 0)

	if dashboardQianchuanAdvertisers == nil {
		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "已授权账户",
			Value: "0",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "有消耗账户",
			Value: "0",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "总单数",
			Value: "0",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "短视频单数",
			Value: "0",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "直播单数",
			Value: "0",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "总成交",
			Value: "0.00¥",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "短视频成交",
			Value: "0.00¥",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "直播成交",
			Value: "0.00¥",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "总消耗",
			Value: "0.00¥",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "短视频消耗",
			Value: "0.00¥",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "直播消耗",
			Value: "0.00¥",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "总广告ROI",
			Value: "0.00",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "短视频ROI",
			Value: "0.00",
		})

		statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
			Key:   "直播ROI",
			Value: "0.00",
		})
	} else {
		for _, dashboardQianchuanAdvertiser := range dashboardQianchuanAdvertisers.Data.Statistics {
			statistics = append(statistics, &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Statistic{
				Key:   dashboardQianchuanAdvertiser.Key,
				Value: dashboardQianchuanAdvertiser.Value,
			})
		}
	}

	return &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.StatisticsCompanyUserDashboardQianchuanAdvertisersReply_Data{
			Statistics: statistics,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserQianchuanReportAwemes(ctx context.Context, in *v1.ListCompanyUserQianchuanReportAwemesRequest) (*v1.ListCompanyUserQianchuanReportAwemesReply, error) {
	qianchuanReportAwemes, err := cs.cuuc.ListCompanyUserQianchuanReportAwemes(ctx, in.UserId, in.CompanyId, in.PageNum, in.PageSize, in.IsDistinction, uint8(in.IsWhite), in.Day, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserQianchuanReportAwemesReply_Aweme, 0)

	if qianchuanReportAwemes == nil {
		return &v1.ListCompanyUserQianchuanReportAwemesReply{
			Code: 200,
			Data: &v1.ListCompanyUserQianchuanReportAwemesReply_Data{
				PageNum:   in.PageNum,
				PageSize:  in.PageSize,
				Total:     0,
				TotalPage: 0,
				List:      list,
			},
		}, nil
	} else {
		for _, qianchuanReportAweme := range qianchuanReportAwemes.Data.List {
			list = append(list, &v1.ListCompanyUserQianchuanReportAwemesReply_Aweme{
				AdvertiserId:            qianchuanReportAweme.AdvertiserId,
				AdvertiserName:          qianchuanReportAweme.AdvertiserName,
				AwemeId:                 qianchuanReportAweme.AwemeId,
				AwemeName:               qianchuanReportAweme.AwemeName,
				AwemeShowId:             qianchuanReportAweme.AwemeShowId,
				AwemeAvatar:             qianchuanReportAweme.AwemeAvatar,
				AwemeUrl:                qianchuanReportAweme.AwemeUrl,
				DyFollow:                qianchuanReportAweme.DyFollow,
				StatCost:                qianchuanReportAweme.StatCost,
				PayOrderCount:           qianchuanReportAweme.PayOrderCount,
				PayOrderAmount:          qianchuanReportAweme.PayOrderAmount,
				PayOrderAveragePrice:    qianchuanReportAweme.PayOrderAveragePrice,
				Roi:                     qianchuanReportAweme.Roi,
				ShowCnt:                 qianchuanReportAweme.ShowCnt,
				ClickCnt:                qianchuanReportAweme.ClickCnt,
				ClickRate:               qianchuanReportAweme.ClickRate,
				ConvertCnt:              qianchuanReportAweme.ConvertCnt,
				ConvertRate:             qianchuanReportAweme.ConvertRate,
				AveragePayOrderStatCost: qianchuanReportAweme.AveragePayOrderStatCost,
			})
		}

		return &v1.ListCompanyUserQianchuanReportAwemesReply{
			Code: 200,
			Data: &v1.ListCompanyUserQianchuanReportAwemesReply_Data{
				PageNum:   qianchuanReportAwemes.Data.PageNum,
				PageSize:  qianchuanReportAwemes.Data.PageSize,
				Total:     qianchuanReportAwemes.Data.Total,
				TotalPage: qianchuanReportAwemes.Data.TotalPage,
				List:      list,
			},
		}, nil
	}
}

func (cs *CompanyService) StatisticsCompanyUserQianchuanReportAwemes(ctx context.Context, in *v1.StatisticsCompanyUserQianchuanReportAwemesRequest) (*v1.StatisticsCompanyUserQianchuanReportAwemesReply, error) {
	qianchuanReportAwemes, err := cs.cuuc.StatisticsCompanyUserQianchuanReportAwemes(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), in.Day)

	if err != nil {
		return nil, err
	}

	statistics := make([]*v1.StatisticsCompanyUserQianchuanReportAwemesReply_Statistic, 0)

	for _, qianchuanReportAweme := range qianchuanReportAwemes.Data.Statistics {
		statistics = append(statistics, &v1.StatisticsCompanyUserQianchuanReportAwemesReply_Statistic{
			Key:   qianchuanReportAweme.Key,
			Value: qianchuanReportAweme.Value,
		})
	}

	return &v1.StatisticsCompanyUserQianchuanReportAwemesReply{
		Code: 200,
		Data: &v1.StatisticsCompanyUserQianchuanReportAwemesReply_Data{
			Statistics: statistics,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserQianchuanReportProducts(ctx context.Context, in *v1.ListCompanyUserQianchuanReportProductsRequest) (*v1.ListCompanyUserQianchuanReportProductsReply, error) {
	qianchuanReportProducts, err := cs.cuuc.ListCompanyUserQianchuanReportProducts(ctx, in.UserId, in.CompanyId, in.AdvertiserId, in.PageNum, in.PageSize, in.IsDistinction, uint8(in.IsWhite), in.Day, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserQianchuanReportProductsReply_Product, 0)

	if qianchuanReportProducts == nil {
		return &v1.ListCompanyUserQianchuanReportProductsReply{
			Code: 200,
			Data: &v1.ListCompanyUserQianchuanReportProductsReply_Data{
				PageNum:   in.PageNum,
				PageSize:  in.PageSize,
				Total:     0,
				TotalPage: 0,
				List:      list,
			},
		}, nil
	} else {
		for _, qianchuanReportProduct := range qianchuanReportProducts.Data.List {
			list = append(list, &v1.ListCompanyUserQianchuanReportProductsReply_Product{
				AdvertiserId:            qianchuanReportProduct.AdvertiserId,
				AdvertiserName:          qianchuanReportProduct.AdvertiserName,
				ProductId:               qianchuanReportProduct.ProductId,
				DiscountPrice:           qianchuanReportProduct.DiscountPrice,
				ProductName:             qianchuanReportProduct.ProductName,
				ProductImg:              qianchuanReportProduct.ProductImg,
				ProductUrl:              qianchuanReportProduct.ProductUrl,
				StatCost:                qianchuanReportProduct.StatCost,
				PayOrderCount:           qianchuanReportProduct.PayOrderCount,
				PayOrderAmount:          qianchuanReportProduct.PayOrderAmount,
				PayOrderAveragePrice:    qianchuanReportProduct.PayOrderAveragePrice,
				Roi:                     qianchuanReportProduct.Roi,
				ShowCnt:                 qianchuanReportProduct.ShowCnt,
				ClickCnt:                qianchuanReportProduct.ClickCnt,
				ClickRate:               qianchuanReportProduct.ClickRate,
				ConvertCnt:              qianchuanReportProduct.ConvertCnt,
				ConvertRate:             qianchuanReportProduct.ConvertRate,
				AveragePayOrderStatCost: qianchuanReportProduct.AveragePayOrderStatCost,
			})
		}

		return &v1.ListCompanyUserQianchuanReportProductsReply{
			Code: 200,
			Data: &v1.ListCompanyUserQianchuanReportProductsReply_Data{
				PageNum:   qianchuanReportProducts.Data.PageNum,
				PageSize:  qianchuanReportProducts.Data.PageSize,
				Total:     qianchuanReportProducts.Data.Total,
				TotalPage: qianchuanReportProducts.Data.TotalPage,
				List:      list,
			},
		}, nil
	}
}

func (cs *CompanyService) StatisticsCompanyUserQianchuanReportProducts(ctx context.Context, in *v1.StatisticsCompanyUserQianchuanReportProductsRequest) (*v1.StatisticsCompanyUserQianchuanReportProductsReply, error) {
	qianchuanReportProducts, err := cs.cuuc.StatisticsCompanyUserQianchuanReportProducts(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), in.Day)

	if err != nil {
		return nil, err
	}

	statistics := make([]*v1.StatisticsCompanyUserQianchuanReportProductsReply_Statistic, 0)

	for _, qianchuanReportProduct := range qianchuanReportProducts.Data.Statistics {
		statistics = append(statistics, &v1.StatisticsCompanyUserQianchuanReportProductsReply_Statistic{
			Key:   qianchuanReportProduct.Key,
			Value: qianchuanReportProduct.Value,
		})
	}

	return &v1.StatisticsCompanyUserQianchuanReportProductsReply{
		Code: 200,
		Data: &v1.StatisticsCompanyUserQianchuanReportProductsReply_Data{
			Statistics: statistics,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserExternalQianchuanAdvertisers(ctx context.Context, in *v1.ListCompanyUserExternalQianchuanAdvertisersRequest) (*v1.ListCompanyUserExternalQianchuanAdvertisersReply, error) {
	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.CompanyValidatorError
		}

		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	}

	qianchuanAdvertisers, err := cs.cuuc.ListCompanyUserExternalQianchuanAdvertisers(ctx, in.UserId, in.CompanyId, in.PageNum, in.PageSize, uint8(in.IsWhite), startTime, endTime, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserExternalQianchuanAdvertisersReply_Advertiser, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
		list = append(list, &v1.ListCompanyUserExternalQianchuanAdvertisersReply_Advertiser{
			AdvertiserId:            qianchuanAdvertiser.AdvertiserId,
			AdvertiserName:          qianchuanAdvertiser.AdvertiserName,
			GeneralTotalBalance:     qianchuanAdvertiser.GeneralTotalBalance,
			StatCost:                qianchuanAdvertiser.StatCost,
			Roi:                     qianchuanAdvertiser.Roi,
			Campaigns:               qianchuanAdvertiser.Campaigns,
			PayOrderCount:           qianchuanAdvertiser.PayOrderCount,
			PayOrderAmount:          qianchuanAdvertiser.PayOrderAmount,
			ClickRate:               qianchuanAdvertiser.ClickRate,
			PayConvertRate:          qianchuanAdvertiser.PayConvertRate,
			AveragePayOrderStatCost: qianchuanAdvertiser.AveragePayOrderStatCost,
			PayOrderAveragePrice:    qianchuanAdvertiser.PayOrderAveragePrice,
			DyFollow:                qianchuanAdvertiser.DyFollow,
		})
	}

	return &v1.ListCompanyUserExternalQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.ListCompanyUserExternalQianchuanAdvertisersReply_Data{
			PageNum:   qianchuanAdvertisers.Data.PageNum,
			PageSize:  qianchuanAdvertisers.Data.PageSize,
			Total:     qianchuanAdvertisers.Data.Total,
			TotalPage: qianchuanAdvertisers.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserExternalSelectQianchuanAdvertisers(ctx context.Context, in *v1.ListCompanyUserExternalSelectQianchuanAdvertisersRequest) (*v1.ListCompanyUserExternalSelectQianchuanAdvertisersReply, error) {
	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Before(nowTime) {
			return nil, biz.CompanyValidatorError
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.CompanyValidatorError
		}

		if !endTime.Before(nowTime) {
			return nil, biz.CompanyValidatorError
		}
	}

	qianchuanAdvertisers, err := cs.cuuc.ListCompanyUserExternalSelectQianchuanAdvertisers(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), startTime, endTime)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserExternalSelectQianchuanAdvertisersReply_Advertiser, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		list = append(list, &v1.ListCompanyUserExternalSelectQianchuanAdvertisersReply_Advertiser{
			AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
			AdvertiserName: qianchuanAdvertiser.AdvertiserName,
		})
	}

	return &v1.ListCompanyUserExternalSelectQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.ListCompanyUserExternalSelectQianchuanAdvertisersReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) StatisticsCompanyUserExternalQianchuanAdvertisers(ctx context.Context, in *v1.StatisticsCompanyUserExternalQianchuanAdvertisersRequest) (*v1.StatisticsCompanyUserExternalQianchuanAdvertisersReply, error) {
	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.CompanyValidatorError
		}

		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	}

	statistics, err := cs.cuuc.StatisticsCompanyUserExternalQianchuanAdvertisers(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), startTime, endTime, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsCompanyUserExternalQianchuanAdvertisersReply_Statistic, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsCompanyUserExternalQianchuanAdvertisersReply_Statistic{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsCompanyUserExternalQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.StatisticsCompanyUserExternalQianchuanAdvertisersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserSelectExternalQianchuanAds(ctx context.Context, in *empty.Empty) (*v1.ListCompanyUserSelectExternalQianchuanAdsReply, error) {
	selects, err := cs.cuuc.ListCompanyUserSelectExternalQianchuanAds(ctx)

	if err != nil {
		return nil, err
	}

	filter := make([]*v1.ListCompanyUserSelectExternalQianchuanAdsReply_Filter, 0)

	for _, lfilter := range selects.Data.Filter {
		filter = append(filter, &v1.ListCompanyUserSelectExternalQianchuanAdsReply_Filter{
			Key:   lfilter.Key,
			Value: lfilter.Value,
		})
	}

	return &v1.ListCompanyUserSelectExternalQianchuanAdsReply{
		Code: 200,
		Data: &v1.ListCompanyUserSelectExternalQianchuanAdsReply_Data{
			Filter: filter,
		},
	}, nil
}

func (cs *CompanyService) GetCompanyUserExternalQianchuanAds(ctx context.Context, in *v1.GetCompanyUserExternalQianchuanAdsRequest) (*v1.GetCompanyUserExternalQianchuanAdsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.CompanyValidatorError
		}

		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	}

	qianchuanAds, err := cs.cuuc.GetCompanyUserExternalQianchuanAds(ctx, in.UserId, in.CompanyId, in.AdId, uint8(in.IsWhite), in.StartDay, in.EndDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetCompanyUserExternalQianchuanAdsReply_Ad, 0)

	for _, qianchuanAd := range qianchuanAds {
		list = append(list, &v1.GetCompanyUserExternalQianchuanAdsReply_Ad{
			Time:                    qianchuanAd.Time,
			StatCost:                qianchuanAd.StatCost,
			Roi:                     qianchuanAd.Roi,
			PayOrderCount:           qianchuanAd.PayOrderCount,
			PayOrderAmount:          qianchuanAd.PayOrderAmount,
			ClickCnt:                qianchuanAd.ClickCnt,
			ShowCnt:                 qianchuanAd.ShowCnt,
			ConvertCnt:              qianchuanAd.ConvertCnt,
			ClickRate:               qianchuanAd.ClickRate,
			CpmPlatform:             qianchuanAd.CpmPlatform,
			DyFollow:                qianchuanAd.DyFollow,
			PayConvertRate:          qianchuanAd.PayConvertRate,
			ConvertCost:             qianchuanAd.ConvertCost,
			AveragePayOrderStatCost: qianchuanAd.AveragePayOrderStatCost,
			PayOrderAveragePrice:    qianchuanAd.PayOrderAveragePrice,
		})
	}

	return &v1.GetCompanyUserExternalQianchuanAdsReply{
		Code: 200,
		Data: &v1.GetCompanyUserExternalQianchuanAdsReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserExternalQianchuanAds(ctx context.Context, in *v1.ListCompanyUserExternalQianchuanAdsRequest) (*v1.ListCompanyUserExternalQianchuanAdsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.CompanyValidatorError
		}

		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	}

	qianchuanAds, err := cs.cuuc.ListCompanyUserExternalQianchuanAds(ctx, in.UserId, in.CompanyId, in.AdvertiserId, in.PageNum, in.PageSize, uint8(in.IsWhite), startTime, endTime, in.Keyword, in.Filter, in.OrderName, in.OrderType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserExternalQianchuanAdsReply_Ad, 0)

	for _, qianchuanAd := range qianchuanAds.Data.List {
		list = append(list, &v1.ListCompanyUserExternalQianchuanAdsReply_Ad{
			AdId:                    qianchuanAd.AdId,
			AdName:                  qianchuanAd.AdName,
			AdvertiserId:            qianchuanAd.AdvertiserId,
			AdvertiserName:          qianchuanAd.AdvertiserName,
			CampaignId:              qianchuanAd.CampaignId,
			CampaignName:            qianchuanAd.CampaignName,
			LabAdType:               qianchuanAd.LabAdType,
			LabAdTypeName:           qianchuanAd.LabAdTypeName,
			MarketingGoal:           qianchuanAd.MarketingGoal,
			MarketingGoalName:       qianchuanAd.MarketingGoalName,
			Status:                  qianchuanAd.Status,
			StatusName:              qianchuanAd.StatusName,
			OptStatus:               qianchuanAd.OptStatus,
			OptStatusName:           qianchuanAd.OptStatusName,
			ExternalAction:          qianchuanAd.ExternalAction,
			ExternalActionName:      qianchuanAd.ExternalActionName,
			DeepExternalAction:      qianchuanAd.DeepExternalAction,
			DeepBidType:             qianchuanAd.DeepBidType,
			PromotionId:             qianchuanAd.PromotionId,
			PromotionShowId:         qianchuanAd.PromotionShowId,
			PromotionName:           qianchuanAd.PromotionName,
			PromotionImg:            qianchuanAd.PromotionImg,
			PromotionAvatar:         qianchuanAd.PromotionAvatar,
			PromotionType:           qianchuanAd.PromotionType,
			StatCost:                qianchuanAd.StatCost,
			Roi:                     qianchuanAd.Roi,
			CpaBid:                  qianchuanAd.CpaBid,
			RoiGoal:                 qianchuanAd.RoiGoal,
			Budget:                  qianchuanAd.Budget,
			BudgetMode:              qianchuanAd.BudgetMode,
			BudgetModeName:          qianchuanAd.BudgetModeName,
			PayOrderCount:           qianchuanAd.PayOrderCount,
			PayOrderAmount:          qianchuanAd.PayOrderAmount,
			ClickCnt:                qianchuanAd.ClickCnt,
			ShowCnt:                 qianchuanAd.ShowCnt,
			ConvertCnt:              qianchuanAd.ConvertCnt,
			ClickRate:               qianchuanAd.ClickRate,
			CpmPlatform:             qianchuanAd.CpmPlatform,
			DyFollow:                qianchuanAd.DyFollow,
			PayConvertRate:          qianchuanAd.PayConvertRate,
			ConvertCost:             qianchuanAd.ConvertCost,
			AveragePayOrderStatCost: qianchuanAd.AveragePayOrderStatCost,
			PayOrderAveragePrice:    qianchuanAd.PayOrderAveragePrice,
			ConvertRate:             qianchuanAd.ConvertRate,
			AdCreateTime:            qianchuanAd.AdCreateTime,
			AdModifyTime:            qianchuanAd.AdModifyTime,
		})
	}

	return &v1.ListCompanyUserExternalQianchuanAdsReply{
		Code: 200,
		Data: &v1.ListCompanyUserExternalQianchuanAdsReply_Data{
			PageNum:   qianchuanAds.Data.PageNum,
			PageSize:  qianchuanAds.Data.PageSize,
			Total:     qianchuanAds.Data.Total,
			TotalPage: qianchuanAds.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) StatisticsCompanyUserExternalQianchuanAds(ctx context.Context, in *v1.StatisticsCompanyUserExternalQianchuanAdsRequest) (*v1.StatisticsCompanyUserExternalQianchuanAdsReply, error) {
	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.CompanyValidatorError
		}

		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.CompanyValidatorError
			}
		}
	}

	statistics, err := cs.cuuc.StatisticsCompanyUserExternalQianchuanAds(ctx, in.UserId, in.CompanyId, in.AdvertiserId, uint8(in.IsWhite), startTime, endTime, in.Keyword, in.Filter)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsCompanyUserExternalQianchuanAdsReply_Statistic, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsCompanyUserExternalQianchuanAdsReply_Statistic{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsCompanyUserExternalQianchuanAdsReply{
		Code: 200,
		Data: &v1.StatisticsCompanyUserExternalQianchuanAdsReply_Data{
			Statistics: list,
		},
	}, nil
}

func (cs *CompanyService) GetCompanyUserIsCompleteSyncData(ctx context.Context, in *v1.GetCompanyUserIsCompleteSyncDataRequest) (*v1.GetCompanyUserIsCompleteSyncDataReply, error) {
	qianchuanAdvertiserHistorys, _ := cs.cuuc.GetCompanyUserIsCompleteSyncData(ctx, in.CompanyId, in.Day)

	var isComplete uint32

	if len(qianchuanAdvertiserHistorys.Data.List) == 0 {
		isComplete = 1
	}

	return &v1.GetCompanyUserIsCompleteSyncDataReply{
		Code: 200,
		Data: &v1.GetCompanyUserIsCompleteSyncDataReply_Data{
			IsComplete: isComplete,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserExternalQianchuanSearch(ctx context.Context, in *v1.ListCompanyUserExternalQianchuanSearchRequest) (*v1.ListCompanyUserExternalQianchuanSearchReply, error) {
	searches, err := cs.cuuc.ListCompanyUserQianchuanSearch(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), in.Day, in.Keyword, in.SearchType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserExternalQianchuanSearchReply_Search, 0)

	for _, search := range searches {
		list = append(list, &v1.ListCompanyUserExternalQianchuanSearchReply_Search{
			Id:   search.Id,
			Name: search.Name,
		})
	}

	return &v1.ListCompanyUserExternalQianchuanSearchReply{
		Code: 200,
		Data: &v1.ListCompanyUserExternalQianchuanSearchReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserExternalHistoryQianchuanSearch(ctx context.Context, in *v1.ListCompanyUserExternalHistoryQianchuanSearchRequest) (*v1.ListCompanyUserExternalHistoryQianchuanSearchReply, error) {
	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Before(nowTime) {
			return nil, biz.CompanyValidatorError
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.CompanyValidatorError
		}

		if !endTime.Before(nowTime) {
			return nil, biz.CompanyValidatorError
		}
	}

	searches, err := cs.cuuc.ListCompanyUserExternalHistoryQianchuanSearch(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), startTime, endTime, in.Keyword, in.SearchType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserExternalHistoryQianchuanSearchReply_Search, 0)

	for _, search := range searches {
		list = append(list, &v1.ListCompanyUserExternalHistoryQianchuanSearchReply_Search{
			Id:   search.Id,
			Name: search.Name,
		})
	}

	return &v1.ListCompanyUserExternalHistoryQianchuanSearchReply{
		Code: 200,
		Data: &v1.ListCompanyUserExternalHistoryQianchuanSearchReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserExternalQianchuanHistorySearch(ctx context.Context, in *v1.ListCompanyUserExternalQianchuanHistorySearchRequest) (*v1.ListCompanyUserExternalQianchuanHistorySearchReply, error) {
	searches, err := cs.cuuc.ListCompanyUserExternalQianchuanHistorySearch(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), in.Day, in.SearchType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserExternalQianchuanHistorySearchReply_Search, 0)

	for _, search := range searches {
		list = append(list, &v1.ListCompanyUserExternalQianchuanHistorySearchReply_Search{
			Id:   search.Id,
			Name: search.Name,
		})
	}

	return &v1.ListCompanyUserExternalQianchuanHistorySearchReply{
		Code: 200,
		Data: &v1.ListCompanyUserExternalQianchuanHistorySearchReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyUserExternalHistoryQianchuanHistorySearch(ctx context.Context, in *v1.ListCompanyUserExternalHistoryQianchuanHistorySearchRequest) (*v1.ListCompanyUserExternalHistoryQianchuanHistorySearchReply, error) {
	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Before(nowTime) {
			return nil, biz.CompanyValidatorError
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.CompanyValidatorError
		}

		if !endTime.Before(nowTime) {
			return nil, biz.CompanyValidatorError
		}
	}

	searches, err := cs.cuuc.ListCompanyUserExternalHistoryQianchuanHistorySearch(ctx, in.UserId, in.CompanyId, uint8(in.IsWhite), startTime, endTime, in.SearchType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUserExternalHistoryQianchuanHistorySearchReply_Search, 0)

	for _, search := range searches {
		list = append(list, &v1.ListCompanyUserExternalHistoryQianchuanHistorySearchReply_Search{
			Id:   search.Id,
			Name: search.Name,
		})
	}

	return &v1.ListCompanyUserExternalHistoryQianchuanHistorySearchReply{
		Code: 200,
		Data: &v1.ListCompanyUserExternalHistoryQianchuanHistorySearchReply_Data{
			List: list,
		},
	}, nil
}
