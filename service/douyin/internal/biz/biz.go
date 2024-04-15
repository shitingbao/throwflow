package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var (
	ProviderSet = wire.NewSet(NewOceanengineConfigUsecase, NewOceanengineAccountUsecase, NewQianchuanAdvertiserUsecase, NewQianchuanCampaignUsecase, NewQianchuanAdUsecase, NewQianchuanReportAwemeUsecase, NewQianchuanReportProductUsecase, NewOceanengineAccountTokenUsecase, NewQianchuanAdvertiserHistoryUsecase, NewOpenDouyinTokenUsecase, NewJinritemaiOrderUsecase, NewJinritemaiStoreUsecase, NewOpenDouyinVideoUsecase, NewOpenDouyinUserInfoUsecase, NewDoukeProductUsecase, NewDoukeOrderUsecase)

	DouyinValidatorError = errors.BadRequest("DOUYIN_VALIDATOR_ERROR", "参数异常")
	DouyinDataError      = errors.BadRequest("DOUYIN_DATA_ERROR", "操作数据异常")
)

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}
