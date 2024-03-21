package service

import (
	v1 "company/api/company/v1"
	"company/internal/pkg/tool"
	"context"
	"fmt"
)

func (cs *CompanyService) ListCompanyPerformanceRebalances(ctx context.Context, in *v1.ListCompanyPerformanceRebalancesRequest) (*v1.ListCompanyPerformanceRebalancesReply, error) {
	companyPerformanceRebalances, err := cs.cpreruc.ListCompanyPerformanceRebalances(ctx, in.UserId, in.CompanyId, in.UpdateDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyPerformanceRebalancesReply_CompanyPerformanceRebalance, 0)

	for _, companyPerformanceRebalance := range companyPerformanceRebalances {
		var cost string

		if companyPerformanceRebalance.RebalanceType == 1 {
			cost = fmt.Sprintf("+%.2f¥", companyPerformanceRebalance.Cost)
		} else {
			cost = fmt.Sprintf("-%.2f¥", companyPerformanceRebalance.Cost)
		}

		list = append(list, &v1.ListCompanyPerformanceRebalancesReply_CompanyPerformanceRebalance{
			Cost:       cost,
			Reason:     companyPerformanceRebalance.Reason,
			Ctype:      uint32(companyPerformanceRebalance.RebalanceType),
			UpdateTime: tool.TimeToString("2006/01/02 15:04", companyPerformanceRebalance.UpdateTime),
		})
	}

	return &v1.ListCompanyPerformanceRebalancesReply{
		Code: 200,
		Data: &v1.ListCompanyPerformanceRebalancesReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) CreateCompanyPerformanceRebalances(ctx context.Context, in *v1.CreateCompanyPerformanceRebalancesRequest) (*v1.CreateCompanyPerformanceRebalancesReply, error) {
	if _, err := cs.cpreruc.CreateCompanyPerformanceRebalances(ctx, in.UserId, in.CompanyId, in.Cost, uint8(in.Ctype), in.Reason); err != nil {
		return nil, err
	}

	return &v1.CreateCompanyPerformanceRebalancesReply{
		Code: 200,
		Data: &v1.CreateCompanyPerformanceRebalancesReply_Data{},
	}, nil

}
