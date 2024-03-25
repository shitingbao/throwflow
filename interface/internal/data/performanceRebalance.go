package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type performanceRebalanceRepo struct {
	data *Data
	log  *log.Helper
}

func NewPerformanceRebalanceRepo(data *Data, logger log.Logger) biz.PerformanceRebalanceRepo {
	return &performanceRebalanceRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (prr *performanceRebalanceRepo) List(ctx context.Context, userId, companyId uint64, updateDay string) (*v1.ListCompanyPerformanceRebalancesReply, error) {
	list, err := prr.data.companyuc.ListCompanyPerformanceRebalances(ctx, &v1.ListCompanyPerformanceRebalancesRequest{
		UserId:    userId,
		CompanyId: companyId,
		UpdateDay: updateDay,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (prr *performanceRebalanceRepo) Save(ctx context.Context, userId, companyId uint64, cost float32, ctype uint32, reason string) (*v1.CreateCompanyPerformanceRebalancesReply, error) {
	performanceRebalance, err := prr.data.companyuc.CreateCompanyPerformanceRebalances(ctx, &v1.CreateCompanyPerformanceRebalancesRequest{
		UserId:    userId,
		CompanyId: companyId,
		Cost:      cost,
		Ctype:     ctype,
		Reason:    reason,
	})

	if err != nil {
		return nil, err
	}

	return performanceRebalance, err
}
