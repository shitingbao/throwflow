package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/company/v1"
	"task/internal/biz"
)

type companyRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyRepo(data *Data, logger log.Logger) biz.CompanyRepo {
	return &companyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *companyRepo) SyncUpdateStatus(ctx context.Context) (*v1.SyncUpdateStatusCompanysReply, error) {
	company, err := cr.data.companyuc.SyncUpdateStatusCompanys(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return company, err
}
