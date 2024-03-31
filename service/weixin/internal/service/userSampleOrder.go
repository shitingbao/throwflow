package service

import (
	"context"
	"math"
	"strconv"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/pkg/tool"
)

func (ws *WeixinService) GetKuaidiInfoUserSampleOrders(ctx context.Context, in *v1.GetKuaidiInfoUserSampleOrdersRequest) (*v1.GetKuaidiInfoUserSampleOrdersReply, error) {
	kuaidiInfo, err := ws.usouc.GetKuaidiInfoUserSampleOrders(ctx, in.UserSampleOrderId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetKuaidiInfoUserSampleOrdersReply_KuaidiInfo, 0)

	for _, lkuaidiInfo := range kuaidiInfo.Data.List {
		list = append(list, &v1.GetKuaidiInfoUserSampleOrdersReply_KuaidiInfo{
			Time:    lkuaidiInfo.Time,
			Content: lkuaidiInfo.Content,
		})
	}

	return &v1.GetKuaidiInfoUserSampleOrdersReply{
		Code: 200,
		Data: &v1.GetKuaidiInfoUserSampleOrdersReply_Data{
			KuaidiCompany:   kuaidiInfo.Data.Name,
			KuaidiCode:      kuaidiInfo.Data.Code,
			KuaidiNum:       kuaidiInfo.Data.Num,
			KuaidiStateName: kuaidiInfo.Data.StateName,
			List:            list,
		},
	}, nil
}

func (ws *WeixinService) ListUserSampleOrders(ctx context.Context, in *v1.ListUserSampleOrdersRequest) (*v1.ListUserSampleOrdersReply, error) {
	userSampleOrders, err := ws.usouc.ListUserSampleOrders(ctx, in.PageNum, in.PageSize, in.Day, in.Keyword, in.SearchType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserSampleOrdersReply_OpenDouyinUser, 0)

	for _, userSampleOrder := range userSampleOrders.List {
		sampleOrders := make([]*v1.ListUserSampleOrdersReply_SampleOrder, 0)

		for _, luserSampleOrder := range userSampleOrder.UserSampleOrders {
			sampleOrders = append(sampleOrders, &v1.ListUserSampleOrdersReply_SampleOrder{
				UserSampleOrderId: luserSampleOrder.Id,
				ProductName:       luserSampleOrder.ProductName,
				ProductImg:        luserSampleOrder.ProductImg,
				OrderSn:           "NO." + strconv.FormatUint(luserSampleOrder.Id, 10) + "-" + luserSampleOrder.OrderSn,
				Name:              luserSampleOrder.Name,
				Phone:             luserSampleOrder.Phone,
				ProvinceAreaName:  luserSampleOrder.ProvinceAreaName,
				CityAreaName:      luserSampleOrder.CityAreaName,
				AreaAreaName:      luserSampleOrder.AreaAreaName,
				AddressInfo:       luserSampleOrder.AddressInfo,
				IsCancel:          uint32(luserSampleOrder.IsCancel),
				CancelNote:        luserSampleOrder.CancelNote,
				KuaidiCompany:     luserSampleOrder.KuaidiCompany,
				KuaidiCode:        luserSampleOrder.KuaidiCode,
				KuaidiNum:         luserSampleOrder.KuaidiNum,
				KuaidiStateName:   luserSampleOrder.KuaidiStateName,
				Note:              luserSampleOrder.Note,
				UpdateTime:        tool.TimeToString("2006-01-02 15:04", luserSampleOrder.CreateTime),
			})
		}

		list = append(list, &v1.ListUserSampleOrdersReply_OpenDouyinUser{
			OpenDouyinUserId: userSampleOrder.Id,
			AccountId:        userSampleOrder.AccountId,
			Nickname:         userSampleOrder.Nickname,
			Avatar:           userSampleOrder.Avatar,
			AvatarLarger:     userSampleOrder.AvatarLarger,
			Fans:             userSampleOrder.Fans,
			FansShow:         userSampleOrder.FansShow,
			SampleOrders:     sampleOrders,
		})
	}

	totalPage := uint64(math.Ceil(float64(userSampleOrders.Total) / float64(userSampleOrders.PageSize)))

	return &v1.ListUserSampleOrdersReply{
		Code: 200,
		Data: &v1.ListUserSampleOrdersReply_Data{
			PageNum:   userSampleOrders.PageNum,
			PageSize:  userSampleOrders.PageSize,
			Total:     userSampleOrders.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) StatisticsUserSampleOrders(ctx context.Context, in *v1.StatisticsUserSampleOrdersRequest) (*v1.StatisticsUserSampleOrdersReply, error) {
	statistics, err := ws.usouc.StatisticsUserSampleOrders(ctx, in.Day, in.Keyword, in.SearchType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsUserSampleOrdersReply_Statistic, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsUserSampleOrdersReply_Statistic{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsUserSampleOrdersReply{
		Code: 200,
		Data: &v1.StatisticsUserSampleOrdersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ws *WeixinService) VerifyUserSampleOrders(ctx context.Context, in *v1.VerifyUserSampleOrdersRequest) (*v1.VerifyUserSampleOrdersReply, error) {
	if err := ws.usouc.VerifyUserSampleOrders(ctx, in.UserId, in.OpenDouyinUserId, in.ProductId); err != nil {
		return nil, err
	}

	return &v1.VerifyUserSampleOrdersReply{
		Code: 200,
		Data: &v1.VerifyUserSampleOrdersReply_Data{},
	}, nil
}

func (ws *WeixinService) CancelUserSampleOrders(ctx context.Context, in *v1.CancelUserSampleOrdersRequest) (*v1.CancelUserSampleOrdersReply, error) {
	if err := ws.usouc.CancelUserSampleOrders(ctx, in.UserSampleOrderId, in.CancelNote); err != nil {
		return nil, err
	}

	return &v1.CancelUserSampleOrdersReply{
		Code: 200,
		Data: &v1.CancelUserSampleOrdersReply_Data{},
	}, nil
}

func (ws *WeixinService) CreateUserSampleOrders(ctx context.Context, in *v1.CreateUserSampleOrdersRequest) (*v1.CreateUserSampleOrdersReply, error) {
	if _, err := ws.usouc.CreateUserSampleOrders(ctx, in.UserId, in.OpenDouyinUserId, in.ProductId, in.UserAddressId, in.Note); err != nil {
		return nil, err
	}

	return &v1.CreateUserSampleOrdersReply{
		Code: 200,
		Data: &v1.CreateUserSampleOrdersReply_Data{},
	}, nil
}

func (ws *WeixinService) UpdateKuaidiUserSampleOrders(ctx context.Context, in *v1.UpdateKuaidiUserSampleOrdersRequest) (*v1.UpdateKuaidiUserSampleOrdersReply, error) {
	if err := ws.usouc.UpdateKuaidiUserSampleOrders(ctx, in.UserSampleOrderId, in.KuaidiCode, in.KuaidiNum); err != nil {
		return nil, err
	}

	return &v1.UpdateKuaidiUserSampleOrdersReply{
		Code: 200,
		Data: &v1.UpdateKuaidiUserSampleOrdersReply_Data{},
	}, nil
}
