package data

import (
	"context"
	"douyin/internal/domain"
	"douyin/internal/pkg/event/event"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"douyin/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// 千川广告账户表
type QianchuanAdvertiser struct {
	Id             uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	AppId          string    `gorm:"column:app_id;type:varchar(20);not null;comment:应用ID"`
	CompanyId      uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:公司ID"`
	AccountId      uint64    `gorm:"column:account_id;type:bigint(20) UNSIGNED;not null;comment:千川账户ID"`
	AdvertiserId   uint64    `gorm:"column:advertiser_id;type:bigint(20) UNSIGNED;not null;comment:广告主ID"`
	AdvertiserName string    `gorm:"column:advertiser_name;type:varchar(250);not null;comment:广告主名称"`
	CompanyName    string    `gorm:"column:company_name;type:varchar(250);not null;comment:公司名称"`
	Status         uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态: 1 启用 0 禁用"`
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (QianchuanAdvertiser) TableName() string {
	return "douyin_qianchuan_advertiser"
}

type qianchuanAdvertiserRepo struct {
	data *Data
	log  *log.Helper
}

func (qa *QianchuanAdvertiser) ToDomain() *domain.QianchuanAdvertiser {
	return &domain.QianchuanAdvertiser{
		Id:             qa.Id,
		AppId:          qa.AppId,
		CompanyId:      qa.CompanyId,
		AccountId:      qa.AccountId,
		AdvertiserId:   qa.AdvertiserId,
		AdvertiserName: qa.AdvertiserName,
		CompanyName:    qa.CompanyName,
		Status:         qa.Status,
		CreateTime:     qa.CreateTime,
		UpdateTime:     qa.UpdateTime,
	}
}

func NewQianchuanAdvertiserRepo(data *Data, logger log.Logger) biz.QianchuanAdvertiserRepo {
	return &qianchuanAdvertiserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qar *qianchuanAdvertiserRepo) List(ctx context.Context, pageNum, pageSize int, companyId uint64, keyword, advertiserIds, status string) ([]*domain.QianchuanAdvertiser, error) {
	var qianchuanAdvertisers []QianchuanAdvertiser
	list := make([]*domain.QianchuanAdvertiser, 0)

	db := qar.data.db.WithContext(ctx)

	if companyId > 0 {
		db = db.Where("company_id = ?", companyId)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		keywords := strings.Split(keyword, " ")

		if len(keywords) == 1 {
			db = db.Where("(advertiser_id = ? or advertiser_name like ?)", keyword, "%"+keyword+"%")
		} else {
			wheres := make([]string, 0)

			for _, key := range keywords {
				wheres = append(wheres, "advertiser_name like '%"+key+"%'")
			}

			db = db.Where("(advertiser_id = ? or ( ? ))", keyword, strings.Join(wheres, " and "))
		}
	}

	if len(advertiserIds) > 0 {
		db = db.Where("advertiser_id in ?", strings.Split(advertiserIds, ","))
	}

	if len(status) > 0 {
		if istatus, err := strconv.Atoi(status); err == nil {
			db = db.Where("status = ?", uint8(istatus))
		}
	}

	if pageNum > 0 {
		if result := db.Order("status DESC,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&qianchuanAdvertisers); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("status DESC,id DESC").
			Find(&qianchuanAdvertisers); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		list = append(list, qianchuanAdvertiser.ToDomain())
	}

	return list, nil
}

func (qar *qianchuanAdvertiserRepo) Count(ctx context.Context, companyId uint64, keyword, advertiserIds, status string) (int64, error) {
	var count int64

	db := qar.data.db.WithContext(ctx)

	if companyId > 0 {
		db = db.Where("company_id = ?", companyId)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		keywords := strings.Split(keyword, " ")

		if len(keywords) == 1 {
			db = db.Where("(advertiser_id = ? or advertiser_name like ?)", keyword, "%"+keyword+"%")
		} else {
			wheres := make([]string, 0)

			for _, key := range keywords {
				wheres = append(wheres, "advertiser_name like '%"+key+"%'")
			}

			db = db.Where("(advertiser_id = ? or ( ? ))", keyword, strings.Join(wheres, " and "))
		}
	}

	if len(advertiserIds) > 0 {
		db = db.Where("advertiser_id in ?", strings.Split(advertiserIds, ","))
	}

	if len(status) > 0 {
		if istatus, err := strconv.Atoi(status); err == nil {
			db = db.Where("status = ?", uint8(istatus))
		}
	}

	if result := db.Model(&QianchuanAdvertiser{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (qar *qianchuanAdvertiserRepo) Statistics(ctx context.Context, companyId uint64, status uint8) (int64, error) {
	var count int64

	if result := qar.data.db.WithContext(ctx).Model(&QianchuanAdvertiser{}).Where("company_id = ?", companyId).Where("status = ?", status).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (qar *qianchuanAdvertiserRepo) Get(ctx context.Context, advertiserId uint64) (*domain.QianchuanAdvertiser, error) {
	qianchuanAdvertiser := &QianchuanAdvertiser{}

	if result := qar.data.db.WithContext(ctx).Where("advertiser_id = ?", advertiserId).First(qianchuanAdvertiser); result.Error != nil {
		return nil, result.Error
	}

	return qianchuanAdvertiser.ToDomain(), nil
}

// /////////////////////需要修改////////////////////////////////
func (qar *qianchuanAdvertiserRepo) GetById(ctx context.Context, companyId, advertiserId uint64) (*domain.QianchuanAdvertiser, error) {
	qianchuanAdvertiser := &QianchuanAdvertiser{}

	if result := qar.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("advertiser_id = ?", advertiserId).First(qianchuanAdvertiser); result.Error != nil {
		return nil, result.Error
	}

	return qianchuanAdvertiser.ToDomain(), nil
}

func (qar *qianchuanAdvertiserRepo) GetByCompanyIdAndAdvertiserId(ctx context.Context, companyId, advertiserId uint64) (*domain.QianchuanAdvertiser, error) {
	qianchuanAdvertiser := &QianchuanAdvertiser{}

	if result := qar.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("advertiser_id = ?", advertiserId).First(qianchuanAdvertiser); result.Error != nil {
		return nil, result.Error
	}

	return qianchuanAdvertiser.ToDomain(), nil
}

func (qar *qianchuanAdvertiserRepo) ListByAdvertiserId(ctx context.Context, advertiserId uint64) ([]*domain.QianchuanAdvertiser, error) {
	var qianchuanAdvertisers []QianchuanAdvertiser
	list := make([]*domain.QianchuanAdvertiser, 0)

	db := qar.data.db.WithContext(ctx)

	if advertiserId > 0 {
		db = db.Where("advertiser_id = ?", advertiserId)
	}

	db = db.Where("status = 1")

	if result := db.Order("update_time DESC,id DESC").
		Find(&qianchuanAdvertisers); result.Error != nil {
		return nil, result.Error
	}

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		list = append(list, qianchuanAdvertiser.ToDomain())
	}

	return list, nil
}

func (qar *qianchuanAdvertiserRepo) ListByAppId(ctx context.Context, appId string, status uint8) ([]*domain.QianchuanAdvertiser, error) {
	var qianchuanAdvertisers []QianchuanAdvertiser
	list := make([]*domain.QianchuanAdvertiser, 0)

	if result := qar.data.db.WithContext(ctx).
		Where("app_id = ?", appId).
		Where("status = ?", status).
		Order("update_time DESC,id DESC").
		Find(&qianchuanAdvertisers); result.Error != nil {
		return nil, result.Error
	}

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		list = append(list, qianchuanAdvertiser.ToDomain())
	}

	return list, nil
}

func (qar *qianchuanAdvertiserRepo) Save(ctx context.Context, in *domain.QianchuanAdvertiser) (*domain.QianchuanAdvertiser, error) {
	qianchuanAdvertiser := &QianchuanAdvertiser{
		AppId:          in.AppId,
		CompanyId:      in.CompanyId,
		AccountId:      in.AccountId,
		AdvertiserId:   in.AdvertiserId,
		AdvertiserName: in.AdvertiserName,
		CompanyName:    in.CompanyName,
		Status:         in.Status,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := qar.data.DB(ctx).Create(qianchuanAdvertiser); result.Error != nil {
		return nil, result.Error
	}

	return qianchuanAdvertiser.ToDomain(), nil
}

func (qar *qianchuanAdvertiserRepo) Update(ctx context.Context, in *domain.QianchuanAdvertiser) (*domain.QianchuanAdvertiser, error) {
	qianchuanAdvertiser := &QianchuanAdvertiser{
		Id:             in.Id,
		AppId:          in.AppId,
		CompanyId:      in.CompanyId,
		AccountId:      in.AccountId,
		AdvertiserId:   in.AdvertiserId,
		AdvertiserName: in.AdvertiserName,
		CompanyName:    in.CompanyName,
		Status:         in.Status,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := qar.data.DB(ctx).Save(qianchuanAdvertiser); result.Error != nil {
		return nil, result.Error
	}

	return qianchuanAdvertiser.ToDomain(), nil
}

func (qar *qianchuanAdvertiserRepo) DeleteByCompanyIdAndAccountId(ctx context.Context, companyId, accountId uint64) error {
	if result := qar.data.DB(ctx).Where("company_id = ?", companyId).Where("account_id = ?", accountId).Delete(&QianchuanAdvertiser{}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (qar *qianchuanAdvertiserRepo) Send(ctx context.Context, message event.Event) error {
	if err := qar.data.kafka.Send(ctx, message); err != nil {
		return err
	}

	return nil
}
