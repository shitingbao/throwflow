package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"math"
	"strconv"
	"time"
)

func (ds *DouyinService) ListJinritemaiStores(ctx context.Context, in *v1.ListJinritemaiStoresRequest) (*v1.ListJinritemaiStoresReply, error) {
	jinritemaiStores, err := ds.jsuc.ListJinritemaiStores(ctx, in.PageNum, in.PageSize, in.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListJinritemaiStoresReply_JinritemaiStores, 0)

	for _, jinritemaiStore := range jinritemaiStores.List {
		list = append(list, &v1.ListJinritemaiStoresReply_JinritemaiStores{
			ClientKey:          jinritemaiStore.ClientKey,
			OpenId:             jinritemaiStore.OpenId,
			ProductId:          jinritemaiStore.ProductId,
			ProductName:        jinritemaiStore.ProductName,
			ProductImg:         jinritemaiStore.ProductImg,
			ProductPrice:       strconv.FormatFloat(float64(jinritemaiStore.ProductPrice), 'f', 2, 64),
			CommissionType:     uint32(jinritemaiStore.CommissionType),
			CommissionTypeName: jinritemaiStore.CommissionTypeName,
			CommissionRatio:    fmt.Sprintf("%.f%%", tool.Decimal(float64(jinritemaiStore.CommissionRatio), 0)),
			PromotionType:      uint32(jinritemaiStore.PromotionType),
			PromotionTypeName:  jinritemaiStore.PromotionTypeName,
		})
	}

	totalPage := uint64(math.Ceil(float64(jinritemaiStores.Total) / float64(jinritemaiStores.PageSize)))

	return &v1.ListJinritemaiStoresReply{
		Code: 200,
		Data: &v1.ListJinritemaiStoresReply_Data{
			PageNum:   jinritemaiStores.PageNum,
			PageSize:  jinritemaiStores.PageSize,
			Total:     jinritemaiStores.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListProductIdJinritemaiStores(ctx context.Context, in *v1.ListProductIdJinritemaiStoresRequest) (*v1.ListProductIdJinritemaiStoresReply, error) {
	jinritemaiStores, err := ds.jsuc.ListProductIdJinritemaiStores(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListProductIdJinritemaiStoresReply_JinritemaiStore, 0)

	for _, jinritemaiStore := range jinritemaiStores.List {
		list = append(list, &v1.ListProductIdJinritemaiStoresReply_JinritemaiStore{
			ProductId: jinritemaiStore.ProductId,
		})
	}

	totalPage := uint64(math.Ceil(float64(jinritemaiStores.Total) / float64(jinritemaiStores.PageSize)))

	return &v1.ListProductIdJinritemaiStoresReply{
		Code: 200,
		Data: &v1.ListProductIdJinritemaiStoresReply_Data{
			PageNum:   jinritemaiStores.PageNum,
			PageSize:  jinritemaiStores.PageSize,
			Total:     jinritemaiStores.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) CreateJinritemaiStores(ctx context.Context, in *v1.CreateJinritemaiStoresRequest) (*v1.CreateJinritemaiStoresReply, error) {
	messages, err := ds.jsuc.CreateJinritemaiStores(ctx, in.UserId, in.CompanyId, in.ProductId, in.OpenDouyinUserIds, in.ActivityUrl)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.CreateJinritemaiStoresReply_Message, 0)

	for _, message := range messages {
		list = append(list, &v1.CreateJinritemaiStoresReply_Message{
			ProductName: message.ProductName,
			AwemeName:   message.AwemeName,
			Content:     message.Content,
		})
	}

	return &v1.CreateJinritemaiStoresReply{
		Code: 200,
		Data: &v1.CreateJinritemaiStoresReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) UpdateJinritemaiStores(ctx context.Context, in *v1.UpdateJinritemaiStoresRequest) (*v1.UpdateJinritemaiStoresReply, error) {
	content, err := ds.jsuc.UpdateJinritemaiStores(ctx, in.UserId, in.Stores)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateJinritemaiStoresReply{
		Code: 200,
		Data: &v1.UpdateJinritemaiStoresReply_Data{
			Content: content,
		},
	}, nil
}

func (ds *DouyinService) DeleteJinritemaiStores(ctx context.Context, in *v1.DeleteJinritemaiStoresRequest) (*v1.DeleteJinritemaiStoresReply, error) {
	content, err := ds.jsuc.DeleteJinritemaiStores(ctx, in.UserId, in.Stores)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteJinritemaiStoresReply{
		Code: 200,
		Data: &v1.DeleteJinritemaiStoresReply_Data{
			Content: content,
		},
	}, nil
}

func (ds *DouyinService) SyncJinritemaiStores(ctx context.Context, in *empty.Empty) (*v1.SyncJinritemaiStoresReply, error) {
	ds.log.Infof("同步精选联盟达人橱窗商品数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.jsuc.SyncJinritemaiStores(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("同步精选联盟达人橱窗商品数据, 结束时间 %s \n", time.Now())

	return &v1.SyncJinritemaiStoresReply{
		Code: 200,
		Data: &v1.SyncJinritemaiStoresReply_Data{},
	}, nil
}
