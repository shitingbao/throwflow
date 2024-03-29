package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 巨量引擎授权账户token表
type OceanengineAccountToken struct {
	Id                    uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	AppId                 string    `gorm:"column:app_id;type:varchar(20);not null;comment:应用ID"`
	CompanyId             uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:公司ID"`
	AccountId             uint64    `gorm:"column:account_id;type:bigint(20) UNSIGNED;not null;comment:千川账户ID"`
	AccessToken           string    `gorm:"column:access_token;type:varchar(50);not null;comment:用于验证权限的token"`
	ExpiresIn             uint32    `gorm:"column:expires_in;type:int(11) UNSIGNED;not null;comment:access_token剩余有效时间,单位(秒)"`
	RefreshToken          string    `gorm:"column:refresh_token;type:varchar(50);not null;comment:刷新access_token,用于获取新的access_token和refresh_token，并且刷新过期时间"`
	RefreshTokenExpiresIn uint32    `gorm:"column:refresh_token_expires_in;type:int(11) UNSIGNED;not null;comment:refresh_token剩余有效时间,单位(秒)"`
	CreateTime            time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime            time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (OceanengineAccountToken) TableName() string {
	return "douyin_oceanengine_account_token"
}

type oceanengineAccountTokenRepo struct {
	data *Data
	log  *log.Helper
}

func (oat *OceanengineAccountToken) ToDomain() *domain.OceanengineAccountToken {
	return &domain.OceanengineAccountToken{
		Id:                    oat.Id,
		AppId:                 oat.AppId,
		CompanyId:             oat.CompanyId,
		AccountId:             oat.AccountId,
		AccessToken:           oat.AccessToken,
		ExpiresIn:             oat.ExpiresIn,
		RefreshToken:          oat.RefreshToken,
		RefreshTokenExpiresIn: oat.RefreshTokenExpiresIn,
		CreateTime:            oat.CreateTime,
		UpdateTime:            oat.UpdateTime,
	}
}

func NewOceanengineAccountTokenRepo(data *Data, logger log.Logger) biz.OceanengineAccountTokenRepo {
	return &oceanengineAccountTokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (oatr *oceanengineAccountTokenRepo) GetByCompanyIdAndAccountId(ctx context.Context, companyId, accountId uint64) (*domain.OceanengineAccountToken, error) {
	oceanengineAccountToken := &OceanengineAccountToken{}

	if result := oatr.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("account_id = ?", accountId).First(oceanengineAccountToken); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineAccountToken.ToDomain(), nil
}

func (oatr *oceanengineAccountTokenRepo) List(ctx context.Context) ([]*domain.OceanengineAccountToken, error) {
	var oceanengineAccountTokens []OceanengineAccountToken
	list := make([]*domain.OceanengineAccountToken, 0)

	if result := oatr.data.db.WithContext(ctx).
		Find(&oceanengineAccountTokens); result.Error != nil {
		return nil, result.Error
	}

	for _, oceanengineAccountToken := range oceanengineAccountTokens {
		list = append(list, oceanengineAccountToken.ToDomain())
	}

	return list, nil
}

func (oatr *oceanengineAccountTokenRepo) ListByAppId(ctx context.Context, appId string) ([]*domain.OceanengineAccountToken, error) {
	var oceanengineAccountTokens []OceanengineAccountToken
	list := make([]*domain.OceanengineAccountToken, 0)

	if result := oatr.data.db.WithContext(ctx).
		Where("app_id = ?", appId).
		Find(&oceanengineAccountTokens); result.Error != nil {
		return nil, result.Error
	}

	for _, oceanengineAccountToken := range oceanengineAccountTokens {
		list = append(list, oceanengineAccountToken.ToDomain())
	}

	return list, nil
}

func (oatr *oceanengineAccountTokenRepo) ListByCompanyId(ctx context.Context, companyId uint64) ([]*domain.OceanengineAccountToken, error) {
	var oceanengineAccountTokens []OceanengineAccountToken
	list := make([]*domain.OceanengineAccountToken, 0)

	if result := oatr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Find(&oceanengineAccountTokens); result.Error != nil {
		return nil, result.Error
	}

	for _, oceanengineAccountToken := range oceanengineAccountTokens {
		list = append(list, oceanengineAccountToken.ToDomain())
	}

	return list, nil
}

func (oatr *oceanengineAccountTokenRepo) Save(ctx context.Context, in *domain.OceanengineAccountToken) (*domain.OceanengineAccountToken, error) {
	oceanengineAccountToken := &OceanengineAccountToken{
		AppId:                 in.AppId,
		CompanyId:             in.CompanyId,
		AccessToken:           in.AccessToken,
		ExpiresIn:             in.ExpiresIn,
		RefreshToken:          in.RefreshToken,
		RefreshTokenExpiresIn: in.RefreshTokenExpiresIn,
		CreateTime:            in.CreateTime,
		UpdateTime:            in.UpdateTime,
		AccountId:             in.AccountId,
	}

	if result := oatr.data.DB(ctx).Create(oceanengineAccountToken); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineAccountToken.ToDomain(), nil
}

func (oatr *oceanengineAccountTokenRepo) Update(ctx context.Context, in *domain.OceanengineAccountToken) (*domain.OceanengineAccountToken, error) {
	oceanengineAccountToken := &OceanengineAccountToken{
		Id:                    in.Id,
		AppId:                 in.AppId,
		CompanyId:             in.CompanyId,
		AccountId:             in.AccountId,
		AccessToken:           in.AccessToken,
		ExpiresIn:             in.ExpiresIn,
		RefreshToken:          in.RefreshToken,
		RefreshTokenExpiresIn: in.RefreshTokenExpiresIn,
		CreateTime:            in.CreateTime,
		UpdateTime:            in.UpdateTime,
	}

	if result := oatr.data.db.WithContext(ctx).Save(oceanengineAccountToken); result.Error != nil {
		return nil, result.Error
	}

	return oceanengineAccountToken.ToDomain(), nil
}

func (oatr *oceanengineAccountTokenRepo) DeleteByCompanyIdAndAccountId(ctx context.Context, companyId, accountId uint64) error {
	if result := oatr.data.DB(ctx).Where("company_id = ?", companyId).Where("account_id = ?", accountId).Delete(&OceanengineAccountToken{}); result.Error != nil {
		return result.Error
	}

	return nil
}
