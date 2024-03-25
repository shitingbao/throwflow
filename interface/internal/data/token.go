package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
	"interface/internal/biz"
)

type tokenRepo struct {
	data *Data
	log  *log.Helper
}

func NewTokenRepo(data *Data, logger log.Logger) biz.TokenRepo {
	return &tokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (tr *tokenRepo) Save(ctx context.Context, key string) (*v1.GetTokenReply, error) {
	token, err := tr.data.commonuc.GetToken(ctx, &v1.GetTokenRequest{
		Key: key,
	})

	if err != nil {
		return nil, err
	}

	return token, err
}

func (tr *tokenRepo) Verify(ctx context.Context, key string) (*v1.VerifyTokenReply, error) {
	token, err := tr.data.commonuc.VerifyToken(ctx, &v1.VerifyTokenRequest{
		Key: key,
	})

	if err != nil {
		return nil, err
	}

	return token, err
}
