package data

import (
	"common/internal/biz"
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 支付异步通知日志表
type PayAsyncNotificationLog struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	Content    string    `gorm:"column:content;type:text;not null;comment:支付异步通知内容"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (PayAsyncNotificationLog) TableName() string {
	return "common_pay_async_notification_log"
}

type payAsyncNotificationLogRepo struct {
	data *Data
	log  *log.Helper
}

func (panl *PayAsyncNotificationLog) ToDomain() *domain.PayAsyncNotificationLog {
	return &domain.PayAsyncNotificationLog{
		Id:         panl.Id,
		Content:    panl.Content,
		CreateTime: panl.CreateTime,
		UpdateTime: panl.UpdateTime,
	}
}

func NewPayAsyncNotificationLogRepo(data *Data, logger log.Logger) biz.PayAsyncNotificationLogRepo {
	return &payAsyncNotificationLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (panlr *payAsyncNotificationLogRepo) Save(ctx context.Context, in *domain.PayAsyncNotificationLog) (*domain.PayAsyncNotificationLog, error) {
	payAsyncNotificationLog := &PayAsyncNotificationLog{
		Content:    in.Content,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := panlr.data.db.WithContext(ctx).Create(payAsyncNotificationLog); result.Error != nil {
		return nil, result.Error
	}

	return payAsyncNotificationLog.ToDomain(), nil
}
