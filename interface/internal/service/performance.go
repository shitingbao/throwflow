package service

import (
	"context"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) ListPerformances(ctx context.Context, in *v1.ListPerformancesRequest) (*v1.ListPerformancesReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	performances, err := is.puc.ListPerformances(ctx, in.PageNum, in.PageSize, companyUser.Data.CurrentCompanyId, in.UpdateDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListPerformancesReply_CompanyPerformances, 0)

	for _, performance := range performances.Data.List {
		advertisers := make([]*v1.ListPerformancesReply_Advertisers, 0)

		for _, ladvertiser := range performance.Advertisers {
			advertisers = append(advertisers, &v1.ListPerformancesReply_Advertisers{
				AdvertiserId:   ladvertiser.AdvertiserId,
				AdvertiserName: ladvertiser.AdvertiserName,
			})
		}

		list = append(list, &v1.ListPerformancesReply_CompanyPerformances{
			UserId:               performance.UserId,
			Username:             performance.Username,
			Job:                  performance.Job,
			QianchuanAdvertisers: performance.QianchuanAdvertisers,
			StatCost:             performance.StatCost,
			StatCostProportion:   performance.StatCostProportion,
			Roi:                  performance.Roi,
			Cost:                 performance.Cost,
			RebalanceCost:        performance.RebalanceCost,
			TotalCost:            performance.TotalCost,
			Advertisers:          advertisers,
		})
	}

	return &v1.ListPerformancesReply{
		Code: 200,
		Data: &v1.ListPerformancesReply_Data{
			PageNum:   performances.Data.PageNum,
			PageSize:  performances.Data.PageSize,
			Total:     performances.Data.Total,
			TotalPage: performances.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) DownPerformances(ctx context.Context, in *v1.DownPerformancesRequest) (*v1.DownPerformancesReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	performances, err := is.puc.DownPerformances(ctx, companyUser.Data.CurrentCompanyId, in.UpdateDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.DownPerformancesReply_CompanyPerformances, 0)

	for _, performance := range performances.Data.List {
		advertisers := make([]*v1.DownPerformancesReply_Advertisers, 0)

		for _, ladvertiser := range performance.Advertisers {
			advertisers = append(advertisers, &v1.DownPerformancesReply_Advertisers{
				AdvertiserId:   ladvertiser.AdvertiserId,
				AdvertiserName: ladvertiser.AdvertiserName,
			})
		}

		list = append(list, &v1.DownPerformancesReply_CompanyPerformances{
			UserId:               performance.UserId,
			Username:             performance.Username,
			Job:                  performance.Job,
			QianchuanAdvertisers: performance.QianchuanAdvertisers,
			StatCost:             performance.StatCost,
			StatCostProportion:   performance.StatCostProportion,
			Roi:                  performance.Roi,
			Cost:                 performance.Cost,
			RebalanceCost:        performance.RebalanceCost,
			TotalCost:            performance.TotalCost,
			Advertisers:          advertisers,
		})
	}

	return &v1.DownPerformancesReply{
		Code: 200,
		Data: &v1.DownPerformancesReply_Data{
			List: list,
		},
	}, nil
}

func (is *InterfaceService) GetPerformances(ctx context.Context, in *v1.GetPerformancesRequest) (*v1.GetPerformancesReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "team")

	if err != nil {
		return nil, err
	}

	performance, err := is.puc.GetPerformances(ctx, in.UserId, companyUser.Data.CurrentCompanyId, in.UpdateDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetPerformancesReply_CompanyPerformances, 0)

	for _, lperformance := range performance.Data.List {
		rebalanceCosts := make([]*v1.GetPerformancesReply_RebalanceCosts, 0)

		for _, rebalanceCost := range lperformance.RebalanceCosts {
			rebalanceCosts = append(rebalanceCosts, &v1.GetPerformancesReply_RebalanceCosts{
				Cost:   rebalanceCost.Cost,
				Reason: rebalanceCost.Reason,
			})
		}

		list = append(list, &v1.GetPerformancesReply_CompanyPerformances{
			QianchuanAdvertisers: lperformance.QianchuanAdvertisers,
			StatCost:             lperformance.StatCost,
			Roi:                  lperformance.Roi,
			Cost:                 lperformance.Cost,
			RebalanceCosts:       rebalanceCosts,
			UpdateDay:            lperformance.UpdateDay,
		})
	}

	return &v1.GetPerformancesReply{
		Code: 200,
		Data: &v1.GetPerformancesReply_Data{
			List: list,
		},
	}, nil
}
