package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
	"interface/internal/pkg/tool"
	"time"
)

func (is *InterfaceService) ListUserOrders(ctx context.Context, in *v1.ListUserOrdersRequest) (*v1.ListUserOrdersReply, error) {
	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.InterfaceValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.InterfaceValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.InterfaceValidatorError
			}
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.InterfaceValidatorError
		}

		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.InterfaceValidatorError
			}
		}
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userOrders, err := is.jouc.ListUserOrders(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId, in.StartDay, in.EndDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserOrdersReply_Order, 0)

	for _, userOrder := range userOrders.Data.List {
		list = append(list, &v1.ListUserOrdersReply_Order{
			ProductId:          userOrder.ProductId,
			ProductName:        userOrder.ProductName,
			ProductImg:         userOrder.ProductImg,
			TotalPayAmount:     userOrder.TotalPayAmount,
			RealCommission:     userOrder.RealCommission,
			RealCommissionRate: userOrder.RealCommissionRate,
			ItemNum:            userOrder.ItemNum,
			MediaType:          userOrder.MediaType,
			MediaTypeName:      userOrder.MediaTypeName,
			MediaId:            userOrder.MediaId,
			MediaCover:         userOrder.MediaCover,
			Avatar:             userOrder.Avatar,
			IsShow:             userOrder.IsShow,
		})
	}

	return &v1.ListUserOrdersReply{
		Code: 200,
		Data: &v1.ListUserOrdersReply_Data{
			PageNum:   userOrders.Data.PageNum,
			PageSize:  userOrders.Data.PageSize,
			Total:     userOrders.Data.Total,
			TotalPage: userOrders.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) StatisticsUserOrders(ctx context.Context, in *v1.StatisticsUserOrdersRequest) (*v1.StatisticsUserOrdersReply, error) {
	startTime, err := tool.StringToTime("2006-01-02", in.StartDay)

	if err != nil {
		return nil, biz.InterfaceValidatorError
	}

	endTime, err := tool.StringToTime("2006-01-02", in.EndDay)

	if err != nil {
		return nil, biz.InterfaceValidatorError
	}

	nowTime, _ := tool.StringToTime("2006-01-02", tool.TimeToString("2006-01-02", time.Now()))

	if startTime.Equal(endTime) {
		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.InterfaceValidatorError
			}
		}
	} else {
		if !startTime.Before(endTime) {
			return nil, biz.InterfaceValidatorError
		}

		if !endTime.Equal(nowTime) {
			if !endTime.Before(nowTime) {
				return nil, biz.InterfaceValidatorError
			}
		}
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	statistics, err := is.jouc.StatisticsUserOrders(ctx, userInfo.Data.UserId, in.StartDay, in.EndDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsUserOrdersReply_Statistics, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsUserOrdersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsUserOrdersReply{
		Code: 200,
		Data: &v1.StatisticsUserOrdersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (is *InterfaceService) GetVideoUrlUserOrders(ctx context.Context, in *v1.GetVideoUrlUserOrdersRequest) (*v1.GetVideoUrlUserOrdersReply, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	videoUrl, err := is.muc.GetVideoUrls(ctx, in.VideoId)

	if err != nil {
		return nil, err
	}

	imageUrls := make([]*v1.GetVideoUrlUserOrdersReply_ImageUrls, 0)

	for _, imageUrl := range videoUrl.Data.ImageUrls {
		imageUrls = append(imageUrls, &v1.GetVideoUrlUserOrdersReply_ImageUrls{
			ImageUrl: imageUrl.ImageUrl,
		})
	}

	return &v1.GetVideoUrlUserOrdersReply{
		Code: 200,
		Data: &v1.GetVideoUrlUserOrdersReply_Data{
			VideoJumpUrl: videoUrl.Data.VideoJumpUrl,
			ImageUrls:    imageUrls,
		},
	}, nil
}
