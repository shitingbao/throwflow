package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库用户当前企业表
type CompanyUserCompany struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	Phone      string    `gorm:"column:phone;type:varchar(20);not null;comment:企业库用户手机号"`
	CompanyId  uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyUserCompany) TableName() string {
	return "company_company_user_company"
}

type companyUserCompanyRepo struct {
	data *Data
	log  *log.Helper
}

func (cuc *CompanyUserCompany) ToDomain() *domain.CompanyUserCompany {
	companyUserCompany := &domain.CompanyUserCompany{
		Id:         cuc.Id,
		Phone:      cuc.Phone,
		CompanyId:  cuc.CompanyId,
		CreateTime: cuc.CreateTime,
		UpdateTime: cuc.UpdateTime,
	}

	return companyUserCompany
}

func NewCompanyUserCompanyRepo(data *Data, logger log.Logger) biz.CompanyUserCompanyRepo {
	return &companyUserCompanyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cucr *companyUserCompanyRepo) GetByPhone(ctx context.Context, phone string) (*domain.CompanyUserCompany, error) {
	companyUserCompany := &CompanyUserCompany{}

	if result := cucr.data.db.WithContext(ctx).Where("phone = ?", phone).First(companyUserCompany); result.Error != nil {
		return nil, result.Error
	}

	return companyUserCompany.ToDomain(), nil
}

func (cucr *companyUserCompanyRepo) Save(ctx context.Context, in *domain.CompanyUserCompany) (*domain.CompanyUserCompany, error) {
	companyUserCompany := &CompanyUserCompany{
		Phone:      in.Phone,
		CompanyId:  in.CompanyId,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := cucr.data.DB(ctx).Create(companyUserCompany); result.Error != nil {
		return nil, result.Error
	}

	return companyUserCompany.ToDomain(), nil
}

func (cucr *companyUserCompanyRepo) Update(ctx context.Context, in *domain.CompanyUserCompany) (*domain.CompanyUserCompany, error) {
	companyUserCompany := &CompanyUserCompany{
		Id:         in.Id,
		Phone:      in.Phone,
		CompanyId:  in.CompanyId,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := cucr.data.DB(ctx).Save(companyUserCompany); result.Error != nil {
		return nil, result.Error
	}

	return companyUserCompany.ToDomain(), nil
}

func (cucr *companyUserCompanyRepo) DeleteByCompanyId(ctx context.Context, companyId uint64) error {
	if result := cucr.data.DB(ctx).Where("company_id = ?", companyId).Delete(&CompanyUserCompany{}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (cucr *companyUserCompanyRepo) Delete(ctx context.Context, in *domain.CompanyUserCompany) error {
	companyUserCompany := &CompanyUserCompany{
		Id:         in.Id,
		Phone:      in.Phone,
		CompanyId:  in.CompanyId,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := cucr.data.DB(ctx).Delete(companyUserCompany); result.Error != nil {
		return result.Error
	}

	return nil
}
