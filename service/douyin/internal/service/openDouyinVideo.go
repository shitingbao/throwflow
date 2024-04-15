package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/domain"
	"math"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

func (ds *DouyinService) ListOpenDouyinVideos(ctx context.Context, in *v1.ListOpenDouyinVideosRequest) (*v1.ListOpenDouyinVideosReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	openDouyinVideos, err := ds.odvuc.ListOpenDouyinVideos(ctx, in.PageNum, in.PageSize, uint8(in.IsExistProduct), in.VideoIds, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListOpenDouyinVideosReply_OpenDouyinVideo, 0)

	for _, openDouyinVideo := range openDouyinVideos.List {
		list = append(list, &v1.ListOpenDouyinVideosReply_OpenDouyinVideo{
			VideoId:     openDouyinVideo.VideoId,
			MediaType:   uint64(openDouyinVideo.MediaType),
			Title:       openDouyinVideo.Title,
			Cover:       openDouyinVideo.Cover,
			CreateTime:  time.Unix(openDouyinVideo.CreateTime, 0).Format("2006-01-02 15:04"),
			AwemeId:     openDouyinVideo.AwemeId,
			AccountId:   openDouyinVideo.AccountId,
			Nickname:    openDouyinVideo.Nickname,
			Avatar:      openDouyinVideo.Avatar,
			ProductId:   openDouyinVideo.ProductId,
			ProductName: openDouyinVideo.ProductName,
			ProductImg:  openDouyinVideo.ProductImg,
		})
	}

	totalPage := uint64(math.Ceil(float64(openDouyinVideos.Total) / float64(openDouyinVideos.PageSize)))

	return &v1.ListOpenDouyinVideosReply{
		Code: 200,
		Data: &v1.ListOpenDouyinVideosReply_Data{
			PageNum:   openDouyinVideos.PageNum,
			PageSize:  openDouyinVideos.PageSize,
			Total:     openDouyinVideos.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListProductOpenDouyinVideos(ctx context.Context, in *v1.ListProductOpenDouyinVideosRequest) (*v1.ListProductOpenDouyinVideosReply, error) {
	productOpenDouyinVideos, err := ds.odvuc.ListProductOpenDouyinVideos(ctx, in.PageNum, in.PageSize, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListProductOpenDouyinVideosReply_ProductOpenDouyinVideo, 0)

	for _, productOpenDouyinVideo := range productOpenDouyinVideos.List {
		list = append(list, &v1.ListProductOpenDouyinVideosReply_ProductOpenDouyinVideo{
			ProductId:   productOpenDouyinVideo.ProductId,
			ProductName: productOpenDouyinVideo.ProductName,
			ProductImg:  productOpenDouyinVideo.ProductImg,
		})
	}

	totalPage := uint64(math.Ceil(float64(productOpenDouyinVideos.Total) / float64(productOpenDouyinVideos.PageSize)))

	return &v1.ListProductOpenDouyinVideosReply{
		Code: 200,
		Data: &v1.ListProductOpenDouyinVideosReply_Data{
			PageNum:   productOpenDouyinVideos.PageNum,
			PageSize:  productOpenDouyinVideos.PageSize,
			Total:     productOpenDouyinVideos.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListVideoIdOpenDouyinVideos(ctx context.Context, in *v1.ListVideoIdOpenDouyinVideosRequest) (*v1.ListVideoIdOpenDouyinVideosReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	openDouyinVideos, err := ds.odvuc.ListVideoIdOpenDouyinVideos(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListVideoIdOpenDouyinVideosReply_OpenDouyinVideo, 0)

	for _, openDouyinVideo := range openDouyinVideos.List {
		list = append(list, &v1.ListVideoIdOpenDouyinVideosReply_OpenDouyinVideo{
			VideoId: openDouyinVideo.VideoId,
		})
	}

	totalPage := uint64(math.Ceil(float64(openDouyinVideos.Total) / float64(openDouyinVideos.PageSize)))

	return &v1.ListVideoIdOpenDouyinVideosReply{
		Code: 200,
		Data: &v1.ListVideoIdOpenDouyinVideosReply_Data{
			PageNum:   openDouyinVideos.PageNum,
			PageSize:  openDouyinVideos.PageSize,
			Total:     openDouyinVideos.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListVideoIdOpenDouyinVideoByClientKeyAndOpenIds(ctx context.Context, in *v1.ListVideoIdOpenDouyinVideoByClientKeyAndOpenIdsRequest) (*v1.ListVideoIdOpenDouyinVideoByClientKeyAndOpenIdsReply, error) {
	openDouyinVideos, err := ds.odvuc.ListVideoIdOpenDouyinVideoByClientKeyAndOpenIds(ctx, in.PageNum, in.PageSize, in.ClientKey, in.OpenId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListVideoIdOpenDouyinVideoByClientKeyAndOpenIdsReply_OpenDouyinVideo, 0)

	for _, openDouyinVideo := range openDouyinVideos.List {
		list = append(list, &v1.ListVideoIdOpenDouyinVideoByClientKeyAndOpenIdsReply_OpenDouyinVideo{
			VideoId: openDouyinVideo.VideoId,
		})
	}

	totalPage := uint64(math.Ceil(float64(openDouyinVideos.Total) / float64(openDouyinVideos.PageSize)))

	return &v1.ListVideoIdOpenDouyinVideoByClientKeyAndOpenIdsReply{
		Code: 200,
		Data: &v1.ListVideoIdOpenDouyinVideoByClientKeyAndOpenIdsReply_Data{
			PageNum:   openDouyinVideos.PageNum,
			PageSize:  openDouyinVideos.PageSize,
			Total:     openDouyinVideos.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListVideoIdOpenDouyinVideoByProductIds(ctx context.Context, in *v1.ListVideoIdOpenDouyinVideoByProductIdsRequest) (*v1.ListVideoIdOpenDouyinVideoByProductIdsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	openDouyinVideos, err := ds.odvuc.ListVideoIdOpenDouyinVideoByProductIds(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListVideoIdOpenDouyinVideoByProductIdsReply_OpenDouyinVideo, 0)

	for _, openDouyinVideo := range openDouyinVideos.List {
		list = append(list, &v1.ListVideoIdOpenDouyinVideoByProductIdsReply_OpenDouyinVideo{
			VideoId: openDouyinVideo.VideoId,
		})
	}

	totalPage := uint64(math.Ceil(float64(openDouyinVideos.Total) / float64(openDouyinVideos.PageSize)))

	return &v1.ListVideoIdOpenDouyinVideoByProductIdsReply{
		Code: 200,
		Data: &v1.ListVideoIdOpenDouyinVideoByProductIdsReply_Data{
			PageNum:   openDouyinVideos.PageNum,
			PageSize:  openDouyinVideos.PageSize,
			Total:     openDouyinVideos.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListVideoTokensOpenDouyinVideos(ctx context.Context, in *v1.ListVideoTokensOpenDouyinVideosRequest) (*v1.ListVideoTokensOpenDouyinVideosReply, error) {
	tokens := []*domain.OpenDouyinToken{}

	for _, v := range in.Tokens {
		tokens = append(tokens, &domain.OpenDouyinToken{
			ClientKey: v.ClientKey,
			OpenId:    v.OpenId,
		})
	}

	openDouyinVideos, err := ds.odvuc.ListProductsByTokens(ctx, in.ProductOutId, in.ClaimTime, tokens)

	if err != nil {
		return nil, err
	}

	list := []*v1.ListVideoTokensOpenDouyinVideosReply_OpenDouyinVideo{}

	for _, openDouyinVideo := range openDouyinVideos.List {
		list = append(list, &v1.ListVideoTokensOpenDouyinVideosReply_OpenDouyinVideo{
			VideoId:     openDouyinVideo.VideoId,
			MediaType:   uint64(openDouyinVideo.MediaType),
			Title:       openDouyinVideo.Title,
			Cover:       openDouyinVideo.Cover,
			CreateTime:  time.Unix(openDouyinVideo.CreateTime, 0).Format("2006-01-02 15:04:05"),
			AwemeId:     openDouyinVideo.AwemeId,
			AccountId:   openDouyinVideo.AccountId,
			Nickname:    openDouyinVideo.Nickname,
			Avatar:      openDouyinVideo.Avatar,
			ProductId:   openDouyinVideo.ProductId,
			ProductName: openDouyinVideo.ProductName,
			ProductImg:  openDouyinVideo.ProductImg,
			Statistics: &v1.ListVideoTokensOpenDouyinVideosReply_VideoStatistics{
				CommentCount:  uint32(openDouyinVideo.Statistics.CommentCount),
				DiggCount:     uint32(openDouyinVideo.Statistics.DiggCount),
				DownloadCount: uint32(openDouyinVideo.Statistics.DownloadCount),
				ForwardCount:  uint32(openDouyinVideo.Statistics.ForwardCount),
				PlayCount:     uint32(openDouyinVideo.Statistics.PlayCount),
				ShareCount:    uint32(openDouyinVideo.Statistics.ShareCount),
			},
			ClientKey:   openDouyinVideo.ClientKey,
			OpenId:      openDouyinVideo.OpenId,
			ItemId:      openDouyinVideo.ItemId,
			IsTop:       openDouyinVideo.IsTop,
			ShareUrl:    openDouyinVideo.ShareUrl,
			VideoStatus: uint32(openDouyinVideo.VideoStatus),
		})
	}

	return &v1.ListVideoTokensOpenDouyinVideosReply{
		Code: 200,
		Data: &v1.ListVideoTokensOpenDouyinVideosReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) UpdateIsUpdateCoverAndProductIdOpenDouyinVideos(ctx context.Context, in *v1.UpdateIsUpdateCoverAndProductIdOpenDouyinVideosRequest) (*v1.UpdateIsUpdateCoverAndProductIdOpenDouyinVideosReply, error) {
	if err := ds.odvuc.UpdateIsUpdateCoverAndProductIdOpenDouyinVideos(ctx, uint8(in.VideoStatus), in.VideoId, in.Cover, in.ProductId); err != nil {
		return nil, err
	}

	return &v1.UpdateIsUpdateCoverAndProductIdOpenDouyinVideosReply{
		Code: 200,
		Data: &v1.UpdateIsUpdateCoverAndProductIdOpenDouyinVideosReply_Data{},
	}, nil
}

func (ds *DouyinService) SyncOpenDouyinVideos(ctx context.Context, in *empty.Empty) (*v1.SyncOpenDouyinVideosReply, error) {
	ds.log.Infof("同步精选联盟达人视频数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.odvuc.SyncOpenDouyinVideos(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("同步精选联盟达人视频数据, 结束时间 %s \n", time.Now())

	return &v1.SyncOpenDouyinVideosReply{
		Code: 200,
		Data: &v1.SyncOpenDouyinVideosReply_Data{},
	}, nil
}
