package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/tool"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
)

var (
	DouyinQianchuanAdvertiserHistoryListError = errors.InternalServer("DOUYIN_QIANCHUAN_ADVERTISER_HISTORY_LIST_ERROR", "千川广告账户历史获取失败")
)

type QianchuanAdvertiserHistoryRepo interface {
	GetByAdvertiserId(context.Context, uint64) (*domain.QianchuanAdvertiserHistory, error)
	List(context.Context, string, uint32) ([]*domain.QianchuanAdvertiserHistory, error)
	Save(context.Context, *domain.QianchuanAdvertiserHistory) (*domain.QianchuanAdvertiserHistory, error)
}

type QianchuanAdvertiserHistoryUsecase struct {
	repo QianchuanAdvertiserHistoryRepo
	conf *conf.Data
	log  *log.Helper
}

func NewQianchuanAdvertiserHistoryUsecase(repo QianchuanAdvertiserHistoryRepo, conf *conf.Data, logger log.Logger) *QianchuanAdvertiserHistoryUsecase {
	return &QianchuanAdvertiserHistoryUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (qahuc *QianchuanAdvertiserHistoryUsecase) ListQianchuanAdvertiserHistorys(ctx context.Context, day, advertiserIds string) ([]*domain.QianchuanAdvertiserHistory, error) {
	tday, err := tool.StringToTime("2006-01-02", day)

	if err != nil {
		return nil, DouyinQianchuanAdvertiserHistoryListError
	}

	uiday, err := strconv.ParseUint(tool.TimeToString("20060102", tday), 10, 64)

	list, err := qahuc.repo.List(ctx, advertiserIds, uint32(uiday))

	if err != nil {
		return nil, DouyinQianchuanAdvertiserHistoryListError
	}

	return list, nil
}
