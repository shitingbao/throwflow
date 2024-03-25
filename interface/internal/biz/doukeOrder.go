package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/douyin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type DoukeOrderRepo interface {
	Statistics(context.Context, uint64) (*v1.StatisticsDoukeOrdersReply, error)
}

type DoukeOrderUsecase struct {
	repo DoukeOrderRepo
	conf *conf.Data
	log  *log.Helper
}

func NewDoukeOrderUsecase(repo DoukeOrderRepo, conf *conf.Data, logger log.Logger) *DoukeOrderUsecase {
	return &DoukeOrderUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (douc *DoukeOrderUsecase) StatisticsUserDoukeOrders(ctx context.Context, userId uint64) (*v1.StatisticsDoukeOrdersReply, error) {
	list, err := douc.repo.Statistics(ctx, userId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_STATISTICS_USER_DOUKE_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}
