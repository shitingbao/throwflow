package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type oceanengineAccountTokenRepo struct {
	data *Data
	log  *log.Helper
}

func NewOceanengineAccountTokenRepo(data *Data, logger log.Logger) biz.OceanengineAccountTokenRepo {
	return &oceanengineAccountTokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (oatr *oceanengineAccountTokenRepo) Refresh(ctx context.Context) (*v1.RefreshOceanengineAccountTokensReply, error) {
	token, err := oatr.data.douyinuc.RefreshOceanengineAccountTokens(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return token, err
}
