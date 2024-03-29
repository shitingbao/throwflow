package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/event/event"
	"douyin/internal/pkg/event/kafka"
	"douyin/internal/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"sync"
	"time"
)

var (
	DouyinQianchuanAdListError = errors.InternalServer("DOUYIN_QIANCHUAN_CAMPAIGN_LIST_ERROR", "千川广告计划获取失败")
)

type QianchuanAdRepo interface {
	GetByAdIdAndDay(context.Context, uint64, string) (*domain.QianchuanAd, error)
	ListNotLadAd(context.Context, string, string) (map[uint64]map[uint64][]*domain.QianchuanCampaign, error)
	ListProductAd(context.Context, string) ([]*domain.QianchuanReportProduct, error)
	ListAwemeAd(context.Context, string) ([]*domain.QianchuanReportAweme, error)
	ListAwemeByAdvertiserId(context.Context, uint64, string) ([]*domain.QianchuanReportAweme, error)
	ListByPromotionId(context.Context, uint64, uint64, string, string) ([]*domain.QianchuanReportAd, error)
	ListByCampaignId(context.Context, uint64, string) ([]*domain.QianchuanAd, error)
	List(context.Context, string, string, string, int64, int64) ([]*domain.QianchuanAd, error)
	AllNotLadAd(context.Context, string) ([]*domain.QianchuanAd, error)
	CountByAdIdAndAdvertiserIdAndDay(context.Context, uint64, string) (int64, error)
	Count(context.Context, string, string, string, string, string) (int64, error)
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.QianchuanAd) error

	Send(context.Context, event.Event) error
}

type QianchuanAdUsecase struct {
	repo     QianchuanAdRepo
	qrarepo  QianchuanReportAdRepo
	qrarrepo QianchuanReportAdRealtimeRepo
	crepo    CompanyRepo
	csrepo   CompanySetRepo
	qairepo  QianchuanAdInfoRepo
	qarepo   QianchuanAdvertiserRepo
	qasrepo  QianchuanAdvertiserStatusRepo
	qainrepo QianchuanAdvertiserInfoRepo
	qcrepo   QianchuanCampaignRepo
	qwrepo   QianchuanWalletRepo
	tlrepo   TaskLogRepo
	conf     *conf.Data
	econf    *conf.Event
	log      *log.Helper
}

func NewQianchuanAdUsecase(repo QianchuanAdRepo, qrarepo QianchuanReportAdRepo, qrarrepo QianchuanReportAdRealtimeRepo, crepo CompanyRepo, csrepo CompanySetRepo, qairepo QianchuanAdInfoRepo, qarepo QianchuanAdvertiserRepo, qasrepo QianchuanAdvertiserStatusRepo, qainrepo QianchuanAdvertiserInfoRepo, qcrepo QianchuanCampaignRepo, qwrepo QianchuanWalletRepo, tlrepo TaskLogRepo, conf *conf.Data, econf *conf.Event, logger log.Logger) *QianchuanAdUsecase {
	return &QianchuanAdUsecase{repo: repo, qrarepo: qrarepo, qrarrepo: qrarrepo, crepo: crepo, csrepo: csrepo, qairepo: qairepo, qarepo: qarepo, qasrepo: qasrepo, qainrepo: qainrepo, qcrepo: qcrepo, qwrepo: qwrepo, tlrepo: tlrepo, conf: conf, econf: econf, log: log.NewHelper(logger)}
}

func (qauc *QianchuanAdUsecase) GetExternalQianchuanAds(ctx context.Context, adId uint64, day string) ([]*domain.ExternalQianchuanReportAdRealtime, error) {
	list := make([]*domain.ExternalQianchuanReportAdRealtime, 0)

	now := time.Now()

	ttime, err := tool.StringToTime("2006-01-02", day)
	tnow, err := tool.StringToTime("2006-01-02", now.Format("2006-01-02"))

	if err != nil {
		return nil, DouyinQianchuanAdListError
	}

	if ttime.Equal(tnow) {
		startTime, _ := tool.StringToTime("2006-01-02 15:04", now.Format("2006-01-02")+" 00:00")
		endTime, _ := tool.StringToTime("2006-01-02 15:04", now.Format("2006-01-02 15:04"))

		for startTime.Before(endTime) {
			list = append(list, &domain.ExternalQianchuanReportAdRealtime{
				AdId: adId,
				Time: startTime.Unix(),
			})

			startTime = startTime.Add(time.Minute * 5)
		}
	} else if ttime.Before(tnow) {
		startTime, _ := tool.StringToTime("2006-01-02 15:04", day+" 00:00")
		endTime, _ := tool.StringToTime("2006-01-02 15:04", day+" 23:59")

		for startTime.Before(endTime) {
			list = append(list, &domain.ExternalQianchuanReportAdRealtime{
				AdId: adId,
				Time: startTime.Unix(),
			})

			startTime = startTime.Add(time.Minute * 5)
		}
	}

	var wg sync.WaitGroup

	for _, l := range list {
		wg.Add(1)

		go func(l *domain.ExternalQianchuanReportAdRealtime) {
			defer wg.Done()

			if qianchuanAd, err := qauc.qrarrepo.GetByTime(ctx, l.AdId, l.Time, time.Unix(l.Time, 0).Format("2006-01-02")); err == nil {
				l.AdId = qianchuanAd.AdId
				l.AdvertiserId = qianchuanAd.AdvertiserId
				l.AwemeId = qianchuanAd.AwemeId
				l.MarketingGoal = qianchuanAd.MarketingGoal
				l.StatCost = qianchuanAd.StatCost
				l.ShowCnt = qianchuanAd.ShowCnt
				l.ClickCnt = qianchuanAd.ClickCnt
				l.PayOrderCount = qianchuanAd.PayOrderCount
				l.CreateOrderAmount = qianchuanAd.CreateOrderAmount
				l.CreateOrderCount = qianchuanAd.CreateOrderCount
				l.PayOrderAmount = qianchuanAd.PayOrderAmount
				l.DyFollow = qianchuanAd.DyFollow
				l.ConvertCnt = qianchuanAd.ConvertCnt

				l.SetRoi(ctx)
				l.SetClickRate(ctx)
				l.SetCpmPlatform(ctx)
				l.SetPayConvertRate(ctx)
				l.SetConvertCost(ctx)
				l.SetConvertRate(ctx)
				l.SetAveragePayOrderStatCost(ctx)
				l.SetPayOrderAveragePrice(ctx)
			}
		}(l)
	}

	wg.Wait()

	return list, nil
}

func (qauc *QianchuanAdUsecase) GetExternalHistoryQianchuanAds(ctx context.Context, adId uint64, startDay, endDay string) ([]*domain.ExternalQianchuanReportAdRealtime, error) {
	list := make([]*domain.ExternalQianchuanReportAdRealtime, 0)

	startTime, _ := tool.StringToTime("2006-01-02", startDay)

	endTime, _ := tool.StringToTime("2006-01-02", endDay)
	endTime = endTime.Add(time.Hour * 24)

	for startTime.Before(endTime) {
		list = append(list, &domain.ExternalQianchuanReportAdRealtime{
			AdId: adId,
			Time: startTime.Unix(),
		})

		startTime = startTime.Add(time.Hour * 24)
	}

	var wg sync.WaitGroup

	for _, l := range list {
		wg.Add(1)

		go func(l *domain.ExternalQianchuanReportAdRealtime) {
			defer wg.Done()

			if qianchuanAd, err := qauc.qrarrepo.Get(ctx, l.AdId, time.Unix(l.Time, 0).Format("2006-01-02")); err == nil {
				l.AdId = qianchuanAd.AdId
				l.AdvertiserId = qianchuanAd.AdvertiserId
				l.AwemeId = qianchuanAd.AwemeId
				l.MarketingGoal = qianchuanAd.MarketingGoal
				l.StatCost = qianchuanAd.StatCost
				l.ShowCnt = qianchuanAd.ShowCnt
				l.ClickCnt = qianchuanAd.ClickCnt
				l.PayOrderCount = qianchuanAd.PayOrderCount
				l.CreateOrderAmount = qianchuanAd.CreateOrderAmount
				l.CreateOrderCount = qianchuanAd.CreateOrderCount
				l.PayOrderAmount = qianchuanAd.PayOrderAmount
				l.DyFollow = qianchuanAd.DyFollow
				l.ConvertCnt = qianchuanAd.ConvertCnt

				l.SetRoi(ctx)
				l.SetClickRate(ctx)
				l.SetCpmPlatform(ctx)
				l.SetPayConvertRate(ctx)
				l.SetConvertCost(ctx)
				l.SetConvertRate(ctx)
				l.SetAveragePayOrderStatCost(ctx)
				l.SetPayOrderAveragePrice(ctx)
			}
		}(l)
	}

	wg.Wait()

	return list, nil
}

func (qauc *QianchuanAdUsecase) ListQianchuanAds(ctx context.Context, pageNum, pageSize uint64, day, keyword, advertiserIds string) (*domain.QianchuanAdList, error) {
	list, err := qauc.repo.List(ctx, advertiserIds, day, keyword, int64(pageNum), int64(pageSize))

	if err != nil {
		return nil, DouyinQianchuanAdListError
	}

	total, err := qauc.repo.Count(ctx, advertiserIds, day, keyword, "", "")

	if err != nil {
		return nil, DouyinQianchuanAdListError
	}

	return &domain.QianchuanAdList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (qauc *QianchuanAdUsecase) ListQianchuanReportAdvertisers(ctx context.Context, day, advertiserIds string) ([]*domain.QianchuanReportAdRealtime, error) {
	list, err := qauc.qrarrepo.ListAdvertisers(ctx, advertiserIds, day)

	if err != nil {
		return nil, DouyinQianchuanAdvertiserListError
	}

	return list, nil
}

func (qauc *QianchuanAdUsecase) ListExternalQianchuanAds(ctx context.Context, pageNum, pageSize uint64, startDay, endDay, keyword, advertiserIds, filter, orderName, orderType string) (*domain.ExternalQianchuanAdList, error) {
	list := make([]*domain.ExternalQianchuanAd, 0)

	qianchuanAdInfos, err := qauc.qairepo.List(ctx, advertiserIds, startDay, endDay, keyword, filter, orderName, orderType, pageNum, pageSize)

	if err != nil {
		return nil, DouyinQianchuanAdListError
	}

	total, err := qauc.qairepo.Count(ctx, advertiserIds, startDay, endDay, keyword, filter)

	if err != nil {
		return nil, DouyinQianchuanAdListError
	}

	var wg sync.WaitGroup

	var qianchuanAdInfoList sync.Map

	for _, lqianchuanAdInfo := range qianchuanAdInfos {
		wg.Add(1)

		go func(adId uint64) {
			defer wg.Done()

			qianchuanAdInfo, _ := qauc.qairepo.GetByDay(ctx, adId, startDay, endDay)

			qianchuanAdInfoList.Store(adId, qianchuanAdInfo)
		}(lqianchuanAdInfo.AdId)
	}

	wg.Wait()

	for _, lqianchuanAdInfo := range qianchuanAdInfos {
		qianchuanAdInfoList.Range(func(k, v interface{}) bool {
			if lqianchuanAdInfo.AdId == k {
				var qianchuanAdInfo *domain.QianchuanAdInfo
				qianchuanAdInfo = v.(*domain.QianchuanAdInfo)

				adCreateTime, _ := tool.StringToTime("2006-01-02 15:04:05", qianchuanAdInfo.AdCreateTime)
				adModifyTime, _ := tool.StringToTime("2006-01-02 15:04:05", qianchuanAdInfo.AdModifyTime)

				qianchuanAdInfo.StatCost = lqianchuanAdInfo.StatCost
				qianchuanAdInfo.ShowCnt = lqianchuanAdInfo.ShowCnt
				qianchuanAdInfo.ClickCnt = lqianchuanAdInfo.ClickCnt
				qianchuanAdInfo.PayOrderCount = lqianchuanAdInfo.PayOrderCount
				qianchuanAdInfo.CreateOrderAmount = lqianchuanAdInfo.CreateOrderAmount
				qianchuanAdInfo.CreateOrderCount = lqianchuanAdInfo.CreateOrderCount
				qianchuanAdInfo.PayOrderAmount = lqianchuanAdInfo.PayOrderAmount
				qianchuanAdInfo.DyFollow = lqianchuanAdInfo.DyFollow
				qianchuanAdInfo.ConvertCnt = lqianchuanAdInfo.ConvertCnt

				qianchuanAdInfo.SetRoi(ctx)
				qianchuanAdInfo.SetClickRate(ctx)
				qianchuanAdInfo.SetCpmPlatform(ctx)
				qianchuanAdInfo.SetPayConvertRate(ctx)
				qianchuanAdInfo.SetConvertCost(ctx)
				qianchuanAdInfo.SetConvertRate(ctx)
				qianchuanAdInfo.SetAveragePayOrderStatCost(ctx)
				qianchuanAdInfo.SetPayOrderAveragePrice(ctx)

				list = append(list, &domain.ExternalQianchuanAd{
					AdId:                    qianchuanAdInfo.AdId,
					AdvertiserId:            qianchuanAdInfo.AdvertiserId,
					AdvertiserName:          qianchuanAdInfo.AdvertiserName,
					CampaignId:              qianchuanAdInfo.CampaignId,
					CampaignName:            qianchuanAdInfo.CampaignName,
					AdName:                  qianchuanAdInfo.Name,
					LabAdType:               qianchuanAdInfo.LabAdType,
					LabAdTypeName:           qianchuanAdInfo.GetLabAdTypeName(ctx),
					MarketingGoal:           qianchuanAdInfo.MarketingGoal,
					MarketingGoalName:       qianchuanAdInfo.GetMarketingGoalName(ctx),
					Status:                  qianchuanAdInfo.Status,
					StatusName:              qianchuanAdInfo.GetStatusName(ctx),
					OptStatus:               qianchuanAdInfo.OptStatus,
					OptStatusName:           qianchuanAdInfo.GetOptStatusName(ctx),
					ExternalAction:          qianchuanAdInfo.DeliverySetting.ExternalAction,
					ExternalActionName:      qianchuanAdInfo.GetExternalActionName(ctx),
					DeepExternalAction:      qianchuanAdInfo.DeliverySetting.DeepExternalAction,
					DeepBidType:             qianchuanAdInfo.DeliverySetting.DeepBidType,
					PromotionId:             qianchuanAdInfo.GetPromotionId(ctx),
					PromotionShowId:         qianchuanAdInfo.GetPromotionShowId(ctx),
					PromotionName:           qianchuanAdInfo.GetPromotionName(ctx),
					PromotionImg:            qianchuanAdInfo.GetPromotionImg(ctx),
					PromotionAvatar:         qianchuanAdInfo.GetPromotionAvatar(ctx),
					PromotionType:           qianchuanAdInfo.GetPromotionType(ctx),
					StatCost:                lqianchuanAdInfo.StatCost,
					Roi:                     qianchuanAdInfo.Roi,
					CpaBid:                  qianchuanAdInfo.DeliverySetting.CpaBid,
					RoiGoal:                 qianchuanAdInfo.DeliverySetting.RoiGoal,
					Budget:                  qianchuanAdInfo.DeliverySetting.Budget,
					BudgetMode:              qianchuanAdInfo.DeliverySetting.BudgetMode,
					BudgetModeName:          qianchuanAdInfo.GetBudgetModeName(ctx),
					PayOrderCount:           qianchuanAdInfo.PayOrderCount,
					PayOrderAmount:          qianchuanAdInfo.PayOrderAmount,
					ClickCnt:                qianchuanAdInfo.ClickCnt,
					ShowCnt:                 qianchuanAdInfo.ShowCnt,
					ConvertCnt:              qianchuanAdInfo.ConvertCnt,
					ClickRate:               qianchuanAdInfo.ClickRate,
					CpmPlatform:             qianchuanAdInfo.CpmPlatform,
					DyFollow:                qianchuanAdInfo.DyFollow,
					PayConvertRate:          qianchuanAdInfo.PayConvertRate,
					ConvertCost:             qianchuanAdInfo.ConvertCost,
					ConvertRate:             qianchuanAdInfo.ConvertRate,
					AveragePayOrderStatCost: qianchuanAdInfo.AveragePayOrderStatCost,
					PayOrderAveragePrice:    qianchuanAdInfo.PayOrderAveragePrice,
					AdCreateTime:            adCreateTime,
					AdModifyTime:            adModifyTime,
				})
			}

			return true
		})
	}

	return &domain.ExternalQianchuanAdList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}, nil
}

func (qauc *QianchuanAdUsecase) StatisticsExternalQianchuanAds(ctx context.Context, startDay, endDay, keyword, advertiserIds, filter string) (*domain.StatisticsExternalQianchuanAds, error) {
	statistics := make([]*domain.StatisticsExternalQianchuanAd, 0)

	if statisticsAd, err := qauc.qairepo.Statistics(ctx, advertiserIds, startDay, endDay, keyword, filter); err == nil {
		statisticsAd.SetRoi(ctx)
		statisticsAd.SetClickRate(ctx)
		statisticsAd.SetCpmPlatform(ctx)
		statisticsAd.SetPayConvertRate(ctx)
		statisticsAd.SetConvertCost(ctx)
		statisticsAd.SetConvertRate(ctx)
		statisticsAd.SetAveragePayOrderStatCost(ctx)
		statisticsAd.SetPayOrderAveragePrice(ctx)

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalStatCost",
			Value: strconv.FormatFloat(statisticsAd.StatCost, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalRoi",
			Value: strconv.FormatFloat(statisticsAd.Roi, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalPayOrderAmount",
			Value: strconv.FormatFloat(statisticsAd.PayOrderAmount, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalPayOrderCount",
			Value: fmt.Sprintf("%d", statisticsAd.PayOrderCount),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalClickCnt",
			Value: fmt.Sprintf("%d", statisticsAd.ClickCnt),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalShowCnt",
			Value: fmt.Sprintf("%d", statisticsAd.ShowCnt),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalDyFollow",
			Value: fmt.Sprintf("%d", statisticsAd.DyFollow),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalConvertCnt",
			Value: fmt.Sprintf("%d", statisticsAd.ConvertCnt),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalClickRate",
			Value: strconv.FormatFloat(statisticsAd.ClickRate, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalCpmPlatform",
			Value: strconv.FormatFloat(statisticsAd.CpmPlatform, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalPayConvertRate",
			Value: strconv.FormatFloat(statisticsAd.PayConvertRate, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalConvertCost",
			Value: strconv.FormatFloat(statisticsAd.ConvertCost, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalConvertRate",
			Value: strconv.FormatFloat(statisticsAd.ConvertRate, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalAveragePayOrderStatCost",
			Value: strconv.FormatFloat(statisticsAd.AveragePayOrderStatCost, 'f', 2, 64),
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalPayOrderAveragePrice",
			Value: strconv.FormatFloat(statisticsAd.PayOrderAveragePrice, 'f', 2, 64),
		})
	} else {
		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalStatCost",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalRoi",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalPayOrderAmount",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalPayOrderCount",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalClickCnt",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalShowCnt",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalDyFollow",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalConvertCnt",
			Value: "0",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalClickRate",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalCpmPlatform",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalPayConvertRate",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalConvertCost",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalConvertRate",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalConvertRate",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalAveragePayOrderStatCost",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsExternalQianchuanAd{
			Key:   "TotalPayOrderAveragePrice",
			Value: "0.00",
		})
	}

	return &domain.StatisticsExternalQianchuanAds{
		Statistics: statistics,
	}, nil
}

func (qauc *QianchuanAdUsecase) ListSelectExternalQianchuanAds(ctx context.Context) (*domain.SelectExternalQianchuanAds, error) {
	return domain.NewSelectExternalQianchuanAds(), nil
}

func (qauc *QianchuanAdUsecase) SyncQianchuanAds(ctx context.Context, day string) error {
	if len(day) == 0 {
		day = tool.TimeToString("2006-01-02", time.Now())
	}

	qauc.qairepo.SaveIndex(ctx, day)

	var wg sync.WaitGroup

	wg.Add(1)

	go qauc.SyncQianchuanAdvertisers(ctx, &wg, day)

	wg.Wait()

	messageAd := domain.MessageAd{
		Type: "qianchuan_ad_data_ready",
	}
	messageAd.Message.Name = "douyin"
	messageAd.Message.SyncDate = day
	messageAd.Message.SendTime = tool.TimeToString("2006-01-02 15:04:05", time.Now())

	bmessageAd, _ := json.Marshal(messageAd)

	qauc.repo.Send(ctx, kafka.NewMessage(strconv.FormatUint(1, 10), bmessageAd))

	return nil
}

func (qauc *QianchuanAdUsecase) SyncQianchuanAdvertisers(ctx context.Context, wg *sync.WaitGroup, day string) error {
	defer wg.Done()

	var sqa sync.WaitGroup

	var qianchuanAdvertiserInfos []*domain.QianchuanAdvertiserInfo
	var qianchuanAdvertiserCampaigns []*domain.QianchuanAdvertiserInfo
	var qianchuanAdvertisers []*domain.QianchuanAdvertiser
	var qianchuanWallets []*domain.QianchuanWallet
	var qerr error
	var qcerr error
	var qaerr error
	var qwerr error

	sqa.Add(4)

	go func() {
		defer sqa.Done()

		qianchuanAdvertiserInfos, qerr = qauc.qairepo.ListAdvertiser(ctx, day)
	}()

	go func() {
		defer sqa.Done()

		qianchuanAdvertiserCampaigns, qcerr = qauc.qairepo.ListAdvertiserCampaigns(ctx, day)
	}()

	go func() {
		defer sqa.Done()

		qianchuanAdvertisers, qaerr = qauc.qarepo.List(ctx, 0, 48, 0, "", "", "")
	}()

	go func() {
		defer sqa.Done()

		qianchuanWallets, qwerr = qauc.qwrepo.List(ctx, day)
	}()

	sqa.Wait()

	if qerr == nil && qcerr == nil && qaerr == nil && qwerr == nil {
		for _, qianchuanAdvertiserInfo := range qianchuanAdvertiserInfos {
			var campaigns uint64 = 0
			name := ""
			var generalTotalBalance float64 = 0.00

			for _, qianchuanAdvertiserCampaign := range qianchuanAdvertiserCampaigns {
				if qianchuanAdvertiserCampaign.Id == qianchuanAdvertiserInfo.Id {
					campaigns = qianchuanAdvertiserCampaign.Campaigns

					break
				}
			}

			for _, qianchuanAdvertiser := range qianchuanAdvertisers {
				if qianchuanAdvertiser.AdvertiserId == qianchuanAdvertiserInfo.Id {
					name = qianchuanAdvertiser.AdvertiserName

					break
				}
			}

			for _, qianchuanWallet := range qianchuanWallets {
				if qianchuanWallet.AdvertiserId == qianchuanAdvertiserInfo.Id {
					generalTotalBalance = qianchuanWallet.GeneralTotalBalance

					break
				}
			}

			inQianchuanAdvertiserInfo := domain.NewQianchuanAdvertiserInfo(ctx, qianchuanAdvertiserInfo.Id, campaigns, generalTotalBalance, qianchuanAdvertiserInfo.StatCost, qianchuanAdvertiserInfo.PayOrderAmount, qianchuanAdvertiserInfo.CreateOrderAmount, qianchuanAdvertiserInfo.PayOrderCount, qianchuanAdvertiserInfo.CreateOrderCount, qianchuanAdvertiserInfo.ClickCnt, qianchuanAdvertiserInfo.ShowCnt, qianchuanAdvertiserInfo.ConvertCnt, qianchuanAdvertiserInfo.DyFollow, name)
			inQianchuanAdvertiserInfo.SetRoi(ctx)
			inQianchuanAdvertiserInfo.SetClickRate(ctx)
			inQianchuanAdvertiserInfo.SetCpmPlatform(ctx)
			inQianchuanAdvertiserInfo.SetPayConvertRate(ctx)
			inQianchuanAdvertiserInfo.SetConvertCost(ctx)
			inQianchuanAdvertiserInfo.SetConvertRate(ctx)
			inQianchuanAdvertiserInfo.SetAveragePayOrderStatCost(ctx)
			inQianchuanAdvertiserInfo.SetPayOrderAveragePrice(ctx)
			inQianchuanAdvertiserInfo.SetCreateTime(ctx)
			inQianchuanAdvertiserInfo.SetUpdateTime(ctx)

			if err := qauc.qainrepo.Upsert(ctx, day, inQianchuanAdvertiserInfo); err != nil {
				sinQianchuanAdvertiserInfo, _ := json.Marshal(inQianchuanAdvertiserInfo)

				inTaskLog := domain.NewTaskLog(ctx, "syncQianchuanAds", fmt.Sprintf("[SyncQianchuanAdsError SyncQianchuanCampaignError] AdvertiserId=%d, Data=%s, Description=%s", qianchuanAdvertiserInfo.Id, sinQianchuanAdvertiserInfo, "同步千川账户，插入数据库失败"))
				inTaskLog.SetCreateTime(ctx)

				qauc.tlrepo.Save(ctx, inTaskLog)
			}
		}
	}

	return nil
}
