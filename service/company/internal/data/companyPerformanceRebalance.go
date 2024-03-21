package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库团队绩效奖罚表
type CompanyPerformanceRebalance struct {
	Id            uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyId     uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	UserId        uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:用户ID"`
	Cost          float32   `gorm:"column:cost;type:decimal(10, 2) UNSIGNED;not null;comment:调整金额"`
	RebalanceType uint8     `gorm:"column:rebalance_type;type:tinyint(3) UNSIGNED;not null;default:1;comment:奖罚类型，1：奖励，2：惩罚"`
	Reason        string    `gorm:"column:reason;type:varchar(250);not null;comment:调整原因"`
	UpdateDay     time.Time `gorm:"column:update_day;type:date;not null;comment:操作时间"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyPerformanceRebalance) TableName() string {
	return "company_company_performance_rebalance"
}

type companyPerformanceRebalanceRepo struct {
	data *Data
	log  *log.Helper
}

func (cpr *CompanyPerformanceRebalance) ToDomain() *domain.CompanyPerformanceRebalance {
	companyPerformanceRebalance := &domain.CompanyPerformanceRebalance{
		Id:            cpr.Id,
		CompanyId:     cpr.CompanyId,
		UserId:        cpr.UserId,
		Cost:          cpr.Cost,
		RebalanceType: cpr.RebalanceType,
		Reason:        cpr.Reason,
		UpdateDay:     cpr.UpdateDay,
		CreateTime:    cpr.CreateTime,
		UpdateTime:    cpr.UpdateTime,
	}

	return companyPerformanceRebalance
}

func NewCompanyPerformanceRebalanceRepo(data *Data, logger log.Logger) biz.CompanyPerformanceRebalanceRepo {
	return &companyPerformanceRebalanceRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cpr *companyPerformanceRebalanceRepo) List(ctx context.Context, userId, companyId uint64, updateDay string) ([]*domain.CompanyPerformanceRebalance, error) {
	var performanceRebalances []CompanyPerformanceRebalance
	list := make([]*domain.CompanyPerformanceRebalance, 0)

	if result := cpr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("company_id = ?", companyId).
		Where("LEFT(update_day, 7) = ?", updateDay).
		Order("update_time DESC").
		Find(&performanceRebalances); result.Error != nil {
		return nil, result.Error
	}

	for _, performanceRebalance := range performanceRebalances {
		list = append(list, performanceRebalance.ToDomain())
	}

	return list, nil
}

func (cpr *companyPerformanceRebalanceRepo) Save(ctx context.Context, in *domain.CompanyPerformanceRebalance) (*domain.CompanyPerformanceRebalance, error) {
	companyPerformanceRebalance := &CompanyPerformanceRebalance{
		CompanyId:     in.CompanyId,
		UserId:        in.UserId,
		Cost:          in.Cost,
		RebalanceType: in.RebalanceType,
		Reason:        in.Reason,
		UpdateDay:     in.UpdateDay,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if result := cpr.data.DB(ctx).Create(companyPerformanceRebalance); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceRebalance.ToDomain(), nil
}

func (cpr *companyPerformanceRebalanceRepo) DeleteByCompanyId(ctx context.Context, companyId uint64) error {
	if result := cpr.data.DB(ctx).Where("company_id = ?", companyId).Delete(&CompanyPerformanceRebalance{}); result.Error != nil {
		return result.Error
	}

	return nil
}
