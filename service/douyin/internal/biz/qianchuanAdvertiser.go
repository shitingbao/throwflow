package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/event/event"
	"douyin/internal/pkg/event/kafka"
	"douyin/internal/pkg/oceanengine/account"
	"douyin/internal/pkg/oceanengine/ad"
	"douyin/internal/pkg/oceanengine/campaign"
	"douyin/internal/pkg/oceanengine/finance"
	"douyin/internal/pkg/oceanengine/product"
	"douyin/internal/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/time/rate"
	"strconv"
	"sync"
	"time"
)

var (
	DouyinQianchuanAdvertiserNotFound    = errors.NotFound("DOUYIN_QIANCHUAN_ADVERTISER_NOT_FOUND", "千川广告账户不存在")
	DouyinQianchuanAdvertiserListError   = errors.InternalServer("DOUYIN_QIANCHUAN_ADVERTISER_LIST_ERROR", "千川广告账户获取失败")
	DouyinQianchuanAdvertiserUpdateError = errors.InternalServer("DOUYIN_QIANCHUAN_ADVERTISER_UPDATE_ERROR", "千川广告账户更新失败")
	DouyinQianchuanAdvertiserSyncError   = errors.InternalServer("DOUYIN_QIANCHUAN_ADVERTISER_SYNC_ERROR", "同步千川广告账户列表失败")
	DouyinQianchuanAdvertiserLimitError  = errors.InternalServer("DOUYIN_QIANCHUAN_ADVERTISER_LIMIT_ERROR", "企业账户总数超出限制")
)

type QianchuanAdvertiserRepo interface {
	GetById(context.Context, uint64, uint64) (*domain.QianchuanAdvertiser, error)
	GetByCompanyIdAndAdvertiserId(context.Context, uint64, uint64) (*domain.QianchuanAdvertiser, error)
	Get(context.Context, uint64) (*domain.QianchuanAdvertiser, error)
	ListByAdvertiserId(context.Context, uint64) ([]*domain.QianchuanAdvertiser, error)
	ListByAppId(context.Context, string, uint8) ([]*domain.QianchuanAdvertiser, error)
	List(context.Context, int, int, uint64, string, string, string) ([]*domain.QianchuanAdvertiser, error)
	Count(context.Context, uint64, string, string, string) (int64, error)
	Statistics(context.Context, uint64, uint8) (int64, error)
	Save(context.Context, *domain.QianchuanAdvertiser) (*domain.QianchuanAdvertiser, error)
	Update(context.Context, *domain.QianchuanAdvertiser) (*domain.QianchuanAdvertiser, error)
	DeleteByCompanyIdAndAccountId(context.Context, uint64, uint64) error

	Send(context.Context, event.Event) error
}

type QianchuanAdvertiserUsecase struct {
	repo     QianchuanAdvertiserRepo
	qainrepo QianchuanAdvertiserInfoRepo
	qasrepo  QianchuanAdvertiserStatusRepo
	ocrepo   OceanengineConfigRepo
	oatrepo  OceanengineAccountTokenRepo
	crepo    CompanyRepo
	curepo   CompanyUserRepo
	tlrepo   TaskLogRepo
	qcrepo   QianchuanCampaignRepo
	qprepo   QianchuanProductRepo
	qarepo   QianchuanAwemeRepo
	qwrepo   QianchuanWalletRepo
	qadrepo  QianchuanAdRepo
	qairepo  QianchuanAdInfoRepo
	qraurepo QianchuanReportAdRepo
	qrarrepo QianchuanReportAdRealtimeRepo
	qrprepo  QianchuanReportProductRepo
	qrawrepo QianchuanReportAwemeRepo
	lrrepo   LianshanRealtimeRepo
	oalrepo  OceanengineApiLogRepo
	tm       Transaction
	conf     *conf.Data
	econf    *conf.Event
	oconf    *conf.Oceanengine
	log      *log.Helper
}

func NewQianchuanAdvertiserUsecase(repo QianchuanAdvertiserRepo, qainrepo QianchuanAdvertiserInfoRepo, qasrepo QianchuanAdvertiserStatusRepo, ocrepo OceanengineConfigRepo, oatrepo OceanengineAccountTokenRepo, crepo CompanyRepo, curepo CompanyUserRepo, tlrepo TaskLogRepo, qcrepo QianchuanCampaignRepo, qprepo QianchuanProductRepo, qarepo QianchuanAwemeRepo, qwrepo QianchuanWalletRepo, qadrepo QianchuanAdRepo, qairepo QianchuanAdInfoRepo, qraurepo QianchuanReportAdRepo, qrarrepo QianchuanReportAdRealtimeRepo, qrprepo QianchuanReportProductRepo, qrawrepo QianchuanReportAwemeRepo, lrrepo LianshanRealtimeRepo, oalrepo OceanengineApiLogRepo, tm Transaction, conf *conf.Data, econf *conf.Event, oconf *conf.Oceanengine, logger log.Logger) *QianchuanAdvertiserUsecase {
	return &QianchuanAdvertiserUsecase{repo: repo, qainrepo: qainrepo, qasrepo: qasrepo, ocrepo: ocrepo, oatrepo: oatrepo, crepo: crepo, curepo: curepo, tlrepo: tlrepo, qcrepo: qcrepo, qprepo: qprepo, qarepo: qarepo, qwrepo: qwrepo, qadrepo: qadrepo, qairepo: qairepo, qraurepo: qraurepo, qrarrepo: qrarrepo, qrprepo: qrprepo, qrawrepo: qrawrepo, lrrepo: lrrepo, oalrepo: oalrepo, tm: tm, conf: conf, econf: econf, oconf: oconf, log: log.NewHelper(logger)}
}

func (qauc *QianchuanAdvertiserUsecase) GetQianchuanAdvertiserByCompanyIdAndAdvertiserIds(ctx context.Context, companyId, advertiserId uint64) (*domain.QianchuanAdvertiser, error) {
	qianchuanAdvertiser, err := qauc.repo.GetByCompanyIdAndAdvertiserId(ctx, companyId, advertiserId)

	if err != nil {
		return nil, DouyinQianchuanAdvertiserNotFound
	}

	return qianchuanAdvertiser, nil
}

func (qauc *QianchuanAdvertiserUsecase) ListQianchuanAdvertisers(ctx context.Context, pageNum, pageSize, companyId uint64, keyword, advertiserIds, status string) (*domain.QianchuanAdvertiserList, error) {
	qianchuanAdvertisers, err := qauc.repo.List(ctx, int(pageNum), int(pageSize), companyId, keyword, advertiserIds, status)

	if err != nil {
		return nil, DouyinDataError
	}

	list := make([]*domain.QianchuanAdvertiser, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		lQianchuanAdvertiser := &domain.QianchuanAdvertiser{
			Id:             qianchuanAdvertiser.Id,
			CompanyId:      qianchuanAdvertiser.CompanyId,
			AccountId:      qianchuanAdvertiser.AccountId,
			AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
			AdvertiserName: qianchuanAdvertiser.AdvertiserName,
			CompanyName:    qianchuanAdvertiser.CompanyName,
			Status:         qianchuanAdvertiser.Status,
			CreateTime:     qianchuanAdvertiser.CreateTime,
			UpdateTime:     qianchuanAdvertiser.UpdateTime,
		}

		list = append(list, lQianchuanAdvertiser)
	}

	total, err := qauc.repo.Count(ctx, companyId, keyword, advertiserIds, status)

	if err != nil {
		return nil, DouyinDataError
	}

	return &domain.QianchuanAdvertiserList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (qauc *QianchuanAdvertiserUsecase) ListQianchuanAdvertiserByDays(ctx context.Context, companyId uint64, day string) ([]*domain.QianchuanAdvertiserStatus, error) {
	tday, err := tool.StringToTime("2006-01-02", day)

	if err != nil {
		return nil, DouyinQianchuanAdvertiserListError
	}

	uiday, err := strconv.ParseUint(tool.TimeToString("20060102", tday), 10, 64)

	if err != nil {
		return nil, DouyinQianchuanAdvertiserListError
	}

	list, err := qauc.qasrepo.List(ctx, companyId, uint32(uiday))

	if err != nil {
		return nil, DouyinQianchuanAdvertiserListError
	}

	return list, nil
}

func (qauc *QianchuanAdvertiserUsecase) ListExternalQianchuanAdvertisers(ctx context.Context, pageNum, pageSize uint64, startDay, endDay, advertiserIds string) (*domain.ExternalQianchuanAdvertiserList, error) {
	list := make([]*domain.ExternalQianchuanAdvertiser, 0)

	qianchuanAdvertiserInfos, err := qauc.qainrepo.List(ctx, advertiserIds, startDay, endDay, pageNum, pageSize)

	if err != nil {
		return nil, DouyinQianchuanAdvertiserListError
	}

	total, err := qauc.qainrepo.Count(ctx, advertiserIds, startDay, endDay)

	if err != nil {
		return nil, DouyinQianchuanAdvertiserListError
	}

	var wg sync.WaitGroup

	var qianchuanAdvertiserInfoList sync.Map

	for _, lqianchuanAdvertiserInfo := range qianchuanAdvertiserInfos {
		wg.Add(1)

		go func(advertiserId uint64) {
			defer wg.Done()

			qianchuanAdvertiserInfo, _ := qauc.qainrepo.GetByDay(ctx, advertiserId, startDay, endDay)

			qianchuanAdvertiserInfoList.Store(advertiserId, qianchuanAdvertiserInfo)
		}(lqianchuanAdvertiserInfo.Id)
	}

	wg.Wait()

	for _, lqianchuanAdvertiserInfo := range qianchuanAdvertiserInfos {
		qianchuanAdvertiserInfoList.Range(func(k, v interface{}) bool {
			if lqianchuanAdvertiserInfo.Id == k {
				var qianchuanAdvertiserInfo *domain.QianchuanAdvertiserInfo
				qianchuanAdvertiserInfo = v.(*domain.QianchuanAdvertiserInfo)

				qianchuanAdvertiserInfo.StatCost = lqianchuanAdvertiserInfo.StatCost
				qianchuanAdvertiserInfo.ShowCnt = lqianchuanAdvertiserInfo.ShowCnt
				qianchuanAdvertiserInfo.ClickCnt = lqianchuanAdvertiserInfo.ClickCnt
				qianchuanAdvertiserInfo.PayOrderCount = lqianchuanAdvertiserInfo.PayOrderCount
				qianchuanAdvertiserInfo.CreateOrderAmount = lqianchuanAdvertiserInfo.CreateOrderAmount
				qianchuanAdvertiserInfo.CreateOrderCount = lqianchuanAdvertiserInfo.CreateOrderCount
				qianchuanAdvertiserInfo.PayOrderAmount = lqianchuanAdvertiserInfo.PayOrderAmount
				qianchuanAdvertiserInfo.DyFollow = lqianchuanAdvertiserInfo.DyFollow
				qianchuanAdvertiserInfo.ConvertCnt = lqianchuanAdvertiserInfo.ConvertCnt

				qianchuanAdvertiserInfo.SetRoi(ctx)
				qianchuanAdvertiserInfo.SetClickRate(ctx)
				qianchuanAdvertiserInfo.SetCpmPlatform(ctx)
				qianchuanAdvertiserInfo.SetPayConvertRate(ctx)
				qianchuanAdvertiserInfo.SetConvertCost(ctx)
				qianchuanAdvertiserInfo.SetConvertRate(ctx)
				qianchuanAdvertiserInfo.SetAveragePayOrderStatCost(ctx)
				qianchuanAdvertiserInfo.SetPayOrderAveragePrice(ctx)

				list = append(list, &domain.ExternalQianchuanAdvertiser{
					AdvertiserId:            qianchuanAdvertiserInfo.Id,
					AdvertiserName:          qianchuanAdvertiserInfo.Name,
					StatCost:                qianchuanAdvertiserInfo.StatCost,
					Roi:                     qianchuanAdvertiserInfo.Roi,
					GeneralTotalBalance:     qianchuanAdvertiserInfo.GeneralTotalBalance,
					Campaigns:               qianchuanAdvertiserInfo.Campaigns,
					PayOrderCount:           qianchuanAdvertiserInfo.PayOrderCount,
					PayOrderAmount:          qianchuanAdvertiserInfo.PayOrderAmount,
					ClickCnt:                qianchuanAdvertiserInfo.ClickCnt,
					ClickRate:               qianchuanAdvertiserInfo.ClickRate,
					PayConvertRate:          qianchuanAdvertiserInfo.PayConvertRate,
					AveragePayOrderStatCost: qianchuanAdvertiserInfo.AveragePayOrderStatCost,
					PayOrderAveragePrice:    qianchuanAdvertiserInfo.PayOrderAveragePrice,
					DyFollow:                qianchuanAdvertiserInfo.DyFollow,
					ShowCnt:                 qianchuanAdvertiserInfo.ShowCnt,
				})
			}

			return true
		})
	}

	return &domain.ExternalQianchuanAdvertiserList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}, nil
}

func (qauc *QianchuanAdvertiserUsecase) StatisticsQianchuanAdvertisers(ctx context.Context, companyId uint64) (*domain.StatisticsQianchuanAdvertisers, error) {
	selects := domain.NewSelectQianchuanAdvertisers()

	statistics := make([]*domain.StatisticsQianchuanAdvertiser, 0)

	for _, status := range selects.Status {
		iStatus, _ := strconv.Atoi(status.Key)

		count, _ := qauc.repo.Statistics(ctx, companyId, uint8(iStatus))

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   status.Value,
			Value: strconv.FormatInt(count, 10),
		})
	}

	return &domain.StatisticsQianchuanAdvertisers{
		Statistics: statistics,
	}, nil
}

func (qauc *QianchuanAdvertiserUsecase) StatisticsDashboardQianchuanAdvertisers(ctx context.Context, companyId uint64, day, advertiserIds string) (*domain.StatisticsQianchuanAdvertisers, error) {
	var totalQianchuanAdvertiserCount int64 = 0
	var totalQianchuanAdvertiserByStatCostCount int64 = 0
	var totalProductPayOrderCount int64 = 0
	var totalAwemePayOrderCount int64 = 0
	var totalProductPayOrderAmount float64 = 0.00
	var totalAwemePayOrderAmount float64 = 0.00
	var totalProductStatCost float64 = 0.00
	var totalAwemeStatCost float64 = 0.00
	var totalRoi float64 = 0.00
	var totalProductRoi float64 = 0.00
	var totalAwemeRoi float64 = 0.00

	var wg sync.WaitGroup

	wg.Add(8)

	go func() {
		totalQianchuanAdvertiserCount, _ = qauc.repo.Count(ctx, companyId, "", advertiserIds, "1")

		wg.Done()
	}()

	go func() {
		totalQianchuanAdvertiserByStatCostCount, _ = qauc.qrarrepo.StatisticsAdvertisers(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalAwemePayOrderCount, _ = qauc.qrawrepo.StatisticsPayOrderCount(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalProductPayOrderCount, _ = qauc.qrprepo.StatisticsPayOrderCount(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalProductPayOrderAmount, _ = qauc.qrprepo.StatisticsPayOrderAmount(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalAwemePayOrderAmount, _ = qauc.qrawrepo.StatisticsPayOrderAmount(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalProductStatCost, _ = qauc.qrprepo.StatisticsStatCost(ctx, advertiserIds, day)

		wg.Done()
	}()

	go func() {
		totalAwemeStatCost, _ = qauc.qrawrepo.StatisticsStatCost(ctx, advertiserIds, day)

		wg.Done()
	}()

	wg.Wait()

	statistics := make([]*domain.StatisticsQianchuanAdvertiser, 0)

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "已授权账户",
		Value: fmt.Sprintf("%d", totalQianchuanAdvertiserCount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "有消耗账户",
		Value: fmt.Sprintf("%d", totalQianchuanAdvertiserByStatCostCount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "总单数",
		Value: fmt.Sprintf("%d", totalProductPayOrderCount+totalAwemePayOrderCount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "短视频单数",
		Value: fmt.Sprintf("%d", totalProductPayOrderCount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "直播单数",
		Value: fmt.Sprintf("%d", totalAwemePayOrderCount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "总成交",
		Value: fmt.Sprintf("%.2f¥", totalProductPayOrderAmount+totalAwemePayOrderAmount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "短视频成交",
		Value: fmt.Sprintf("%.2f¥", totalProductPayOrderAmount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "直播成交",
		Value: fmt.Sprintf("%.2f¥", totalAwemePayOrderAmount),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "总消耗",
		Value: fmt.Sprintf("%.2f¥", totalProductStatCost+totalAwemeStatCost),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "短视频消耗",
		Value: fmt.Sprintf("%.2f¥", totalProductStatCost),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "直播消耗",
		Value: fmt.Sprintf("%.2f¥", totalAwemeStatCost),
	})

	if totalProductStatCost > 0 {
		totalProductRoi = totalProductPayOrderAmount / totalProductStatCost
	}

	if totalAwemeStatCost > 0 {
		totalAwemeRoi = totalAwemePayOrderAmount / totalAwemeStatCost
	}

	if (totalProductStatCost + totalAwemeStatCost) > 0 {
		totalRoi = (totalProductPayOrderAmount + totalAwemePayOrderAmount) / (totalProductStatCost + totalAwemeStatCost)
	}

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "总广告ROI",
		Value: fmt.Sprintf("%.2f", totalRoi),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "短视频ROI",
		Value: fmt.Sprintf("%.2f", totalProductRoi),
	})

	statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
		Key:   "直播ROI",
		Value: fmt.Sprintf("%.2f", totalAwemeRoi),
	})

	return &domain.StatisticsQianchuanAdvertisers{
		Statistics: statistics,
	}, nil
}

func (qauc *QianchuanAdvertiserUsecase) StatisticsExternalQianchuanAdvertisers(ctx context.Context, startDay, endDay, advertiserIds string) (*domain.StatisticsQianchuanAdvertisers, error) {
	statistics := make([]*domain.StatisticsQianchuanAdvertiser, 0)

	if statisticsAdvertiser, err := qauc.qainrepo.Statistics(ctx, advertiserIds, startDay, endDay); err == nil {
		statisticsAdvertiser.SetRoi(ctx)
		statisticsAdvertiser.SetClickRate(ctx)
		statisticsAdvertiser.SetCpmPlatform(ctx)
		statisticsAdvertiser.SetPayConvertRate(ctx)
		statisticsAdvertiser.SetConvertCost(ctx)
		statisticsAdvertiser.SetConvertRate(ctx)
		statisticsAdvertiser.SetAveragePayOrderStatCost(ctx)
		statisticsAdvertiser.SetPayOrderAveragePrice(ctx)

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalGeneralTotalBalance",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.GeneralTotalBalance, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalCampaigns",
			Value: strconv.FormatUint(statisticsAdvertiser.Campaigns, 10),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalStatCost",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.StatCost, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalPayOrderCount",
			Value: strconv.FormatInt(statisticsAdvertiser.PayOrderCount, 10),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalPayOrderAmount",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.PayOrderAmount, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalCreateOrderCount",
			Value: strconv.FormatInt(statisticsAdvertiser.CreateOrderCount, 10),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalCreateOrderAmount",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.CreateOrderAmount, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalClickCnt",
			Value: strconv.FormatInt(statisticsAdvertiser.ClickCnt, 10),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalShowCnt",
			Value: strconv.FormatInt(statisticsAdvertiser.ShowCnt, 10),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalConvertCnt",
			Value: strconv.FormatInt(statisticsAdvertiser.ConvertCnt, 10),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalDyFollow",
			Value: strconv.FormatInt(statisticsAdvertiser.DyFollow, 10),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalRoi",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.Roi, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalClickRate",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.ClickRate, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalCpmPlatform",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.CpmPlatform, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalPayConvertRate",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.PayConvertRate, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalConvertCost",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.ConvertCost, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalConvertRate",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.ConvertRate, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalAveragePayOrderStatCost",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.AveragePayOrderStatCost, 2), 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalPayOrderAveragePrice",
			Value: strconv.FormatFloat(tool.Decimal(statisticsAdvertiser.PayOrderAveragePrice, 2), 'f', 2, 64),
		})
	} else {
		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalGeneralTotalBalance",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalCampaigns",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalStatCost",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalPayOrderCount",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalPayOrderAmount",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalCreateOrderCount",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalCreateOrderAmount",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalClickCnt",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalShowCnt",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalConvertCnt",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalDyFollow",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalRoi",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalClickRate",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalCpmPlatform",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalPayConvertRate",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalConvertCost",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalConvertRate",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalAveragePayOrderStatCost",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsQianchuanAdvertiser{
			Key:   "TotalPayOrderAveragePrice",
			Value: "0.00",
		})
	}

	return &domain.StatisticsQianchuanAdvertisers{
		Statistics: statistics,
	}, nil
}

func (qauc *QianchuanAdvertiserUsecase) UpdateStatusQianchuanAdvertisers(ctx context.Context, companyId, advertiserId uint64, status uint8) (*domain.QianchuanAdvertiser, error) {
	inQianchuanAdvertiser, err := qauc.repo.GetById(ctx, companyId, advertiserId)

	if err != nil {
		return nil, DouyinQianchuanAdvertiserNotFound
	}

	if status == 1 && inQianchuanAdvertiser.Status == 0 {
		if err := qauc.verifyQianchuanAdvertiserLimit(ctx, companyId); err != nil {
			return nil, err
		}
	}

	oldStatus := inQianchuanAdvertiser.Status

	inQianchuanAdvertiser.SetStatus(ctx, status)

	var qianchuanAdvertiser *domain.QianchuanAdvertiser

	day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

	err = qauc.tm.InTx(ctx, func(ctx context.Context) error {
		qianchuanAdvertiser, err = qauc.repo.Update(ctx, inQianchuanAdvertiser)

		if err != nil {
			return err
		}

		if qianchuanAdvertiser.Status != oldStatus {
			if qianchuanAdvertiserStatus, err := qauc.qasrepo.GetByCompanyIdAndAdvertiserIdAndDay(ctx, companyId, advertiserId, uint32(day)); err == nil {
				if qianchuanAdvertiserStatus.Day == uint32(day) {
					count, err := qauc.qasrepo.Count(ctx, companyId, advertiserId, uint32(day))

					if err != nil {
						return err
					}

					if count == 1 {
						inQianchuanAdvertiserStatus := qianchuanAdvertiserStatus
						inQianchuanAdvertiserStatus.SetStatus(ctx, qianchuanAdvertiser.Status)
						inQianchuanAdvertiserStatus.SetUpdateTime(ctx)

						if _, err := qauc.qasrepo.Update(ctx, inQianchuanAdvertiserStatus); err != nil {
							return err
						}
					} else {
						if err := qauc.qasrepo.DeleteByCompanyIdAndAdvertiserIdAndDay(ctx, companyId, advertiserId, uint32(day)); err != nil {
							return err
						}
					}
				} else {
					inQianchuanAdvertiserStatus := domain.NewQianchuanAdvertiserStatus(ctx, companyId, advertiserId, uint32(day), status)
					inQianchuanAdvertiserStatus.SetCreateTime(ctx)
					inQianchuanAdvertiserStatus.SetUpdateTime(ctx)

					if _, err := qauc.qasrepo.Save(ctx, inQianchuanAdvertiserStatus); err != nil {
						return err
					}
				}
			} else {
				return err
			}
		}

		if qianchuanAdvertiser.Status == 1 {
			if qianchuanAdvertisers, err := qauc.repo.ListByAdvertiserId(ctx, qianchuanAdvertiser.AdvertiserId); err == nil {
				for _, iqianchuanAdvertiser := range qianchuanAdvertisers {
					if iqianchuanAdvertiser.CompanyId != qianchuanAdvertiser.CompanyId {
						oldiStatus := iqianchuanAdvertiser.Status

						iqianchuanAdvertiser.SetStatus(ctx, 0)
						iqianchuanAdvertiser.SetUpdateTime(ctx)

						if _, err := qauc.repo.Update(ctx, iqianchuanAdvertiser); err != nil {
							return err
						}

						if oldiStatus == 1 {
							if iqianchuanAdvertiserStatus, err := qauc.qasrepo.GetByCompanyIdAndAdvertiserIdAndDay(ctx, iqianchuanAdvertiser.CompanyId, iqianchuanAdvertiser.AdvertiserId, uint32(day)); err == nil {
								if iqianchuanAdvertiserStatus.Day == uint32(day) {
									count, err := qauc.qasrepo.Count(ctx, iqianchuanAdvertiser.CompanyId, iqianchuanAdvertiser.AdvertiserId, uint32(day))

									if err != nil {
										return err
									}

									if count == 1 {
										inQianchuanAdvertiserStatus := iqianchuanAdvertiserStatus
										inQianchuanAdvertiserStatus.SetStatus(ctx, 0)
										inQianchuanAdvertiserStatus.SetUpdateTime(ctx)

										if _, err := qauc.qasrepo.Update(ctx, inQianchuanAdvertiserStatus); err != nil {
											return err
										}
									} else {
										if err := qauc.qasrepo.DeleteByCompanyIdAndAdvertiserIdAndDay(ctx, iqianchuanAdvertiser.CompanyId, iqianchuanAdvertiser.AdvertiserId, uint32(day)); err != nil {
											return err
										}
									}
								} else {
									inQianchuanAdvertiserStatus := domain.NewQianchuanAdvertiserStatus(ctx, iqianchuanAdvertiser.CompanyId, iqianchuanAdvertiser.AdvertiserId, uint32(day), 0)
									inQianchuanAdvertiserStatus.SetCreateTime(ctx)
									inQianchuanAdvertiserStatus.SetUpdateTime(ctx)

									if _, err := qauc.qasrepo.Save(ctx, inQianchuanAdvertiserStatus); err != nil {
										return err
									}
								}
							} else {
								return err
							}
						}
					}
				}
			}
		}

		if qianchuanAdvertiser.Status == 0 {
			if _, err := qauc.curepo.DeleteByCompanyIdAndAdvertiserId(ctx, companyId, advertiserId); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, DouyinQianchuanAdvertiserUpdateError
	}

	return qianchuanAdvertiser, nil
}

func (qauc *QianchuanAdvertiserUsecase) UpdateStatusQianchuanAdvertisersByCompanyId(ctx context.Context, companyId uint64, status uint8) error {
	day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

	err := qauc.tm.InTx(ctx, func(ctx context.Context) error {
		if qianchuanAdvertisers, err := qauc.repo.List(ctx, 0, 0, companyId, "", "", "1"); err != nil {
			return err
		} else {
			for _, inQianchuanAdvertiser := range qianchuanAdvertisers {
				inQianchuanAdvertiser.SetStatus(ctx, status)
				inQianchuanAdvertiser.SetUpdateTime(ctx)

				if _, err := qauc.repo.Update(ctx, inQianchuanAdvertiser); err != nil {
					return err
				}

				if iqianchuanAdvertiserStatus, err := qauc.qasrepo.GetByCompanyIdAndAdvertiserIdAndDay(ctx, inQianchuanAdvertiser.CompanyId, inQianchuanAdvertiser.AdvertiserId, uint32(day)); err == nil {
					if iqianchuanAdvertiserStatus.Day == uint32(day) {
						count, err := qauc.qasrepo.Count(ctx, inQianchuanAdvertiser.CompanyId, inQianchuanAdvertiser.AdvertiserId, uint32(day))

						if err != nil {
							return err
						}

						if count == 1 {
							inQianchuanAdvertiserStatus := iqianchuanAdvertiserStatus
							inQianchuanAdvertiserStatus.SetStatus(ctx, 0)
							inQianchuanAdvertiserStatus.SetUpdateTime(ctx)

							if _, err := qauc.qasrepo.Update(ctx, inQianchuanAdvertiserStatus); err != nil {
								return err
							}
						} else {
							if err := qauc.qasrepo.DeleteByCompanyIdAndAdvertiserIdAndDay(ctx, inQianchuanAdvertiser.CompanyId, inQianchuanAdvertiser.AdvertiserId, uint32(day)); err != nil {
								return err
							}
						}
					} else {
						inQianchuanAdvertiserStatus := domain.NewQianchuanAdvertiserStatus(ctx, inQianchuanAdvertiser.CompanyId, inQianchuanAdvertiser.AdvertiserId, uint32(day), 0)
						inQianchuanAdvertiserStatus.SetCreateTime(ctx)
						inQianchuanAdvertiserStatus.SetUpdateTime(ctx)

						if _, err := qauc.qasrepo.Save(ctx, inQianchuanAdvertiserStatus); err != nil {
							return err
						}
					}
				} else {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return DouyinQianchuanAdvertiserUpdateError
	}

	return nil
}

func (qauc *QianchuanAdvertiserUsecase) SyncQianchuanDatas(ctx context.Context) error {
	day := tool.TimeToString("2006-01-02", time.Now())

	list, err := qauc.ocrepo.List(ctx, 1, 0, 0)

	if err != nil {
		return DouyinOceanengineConfigNotFound
	}

	var wg sync.WaitGroup

	for _, oceanengineConfig := range list {
		wg.Add(1)

		go qauc.syncQianchuanData(ctx, &wg, oceanengineConfig, day)
	}

	wg.Wait()

	wg.Add(1)

	go qauc.syncQianchuanAwemwAndProducts(ctx, &wg, day)

	wg.Wait()

	return nil
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanData(ctx context.Context, wg *sync.WaitGroup, oceanengineConfig *domain.OceanengineConfig, day string) error {
	defer wg.Done()

	var sqdwg sync.WaitGroup
	limiter := rate.NewLimiter(0, int(oceanengineConfig.Concurrents))
	limiter.SetLimit(rate.Limit(oceanengineConfig.Concurrents))

	accessTokens, err := qauc.oatrepo.ListByAppId(ctx, oceanengineConfig.AppId)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[syncQianchuanDataError] AppId=%s, Description=%s", oceanengineConfig.AppId, "获取巨量引擎配置文件失败"))
		inTaskLog.SetCreateTime(ctx)

		qauc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	qianchuanAdvertisers, err := qauc.repo.ListByAppId(ctx, oceanengineConfig.AppId, 1)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[syncQianchuanDataError] AppId=%s, Description=%s", oceanengineConfig.AppId, "获取千川账户数据失败"))
		inTaskLog.SetCreateTime(ctx)

		qauc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	sqdwg.Add(5)

	go qauc.syncQianchuanCampaigns(ctx, &sqdwg, limiter, day, qianchuanAdvertisers, accessTokens)
	go qauc.syncQianchuanAds(ctx, &sqdwg, limiter, day, qianchuanAdvertisers, accessTokens)
	go qauc.syncQianchuanProducts(ctx, &sqdwg, limiter, day, qianchuanAdvertisers, accessTokens)
	go qauc.syncQianchuanAwemes(ctx, &sqdwg, limiter, day, qianchuanAdvertisers, accessTokens)
	go qauc.syncQianchuanWallets(ctx, &sqdwg, limiter, day, qianchuanAdvertisers, accessTokens)

	sqdwg.Wait()

	return nil
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanCampaigns(ctx context.Context, sqdwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertisers []*domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken) {
	defer sqdwg.Done()

	qauc.qcrepo.SaveIndex(ctx, day)

	var sqcwg sync.WaitGroup

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		sqcwg.Add(2)

		go qauc.getQianchuanCampaign(ctx, &sqcwg, limiter, day, qianchuanAdvertiser, accessTokens, "VIDEO_PROM_GOODS", "FEED")
		go qauc.getQianchuanCampaign(ctx, &sqcwg, limiter, day, qianchuanAdvertiser, accessTokens, "LIVE_PROM_GOODS", "FEED")
	}

	sqcwg.Wait()
}

func (qauc *QianchuanAdvertiserUsecase) getQianchuanCampaign(ctx context.Context, sqcwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken, marketingGoal, marketingScene string) {
	defer sqcwg.Done()

	for _, accessToken := range accessTokens {
		if accessToken.AccountId == qianchuanAdvertiser.AccountId && accessToken.AppId == qianchuanAdvertiser.AppId && accessToken.CompanyId == qianchuanAdvertiser.CompanyId {
			var gqcwg sync.WaitGroup

			gqcwg.Add(1)

			campaigns, err := qauc.listCampaigns(ctx, &gqcwg, limiter, day, qianchuanAdvertiser, accessToken, marketingGoal, marketingScene, 1)

			if err == nil {
				if campaigns.Data.PageInfo.TotalPage > 1 {
					var page uint32

					for page = 2; page <= campaigns.Data.PageInfo.TotalPage; page++ {
						gqcwg.Add(1)

						go qauc.listCampaigns(ctx, &gqcwg, limiter, day, qianchuanAdvertiser, accessToken, marketingGoal, marketingScene, page)
					}
				}
			}

			gqcwg.Wait()

			break
		}
	}
}

func (qauc *QianchuanAdvertiserUsecase) listCampaigns(ctx context.Context, gqcwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessToken *domain.OceanengineAccountToken, marketingGoal, marketingScene string, page uint32) (*campaign.ListCampaignResponse, error) {
	defer gqcwg.Done()

	var campaigns *campaign.ListCampaignResponse
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		limiter.Wait(ctx)

		campaigns, err = campaign.ListCampaign(qianchuanAdvertiser.AdvertiserId, accessToken.AccessToken, marketingGoal, marketingScene, page)

		if err != nil {
			if retryNum == 2 {
				inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, qianchuanAdvertiser.CompanyId, qianchuanAdvertiser.AccountId, qianchuanAdvertiser.AdvertiserId, 0, 0, qianchuanAdvertiser.AppId, accessToken.AccessToken, err.Error())
				inOceanengineApiLog.SetCreateTime(ctx)

				qauc.oalrepo.Save(ctx, inOceanengineApiLog)
			}
		} else {
			for _, campaign := range campaigns.Data.List {
				inQianchuanCampaign := domain.NewQianchuanCampaign(ctx, campaign.Id, qianchuanAdvertiser.AdvertiserId, campaign.Budget, campaign.Name, campaign.BudgetMode, campaign.MarketingGoal, campaign.MarketingScene, campaign.Status, campaign.CreateDate)
				inQianchuanCampaign.SetCreateTime(ctx)
				inQianchuanCampaign.SetUpdateTime(ctx)

				if err := qauc.qcrepo.Upsert(ctx, day, inQianchuanCampaign); err != nil {
					sinQianchuanCampaign, _ := json.Marshal(inQianchuanCampaign)

					inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanCampaignError] CampaignId=%d, AdvertiserId=%d, Data=%s, Description=%s", campaign.Id, qianchuanAdvertiser.AdvertiserId, sinQianchuanCampaign, "同步千川广告组，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					qauc.tlrepo.Save(ctx, inTaskLog)
				}

				inQianchuanAdInfo := domain.NewQianchuanAdQianchuanCampaign(ctx, qianchuanAdvertiser.AdvertiserId, campaign.Id, campaign.Budget, campaign.Name, campaign.BudgetMode, campaign.Status, campaign.CreateDate)
				inQianchuanAdInfo.SetUpdateTime(ctx)

				if err := qauc.qairepo.UpsertQianchuanCampaign(ctx, day, inQianchuanAdInfo); err != nil {
					sinQianchuanAdInfo, _ := json.Marshal(inQianchuanAdInfo)

					inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAdInfoError] CampaignId=%d, AdvertiserId=%d, Data=%s, Description=%s", campaign.Id, qianchuanAdvertiser.AdvertiserId, sinQianchuanAdInfo, "同步千川广告组数据到千川计划信息表，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					qauc.tlrepo.Save(ctx, inTaskLog)
				}
			}

			break
		}
	}

	return campaigns, err
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanAds(ctx context.Context, sqdwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertisers []*domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken) {
	defer sqdwg.Done()

	qauc.qadrepo.SaveIndex(ctx, day)
	qauc.qairepo.SaveIndex(ctx, day)

	var sqawg sync.WaitGroup

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		sqawg.Add(2)

		go qauc.getQianchuanAd(ctx, &sqawg, limiter, day, qianchuanAdvertiser, accessTokens, "VIDEO_PROM_GOODS")
		go qauc.getQianchuanAd(ctx, &sqawg, limiter, day, qianchuanAdvertiser, accessTokens, "LIVE_PROM_GOODS")
	}

	sqawg.Wait()
}

func (qauc *QianchuanAdvertiserUsecase) getQianchuanAd(ctx context.Context, sqcwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken, marketingGoal string) {
	defer sqcwg.Done()

	for _, accessToken := range accessTokens {
		if accessToken.AccountId == qianchuanAdvertiser.AccountId && accessToken.AppId == qianchuanAdvertiser.AppId && accessToken.CompanyId == qianchuanAdvertiser.CompanyId {
			var gqawg sync.WaitGroup
			var modifyTime string

			gqawg.Add(1)

			/*if total, err := qauc.qadrepo.CountByAdIdAndAdvertiserIdAndDay(ctx, qianchuanAdvertiser.AdvertiserId, day); err == nil {
				if total > 0 {
					modifyTime = tool.TimeToString("2006-01-02 15", time.Now().Add(-2*time.Hour))
				}
			}*/

			ads, err := qauc.listAds(ctx, &gqawg, limiter, day, qianchuanAdvertiser, accessToken, marketingGoal, modifyTime, 1)

			if err == nil {
				if ads.Data.PageInfo.TotalPage > 1 {
					var page uint32

					for page = 2; page <= ads.Data.PageInfo.TotalPage; page++ {
						gqawg.Add(1)

						go qauc.listAds(ctx, &gqawg, limiter, day, qianchuanAdvertiser, accessToken, marketingGoal, modifyTime, page)
					}
				}
			}

			gqawg.Wait()

			break
		}
	}
}

func (qauc *QianchuanAdvertiserUsecase) listAds(ctx context.Context, gqcwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessToken *domain.OceanengineAccountToken, marketingGoal, modifyTime string, page uint32) (*ad.ListAdResponse, error) {
	defer gqcwg.Done()

	var ads *ad.ListAdResponse
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		limiter.Wait(ctx)

		ads, err = ad.ListAd(qianchuanAdvertiser.AdvertiserId, accessToken.AccessToken, marketingGoal, modifyTime, page)

		if err != nil {
			if retryNum == 2 {
				inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, qianchuanAdvertiser.CompanyId, qianchuanAdvertiser.AccountId, qianchuanAdvertiser.AdvertiserId, 0, 0, qianchuanAdvertiser.AppId, accessToken.AccessToken, err.Error())
				inOceanengineApiLog.SetCreateTime(ctx)

				qauc.oalrepo.Save(ctx, inOceanengineApiLog)
			}
		} else {
			for _, ad := range ads.Data.List {
				productInfo := make([]*domain.ProductInfo, 0)
				awemeInfo := make([]*domain.AwemeInfo, 0)

				deliverySetting := domain.DeliverySetting{
					DeepExternalAction: ad.DeliverySetting.DeepExternalAction,
					DeepBidType:        ad.DeliverySetting.DeepBidType,
					RoiGoal:            ad.DeliverySetting.RoiGoal,
					SmartBidType:       ad.DeliverySetting.SmartBidType,
					ExternalAction:     ad.DeliverySetting.ExternalAction,
					Budget:             ad.DeliverySetting.Budget,
					ReviveBudget:       ad.DeliverySetting.ReviveBudget,
					BudgetMode:         ad.DeliverySetting.BudgetMode,
					CpaBid:             ad.DeliverySetting.CpaBid,
					StartTime:          ad.DeliverySetting.StartTime,
					EndTime:            ad.DeliverySetting.EndTime,
				}

				for _, lproductInfo := range ad.ProductInfo {
					productInfo = append(productInfo, &domain.ProductInfo{
						Id:                  lproductInfo.Id,
						Name:                lproductInfo.Name,
						DiscountPrice:       lproductInfo.DiscountPrice,
						Img:                 lproductInfo.Img,
						MarketPrice:         lproductInfo.MarketPrice,
						DiscountLowerPrice:  lproductInfo.DiscountLowerPrice,
						DiscountHigherPrice: lproductInfo.DiscountHigherPrice,
					})
				}

				for _, lawemeInfo := range ad.AwemeInfo {
					awemeInfo = append(awemeInfo, &domain.AwemeInfo{
						AwemeId:     lawemeInfo.AwemeId,
						AwemeName:   lawemeInfo.AwemeName,
						AwemeShowId: lawemeInfo.AwemeShowId,
						AwemeAvatar: lawemeInfo.AwemeAvatar,
					})
				}

				inQianchuanAd := domain.NewQianchuanAd(ctx, ad.AdId, qianchuanAdvertiser.AdvertiserId, ad.CampaignId, ad.PromotionWay, ad.MarketingGoal, ad.MarketingScene, ad.Name, ad.Status, ad.OptStatus, ad.AdCreateTime, ad.AdModifyTime, ad.LabAdType, productInfo, awemeInfo, deliverySetting)
				inQianchuanAd.SetCreateTime(ctx)
				inQianchuanAd.SetUpdateTime(ctx)

				if err := qauc.qadrepo.Upsert(ctx, day, inQianchuanAd); err != nil {
					sinQianchuanAd, _ := json.Marshal(inQianchuanAd)

					inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAdError] AdId=%d, AdvertiserId=%d, Data=%s, Description=%s", ad.AdId, qianchuanAdvertiser.AdvertiserId, sinQianchuanAd, "同步千川广告计划，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					qauc.tlrepo.Save(ctx, inTaskLog)
				}

				inQianchuanAdInfo := domain.NewQianchuanAdInfoQianchuanAd(ctx, ad.AdId, qianchuanAdvertiser.AdvertiserId, qianchuanAdvertiser.AdvertiserName, ad.PromotionWay, ad.MarketingGoal, ad.MarketingScene, ad.Name, ad.Status, ad.OptStatus, ad.AdCreateTime, ad.AdModifyTime, ad.LabAdType, productInfo, awemeInfo, deliverySetting)
				inQianchuanAdInfo.SetCreateTime(ctx)
				inQianchuanAdInfo.SetUpdateTime(ctx)

				if ad.CampaignId > 0 {
					inQianchuanAdInfo.SetCampaignId(ctx, ad.CampaignId)
				}

				if err := qauc.qairepo.UpsertQianchuanAd(ctx, day, inQianchuanAdInfo); err != nil {
					sinQianchuanAdInfo, _ := json.Marshal(inQianchuanAdInfo)

					inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAdError] AdId=%d, AdvertiserId=%d, Data=%s, Description=%s", ad.AdId, qianchuanAdvertiser.AdvertiserId, sinQianchuanAdInfo, "同步连山云RDS数据到千川计划信息表，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					qauc.tlrepo.Save(ctx, inTaskLog)
				}
			}

			break
		}
	}

	return ads, err
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanProducts(ctx context.Context, sqdwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertisers []*domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken) {
	defer sqdwg.Done()

	qauc.qprepo.SaveIndex(ctx, day)

	var sqpwg sync.WaitGroup

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		sqpwg.Add(1)

		go qauc.getQianchuanProduct(ctx, &sqpwg, limiter, day, "FEED", qianchuanAdvertiser, accessTokens)
	}

	sqpwg.Wait()
}

func (qauc *QianchuanAdvertiserUsecase) getQianchuanProduct(ctx context.Context, sqpwg *sync.WaitGroup, limiter *rate.Limiter, day, marketingScene string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken) {
	defer sqpwg.Done()

	for _, accessToken := range accessTokens {
		if accessToken.AccountId == qianchuanAdvertiser.AccountId && accessToken.AppId == qianchuanAdvertiser.AppId && accessToken.CompanyId == qianchuanAdvertiser.CompanyId {
			var gqpwg sync.WaitGroup

			gqpwg.Add(1)

			campaigns, err := qauc.listQianchuanProducts(ctx, &gqpwg, limiter, day, marketingScene, qianchuanAdvertiser, accessToken, 1)

			if err == nil {
				if campaigns.Data.PageInfo.TotalPage > 1 {
					var page uint32

					for page = 2; page <= campaigns.Data.PageInfo.TotalPage; page++ {
						gqpwg.Add(1)

						go qauc.listQianchuanProducts(ctx, &gqpwg, limiter, day, marketingScene, qianchuanAdvertiser, accessToken, page)
					}
				}
			}

			gqpwg.Wait()

			break
		}
	}
}

func (qauc *QianchuanAdvertiserUsecase) listQianchuanProducts(ctx context.Context, gqpwg *sync.WaitGroup, limiter *rate.Limiter, day, marketingScene string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessToken *domain.OceanengineAccountToken, page uint32) (*product.ListAvailableProductResponse, error) {
	defer gqpwg.Done()

	var qianchuanProducts *product.ListAvailableProductResponse
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		limiter.Wait(ctx)

		qianchuanProducts, err = product.ListAvailableProduct(qianchuanAdvertiser.AdvertiserId, accessToken.AccessToken, marketingScene, page)

		if err != nil {
			if retryNum == 2 {
				inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, qianchuanAdvertiser.CompanyId, qianchuanAdvertiser.AccountId, qianchuanAdvertiser.AdvertiserId, 0, 0, qianchuanAdvertiser.AppId, accessToken.AccessToken, err.Error())
				inOceanengineApiLog.SetCreateTime(ctx)

				qauc.oalrepo.Save(ctx, inOceanengineApiLog)
			}
		} else {
			for _, qianchuanProduct := range qianchuanProducts.Data.ProductList {
				imgList := make([]*domain.ImgUrl, 0)

				for _, imgUrl := range qianchuanProduct.ImgList {
					imgList = append(imgList, &domain.ImgUrl{
						ImgUrl: imgUrl.ImgUrl,
					})
				}

				inQianchuanProduct := domain.NewQianchuanProduct(ctx, qianchuanProduct.Id, qianchuanAdvertiser.AdvertiserId, qianchuanProduct.Inventory, qianchuanProduct.DiscountPrice, qianchuanProduct.DiscountLowerPrice, qianchuanProduct.DiscountHigherPrice, qianchuanProduct.MarketPrice, qianchuanProduct.ProductRate, qianchuanProduct.Name, qianchuanProduct.Img, qianchuanProduct.CategoryName, qianchuanProduct.SaleTime, qianchuanProduct.Tags, imgList)
				inQianchuanProduct.SetCreateTime(ctx)
				inQianchuanProduct.SetUpdateTime(ctx)

				if err := qauc.qprepo.Upsert(ctx, day, inQianchuanProduct); err != nil {
					sinQianchuanProduct, _ := json.Marshal(inQianchuanProduct)

					inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanProductError] ProductId=%d, AdvertiserId=%d, Data=%s, Description=%s", qianchuanProduct.Id, qianchuanAdvertiser.AdvertiserId, sinQianchuanProduct, "同步千川商家可投商品，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					qauc.tlrepo.Save(ctx, inTaskLog)
				}

				inQianchuanReportProduct := domain.NewQianchuanReportProduct(ctx, qianchuanAdvertiser.AdvertiserId, qianchuanProduct.Id, 0, 0, 0, 0, 0, qianchuanProduct.MarketPrice, 0.00, 0.00, qianchuanAdvertiser.AdvertiserName, qianchuanProduct.Name, qianchuanProduct.Img)
				inQianchuanReportProduct.SetCreateTime(ctx)
				inQianchuanReportProduct.SetUpdateTime(ctx)

				if err := qauc.qrprepo.UpsertQianchuanReportProductInfo(ctx, day, inQianchuanReportProduct); err != nil {
					sinQianchuanReportProduct, _ := json.Marshal(inQianchuanReportProduct)

					inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanProductError] ProductId=%d, AdvertiserId=%d, Data=%s, Description=%s", qianchuanProduct.Id, qianchuanAdvertiser.AdvertiserId, sinQianchuanReportProduct, "同步千川商家可投商品，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					qauc.tlrepo.Save(ctx, inTaskLog)
				}
			}

			break
		}
	}

	return qianchuanProducts, err
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanAwemes(ctx context.Context, sqdwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertisers []*domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken) {
	defer sqdwg.Done()

	qauc.qarepo.SaveIndex(ctx, day)

	var sqawg sync.WaitGroup

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		sqawg.Add(1)

		go qauc.getQianchuanAweme(ctx, &sqawg, limiter, day, qianchuanAdvertiser, accessTokens)
	}

	sqawg.Wait()
}

func (qauc *QianchuanAdvertiserUsecase) getQianchuanAweme(ctx context.Context, sqpwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken) {
	defer sqpwg.Done()

	for _, accessToken := range accessTokens {
		if accessToken.AccountId == qianchuanAdvertiser.AccountId && accessToken.AppId == qianchuanAdvertiser.AppId && accessToken.CompanyId == qianchuanAdvertiser.CompanyId {
			var gqawg sync.WaitGroup

			gqawg.Add(1)

			qianchuanAwemes, err := qauc.listQianchuanAwemes(ctx, &gqawg, limiter, day, qianchuanAdvertiser, accessToken, 1)

			if err == nil {
				if qianchuanAwemes.Data.PageInfo.TotalPage > 1 {
					var page uint32

					for page = 2; page <= qianchuanAwemes.Data.PageInfo.TotalPage; page++ {
						gqawg.Add(1)

						go qauc.listQianchuanAwemes(ctx, &gqawg, limiter, day, qianchuanAdvertiser, accessToken, page)
					}
				}
			}

			gqawg.Wait()

			break
		}
	}
}

func (qauc *QianchuanAdvertiserUsecase) listQianchuanAwemes(ctx context.Context, gqpwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessToken *domain.OceanengineAccountToken, page uint32) (*account.ListAwemeResponse, error) {
	defer gqpwg.Done()

	var qianchuanAwemes *account.ListAwemeResponse
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		limiter.Wait(ctx)

		qianchuanAwemes, err = account.ListAweme(qianchuanAdvertiser.AdvertiserId, accessToken.AccessToken, page)

		if err != nil {
			if retryNum == 2 {
				inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, qianchuanAdvertiser.CompanyId, qianchuanAdvertiser.AccountId, qianchuanAdvertiser.AdvertiserId, 0, 0, qianchuanAdvertiser.AppId, accessToken.AccessToken, err.Error())
				inOceanengineApiLog.SetCreateTime(ctx)

				qauc.oalrepo.Save(ctx, inOceanengineApiLog)
			}
		} else {
			for _, qianchuanAweme := range qianchuanAwemes.Data.AwemeIdList {
				bindType := make([]*string, 0)

				for _, lBindType := range qianchuanAweme.BindType {
					bindType = append(bindType, &lBindType)
				}

				inQianchuanAweme := domain.NewQianchuanAweme(ctx, qianchuanAweme.AwemeId, qianchuanAdvertiser.AdvertiserId, qianchuanAweme.AwemeAvatar, qianchuanAweme.AwemeShowId, qianchuanAweme.AwemeName, qianchuanAweme.AwemeStatus, qianchuanAweme.AwemeHasVideoPermission, qianchuanAweme.AwemeHasLivePermission, qianchuanAweme.AwemeHasUniProm, bindType)
				inQianchuanAweme.SetCreateTime(ctx)
				inQianchuanAweme.SetUpdateTime(ctx)

				if err := qauc.qarepo.Upsert(ctx, day, inQianchuanAweme); err != nil {
					sinQianchuanAweme, _ := json.Marshal(inQianchuanAweme)

					inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAwemeError] AwemeId=%d, AdvertiserId=%d, Data=%s, Description=%s", qianchuanAweme.AwemeId, qianchuanAdvertiser.AdvertiserId, sinQianchuanAweme, "同步千川账户已授权抖音号，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					qauc.tlrepo.Save(ctx, inTaskLog)
				}

				inQianchuanReportAweme := domain.NewQianchuanReportAweme(ctx, qianchuanAdvertiser.AdvertiserId, qianchuanAweme.AwemeId, 0, 0, 0, 0, 0, 0.00, 0.00, qianchuanAdvertiser.AdvertiserName, qianchuanAweme.AwemeName, qianchuanAweme.AwemeShowId, qianchuanAweme.AwemeAvatar)
				inQianchuanReportAweme.SetCreateTime(ctx)
				inQianchuanReportAweme.SetUpdateTime(ctx)

				if err := qauc.qrawrepo.UpsertQianchuanReportAwemeInfo(ctx, day, inQianchuanReportAweme); err != nil {
					sinQianchuanAweme, _ := json.Marshal(inQianchuanAweme)

					inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAwemeError] AwemeId=%d, AdvertiserId=%d, Data=%s, Description=%s", qianchuanAweme.AwemeId, qianchuanAdvertiser.AdvertiserId, sinQianchuanAweme, "同步千川账户已授权抖音号，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					qauc.tlrepo.Save(ctx, inTaskLog)
				}
			}

			break
		}
	}

	return qianchuanAwemes, err
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanWallets(ctx context.Context, sqdwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertisers []*domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken) {
	defer sqdwg.Done()

	qauc.qwrepo.SaveIndex(ctx, day)

	var sqwwg sync.WaitGroup

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		sqwwg.Add(1)

		go qauc.getQianchuanWallet(ctx, &sqwwg, limiter, day, qianchuanAdvertiser, accessTokens)
	}

	sqwwg.Wait()
}

func (qauc *QianchuanAdvertiserUsecase) getQianchuanWallet(ctx context.Context, sqpwg *sync.WaitGroup, limiter *rate.Limiter, day string, qianchuanAdvertiser *domain.QianchuanAdvertiser, accessTokens []*domain.OceanengineAccountToken) {
	defer sqpwg.Done()

	for _, accessToken := range accessTokens {
		if accessToken.AccountId == qianchuanAdvertiser.AccountId && accessToken.AppId == qianchuanAdvertiser.AppId && accessToken.CompanyId == qianchuanAdvertiser.CompanyId {
			var qianchuanwallet *finance.GetWalletResponse
			var err error

			for retryNum := 0; retryNum < 3; retryNum++ {
				limiter.Wait(ctx)

				qianchuanwallet, err = finance.GetWallet(qianchuanAdvertiser.AdvertiserId, accessToken.AccessToken)

				if err != nil {
					if retryNum == 2 {
						inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, qianchuanAdvertiser.CompanyId, qianchuanAdvertiser.AccountId, qianchuanAdvertiser.AdvertiserId, 0, 0, qianchuanAdvertiser.AppId, accessToken.AccessToken, err.Error())
						inOceanengineApiLog.SetCreateTime(ctx)

						qauc.oalrepo.Save(ctx, inOceanengineApiLog)
					}
				} else {
					shareExpiringDetailList := make([]*domain.ShareExpiringDetailList, 0)

					inQianchuanWallet := domain.NewQianchuanWallet(ctx, qianchuanAdvertiser.AdvertiserId,
						qianchuanwallet.Data.TotalBalanceAbs,
						qianchuanwallet.Data.GrantBalance,
						qianchuanwallet.Data.UnionValidGrantBalance,
						qianchuanwallet.Data.SearchValidGrantBalance,
						qianchuanwallet.Data.CommonValidGrantBalance,
						qianchuanwallet.Data.DefaultValidGrantBalance,
						qianchuanwallet.Data.GeneralTotalBalance,
						qianchuanwallet.Data.GeneralBalanceValid,
						qianchuanwallet.Data.GeneralBalanceValidNonGrant,
						qianchuanwallet.Data.GeneralBalanceValidGrantUnion,
						qianchuanwallet.Data.GeneralBalanceValidGrantSearch,
						qianchuanwallet.Data.GeneralBalanceValidGrantCommon,
						qianchuanwallet.Data.GeneralBalanceValidGrantDefault,
						qianchuanwallet.Data.GeneralBalanceInvalid,
						qianchuanwallet.Data.GeneralBalanceInvalidOrder,
						qianchuanwallet.Data.GeneralBalanceInvalidFrozen,
						qianchuanwallet.Data.BrandBalance,
						qianchuanwallet.Data.BrandBalanceValid,
						qianchuanwallet.Data.BrandBalanceValidNonGrant,
						qianchuanwallet.Data.BrandBalanceValidGrant,
						qianchuanwallet.Data.BrandBalanceInvalid,
						qianchuanwallet.Data.BrandBalanceInvalidFrozen,
						qianchuanwallet.Data.DeductionCouponBalance,
						qianchuanwallet.Data.DeductionCouponBalanceAll,
						qianchuanwallet.Data.DeductionCouponBalanceOther,
						qianchuanwallet.Data.DeductionCouponBalanceSelf,
						qianchuanwallet.Data.GrantExpiring,
						qianchuanwallet.Data.ShareBalance,
						qianchuanwallet.Data.ShareBalanceValidGrantUnion,
						qianchuanwallet.Data.ShareBalanceValidGrantSearch,
						qianchuanwallet.Data.ShareBalanceValidGrantCommon,
						qianchuanwallet.Data.ShareBalanceValidGrantDefault,
						qianchuanwallet.Data.ShareBalanceValid,
						qianchuanwallet.Data.ShareBalanceExpiring,
						shareExpiringDetailList)
					inQianchuanWallet.SetCreateTime(ctx)
					inQianchuanWallet.SetUpdateTime(ctx)

					if err := qauc.qwrepo.Upsert(ctx, day, inQianchuanWallet); err != nil {
						sinQianchuanWallet, _ := json.Marshal(inQianchuanWallet)

						inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanWalletError] AdvertiserId=%d, Data=%s, Description=%s", qianchuanAdvertiser.AdvertiserId, sinQianchuanWallet, "同步账户钱包信息，插入数据库失败"))
						inTaskLog.SetCreateTime(ctx)

						qauc.tlrepo.Save(ctx, inTaskLog)
					}

					break
				}
			}
		}
	}
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanAwemwAndProducts(ctx context.Context, wg *sync.WaitGroup, day string) {
	defer wg.Done()

	qauc.qrawrepo.SaveIndex(ctx, day)
	qauc.qrprepo.SaveIndex(ctx, day)

	var sqaapwg sync.WaitGroup

	sqaapwg.Add(2)

	go qauc.syncQianchuanReportAwemw(ctx, &sqaapwg, day)
	go qauc.syncQianchuanReportProduct(ctx, &sqaapwg, day)

	sqaapwg.Wait()
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanReportAwemw(ctx context.Context, sqaapwg *sync.WaitGroup, day string) {
	defer sqaapwg.Done()

	var sqrawg sync.WaitGroup
	var lerr error
	var rerr error
	var listAwemeAds []*domain.QianchuanReportAweme
	var listAwemeReportAds []*domain.QianchuanReportAd

	sqrawg.Add(2)

	go func() {
		defer sqrawg.Done()

		listAwemeAds, lerr = qauc.qadrepo.ListAwemeAd(ctx, day)
	}()

	go func() {
		defer sqrawg.Done()

		listAwemeReportAds, rerr = qauc.qraurepo.ListByMarketingGoal(ctx, "LIVE_PROM_GOODS", day)
	}()

	sqrawg.Wait()

	if lerr != nil || rerr != nil {
		inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAwemwAndProductsError] Description=%s", "同步千川达人，查询数据库失败"))
		inTaskLog.SetCreateTime(ctx)

		qauc.tlrepo.Save(ctx, inTaskLog)
	}

	for _, listAwemeAd := range listAwemeAds {
		for _, adId := range listAwemeAd.AdIds {
			for _, listAwemeReportAd := range listAwemeReportAds {
				if listAwemeReportAd.AdId == adId && listAwemeReportAd.AdvertiserId == listAwemeAd.AdvertiserId {
					listAwemeAd.DyFollow += listAwemeReportAd.DyFollow
					listAwemeAd.StatCost += listAwemeReportAd.StatCost
					listAwemeAd.PayOrderCount += listAwemeReportAd.PayOrderCount
					listAwemeAd.PayOrderAmount += listAwemeReportAd.PayOrderAmount
					listAwemeAd.ShowCnt += listAwemeReportAd.ShowCnt
					listAwemeAd.ClickCnt += listAwemeReportAd.ClickCnt
					listAwemeAd.ConvertCnt += listAwemeReportAd.ConvertCnt

					break
				}
			}
		}

		inQianchuanReportAweme := domain.NewQianchuanReportAweme(ctx, listAwemeAd.AdvertiserId, listAwemeAd.AwemeId, listAwemeAd.DyFollow, listAwemeAd.PayOrderCount, listAwemeAd.ShowCnt, listAwemeAd.ClickCnt, listAwemeAd.ConvertCnt, listAwemeAd.StatCost, listAwemeAd.PayOrderAmount, "", listAwemeAd.AwemeName, listAwemeAd.AwemeShowId, listAwemeAd.AwemeAvatar)
		inQianchuanReportAweme.SetCreateTime(ctx)
		inQianchuanReportAweme.SetUpdateTime(ctx)

		if err := qauc.qrawrepo.UpsertQianchuanReportAweme(ctx, day, inQianchuanReportAweme); err != nil {
			sinQianchuanReportAweme, _ := json.Marshal(inQianchuanReportAweme)

			inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAwemwAndProductsError] Data=%s, Description=%s", sinQianchuanReportAweme, "同步千川达人，插入数据库失败"))
			inTaskLog.SetCreateTime(ctx)

			qauc.tlrepo.Save(ctx, inTaskLog)
		}
	}
}

func (qauc *QianchuanAdvertiserUsecase) syncQianchuanReportProduct(ctx context.Context, sqaapwg *sync.WaitGroup, day string) {
	defer sqaapwg.Done()

	var sqrpwg sync.WaitGroup
	var lerr error
	var rerr error
	var listProductAds []*domain.QianchuanReportProduct
	var listProductReportAds []*domain.QianchuanReportAd

	sqrpwg.Add(2)

	go func() {
		defer sqrpwg.Done()

		listProductAds, lerr = qauc.qadrepo.ListProductAd(ctx, day)
	}()

	go func() {
		defer sqrpwg.Done()

		listProductReportAds, rerr = qauc.qraurepo.ListByMarketingGoal(ctx, "VIDEO_PROM_GOODS", day)
	}()

	sqrpwg.Wait()

	if lerr != nil || rerr != nil {
		inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAwemwAndProductsError] Description=%s", "同步千川商品，查询数据库失败"))
		inTaskLog.SetCreateTime(ctx)

		qauc.tlrepo.Save(ctx, inTaskLog)
	}

	for _, listProductAd := range listProductAds {
		for _, adId := range listProductAd.AdIds {
			for _, listProductReportAd := range listProductReportAds {
				if listProductReportAd.AdId == adId && listProductReportAd.AdvertiserId == listProductAd.AdvertiserId {
					listProductAd.DyFollow += listProductReportAd.DyFollow
					listProductAd.StatCost += listProductReportAd.StatCost
					listProductAd.PayOrderCount += listProductReportAd.PayOrderCount
					listProductAd.PayOrderAmount += listProductReportAd.PayOrderAmount
					listProductAd.ShowCnt += listProductReportAd.ShowCnt
					listProductAd.ClickCnt += listProductReportAd.ClickCnt
					listProductAd.ConvertCnt += listProductReportAd.ConvertCnt

					break
				}
			}
		}

		inQianchuanReportProduct := domain.NewQianchuanReportProduct(ctx, listProductAd.AdvertiserId, listProductAd.ProductId, listProductAd.PayOrderCount, listProductAd.ShowCnt, listProductAd.ClickCnt, listProductAd.ConvertCnt, listProductAd.DyFollow, listProductAd.DiscountPrice, listProductAd.StatCost, listProductAd.PayOrderAmount, "", listProductAd.ProductName, listProductAd.ProductImg)
		inQianchuanReportProduct.SetCreateTime(ctx)
		inQianchuanReportProduct.SetUpdateTime(ctx)

		if err := qauc.qrprepo.UpsertQianchuanReportProduct(ctx, day, inQianchuanReportProduct); err != nil {
			sinQianchuanReportProduct, _ := json.Marshal(inQianchuanReportProduct)

			inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanData", fmt.Sprintf("[SyncQianchuanDataError SyncQianchuanAwemwAndProductsError] Data=%s, Description=%s", sinQianchuanReportProduct, "同步千川商品，插入数据库失败"))
			inTaskLog.SetCreateTime(ctx)

			qauc.tlrepo.Save(ctx, inTaskLog)
		}
	}
}

func (qauc *QianchuanAdvertiserUsecase) SyncRdsDatas(ctx context.Context) error {
	list, err := qauc.ocrepo.List(ctx, 1, 0, 0)

	if err != nil {
		return DouyinOceanengineConfigNotFound
	}

	t := time.Now()

	var swg sync.WaitGroup

	for _, oceanengineConfig := range list {
		swg.Add(1)

		go qauc.SyncRdsData(ctx, &swg, oceanengineConfig, t)
	}

	swg.Wait()

	var stime string

	if t.Minute()%10 < 5 {
		stime = tool.TimeToString("2006-01-02 15", t) + ":" + strconv.Itoa(t.Minute()/10) + "0:00"
	} else {
		stime = tool.TimeToString("2006-01-02 15", t) + ":" + strconv.Itoa(t.Minute()/10) + "5:00"
	}

	messageAd := domain.MessageAd{
		Type: "qianchuan_report_ad_realtime_data_ready",
	}
	messageAd.Message.Name = "douyin"
	messageAd.Message.SyncTime = stime
	messageAd.Message.SendTime = tool.TimeToString("2006-01-02 15:04:05", time.Now())

	bmessageAd, _ := json.Marshal(messageAd)

	qauc.repo.Send(ctx, kafka.NewMessage(strconv.FormatUint(1, 10), bmessageAd))

	return nil
}

func (qauc *QianchuanAdvertiserUsecase) SyncRdsData(ctx context.Context, swg *sync.WaitGroup, oceanengineConfig *domain.OceanengineConfig, t time.Time) error {
	defer swg.Done()

	var srdwg sync.WaitGroup

	srdwg.Add(1)

	go qauc.SyncRdsAdData(ctx, &srdwg, oceanengineConfig, t)

	srdwg.Wait()

	return nil

}

func (qauc *QianchuanAdvertiserUsecase) SyncRdsAdData(ctx context.Context, srdwg *sync.WaitGroup, oceanengineConfig *domain.OceanengineConfig, t time.Time) error {
	defer srdwg.Done()

	day := tool.TimeToString("2006-01-02", t)

	rdsDatas, err := qauc.lrrepo.List(ctx, oceanengineConfig.AppId, tool.GetMondayOfWeek("20060102", t), tool.TimeToString("2006-01-02", t), "ad_id")

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "syncRdsData", fmt.Sprintf("[SyncRdsDataError SyncRdsAdDataError] AppId=%d, Description=%s", oceanengineConfig.AppId, "查询连山云RDS数据失败"))
		inTaskLog.SetCreateTime(ctx)

		qauc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	qauc.qraurepo.SaveIndex(ctx, day)
	qauc.qairepo.SaveIndex(ctx, day)
	qauc.qrarrepo.SaveIndex(ctx, day)

	var stime string

	if t.Minute()%10 < 5 {
		stime = tool.TimeToString("2006-01-02 15", t) + ":" + strconv.Itoa(t.Minute()/10) + "0:00"
	} else {
		stime = tool.TimeToString("2006-01-02 15", t) + ":" + strconv.Itoa(t.Minute()/10) + "5:00"
	}

	ttime, _ := tool.StringToTime("2006-01-02 15:04:05", stime)

	for _, rdsData := range rdsDatas {
		inQianchuanReportAd := domain.NewQianchuanReportAd(ctx, uint64(rdsData.AdId), uint64(rdsData.AdvertiserId), uint64(rdsData.AwemeId), rdsData.ShowCnt, rdsData.ClickCnt, rdsData.PayOrderCount, rdsData.CreateOrderCount, rdsData.DyFollow, rdsData.ConvertCnt, rdsData.MarketingGoal, rdsData.StatCost, rdsData.CreateOrderAmount, rdsData.PayOrderAmount)
		inQianchuanReportAd.SetCreateTime(ctx)
		inQianchuanReportAd.SetUpdateTime(ctx)

		if err := qauc.qraurepo.Upsert(ctx, day, inQianchuanReportAd); err != nil {
			sinQianchuanReportAd, _ := json.Marshal(inQianchuanReportAd)

			inTaskLog := domain.NewTaskLog(ctx, "syncRdsData", fmt.Sprintf("[SyncRdsDataError] AppId=%d, AdId=%d, AdvertiserId=%d, Data=%s, Description=%s", oceanengineConfig.AppId, uint64(rdsData.AdId), uint64(rdsData.AdvertiserId), sinQianchuanReportAd, "同步连山云RDS数据失，插入数据库失败"))
			inTaskLog.SetCreateTime(ctx)

			qauc.tlrepo.Save(ctx, inTaskLog)
		}

		inQianchuanAdInfo := domain.NewQianchuanAdInfoQianchuanReportAd(ctx, uint64(rdsData.AdId), uint64(rdsData.AdvertiserId), rdsData.StatCost, rdsData.PayOrderAmount, rdsData.CreateOrderAmount, rdsData.PayOrderCount, rdsData.CreateOrderCount, rdsData.ClickCnt, rdsData.ShowCnt, rdsData.ConvertCnt, rdsData.DyFollow)
		inQianchuanAdInfo.SetRoi(ctx)
		inQianchuanAdInfo.SetClickRate(ctx)
		inQianchuanAdInfo.SetCpmPlatform(ctx)
		inQianchuanAdInfo.SetPayConvertRate(ctx)
		inQianchuanAdInfo.SetConvertCost(ctx)
		inQianchuanAdInfo.SetConvertRate(ctx)
		inQianchuanAdInfo.SetAveragePayOrderStatCost(ctx)
		inQianchuanAdInfo.SetPayOrderAveragePrice(ctx)
		inQianchuanAdInfo.SetCreateTime(ctx)
		inQianchuanAdInfo.SetUpdateTime(ctx)

		if err := qauc.qairepo.UpsertQianchuanReportAd(ctx, day, inQianchuanAdInfo); err != nil {
			sinQianchuanAdInfo, _ := json.Marshal(inQianchuanAdInfo)

			inTaskLog := domain.NewTaskLog(ctx, "syncRdsData", fmt.Sprintf("[SyncRdsDataError] AppId=%d, AdId=%d, AdvertiserId=%d, Data=%s, Description=%s", oceanengineConfig.AppId, uint64(rdsData.AdId), uint64(rdsData.AdvertiserId), sinQianchuanAdInfo, "同步连山云RDS数据到千川计划信息表，插入数据库失败"))
			inTaskLog.SetCreateTime(ctx)

			qauc.tlrepo.Save(ctx, inTaskLog)
		}

		inQianchuanReportAdRealtime := domain.NewQianchuanReportAdRealtime(ctx, uint64(rdsData.AdId), uint64(rdsData.AdvertiserId), uint64(rdsData.AwemeId), rdsData.MarketingGoal, rdsData.ShowCnt, rdsData.ClickCnt, rdsData.PayOrderCount, rdsData.CreateOrderCount, rdsData.DyFollow, rdsData.ConvertCnt, ttime.Unix(), rdsData.StatCost, rdsData.CreateOrderAmount, rdsData.PayOrderAmount)
		inQianchuanReportAdRealtime.SetCreateTime(ctx)
		inQianchuanReportAdRealtime.SetUpdateTime(ctx)

		if err := qauc.qrarrepo.Upsert(ctx, day, inQianchuanReportAdRealtime); err != nil {
			sinQianchuanReportAdRealtime, _ := json.Marshal(inQianchuanReportAdRealtime)

			inTaskLog := domain.NewTaskLog(ctx, "syncRdsData", fmt.Sprintf("[SyncRdsDataError] AppId=%d, AdId=%d, AdvertiserId=%d, Data=%s, Description=%s", oceanengineConfig.AppId, uint64(rdsData.AdId), uint64(rdsData.AdvertiserId), sinQianchuanReportAdRealtime, "同步连山云RDS数据到千川计划信息实时表，插入数据库失败"))
			inTaskLog.SetCreateTime(ctx)

			qauc.tlrepo.Save(ctx, inTaskLog)
		}
	}

	return nil
}

func (qauc *QianchuanAdvertiserUsecase) verifyQianchuanAdvertiserLimit(ctx context.Context, companyId uint64) error {
	company, err := qauc.crepo.GetById(ctx, companyId)

	if err != nil {
		return DouyinDataError
	}

	var totalNums int64 = 0

	if total, err := qauc.repo.Count(ctx, company.Data.Id, "", "", "1"); err == nil {
		totalNums = total
	}

	if totalNums >= int64(company.Data.QianchuanAdvertisers) {
		return DouyinQianchuanAdvertiserLimitError
	}

	return nil
}
