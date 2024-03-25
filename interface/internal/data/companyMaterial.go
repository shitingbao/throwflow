package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type companyMaterialRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyMaterialRepo(data *Data, logger log.Logger) biz.CompanyMaterialRepo {
	return &companyMaterialRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cmr *companyMaterialRepo) List(ctx context.Context, companyId uint64) (*v1.ListCompanyMaterialsReply, error) {
	list, err := cmr.data.companyuc.ListCompanyMaterials(ctx, &v1.ListCompanyMaterialsRequest{
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
