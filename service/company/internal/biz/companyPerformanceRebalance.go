package biz

import (
	"company/internal/conf"
	"company/internal/domain"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
)

var (
	CompanyCompanyPerformanceRebalanceNotFound    = errors.NotFound("COMPANY_COMPANY_PERFORMANCE_REBALANCE_NOT_FOUND", "绩效奖罚不存在")
	CompanyCompanyPerformanceRebalanceCreateError = errors.InternalServer("COMPANY_COMPANY_PERFORMANCE_REBALANCE_CREATE_ERROR", "绩效奖罚创建失败")
)

type CompanyPerformanceRebalanceRepo interface {
	List(context.Context, uint64, uint64, string) ([]*domain.CompanyPerformanceRebalance, error)
	Save(context.Context, *domain.CompanyPerformanceRebalance) (*domain.CompanyPerformanceRebalance, error)
	DeleteByCompanyId(context.Context, uint64) error
}

type CompanyPerformanceRebalanceUsecase struct {
	repo    CompanyPerformanceRebalanceRepo
	cpmrepo CompanyPerformanceMonthlyRepo
	curepo  CompanyUserRepo
	tm      Transaction
	conf    *conf.Data
	log     *log.Helper
}

func NewCompanyPerformanceRebalanceUsecase(repo CompanyPerformanceRebalanceRepo, cpmrepo CompanyPerformanceMonthlyRepo, curepo CompanyUserRepo, tm Transaction, conf *conf.Data, logger log.Logger) *CompanyPerformanceRebalanceUsecase {
	return &CompanyPerformanceRebalanceUsecase{repo: repo, cpmrepo: cpmrepo, curepo: curepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (cpruc *CompanyPerformanceRebalanceUsecase) ListCompanyPerformanceRebalances(ctx context.Context, userId, companyId uint64, updateDay string) ([]*domain.CompanyPerformanceRebalance, error) {
	if _, err := cpruc.curepo.GetById(ctx, userId, companyId); err != nil {
		return nil, CompanyCompanyUserNotFound
	}

	list, err := cpruc.repo.List(ctx, userId, companyId, updateDay)

	if err != nil {
		return nil, CompanyCompanyPerformanceRebalanceNotFound
	}

	return list, nil
}

func (cpruc *CompanyPerformanceRebalanceUsecase) CreateCompanyPerformanceRebalances(ctx context.Context, userId, companyId uint64, cost float32, ctype uint8, reason string) (*domain.CompanyPerformanceRebalance, error) {
	if _, err := cpruc.curepo.GetById(ctx, userId, companyId); err != nil {
		return nil, CompanyCompanyUserNotFound
	}

	var companyPerformanceRebalance *domain.CompanyPerformanceRebalance

	err := cpruc.tm.InTx(ctx, func(ctx context.Context) error {
		var err error

		inCompanyPerformanceRebalance := domain.NewCompanyPerformanceRebalance(ctx, userId, companyId, cost, ctype, reason)
		inCompanyPerformanceRebalance.SetUpdateDay(ctx)
		inCompanyPerformanceRebalance.SetCreateTime(ctx)
		inCompanyPerformanceRebalance.SetUpdateTime(ctx)

		companyPerformanceRebalance, err = cpruc.repo.Save(ctx, inCompanyPerformanceRebalance)

		if err != nil {
			return err
		}

		currentMonth := fmt.Sprintf("%d", time.Now().Year()) + time.Now().Format("01")
		uicurrentMonth, _ := strconv.ParseUint(currentMonth, 10, 64)

		if inCompanyPerformanceMonthly, err := cpruc.cpmrepo.GetByUserIdAndCompanyIdAndUpdateDay(ctx, userId, companyId, uint32(uicurrentMonth)); err == nil {
			var rebalanceCost float32
			var totalCost float32

			if ctype == 1 {
				rebalanceCost = inCompanyPerformanceMonthly.RebalanceCost + cost
				totalCost = inCompanyPerformanceMonthly.TotalCost + cost
			} else {
				rebalanceCost = inCompanyPerformanceMonthly.RebalanceCost - cost
				totalCost = inCompanyPerformanceMonthly.TotalCost - cost
			}

			inCompanyPerformanceMonthly.SetRebalanceCost(ctx, rebalanceCost)
			inCompanyPerformanceMonthly.SetTotalCost(ctx, totalCost)

			if _, err := cpruc.cpmrepo.Update(ctx, inCompanyPerformanceMonthly); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, CompanyCompanyPerformanceRebalanceCreateError
	}

	return companyPerformanceRebalance, nil
}
