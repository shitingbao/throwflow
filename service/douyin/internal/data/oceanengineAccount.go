package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 巨量引擎授权账户表
type OceanengineAccount struct {
	Id          uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	AppId       string    `gorm:"column:app_id;type:varchar(20);not null;comment:应用ID"`
	CompanyId   uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:公司ID"`
	AccountId   uint64    `gorm:"column:account_id;type:bigint(20) UNSIGNED;not null;comment:千川账户ID"`
	AccountName string    `gorm:"column:account_name;type:varchar(250);not null;comment:千川账户名称"`
	AccountRole string    `gorm:"column:account_role;type:varchar(50);not null;comment:授权账号角色，返回值：PLATFORM_ROLE_QIANCHUAN_AGENT代理商账户、PLATFORM_ROLE_SHOP_ACCOUNT 店铺账户"`
	IsValid     uint8     `gorm:"column:is_valid;type:tinyint(3) UNSIGNED;not null;default:0;comment:授权有效性，1：授权有效，0：授权无效"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (OceanengineAccount) TableName() string {
	return "douyin_oceanengine_account"
}

type oceanengineAccountRepo struct {
	data *Data
	log  *log.Helper
}

func (oa *OceanengineAccount) ToDomain() *domain.OceanengineAccount {
	return &domain.OceanengineAccount{
		Id:          oa.Id,
		AppId:       oa.AppId,
		CompanyId:   oa.CompanyId,
		AccountId:   oa.AccountId,
		AccountName: oa.AccountName,
		AccountRole: oa.AccountRole,
		IsValid:     oa.IsValid,
		CreateTime:  oa.CreateTime,
		UpdateTime:  oa.UpdateTime,
	}
}

func NewOceanengineAccountRepo(data *Data, logger log.Logger) biz.OceanengineAccountRepo {
	return &oceanengineAccountRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (oar *oceanengineAccountRepo) GetById(ctx context.Context, companyId, accountId uint64) (*domain.OceanengineAccount, error) {
	oceanengineAccount := &OceanengineAccount{}

	if result := oar.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("account_id = ?", accountId).First(oceanengineAccount); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineAccount.ToDomain(), nil
}

func (oar *oceanengineAccountRepo) GetBycompanyIdAndAccountId(ctx context.Context, companyId, accountId uint64) (*domain.OceanengineAccount, error) {
	oceanengineAccount := &OceanengineAccount{}

	if result := oar.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("account_id = ?", accountId).First(oceanengineAccount); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineAccount.ToDomain(), nil
}

func (oar *oceanengineAccountRepo) Save(ctx context.Context, in *domain.OceanengineAccount) (*domain.OceanengineAccount, error) {
	oceanengineAccount := &OceanengineAccount{
		AppId:       in.AppId,
		CompanyId:   in.CompanyId,
		AccountId:   in.AccountId,
		AccountName: in.AccountName,
		AccountRole: in.AccountRole,
		IsValid:     in.IsValid,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
	}

	if result := oar.data.DB(ctx).Create(oceanengineAccount); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineAccount.ToDomain(), nil
}

func (oar *oceanengineAccountRepo) Update(ctx context.Context, in *domain.OceanengineAccount) (*domain.OceanengineAccount, error) {
	oceanengineAccount := &OceanengineAccount{
		Id:          in.Id,
		AppId:       in.AppId,
		CompanyId:   in.CompanyId,
		AccountId:   in.AccountId,
		AccountName: in.AccountName,
		AccountRole: in.AccountRole,
		IsValid:     in.IsValid,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
	}

	if result := oar.data.db.WithContext(ctx).Save(oceanengineAccount); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineAccount.ToDomain(), nil
}

func (oar *oceanengineAccountRepo) DeleteByCompanyIdAndAccountId(ctx context.Context, companyId, accountId uint64) error {
	if result := oar.data.DB(ctx).Where("company_id = ?", companyId).Where("account_id = ?", accountId).Delete(&OceanengineAccount{}); result.Error != nil {
		return result.Error
	}

	return nil
}
