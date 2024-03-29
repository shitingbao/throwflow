package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	DouyinQianchuanCampaignListError = errors.InternalServer("DOUYIN_QIANCHUAN_CAMPAIGN_LIST_ERROR", "千川广告组获取失败")
	DouyinQianchuanCampaignNotFound  = errors.NotFound("DOUYIN_QIANCHUAN_CAMPAIGN_NOT_FOUND", "千川广告组不存在")
)

type QianchuanCampaignRepo interface {
	GetByAdvertiserId(context.Context, uint64, uint64, string) (*domain.QianchuanCampaign, error)
	ListByAdvertiserId(context.Context, uint64, string, string) ([]*domain.QianchuanCampaign, error)
	List(context.Context, string, string, string, string, string, int64, int64) ([]*domain.QianchuanCampaign, error)
	All(context.Context, string) ([]*domain.QianchuanCampaign, error)
	CountByAdvertiserIds(context.Context, string, string) (map[uint64]int64, error)
	Count(context.Context, string, string, string, string, string) (int64, error)
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.QianchuanCampaign) error
}

type QianchuanCampaignUsecase struct {
	repo     QianchuanCampaignRepo
	qrarrepo QianchuanReportAdRealtimeRepo
	csrepo   CompanySetRepo
	qarepo   QianchuanAdRepo
	qrarepo  QianchuanReportAdRepo
	qadrepo  QianchuanAdvertiserRepo
	conf     *conf.Data
	log      *log.Helper
}

func NewQianchuanCampaignUsecase(repo QianchuanCampaignRepo, qrarrepo QianchuanReportAdRealtimeRepo, csrepo CompanySetRepo, qarepo QianchuanAdRepo, qrarepo QianchuanReportAdRepo, qadrepo QianchuanAdvertiserRepo, conf *conf.Data, logger log.Logger) *QianchuanCampaignUsecase {
	return &QianchuanCampaignUsecase{repo: repo, qrarrepo: qrarrepo, csrepo: csrepo, qarepo: qarepo, qrarepo: qrarepo, qadrepo: qadrepo, conf: conf, log: log.NewHelper(logger)}
}

func (qcuc *QianchuanCampaignUsecase) ListQianchuanCampaigns(ctx context.Context, pageNum, pageSize uint64, day, keyword, advertiserIds string) (*domain.QianchuanCampaignList, error) {
	list, err := qcuc.repo.List(ctx, advertiserIds, day, keyword, "", "", int64(pageNum), int64(pageSize))

	if err != nil {
		return nil, DouyinQianchuanCampaignListError
	}

	total, err := qcuc.repo.Count(ctx, advertiserIds, day, keyword, "", "")

	if err != nil {
		return nil, DouyinQianchuanCampaignListError
	}

	return &domain.QianchuanCampaignList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}
