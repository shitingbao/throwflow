package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库团队绩效规则设置表
type CompanyPerformanceRule struct {
	Id               uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyId        uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	PerformanceName  string    `gorm:"column:performance_name;type:varchar(50);not null;comment:规则名称"`
	AdvertiserIds    string    `gorm:"column:advertiser_ids;type:text;not null;comment:使用账户对象"`
	PerformanceRules string    `gorm:"column:performance_rules;type:text;not null;comment:提成算法公式"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyPerformanceRule) TableName() string {
	return "company_company_performance_rule"
}

type companyPerformanceRuleRepo struct {
	data *Data
	log  *log.Helper
}

func (cpr *CompanyPerformanceRule) ToDomain() *domain.CompanyPerformanceRule {
	companyPerformanceRule := &domain.CompanyPerformanceRule{
		Id:               cpr.Id,
		CompanyId:        cpr.CompanyId,
		PerformanceName:  cpr.PerformanceName,
		AdvertiserIds:    cpr.AdvertiserIds,
		PerformanceRules: cpr.PerformanceRules,
		CreateTime:       cpr.CreateTime,
		UpdateTime:       cpr.UpdateTime,
	}

	return companyPerformanceRule
}

func NewCompanyPerformanceRuleRepo(data *Data, logger log.Logger) biz.CompanyPerformanceRuleRepo {
	return &companyPerformanceRuleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cpr *companyPerformanceRuleRepo) GetById(ctx context.Context, id, companyId uint64) (*domain.CompanyPerformanceRule, error) {
	companyPerformanceRule := &CompanyPerformanceRule{}

	if result := cpr.data.db.WithContext(ctx).Where("company_id = ?", companyId).First(companyPerformanceRule, id); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceRule.ToDomain(), nil
}

func (cpr *companyPerformanceRuleRepo) List(ctx context.Context, companyId uint64) ([]*domain.CompanyPerformanceRule, error) {
	var companyPerformanceRules []CompanyPerformanceRule
	list := make([]*domain.CompanyPerformanceRule, 0)

	if result := cpr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Order("create_time DESC,id DESC").
		Find(&companyPerformanceRules); result.Error != nil {
		return nil, result.Error
	}

	for _, companyPerformanceRule := range companyPerformanceRules {
		list = append(list, companyPerformanceRule.ToDomain())
	}

	return list, nil
}

func (cpr *companyPerformanceRuleRepo) ListAdvertiserIdsAndPerformanceRuleNotNull(ctx context.Context, companyId uint64) ([]*domain.CompanyPerformanceRule, error) {
	var companyPerformanceRules []CompanyPerformanceRule
	list := make([]*domain.CompanyPerformanceRule, 0)

	if result := cpr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("advertiser_ids != ''").
		Where("performance_rules != ''").
		Order("create_time,id DESC").
		Find(&companyPerformanceRules); result.Error != nil {
		return nil, result.Error
	}

	for _, companyPerformanceRule := range companyPerformanceRules {
		list = append(list, companyPerformanceRule.ToDomain())
	}

	return list, nil
}

func (cpr *companyPerformanceRuleRepo) Save(ctx context.Context, in *domain.CompanyPerformanceRule) (*domain.CompanyPerformanceRule, error) {
	companyPerformanceRule := &CompanyPerformanceRule{
		CompanyId:        in.CompanyId,
		PerformanceName:  in.PerformanceName,
		AdvertiserIds:    in.AdvertiserIds,
		PerformanceRules: in.PerformanceRules,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := cpr.data.DB(ctx).Create(companyPerformanceRule); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceRule.ToDomain(), nil
}

func (cpr *companyPerformanceRuleRepo) Update(ctx context.Context, in *domain.CompanyPerformanceRule) (*domain.CompanyPerformanceRule, error) {
	companyPerformanceRule := &CompanyPerformanceRule{
		Id:               in.Id,
		CompanyId:        in.CompanyId,
		PerformanceName:  in.PerformanceName,
		AdvertiserIds:    in.AdvertiserIds,
		PerformanceRules: in.PerformanceRules,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := cpr.data.DB(ctx).Save(companyPerformanceRule); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceRule.ToDomain(), nil
}

func (cpr *companyPerformanceRuleRepo) DeleteByCompanyId(ctx context.Context, companyId uint64) error {
	if result := cpr.data.DB(ctx).Where("company_id = ?", companyId).Delete(&CompanyPerformanceRule{}); result.Error != nil {
		return result.Error
	}

	return nil
}
