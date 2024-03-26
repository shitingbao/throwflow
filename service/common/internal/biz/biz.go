package biz

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var (
	ProviderSet = wire.NewSet(NewSmsUsecase, NewPayUsecase, NewTokenUsecase, NewAreaUsecase, NewUpdateLogUsecase, NewShortUrlUsecase, NewShortCodeUsecase, NewKuaidiCompanyUsecase, NewKuaidiInfoUsecase)

	CommonDataError      = errors.BadRequest("COMMON_DATA_ERROR", "操作数据异常")
	CommonValidatorError = errors.BadRequest("COMMON_VALIDATOR_ERROR", "参数异常")
)
