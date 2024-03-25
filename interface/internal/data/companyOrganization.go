package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type companyOrganizationRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyOrganizationRepo(data *Data, logger log.Logger) biz.CompanyOrganizationRepo {
	return &companyOrganizationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cor *companyOrganizationRepo) List(ctx context.Context, pageNum, pageSize uint64) (*v1.ListCompanyOrganizationsReply, error) {
	list, err := cor.data.companyuc.ListCompanyOrganizations(ctx, &v1.ListCompanyOrganizationsRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cor *companyOrganizationRepo) ListSelect(ctx context.Context) (*v1.ListSelectCompanyOrganizationsReply, error) {
	list, err := cor.data.companyuc.ListSelectCompanyOrganizations(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cor *companyOrganizationRepo) Create(ctx context.Context, organizationName, organizationLogo, organizationMcn, companyName, bankCode, bankDeposit string) (*v1.CreateCompanyOrganizationsReply, error) {
	companyOrganization, err := cor.data.companyuc.CreateCompanyOrganizations(ctx, &v1.CreateCompanyOrganizationsRequest{
		OrganizationName: organizationName,
		OrganizationLogo: organizationLogo,
		OrganizationMcn:  organizationMcn,
		CompanyName:      companyName,
		BankCode:         bankCode,
		BankDeposit:      bankDeposit,
	})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}

func (cor *companyOrganizationRepo) Update(ctx context.Context, organizationId uint64, organizationName, organizationLogo, organizationMcn, companyName, bankCode, bankDeposit string) (*v1.UpdateCompanyOrganizationsReply, error) {
	companyOrganization, err := cor.data.companyuc.UpdateCompanyOrganizations(ctx, &v1.UpdateCompanyOrganizationsRequest{
		OrganizationId:   organizationId,
		OrganizationName: organizationName,
		OrganizationLogo: organizationLogo,
		OrganizationMcn:  organizationMcn,
		CompanyName:      companyName,
		BankCode:         bankCode,
		BankDeposit:      bankDeposit,
	})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}

func (cor *companyOrganizationRepo) UpdateTeam(ctx context.Context, organizationId uint64, organizationUser string) (*v1.UpdateTeamCompanyOrganizationsReply, error) {
	companyOrganization, err := cor.data.companyuc.UpdateTeamCompanyOrganizations(ctx, &v1.UpdateTeamCompanyOrganizationsRequest{
		OrganizationId:   organizationId,
		OrganizationUser: organizationUser,
	})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}

func (cor *companyOrganizationRepo) UpdateCommission(ctx context.Context, organizationId uint64, organizationCommission string) (*v1.UpdateCommissionCompanyOrganizationsReply, error) {
	companyOrganization, err := cor.data.companyuc.UpdateCommissionCompanyOrganizations(ctx, &v1.UpdateCommissionCompanyOrganizationsRequest{
		OrganizationId:         organizationId,
		OrganizationCommission: organizationCommission,
	})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}

func (cor *companyOrganizationRepo) UpdateColonelCommission(ctx context.Context, organizationId uint64, organizationColonelCommission string) (*v1.UpdateColonelCommissionCompanyOrganizationsReply, error) {
	companyOrganization, err := cor.data.companyuc.UpdateColonelCommissionCompanyOrganizations(ctx, &v1.UpdateColonelCommissionCompanyOrganizationsRequest{
		OrganizationId:                organizationId,
		OrganizationColonelCommission: organizationColonelCommission,
	})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}

func (cor *companyOrganizationRepo) UpdateCourse(ctx context.Context, organizationId uint64, organizationCourse string) (*v1.UpdateCourseCompanyOrganizationsReply, error) {
	companyOrganization, err := cor.data.companyuc.UpdateCourseCompanyOrganizations(ctx, &v1.UpdateCourseCompanyOrganizationsRequest{
		OrganizationId:     organizationId,
		OrganizationCourse: organizationCourse,
	})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}
