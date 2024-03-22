package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type doukeTokenRepo struct {
	data *Data
	log  *log.Helper
}

func NewDoukeTokenRepo(data *Data, logger log.Logger) biz.DoukeTokenRepo {
	return &doukeTokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (dtr *doukeTokenRepo) Refresh(ctx context.Context) (*v1.RefreshDoukeTokensReply, error) {
	token, err := dtr.data.douyinuc.RefreshDoukeTokens(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return token, err
}
