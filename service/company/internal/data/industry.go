package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 行业分类表
type Industry struct {
	Id           uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	IndustryName string    `gorm:"column:industry_name;type:varchar(20);not null;comment:行业名称"`
	Status       uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态：1：启用，0：禁用"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (Industry) TableName() string {
	return "company_industry"
}

type industryRepo struct {
	data *Data
	log  *log.Helper
}

func (i *Industry) ToDomain() *domain.Industry {
	return &domain.Industry{
		Id:           i.Id,
		IndustryName: i.IndustryName,
		Status:       i.Status,
		CreateTime:   i.CreateTime,
		UpdateTime:   i.UpdateTime,
	}
}

func NewIndustryRepo(data *Data, logger log.Logger) biz.IndustryRepo {
	return &industryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ir *industryRepo) GetById(ctx context.Context, id uint64) (*domain.Industry, error) {
	industry := &Industry{}

	if result := ir.data.db.WithContext(ctx).First(industry, id); result.Error != nil {
		return nil, result.Error
	}

	return industry.ToDomain(), nil
}

func (ir *industryRepo) List(ctx context.Context, status uint8) ([]*domain.Industry, error) {
	var industries []Industry
	list := make([]*domain.Industry, 0)

	if result := ir.data.db.WithContext(ctx).
		Where("status = ?", status).
		Order("id ASC").
		Find(&industries); result.Error != nil {
		return nil, result.Error
	}

	for _, industry := range industries {
		list = append(list, industry.ToDomain())
	}

	return list, nil
}
