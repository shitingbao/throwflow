package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
)

func (ds *DouyinService) ListOceanengineConfigs(ctx context.Context, in *v1.ListOceanengineConfigsRequest) (*v1.ListOceanengineConfigsReply, error) {
	oceanengineConfigs, err := ds.ocuc.ListOceanengineConfigs(ctx, uint8(in.OceanengineType), in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListOceanengineConfigsReply_Oceanengines, 0)

	for _, qianchuanConfig := range oceanengineConfigs.List {
		list = append(list, &v1.ListOceanengineConfigsReply_Oceanengines{
			Id:              qianchuanConfig.Id,
			OceanengineType: uint32(qianchuanConfig.OceanengineType),
			AppId:           qianchuanConfig.AppId,
			AppName:         qianchuanConfig.AppName,
			AppSecret:       qianchuanConfig.AppSecret,
			RedirectUrl:     qianchuanConfig.RedirectUrl,
			Concurrents:     uint32(qianchuanConfig.Concurrents),
			Status:          uint32(qianchuanConfig.Status),
			CreateTime:      tool.TimeToString("2006-01-02 15:04:05", qianchuanConfig.CreateTime),
			UpdateTime:      tool.TimeToString("2006-01-02 15:04:05", qianchuanConfig.UpdateTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(oceanengineConfigs.Total) / float64(oceanengineConfigs.PageSize)))

	return &v1.ListOceanengineConfigsReply{
		Code: 200,
		Data: &v1.ListOceanengineConfigsReply_Data{
			PageNum:   oceanengineConfigs.PageNum,
			PageSize:  oceanengineConfigs.PageSize,
			Total:     oceanengineConfigs.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListSelectOceanengineConfigs(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectOceanengineConfigsReply, error) {
	selects, err := ds.ocuc.ListSelectOceanengineConfigs(ctx)

	if err != nil {
		return nil, err
	}

	oceanengineType := make([]*v1.ListSelectOceanengineConfigsReply_OceanengineType, 0)

	for _, loceanengineType := range selects.OceanengineType {
		oceanengineType = append(oceanengineType, &v1.ListSelectOceanengineConfigsReply_OceanengineType{
			Key:   loceanengineType.Key,
			Value: loceanengineType.Value,
		})
	}

	return &v1.ListSelectOceanengineConfigsReply{
		Code: 200,
		Data: &v1.ListSelectOceanengineConfigsReply_Data{
			OceanengineType: oceanengineType,
		},
	}, nil
}

func (ds *DouyinService) RandOceanengineConfigs(ctx context.Context, in *v1.RandOceanengineConfigsRequest) (*v1.RandOceanengineConfigsReply, error) {
	oceanengineConfig, err := ds.ocuc.RandOceanengineConfigs(ctx, uint8(in.OceanengineType))

	if err != nil {
		return nil, err
	}

	return &v1.RandOceanengineConfigsReply{
		Code: 200,
		Data: &v1.RandOceanengineConfigsReply_Data{
			Id:              oceanengineConfig.Id,
			OceanengineType: uint32(oceanengineConfig.OceanengineType),
			AppId:           oceanengineConfig.AppId,
			AppName:         oceanengineConfig.AppName,
			AppSecret:       oceanengineConfig.AppSecret,
			RedirectUrl:     oceanengineConfig.RedirectUrl,
			Concurrents:     uint32(oceanengineConfig.Concurrents),
			Status:          uint32(oceanengineConfig.Status),
			CreateTime:      tool.TimeToString("2006-01-02 15:04:05", oceanengineConfig.CreateTime),
			UpdateTime:      tool.TimeToString("2006-01-02 15:04:05", oceanengineConfig.UpdateTime),
		},
	}, nil
}

func (ds *DouyinService) CreateOceanengineConfigs(ctx context.Context, in *v1.CreateOceanengineConfigsRequest) (*v1.CreateOceanengineConfigsReply, error) {
	oceanengineConfig, err := ds.ocuc.CreateOceanengineConfigs(ctx, in.AppId, in.AppName, in.AppSecret, in.RedirectUrl, uint8(in.OceanengineType), uint8(in.Concurrents), uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.CreateOceanengineConfigsReply{
		Code: 200,
		Data: &v1.CreateOceanengineConfigsReply_Data{
			Id:              oceanengineConfig.Id,
			OceanengineType: uint32(oceanengineConfig.OceanengineType),
			AppId:           oceanengineConfig.AppId,
			AppName:         oceanengineConfig.AppName,
			AppSecret:       oceanengineConfig.AppSecret,
			RedirectUrl:     oceanengineConfig.RedirectUrl,
			Concurrents:     uint32(oceanengineConfig.Concurrents),
			Status:          uint32(oceanengineConfig.Status),
			CreateTime:      tool.TimeToString("2006-01-02 15:04:05", oceanengineConfig.CreateTime),
			UpdateTime:      tool.TimeToString("2006-01-02 15:04:05", oceanengineConfig.UpdateTime),
		},
	}, nil
}

func (ds *DouyinService) UpdateOceanengineConfigs(ctx context.Context, in *v1.UpdateOceanengineConfigsRequest) (*v1.UpdateOceanengineConfigsReply, error) {
	oceanengineConfig, err := ds.ocuc.UpdateOceanengineConfigs(ctx, in.Id, in.AppId, in.AppName, in.AppSecret, in.RedirectUrl, uint8(in.OceanengineType), uint8(in.Concurrents), uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateOceanengineConfigsReply{
		Code: 200,
		Data: &v1.UpdateOceanengineConfigsReply_Data{
			Id:              oceanengineConfig.Id,
			OceanengineType: uint32(oceanengineConfig.OceanengineType),
			AppId:           oceanengineConfig.AppId,
			AppName:         oceanengineConfig.AppName,
			AppSecret:       oceanengineConfig.AppSecret,
			RedirectUrl:     oceanengineConfig.RedirectUrl,
			Concurrents:     uint32(oceanengineConfig.Concurrents),
			Status:          uint32(oceanengineConfig.Status),
			CreateTime:      tool.TimeToString("2006-01-02 15:04:05", oceanengineConfig.CreateTime),
			UpdateTime:      tool.TimeToString("2006-01-02 15:04:05", oceanengineConfig.UpdateTime),
		},
	}, nil
}

func (ds *DouyinService) UpdateStatusOceanengineConfigs(ctx context.Context, in *v1.UpdateStatusOceanengineConfigsRequest) (*v1.UpdateStatusOceanengineConfigsReply, error) {
	oceanengineConfig, err := ds.ocuc.UpdateStatusOceanengineConfigs(ctx, in.Id, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusOceanengineConfigsReply{
		Code: 200,
		Data: &v1.UpdateStatusOceanengineConfigsReply_Data{
			Id:              oceanengineConfig.Id,
			OceanengineType: uint32(oceanengineConfig.OceanengineType),
			AppId:           oceanengineConfig.AppId,
			AppName:         oceanengineConfig.AppName,
			AppSecret:       oceanengineConfig.AppSecret,
			RedirectUrl:     oceanengineConfig.RedirectUrl,
			Concurrents:     uint32(oceanengineConfig.Concurrents),
			Status:          uint32(oceanengineConfig.Status),
			CreateTime:      tool.TimeToString("2006-01-02 15:04:05", oceanengineConfig.CreateTime),
			UpdateTime:      tool.TimeToString("2006-01-02 15:04:05", oceanengineConfig.UpdateTime),
		},
	}, nil
}

func (ds *DouyinService) DeleteOceanengineConfigs(ctx context.Context, in *v1.DeleteOceanengineConfigsRequest) (*v1.DeleteOceanengineConfigsReply, error) {
	err := ds.ocuc.DeleteOceanengineConfigs(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteOceanengineConfigsReply{
		Code: 200,
		Data: &v1.DeleteOceanengineConfigsReply_Data{},
	}, nil
}
