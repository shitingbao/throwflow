package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
)

func (as *AdminService) ListOceanengineConfigs(ctx context.Context, in *v1.ListOceanengineConfigsRequest) (*v1.ListOceanengineConfigsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:oceanengineConfig:list"); err != nil {
		return nil, err
	}

	oceanengineConfigs, err := as.ocuc.ListOceanengineConfigs(ctx, in.OceanengineType, in.PageNum)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListOceanengineConfigsReply_OceanengineConfigs, 0)

	for _, oceanengineConfig := range oceanengineConfigs.Data.List {
		list = append(list, &v1.ListOceanengineConfigsReply_OceanengineConfigs{
			Id:              oceanengineConfig.Id,
			OceanengineType: oceanengineConfig.OceanengineType,
			AppId:           oceanengineConfig.AppId,
			AppName:         oceanengineConfig.AppName,
			AppSecret:       oceanengineConfig.AppSecret,
			RedirectUrl:     oceanengineConfig.RedirectUrl,
			Concurrents:     oceanengineConfig.Concurrents,
			Status:          oceanengineConfig.Status,
			CreateTime:      oceanengineConfig.CreateTime,
			UpdateTime:      oceanengineConfig.UpdateTime,
		})
	}

	totalPage := uint64(math.Ceil(float64(oceanengineConfigs.Data.Total) / float64(oceanengineConfigs.Data.PageSize)))

	return &v1.ListOceanengineConfigsReply{
		Code: 200,
		Data: &v1.ListOceanengineConfigsReply_Data{
			PageNum:   oceanengineConfigs.Data.PageNum,
			PageSize:  oceanengineConfigs.Data.PageSize,
			Total:     oceanengineConfigs.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (as *AdminService) ListSelectOceanengineConfigs(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectOceanengineConfigsReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	selects, err := as.ocuc.ListSelectOceanengineConfigs(ctx)

	if err != nil {
		return nil, err
	}

	oceanengineType := make([]*v1.ListSelectOceanengineConfigsReply_OceanengineType, 0)

	for _, loceanengineType := range selects.Data.OceanengineType {
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

func (as *AdminService) CreateOceanengineConfigs(ctx context.Context, in *v1.CreateOceanengineConfigsRequest) (*v1.CreateOceanengineConfigsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:oceanengineConfig:create"); err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	oceanengineConfig, err := as.ocuc.CreateOceanengineConfigs(ctx, in.AppId, in.AppName, in.AppSecret, in.RedirectUrl, in.OceanengineType, in.Concurrents, in.Status)

	if err != nil {
		return nil, err
	}

	return &v1.CreateOceanengineConfigsReply{
		Code: 200,
		Data: &v1.CreateOceanengineConfigsReply_Data{
			Id:              oceanengineConfig.Data.Id,
			OceanengineType: oceanengineConfig.Data.OceanengineType,
			AppId:           oceanengineConfig.Data.AppId,
			AppName:         oceanengineConfig.Data.AppName,
			AppSecret:       oceanengineConfig.Data.AppSecret,
			RedirectUrl:     oceanengineConfig.Data.RedirectUrl,
			Concurrents:     oceanengineConfig.Data.Concurrents,
			Status:          oceanengineConfig.Data.Status,
			CreateTime:      oceanengineConfig.Data.CreateTime,
			UpdateTime:      oceanengineConfig.Data.UpdateTime,
		},
	}, nil
}

func (as *AdminService) UpdateOceanengineConfigs(ctx context.Context, in *v1.UpdateOceanengineConfigsRequest) (*v1.UpdateOceanengineConfigsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:oceanengineConfig:update"); err != nil {
		return nil, err
	}

	oceanengineConfig, err := as.ocuc.UpdateOceanengineConfigs(ctx, in.Id, in.AppId, in.AppName, in.AppSecret, in.RedirectUrl, in.OceanengineType, in.Concurrents, in.Status)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateOceanengineConfigsReply{
		Code: 200,
		Data: &v1.UpdateOceanengineConfigsReply_Data{
			Id:          oceanengineConfig.Data.Id,
			AppId:       oceanengineConfig.Data.AppId,
			AppName:     oceanengineConfig.Data.AppName,
			AppSecret:   oceanengineConfig.Data.AppSecret,
			RedirectUrl: oceanengineConfig.Data.RedirectUrl,
			Concurrents: oceanengineConfig.Data.Concurrents,
			Status:      oceanengineConfig.Data.Status,
			CreateTime:  oceanengineConfig.Data.CreateTime,
			UpdateTime:  oceanengineConfig.Data.UpdateTime,
		},
	}, nil
}

func (as *AdminService) UpdateStatusOceanengineConfigs(ctx context.Context, in *v1.UpdateStatusOceanengineConfigsRequest) (*v1.UpdateStatusOceanengineConfigsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:oceanengineConfig:updateStatus"); err != nil {
		return nil, err
	}

	oceanengineConfig, err := as.ocuc.UpdateStatusOceanengineConfigs(ctx, in.Id, in.Status)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusOceanengineConfigsReply{
		Code: 200,
		Data: &v1.UpdateStatusOceanengineConfigsReply_Data{
			Id:          oceanengineConfig.Data.Id,
			AppId:       oceanengineConfig.Data.AppId,
			AppName:     oceanengineConfig.Data.AppName,
			AppSecret:   oceanengineConfig.Data.AppSecret,
			RedirectUrl: oceanengineConfig.Data.RedirectUrl,
			Concurrents: oceanengineConfig.Data.Concurrents,
			Status:      oceanengineConfig.Data.Status,
			CreateTime:  oceanengineConfig.Data.CreateTime,
			UpdateTime:  oceanengineConfig.Data.UpdateTime,
		},
	}, nil
}

func (as *AdminService) DeleteOceanengineConfigs(ctx context.Context, in *v1.DeleteOceanengineConfigsRequest) (*v1.DeleteOceanengineConfigsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:oceanengineConfig:delete"); err != nil {
		return nil, err
	}

	err := as.ocuc.DeleteOceanengineConfigs(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteOceanengineConfigsReply{
		Code: 200,
		Data: &v1.DeleteOceanengineConfigsReply_Data{},
	}, nil
}
