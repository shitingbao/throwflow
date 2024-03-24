package biz

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"unicode/utf8"
)

type CompanyRepo interface {
	List(context.Context, uint64, uint64, string, string, uint32) (*v1.ListCompanysReply, error)
	ListSelect(context.Context) (*v1.ListSelectCompanysReply, error)
	Statistics(context.Context) (*v1.StatisticsCompanysReply, error)
	Save(context.Context, string, string, string, string, string, string, string, string, string, uint64, uint64, uint64, uint32, uint32, uint32) (*v1.CreateCompanysReply, error)
	Update(context.Context, uint64, string, string, string, string, string, string, string, string, uint64, uint64, uint32, uint32, uint32) (*v1.UpdateCompanysReply, error)
	UpdateStatus(context.Context, uint64, uint32) (*v1.UpdateStatusCompanysReply, error)
	UpdateRole(context.Context, uint64, uint64, string, string, string, uint32, uint32, uint32, uint32) (*v1.UpdateRoleCompanysReply, error)
	Delete(context.Context, uint64) (*v1.DeleteCompanysReply, error)
}

type CompanyUsecase struct {
	repo CompanyRepo
	log  *log.Helper
}

func NewCompanyUsecase(repo CompanyRepo, logger log.Logger) *CompanyUsecase {
	return &CompanyUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (cuc *CompanyUsecase) ListCompanys(ctx context.Context, pageNum, industryId uint64, keyword, status string, companyType uint32) (*v1.ListCompanysReply, error) {
	list, err := cuc.repo.List(ctx, pageNum, industryId, keyword, status, companyType)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cuc *CompanyUsecase) ListSelectCompanys(ctx context.Context) (*v1.ListSelectCompanysReply, error) {
	list, err := cuc.repo.ListSelect(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cuc *CompanyUsecase) StatisticsCompanys(ctx context.Context) (*v1.StatisticsCompanysReply, error) {
	list, err := cuc.repo.Statistics(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cuc *CompanyUsecase) CreateCompanys(ctx context.Context, companyName, contactInformation, seller, facilitator, adminName, adminPhone, address, industryId string, companyType, qianchuanUse uint32, userId, clueId, areaCode uint64) (*v1.CreateCompanysReply, error) {
	var status uint32 = 1
	source := "录入"

	if l := utf8.RuneCountInString(seller); l > 1 {
		status = 2
	}

	company, err := cuc.repo.Save(ctx, companyName, contactInformation, source, seller, facilitator, adminName, adminPhone, address, industryId, userId, clueId, areaCode, companyType, qianchuanUse, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_CREATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return company, nil
}

func (cuc *CompanyUsecase) UpdateCompanys(ctx context.Context, id uint64, companyName, contactInformation, seller, facilitator, adminName, adminPhone, address, industryId string, companyType, qianchuanUse uint32, clueId, areaCode uint64) (*v1.UpdateCompanysReply, error) {
	var status uint32 = 1

	if l := utf8.RuneCountInString(seller); l > 1 {
		status = 2
	}

	company, err := cuc.repo.Update(ctx, id, companyName, contactInformation, seller, facilitator, adminName, adminPhone, address, industryId, clueId, areaCode, companyType, qianchuanUse, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return company, nil
}

func (cuc *CompanyUsecase) UpdateStatusCompanys(ctx context.Context, id uint64, status uint32) (*v1.UpdateStatusCompanysReply, error) {
	company, err := cuc.repo.UpdateStatus(ctx, id, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return company, nil
}

func (cuc *CompanyUsecase) UpdateRoleCompanys(ctx context.Context, id, userId uint64, menuIds, startTime, endTime string, accounts, qianchuanAdvertisers, companyType, isTermwork uint32) (*v1.UpdateRoleCompanysReply, error) {
	company, err := cuc.repo.UpdateRole(ctx, id, userId, menuIds, startTime, endTime, accounts, qianchuanAdvertisers, companyType, isTermwork)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_ROLE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return company, nil
}

func (cuc *CompanyUsecase) DeleteCompanys(ctx context.Context, id uint64) (*v1.DeleteCompanysReply, error) {
	company, err := cuc.repo.Delete(ctx, id)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_ROLE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return company, nil
}
