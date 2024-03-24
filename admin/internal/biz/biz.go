package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var (
	ProviderSet = wire.NewSet(NewUserUsecase, NewMenuUsecase, NewRoleUsecase, NewTokenUsecase, NewSmsLogUsecase, NewOceanengineConfigUsecase, NewCompanyMenuUsecase, NewClueUsecase, NewIndustryUsecase, NewCompanyUsecase, NewCompanyUserUsecase, NewAreaUsecase, NewUpdateLogUsecase)

	AdminValidatorError = errors.BadRequest("ADMIN_VALIDATOR_ERROR", "参数异常")
	AdminDataError      = errors.BadRequest("ADMIN_DATA_ERROR", "操作数据异常")
)

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}
