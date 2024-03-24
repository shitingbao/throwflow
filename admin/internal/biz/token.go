package biz

import (
	v1 "admin/api/service/common/v1"
	"admin/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	AdminTokenCreateError = errors.InternalServer("ADMIN_TOKEN_CREATE_ERROR", "TOKEN创建失败")
	AdminTokenVerifyError = errors.InternalServer("ADMIN_TOKEN_VERIFY_ERROR", "TOKEN验证失败")
)

type TokenRepo interface {
	Save(context.Context, string) (*v1.GetTokenReply, error)
	Verify(context.Context, string) (*v1.VerifyTokenReply, error)
}

type TokenUsecase struct {
	repo TokenRepo
	conf *conf.Data
	log  *log.Helper
}

func NewTokenUsecase(repo TokenRepo, conf *conf.Data, logger log.Logger) *TokenUsecase {
	return &TokenUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (tuc *TokenUsecase) CreateToken(ctx context.Context) (string, error) {
	token, err := tuc.repo.Save(ctx, "admin:token:")

	if err != nil {
		return "", AdminTokenCreateError
	}

	return token.Data.Token, nil
}

func (tuc *TokenUsecase) VerifyToken(ctx context.Context, token string) bool {
	if _, err := tuc.repo.Verify(ctx, "admin:token:"+token); err != nil {
		return false
	}

	return true
}
