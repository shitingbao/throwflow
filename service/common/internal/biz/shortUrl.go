package biz

import (
	"common/internal/conf"
	"common/internal/pkg/suolink/suolink"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	CommonShortUrlCreateError = errors.InternalServer("COMMON_SHORT_URL_CREATE_ERROR", "短链接创建失败")
)

type ShortUrlUsecase struct {
	conf  *conf.Data
	sconf *conf.Suolink
	log   *log.Helper
}

func NewShortUrlUsecase(conf *conf.Data, sconf *conf.Suolink, logger log.Logger) *ShortUrlUsecase {
	return &ShortUrlUsecase{conf: conf, sconf: sconf, log: log.NewHelper(logger)}
}

func (suuc *ShortUrlUsecase) CreateShortUrl(ctx context.Context, content string) (*suolink.GetSuoLinkResponse, error) {
	shortUrl, err := suolink.GetSuoLink(suuc.sconf.Key, suuc.sconf.Url+"?"+content)

	if err != nil {
		return nil, err
	}

	return shortUrl, nil
}