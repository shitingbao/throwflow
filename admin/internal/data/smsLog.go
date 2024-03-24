package data

import (
	v1 "admin/api/service/common/v1"
	"admin/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type smsLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewSmsLogRepo(data *Data, logger log.Logger) biz.SmsLogRepo {
	return &smsLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (slr *smsLogRepo) List(ctx context.Context, pageNum uint64) (*v1.ListSmsReply, error) {
	smsLogList, err := slr.data.commonuc.ListSms(ctx, &v1.ListSmsRequest{
		PageNum:  pageNum,
		PageSize: uint64(slr.data.conf.Database.PageSize),
	})

	if err != nil {
		return nil, err
	}

	return smsLogList, nil
}
