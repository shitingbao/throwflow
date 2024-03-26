package data

import (
	"common/internal/biz"
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 更新日志表
type UpdateLog struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	Name       string    `gorm:"column:name;type:char(255);not null;index:phone;comment:标题"`
	Content    string    `gorm:"column:content;type:text;not null;comment:更新内容"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UpdateLog) TableName() string {
	return "common_update_log"
}

type updateLogRepo struct {
	data *Data
	log  *log.Helper
}

func (ul *UpdateLog) ToDomain() *domain.UpdateLog {
	return &domain.UpdateLog{
		Id:         ul.Id,
		Name:       ul.Name,
		Content:    ul.Content,
		CreateTime: ul.CreateTime,
		UpdateTime: ul.UpdateTime,
	}
}

func NewUpdateLogRepo(data *Data, logger log.Logger) biz.UpdateLogRepo {
	return &updateLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ulr *updateLogRepo) Get(ctx context.Context, id uint64) (*domain.UpdateLog, error) {
	updateLog := &UpdateLog{}

	if result := ulr.data.db.WithContext(ctx).First(updateLog, id); result.Error != nil {
		return nil, result.Error
	}

	return updateLog.ToDomain(), nil
}

func (ulr *updateLogRepo) List(ctx context.Context) ([]*domain.UpdateLog, error) {
	var updateLogs []UpdateLog
	list := make([]*domain.UpdateLog, 0)

	if result := ulr.data.db.WithContext(ctx).
		Order("create_time DESC").
		Find(&updateLogs); result.Error != nil {
		return nil, result.Error
	}

	for _, updateLog := range updateLogs {
		list = append(list, updateLog.ToDomain())
	}

	return list, nil
}

func (ulr *updateLogRepo) Save(ctx context.Context, in *domain.UpdateLog) (*domain.UpdateLog, error) {
	updateLog := &UpdateLog{
		Name:       in.Name,
		Content:    in.Content,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := ulr.data.db.WithContext(ctx).Create(updateLog); result.Error != nil {
		return nil, result.Error
	}

	return updateLog.ToDomain(), nil
}

func (ulr *updateLogRepo) Update(ctx context.Context, in *domain.UpdateLog) (*domain.UpdateLog, error) {
	updateLog := &UpdateLog{
		Id:         in.Id,
		Name:       in.Name,
		Content:    in.Content,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := ulr.data.db.WithContext(ctx).Save(updateLog); result.Error != nil {
		return nil, result.Error
	}

	return updateLog.ToDomain(), nil
}
