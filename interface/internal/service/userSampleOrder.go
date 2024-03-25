package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) ListSampleOrders(ctx context.Context, in *v1.ListSampleOrdersRequest) (*v1.ListSampleOrdersReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	sampleOrders, err := is.usouc.ListSampleOrders(ctx, in.PageNum, in.PageSize, in.Day, in.Keyword, in.SearchType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListSampleOrdersReply_OpenDouyinUser, 0)

	for _, lsampleOrder := range sampleOrders.Data.List {
		externalSampleOrders := make([]*v1.ListSampleOrdersReply_SampleOrder, 0)

		for _, l := range lsampleOrder.SampleOrders {
			externalSampleOrders = append(externalSampleOrders, &v1.ListSampleOrdersReply_SampleOrder{
				UserSampleOrderId: l.UserSampleOrderId,
				ProductName:       l.ProductName,
				ProductImg:        l.ProductImg,
				OrderSn:           l.OrderSn,
				Name:              l.Name,
				Phone:             l.Phone,
				ProvinceAreaName:  l.ProvinceAreaName,
				CityAreaName:      l.CityAreaName,
				AreaAreaName:      l.AreaAreaName,
				AddressInfo:       l.AddressInfo,
				IsCancel:          l.IsCancel,
				CancelNote:        l.CancelNote,
				KuaidiCompany:     l.KuaidiCompany,
				KuaidiCode:        l.KuaidiCode,
				KuaidiNum:         l.KuaidiNum,
				KuaidiStateName:   l.KuaidiStateName,
				Note:              l.Note,
				UpdateTime:        l.UpdateTime,
			})
		}

		list = append(list, &v1.ListSampleOrdersReply_OpenDouyinUser{
			OpenDouyinUserId: lsampleOrder.OpenDouyinUserId,
			AccountId:        lsampleOrder.AccountId,
			Nickname:         lsampleOrder.Nickname,
			Avatar:           lsampleOrder.Avatar,
			AvatarLarger:     lsampleOrder.AvatarLarger,
			Fans:             lsampleOrder.Fans,
			FansShow:         lsampleOrder.FansShow,
			SampleOrders:     externalSampleOrders,
		})
	}

	return &v1.ListSampleOrdersReply{
		Code: 200,
		Data: &v1.ListSampleOrdersReply_Data{
			PageNum:   sampleOrders.Data.PageNum,
			PageSize:  sampleOrders.Data.PageSize,
			Total:     sampleOrders.Data.Total,
			TotalPage: sampleOrders.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) GetKuaidiInfoSampleOrders(ctx context.Context, in *v1.GetKuaidiInfoSampleOrdersRequest) (*v1.GetKuaidiInfoSampleOrdersReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	kuaidiInfoData, err := is.usouc.GetKuaidiInfoSampleOrders(ctx, in.UserSampleOrderId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetKuaidiInfoSampleOrdersReply_KuaidiInfo, 0)

	for _, kuaidiInfo := range kuaidiInfoData.Data.List {
		list = append(list, &v1.GetKuaidiInfoSampleOrdersReply_KuaidiInfo{
			Time:    kuaidiInfo.Time,
			Content: kuaidiInfo.Content,
		})
	}

	return &v1.GetKuaidiInfoSampleOrdersReply{
		Code: 200,
		Data: &v1.GetKuaidiInfoSampleOrdersReply_Data{
			KuaidiCompany:   kuaidiInfoData.Data.KuaidiCompany,
			KuaidiCode:      kuaidiInfoData.Data.KuaidiCode,
			KuaidiNum:       kuaidiInfoData.Data.KuaidiNum,
			KuaidiStateName: kuaidiInfoData.Data.KuaidiStateName,
			List:            list,
		},
	}, nil
}

func (is *InterfaceService) StatisticsSampleOrders(ctx context.Context, in *v1.StatisticsSampleOrdersRequest) (*v1.StatisticsSampleOrdersReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	statistics, err := is.usouc.StatisticsSampleOrders(ctx, in.Day, in.Keyword, in.SearchType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsSampleOrdersReply_Statistic, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsSampleOrdersReply_Statistic{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsSampleOrdersReply{
		Code: 200,
		Data: &v1.StatisticsSampleOrdersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (is *InterfaceService) VerifyUserSampleOrders(ctx context.Context, in *v1.VerifyUserSampleOrdersRequest) (*v1.VerifyUserSampleOrdersReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.usouc.VerifyUserSampleOrders(ctx, userInfo.Data.UserId, in.OpenDouyinUserId, in.ProductId); err != nil {
		return nil, err
	}

	return &v1.VerifyUserSampleOrdersReply{
		Code: 200,
		Data: &v1.VerifyUserSampleOrdersReply_Data{},
	}, nil
}

func (is *InterfaceService) CancelSampleOrders(ctx context.Context, in *v1.CancelSampleOrdersRequest) (*v1.CancelSampleOrdersReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	if _, err := is.usouc.CancelSampleOrders(ctx, in.UserSampleOrderId, in.CancelNote); err != nil {
		return nil, err
	}

	return &v1.CancelSampleOrdersReply{
		Code: 200,
		Data: &v1.CancelSampleOrdersReply_Data{},
	}, nil
}

func (is *InterfaceService) CreateUserSampleOrders(ctx context.Context, in *v1.CreateUserSampleOrdersRequest) (*v1.CreateUserSampleOrdersReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.usouc.CreateUserSampleOrders(ctx, userInfo.Data.UserId, in.OpenDouyinUserId, in.ProductId, in.UserAddressId, in.Note); err != nil {
		return nil, err
	}

	return &v1.CreateUserSampleOrdersReply{
		Code: 200,
		Data: &v1.CreateUserSampleOrdersReply_Data{},
	}, nil
}

func (is *InterfaceService) UpdateKuaidiSampleOrders(ctx context.Context, in *v1.UpdateKuaidiSampleOrdersRequest) (*v1.UpdateKuaidiSampleOrdersReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	if _, err := is.usouc.UpdateKuaidiSampleOrders(ctx, in.UserSampleOrderId, in.KuaidiCode, in.KuaidiNum); err != nil {
		return nil, err
	}

	return &v1.UpdateKuaidiSampleOrdersReply{
		Code: 200,
		Data: &v1.UpdateKuaidiSampleOrdersReply_Data{},
	}, nil
}
