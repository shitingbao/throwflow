package service

import (
	"context"
	v1 "interface/api/interface/v1"
	cv1 "interface/api/service/company/v1"
	wv1 "interface/api/service/weixin/v1"
	"interface/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewInterfaceService)

type InterfaceService struct {
	v1.UnimplementedInterfaceServer

	iuc   *biz.IndustryUsecase
	auc   *biz.AreaUsecase
	suc   *biz.SmsUsecase
	tuc   *biz.TokenUsecase
	kcuc  *biz.KuaidiCompanyUsecase
	luc   *biz.LoginUsecase
	cuuc  *biz.CompanyUserUsecase
	couc  *biz.CompanyUsecase
	csuc  *biz.CompanySetUsecase
	cuc   *biz.ClueUsecase
	muc   *biz.MaterialUsecase
	pruc  *biz.PerformanceRuleUsecase
	puc   *biz.PerformanceUsecase
	preuc *biz.PerformanceRebalanceUsecase
	uluc  *biz.UpdateLogUsecase
	uuc   *biz.UserUsecase
	prouc *biz.ProductUsecase
	cmuc  *biz.CompanyMaterialUsecase
	uauc  *biz.UserAddressUsecase
	usuc  *biz.UserStoreUsecase
	uoduc *biz.UserOpenDouyinUsecase
	usruc *biz.UserScanRecordUsecase
	jouc  *biz.JinritemailOrderUsecase
	douc  *biz.DoukeOrderUsecase
	usouc *biz.UserSampleOrderUsecase
	uouc  *biz.UserOrganizationUsecase
	ucuc  *biz.UserCouponUsecase
	ubuc  *biz.UserBalanceUsecase
	coruc *biz.CompanyOrganizationUsecase
	ctuc  *biz.CompanyTaskUsecase
	ubauc *biz.UserBankUsecase
	ucouc *biz.UserContractUsecase
}

func NewInterfaceService(iuc *biz.IndustryUsecase, auc *biz.AreaUsecase, suc *biz.SmsUsecase, tuc *biz.TokenUsecase, kcuc *biz.KuaidiCompanyUsecase, luc *biz.LoginUsecase, cuuc *biz.CompanyUserUsecase, couc *biz.CompanyUsecase, csuc *biz.CompanySetUsecase, cuc *biz.ClueUsecase, muc *biz.MaterialUsecase, pruc *biz.PerformanceRuleUsecase, puc *biz.PerformanceUsecase, preuc *biz.PerformanceRebalanceUsecase, uluc *biz.UpdateLogUsecase, uuc *biz.UserUsecase, prouc *biz.ProductUsecase, cmuc *biz.CompanyMaterialUsecase, uauc *biz.UserAddressUsecase, usuc *biz.UserStoreUsecase, uoduc *biz.UserOpenDouyinUsecase, usruc *biz.UserScanRecordUsecase, jouc *biz.JinritemailOrderUsecase, douc *biz.DoukeOrderUsecase, usouc *biz.UserSampleOrderUsecase, uouc *biz.UserOrganizationUsecase, ucuc *biz.UserCouponUsecase, ubuc *biz.UserBalanceUsecase, ubauc *biz.UserBankUsecase, ucouc *biz.UserContractUsecase, coruc *biz.CompanyOrganizationUsecase, ctuc *biz.CompanyTaskUsecase) *InterfaceService {
	return &InterfaceService{iuc: iuc, auc: auc, suc: suc, tuc: tuc, kcuc: kcuc, luc: luc, cuuc: cuuc, couc: couc, csuc: csuc, cuc: cuc, muc: muc, pruc: pruc, puc: puc, preuc: preuc, uluc: uluc, uuc: uuc, prouc: prouc, cmuc: cmuc, uauc: uauc, usuc: usuc, uoduc: uoduc, usruc: usruc, jouc: jouc, douc: douc, usouc: usouc, uouc: uouc, ucuc: ucuc, ubuc: ubuc, ubauc: ubauc, ucouc: ucouc, coruc: coruc, ctuc: ctuc}
}

func (is *InterfaceService) verifyToken(ctx context.Context, token string) bool {
	return is.tuc.VerifyToken(ctx, token)
}

func (is *InterfaceService) verifyCode(ctx context.Context, phone, types, code string) bool {
	return is.suc.VerifyCode(ctx, phone, types, code)
}

func (is *InterfaceService) verifyLogin(ctx context.Context, verifyLoginStatus, verifyLoginMenu bool, permissionCode string) (*cv1.GetCompanyUserReply, error) {
	companyUser, err := is.cuuc.GetCompanyUser(ctx, verifyLoginStatus)

	if err != nil {
		return nil, err
	}

	if verifyLoginMenu {
		menus, err := is.cuuc.ListCompanyUserMenu(ctx, companyUser.Data.CurrentCompanyId)

		if err != nil {
			return nil, err
		}

		for _, menu := range menus.Data.List {
			if menu.PermissionCode == permissionCode {
				return companyUser, nil
			}

			for _, cmenu := range menu.ChildList {
				if cmenu.PermissionCode == permissionCode {
					return companyUser, nil
				}
			}
		}

		return nil, errors.InternalServer("INTERFACE_NOT_PERMISSION", "没有权限")
	}

	return companyUser, nil
}

func (is *InterfaceService) verifyMiniUserLogin(ctx context.Context) (*wv1.GetUsersReply, error) {
	user, err := is.uuc.GetUser(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}
