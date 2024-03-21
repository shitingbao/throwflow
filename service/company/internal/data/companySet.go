package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库设置表
type CompanySet struct {
	CompanyId  uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	Day        uint32    `gorm:"column:day;type:int(11);not null;comment:修改日期"`
	SetKey     string    `gorm:"column:set_key;type:varchar(20);not null;comment:设置键,sampleThreshold：寄样门槛"`
	SetValue   string    `gorm:"column:set_value;type:varchar(250);not null;comment:设置值"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanySet) TableName() string {
	return "company_company_set"
}

type companySetRepo struct {
	data *Data
	log  *log.Helper
}

func (cs *CompanySet) ToDomain() *domain.CompanySet {
	companySet := &domain.CompanySet{
		CompanyId:  cs.CompanyId,
		Day:        cs.Day,
		SetKey:     cs.SetKey,
		SetValue:   cs.SetValue,
		CreateTime: cs.CreateTime,
		UpdateTime: cs.UpdateTime,
	}

	return companySet
}

func NewCompanySetRepo(data *Data, logger log.Logger) biz.CompanySetRepo {
	return &companySetRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (csr *companySetRepo) GetByCompanyIdAndDayAndSetKey(ctx context.Context, companyId uint64, day uint32, setKey string) (*domain.CompanySet, error) {
	companySet := &CompanySet{}

	if result := csr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("day <= ?", day).
		Where("set_key = ?", setKey).
		Order("day desc").
		First(companySet); result.Error != nil {
		return nil, result.Error
	}

	return companySet.ToDomain(), nil
}

func (csr *companySetRepo) Save(ctx context.Context, in *domain.CompanySet) (*domain.CompanySet, error) {
	companySet := &CompanySet{
		CompanyId:  in.CompanyId,
		Day:        in.Day,
		SetKey:     in.SetKey,
		SetValue:   in.SetValue,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := csr.data.DB(ctx).Create(companySet); result.Error != nil {
		return nil, result.Error
	}

	return companySet.ToDomain(), nil
}

func (csr *companySetRepo) Update(ctx context.Context, in *domain.CompanySet) error {
	if result := csr.data.DB(ctx).Model(CompanySet{}).
		Where("company_id = ?", in.CompanyId).
		Where("set_key = ?", in.SetKey).
		Where("day = ?", in.Day).
		Updates(map[string]interface{}{"set_value": in.SetValue, "update_time": in.UpdateTime}); result.Error != nil {
		return result.Error
	}

	return nil
}
