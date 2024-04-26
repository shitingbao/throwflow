package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"time"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/pkg/tool"
)

func (ws *WeixinService) ListUserCommissions(ctx context.Context, in *v1.ListUserCommissionsRequest) (*v1.ListUserCommissionsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userCommissions, err := ws.ucuc.ListUserCommissions(ctx, in.PageNum, in.PageSize, in.UserId, in.OrganizationId, uint8(in.CommissionType), in.Month, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserCommissionsReply_UserCommission, 0)

	for _, userCommission := range userCommissions.List {
		list = append(list, &v1.ListUserCommissionsReply_UserCommission{
			NickName:           userCommission.ChildNickName,
			AvatarUrl:          userCommission.ChildAvatarUrl,
			Phone:              tool.FormatPhone(userCommission.ChildPhone),
			ActivationTime:     userCommission.ActivationTime,
			RelationName:       userCommission.RelationName,
			CommissionTypeName: userCommission.CommissionTypeName,
			TotalPayAmount:     tool.Decimal(float64(userCommission.TotalPayAmount), 2),
			CommissionPool:     tool.Decimal(float64(userCommission.CommissionPool), 2),
			CommissionRatio:    tool.Decimal(float64(userCommission.CommissionRatio), 2),
			CommissionAmount:   tool.Decimal(float64(userCommission.CommissionAmount), 2),
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

func (ws *WeixinService) StatisticsDetailUserCommissions(ctx context.Context, in *v1.StatisticsDetailUserCommissionsRequest) (*v1.StatisticsDetailUserCommissionsReply, error) {
	statisticsUserOrganizations, err := ws.ucuc.StatisticsDetailUserCommissions(ctx, in.UserId, in.OrganizationId, uint8(in.CommissionType), in.Month, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsDetailUserCommissionsReply_Statistic, 0)

	for _, statisticsUserOrganization := range statisticsUserOrganizations.Statistics {
		list = append(list, &v1.StatisticsDetailUserCommissionsReply_Statistic{
			Key:   statisticsUserOrganization.Key,
			Value: statisticsUserOrganization.Value,
		})
	}

	return &v1.StatisticsDetailUserCommissionsReply{
		Code: 200,
		Data: &v1.StatisticsDetailUserCommissionsReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ws *WeixinService) CreateOrderUserCommissions(ctx context.Context, in *v1.CreateOrderUserCommissionsRequest) (*v1.CreateOrderUserCommissionsReply, error) {
	if err := ws.ucuc.CreateOrderUserCommissions(ctx, in.TotalPayAmount, in.Commission, in.ClientKey, in.OpenId, in.OrderId, in.FlowPoint, in.PaySuccessTime); err != nil {
		return nil, err
	}

	return &v1.CreateOrderUserCommissionsReply{
		Code: 200,
		Data: &v1.CreateOrderUserCommissionsReply_Data{},
	}, nil
}

func (ws *WeixinService) CreateCostOrderUserCommissions(ctx context.Context, in *v1.CreateCostOrderUserCommissionsRequest) (*v1.CreateCostOrderUserCommissionsReply, error) {
	if err := ws.ucuc.CreateCostOrderUserCommissions(ctx, in.UserId, in.TotalPayAmount, in.Commission, in.OrderId, in.ProductId, in.FlowPoint, in.PaySuccessTime); err != nil {
		return nil, err
	}

	return &v1.CreateCostOrderUserCommissionsReply{
		Code: 200,
		Data: &v1.CreateCostOrderUserCommissionsReply_Data{},
	}, nil
}

func (ws *WeixinService) CreateTaskUserCommissions(ctx context.Context, in *v1.CreateTaskUserCommissionsRequest) (*v1.CreateTaskUserCommissionsReply, error) {
	if err := ws.ucuc.CreateTaskUserCommissions(ctx, in.UserId, in.TaskRelationId, in.Commission, in.FlowPoint, in.SuccessTime); err != nil {
		return nil, err
	}

	return &v1.CreateTaskUserCommissionsReply{
		Code: 200,
		Data: &v1.CreateTaskUserCommissionsReply_Data{},
	}, nil
}

func (ws *WeixinService) SyncTaskUserCommissions(ctx context.Context, in *emptypb.Empty) (*v1.SyncTaskUserCommissionsReply, error) {
	ws.log.Infof("同步微信用户任务分佣状态数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.ucuc.SyncTaskUserCommissions(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("同步微信用户任务分佣状态数据, 结束时间 %s \n", time.Now())

	return &v1.SyncTaskUserCommissionsReply{
		Code: 200,
		Data: &v1.SyncTaskUserCommissionsReply_Data{},
	}, nil
}
