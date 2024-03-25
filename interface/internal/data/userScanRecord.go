package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userScanRecordRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserScanRecordRepo(data *Data, logger log.Logger) biz.UserScanRecordRepo {
	return &userScanRecordRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (usrr *userScanRecordRepo) Create(ctx context.Context, userId, organizationId, parentUserId uint64) (*v1.CreateUserScanRecordsReply, error) {
	list, err := usrr.data.weixinuc.CreateUserScanRecords(ctx, &v1.CreateUserScanRecordsRequest{
		UserId:         userId,
		OrganizationId: organizationId,
		ParentUserId:   parentUserId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
