package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) ListPerformanceRules(ctx context.Context, in *emptypb.Empty) (*v1.ListPerformanceRulesReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	performanceRules, err := is.pruc.ListPerformanceRules(ctx, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListPerformanceRulesReply_DPerformanceRules, 0)

	for _, lperformanceRule := range performanceRules.Data.List {
		advertisers := make([]*v1.ListPerformanceRulesReply_Advertisers, 0)

		for _, advertiser := range lperformanceRule.Advertisers {
			advertisers = append(advertisers, &v1.ListPerformanceRulesReply_Advertisers{
				AdvertiserId:   advertiser.AdvertiserId,
				AdvertiserName: advertiser.AdvertiserName,
				CompanyName:    advertiser.CompanyName,
			})
		}

		dperformanceRules := make([]*v1.ListPerformanceRulesReply_PerformanceRules, 0)

		for _, performanceRule := range lperformanceRule.PerformanceRules {
			rules := make([]*v1.ListPerformanceRulesReply_Rules, 0)

			for _, rule := range performanceRule.Rules {
				rules = append(rules, &v1.ListPerformanceRulesReply_Rules{
					Condition: &v1.ListPerformanceRulesReply_Condition{
						Min:   rule.Condition.Min,
						Max:   rule.Condition.Max,
						Ctype: rule.Condition.Ctype,
					},
					Commission: &v1.ListPerformanceRulesReply_Commission{
						Ctype:      rule.Commission.Ctype,
						Percentage: rule.Commission.Percentage,
					},
				})
			}

			dperformanceRules = append(dperformanceRules, &v1.ListPerformanceRulesReply_PerformanceRules{
				Ptype: performanceRule.Ptype,
				Rules: rules,
			})
		}

		list = append(list, &v1.ListPerformanceRulesReply_DPerformanceRules{
			Id:               lperformanceRule.Id,
			PerformanceName:  lperformanceRule.PerformanceName,
			Advertisers:      advertisers,
			PerformanceRules: dperformanceRules,
		})
	}

	return &v1.ListPerformanceRulesReply{
		Code: 200,
		Data: &v1.ListPerformanceRulesReply_Data{
			List: list,
		},
	}, nil
}

func (is *InterfaceService) ListQianchuanAdvertisersPerformanceRules(ctx context.Context, in *v1.ListQianchuanAdvertisersPerformanceRulesRequest) (*v1.ListQianchuanAdvertisersPerformanceRulesReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	qianchuanAdvertisersPerformanceRules, err := is.pruc.ListQianchuanAdvertisersPerformanceRules(ctx, in.Id, companyUser.Data.CurrentCompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdvertisersPerformanceRulesReply_Advertisers, 0)

	for _, qianchuanAdvertisersPerformanceRule := range qianchuanAdvertisersPerformanceRules.Data.List {
		list = append(list, &v1.ListQianchuanAdvertisersPerformanceRulesReply_Advertisers{
			AdvertiserId:   qianchuanAdvertisersPerformanceRule.AdvertiserId,
			AdvertiserName: qianchuanAdvertisersPerformanceRule.AdvertiserName,
			CompanyName:    qianchuanAdvertisersPerformanceRule.CompanyName,
			IsSelect:       qianchuanAdvertisersPerformanceRule.IsSelect,
		})
	}

	return &v1.ListQianchuanAdvertisersPerformanceRulesReply{
		Code: 200,
		Data: &v1.ListQianchuanAdvertisersPerformanceRulesReply_Data{
			List: list,
		},
	}, nil
}

func (is *InterfaceService) CreatePerformanceRules(ctx context.Context, in *v1.CreatePerformanceRulesRequest) (*v1.CreatePerformanceRulesReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	performanceRule, err := is.pruc.CreatePerformanceRules(ctx, companyUser.Data.CurrentCompanyId, in.PerformanceName, in.AdvertiserIds, in.PerformanceRules)

	if err != nil {
		return nil, err
	}

	advertisers := make([]*v1.CreatePerformanceRulesReply_Advertisers, 0)

	for _, advertiser := range performanceRule.Data.Advertisers {
		advertisers = append(advertisers, &v1.CreatePerformanceRulesReply_Advertisers{
			AdvertiserId:   advertiser.AdvertiserId,
			AdvertiserName: advertiser.AdvertiserName,
			CompanyName:    advertiser.CompanyName,
		})
	}

	performanceRules := make([]*v1.CreatePerformanceRulesReply_PerformanceRules, 0)

	for _, lperformanceRule := range performanceRule.Data.PerformanceRules {
		rules := make([]*v1.CreatePerformanceRulesReply_Rules, 0)

		for _, rule := range lperformanceRule.Rules {
			rules = append(rules, &v1.CreatePerformanceRulesReply_Rules{
				Condition: &v1.CreatePerformanceRulesReply_Condition{
					Min:   rule.Condition.Min,
					Max:   rule.Condition.Max,
					Ctype: rule.Condition.Ctype,
				},
				Commission: &v1.CreatePerformanceRulesReply_Commission{
					Ctype:      rule.Commission.Ctype,
					Percentage: rule.Commission.Percentage,
				},
			})
		}

		performanceRules = append(performanceRules, &v1.CreatePerformanceRulesReply_PerformanceRules{
			Ptype: lperformanceRule.Ptype,
			Rules: rules,
		})
	}

	return &v1.CreatePerformanceRulesReply{
		Code: 200,
		Data: &v1.CreatePerformanceRulesReply_Data{
			Id:              performanceRule.Data.Id,
			PerformanceName: performanceRule.Data.PerformanceName,
			Advertisers:     advertisers,
		},
	}, nil
}

func (is *InterfaceService) UpdatePerformanceRules(ctx context.Context, in *v1.UpdatePerformanceRulesRequest) (*v1.UpdatePerformanceRulesReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	performanceRule, err := is.pruc.UpdatePerformanceRules(ctx, in.Id, companyUser.Data.CurrentCompanyId, in.PerformanceName, in.AdvertiserIds, in.PerformanceRules)

	if err != nil {
		return nil, err
	}

	advertisers := make([]*v1.UpdatePerformanceRulesReply_Advertisers, 0)

	for _, advertiser := range performanceRule.Data.Advertisers {
		advertisers = append(advertisers, &v1.UpdatePerformanceRulesReply_Advertisers{
			AdvertiserId:   advertiser.AdvertiserId,
			AdvertiserName: advertiser.AdvertiserName,
			CompanyName:    advertiser.CompanyName,
		})
	}

	performanceRules := make([]*v1.UpdatePerformanceRulesReply_PerformanceRules, 0)

	for _, lperformanceRule := range performanceRule.Data.PerformanceRules {
		rules := make([]*v1.UpdatePerformanceRulesReply_Rules, 0)

		for _, rule := range lperformanceRule.Rules {
			rules = append(rules, &v1.UpdatePerformanceRulesReply_Rules{
				Condition: &v1.UpdatePerformanceRulesReply_Condition{
					Min:   rule.Condition.Min,
					Max:   rule.Condition.Max,
					Ctype: rule.Condition.Ctype,
				},
				Commission: &v1.UpdatePerformanceRulesReply_Commission{
					Ctype:      rule.Commission.Ctype,
					Percentage: rule.Commission.Percentage,
				},
			})
		}

		performanceRules = append(performanceRules, &v1.UpdatePerformanceRulesReply_PerformanceRules{
			Ptype: lperformanceRule.Ptype,
			Rules: rules,
		})
	}

	return &v1.UpdatePerformanceRulesReply{
		Code: 200,
		Data: &v1.UpdatePerformanceRulesReply_Data{
			Id:               performanceRule.Data.Id,
			PerformanceName:  performanceRule.Data.PerformanceName,
			Advertisers:      advertisers,
			PerformanceRules: performanceRules,
		},
	}, nil
}
