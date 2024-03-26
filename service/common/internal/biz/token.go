package biz

import (
	"common/internal/conf"
	"common/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

var (
	CommonTokenCreateError = errors.InternalServer("COMMON_TOKEN_CREATE_ERROR", "TOKEN创建失败")
	CommonTokenVerifyError = errors.InternalServer("COMMON_TOKEN_VERIFY_ERROR", "TOKEN验证失败")
)

type TokenRepo interface {
	Save(context.Context, string, string, time.Duration) error
	Verify(context.Context, string) error
}

type TokenUsecase struct {
	repo TokenRepo
	conf *conf.Data
	log  *log.Helper
}

func NewTokenUsecase(repo TokenRepo, conf *conf.Data, logger log.Logger) *TokenUsecase {
	return &TokenUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (tuc *TokenUsecase) CreateToken(ctx context.Context, key string) (string, error) {
	token := tool.GetToken()

	if err := tuc.repo.Save(ctx, key+token, tool.TimeToString("2006-01-02 15:04:05", time.Now()), tuc.conf.Redis.PostTokenTimeout.AsDuration()); err != nil {
		return "", CommonTokenCreateError
	}

	return token, nil
}

func (tuc *TokenUsecase) VerifyToken(ctx context.Context, key string) error {
	if err := tuc.repo.Verify(ctx, key); err != nil {
		return CommonTokenVerifyError
	}

	return nil
}
