package data

import (
	"common/internal/biz"
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// 行政区划代码表
type Area struct {
	AreaName       string `gorm:"column:area_name;type:varchar(50);not null;comment:名称"`
	AreaCode       uint64 `gorm:"column:area_code;type:bigint(20) UNSIGNED;not null;default:0;comment:区划代码"`
	ParentAreaCode uint64 `gorm:"column:parent_area_code;type:bigint(20) UNSIGNED;not null;default:0;comment:父级区划代码"`
}

func (Area) TableName() string {
	return "common_area"
}

type areaRepo struct {
	data *Data
	log  *log.Helper
}

func (a *Area) ToDomain() *domain.Area {
	return &domain.Area{
		AreaName:       a.AreaName,
		AreaCode:       a.AreaCode,
		ParentAreaCode: a.ParentAreaCode,
	}
}

func NewAreaRepo(data *Data, logger log.Logger) biz.AreaRepo {
	return &areaRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ar *areaRepo) GetByAreaCode(ctx context.Context, areaCode uint64) (*domain.Area, error) {
	area := &Area{}

	if result := ar.data.db.WithContext(ctx).Where("area_code = ?", areaCode).First(area); result.Error != nil {
		return nil, result.Error
	}

	return area.ToDomain(), nil
}

func (ar *areaRepo) ListByParentAreaCode(ctx context.Context, parentAreaCode uint64) ([]*domain.Area, error) {
	var areas []Area
	list := make([]*domain.Area, 0)

	if result := ar.data.db.WithContext(ctx).Where("parent_area_code = ?", parentAreaCode).Find(&areas); result.Error != nil {
		return nil, result.Error
	}

	for _, area := range areas {
		list = append(list, area.ToDomain())
	}

	return list, nil
}
