package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/weixin/v1"
	"task/internal/biz"
)

type weixinUserOrganizationRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserOrganizationRepo(data *Data, logger log.Logger) biz.WeixinUserOrganizationRepo {
	return &weixinUserOrganizationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (wuor *weixinUserOrganizationRepo) Sync(ctx context.Context) (*v1.SyncUserOrganizationRelationsReply, error) {
	userOrganization, err := wuor.data.weixinuc.SyncUserOrganizationRelations(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return userOrganization, err
}
