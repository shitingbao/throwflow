package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) ListPerformanceRebalances(ctx context.Context, in *v1.ListPerformanceRebalancesRequest) (*v1.ListPerformanceRebalancesReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	performanceRebalances, err := is.preuc.ListPerformanceRebalances(ctx, in.UserId, companyUser.Data.CurrentCompanyId, in.UpdateDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListPerformanceRebalancesReply_CompanyPerformanceRebalances, 0)

	for _, performanceRebalance := range performanceRebalances.Data.List {
		list = append(list, &v1.ListPerformanceRebalancesReply_CompanyPerformanceRebalances{
			Cost:       performanceRebalance.Cost,
			Reason:     performanceRebalance.Reason,
			Ctype:      performanceRebalance.Ctype,
			UpdateTime: performanceRebalance.UpdateTime,
		})
	}

	return &v1.ListPerformanceRebalancesReply{
		Code: 200,
		Data: &v1.ListPerformanceRebalancesReply_Data{
			List: list,
		},
	}, nil
}

func (is *InterfaceService) CreatePerformanceRebalances(ctx context.Context, in *v1.CreatePerformanceRebalancesRequest) (*v1.CreatePerformanceRebalancesReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	if _, err := is.preuc.CreatePerformanceRebalances(ctx, in.UserId, companyUser.Data.CurrentCompanyId, in.Cost, in.Ctype, in.Reason); err != nil {
		return nil, err
	}

	return &v1.CreatePerformanceRebalancesReply{
		Code: 200,
		Data: &v1.CreatePerformanceRebalancesReply_Data{},
	}, nil
}
