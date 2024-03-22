package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/company/v1"
	"task/internal/pkg/tool"
)

type CompanyOrganizationRepo interface {
	SyncUpdateQrCode(context.Context) (*v1.SyncUpdateQrCodeCompanyOrganizationsReply, error)
}

type CompanyOrganizationUsecase struct {
	repo CompanyOrganizationRepo
	log  *log.Helper
}

func NewCompanyOrganizationUsecase(repo CompanyOrganizationRepo, logger log.Logger) *CompanyOrganizationUsecase {
	return &CompanyOrganizationUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (couc *CompanyOrganizationUsecase) SyncUpdateQrCodeCompanyOrganizations(ctx context.Context) (*v1.SyncUpdateQrCodeCompanyOrganizationsReply, error) {
	companyOrganization, err := couc.repo.SyncUpdateQrCode(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_UPDATE_QR_CODE_COMPANY_ORGANIZATION_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyOrganization, nil
}
