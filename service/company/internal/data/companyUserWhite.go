package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库用户白名单表
type CompanyUserWhite struct {
	Phone      string    `gorm:"column:phone;type:varchar(20);not null;comment:手机号"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyUserWhite) TableName() string {
	return "company_company_user_white"
}

type companyUserWhiteRepo struct {
	data *Data
	log  *log.Helper
}

func (cuw *CompanyUserWhite) ToDomain() *domain.CompanyUserWhite {
	companyUserWhite := &domain.CompanyUserWhite{
		Phone:      cuw.Phone,
		CreateTime: cuw.CreateTime,
		UpdateTime: cuw.UpdateTime,
	}

	return companyUserWhite
}

func NewCompanyUserWhiteRepo(data *Data, logger log.Logger) biz.CompanyUserWhiteRepo {
	return &companyUserWhiteRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cuw *companyUserWhiteRepo) GetByPhone(ctx context.Context, phone string) (*domain.CompanyUserWhite, error) {
	companyUserWhite := &CompanyUserWhite{}

	if result := cuw.data.db.WithContext(ctx).Where("phone = ?", phone).First(companyUserWhite); result.Error != nil {
		return nil, result.Error
	}

	return companyUserWhite.ToDomain(), nil
}

func (cuw *companyUserWhiteRepo) Save(ctx context.Context, in *domain.CompanyUserWhite) (*domain.CompanyUserWhite, error) {
	companyUserWhite := &CompanyUserWhite{
		Phone:      in.Phone,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := cuw.data.DB(ctx).Create(companyUserWhite); result.Error != nil {
		return nil, result.Error
	}

	return companyUserWhite.ToDomain(), nil
}

func (cuw *companyUserWhiteRepo) Delete(ctx context.Context, phone string) error {
	if result := cuw.data.DB(ctx).Where("phone = ?", phone).Delete(&CompanyUserWhite{}); result.Error != nil {
		return result.Error
	}

	return nil
}
