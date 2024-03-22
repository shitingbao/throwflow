package data

import (
	"context"
	v1 "task/api/service/company/v1"
	"task/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type companyTaskRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyTaskRepo(data *Data, logger log.Logger) biz.CompanyTaskRepo {
	return &companyTaskRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *companyTaskRepo) Sync(ctx context.Context) (*v1.SyncUpdateCompanyTaskDetailReply, error) {
	companyTask, err := cr.data.companyuc.SyncUpdateCompanyTaskDetail(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return companyTask, err
}
