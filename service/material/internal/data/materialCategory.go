package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"material/internal/biz"
	"material/internal/domain"
)

// 市场素材分类信息表
type MaterialCategory struct {
	CategoryId   uint64 `gorm:"column:category_id;type:bigint(20) UNSIGNED;not null;default:0;comment:分类ID"`
	ParentId     uint64 `gorm:"column:parent_id;type:bigint(20) UNSIGNED;not null;default:0;comment:父级ID"`
	CategoryName string `gorm:"column:category_name;type:varchar(50);not null;default:'';comment:分类名称"`
	Sort         uint64 `gorm:"column:sort;type:bigint(20) UNSIGNED;not null;default:0;comment:排序数值"`
}

func (MaterialCategory) TableName() string {
	return "material_material_category"
}

type materialCategoryRepo struct {
	data *Data
	log  *log.Helper
}

func (mc *MaterialCategory) ToDomain() *domain.MaterialCategory {
	materialCategory := &domain.MaterialCategory{
		CategoryId:   mc.CategoryId,
		ParentId:     mc.ParentId,
		CategoryName: mc.CategoryName,
		Sort:         mc.Sort,
	}

	return materialCategory
}

func NewMaterialCategoryRepo(data *Data, logger log.Logger) biz.MaterialCategoryRepo {
	return &materialCategoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (mcr *materialCategoryRepo) List(ctx context.Context) ([]*domain.MaterialCategory, error) {
	var materialCategories []MaterialCategory
	list := make([]*domain.MaterialCategory, 0)

	if result := mcr.data.db.WithContext(ctx).Order("parent_id ASC, sort DESC, category_id DESC").
		Find(&materialCategories); result.Error != nil {
		return nil, result.Error
	}

	for _, materialCategory := range materialCategories {
		list = append(list, materialCategory.ToDomain())
	}

	return list, nil
}
