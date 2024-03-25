package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/service/common/v1"
	"interface/internal/biz"
)

type updateLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewUpdateLogRepo(data *Data, logger log.Logger) biz.UpdateLogRepo {
	return &updateLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ulr *updateLogRepo) List(ctx context.Context) (*v1.ListUpdateLogsReply, error) {
	list, err := ulr.data.commonuc.ListUpdateLogs(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}
