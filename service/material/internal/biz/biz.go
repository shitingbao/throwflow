package biz

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var (
	ProviderSet = wire.NewSet(NewMaterialUsecase, NewMaterialContentUsecase, NewCollectUsecase)

	MaterialValidatorError = errors.BadRequest("MATERIAL_VALIDATOR_ERROR", "参数异常")
	MaterialDataError      = errors.BadRequest("MATERIAL_DATA_ERROR", "操作数据异常")
)
