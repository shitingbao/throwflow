package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type performanceRepo struct {
	data *Data
	log  *log.Helper
}

func NewPerformanceRepo(data *Data, logger log.Logger) biz.PerformanceRepo {
	return &performanceRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (pr *performanceRepo) List(ctx context.Context, companyId, pageNum, pageSize uint64, updateDay string) (*v1.ListCompanyPerformancesReply, error) {
	list, err := pr.data.companyuc.ListCompanyPerformances(ctx, &v1.ListCompanyPerformancesRequest{
		PageNum:   pageNum,
		PageSize:  pageSize,
		CompanyId: companyId,
		UpdateDay: updateDay,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (pr *performanceRepo) Get(ctx context.Context, userId, companyId uint64, updateDay string) (*v1.GetCompanyPerformancesReply, error) {
	performance, err := pr.data.companyuc.GetCompanyPerformances(ctx, &v1.GetCompanyPerformancesRequest{
		UserId:    userId,
		CompanyId: companyId,
		UpdateDay: updateDay,
	})

	if err != nil {
		return nil, err
	}

	return performance, err
}
