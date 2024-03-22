package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewOceanengineAccountTokenUsecase, NewOpenDouyinTokenUsecase, NewQianchuanAdvertiserUsecase, NewQianchuanAdUsecase, NewCompanyUsecase, NewCompanyOrganizationUsecase, NewJinritemaiOrderUsecase, NewJinritemaiStoreUsecase, NewOpenDouyinVideoUsecase, NewWeixinUserUsecase, NewWeixinUserOrganizationUsecase, NewWeixinUserCommissionUsecase, NewWeixinUserCouponUsecase, NewWeixinUserBalanceUsecase, NewDoukeTokenUsecase, NewCompanyTaskUsecase)
