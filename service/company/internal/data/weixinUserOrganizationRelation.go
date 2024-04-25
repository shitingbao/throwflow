package data

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type weixinUserOrganizationRelationRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserOrganizationRelationRepo(data *Data, logger log.Logger) biz.WeixinUserOrganizationRelationRepo {
	return &weixinUserOrganizationRelationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ctr *weixinUserOrganizationRelationRepo) GetUserOrganizationRelations(ctx context.Context, userId uint64) (*v1.GetUserOrganizationRelationsReply, error) {
	return ctr.data.weixinuc.GetUserOrganizationRelations(ctx, &v1.GetUserOrganizationRelationsRequest{
		UserId: userId,
	})
}
