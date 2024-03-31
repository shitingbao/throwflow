package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/douyin/v1"
	"weixin/internal/biz"
)

type openDouyinUserInfoRepo struct {
	data *Data
	log  *log.Helper
}

func NewOpenDouyinUserInfoRepo(data *Data, logger log.Logger) biz.OpenDouyinUserInfoRepo {
	return &openDouyinUserInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (oduir *openDouyinUserInfoRepo) UpdateCooperativeCodes(ctx context.Context, clientKey, openId, cooperativeCode string) (*v1.UpdateCooperativeCodeDouyinTokensReply, error) {
	douyinTokens, err := oduir.data.douyinuc.UpdateCooperativeCodeDouyinTokens(ctx, &v1.UpdateCooperativeCodeDouyinTokensRequest{
		ClientKey:       clientKey,
		OpenId:          openId,
		CooperativeCode: cooperativeCode,
	})

	if err != nil {
		return nil, err
	}

	return douyinTokens, err
}
