package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"time"
)

// 企业机构表
type CompanyOrganization struct {
	Id                            uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	OrganizationName              string    `gorm:"column:organization_name;type:varchar(150);not null;comment:机构名称"`
	OrganizationMcn               string    `gorm:"column:organization_mcn;type:varchar(250);not null;comment:mcn机构"`
	CompanyName                   string    `gorm:"column:company_name;type:varchar(150);not null;comment:公司名称"`
	BankCode                      string    `gorm:"column:bank_code;type:varchar(150);not null;comment:银行账户"`
	BankDeposit                   string    `gorm:"column:bank_deposit;type:varchar(150);not null;comment:开户行"`
	OrganizationLogoUrl           string    `gorm:"column:organization_logo_url;type:varchar(250);not null;comment:机构LOGO URL"`
	OrganizationCode              string    `gorm:"column:organization_code;type:varchar(10);not null;comment:机构简码"`
	OrganizationQrCodeUrl         string    `gorm:"column:organization_qr_code_url;type:varchar(250);not null;comment:机构小程序码 URL"`
	OrganizationShortUrl          string    `gorm:"column:organization_short_url;type:varchar(250);not null;comment:机构短链接"`
	OrganizationCommission        string    `gorm:"column:organization_commission;type:text;not null;comment:机构分佣设置"`
	OrganizationColonelCommission string    `gorm:"column:organization_colonel_commission;type:text;not null;comment:团长分佣设置"`
	OrganizationCourse            string    `gorm:"column:organization_course;type:text;not null;comment:机构课程"`
	CreateTime                    time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime                    time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyOrganization) TableName() string {
	return "company_organization"
}

type companyOrganizationRepo struct {
	data *Data
	log  *log.Helper
}

func (co *CompanyOrganization) ToDomain() *domain.CompanyOrganization {
	companyOrganization := &domain.CompanyOrganization{
		Id:                            co.Id,
		OrganizationName:              co.OrganizationName,
		OrganizationMcn:               co.OrganizationMcn,
		CompanyName:                   co.CompanyName,
		BankCode:                      co.BankCode,
		BankDeposit:                   co.BankDeposit,
		OrganizationLogoUrl:           co.OrganizationLogoUrl,
		OrganizationCode:              co.OrganizationCode,
		OrganizationQrCodeUrl:         co.OrganizationQrCodeUrl,
		OrganizationShortUrl:          co.OrganizationShortUrl,
		OrganizationCommission:        co.OrganizationCommission,
		OrganizationColonelCommission: co.OrganizationColonelCommission,
		OrganizationCourse:            co.OrganizationCourse,
		CreateTime:                    co.CreateTime,
		UpdateTime:                    co.UpdateTime,
	}

	return companyOrganization
}

func NewCompanyOrganizationRepo(data *Data, logger log.Logger) biz.CompanyOrganizationRepo {
	return &companyOrganizationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cor *companyOrganizationRepo) GetById(ctx context.Context, organizationId uint64) (*domain.CompanyOrganization, error) {
	companyOrganization := &CompanyOrganization{}

	db := cor.data.db.WithContext(ctx).
		Where("id = ?", organizationId)

	if result := db.First(companyOrganization); result.Error != nil {
		return nil, result.Error
	}

	return companyOrganization.ToDomain(), nil
}

func (cor *companyOrganizationRepo) GetByOrganizationCode(ctx context.Context, organizationCode string) (*domain.CompanyOrganization, error) {
	companyOrganization := &CompanyOrganization{}

	db := cor.data.db.WithContext(ctx).
		Where("organization_code = ?", organizationCode)

	if result := db.First(companyOrganization); result.Error != nil {
		return nil, result.Error
	}

	return companyOrganization.ToDomain(), nil
}

func (cor *companyOrganizationRepo) List(ctx context.Context, pageNum, pageSize int) ([]*domain.CompanyOrganization, error) {
	var companyOrganizations []CompanyOrganization
	list := make([]*domain.CompanyOrganization, 0)

	db := cor.data.db.WithContext(ctx)

	if pageNum == 0 {
		if result := db.Order("create_time DESC,id DESC").
			Find(&companyOrganizations); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("create_time DESC,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&companyOrganizations); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, companyOrganization := range companyOrganizations {
		list = append(list, companyOrganization.ToDomain())
	}

	return list, nil
}

func (cor *companyOrganizationRepo) Count(ctx context.Context) (int64, error) {
	var count int64

	if result := cor.data.db.WithContext(ctx).Model(&CompanyOrganization{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cor *companyOrganizationRepo) Save(ctx context.Context, in *domain.CompanyOrganization) (*domain.CompanyOrganization, error) {
	companyOrganization := &CompanyOrganization{
		OrganizationName:              in.OrganizationName,
		OrganizationMcn:               in.OrganizationMcn,
		CompanyName:                   in.CompanyName,
		BankCode:                      in.BankCode,
		BankDeposit:                   in.BankDeposit,
		OrganizationLogoUrl:           in.OrganizationLogoUrl,
		OrganizationCode:              in.OrganizationCode,
		OrganizationQrCodeUrl:         in.OrganizationQrCodeUrl,
		OrganizationShortUrl:          in.OrganizationShortUrl,
		OrganizationCommission:        in.OrganizationCommission,
		OrganizationColonelCommission: in.OrganizationColonelCommission,
		OrganizationCourse:            in.OrganizationCourse,
		CreateTime:                    in.CreateTime,
		UpdateTime:                    in.UpdateTime,
	}

	if result := cor.data.DB(ctx).Create(companyOrganization); result.Error != nil {
		return nil, result.Error
	}

	return companyOrganization.ToDomain(), nil
}

func (cor *companyOrganizationRepo) Update(ctx context.Context, in *domain.CompanyOrganization) (*domain.CompanyOrganization, error) {
	companyOrganization := &CompanyOrganization{
		Id:                            in.Id,
		OrganizationName:              in.OrganizationName,
		OrganizationMcn:               in.OrganizationMcn,
		CompanyName:                   in.CompanyName,
		BankCode:                      in.BankCode,
		BankDeposit:                   in.BankDeposit,
		OrganizationLogoUrl:           in.OrganizationLogoUrl,
		OrganizationCode:              in.OrganizationCode,
		OrganizationQrCodeUrl:         in.OrganizationQrCodeUrl,
		OrganizationShortUrl:          in.OrganizationShortUrl,
		OrganizationCommission:        in.OrganizationCommission,
		OrganizationColonelCommission: in.OrganizationColonelCommission,
		OrganizationCourse:            in.OrganizationCourse,
		CreateTime:                    in.CreateTime,
		UpdateTime:                    in.UpdateTime,
	}

	if result := cor.data.DB(ctx).Save(companyOrganization); result.Error != nil {
		return nil, result.Error
	}

	return companyOrganization.ToDomain(), nil
}

func (cor *companyOrganizationRepo) PutContent(ctx context.Context, fileName string, content io.Reader) (*ctos.PutObjectV2Output, error) {
	for _, ltos := range cor.data.toses {
		if ltos.name == "organization" {
			output, err := ltos.tos.PutContent(ctx, fileName, content)

			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("tos is not exist")
}
