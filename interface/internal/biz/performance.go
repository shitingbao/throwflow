package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type PerformanceRepo interface {
	List(context.Context, uint64, uint64, uint64, string) (*v1.ListCompanyPerformancesReply, error)
	Get(context.Context, uint64, uint64, string) (*v1.GetCompanyPerformancesReply, error)
}

type PerformanceUsecase struct {
	repo PerformanceRepo
	conf *conf.Data
	log  *log.Helper
}

func NewPerformanceUsecase(repo PerformanceRepo, conf *conf.Data, logger log.Logger) *PerformanceUsecase {
	return &PerformanceUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (puc *PerformanceUsecase) ListPerformances(ctx context.Context, pageNum, pageSize, companyId uint64, updateDay string) (*v1.ListCompanyPerformancesReply, error) {
	if pageSize == 0 {
		pageSize = uint64(puc.conf.Database.PageSize)
	}

	performances, err := puc.repo.List(ctx, companyId, pageNum, pageSize, updateDay)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PERFORMANCE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return performances, nil
}

func (puc *PerformanceUsecase) DownPerformances(ctx context.Context, companyId uint64, updateDay string) (*v1.ListCompanyPerformancesReply, error) {
	performances, err := puc.repo.List(ctx, companyId, 0, uint64(puc.conf.Database.PageSize), updateDay)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PERFORMANCE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return performances, nil
}

func (puc *PerformanceUsecase) GetPerformances(ctx context.Context, userId, companyId uint64, updateDay string) (*v1.GetCompanyPerformancesReply, error) {
	performance, err := puc.repo.Get(ctx, userId, companyId, updateDay)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_PERFORMANCE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return performance, nil
}
