package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var (
	ProviderSet = wire.NewSet(NewMenuUsecase, NewIndustryUsecase, NewClueUsecase, NewCompanyUsecase, NewCompanyUserUsecase, NewCompanyPerformanceRuleUsecase, NewCompanyPerformanceRebalanceUsecase, NewCompanyPerformanceUsecase, NewCompanySetUsecase, NewCompanyProductUsecase, NewCompanyMaterialUsecase, NewCompanyOrganizationUsecase, NewAreaUsecase, NewCompanyTaskUsecase, NewCompanyTaskAccountRelationUsecase, NewCompanyTaskDetailUsecase)

	CompanyValidatorError = errors.BadRequest("COMPANY_VALIDATOR_ERROR", "参数异常")
	CompanyDataError      = errors.BadRequest("COMPANY_DATA_ERROR", "操作数据异常")
	CompanyLoginAbnormal  = errors.Unauthorized("COMPANY_lOGIN_ABNORMAL", "登录异常")
)

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}
