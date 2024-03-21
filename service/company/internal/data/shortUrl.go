package data

import (
	v1 "company/api/service/common/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
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
