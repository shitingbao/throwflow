package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "weixin/api/service/common/v1"
	"weixin/internal/biz"
)

type shortCodeRepo struct {
	data *Data
	log  *log.Helper
}

func NewShortCodeRepo(data *Data, logger log.Logger) biz.ShortCodeRepo {
	return &shortCodeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (scr *shortCodeRepo) Create(ctx context.Context) (*v1.CreateShortCodeReply, error) {
	shortCode, err := scr.data.commonuc.CreateShortCode(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return shortCode, err
}
