package service

import (
	"context"
	"math"
	"time"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/pkg/tool"
)

func (ws *WeixinService) ListUserCommissions(ctx context.Context, in *v1.ListUserCommissionsRequest) (*v1.ListUserCommissionsReply, error) {
	userCommissions, err := ws.ucuc.ListUserCommissions(ctx, in.PageNum, in.PageSize, in.UserId, in.OrganizationId, uint8(in.IsDirect), in.Month, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserCommissionsReply_UserCommission, 0)

	for _, userCommission := range userCommissions.List {
		list = append(list, &v1.ListUserCommissionsReply_UserCommission{
			NickName:                userCommission.ChildNickName,
			AvatarUrl:               userCommission.ChildAvatarUrl,
			Phone:                   tool.FormatPhone(userCommission.ChildPhone),
			ActivationTime:          userCommission.ActivationTime,
			RelationName:            userCommission.RelationName,
			TotalPayAmount:          tool.Decimal(float64(userCommission.TotalPayAmount), 2),
			CommissionPool:          tool.Decimal(float64(userCommission.CommissionPool), 2),
			EstimatedUserCommission: tool.Decimal(float64(userCommission.EstimatedUserCommission), 2),
			CommissionRatio:         tool.Decimal(float64(userCommission.CommissionRatio), 2),
			RealUserCommission:      tool.Decimal(float64(userCommission.RealUserCommission), 2),
			UserCommissionTypeName:  userCommission.UserCommissionTypeName,
		})
	}

	totalPage := uint64(math.Ceil(float64(userCommissions.Total) / float64(userCommissions.PageSize)))

	return &v1.ListUserCommissionsReply{
		Code: 200,
		Data: &v1.ListUserCommissionsReply_Data{
			PageNum:   userCommissions.PageNum,
			PageSize:  userCommissions.PageSize,
			Total:     userCommissions.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) StatisticsUserCommissions(ctx context.Context, in *v1.StatisticsUserCommissionsRequest) (*v1.StatisticsUserCommissionsReply, error) {
	statisticsUserOrganizations, err := ws.ucuc.StatisticsUserCommissions(ctx, in.UserId, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsUserCommissionsReply_Statistic, 0)

	for _, statisticsUserOrganization := range statisticsUserOrganizations.Statistics {
		list = append(list, &v1.StatisticsUserCommissionsReply_Statistic{
			Key:   statisticsUserOrganization.Key,
			Value: statisticsUserOrganization.Value,
		})
	}

	return &v1.StatisticsUserCommissionsReply{
		Code: 200,
		Data: &v1.StatisticsUserCommissionsReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ws *WeixinService) SyncOrderUserCommissions(ctx context.Context, in *v1.SyncOrderUserCommissionsRequest) (*v1.SyncOrderUserCommissionsReply, error) {
	ws.log.Infof("同步微信用户电商分佣数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.ucuc.SyncOrderUserCommissions(ctx, in.Day); err != nil {
		return nil, err
	}

	ws.log.Infof("同步微信用户电商分佣数据, 结束时间 %s \n", time.Now())

	return &v1.SyncOrderUserCommissionsReply{
		Code: 200,
		Data: &v1.SyncOrderUserCommissionsReply_Data{},
	}, nil
}

func (ws *WeixinService) SyncCostOrderUserCommissions(ctx context.Context, in *v1.SyncCostOrderUserCommissionsRequest) (*v1.SyncCostOrderUserCommissionsReply, error) {
	ws.log.Infof("同步微信用户成本购分佣数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.ucuc.SyncCostOrderUserCommissions(ctx, in.Day); err != nil {
		return nil, err
	}

	ws.log.Infof("同步微信用户成本购分佣数据, 结束时间 %s \n", time.Now())

	return &v1.SyncCostOrderUserCommissionsReply{
		Code: 200,
		Data: &v1.SyncCostOrderUserCommissionsReply_Data{},
	}, nil
}
