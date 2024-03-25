package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserScanRecordRepo interface {
	Create(context.Context, uint64, uint64, uint64) (*v1.CreateUserScanRecordsReply, error)
}

type UserScanRecordUsecase struct {
	repo UserScanRecordRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserScanRecordUsecase(repo UserScanRecordRepo, conf *conf.Data, logger log.Logger) *UserScanRecordUsecase {
	return &UserScanRecordUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (usruc *UserScanRecordUsecase) CreateUserScanRecords(ctx context.Context, userId, organizationId, parentUserId uint64) (*v1.CreateUserScanRecordsReply, error) {
	scanRecord, err := usruc.repo.Create(ctx, userId, organizationId, parentUserId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_SCAN_RECORD_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return scanRecord, nil
}
