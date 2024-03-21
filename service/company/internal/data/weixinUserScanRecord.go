package data

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type weixinUserScanRecordRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserScanRecordRepo(data *Data, logger log.Logger) biz.WeixinUserScanRecordRepo {
	return &weixinUserScanRecordRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (wusrr *weixinUserScanRecordRepo) Create(ctx context.Context, userId uint64) (*v1.CreateUserScanRecordsReply, error) {
	scanRecord, err := wusrr.data.weixinuc.CreateUserScanRecords(ctx, &v1.CreateUserScanRecordsRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return scanRecord, err
}
