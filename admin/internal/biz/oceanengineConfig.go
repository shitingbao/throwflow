package biz

import (
	v1 "admin/api/service/douyin/v1"
	"admin/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type OceanengineConfigRepo interface {
	List(context.Context, uint32, uint64) (*v1.ListOceanengineConfigsReply, error)
	ListSelect(context.Context) (*v1.ListSelectOceanengineConfigsReply, error)
	Save(context.Context, string, string, string, string, uint32, uint32, uint32) (*v1.CreateOceanengineConfigsReply, error)
	Update(context.Context, uint64, string, string, string, string, uint32, uint32, uint32) (*v1.UpdateOceanengineConfigsReply, error)
	UpdateStatus(context.Context, uint64, uint32) (*v1.UpdateStatusOceanengineConfigsReply, error)
	Delete(context.Context, uint64) (*v1.DeleteOceanengineConfigsReply, error)
}

type OceanengineConfigUsecase struct {
	repo OceanengineConfigRepo
	log  *log.Helper
}

func NewOceanengineConfigUsecase(repo OceanengineConfigRepo, logger log.Logger) *OceanengineConfigUsecase {
	return &OceanengineConfigUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (ocuc *OceanengineConfigUsecase) ListOceanengineConfigs(ctx context.Context, oceanengineType uint32, pageNum uint64) (*v1.ListOceanengineConfigsReply, error) {
	list, err := ocuc.repo.List(ctx, oceanengineType, pageNum)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (ocuc *OceanengineConfigUsecase) ListSelectOceanengineConfigs(ctx context.Context) (*v1.ListSelectOceanengineConfigsReply, error) {
	ListSelect, err := ocuc.repo.ListSelect(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return ListSelect, nil
}

func (ocuc *OceanengineConfigUsecase) CreateOceanengineConfigs(ctx context.Context, appId, appName, appSecret, redirectUrl string, oceanengineType, concurrents, status uint32) (*v1.CreateOceanengineConfigsReply, error) {
	oceanengineConfig, err := ocuc.repo.Save(ctx, appId, appName, appSecret, redirectUrl, oceanengineType, concurrents, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_OCEANENGINE_CONFIG_CREATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return oceanengineConfig, nil
}

func (ocuc *OceanengineConfigUsecase) UpdateOceanengineConfigs(ctx context.Context, id uint64, appId, appName, appSecret, redirectUrl string, oceanengineType, concurrents, status uint32) (*v1.UpdateOceanengineConfigsReply, error) {
	oceanengineConfig, err := ocuc.repo.Update(ctx, id, appId, appName, appSecret, redirectUrl, oceanengineType, concurrents, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_OCEANENGINE_CONFIG_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return oceanengineConfig, nil
}

func (ocuc *OceanengineConfigUsecase) UpdateStatusOceanengineConfigs(ctx context.Context, id uint64, status uint32) (*v1.UpdateStatusOceanengineConfigsReply, error) {
	oceanengineConfig, err := ocuc.repo.UpdateStatus(ctx, id, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_OCEANENGINE_CONFIG_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return oceanengineConfig, nil
}

func (ocuc *OceanengineConfigUsecase) DeleteOceanengineConfigs(ctx context.Context, id uint64) error {
	_, err := ocuc.repo.Delete(ctx, id)

	if err != nil {
		return errors.InternalServer("ADMIN_OCEANENGINE_CONFIG_DELETE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return nil
}
