package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库团队绩效月表
type CompanyPerformanceMonthly struct {
	Id            uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyId     uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	UserId        uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:用户ID"`
	Username      string    `gorm:"column:username;type:varchar(50);not null;comment:姓名"`
	Job           string    `gorm:"column:job;type:varchar(50);not null;comment:职位"`
	Advertisers   string    `gorm:"column:advertisers;type:text;not null;comment:负责千川账户数"`
	StatCost      float32   `gorm:"column:stat_cost;type:decimal(10, 2) UNSIGNED;not null;comment:月底总消耗(元)"`
	Roi           float32   `gorm:"column:roi;type:decimal(10, 2) UNSIGNED;not null;comment:平均ROI"`
	Cost          float32   `gorm:"column:cost;type:decimal(10, 2) UNSIGNED;not null;comment:绩效提成"`
	RebalanceCost float32   `gorm:"column:rebalance_cost;type:decimal(10, 2);not null;comment:奖罚"`
	TotalCost     float32   `gorm:"column:total_cost;type:decimal(10, 2);not null;comment:合计收入"`
	UpdateDay     uint32    `gorm:"column:update_day;type:int UNSIGNED;not null;comment:操作时间"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyPerformanceMonthly) TableName() string {
	return "company_company_performance_monthly"
}

type companyPerformanceMonthlyRepo struct {
	data *Data
	log  *log.Helper
}

func (cpm *CompanyPerformanceMonthly) ToDomain() *domain.CompanyPerformanceMonthly {
	companyPerformanceMonthly := &domain.CompanyPerformanceMonthly{
		Id:            cpm.Id,
		CompanyId:     cpm.CompanyId,
		UserId:        cpm.UserId,
		Username:      cpm.Username,
		Job:           cpm.Job,
		Advertisers:   cpm.Advertisers,
		StatCost:      cpm.StatCost,
		Roi:           cpm.Roi,
		Cost:          cpm.Cost,
		RebalanceCost: cpm.RebalanceCost,
		TotalCost:     cpm.TotalCost,
		UpdateDay:     cpm.UpdateDay,
		CreateTime:    cpm.CreateTime,
		UpdateTime:    cpm.UpdateTime,
	}

	return companyPerformanceMonthly
}

func NewCompanyPerformanceMonthlyRepo(data *Data, logger log.Logger) biz.CompanyPerformanceMonthlyRepo {
	return &companyPerformanceMonthlyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cpmr *companyPerformanceMonthlyRepo) GetByUserIdAndCompanyIdAndUpdateDay(ctx context.Context, userId, companyId uint64, updateDay uint32) (*domain.CompanyPerformanceMonthly, error) {
	companyPerformanceMonthly := &CompanyPerformanceMonthly{}

	if result := cpmr.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("user_id = ?", userId).Where("update_day = ?", updateDay).First(companyPerformanceMonthly); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceMonthly.ToDomain(), nil
}

func (cpmr *companyPerformanceMonthlyRepo) List(ctx context.Context, pageNum, pageSize int, companyId uint64, updateDay uint32) ([]*domain.CompanyPerformanceMonthly, error) {
	var companyPerformanceMonthlies []CompanyPerformanceMonthly
	list := make([]*domain.CompanyPerformanceMonthly, 0)

	if pageNum > 0 {
		if result := cpmr.data.db.WithContext(ctx).
			Where("company_id = ?", companyId).
			Where("update_day = ?", updateDay).
			Order("total_cost,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&companyPerformanceMonthlies); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := cpmr.data.db.WithContext(ctx).
			Where("company_id = ?", companyId).
			Where("update_day = ?", updateDay).
			Order("total_cost,id DESC").
			Find(&companyPerformanceMonthlies); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, companyPerformanceMonthly := range companyPerformanceMonthlies {
		list = append(list, companyPerformanceMonthly.ToDomain())
	}

	return list, nil
}

func (cpmr *companyPerformanceMonthlyRepo) Count(ctx context.Context, companyId uint64, updateDay uint32) (int64, error) {
	var count int64

	if result := cpmr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("update_day = ?", updateDay).
		Model(&CompanyPerformanceMonthly{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cpmr *companyPerformanceMonthlyRepo) Sum(ctx context.Context, companyId uint64, updateDay uint32) (float32, error) {
	var list []float32
	var sum float32

	if result := cpmr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("update_day = ?", updateDay).
		Model(&CompanyPerformanceMonthly{}).Pluck("stat_cost", &list); result.Error != nil {
		return 0.00, result.Error
	}

	for _, v := range list {
		sum += v
	}

	return sum, nil
}

func (cpmr *companyPerformanceMonthlyRepo) Save(ctx context.Context, in *domain.CompanyPerformanceMonthly) (*domain.CompanyPerformanceMonthly, error) {
	companyPerformanceMonthly := &CompanyPerformanceMonthly{
		CompanyId:     in.CompanyId,
		UserId:        in.UserId,
		Username:      in.Username,
		Job:           in.Job,
		Advertisers:   in.Advertisers,
		StatCost:      in.StatCost,
		Roi:           in.Roi,
		Cost:          in.Cost,
		RebalanceCost: in.RebalanceCost,
		TotalCost:     in.TotalCost,
		UpdateDay:     in.UpdateDay,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if result := cpmr.data.DB(ctx).Create(companyPerformanceMonthly); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceMonthly.ToDomain(), nil
}

func (cpmr *companyPerformanceMonthlyRepo) Update(ctx context.Context, in *domain.CompanyPerformanceMonthly) (*domain.CompanyPerformanceMonthly, error) {
	companyPerformanceMonthly := &CompanyPerformanceMonthly{
		Id:            in.Id,
		CompanyId:     in.CompanyId,
		UserId:        in.UserId,
		Username:      in.Username,
		Job:           in.Job,
		Advertisers:   in.Advertisers,
		StatCost:      in.StatCost,
		Roi:           in.Roi,
		Cost:          in.Cost,
		RebalanceCost: in.RebalanceCost,
		TotalCost:     in.TotalCost,
		UpdateDay:     in.UpdateDay,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if result := cpmr.data.DB(ctx).Save(companyPerformanceMonthly); result.Error != nil {
		return nil, result.Error
	}

	return companyPerformanceMonthly.ToDomain(), nil
}

func (cpmr *companyPerformanceMonthlyRepo) DeleteByCompanyId(ctx context.Context, companyId uint64) error {
	if result := cpmr.data.DB(ctx).Where("company_id = ?", companyId).Delete(&CompanyPerformanceMonthly{}); result.Error != nil {
		return result.Error
	}

	return nil
}
