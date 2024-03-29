package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 抖音开放平台达人信息新增日志表
type OpenDouyinUserInfoCreateLog struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ClientKey  string    `gorm:"column:client_key;type:varchar(50);not null;uniqueIndex:client_key_open_id;comment:应用Client Key"`
	OpenId     string    `gorm:"column:open_id;type:varchar(100);not null;uniqueIndex:client_key_open_id;comment:授权用户唯一标识"`
	IsHandle   uint8     `gorm:"column:is_handle;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否处理，1：正在处理，0：未处理"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (OpenDouyinUserInfoCreateLog) TableName() string {
	return "douyin_open_douyin_user_info_create_log"
}

type openDouyinUserInfoCreateLogRepo struct {
	data *Data
	log  *log.Helper
}

func (oduicl *OpenDouyinUserInfoCreateLog) ToDomain() *domain.OpenDouyinUserInfoCreateLog {
	return &domain.OpenDouyinUserInfoCreateLog{
		Id:         oduicl.Id,
		ClientKey:  oduicl.ClientKey,
		OpenId:     oduicl.OpenId,
		IsHandle:   oduicl.IsHandle,
		CreateTime: oduicl.CreateTime,
		UpdateTime: oduicl.UpdateTime,
	}
}

func NewOpenDouyinUserInfoCreateLogRepo(data *Data, logger log.Logger) biz.OpenDouyinUserInfoCreateLogRepo {
	return &openDouyinUserInfoCreateLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (oduiclr *openDouyinUserInfoCreateLogRepo) List(ctx context.Context, isHandle string) ([]*domain.OpenDouyinUserInfoCreateLog, error) {
	var openDouyinUserInfoCreateLogs []OpenDouyinUserInfoCreateLog
	list := make([]*domain.OpenDouyinUserInfoCreateLog, 0)

	if result := oduiclr.data.db.WithContext(ctx).Where("is_handle = ?", isHandle).
		Find(&openDouyinUserInfoCreateLogs); result.Error != nil {
		return nil, result.Error
	}

	for _, openDouyinUserInfoCreateLog := range openDouyinUserInfoCreateLogs {
		list = append(list, openDouyinUserInfoCreateLog.ToDomain())
	}

	return list, nil
}

func (oduiclr *openDouyinUserInfoCreateLogRepo) Save(ctx context.Context, in *domain.OpenDouyinUserInfoCreateLog) (*domain.OpenDouyinUserInfoCreateLog, error) {
	openDouyinUserInfoCreateLog := &OpenDouyinUserInfoCreateLog{
		ClientKey:  in.ClientKey,
		OpenId:     in.OpenId,
		IsHandle:   in.IsHandle,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := oduiclr.data.DB(ctx).Create(openDouyinUserInfoCreateLog); result.Error != nil {
		return nil, result.Error
	}

	return openDouyinUserInfoCreateLog.ToDomain(), nil
}

func (oduiclr *openDouyinUserInfoCreateLogRepo) Update(ctx context.Context, in *domain.OpenDouyinUserInfoCreateLog) (*domain.OpenDouyinUserInfoCreateLog, error) {
	openDouyinUserInfoCreateLog := &OpenDouyinUserInfoCreateLog{
		Id:         in.Id,
		ClientKey:  in.ClientKey,
		OpenId:     in.OpenId,
		IsHandle:   in.IsHandle,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := oduiclr.data.db.WithContext(ctx).Save(openDouyinUserInfoCreateLog); result.Error != nil {
		return nil, result.Error
	}

	return openDouyinUserInfoCreateLog.ToDomain(), nil
}

func (oduiclr *openDouyinUserInfoCreateLogRepo) Delete(ctx context.Context, in *domain.OpenDouyinUserInfoCreateLog) error {
	openDouyinUserInfoCreateLog := &OpenDouyinUserInfoCreateLog{
		Id:         in.Id,
		ClientKey:  in.ClientKey,
		OpenId:     in.OpenId,
		IsHandle:   in.IsHandle,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := oduiclr.data.DB(ctx).Delete(openDouyinUserInfoCreateLog); result.Error != nil {
		return result.Error
	}

	return nil
}
