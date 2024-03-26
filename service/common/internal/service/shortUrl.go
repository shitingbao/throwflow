package service

import (
	v1 "common/api/common/v1"
	"context"
)

func (cs *CommonService) CreateShortUrl(ctx context.Context, in *v1.CreateShortUrlRequest) (*v1.CreateShortUrlReply, error) {
	shortUrl, err := cs.suuc.CreateShortUrl(ctx, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.CreateShortUrlReply{
		Code: 200,
		Data: &v1.CreateShortUrlReply_Data{
			ShortUrl: shortUrl.Url,
		},
	}, nil
}
