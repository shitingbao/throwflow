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
	DouyinQianchuanReportProductListError = errors.InternalServer("DOUYIN_QIANCHUAN_REPORT_PRODUCT_LIST_ERROR", "短视频带货获取失败")
)

type QianchuanReportProductRepo interface {
	List(context.Context, string, string, string, uint8, uint64, uint64) ([]*domain.QianchuanReportProduct, error)
	Count(context.Context, string, string, string, uint8) (int64, error)
	StatisticsPayOrderCount(context.Context, string, string) (int64, error)
	StatisticsPayOrderAmount(context.Context, string, string) (float64, error)
	StatisticsStatCost(context.Context, string, string) (float64, error)
	SaveIndex(context.Context, string)
	UpsertQianchuanReportProductInfo(context.Context, string, *domain.QianchuanReportProduct) error
	UpsertQianchuanReportProduct(context.Context, string, *domain.QianchuanReportProduct) error
}

type QianchuanReportProductUsecase struct {
	repo QianchuanReportProductRepo
	conf *conf.Data
	log  *log.Helper
}

func NewQianchuanReportProductUsecase(repo QianchuanReportProductRepo, conf *conf.Data, logger log.Logger) *QianchuanReportProductUsecase {
	return &QianchuanReportProductUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (qrpuc *QianchuanReportProductUsecase) ListQianchuanReportProducts(ctx context.Context, pageNum, pageSize uint64, isDistinction uint8, day, keyword, advertiserIds string) (*domain.QianchuanReportProductList, error) {
	list := make([]*domain.QianchuanReportProduct, 0)
	var total int64 = 0
	var lerr error
	var terr error

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		list, lerr = qrpuc.repo.List(ctx, advertiserIds, day, keyword, isDistinction, pageNum, pageSize)

		wg.Done()
	}()

	go func() {
		total, terr = qrpuc.repo.Count(ctx, advertiserIds, day, keyword, isDistinction)

		wg.Done()
	}()

	wg.Wait()

	if lerr != nil || terr != nil {
		return nil, DouyinQianchuanReportProductListError
	}

	for _, l := range list {
		l.SetRoi(ctx)
		l.SetPayOrderAveragePrice(ctx)
		l.SetClickRate(ctx)
		l.SetConvertRate(ctx)
		l.SetAveragePayOrderStatCost(ctx)
	}

	return &domain.QianchuanReportProductList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (qrpuc *QianchuanReportProductUsecase) StatisticsQianchuanReportProducts(ctx context.Context, day, advertiserIds string) (*domain.StatisticsQianchuanReportProducts, error) {
	var totalProduct int64 = 0
	var totalPayOrderCount int64 = 0
	var totalPayOrderAmount float64 = 0.00
	var totalStatCost float64 = 0.00
	var totalRoi float64 = 0.00

	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		totalProduct, _ = qrpuc.repo.Count(ctx, advertiserIds, day, "", 0)

		wg.Done()
	}()

	go func() {
		totalPayOrderCount, _ = qrpuc.repo.StatisticsPayOrderCount(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalPayOrderAmount, _ = qrpuc.repo.StatisticsPayOrderAmount(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalStatCost, _ = qrpuc.repo.StatisticsStatCost(ctx, advertiserIds, day)

		wg.Done()
	}()

	wg.Wait()

	statistics := make([]*domain.StatisticsQianchuanReportProduct, 0)

	statistics = append(statistics, &domain.StatisticsQianchuanReportProduct{
		Key:   "商品数",
		Value: fmt.Sprintf("%d", totalProduct),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanReportProduct{
		Key:   "广告单数",
		Value: fmt.Sprintf("%d", totalPayOrderCount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanReportProduct{
		Key:   "总广告成交",
		Value: fmt.Sprintf("%.2f¥", totalPayOrderAmount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanReportProduct{
		Key:   "广告消耗",
		Value: fmt.Sprintf("%.2f", totalStatCost),
	})

	if totalStatCost > 0 {
		totalRoi = totalPayOrderAmount / totalStatCost
	}

	statistics = append(statistics, &domain.StatisticsQianchuanReportProduct{
		Key:   "广告ROI",
		Value: fmt.Sprintf("%.2f", totalRoi),
	})

	return &domain.StatisticsQianchuanReportProducts{
		Statistics: statistics,
	}, nil
}
