package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type CompanyOrganizationRepo interface {
	List(context.Context, uint64, uint64) (*v1.ListCompanyOrganizationsReply, error)
	ListSelect(context.Context) (*v1.ListSelectCompanyOrganizationsReply, error)
	Create(context.Context, string, string, string, string, string, string) (*v1.CreateCompanyOrganizationsReply, error)
	Update(context.Context, uint64, string, string, string, string, string, string) (*v1.UpdateCompanyOrganizationsReply, error)
	UpdateTeam(context.Context, uint64, string) (*v1.UpdateTeamCompanyOrganizationsReply, error)
	UpdateCommission(context.Context, uint64, string) (*v1.UpdateCommissionCompanyOrganizationsReply, error)
	UpdateColonelCommission(context.Context, uint64, string) (*v1.UpdateColonelCommissionCompanyOrganizationsReply, error)
	UpdateCourse(context.Context, uint64, string) (*v1.UpdateCourseCompanyOrganizationsReply, error)
}

type CompanyOrganizationUsecase struct {
	repo CompanyOrganizationRepo
	conf *conf.Data
	log  *log.Helper
}

func NewCompanyOrganizationUsecase(repo CompanyOrganizationRepo, conf *conf.Data, logger log.Logger) *CompanyOrganizationUsecase {
	return &CompanyOrganizationUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (couc *CompanyOrganizationUsecase) ListCompanyOrganizations(ctx context.Context, pageNum, pageSize uint64) (*v1.ListCompanyOrganizationsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(couc.conf.Database.PageSize)
	}

	list, err := couc.repo.List(ctx, pageNum, pageSize)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_COMPANY_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (couc *CompanyOrganizationUsecase) ListSelectCompanyOrganizations(ctx context.Context) (*v1.ListSelectCompanyOrganizationsReply, error) {
	list, err := couc.repo.ListSelect(ctx)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_SELECT_COMPANY_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (couc *CompanyOrganizationUsecase) CreateCompanyOrganizations(ctx context.Context, organizationName, organizationLogo, organizationMcn, companyName, bankCode, bankDeposit string) (*v1.CreateCompanyOrganizationsReply, error) {
	companyOrganization, err := couc.repo.Create(ctx, organizationName, organizationLogo, organizationMcn, companyName, bankCode, bankDeposit)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_COMPANY_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateCompanyOrganizations(ctx context.Context, organizationId uint64, organizationName, organizationLogo, organizationMcn, companyName, bankCode, bankDeposit string) (*v1.UpdateCompanyOrganizationsReply, error) {
	companyOrganization, err := couc.repo.Update(ctx, organizationId, organizationName, organizationLogo, organizationMcn, companyName, bankCode, bankDeposit)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_COMPANY_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateTeamCompanyOrganizations(ctx context.Context, organizationId uint64, organizationUser string) (*v1.UpdateTeamCompanyOrganizationsReply, error) {
	companyOrganization, err := couc.repo.UpdateTeam(ctx, organizationId, organizationUser)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_TEAM_COMPANY_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateCommissionCompanyOrganizations(ctx context.Context, organizationId uint64, organizationCommission string) (*v1.UpdateCommissionCompanyOrganizationsReply, error) {
	companyOrganization, err := couc.repo.UpdateCommission(ctx, organizationId, organizationCommission)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_COMMISSION_COMPANY_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateColonelCommissionCompanyOrganizations(ctx context.Context, organizationId uint64, organizationColonelCommission string) (*v1.UpdateColonelCommissionCompanyOrganizationsReply, error) {
	companyOrganization, err := couc.repo.UpdateColonelCommission(ctx, organizationId, organizationColonelCommission)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_COLONE_COMMISSION_COMPANY_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateCourseCompanyOrganizations(ctx context.Context, organizationId uint64, organizationCourse string) (*v1.UpdateCourseCompanyOrganizationsReply, error) {
	companyOrganization, err := couc.repo.UpdateCourse(ctx, organizationId, organizationCourse)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_COURSE_COMPANY_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companyOrganization, nil
}
