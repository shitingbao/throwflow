package service

import (
	v1 "company/api/company/v1"
	"company/internal/biz"
	"company/internal/pkg/tool"
	"context"
	"fmt"
	"math"
	"strconv"
	"time"
	"unicode/utf8"
)

func (cs *CompanyService) ListCompanyPerformances(ctx context.Context, in *v1.ListCompanyPerformancesRequest) (*v1.ListCompanyPerformancesReply, error) {
	updateTime, err := tool.StringToTime("2006-01", in.UpdateDay)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	updateDay, err := strconv.ParseUint(tool.TimeToString("200601", updateTime), 10, 64)

	if err != nil {
		return nil, biz.CompanyValidatorError
	}

	companyPerformances, err := cs.cpuc.ListCompanyPerformances(ctx, in.PageNum, in.PageSize, in.CompanyId, uint32(updateDay))

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyPerformancesReply_CompanyPerformance, 0)

	for _, companyPerformance := range companyPerformances.List {
		advertisers := make([]*v1.ListCompanyPerformancesReply_Advertiser, 0)

		for _, ladvertiser := range companyPerformance.Advertiserst {
			advertisers = append(advertisers, &v1.ListCompanyPerformancesReply_Advertiser{
				AdvertiserId:   ladvertiser.AdvertiserId,
				AdvertiserName: ladvertiser.AdvertiserName,
			})
		}

		list = append(list, &v1.ListCompanyPerformancesReply_CompanyPerformance{
			UserId:               companyPerformance.UserId,
			Username:             companyPerformance.Username,
			Job:                  companyPerformance.Job,
			QianchuanAdvertisers: uint32(companyPerformance.QianchuanAdvertisers),
			StatCost:             fmt.Sprintf("%.2f", companyPerformance.StatCost),
			StatCostProportion:   fmt.Sprintf("%.2f%%", companyPerformance.StatCostProportion*100),
			Roi:                  fmt.Sprintf("%.2f", companyPerformance.Roi),
			Cost:                 fmt.Sprintf("%.2f", companyPerformance.Cost),
			RebalanceCost:        fmt.Sprintf("%.2f¥", companyPerformance.RebalanceCost),
			TotalCost:            fmt.Sprintf("%.2f", companyPerformance.TotalCost),
			Advertisers:          advertisers,
		})
	}

	totalPage := uint64(math.Ceil(float64(companyPerformances.Total) / float64(companyPerformances.PageSize)))

	return &v1.ListCompanyPerformancesReply{
		Code: 200,
		Data: &v1.ListCompanyPerformancesReply_Data{
			PageNum:   companyPerformances.PageNum,
			PageSize:  companyPerformances.PageSize,
			Total:     companyPerformances.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) GetCompanyPerformances(ctx context.Context, in *v1.GetCompanyPerformancesRequest) (*v1.GetCompanyPerformancesReply, error) {
	companyPerformances, err := cs.cpuc.GetCompanyPerformances(ctx, in.UserId, in.CompanyId, in.UpdateDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetCompanyPerformancesReply_CompanyPerformance, 0)

	for _, companyPerformance := range companyPerformances.List {
		rebalanceCosts := make([]*v1.GetCompanyPerformancesReply_RebalanceCost, 0)

		for _, rebalanceCost := range companyPerformance.RebalanceCosts {
			var cost string

			if rebalanceCost.RebalanceType == 1 {
				cost = fmt.Sprintf("+%.2f¥", rebalanceCost.Cost)
			} else {
				cost = fmt.Sprintf("-%.2f¥", rebalanceCost.Cost)
			}

			rebalanceCosts = append(rebalanceCosts, &v1.GetCompanyPerformancesReply_RebalanceCost{
				Cost:   cost,
				Reason: rebalanceCost.Reason,
			})
		}

		list = append(list, &v1.GetCompanyPerformancesReply_CompanyPerformance{
			QianchuanAdvertisers: uint32(len(companyPerformance.Advertisers)),
			StatCost:             fmt.Sprintf("%.2f", companyPerformance.StatCost),
			Roi:                  fmt.Sprintf("%.2f", companyPerformance.Roi),
			Cost:                 fmt.Sprintf("%.2f", companyPerformance.Cost),
			RebalanceCosts:       rebalanceCosts,
			UpdateDay:            tool.TimeToString("2006-01-02", companyPerformance.UpdateDay),
		})
	}

	return &v1.GetCompanyPerformancesReply{
		Code: 200,
		Data: &v1.GetCompanyPerformancesReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) RobotCompanyPerformances(ctx context.Context, in *v1.RobotCompanyPerformancesRequest) (*v1.RobotCompanyPerformancesReply, error) {
	ctx = context.Background()

	if l := utf8.RuneCountInString(in.UpdateDay); l == 0 {
		in.UpdateDay = tool.TimeToString("2006-01-02", time.Now().AddDate(0, 0, -1))
	}

	if err := cs.cpuc.RobotCompanyPerformances(ctx, in.UpdateDay); err != nil {
		return nil, err
	}

	return &v1.RobotCompanyPerformancesReply{
		Code: 200,
		Data: &v1.RobotCompanyPerformancesReply_Data{},
	}, nil
}
