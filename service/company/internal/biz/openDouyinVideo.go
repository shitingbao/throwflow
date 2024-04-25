package biz

import (
	douyinv1 "company/api/service/douyin/v1"
	"company/internal/domain"
	"context"
	"time"
)

type OpenDouyinVideoRepo interface {
	ListVideoTokensOpenDouyinVideos(context.Context, uint64, time.Time, []*domain.CompanyTaskClientKeyAndOpenId) ([]*douyinv1.ListVideoTokensOpenDouyinVideosReply_OpenDouyinVideo, error)
}
