package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库团队绩效日表
type CompanyPerformanceDaily struct {
	Id             uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyId      uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	UserId         uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:用户ID"`
	Advertisers    string    `gorm:"column:advertisers;type:text;not null;comment:负责千川账户数"`
	StatCost       float32   `gorm:"column:stat_cost;type:decimal(10, 2) UNSIGNED;not null;comment:消耗"`
	PayOrderAmount float32   `gorm:"column:pay_order_amount;type:decimal(10, 2) UNSIGNED;not null;comment:成交额"`
	Roi            float32   `gorm:"column:roi;type:decimal(10, 2) UNSIGNED;not null;comment:ROI"`
	Cost           float32   `gorm:"column:cost;type:decimal(10, 2) UNSIGNED;not null;comment:绩效提成"`
	UpdateDay      time.Time `gorm:"column:update_day;type:date;not null;comment:操作时间"`
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyPerformanceDaily) TableName() string {
	return "company_company_performance_daily"
}

type companyPerformanceDailyRepo struct {
	data *Data
	log  *log.Helper
}

func (c *CompanyPerformanceDaily) ToDomain() *domain.CompanyPerformanceDaily {
	companyPerformanceDaily := &domain.CompanyPerformanceDaily{
		Id:             c.Id,
		CompanyId:      c.CompanyId,
		UserId:         c.UserId,
		Advertisers:    c.Advertisers,
		StatCost:       c.StatCost,
		PayOrderAmount: c.PayOrderAmount,
		Roi:            c.Roi,
		Cost:           c.Cost,
		UpdateDay:      c.UpdateDay,
		CreateTime:     c.CreateTime,
		UpdateTime:     c.UpdateTime,
	}

	return companyPerformanceDaily
}

func NewCompanyPerformanceDailyRepo(data *Data, logger log.Logger) biz.CompanyPerformanceDailyRepo {
	return &companyPerformanceDailyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cpdr *companyPerformanceDailyRepo) GetByUserIdAndCompanyIdAndUpdateDay(ctx context.Context, userId, companyId uint64, updateDay string) (*domain.CompanyPerformanceDaily, error) {
	companyPerformanceDaily := &CompanyPerformanceDaily{}

	if result := cpdr.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("user_id = ?", userId).Where("update_day = ?", updateDay).First(companyPerformanceDaily); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceDaily.ToDomain(), nil
}

func (cpdr *companyPerformanceDailyRepo) List(ctx context.Context, userId, companyId uint64, updateDay string) ([]*domain.CompanyPerformanceDaily, error) {
	var companyPerformanceDailies []CompanyPerformanceDaily
	list := make([]*domain.CompanyPerformanceDaily, 0)

	if result := cpdr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("company_id = ?", companyId).
		Where("LEFT(update_day, 7) = ?", updateDay).
		Order("update_day DESC").
		Find(&companyPerformanceDailies); result.Error != nil {
		return nil, result.Error
	}

	for _, companyPerformanceDaily := range companyPerformanceDailies {
		list = append(list, companyPerformanceDaily.ToDomain())
	}

	return list, nil
}

func (cpdr *companyPerformanceDailyRepo) Save(ctx context.Context, in *domain.CompanyPerformanceDaily) (*domain.CompanyPerformanceDaily, error) {
	companyPerformanceDaily := &CompanyPerformanceDaily{
		CompanyId:      in.CompanyId,
		UserId:         in.UserId,
		Advertisers:    in.Advertisers,
		StatCost:       in.StatCost,
		PayOrderAmount: in.PayOrderAmount,
		Roi:            in.Roi,
		Cost:           in.Cost,
		UpdateDay:      in.UpdateDay,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := cpdr.data.DB(ctx).Create(companyPerformanceDaily); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceDaily.ToDomain(), nil
}

func (cpdr *companyPerformanceDailyRepo) Update(ctx context.Context, in *domain.CompanyPerformanceDaily) (*domain.CompanyPerformanceDaily, error) {
	companyPerformanceDaily := &CompanyPerformanceDaily{
		Id:             in.Id,
		CompanyId:      in.CompanyId,
		UserId:         in.UserId,
		Advertisers:    in.Advertisers,
		StatCost:       in.StatCost,
		PayOrderAmount: in.PayOrderAmount,
		Roi:            in.Roi,
		Cost:           in.Cost,
		UpdateDay:      in.UpdateDay,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := cpdr.data.DB(ctx).Save(companyPerformanceDaily); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceDaily.ToDomain(), nil
}

func (cpdr *companyPerformanceDailyRepo) DeleteByCompanyId(ctx context.Context, companyId uint64) error {
	if result := cpdr.data.DB(ctx).Where("company_id = ?", companyId).Delete(&CompanyPerformanceDaily{}); result.Error != nil {
		return result.Error
	}

	return nil
}
