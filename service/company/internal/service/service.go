package service

import (
	v1 "company/api/company/v1"
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"strconv"
	"strings"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewCompanyService)

type CompanyService struct {
	v1.UnimplementedCompanyServer

	muc     *biz.MenuUsecase
	iuc     *biz.IndustryUsecase
	cuc     *biz.ClueUsecase
	couc    *biz.CompanyUsecase
	cuuc    *biz.CompanyUserUsecase
	cpruc   *biz.CompanyPerformanceRuleUsecase
	cpreruc *biz.CompanyPerformanceRebalanceUsecase
	cpuc    *biz.CompanyPerformanceUsecase
	csuc    *biz.CompanySetUsecase
	cprouc  *biz.CompanyProductUsecase
	cmuc    *biz.CompanyMaterialUsecase
	coruc   *biz.CompanyOrganizationUsecase
	ctuc    *biz.CompanyTaskUsecase
	ctaruc  *biz.CompanyTaskAccountRelationUsecase
	ctduc   *biz.CompanyTaskDetailUsecase
}

func NewCompanyService(muc *biz.MenuUsecase, iuc *biz.IndustryUsecase, cuc *biz.ClueUsecase, couc *biz.CompanyUsecase, cuuc *biz.CompanyUserUsecase, cpruc *biz.CompanyPerformanceRuleUsecase, cpreruc *biz.CompanyPerformanceRebalanceUsecase, cpuc *biz.CompanyPerformanceUsecase, csuc *biz.CompanySetUsecase, cprouc *biz.CompanyProductUsecase, cmuc *biz.CompanyMaterialUsecase, coruc *biz.CompanyOrganizationUsecase, ctuc *biz.CompanyTaskUsecase, ctaruc *biz.CompanyTaskAccountRelationUsecase, ctduc *biz.CompanyTaskDetailUsecase) *CompanyService {
	return &CompanyService{muc: muc, iuc: iuc, cuc: cuc, couc: couc, cuuc: cuuc, cpruc: cpruc, cpreruc: cpreruc, cpuc: cpuc, csuc: csuc, cprouc: cprouc, cmuc: cmuc, coruc: coruc, ctuc: ctuc, ctaruc: ctaruc, ctduc: ctduc}
}

func (cs *CompanyService) verifyMenu(ctx context.Context, ids string) ([]string, error) {
	inIds := strings.Split(ids, ",")
	list := make([]string, 0)

	for _, id := range inIds {
		idUint, _ := strconv.ParseUint(id, 10, 64)

		if _, err := cs.muc.GetMenuById(ctx, idUint); err != nil {
			continue
		}

		isExist := true

		for _, lid := range list {
			if lid == id {
				isExist = false
				break
			}
		}

		if isExist {
			list = append(list, id)
		}
	}

	if len(list) == 0 {
		return nil, biz.CompanyRoleMenuNotFound
	}

	return list, nil
}

func (cs *CompanyService) verifyPermissionCode(ctx context.Context, permissionCode string) bool {
	permissionCodes := make([]*domain.PermissionCodes, 0)
	permissionCodes = domain.NewPermissionCodes()

	isExist := false

	for _, pc := range permissionCodes {
		if pc.Code == permissionCode {
			isExist = true
			break
		}
	}

	return isExist
}
