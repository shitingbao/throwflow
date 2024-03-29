package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业千川广告账户状态表
type QianchuanAdvertiserStatus struct {
	CompanyId    uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:公司ID"`
	AdvertiserId uint64    `gorm:"column:advertiser_id;type:bigint(20) UNSIGNED;not null;comment:广告主ID"`
	Status       uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态: 1 启用 0 禁用"`
	Day          uint32    `gorm:"column:day;type:int(11);not null;comment:修改日期"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (QianchuanAdvertiserStatus) TableName() string {
	return "douyin_qianchuan_advertiser_status"
}

type qianchuanAdvertiserStatusRepo struct {
	data *Data
	log  *log.Helper
}

func (qas *QianchuanAdvertiserStatus) ToDomain() *domain.QianchuanAdvertiserStatus {
	return &domain.QianchuanAdvertiserStatus{
		CompanyId:    qas.CompanyId,
		AdvertiserId: qas.AdvertiserId,
		Status:       qas.Status,
		Day:          qas.Day,
		CreateTime:   qas.CreateTime,
		UpdateTime:   qas.UpdateTime,
	}
}

func NewQianchuanAdvertiserStatusRepo(data *Data, logger log.Logger) biz.QianchuanAdvertiserStatusRepo {
	return &qianchuanAdvertiserStatusRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qasr *qianchuanAdvertiserStatusRepo) GetByCompanyIdAndAdvertiserIdAndDay(ctx context.Context, companyId, advertiserId uint64, day uint32) (*domain.QianchuanAdvertiserStatus, error) {
	var qianchuanAdvertiserStatuses []QianchuanAdvertiserStatus
	var qianchuanAdvertiserStatus *domain.QianchuanAdvertiserStatus

	if result := qasr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("advertiser_id = ?", advertiserId).
		Where("day <= ?", day).
		Order("day ASC").
		Find(&qianchuanAdvertiserStatuses); result.Error != nil {
		return nil, result.Error
	}

	for _, lqianchuanAdvertiserStatus := range qianchuanAdvertiserStatuses {
		if qianchuanAdvertiserStatus == nil {
			qianchuanAdvertiserStatus = lqianchuanAdvertiserStatus.ToDomain()
		}

		qianchuanAdvertiserStatus.Status = lqianchuanAdvertiserStatus.Status
		qianchuanAdvertiserStatus.Day = lqianchuanAdvertiserStatus.Day
		qianchuanAdvertiserStatus.CreateTime = lqianchuanAdvertiserStatus.CreateTime
		qianchuanAdvertiserStatus.UpdateTime = lqianchuanAdvertiserStatus.UpdateTime
	}

	return qianchuanAdvertiserStatus, nil
}

func (qasr *qianchuanAdvertiserStatusRepo) List(ctx context.Context, companyId uint64, day uint32) ([]*domain.QianchuanAdvertiserStatus, error) {
	list := make([]*domain.QianchuanAdvertiserStatus, 0)

	var qianchuanAdvertiserStatuses []QianchuanAdvertiserStatus

	if result := qasr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("day <= ?", day).
		Order("day ASC").
		Find(&qianchuanAdvertiserStatuses); result.Error != nil {
		return nil, result.Error
	}

	for _, qianchuanAdvertiserStatus := range qianchuanAdvertiserStatuses {
		isNotExist := true

		for _, l := range list {
			if l.AdvertiserId == qianchuanAdvertiserStatus.AdvertiserId {
				l.Status = qianchuanAdvertiserStatus.Status

				isNotExist = false

				break
			}
		}

		if isNotExist {
			list = append(list, qianchuanAdvertiserStatus.ToDomain())
		}
	}

	return list, nil
}

func (qasr *qianchuanAdvertiserStatusRepo) Count(ctx context.Context, companyId, advertiserId uint64, day uint32) (int64, error) {
	var count int64

	if result := qasr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("advertiser_id = ?", advertiserId).
		Where("day <= ?", day).
		Model(&QianchuanAdvertiserStatus{}).
		Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (qasr *qianchuanAdvertiserStatusRepo) Save(ctx context.Context, in *domain.QianchuanAdvertiserStatus) (*domain.QianchuanAdvertiserStatus, error) {
	qianchuanAdvertiserStatus := &QianchuanAdvertiserStatus{
		CompanyId:    in.CompanyId,
		AdvertiserId: in.AdvertiserId,
		Status:       in.Status,
		Day:          in.Day,
		CreateTime:   in.CreateTime,
		UpdateTime:   in.UpdateTime,
	}

	if result := qasr.data.DB(ctx).Create(qianchuanAdvertiserStatus); result.Error != nil {
		return nil, result.Error
	}

	return qianchuanAdvertiserStatus.ToDomain(), nil
}

func (qasr *qianchuanAdvertiserStatusRepo) Update(ctx context.Context, in *domain.QianchuanAdvertiserStatus) (*domain.QianchuanAdvertiserStatus, error) {
	qianchuanAdvertiserStatus := map[string]interface{}{
		"status":      in.Status,
		"update_time": in.UpdateTime,
	}

	if result := qasr.data.DB(ctx).
		Model(QianchuanAdvertiserStatus{}).
		Where("company_id = ?", in.CompanyId).
		Where("advertiser_id = ?", in.AdvertiserId).
		Where("day = ?", in.Day).
		Updates(qianchuanAdvertiserStatus); result.Error != nil {
		return nil, result.Error
	}

	return in, nil
}

func (qasr *qianchuanAdvertiserStatusRepo) DeleteByCompanyIdAndAdvertiserIdAndDay(ctx context.Context, companyId, advertiserId uint64, day uint32) error {
	if result := qasr.data.DB(ctx).
		Where("company_id = ?", companyId).
		Where("advertiser_id = ?", advertiserId).
		Where("day = ?", day).
		Delete(&QianchuanAdvertiserStatus{}); result.Error != nil {
		return result.Error
	}

	return nil
}
