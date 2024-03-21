package biz

import (
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
)

var (
	CompanyCompanyPerformanceRuleNotFound    = errors.NotFound("COMPANY_COMPANY_PERFORMANCE_RULE_NOT_FOUND", "绩效方案不存在")
	CompanyCompanyPerformanceRuleCreateError = errors.InternalServer("COMPANY_COMPANY_PERFORMANCE_RULE_CREATE_ERROR", "绩效方案创建失败")
	CompanyCompanyPerformanceRuleUpdateError = errors.InternalServer("COMPANY_COMPANY_PERFORMANCE_RULE_UPDATE_ERROR", "绩效方案更新失败")
)

type CompanyPerformanceRuleRepo interface {
	GetById(context.Context, uint64, uint64) (*domain.CompanyPerformanceRule, error)
	List(context.Context, uint64) ([]*domain.CompanyPerformanceRule, error)
	ListAdvertiserIdsAndPerformanceRuleNotNull(context.Context, uint64) ([]*domain.CompanyPerformanceRule, error)
	Save(context.Context, *domain.CompanyPerformanceRule) (*domain.CompanyPerformanceRule, error)
	Update(context.Context, *domain.CompanyPerformanceRule) (*domain.CompanyPerformanceRule, error)
	DeleteByCompanyId(context.Context, uint64) error
}

type CompanyPerformanceRuleUsecase struct {
	repo   CompanyPerformanceRuleRepo
	crepo  CompanyRepo
	qarepo QianchuanAdvertiserRepo
	conf   *conf.Data
	log    *log.Helper
}

func NewCompanyPerformanceRuleUsecase(repo CompanyPerformanceRuleRepo, crepo CompanyRepo, qarepo QianchuanAdvertiserRepo, conf *conf.Data, logger log.Logger) *CompanyPerformanceRuleUsecase {
	return &CompanyPerformanceRuleUsecase{repo: repo, crepo: crepo, qarepo: qarepo, conf: conf, log: log.NewHelper(logger)}
}

func (cpruc *CompanyPerformanceRuleUsecase) ListCompanyPerformanceRules(ctx context.Context, companyId uint64) ([]*domain.CompanyPerformanceRule, error) {
	companyPerformanceRules, err := cpruc.repo.List(ctx, companyId)

	if err != nil {
		return nil, CompanyDataError
	}

	list := make([]*domain.CompanyPerformanceRule, 0)

	for _, companyPerformanceRule := range companyPerformanceRules {
		if companyPerformanceRule, err = cpruc.getCompanyPerformanceRule(ctx, companyPerformanceRule.Id, companyPerformanceRule.CompanyId); err == nil {
			list = append(list, companyPerformanceRule)
		}
	}

	return list, err
}

func (cpruc *CompanyPerformanceRuleUsecase) ListQianchuanAdvertisersCompanyPerformanceRules(ctx context.Context, id, companyId uint64) ([]*domain.QianchuanAdvertiserList, error) {
	if _, err := cpruc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	if id > 0 {
		if _, err := cpruc.repo.GetById(ctx, id, companyId); err != nil {
			return nil, CompanyCompanyPerformanceRuleNotFound
		}
	}

	qianchuanAdvertisers, err := cpruc.qarepo.List(ctx, companyId, 0, 0, "", "", "1")

	if err != nil {
		return nil, CompanyDataError
	}

	list := make([]*domain.QianchuanAdvertiserList, 0)

	advertiserIds := make([]string, 0)
	oAdvertiserIds := make([]string, 0)

	if companyPerformanceRules, err := cpruc.repo.List(ctx, companyId); err == nil {
		for _, acompanyPerformanceRule := range companyPerformanceRules {
			tAdvertiserIds := tool.RemoveEmptyString(strings.Split(acompanyPerformanceRule.AdvertiserIds, ","))

			if acompanyPerformanceRule.Id == id {
				advertiserIds = append(advertiserIds, tAdvertiserIds...)
			} else {
				oAdvertiserIds = append(oAdvertiserIds, tAdvertiserIds...)
			}
		}
	}

	for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
		lqianchuanAdvertiser := &domain.QianchuanAdvertiserList{
			AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
			Status:         qianchuanAdvertiser.Status,
			IsSelect:       0,
			AdvertiserName: qianchuanAdvertiser.AdvertiserName,
			CompanyName:    qianchuanAdvertiser.CompanyName,
		}

		for _, advertiserId := range advertiserIds {
			if strconv.FormatUint(lqianchuanAdvertiser.AdvertiserId, 10) == advertiserId {
				lqianchuanAdvertiser.IsSelect = 1
			}
		}

		for _, oAdvertiserId := range oAdvertiserIds {
			if strconv.FormatUint(lqianchuanAdvertiser.AdvertiserId, 10) == oAdvertiserId {
				lqianchuanAdvertiser.IsSelect = 2
			}
		}

		list = append(list, lqianchuanAdvertiser)
	}

	return list, nil
}

func (cpruc *CompanyPerformanceRuleUsecase) CreateCompanyPerformanceRules(ctx context.Context, companyId uint64, performanceName, advertiserIds, performanceRules string) (*domain.CompanyPerformanceRule, error) {
	if _, err := cpruc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	inCompanyPerformanceRule := domain.NewCompanyPerformanceRule(ctx, companyId, performanceName, performanceRules)
	inCompanyPerformanceRule.SetCreateTime(ctx)
	inCompanyPerformanceRule.SetUpdateTime(ctx)

	if ok := inCompanyPerformanceRule.VerifyPerformanceRules(ctx); !ok {
		return nil, CompanyDataError
	}

	qianchuanAdvertisers, err := cpruc.qarepo.List(ctx, companyId, 0, 0, "", "", "1")

	if err != nil {
		return nil, CompanyDataError
	}

	companyPerformanceRules, err := cpruc.repo.List(ctx, companyId)

	if err != nil {
		return nil, CompanyDataError
	}

	inCompanyPerformanceRule.SetAdvertiserIds(ctx, advertiserIds, qianchuanAdvertisers, companyPerformanceRules)

	companyPerformanceRule, err := cpruc.repo.Save(ctx, inCompanyPerformanceRule)

	if err != nil {
		return nil, CompanyCompanyPerformanceRuleCreateError
	}

	companyPerformanceRule, err = cpruc.getCompanyPerformanceRule(ctx, companyPerformanceRule.Id, companyPerformanceRule.CompanyId)

	if err != nil {
		return nil, CompanyCompanyPerformanceRuleCreateError
	}

	return companyPerformanceRule, nil
}

func (cpruc *CompanyPerformanceRuleUsecase) UpdateCompanyPerformanceRules(ctx context.Context, id, companyId uint64, performanceName, advertiserIds, performanceRules string) (*domain.CompanyPerformanceRule, error) {
	if _, err := cpruc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	inCompanyPerformanceRule, err := cpruc.repo.GetById(ctx, id, companyId)

	if err != nil {
		return nil, CompanyCompanyPerformanceRuleNotFound
	}

	inCompanyPerformanceRule.SetCompanyId(ctx, companyId)
	inCompanyPerformanceRule.SetPerformanceName(ctx, performanceName)
	inCompanyPerformanceRule.SetPerformanceRules(ctx, performanceRules)
	inCompanyPerformanceRule.SetUpdateTime(ctx)

	if ok := inCompanyPerformanceRule.VerifyPerformanceRules(ctx); !ok {
		return nil, CompanyDataError
	}

	qianchuanAdvertisers, err := cpruc.qarepo.List(ctx, companyId, 0, 0, "", "", "1")

	if err != nil {
		return nil, CompanyDataError
	}

	companyPerformanceRules, err := cpruc.repo.List(ctx, companyId)

	if err != nil {
		return nil, CompanyDataError
	}

	inCompanyPerformanceRule.SetAdvertiserIds(ctx, advertiserIds, qianchuanAdvertisers, companyPerformanceRules)

	companyPerformanceRule, err := cpruc.repo.Update(ctx, inCompanyPerformanceRule)

	if err != nil {
		return nil, CompanyCompanyPerformanceRuleUpdateError
	}

	companyPerformanceRule, err = cpruc.getCompanyPerformanceRule(ctx, companyPerformanceRule.Id, companyPerformanceRule.CompanyId)

	if err != nil {
		return nil, CompanyCompanyPerformanceRuleUpdateError
	}

	return companyPerformanceRule, nil
}

func (cpruc *CompanyPerformanceRuleUsecase) getCompanyPerformanceRule(ctx context.Context, id, companyId uint64) (*domain.CompanyPerformanceRule, error) {
	companyPerformanceRule, err := cpruc.repo.GetById(ctx, id, companyId)

	if err != nil {
		return nil, err
	}

	qianchuanAdvertisers, err := cpruc.qarepo.List(ctx, companyPerformanceRule.CompanyId, 0, 0, "", "", "1")

	if err != nil {
		return nil, err
	}

	companyPerformanceRule.GetAdvertiserIds(ctx, companyPerformanceRule.AdvertiserIds, qianchuanAdvertisers)
	companyPerformanceRule.GetPerformanceRules(ctx)

	return companyPerformanceRule, nil
}
