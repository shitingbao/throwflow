package data

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type industryRepo struct {
	data *Data
	log  *log.Helper
}

func NewIndustryRepo(data *Data, logger log.Logger) biz.IndustryRepo {
	return &industryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ir *industryRepo) List(ctx context.Context) (*v1.ListIndustriesReply, error) {
	list, err := ir.data.companyuc.ListIndustries(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}
