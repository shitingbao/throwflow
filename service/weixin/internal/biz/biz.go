package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var (
	ProviderSet = wire.NewSet(NewUserUsecase, NewUserOpenDouyinUsecase, NewQrCodeUsecase, NewUserAddressUsecase, NewUserSampleOrderUsecase, NewUserOrderUsecase, NewUserOrganizationRelationUsecase, NewUserScanRecordUsecase, NewUserCommissionUsecase, NewUserCouponUsecase, NewUserBalanceUsecase, NewUserContractUsecase, NewUserBankUsecase)

	WeixinValidatorError = errors.BadRequest("WEIXIN_VALIDATOR_ERROR", "参数异常")
)

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}
