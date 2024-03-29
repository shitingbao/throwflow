package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 企业千川广告账户授权历史表
type QianchuanAdvertiserHistory struct {
	Id           uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	AdvertiserId uint64    `gorm:"column:advertiser_id;type:bigint(20) UNSIGNED;not null;comment:广告主ID"`
	Day          uint32    `gorm:"column:day;type:int(11);not null;comment:修改日期"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (QianchuanAdvertiserHistory) TableName() string {
	return "douyin_qianchuan_advertiser_history"
}

type qianchuanAdvertiserHistoryRepo struct {
	data *Data
	log  *log.Helper
}

func (qah *QianchuanAdvertiserHistory) ToDomain() *domain.QianchuanAdvertiserHistory {
	return &domain.QianchuanAdvertiserHistory{
		Id:           qah.Id,
		AdvertiserId: qah.AdvertiserId,
		Day:          qah.Day,
		CreateTime:   qah.CreateTime,
		UpdateTime:   qah.UpdateTime,
	}
}

func NewQianchuanAdvertiserHistoryRepo(data *Data, logger log.Logger) biz.QianchuanAdvertiserHistoryRepo {
	return &qianchuanAdvertiserHistoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qahr *qianchuanAdvertiserHistoryRepo) GetByAdvertiserId(ctx context.Context, advertiserId uint64) (*domain.QianchuanAdvertiserHistory, error) {
	qianchuanAdvertiserHistory := &QianchuanAdvertiserHistory{}

	if result := qahr.data.db.WithContext(ctx).Where("advertiser_id = ?", advertiserId).First(qianchuanAdvertiserHistory); result.Error != nil {
		return nil, result.Error
	}

	return qianchuanAdvertiserHistory.ToDomain(), nil
}

func (qahr *qianchuanAdvertiserHistoryRepo) List(ctx context.Context, advertiserIds string, day uint32) ([]*domain.QianchuanAdvertiserHistory, error) {
	list := make([]*domain.QianchuanAdvertiserHistory, 0)

	var qianchuanAdvertiserHistorys []QianchuanAdvertiserHistory

	if result := qahr.data.db.WithContext(ctx).
		Where("advertiser_id in ?", strings.Split(advertiserIds, ",")).
		Where("day = ?", day).
		Find(&qianchuanAdvertiserHistorys); result.Error != nil {
		return nil, result.Error
	}

	for _, qianchuanAdvertiserHistory := range qianchuanAdvertiserHistorys {
		list = append(list, qianchuanAdvertiserHistory.ToDomain())
	}

	return list, nil
}

func (qahr *qianchuanAdvertiserHistoryRepo) Save(ctx context.Context, in *domain.QianchuanAdvertiserHistory) (*domain.QianchuanAdvertiserHistory, error) {
	qianchuanAdvertiserHistory := &QianchuanAdvertiserHistory{
		AdvertiserId: in.AdvertiserId,
		Day:          in.Day,
		CreateTime:   in.CreateTime,
		UpdateTime:   in.UpdateTime,
	}

	if result := qahr.data.DB(ctx).Create(qianchuanAdvertiserHistory); result.Error != nil {
		return nil, result.Error
	}

	return qianchuanAdvertiserHistory.ToDomain(), nil
}
