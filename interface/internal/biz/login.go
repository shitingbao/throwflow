package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/pkg/tool"
)

type LoginRepo interface {
	Login(context.Context, string) (*v1.LoginCompanyUserReply, error)
	Logout(context.Context, string) (*v1.LogoutCompanyUserReply, error)
}

type LoginUsecase struct {
	repo LoginRepo
	log  *log.Helper
}

func NewLoginUsecase(repo LoginRepo, logger log.Logger) *LoginUsecase {
	return &LoginUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (luc *LoginUsecase) Login(ctx context.Context, phone string) (*v1.LoginCompanyUserReply, error) {
	companyUser, err := luc.repo.Login(ctx, phone)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LOGIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companyUser, nil
}

func (luc *LoginUsecase) Logout(ctx context.Context) error {
	token := ctx.Value("token")

	if _, err := luc.repo.Logout(ctx, token.(string)); err != nil {
		return errors.InternalServer("INTERFACE_LOGIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return nil
}
