package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"fmt"
	"math"
	"time"
)

func (ds *DouyinService) ListQianchuanCampaigns(ctx context.Context, in *v1.ListQianchuanCampaignsRequest) (*v1.ListQianchuanCampaignsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	qianchuanCampaigns, err := ds.qcuc.ListQianchuanCampaigns(ctx, in.PageNum, in.PageSize, in.Day, in.Keyword, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanCampaignsReply_QianchuanCampaigns, 0)

	for _, qianchuanCampaign := range qianchuanCampaigns.List {
		list = append(list, &v1.ListQianchuanCampaignsReply_QianchuanCampaigns{
			Id:             qianchuanCampaign.Id,
			AdvertiserId:   qianchuanCampaign.AdvertiserId,
			Name:           qianchuanCampaign.Name,
			Budget:         fmt.Sprintf("%.2f", qianchuanCampaign.Budget),
			BudgetMode:     qianchuanCampaign.BudgetMode,
			MarketingGoal:  qianchuanCampaign.MarketingGoal,
			MarketingScene: qianchuanCampaign.MarketingScene,
		})
	}

	totalPage := uint64(math.Ceil(float64(qianchuanCampaigns.Total) / float64(qianchuanCampaigns.PageSize)))

	return &v1.ListQianchuanCampaignsReply{
		Code: 200,
		Data: &v1.ListQianchuanCampaignsReply_Data{
			PageNum:   qianchuanCampaigns.PageNum,
			PageSize:  qianchuanCampaigns.PageSize,
			Total:     qianchuanCampaigns.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}
