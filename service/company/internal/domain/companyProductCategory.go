package domain

import (
	"context"
)

type CompanyProductCategory struct {
	CategoryId   uint64
	ParentId     uint64
	CategoryName string
	Sort         uint64
	ChildList    []*CompanyProductCategory
}

func NewCompanyProductCategory(ctx context.Context, categoryId, parentId uint64, categoryName string) *CompanyProductCategory {
	return &CompanyProductCategory{
		CategoryId:   categoryId,
		ParentId:     parentId,
		CategoryName: categoryName,
	}
}

func (cpc *CompanyProductCategory) SetCategoryId(ctx context.Context, categoryId uint64) {
	cpc.CategoryId = categoryId
}

func (cpc *CompanyProductCategory) SetParentId(ctx context.Context, parentId uint64) {
	cpc.ParentId = parentId
}

func (cpc *CompanyProductCategory) SetCategoryName(ctx context.Context, categoryName string) {
	cpc.CategoryName = categoryName
}

func (cpc *CompanyProductCategory) SetSort(ctx context.Context, sort uint64) {
	cpc.Sort = sort
}
