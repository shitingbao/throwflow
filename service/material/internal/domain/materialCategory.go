package domain

import (
	"context"
)

type MaterialCategory struct {
	CategoryId   uint64
	ParentId     uint64
	CategoryName string
	Sort         uint64
	ChildList    []*MaterialCategory
}

func NewMaterialCategory(ctx context.Context, categoryId, parentId uint64, categoryName string) *MaterialCategory {
	return &MaterialCategory{
		CategoryId:   categoryId,
		ParentId:     parentId,
		CategoryName: categoryName,
	}
}

func (mc *MaterialCategory) SetCategoryId(ctx context.Context, categoryId uint64) {
	mc.CategoryId = categoryId
}

func (mc *MaterialCategory) SetParentId(ctx context.Context, parentId uint64) {
	mc.ParentId = parentId
}

func (mc *MaterialCategory) SetCategoryName(ctx context.Context, categoryName string) {
	mc.CategoryName = categoryName
}

func (mc *MaterialCategory) SetSort(ctx context.Context, sort uint64) {
	mc.Sort = sort
}
