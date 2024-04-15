package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/company/v1"
	"weixin/internal/biz"
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

func (ctr *companyTaskRepo) GetCompanyTaskAccountRelation(ctx context.Context, taskAccountRelationId uint64) (*v1.GetCompanyTaskAccountRelationsReply, error) {
	companyTaskAccountRelation, err := ctr.data.companyuc.GetCompanyTaskAccountRelations(ctx, &v1.GetCompanyTaskAccountRelationsRequest{
		TaskAccountRelationId: taskAccountRelationId,
	})

	if err != nil {
		return nil, err
	}

	return companyTaskAccountRelation, err
}
