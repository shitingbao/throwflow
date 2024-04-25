package data

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type weixinUserOpenDouyinRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserOpenDouyinRepo(data *Data, logger log.Logger) biz.WeixinUserOpenDouyinRepo {
	return &weixinUserOpenDouyinRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// UserOpenDouyin
func (wur *weixinUserOpenDouyinRepo) ListByClientKeyAndOpenIds(ctx context.Context, pageNum, pageSize uint64, clientKeyAndOpenIds, keyword string) (*v1.ListByClientKeyAndOpenIdsReply, error) {
	return wur.data.weixinuc.ListByClientKeyAndOpenIds(ctx, &v1.ListByClientKeyAndOpenIdsRequest{
		PageNum:             pageNum,
		PageSize:            pageSize,
		ClientKeyAndOpenIds: clientKeyAndOpenIds,
		Keyword:             keyword,
	})
}

func (ctr *weixinUserOpenDouyinRepo) ListOpenDouyinUsers(ctx context.Context, userId, pageNum, pageSize uint64, keyword string) (*v1.ListOpenDouyinUsersReply, error) {
	return ctr.data.weixinuc.ListOpenDouyinUsers(ctx, &v1.ListOpenDouyinUsersRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		UserId:   userId,
		Keyword:  keyword,
	})
}
