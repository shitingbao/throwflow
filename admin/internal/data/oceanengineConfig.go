package data

import (
	v1 "admin/api/service/douyin/v1"
	"admin/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type oceanengineConfigRepo struct {
	data *Data
	log  *log.Helper
}

func NewOceanengineConfigRepo(data *Data, logger log.Logger) biz.OceanengineConfigRepo {
	return &oceanengineConfigRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ocr *oceanengineConfigRepo) List(ctx context.Context, oceanengineType uint32, pageNum uint64) (*v1.ListOceanengineConfigsReply, error) {
	list, err := ocr.data.douyinuc.ListOceanengineConfigs(ctx, &v1.ListOceanengineConfigsRequest{
		OceanengineType: oceanengineType,
		PageNum:         pageNum,
		PageSize:        uint64(ocr.data.conf.Database.PageSize),
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (ocr *oceanengineConfigRepo) ListSelect(ctx context.Context) (*v1.ListSelectOceanengineConfigsReply, error) {
	listSelect, err := ocr.data.douyinuc.ListSelectOceanengineConfigs(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return listSelect, err
}

func (ocr *oceanengineConfigRepo) Save(ctx context.Context, appId, appName, appSecret, redirectUrl string, oceanengineType, concurrents, status uint32) (*v1.CreateOceanengineConfigsReply, error) {
	oceanengineConfig, err := ocr.data.douyinuc.CreateOceanengineConfigs(ctx, &v1.CreateOceanengineConfigsRequest{
		OceanengineType: oceanengineType,
		AppId:           appId,
		AppName:         appName,
		AppSecret:       appSecret,
		RedirectUrl:     redirectUrl,
		Concurrents:     concurrents,
		Status:          status,
	})

	if err != nil {
		return nil, err
	}

	return oceanengineConfig, err
}

func (ocr *oceanengineConfigRepo) Update(ctx context.Context, id uint64, appId, appName, appSecret, redirectUrl string, oceanengineType, concurrents, status uint32) (*v1.UpdateOceanengineConfigsReply, error) {
	oceanengineConfig, err := ocr.data.douyinuc.UpdateOceanengineConfigs(ctx, &v1.UpdateOceanengineConfigsRequest{
		Id:              id,
		OceanengineType: oceanengineType,
		AppId:           appId,
		AppName:         appName,
		AppSecret:       appSecret,
		RedirectUrl:     redirectUrl,
		Concurrents:     concurrents,
		Status:          status,
	})

	if err != nil {
		return nil, err
	}

	return oceanengineConfig, err
}

func (ocr *oceanengineConfigRepo) UpdateStatus(ctx context.Context, id uint64, status uint32) (*v1.UpdateStatusOceanengineConfigsReply, error) {
	oceanengineConfig, err := ocr.data.douyinuc.UpdateStatusOceanengineConfigs(ctx, &v1.UpdateStatusOceanengineConfigsRequest{
		Id:     id,
		Status: status,
	})

	if err != nil {
		return nil, err
	}

	return oceanengineConfig, err
}

func (ocr *oceanengineConfigRepo) Delete(ctx context.Context, id uint64) (*v1.DeleteOceanengineConfigsReply, error) {
	oceanengineConfig, err := ocr.data.douyinuc.DeleteOceanengineConfigs(ctx, &v1.DeleteOceanengineConfigsRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return oceanengineConfig, err
}
