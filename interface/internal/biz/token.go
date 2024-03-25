package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
)

var (
	InterfaceTokenCreateError = errors.InternalServer("INTERFACE_TOKEN_CREATE_ERROR", "TOKEN创建失败")
	InterfaceTokenVerifyError = errors.InternalServer("INTERFACE_TOKEN_VERIFY_ERROR", "TOKEN验证失败")
)

type TokenRepo interface {
	Save(context.Context, string) (*v1.GetTokenReply, error)
	Verify(context.Context, string) (*v1.VerifyTokenReply, error)
}

type TokenUsecase struct {
	repo TokenRepo
	log  *log.Helper
}

func NewTokenUsecase(repo TokenRepo, logger log.Logger) *TokenUsecase {
	return &TokenUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (tuc *TokenUsecase) CreateToken(ctx context.Context) (string, error) {
	token, err := tuc.repo.Save(ctx, "interface:token:")

	if err != nil {
		return "", InterfaceTokenCreateError
	}

	return token.Data.Token, nil
}

func (tuc *TokenUsecase) VerifyToken(ctx context.Context, token string) bool {
	if _, err := tuc.repo.Verify(ctx, "interface:token:"+token); err != nil {
		return false
	}

	return true
}
