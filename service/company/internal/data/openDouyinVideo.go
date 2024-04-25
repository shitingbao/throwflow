package data

import (
	douyinv1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type openDouyinVideoRepo struct {
	data *Data
	log  *log.Helper
}

func NewOpenDouyinVideoRepo(data *Data, logger log.Logger) biz.OpenDouyinVideoRepo {
	return &openDouyinVideoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ctr *openDouyinVideoRepo) ListVideoTokensOpenDouyinVideos(ctx context.Context, productOutId uint64, claimTime time.Time, tokens []*domain.CompanyTaskClientKeyAndOpenId) ([]*douyinv1.ListVideoTokensOpenDouyinVideosReply_OpenDouyinVideo, error) {
	tks := []*douyinv1.ListVideoTokensOpenDouyinVideosRequestToken{}

	for _, v := range tokens {
		tks = append(tks, &douyinv1.ListVideoTokensOpenDouyinVideosRequestToken{
			ClientKey: v.ClientKey,
			OpenId:    v.OpenId,
		})
	}

	res, err := ctr.data.douyinuc.ListVideoTokensOpenDouyinVideos(ctx, &douyinv1.ListVideoTokensOpenDouyinVideosRequest{
		ProductOutId: productOutId,
		Tokens:       tks,
		ClaimTime:    tool.TimeToString("2006-01-02 15:04:05", claimTime),
	})

	if err != nil {
		return nil, err
	}

	return res.Data.List, nil
}
