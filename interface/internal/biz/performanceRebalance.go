package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type PerformanceRebalanceRepo interface {
	List(context.Context, uint64, uint64, string) (*v1.ListCompanyPerformanceRebalancesReply, error)
	Save(context.Context, uint64, uint64, float32, uint32, string) (*v1.CreateCompanyPerformanceRebalancesReply, error)
}

type PerformanceRebalanceUsecase struct {
	repo PerformanceRebalanceRepo
	conf *conf.Data
	log  *log.Helper
}

func NewPerformanceRebalanceUsecase(repo PerformanceRebalanceRepo, conf *conf.Data, logger log.Logger) *PerformanceRebalanceUsecase {
	return &PerformanceRebalanceUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (pruc *PerformanceRebalanceUsecase) ListPerformanceRebalances(ctx context.Context, userId, companyId uint64, updateDay string) (*v1.ListCompanyPerformanceRebalancesReply, error) {
	performanceRebalances, err := pruc.repo.List(ctx, userId, companyId, updateDay)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PERFORMANCE_REBALANCE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return performanceRebalances, nil
}

func (pruc *PerformanceRebalanceUsecase) CreatePerformanceRebalances(ctx context.Context, userId, companyId uint64, cost float32, ctype uint32, reason string) (*v1.CreateCompanyPerformanceRebalancesReply, error) {
	performanceRebalance, err := pruc.repo.Save(ctx, userId, companyId, cost, ctype, reason)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_PERFORMANCE_REBALANCE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return performanceRebalance, nil
}
