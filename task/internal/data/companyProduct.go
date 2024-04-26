package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/company/v1"
	"task/internal/biz"
)

type companyProductRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyProductRepo(data *Data, logger log.Logger) biz.CompanyProductRepo {
	return &companyProductRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cpr *companyProductRepo) Sync(ctx context.Context) (*v1.SyncCompanyProductsReply, error) {
	companyProduct, err := cpr.data.companyuc.SyncCompanyProducts(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return companyProduct, err
}
