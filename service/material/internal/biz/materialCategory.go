package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"material/internal/domain"
)

var (
	MaterialMaterialCategoryNotFound = errors.NotFound("MATERIAL_MATERIAL_CATEGORY_NOT_FOUND", "素材分类不存在")
)

type MaterialCategoryRepo interface {
	List(context.Context) ([]*domain.MaterialCategory, error)
}
