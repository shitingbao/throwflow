package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/company/v1"
	"task/internal/biz"
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

func (cor *companyOrganizationRepo) SyncUpdateQrCode(ctx context.Context) (*v1.SyncUpdateQrCodeCompanyOrganizationsReply, error) {
	companyOrganization, err := cor.data.companyuc.SyncUpdateQrCodeCompanyOrganizations(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return companyOrganization, err
}
