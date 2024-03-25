package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type PerformanceRuleRepo interface {
	List(context.Context, uint64) (*v1.ListCompanyPerformanceRulesReply, error)
	ListQianchuanAdvertisers(context.Context, uint64, uint64) (*v1.ListQianchuanAdvertisersCompanyPerformanceRulesReply, error)
	Save(context.Context, uint64, string, string, string) (*v1.CreateCompanyPerformanceRulesReply, error)
	Update(context.Context, uint64, uint64, string, string, string) (*v1.UpdateCompanyPerformanceRulesReply, error)
}

type PerformanceRuleUsecase struct {
	repo PerformanceRuleRepo
	conf *conf.Data
	log  *log.Helper
}

func NewPerformanceRuleUsecase(repo PerformanceRuleRepo, conf *conf.Data, logger log.Logger) *PerformanceRuleUsecase {
	return &PerformanceRuleUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (pruc *PerformanceRuleUsecase) ListPerformanceRules(ctx context.Context, companyId uint64) (*v1.ListCompanyPerformanceRulesReply, error) {
	performanceRules, err := pruc.repo.List(ctx, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PERFORMANCE_RULE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return performanceRules, nil
}

func (pruc *PerformanceRuleUsecase) ListQianchuanAdvertisersPerformanceRules(ctx context.Context, id, companyId uint64) (*v1.ListQianchuanAdvertisersCompanyPerformanceRulesReply, error) {
	qianchuanAdvertisersPerformanceRules, err := pruc.repo.ListQianchuanAdvertisers(ctx, id, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_QIANCHUAN_ADVERTISERS_PERFORMANCE_RULE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return qianchuanAdvertisersPerformanceRules, nil
}

func (pruc *PerformanceRuleUsecase) CreatePerformanceRules(ctx context.Context, companyId uint64, performanceName, advertiserIds, rules string) (*v1.CreateCompanyPerformanceRulesReply, error) {
	performanceRule, err := pruc.repo.Save(ctx, companyId, performanceName, advertiserIds, rules)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_PERFORMANCE_RULE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return performanceRule, nil
}

func (pruc *PerformanceRuleUsecase) UpdatePerformanceRules(ctx context.Context, id, companyId uint64, performanceName, advertiserIds, rules string) (*v1.UpdateCompanyPerformanceRulesReply, error) {
	performanceRule, err := pruc.repo.Update(ctx, id, companyId, performanceName, advertiserIds, rules)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_PERFORMANCE_RULE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return performanceRule, nil
}
