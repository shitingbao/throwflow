package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type loginRepo struct {
	data *Data
	log  *log.Helper
}

func NewLoginRepo(data *Data, logger log.Logger) biz.LoginRepo {
	return &loginRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (lr *loginRepo) Login(ctx context.Context, phone string) (*v1.LoginCompanyUserReply, error) {
	companyUser, err := lr.data.companyuc.LoginCompanyUser(ctx, &v1.LoginCompanyUserRequest{
		Phone: phone,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

func (lr *loginRepo) Logout(ctx context.Context, token string) (*v1.LogoutCompanyUserReply, error) {
	companyUser, err := lr.data.companyuc.LogoutCompanyUser(ctx, &v1.LogoutCompanyUserRequest{
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}
