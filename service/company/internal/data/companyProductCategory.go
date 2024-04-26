package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/clause"
)

// 企业爆品库商品分类信息表
type CompanyProductCategory struct {
	CategoryId   uint64 `gorm:"column:category_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:category_id;default:0;comment:分类ID"`
	ParentId     uint64 `gorm:"column:parent_id;type:bigint(20) UNSIGNED;not null;default:0;comment:父级ID"`
	CategoryName string `gorm:"column:category_name;type:varchar(50);not null;default:'';comment:分类名称"`
	Sort         uint64 `gorm:"column:sort;type:bigint(20) UNSIGNED;not null;default:0;comment:排序数值"`
}

func (CompanyProductCategory) TableName() string {
	return "company_product_category"
}

type companyProductCategoryRepo struct {
	data *Data
	log  *log.Helper
}

func (cpc *CompanyProductCategory) ToDomain() *domain.CompanyProductCategory {
	companyProductCategory := &domain.CompanyProductCategory{
		CategoryId:   cpc.CategoryId,
		ParentId:     cpc.ParentId,
		CategoryName: cpc.CategoryName,
		Sort:         cpc.Sort,
	}

	return companyProductCategory
}

func NewCompanyProductCategoryRepo(data *Data, logger log.Logger) biz.CompanyProductCategoryRepo {
	return &companyProductCategoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cpcr *companyProductCategoryRepo) List(ctx context.Context) ([]*domain.CompanyProductCategory, error) {
	var companyProductCategories []CompanyProductCategory
	list := make([]*domain.CompanyProductCategory, 0)

	if result := cpcr.data.db.WithContext(ctx).Order("parent_id ASC, sort DESC").
		Find(&companyProductCategories); result.Error != nil {
		return nil, result.Error
	}

	for _, companyProductCategory := range companyProductCategories {
		list = append(list, companyProductCategory.ToDomain())
	}

	return list, nil
}

func (cpcr *companyProductCategoryRepo) Upsert(ctx context.Context, in *domain.CompanyProductCategory) error {
	if result := cpcr.data.DB(ctx).Table("company_product_category").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "category_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"parent_id", "category_name", "sort"}),
	}).Create(&domain.CompanyProductCategoryGorm{
		CategoryId:   in.CategoryId,
		ParentId:     in.ParentId,
		CategoryName: in.CategoryName,
		Sort:         in.Sort,
	}); result.Error != nil {
		return result.Error
	}

	return nil
}
