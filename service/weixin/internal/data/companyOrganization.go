package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/company/v1"
	"weixin/internal/biz"
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

func (cor *companyOrganizationRepo) Get(ctx context.Context, organizationId uint64) (*v1.GetCompanyOrganizationsReply, error) {
	companyOrganization, err := cor.data.companyuc.GetCompanyOrganizations(ctx, &v1.GetCompanyOrganizationsRequest{
		OrganizationId: organizationId,
	})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}

func (cor *companyOrganizationRepo) GetByOrganizationCode(ctx context.Context, organizationCode string) (*v1.GetCompanyOrganizationByOrganizationCodesReply, error) {
	companyOrganization, err := cor.data.companyuc.GetCompanyOrganizationByOrganizationCodes(ctx, &v1.GetCompanyOrganizationByOrganizationCodesRequest{
		OrganizationCode: organizationCode,
	})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}

func (cor *companyOrganizationRepo) List(ctx context.Context) (*v1.ListCompanyOrganizationsReply, error) {
	list, err := cor.data.companyuc.ListCompanyOrganizations(ctx, &v1.ListCompanyOrganizationsRequest{
		PageNum:  0,
		PageSize: 40,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
