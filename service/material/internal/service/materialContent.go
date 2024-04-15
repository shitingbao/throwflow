package service

import (
	"context"
	v1 "material/api/material/v1"
	"material/internal/pkg/tool"
	"math"
)

func (ms *MaterialService) ListMaterialContent(ctx context.Context, in *v1.ListMaterialContentRequest) (*v1.ListMaterialContentReply, error) {
	materials, err := ms.mcuc.ListMaterialContent(ctx, in.PageNum, in.PageSize, in.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMaterialContentReply_MaterialContent, 0)

	for _, material := range materials.List {
		list = append(list, &v1.ListMaterialContentReply_MaterialContent{
			Id:         material.Id,
			ProductId:  material.ProductId,
			UserId:     material.UserId,
			VideoId:    material.VideoId,
			Content:    material.Content,
			VideoName:  material.VideoName,
			VideoUrl:   material.VideoUrl,
			VideoCover: material.VideoCover,
			CreateTime: tool.TimeToString("2006-01-02 15:04:05", material.UpdateTime),
			UpdateTime: tool.TimeToString("2006-01-02 15:04:05", material.UpdateTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(materials.Total) / float64(materials.PageSize)))

	return &v1.ListMaterialContentReply{
		Code: 200,
		Data: &v1.ListMaterialContentReply_Data{
			PageNum:   materials.PageNum,
			PageSize:  materials.PageSize,
			Total:     materials.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ms *MaterialService) CreateMaterialContent(ctx context.Context, in *v1.CreateMaterialContentRequest) (*v1.CreateMaterialContentReply, error) {
	material, err := ms.mcuc.CreateMaterialContent(ctx, in.UserId, in.VideoId)

	if err != nil {
		return nil, err
	}

	return &v1.CreateMaterialContentReply{
		Code: 200,
		Data: &v1.CreateMaterialContentReply_Data{
			Id:         material.Id,
			ProductId:  material.ProductId,
			UserId:     material.UserId,
			VideoId:    material.VideoId,
			Content:    material.Content,
			VideoName:  material.VideoName,
			VideoUrl:   material.VideoUrl,
			VideoCover: material.VideoCover,
			CreateTime: tool.TimeToString("2006-01-02 15:04:05", material.UpdateTime),
			UpdateTime: tool.TimeToString("2006-01-02 15:04:05", material.UpdateTime),
		},
	}, nil
}

func (ms *MaterialService) UpdateMaterialContent(ctx context.Context, in *v1.UpdateMaterialContentRequest) (*v1.UpdateMaterialContentReply, error) {
	material, err := ms.mcuc.UpdateMaterialContent(ctx, in.UserId, in.VideoId, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateMaterialContentReply{
		Code: 200,
		Data: &v1.UpdateMaterialContentReply_Data{
			Id:         material.Id,
			ProductId:  material.ProductId,
			UserId:     material.UserId,
			VideoId:    material.VideoId,
			Content:    material.Content,
			VideoName:  material.VideoName,
			VideoUrl:   material.VideoUrl,
			VideoCover: material.VideoCover,
			CreateTime: tool.TimeToString("2006-01-02 15:04:05", material.UpdateTime),
			UpdateTime: tool.TimeToString("2006-01-02 15:04:05", material.UpdateTime),
		},
	}, nil
}

func (ms *MaterialService) RecoveMaterialContent(ctx context.Context, in *v1.RecoveMaterialContentRequest) (*v1.RecoveMaterialContentReply, error) {
	material, err := ms.mcuc.RecoveMaterialContent(ctx, in.UserId, in.VideoId)

	if err != nil {
		return nil, err
	}

	return &v1.RecoveMaterialContentReply{
		Code: 200,
		Data: &v1.RecoveMaterialContentReply_Data{
			Id:         material.Id,
			ProductId:  material.ProductId,
			UserId:     material.UserId,
			VideoId:    material.VideoId,
			Content:    material.Content,
			VideoName:  material.VideoName,
			VideoUrl:   material.VideoUrl,
			VideoCover: material.VideoCover,
			CreateTime: tool.TimeToString("2006-01-02 15:04:05", material.UpdateTime),
			UpdateTime: tool.TimeToString("2006-01-02 15:04:05", material.UpdateTime),
		},
	}, nil
}
