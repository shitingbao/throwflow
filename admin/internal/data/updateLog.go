package data

import (
	v1 "admin/api/service/common/v1"
	"admin/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type updateLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewupdateLogRepo(data *Data, logger log.Logger) biz.UpdateLogRepo {
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

func (ulr *updateLogRepo) Save(ctx context.Context, name, content string) (*v1.CreateUpdateLogsReply, error) {
	updateLog, err := ulr.data.commonuc.CreateUpdateLogs(ctx, &v1.CreateUpdateLogsRequest{
		Name:    name,
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return updateLog, err
}

func (ulr *updateLogRepo) Update(ctx context.Context, id uint64, name, content string) (*v1.UpdateUpdateLogsReply, error) {
	updateLog, err := ulr.data.commonuc.UpdateUpdateLogs(ctx, &v1.UpdateUpdateLogsRequest{
		Id:      id,
		Name:    name,
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return updateLog, err
}
