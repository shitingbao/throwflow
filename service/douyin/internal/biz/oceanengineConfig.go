package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	DouyinOceanengineConfigCreateError = errors.InternalServer("DOUYIN_OCEANENGINE_CONFIG_CREATE_ERROR", "巨量千川应用配置创建失败")
	DouyinOceanengineConfigUpdateError = errors.InternalServer("DOUYIN_OCEANENGINE_CONFIG_UPDATE_ERROR", "巨量千川应用配置更新失败")
	DouyinOceanengineConfigDeleteError = errors.InternalServer("DOUYIN_OCEANENGINE_CONFIG_DELETE_ERROR", "巨量千川应用配置删除失败")
	DouyinOceanengineConfigNotFound    = errors.NotFound("DOUYIN_OCEANENGINE_CONFIG_NOT_FOUND", "巨量千川应用配置不存在")
)

type OceanengineConfigRepo interface {
	GetById(context.Context, uint64) (*domain.OceanengineConfig, error)
	Rand(context.Context, uint8) (*domain.OceanengineConfig, error)
	GetByAppId(context.Context, string) (*domain.OceanengineConfig, error)
	List(context.Context, uint8, int, int) ([]*domain.OceanengineConfig, error)
	Count(context.Context, uint8) (int64, error)
	Save(context.Context, *domain.OceanengineConfig) (*domain.OceanengineConfig, error)
	Update(context.Context, *domain.OceanengineConfig) (*domain.OceanengineConfig, error)
	Delete(context.Context, *domain.OceanengineConfig) error
}

type OceanengineConfigUsecase struct {
	repo OceanengineConfigRepo
	conf *conf.Data
	log  *log.Helper
}

func NewOceanengineConfigUsecase(repo OceanengineConfigRepo, conf *conf.Data, logger log.Logger) *OceanengineConfigUsecase {
	return &OceanengineConfigUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (ocuc *OceanengineConfigUsecase) ListOceanengineConfigs(ctx context.Context, oceanengineType uint8, pageNum, pageSize uint64) (*domain.OceanengineConfigList, error) {
	list, err := ocuc.repo.List(ctx, oceanengineType, int(pageNum), int(pageSize))

	if err != nil {
		return nil, DouyinDataError
	}

	total, err := ocuc.repo.Count(ctx, oceanengineType)

	if err != nil {
		return nil, DouyinDataError
	}

	return &domain.OceanengineConfigList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (ocuc *OceanengineConfigUsecase) ListSelectOceanengineConfigs(ctx context.Context) (*domain.SelectOceanengineConfigs, error) {
	return domain.NewSelectOceanengineConfigs(), nil
}

/*func (o *OceanengineConfigUsecase) GetQianchuanConfigs(ctx context.Context, appId string) (*domain.OceanengineConfig, error) {
	oceanengineConfig, err := o.repo.GetByAppId(ctx, appId)

	if err != nil {
		return nil, DouyinDataError
	}

	return oceanengineConfig, nil
}*/

func (ocuc *OceanengineConfigUsecase) RandOceanengineConfigs(ctx context.Context, oceanengineType uint8) (*domain.OceanengineConfig, error) {
	oceanengineConfig, err := ocuc.repo.Rand(ctx, oceanengineType)

	if err != nil {
		return nil, DouyinDataError
	}

	return oceanengineConfig, nil
}

func (ocuc *OceanengineConfigUsecase) CreateOceanengineConfigs(ctx context.Context, appId, appName, appSecret, redirectUrl string, oceanengineType, concurrents, status uint8) (*domain.OceanengineConfig, error) {
	inOceanengineConfig := domain.NewOceanengineConfig(ctx, appId, appName, appSecret, redirectUrl, oceanengineType, concurrents, status)
	inOceanengineConfig.SetCreateTime(ctx)
	inOceanengineConfig.SetUpdateTime(ctx)

	oceanengineConfig, err := ocuc.repo.Save(ctx, inOceanengineConfig)

	if err != nil {
		return nil, DouyinOceanengineConfigCreateError
	}

	return oceanengineConfig, nil
}

func (ocuc *OceanengineConfigUsecase) UpdateOceanengineConfigs(ctx context.Context, id uint64, appId, appName, appSecret, redirectUrl string, oceanengineType, concurrents, status uint8) (*domain.OceanengineConfig, error) {
	inOceanengineConfig, err := ocuc.getOceanengineConfigById(ctx, id)

	if err != nil {
		return nil, DouyinOceanengineConfigNotFound
	}

	inOceanengineConfig.SetOceanengineType(ctx, oceanengineType)
	inOceanengineConfig.SetAppId(ctx, appId)
	inOceanengineConfig.SetAppName(ctx, appName)
	inOceanengineConfig.SetAppSecret(ctx, appSecret)
	inOceanengineConfig.SetRedirectUrl(ctx, redirectUrl)
	inOceanengineConfig.SetConcurrents(ctx, concurrents)
	inOceanengineConfig.SetStatus(ctx, status)
	inOceanengineConfig.SetUpdateTime(ctx)

	oceanengineConfig, err := ocuc.repo.Update(ctx, inOceanengineConfig)

	if err != nil {
		return nil, DouyinOceanengineConfigUpdateError
	}

	return oceanengineConfig, nil
}

func (ocuc *OceanengineConfigUsecase) UpdateStatusOceanengineConfigs(ctx context.Context, id uint64, status uint8) (*domain.OceanengineConfig, error) {
	inOceanengineConfig, err := ocuc.getOceanengineConfigById(ctx, id)

	if err != nil {
		return nil, DouyinOceanengineConfigNotFound
	}

	inOceanengineConfig.SetStatus(ctx, status)
	inOceanengineConfig.SetUpdateTime(ctx)

	oceanengineConfig, err := ocuc.repo.Update(ctx, inOceanengineConfig)

	if err != nil {
		return nil, DouyinOceanengineConfigUpdateError
	}

	return oceanengineConfig, nil
}

func (ocuc *OceanengineConfigUsecase) DeleteOceanengineConfigs(ctx context.Context, id uint64) error {
	inOceanengineConfig, err := ocuc.getOceanengineConfigById(ctx, id)

	if err != nil {
		return DouyinOceanengineConfigNotFound
	}

	if err := ocuc.repo.Delete(ctx, inOceanengineConfig); err != nil {
		return DouyinOceanengineConfigDeleteError
	}

	return nil
}

func (ocuc *OceanengineConfigUsecase) getOceanengineConfigById(ctx context.Context, id uint64) (*domain.OceanengineConfig, error) {
	oceanengineConfig, err := ocuc.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return oceanengineConfig, nil
}
