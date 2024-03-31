package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/common/v1"
	"weixin/internal/biz"
)

type shortUrlRepo struct {
	data *Data
	log  *log.Helper
}

func NewShortUrlRepo(data *Data, logger log.Logger) biz.ShortUrlRepo {
	return &shortUrlRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (sur *shortUrlRepo) Create(ctx context.Context, content string) (*v1.CreateShortUrlReply, error) {
	shortUrl, err := sur.data.commonuc.CreateShortUrl(ctx, &v1.CreateShortUrlRequest{
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return shortUrl, err
}
