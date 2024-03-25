package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type performanceRuleRepo struct {
	data *Data
	log  *log.Helper
}

func NewPerformanceRuleRepo(data *Data, logger log.Logger) biz.PerformanceRuleRepo {
	return &performanceRuleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (prr *performanceRuleRepo) List(ctx context.Context, companyId uint64) (*v1.ListCompanyPerformanceRulesReply, error) {
	list, err := prr.data.companyuc.ListCompanyPerformanceRules(ctx, &v1.ListCompanyPerformanceRulesRequest{
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (prr *performanceRuleRepo) ListQianchuanAdvertisers(ctx context.Context, id, companyId uint64) (*v1.ListQianchuanAdvertisersCompanyPerformanceRulesReply, error) {
	list, err := prr.data.companyuc.ListQianchuanAdvertisersCompanyPerformanceRules(ctx, &v1.ListQianchuanAdvertisersCompanyPerformanceRulesRequest{
		Id:        id,
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (prr *performanceRuleRepo) Save(ctx context.Context, companyId uint64, performanceName, advertiserIds, rules string) (*v1.CreateCompanyPerformanceRulesReply, error) {
	performanceRule, err := prr.data.companyuc.CreateCompanyPerformanceRules(ctx, &v1.CreateCompanyPerformanceRulesRequest{
		CompanyId:        companyId,
		PerformanceName:  performanceName,
		AdvertiserIds:    advertiserIds,
		PerformanceRules: rules,
	})

	if err != nil {
		return nil, err
	}

	return performanceRule, err
}

func (prr *performanceRuleRepo) Update(ctx context.Context, id, companyId uint64, performanceName, advertiserIds, rules string) (*v1.UpdateCompanyPerformanceRulesReply, error) {
	performanceRule, err := prr.data.companyuc.UpdateCompanyPerformanceRules(ctx, &v1.UpdateCompanyPerformanceRulesRequest{
		Id:               id,
		CompanyId:        companyId,
		PerformanceName:  performanceName,
		AdvertiserIds:    advertiserIds,
		PerformanceRules: rules,
	})

	if err != nil {
		return nil, err
	}

	return performanceRule, err
}
