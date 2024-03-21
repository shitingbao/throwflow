package service

import (
	v1 "company/api/company/v1"
	"context"
)

func (cs *CompanyService) ListCompanyPerformanceRules(ctx context.Context, in *v1.ListCompanyPerformanceRulesRequest) (*v1.ListCompanyPerformanceRulesReply, error) {
	companyPerformanceRules, err := cs.cpruc.ListCompanyPerformanceRules(ctx, in.CompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyPerformanceRulesReply_DPerformanceRule, 0)

	for _, lcompanyPerformanceRule := range companyPerformanceRules {
		advertisers := make([]*v1.ListCompanyPerformanceRulesReply_Advertiser, 0)

		for _, advertiser := range lcompanyPerformanceRule.Advertisers {
			advertisers = append(advertisers, &v1.ListCompanyPerformanceRulesReply_Advertiser{
				AdvertiserId:   advertiser.AdvertiserId,
				AdvertiserName: advertiser.AdvertiserName,
				CompanyName:    advertiser.CompanyName,
			})
		}

		performanceRules := make([]*v1.ListCompanyPerformanceRulesReply_PerformanceRule, 0)

		for _, performanceRule := range lcompanyPerformanceRule.PerformanceRulest {
			rules := make([]*v1.ListCompanyPerformanceRulesReply_Rule, 0)

			for _, rule := range performanceRule.Rules {
				rules = append(rules, &v1.ListCompanyPerformanceRulesReply_Rule{
					Condition: &v1.ListCompanyPerformanceRulesReply_Condition{
						Min:   rule.Condition.Min,
						Max:   rule.Condition.Max,
						Ctype: rule.Condition.Ctype,
					},
					Commission: &v1.ListCompanyPerformanceRulesReply_Commission{
						Ctype:      rule.Commission.Ctype,
						Percentage: rule.Commission.Percentage,
					},
				})
			}

			performanceRules = append(performanceRules, &v1.ListCompanyPerformanceRulesReply_PerformanceRule{
				Ptype: performanceRule.Ptype,
				Rules: rules,
			})
		}

		list = append(list, &v1.ListCompanyPerformanceRulesReply_DPerformanceRule{
			Id:               lcompanyPerformanceRule.Id,
			PerformanceName:  lcompanyPerformanceRule.PerformanceName,
			Advertisers:      advertisers,
			PerformanceRules: performanceRules,
		})
	}

	return &v1.ListCompanyPerformanceRulesReply{
		Code: 200,
		Data: &v1.ListCompanyPerformanceRulesReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListQianchuanAdvertisersCompanyPerformanceRules(ctx context.Context, in *v1.ListQianchuanAdvertisersCompanyPerformanceRulesRequest) (*v1.ListQianchuanAdvertisersCompanyPerformanceRulesReply, error) {
	qianchuanAdvertisersCompanyPerformanceRules, err := cs.cpruc.ListQianchuanAdvertisersCompanyPerformanceRules(ctx, in.Id, in.CompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdvertisersCompanyPerformanceRulesReply_Advertiser, 0)

	for _, qianchuanAdvertisersCompanyPerformanceRule := range qianchuanAdvertisersCompanyPerformanceRules {
		list = append(list, &v1.ListQianchuanAdvertisersCompanyPerformanceRulesReply_Advertiser{
			AdvertiserId:   qianchuanAdvertisersCompanyPerformanceRule.AdvertiserId,
			AdvertiserName: qianchuanAdvertisersCompanyPerformanceRule.AdvertiserName,
			CompanyName:    qianchuanAdvertisersCompanyPerformanceRule.CompanyName,
			IsSelect:       qianchuanAdvertisersCompanyPerformanceRule.IsSelect,
		})
	}

	return &v1.ListQianchuanAdvertisersCompanyPerformanceRulesReply{
		Code: 200,
		Data: &v1.ListQianchuanAdvertisersCompanyPerformanceRulesReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) CreateCompanyPerformanceRules(ctx context.Context, in *v1.CreateCompanyPerformanceRulesRequest) (*v1.CreateCompanyPerformanceRulesReply, error) {
	companyPerformanceRule, err := cs.cpruc.CreateCompanyPerformanceRules(ctx, in.CompanyId, in.PerformanceName, in.AdvertiserIds, in.PerformanceRules)

	if err != nil {
		return nil, err
	}

	advertisers := make([]*v1.CreateCompanyPerformanceRulesReply_Advertiser, 0)

	for _, advertiser := range companyPerformanceRule.Advertisers {
		advertisers = append(advertisers, &v1.CreateCompanyPerformanceRulesReply_Advertiser{
			AdvertiserId:   advertiser.AdvertiserId,
			AdvertiserName: advertiser.AdvertiserName,
			CompanyName:    advertiser.CompanyName,
		})
	}

	performanceRules := make([]*v1.CreateCompanyPerformanceRulesReply_PerformanceRule, 0)

	for _, performanceRule := range companyPerformanceRule.PerformanceRulest {
		rules := make([]*v1.CreateCompanyPerformanceRulesReply_Rule, 0)

		for _, rule := range performanceRule.Rules {
			rules = append(rules, &v1.CreateCompanyPerformanceRulesReply_Rule{
				Condition: &v1.CreateCompanyPerformanceRulesReply_Condition{
					Min:   rule.Condition.Min,
					Max:   rule.Condition.Max,
					Ctype: rule.Condition.Ctype,
				},
				Commission: &v1.CreateCompanyPerformanceRulesReply_Commission{
					Ctype:      rule.Commission.Ctype,
					Percentage: rule.Commission.Percentage,
				},
			})
		}

		performanceRules = append(performanceRules, &v1.CreateCompanyPerformanceRulesReply_PerformanceRule{
			Ptype: performanceRule.Ptype,
			Rules: rules,
		})
	}

	return &v1.CreateCompanyPerformanceRulesReply{
		Code: 200,
		Data: &v1.CreateCompanyPerformanceRulesReply_Data{
			Id:               companyPerformanceRule.Id,
			PerformanceName:  companyPerformanceRule.PerformanceName,
			Advertisers:      advertisers,
			PerformanceRules: performanceRules,
		},
	}, nil
}

func (cs *CompanyService) UpdateCompanyPerformanceRules(ctx context.Context, in *v1.UpdateCompanyPerformanceRulesRequest) (*v1.UpdateCompanyPerformanceRulesReply, error) {
	companyPerformanceRule, err := cs.cpruc.UpdateCompanyPerformanceRules(ctx, in.Id, in.CompanyId, in.PerformanceName, in.AdvertiserIds, in.PerformanceRules)

	if err != nil {
		return nil, err
	}

	advertisers := make([]*v1.UpdateCompanyPerformanceRulesReply_Advertiser, 0)

	for _, advertiser := range companyPerformanceRule.Advertisers {
		advertisers = append(advertisers, &v1.UpdateCompanyPerformanceRulesReply_Advertiser{
			AdvertiserId:   advertiser.AdvertiserId,
			AdvertiserName: advertiser.AdvertiserName,
			CompanyName:    advertiser.CompanyName,
		})
	}

	performanceRules := make([]*v1.UpdateCompanyPerformanceRulesReply_PerformanceRule, 0)

	for _, performanceRule := range companyPerformanceRule.PerformanceRulest {
		rules := make([]*v1.UpdateCompanyPerformanceRulesReply_Rule, 0)

		for _, rule := range performanceRule.Rules {
			rules = append(rules, &v1.UpdateCompanyPerformanceRulesReply_Rule{
				Condition: &v1.UpdateCompanyPerformanceRulesReply_Condition{
					Min:   rule.Condition.Min,
					Max:   rule.Condition.Max,
					Ctype: rule.Condition.Ctype,
				},
				Commission: &v1.UpdateCompanyPerformanceRulesReply_Commission{
					Ctype:      rule.Commission.Ctype,
					Percentage: rule.Commission.Percentage,
				},
			})
		}

		performanceRules = append(performanceRules, &v1.UpdateCompanyPerformanceRulesReply_PerformanceRule{
			Ptype: performanceRule.Ptype,
			Rules: rules,
		})
	}

	return &v1.UpdateCompanyPerformanceRulesReply{
		Code: 200,
		Data: &v1.UpdateCompanyPerformanceRulesReply_Data{
			Id:               companyPerformanceRule.Id,
			PerformanceName:  companyPerformanceRule.PerformanceName,
			Advertisers:      advertisers,
			PerformanceRules: performanceRules,
		},
	}, nil
}
