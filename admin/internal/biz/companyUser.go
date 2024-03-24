package biz

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type CompanyUserRepo interface {
	List(context.Context, uint64, uint64, string) (*v1.ListCompanyUsersReply, error)
	ListSelect(context.Context) (*v1.ListSelectCompanyUsersReply, error)
	ListQianchuanAdvertisers(context.Context, uint64, uint64) (*v1.ListQianchuanAdvertisersCompanyUsersReply, error)
	Save(context.Context, uint64, string, string, string, uint32) (*v1.CreateCompanyUsersReply, error)
	Update(context.Context, uint64, uint64, string, string, string, uint32) (*v1.UpdateCompanyUsersReply, error)
	UpdateStatus(context.Context, uint64, uint64, uint32) (*v1.UpdateStatusCompanyUsersReply, error)
	UpdateWhite(context.Context, uint64, uint64, uint32) (*v1.UpdateWhiteCompanyUsersReply, error)
	UpdateRole(context.Context, uint64, uint64, string) (*v1.UpdateRoleCompanyUsersReply, error)
	Delete(context.Context, uint64, uint64) (*v1.DeleteCompanyUsersReply, error)
}

type CompanyUserUsecase struct {
	repo CompanyUserRepo
	log  *log.Helper
}

func NewCompanyUserUsecase(repo CompanyUserRepo, logger log.Logger) *CompanyUserUsecase {
	return &CompanyUserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (cuuc *CompanyUserUsecase) ListCompanyUsers(ctx context.Context, companyId, pageNum uint64, keyword string) (*v1.ListCompanyUsersReply, error) {
	list, err := cuuc.repo.List(ctx, companyId, pageNum, keyword)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListSelectCompanyUsers(ctx context.Context) (*v1.ListSelectCompanyUsersReply, error) {
	list, err := cuuc.repo.ListSelect(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListQianchuanAdvertisersCompanyUsersRequest(ctx context.Context, id, companyId uint64) (*v1.ListQianchuanAdvertisersCompanyUsersReply, error) {
	list, err := cuuc.repo.ListQianchuanAdvertisers(ctx, id, companyId)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) CreateCompanyUsers(ctx context.Context, companyId uint64, username, job, phone string, role uint32) (*v1.CreateCompanyUsersReply, error) {
	companyUser, err := cuuc.repo.Save(ctx, companyId, username, job, phone, role)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_USER_CREATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateCompanyUsers(ctx context.Context, id, companyId uint64, username, job, phone string, role uint32) (*v1.UpdateCompanyUsersReply, error) {
	companyUser, err := cuuc.repo.Update(ctx, id, companyId, username, job, phone, role)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateStatusCompanyUsers(ctx context.Context, id, companyId uint64, status uint32) (*v1.UpdateStatusCompanyUsersReply, error) {
	companyUser, err := cuuc.repo.UpdateStatus(ctx, id, companyId, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateWhiteCompanyUsers(ctx context.Context, id, companyId uint64, isWhite uint32) (*v1.UpdateWhiteCompanyUsersReply, error) {
	companyUser, err := cuuc.repo.UpdateWhite(ctx, id, companyId, isWhite)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateRoleCompanyUsers(ctx context.Context, id, companyId uint64, roleIds string) (*v1.UpdateRoleCompanyUsersReply, error) {
	if len(roleIds) == 0 {
		return nil, errors.InternalServer("ADMIN_COMPANY_USER_UPDATE_ERROR", "至少选择一个广告权限")
	}

	sroleIds := strings.Split(roleIds, ",")
	droleIds := make([]string, 0)

	for _, sroleId := range sroleIds {
		isExit := true

		for _, droleId := range droleIds {
			if droleId == sroleId {
				isExit = false
				break
			}
		}

		if isExit {
			droleIds = append(droleIds, sroleId)
		}
	}

	companyUser, err := cuuc.repo.UpdateRole(ctx, id, companyId, strings.Join(droleIds, ","))

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) DeleteCompanyUsers(ctx context.Context, id, companyId uint64) (*v1.DeleteCompanyUsersReply, error) {
	companyUser, err := cuuc.repo.Delete(ctx, id, companyId)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_USER_DELETE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyUser, nil
}
