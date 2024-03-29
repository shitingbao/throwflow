package data

import (
	"context"
	"douyin/internal/domain"
	"time"

	"douyin/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// 巨量引擎应用配置表
type OceanengineConfig struct {
	Id              uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	OceanengineType uint8     `gorm:"column:oceanengine_type;type:tinyint(3) UNSIGNED;not null;default:1;comment:巨量引擎应用，1：巨量千川，2：巨量广告"`
	AppId           string    `gorm:"column:app_id;uniqueIndex;type:varchar(20);not null;comment:应用ID"`
	AppName         string    `gorm:"column:app_name;type:varchar(200);not null;comment:应用名称"`
	AppSecret       string    `gorm:"column:app_secret;type:varchar(255);not null;comment:应用密钥"`
	RedirectUrl     string    `gorm:"column:redirect_url;type:text;not null;comment:授权跳转链接"`
	Concurrents     uint8     `gorm:"column:concurrents;type:tinyint(3) UNSIGNED;not null;default:0;comment:并发数"`
	Status          uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态：1：启用，0：禁用"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (OceanengineConfig) TableName() string {
	return "douyin_oceanengine_config"
}

type oceanengineConfigRepo struct {
	data *Data
	log  *log.Helper
}

func (oc *OceanengineConfig) ToDomain() *domain.OceanengineConfig {
	return &domain.OceanengineConfig{
		Id:              oc.Id,
		OceanengineType: oc.OceanengineType,
		AppId:           oc.AppId,
		AppName:         oc.AppName,
		AppSecret:       oc.AppSecret,
		RedirectUrl:     oc.RedirectUrl,
		Concurrents:     oc.Concurrents,
		Status:          oc.Status,
		CreateTime:      oc.CreateTime,
		UpdateTime:      oc.UpdateTime,
	}
}

func NewOceanengineConfigRepo(data *Data, logger log.Logger) biz.OceanengineConfigRepo {
	return &oceanengineConfigRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ocr *oceanengineConfigRepo) GetById(ctx context.Context, id uint64) (*domain.OceanengineConfig, error) {
	oceanengineConfig := &OceanengineConfig{}

	if result := ocr.data.db.WithContext(ctx).First(oceanengineConfig, id); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineConfig.ToDomain(), nil
}

func (ocr *oceanengineConfigRepo) Rand(ctx context.Context, oceanengineType uint8) (*domain.OceanengineConfig, error) {
	oceanengineConfig := &OceanengineConfig{}

	if result := ocr.data.db.WithContext(ctx).Where("oceanengine_type = ?", oceanengineType).Order("RAND()").First(oceanengineConfig); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineConfig.ToDomain(), nil
}

func (ocr *oceanengineConfigRepo) GetByAppId(ctx context.Context, appId string) (*domain.OceanengineConfig, error) {
	oceanengineConfig := &OceanengineConfig{}

	if result := ocr.data.db.WithContext(ctx).Where("app_id = ?", appId).First(oceanengineConfig); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineConfig.ToDomain(), nil
}

func (ocr *oceanengineConfigRepo) List(ctx context.Context, oceanengineType uint8, pageNum, pageSize int) ([]*domain.OceanengineConfig, error) {
	var oceanengineConfigs []OceanengineConfig
	list := make([]*domain.OceanengineConfig, 0)

	db := ocr.data.db.WithContext(ctx)

	if oceanengineType > 0 {
		db = db.Where("oceanengine_type = ?", oceanengineType)
	}

	if pageNum == 0 {
		if result := db.
			Order("update_time DESC,id DESC").
			Find(&oceanengineConfigs); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.
			Order("update_time DESC,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&oceanengineConfigs); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, oceanengineConfig := range oceanengineConfigs {
		list = append(list, oceanengineConfig.ToDomain())
	}

	return list, nil
}

func (ocr *oceanengineConfigRepo) Count(ctx context.Context, oceanengineType uint8) (int64, error) {
	var count int64

	db := ocr.data.db.WithContext(ctx)

	if oceanengineType > 0 {
		db = db.Where("oceanengine_type = ?", oceanengineType)
	}

	if result := db.Model(&OceanengineConfig{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ocr *oceanengineConfigRepo) Save(ctx context.Context, in *domain.OceanengineConfig) (*domain.OceanengineConfig, error) {
	oceanengineConfig := &OceanengineConfig{
		OceanengineType: in.OceanengineType,
		AppId:           in.AppId,
		AppName:         in.AppName,
		AppSecret:       in.AppSecret,
		RedirectUrl:     in.RedirectUrl,
		Concurrents:     in.Concurrents,
		Status:          in.Status,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := ocr.data.db.WithContext(ctx).Create(oceanengineConfig); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineConfig.ToDomain(), nil
}

func (ocr *oceanengineConfigRepo) Update(ctx context.Context, in *domain.OceanengineConfig) (*domain.OceanengineConfig, error) {
	oceanengineConfig := &OceanengineConfig{
		Id:              in.Id,
		OceanengineType: in.OceanengineType,
		AppId:           in.AppId,
		AppName:         in.AppName,
		AppSecret:       in.AppSecret,
		RedirectUrl:     in.RedirectUrl,
		Concurrents:     in.Concurrents,
		Status:          in.Status,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := ocr.data.db.WithContext(ctx).Save(oceanengineConfig); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineConfig.ToDomain(), nil
}

func (ocr *oceanengineConfigRepo) Delete(ctx context.Context, in *domain.OceanengineConfig) error {
	oceanengineConfig := &OceanengineConfig{
		Id:              in.Id,
		OceanengineType: in.OceanengineType,
		AppId:           in.AppId,
		AppName:         in.AppName,
		AppSecret:       in.AppSecret,
		RedirectUrl:     in.RedirectUrl,
		Concurrents:     in.Concurrents,
		Status:          in.Status,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := ocr.data.db.WithContext(ctx).Delete(oceanengineConfig); result.Error != nil {
		return result.Error
	}

	return nil
}
