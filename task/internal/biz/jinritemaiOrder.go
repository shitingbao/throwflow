package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
)

type JinritemaiOrderRepo interface {
	Sync90Day(context.Context) (*v1.Sync90DayJinritemaiOrdersReply, error)
}

type JinritemaiOrderUsecase struct {
	repo JinritemaiOrderRepo
	log  *log.Helper
}

func NewJinritemaiOrderUsecase(repo JinritemaiOrderRepo, logger log.Logger) *JinritemaiOrderUsecase {
	return &JinritemaiOrderUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (jouc *JinritemaiOrderUsecase) Sync90DayJinritemaiOrders(ctx context.Context) (*v1.Sync90DayJinritemaiOrdersReply, error) {
	jinritemaiOrder, err := jouc.repo.Sync90Day(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC90DAY_JINRITEMAI_ORDER_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return jinritemaiOrder, nil
}
