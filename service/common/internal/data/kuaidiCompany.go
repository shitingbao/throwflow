package data

import (
	"common/internal/biz"
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"unicode/utf8"
)

// 快递公司编码表
type KuaidiCompany struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	Name       string    `gorm:"column:name;type:char(150);not null;comment:快递公司名称"`
	Code       string    `gorm:"column:code;type:char(100);not null;comment:快递公司简码"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (KuaidiCompany) TableName() string {
	return "common_kuaidi_company"
}

type kuaidiCompanyRepo struct {
	data *Data
	log  *log.Helper
}

func (kc *KuaidiCompany) ToDomain() *domain.KuaidiCompany {
	return &domain.KuaidiCompany{
		Id:         kc.Id,
		Name:       kc.Name,
		Code:       kc.Code,
		CreateTime: kc.CreateTime,
		UpdateTime: kc.UpdateTime,
	}
}

func NewKuaidiCompanyRepo(data *Data, logger log.Logger) biz.KuaidiCompanyRepo {
	return &kuaidiCompanyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (kcr *kuaidiCompanyRepo) Get(ctx context.Context, code string) (*domain.KuaidiCompany, error) {
	kuaidiCompany := &KuaidiCompany{}

	if result := kcr.data.db.WithContext(ctx).
		Where("code = ?", code).
		First(kuaidiCompany); result.Error != nil {
		return nil, result.Error
	}

	return kuaidiCompany.ToDomain(), nil
}

func (kcr *kuaidiCompanyRepo) List(ctx context.Context, pageNum, pageSize int, keyword string) ([]*domain.KuaidiCompany, error) {
	var kuaidiCompanies []KuaidiCompany
	list := make([]*domain.KuaidiCompany, 0)

	db := kcr.data.db.WithContext(ctx)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("name like ?", "%"+keyword+"%")
	}

	if result := db.Order("id ASC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&kuaidiCompanies); result.Error != nil {
		return nil, result.Error
	}

	for _, kuaidiCompany := range kuaidiCompanies {
		list = append(list, kuaidiCompany.ToDomain())
	}

	return list, nil
}

func (kcr *kuaidiCompanyRepo) Count(ctx context.Context, keyword string) (int64, error) {
	var count int64

	db := kcr.data.db.WithContext(ctx)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("name like ?", "%"+keyword+"%")
	}

	if result := db.Model(&KuaidiCompany{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}
