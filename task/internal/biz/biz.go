package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewOceanengineAccountTokenUsecase, NewOpenDouyinTokenUsecase, NewQianchuanAdvertiserUsecase, NewQianchuanAdUsecase, NewCompanyUsecase, NewCompanyOrganizationUsecase, NewJinritemaiOrderUsecase, NewJinritemaiStoreUsecase, NewOpenDouyinVideoUsecase, NewWeixinUserUsecase, NewWeixinUserCommissionUsecase, NewWeixinUserCouponUsecase, NewWeixinUserBalanceUsecase, NewDoukeOrderUsecase, NewCompanyTaskUsecase, NewCompanyProductUsecase)
