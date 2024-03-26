package biz

import (
	"common/internal/conf"
	"common/internal/domain"
	"common/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type ShortCodeUsecase struct {
	sclrepo ShortCodeLogRepo
	dconf   *conf.Data
	log     *log.Helper
}

func NewShortCodeUsecase(sclrepo ShortCodeLogRepo, dconf *conf.Data, logger log.Logger) *ShortCodeUsecase {
	return &ShortCodeUsecase{sclrepo: sclrepo, dconf: dconf, log: log.NewHelper(logger)}
}

func (scuc *ShortCodeUsecase) CreateShortCode(ctx context.Context) (string, error) {
	shortCode := ""
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		shortCode, err = tool.GetShortCode()

		if err == nil {
			if _, err = scuc.sclrepo.Get(ctx, shortCode); err != nil {
				break
			} else {
				shortCode = ""
			}
		}
	}

	if len(shortCode) == 0 {
		return "", CommonShortCodeCreateError
	}

	inShortCodeLog := domain.NewShortCodeLog(ctx, shortCode)
	inShortCodeLog.SetCreateTime(ctx)
	inShortCodeLog.SetUpdateTime(ctx)

	if _, err := scuc.sclrepo.Save(ctx, inShortCodeLog); err != nil {
		return "", CommonShortCodeCreateError
	}

	return shortCode, nil

}
