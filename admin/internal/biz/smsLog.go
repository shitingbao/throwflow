package biz

import (
	v1 "admin/api/service/common/v1"
	"admin/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type SmsLogRepo interface {
	List(context.Context, uint64) (*v1.ListSmsReply, error)
}

type SmsLogUsecase struct {
	repo SmsLogRepo
	conf *conf.Data
	log  *log.Helper
}

func NewSmsLogUsecase(repo SmsLogRepo, conf *conf.Data, logger log.Logger) *SmsLogUsecase {
	return &SmsLogUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (sluc *SmsLogUsecase) ListSmsLogs(ctx context.Context, pageNum uint64) (*v1.ListSmsReply, error) {
	list, err := sluc.repo.List(ctx, pageNum)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}
