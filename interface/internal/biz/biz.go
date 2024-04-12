package biz

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var (
	ProviderSet = wire.NewSet(NewIndustryUsecase, NewAreaUsecase, NewSmsUsecase, NewTokenUsecase, NewKuaidiCompanyUsecase, NewLoginUsecase, NewCompanyUserUsecase, NewCompanySetUsecase, NewClueUsecase, NewMaterialUsecase, NewPerformanceRuleUsecase, NewPerformanceUsecase, NewPerformanceRebalanceUsecase, NewUpdateLogUsecase, NewUserUsecase, NewProductUsecase, NewCompanyUsecase, NewCompanyMaterialUsecase, NewCompanyOrganizationUsecase, NewUserAddressUsecase, NewUserStoreUsecase, NewUserOpenDouyinUsecase, NewUserScanRecordUsecase, NewUserOrganizationUsecase, NewJinritemailOrderUsecase, NewDoukeOrderUsecase, NewUserSampleOrderUsecase, NewUserCouponUsecase, NewUserBalanceUsecase, NewUserContractUsecase, NewUserBankUsecase, NewCompanyTaskUsecase, NewCourseUsecase)

	InterfaceDataError       = errors.BadRequest("INTERFACE_DATA_ERROR", "操作数据异常")
	InterfaceValidatorError  = errors.BadRequest("INTERFACE_VALIDATOR_ERROR", "参数异常")
	InterfaceLoginTokenError = errors.InternalServer("INTERFACE_LOGIN_TOKEN_ERROR", "token验证失败")
)
