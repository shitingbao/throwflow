package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/douyin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type JinritemailOrderRepo interface {
	GetStorePreference(context.Context, uint64) (*v1.GetStorePreferenceJinritemaiOrdersReply, error)
	List(context.Context, uint64, uint64, uint64, string, string) (*v1.ListJinritemaiOrdersReply, error)
	Statistics(context.Context, uint64, string, string) (*v1.StatisticsJinritemaiOrdersReply, error)
}

type JinritemailOrderUsecase struct {
	repo JinritemailOrderRepo
	conf *conf.Data
	log  *log.Helper
}

func NewJinritemailOrderUsecase(repo JinritemailOrderRepo, conf *conf.Data, logger log.Logger) *JinritemailOrderUsecase {
	return &JinritemailOrderUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (jouc *JinritemailOrderUsecase) GetStorePreferenceUserOpenDouyins(ctx context.Context, userId uint64) (*v1.GetStorePreferenceJinritemaiOrdersReply, error) {
	user, err := jouc.repo.GetStorePreference(ctx, userId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_STORE_PREFERENCE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return user, nil
}

func (jouc *JinritemailOrderUsecase) ListUserOrders(ctx context.Context, pageNum, pageSize, userId uint64, startDay, endDay string) (*v1.ListJinritemaiOrdersReply, error) {
	if pageSize == 0 {
		pageSize = uint64(jouc.conf.Database.PageSize)
	}

	list, err := jouc.repo.List(ctx, userId, pageNum, pageSize, startDay, endDay)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (jouc *JinritemailOrderUsecase) StatisticsUserOrders(ctx context.Context, userId uint64, startDay, endDay string) (*v1.StatisticsJinritemaiOrdersReply, error) {
	list, err := jouc.repo.Statistics(ctx, userId, startDay, endDay)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_STATISTICS_USER_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}
