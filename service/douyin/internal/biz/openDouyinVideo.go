package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/openDouyin/video"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/time/rate"
)

var (
	DouyinOpenDouyinVideoListError = errors.InternalServer("DOUYIN_OPEN_DOUYIN_VIDEO_LIST_ERROR", "达人视频列表获取失败")
)

type OpenDouyinVideoRepo interface {
	GetByClientKeyAndOpenId(context.Context, string, string, string, int32) (*domain.OpenDouyinVideo, error)
	List(context.Context, uint8, string, string, string, string, []*domain.OpenDouyinToken, int64, int64) ([]*domain.OpenDouyinVideo, error)
	ListProduct(context.Context, string, []*domain.OpenDouyinToken, int64, int64) ([]*domain.OpenDouyinVideo, error)
	ListVideoId(context.Context, int64, int64) ([]*domain.OpenDouyinVideo, error)
	ListProductsByTokens(context.Context, uint64, string, []*domain.OpenDouyinToken) ([]*domain.OpenDouyinVideo, error)
	Count(context.Context, uint8, string, string, string, string, []*domain.OpenDouyinToken) (int64, error)
	CountProduct(context.Context, string, []*domain.OpenDouyinToken) (int64, error)
	CountVideoId(context.Context) (int64, error)
	SaveIndex(context.Context)
	Upsert(context.Context, *domain.OpenDouyinVideo) error
	UpdateIsUpdateCoverAndProductId(context.Context, string, string, string) error
	UpdateVideoStatus(context.Context, int32, []string) error
}

type OpenDouyinVideoUsecase struct {
	repo     OpenDouyinVideoRepo
	odtrepo  OpenDouyinTokenRepo
	oduirepo OpenDouyinUserInfoRepo
	tlrepo   TaskLogRepo
	odalrepo OpenDouyinApiLogRepo
	joirepo  JinritemaiOrderInfoRepo
	wuodrepo WeixinUserOpenDouyinRepo
	conf     *conf.Data
	log      *log.Helper
}

func NewOpenDouyinVideoUsecase(repo OpenDouyinVideoRepo, odtrepo OpenDouyinTokenRepo, oduirepo OpenDouyinUserInfoRepo, tlrepo TaskLogRepo, odalrepo OpenDouyinApiLogRepo, joirepo JinritemaiOrderInfoRepo, wuodrepo WeixinUserOpenDouyinRepo, conf *conf.Data, logger log.Logger) *OpenDouyinVideoUsecase {
	return &OpenDouyinVideoUsecase{repo: repo, odtrepo: odtrepo, oduirepo: oduirepo, tlrepo: tlrepo, odalrepo: odalrepo, joirepo: joirepo, wuodrepo: wuodrepo, conf: conf, log: log.NewHelper(logger)}
}

func (odvuc *OpenDouyinVideoUsecase) ListOpenDouyinVideos(ctx context.Context, pageNum, pageSize uint64, isExistProduct uint8, videoIds, keyword string) (*domain.OpenDouyinVideoList, error) {
	openDouyinTokens := make([]*domain.OpenDouyinToken, 0)

	list, err := odvuc.repo.List(ctx, isExistProduct, videoIds, keyword, "1", "", openDouyinTokens, int64(pageNum), int64(pageSize))

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	total, err := odvuc.repo.Count(ctx, isExistProduct, videoIds, keyword, "1", "", openDouyinTokens)

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	return &domain.OpenDouyinVideoList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (odvuc *OpenDouyinVideoUsecase) ListProductOpenDouyinVideos(ctx context.Context, pageNum, pageSize uint64, keyword string) (*domain.OpenDouyinVideoList, error) {
	openDouyinTokens := make([]*domain.OpenDouyinToken, 0)

	list, err := odvuc.repo.ListProduct(ctx, keyword, openDouyinTokens, int64(pageNum), int64(pageSize))

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	total, err := odvuc.repo.CountProduct(ctx, keyword, openDouyinTokens)

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	return &domain.OpenDouyinVideoList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (odvuc *OpenDouyinVideoUsecase) ListVideoIdOpenDouyinVideos(ctx context.Context, pageNum, pageSize uint64) (*domain.OpenDouyinVideoList, error) {
	list, err := odvuc.repo.ListVideoId(ctx, int64(pageNum), int64(pageSize))

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	total, err := odvuc.repo.CountVideoId(ctx)

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	return &domain.OpenDouyinVideoList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (odvuc *OpenDouyinVideoUsecase) ListVideoIdOpenDouyinVideoByClientKeyAndOpenIds(ctx context.Context, pageNum, pageSize uint64, clientKey, openId string) (*domain.OpenDouyinVideoList, error) {
	openDouyinTokens := make([]*domain.OpenDouyinToken, 0)
	openDouyinTokens = append(openDouyinTokens, &domain.OpenDouyinToken{
		ClientKey: clientKey,
		OpenId:    openId,
	})

	list, err := odvuc.repo.List(ctx, 0, "", "", "1", "2,4", openDouyinTokens, int64(pageNum), int64(pageSize))

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	total, err := odvuc.repo.Count(ctx, 0, "", "", "1", "2,4", openDouyinTokens)

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	return &domain.OpenDouyinVideoList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (odvuc *OpenDouyinVideoUsecase) UpdateIsUpdateCoverAndProductIdOpenDouyinVideos(ctx context.Context, videoStatus uint8, videoId, cover, productId string) error {
	if videoStatus == 99 {
		if err := odvuc.repo.UpdateVideoStatus(ctx, int32(videoStatus), []string{videoId}); err != nil {
			return DouyinOpenDouyinUserInfoUpdateError
		}
	} else {
		if err := odvuc.repo.UpdateIsUpdateCoverAndProductId(ctx, videoId, cover, productId); err != nil {
			return DouyinOpenDouyinUserInfoUpdateError
		}
	}

	return nil
}

func (odvuc *OpenDouyinVideoUsecase) SyncOpenDouyinVideos(ctx context.Context) error {
	openDouyinTokens, err := odvuc.odtrepo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncOpenDouyinVideos", fmt.Sprintf("[SyncOpenDouyinVideosError] Description=%s", "获取抖音开放平台token列表失败"))
		inTaskLog.SetCreateTime(ctx)

		odvuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	odvuc.repo.SaveIndex(ctx)

	var wg sync.WaitGroup

	limiter := rate.NewLimiter(0, 60)
	limiter.SetLimit(rate.Limit(60))

	for _, openDouyinToken := range openDouyinTokens {
		wg.Add(1)

		openDouyinUserInfo, _ := odvuc.oduirepo.GetByClientKeyAndOpenId(ctx, openDouyinToken.ClientKey, openDouyinToken.OpenId)
		jinritemaiOrderInfos, _ := odvuc.joirepo.ListProductByClientKeyAndOpenIdAndMediaType(ctx, openDouyinToken.ClientKey, openDouyinToken.OpenId, "video")

		go odvuc.SyncOpenDouyinVideo(ctx, &wg, limiter, openDouyinToken, openDouyinUserInfo, jinritemaiOrderInfos)
	}

	wg.Wait()

	return nil
}

func (odvuc *OpenDouyinVideoUsecase) SyncOpenDouyinVideo(ctx context.Context, wg *sync.WaitGroup, limiter *rate.Limiter, openDouyinToken *domain.OpenDouyinToken, openDouyinUserInfo *domain.OpenDouyinUserInfo, jinritemaiOrderInfos []*domain.JinritemaiOrderInfo) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("获取抖音开放平台达人视频列表异常，", err)
		}

		wg.Done()
	}()

	videoIds := make([]string, 0)

	videos, err := odvuc.listVideos(ctx, limiter, 0, openDouyinToken, openDouyinUserInfo, jinritemaiOrderInfos)

	if err == nil {
		for _, video := range videos.Data.List {
			videoIds = append(videoIds, video.VideoId)
		}

		for videos.Data.HasMore {
			videos, err = odvuc.listVideos(ctx, limiter, videos.Data.Cursor, openDouyinToken, openDouyinUserInfo, jinritemaiOrderInfos)

			if err != nil {
				continue
			} else {
				for _, video := range videos.Data.List {
					videoIds = append(videoIds, video.VideoId)
				}
			}
		}
	}

	odvuc.repo.UpdateVideoStatus(ctx, 99, videoIds)
}

func (odvuc *OpenDouyinVideoUsecase) listVideos(ctx context.Context, limiter *rate.Limiter, cursor int64, openDouyinToken *domain.OpenDouyinToken, openDouyinUserInfo *domain.OpenDouyinUserInfo, jinritemaiOrderInfos []*domain.JinritemaiOrderInfo) (*video.ListVideoResponse, error) {
	var videos *video.ListVideoResponse
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		limiter.Wait(ctx)

		videos, err = video.ListVideo(openDouyinToken.AccessToken, openDouyinToken.OpenId, cursor)

		if err != nil {
			if retryNum == 2 {
				inOpenDouyinApiLog := domain.NewOpenDouyinApiLog(ctx, openDouyinToken.ClientKey, openDouyinToken.OpenId, openDouyinToken.AccessToken, err.Error())
				inOpenDouyinApiLog.SetCreateTime(ctx)

				odvuc.odalrepo.Save(ctx, inOpenDouyinApiLog)
			}
		} else {
			for _, video := range videos.Data.List {
				statistics := domain.VideoStatistics{
					CommentCount:  video.Statistics.CommentCount,
					DiggCount:     video.Statistics.DiggCount,
					DownloadCount: video.Statistics.DownloadCount,
					ForwardCount:  video.Statistics.ForwardCount,
					PlayCount:     video.Statistics.PlayCount,
					ShareCount:    video.Statistics.ShareCount,
				}

				var inOpenDouyinVideo *domain.OpenDouyinVideo

				if openDouyinUserInfo == nil {
					inOpenDouyinVideo = domain.NewOpenDouyinVideo(ctx, openDouyinToken.ClientKey, openDouyinToken.OpenId, "", "", "", video.Title, video.Cover, video.ItemId, video.ShareUrl, video.VideoId, 0, video.CreateTime, video.MediaType, video.VideoStatus, video.IsReviewed, video.IsTop, statistics)
				} else {
					inOpenDouyinVideo = domain.NewOpenDouyinVideo(ctx, openDouyinToken.ClientKey, openDouyinToken.OpenId, openDouyinUserInfo.AccountId, openDouyinUserInfo.Nickname, openDouyinUserInfo.Avatar, video.Title, video.Cover, video.ItemId, video.ShareUrl, video.VideoId, openDouyinUserInfo.AwemeId, video.CreateTime, video.MediaType, video.VideoStatus, video.IsReviewed, video.IsTop, statistics)
				}

				for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
					if strconv.FormatUint(jinritemaiOrderInfo.MediaId, 10) == video.VideoId {
						inOpenDouyinVideo.SetProductId(ctx, jinritemaiOrderInfo.ProductId)
						inOpenDouyinVideo.SetProductName(ctx, jinritemaiOrderInfo.ProductName)
						inOpenDouyinVideo.SetProductImg(ctx, jinritemaiOrderInfo.ProductImg)
					}
				}

				if err := odvuc.repo.Upsert(ctx, inOpenDouyinVideo); err != nil {
					sinOpenDouyinVideo, _ := json.Marshal(inOpenDouyinVideo)

					inTaskLog := domain.NewTaskLog(ctx, "SyncOpenDouyinVideos", fmt.Sprintf("[SyncOpenDouyinVideosError SyncOpenDouyinVideoError] ClientKey=%d, OpenId=%d, Data=%s, Description=%s", openDouyinToken.ClientKey, openDouyinToken.OpenId, sinOpenDouyinVideo, "同步精选联盟达人视频数据，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					odvuc.tlrepo.Save(ctx, inTaskLog)
				}
			}

			break
		}
	}

	return videos, err
}

func (odvuc *OpenDouyinVideoUsecase) ListProductsByTokens(ctx context.Context, productOutId uint64, claimTime string, tokens []*domain.OpenDouyinToken) (*domain.OpenDouyinVideoList, error) {
	list, err := odvuc.repo.ListProductsByTokens(ctx, productOutId, claimTime, tokens)

	if err != nil {
		return nil, DouyinOpenDouyinVideoListError
	}

	return &domain.OpenDouyinVideoList{
		List: list,
	}, nil
}
