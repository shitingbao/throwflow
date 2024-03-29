package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
)

var (
	DouyinQianchuanReportAwemeListError = errors.InternalServer("DOUYIN_QIANCHUAN_REPORT_AWEME_LIST_ERROR", "直播带货获取失败")
)

type QianchuanReportAwemeRepo interface {
	List(context.Context, string, string, string, uint8, uint64, uint64) ([]*domain.QianchuanReportAweme, error)
	Count(context.Context, string, string, string, uint8) (int64, error)
	StatisticsPayOrderCount(context.Context, string, string) (int64, error)
	StatisticsPayOrderAmount(context.Context, string, string) (float64, error)
	StatisticsStatCost(context.Context, string, string) (float64, error)
	SaveIndex(context.Context, string)
	UpsertQianchuanReportAwemeInfo(context.Context, string, *domain.QianchuanReportAweme) error
	UpsertQianchuanReportAweme(context.Context, string, *domain.QianchuanReportAweme) error
}

type QianchuanReportAwemeUsecase struct {
	repo QianchuanReportAwemeRepo
	conf *conf.Data
	log  *log.Helper
}

func NewQianchuanReportAwemeUsecase(repo QianchuanReportAwemeRepo, conf *conf.Data, logger log.Logger) *QianchuanReportAwemeUsecase {
	return &QianchuanReportAwemeUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (qrauc *QianchuanReportAwemeUsecase) ListQianchuanReportAwemes(ctx context.Context, pageNum, pageSize uint64, isDistinction uint8, day, keyword, advertiserIds string) (*domain.QianchuanReportAwemeList, error) {
	list := make([]*domain.QianchuanReportAweme, 0)
	var total int64 = 0
	var lerr error
	var terr error

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		list, lerr = qrauc.repo.List(ctx, advertiserIds, day, keyword, isDistinction, pageNum, pageSize)

		wg.Done()
	}()

	go func() {
		total, terr = qrauc.repo.Count(ctx, advertiserIds, day, keyword, isDistinction)

		wg.Done()
	}()

	wg.Wait()

	if lerr != nil || terr != nil {
		return nil, DouyinQianchuanReportAwemeListError
	}

	for _, l := range list {
		l.SetRoi(ctx)
		l.SetPayOrderAveragePrice(ctx)
		l.SetClickRate(ctx)
		l.SetConvertRate(ctx)
		l.SetAveragePayOrderStatCost(ctx)
	}

	return &domain.QianchuanReportAwemeList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (qrauc *QianchuanReportAwemeUsecase) StatisticsQianchuanReportAwemes(ctx context.Context, day, advertiserIds string) (*domain.StatisticsQianchuanReportAwemes, error) {
	var totalAweme int64 = 0
	var totalPayOrderCount int64 = 0
	var totalPayOrderAmount float64 = 0.00
	var totalStatCost float64 = 0.00
	var totalRoi float64 = 0.00

	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		totalAweme, _ = qrauc.repo.Count(ctx, advertiserIds, day, "", 0)

		wg.Done()
	}()

	go func() {
		totalPayOrderCount, _ = qrauc.repo.StatisticsPayOrderCount(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalPayOrderAmount, _ = qrauc.repo.StatisticsPayOrderAmount(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalStatCost, _ = qrauc.repo.StatisticsStatCost(ctx, advertiserIds, day)

		wg.Done()
	}()

	wg.Wait()

	statistics := make([]*domain.StatisticsQianchuanReportAweme, 0)

	statistics = append(statistics, &domain.StatisticsQianchuanReportAweme{
		Key:   "抖音号",
		Value: fmt.Sprintf("%d", totalAweme),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanReportAweme{
		Key:   "广告单数",
		Value: fmt.Sprintf("%d", totalPayOrderCount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanReportAweme{
		Key:   "总广告成交",
		Value: fmt.Sprintf("%.2f¥", totalPayOrderAmount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanReportAweme{
		Key:   "广告消耗",
		Value: fmt.Sprintf("%.2f", totalStatCost),
	})

	if totalStatCost > 0 {
		totalRoi = totalPayOrderAmount / totalStatCost
	}

	statistics = append(statistics, &domain.StatisticsQianchuanReportAweme{
		Key:   "广告ROI",
		Value: fmt.Sprintf("%.2f", totalRoi),
	})

	return &domain.StatisticsQianchuanReportAwemes{
		Statistics: statistics,
	}, nil
}
