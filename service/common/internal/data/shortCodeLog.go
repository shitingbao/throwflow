package data

import (
	"common/internal/biz"
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 短码生成日志表
type ShortCodeLog struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ShortCode  string    `gorm:"column:short_code;type:varchar(10);uniqueIndex:short_code;not null;comment:短码"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (ShortCodeLog) TableName() string {
	return "common_short_code_log"
}

type shortCodeLogRepo struct {
	data *Data
	log  *log.Helper
}

func (scl *ShortCodeLog) ToDomain() *domain.ShortCodeLog {
	return &domain.ShortCodeLog{
		Id:         scl.Id,
		ShortCode:  scl.ShortCode,
		CreateTime: scl.CreateTime,
		UpdateTime: scl.UpdateTime,
	}
}

func NewShortCodeLogRepo(data *Data, logger log.Logger) biz.ShortCodeLogRepo {
	return &shortCodeLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (sclr *shortCodeLogRepo) Get(ctx context.Context, shortCode string) (*domain.ShortCodeLog, error) {
	shortCodeLog := &ShortCodeLog{}

	if result := sclr.data.db.WithContext(ctx).
		Where("short_code = ?", shortCode).
		First(shortCodeLog); result.Error != nil {
		return nil, result.Error
	}

	return shortCodeLog.ToDomain(), nil
}

func (sclr *shortCodeLogRepo) Save(ctx context.Context, in *domain.ShortCodeLog) (*domain.ShortCodeLog, error) {
	shortCodeLog := &ShortCodeLog{
		ShortCode:  in.ShortCode,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := sclr.data.db.WithContext(ctx).Create(shortCodeLog); result.Error != nil {
		return nil, result.Error
	}

	return shortCodeLog.ToDomain(), nil
}
