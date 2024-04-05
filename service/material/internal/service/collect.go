package service

import (
	"context"
	v1 "material/api/material/v1"
	"material/internal/pkg/tool"
	"math"
	"time"
)

func (ms *MaterialService) ListCollectMaterials(ctx context.Context, in *v1.ListCollectMaterialsRequest) (*v1.ListCollectMaterialsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	materials, err := ms.cuc.ListCollectMaterials(ctx, in.PageNum, in.PageSize, in.CompanyId, in.Phone, in.Category, in.Keyword, in.Search, in.Msort, in.Mplatform)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCollectMaterialsReply_Materials, 0)

	for _, material := range materials.List {
		var videoCategory string

		if material.CategoryName == "" {
			videoCategory = material.IndustryName
		} else {
			videoCategory = material.IndustryName + "-" + material.CategoryName
		}

		list = append(list, &v1.ListCollectMaterialsReply_Materials{
			VideoId:            material.VideoId,
			VideoName:          material.VideoName,
			VideoUrl:           material.VideoUrl,
			VideoCover:         material.VideoCover,
			VideoCategory:      videoCategory,
			AwemeId:            material.AwemeId,
			AwemeName:          material.AwemeName,
			AwemeAccount:       material.AwemeAccount,
			AwemeFollowers:     material.AwemeFollowers,
			AwemeFollowersShow: material.AwemeFollowersShow,
			AwemeImg:           material.AwemeImg,
			AwemeLandingPage:   material.AwemeLandingPage,
			ProductId:          material.ProductId,
			ProductName:        material.ProductName,
			ProductImg:         material.ProductImg,
			ProductLandingPage: material.ProductLandingPage,
			ProductPrice:       material.ProductPrice,
			PlatformName:       material.PlatformName,
			UpdateDay:          tool.TimeToString("2006-01-02", material.UpdateDay),
			IsCollect:          uint32(material.IsCollect),
		})
	}

	totalPage := uint64(math.Ceil(float64(materials.Total) / float64(materials.PageSize)))

	return &v1.ListCollectMaterialsReply{
		Code: 200,
		Data: &v1.ListCollectMaterialsReply_Data{
			PageNum:   materials.PageNum,
			PageSize:  materials.PageSize,
			Total:     materials.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ms *MaterialService) UpdateCollects(ctx context.Context, in *v1.UpdateCollectsRequest) (*v1.UpdateCollectsReply, error) {
	if err := ms.cuc.UpdateCollects(ctx, in.CompanyId, in.VideoId, in.Phone); err != nil {
		return nil, err
	}

	return &v1.UpdateCollectsReply{
		Code: 200,
		Data: &v1.UpdateCollectsReply_Data{},
	}, nil
}
