package data

import (
	v1 "company/api/service/material/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type materialMaterialRepo struct {
	data *Data
	log  *log.Helper
}

func NewMaterialMaterialRepo(data *Data, logger log.Logger) biz.MaterialMaterialRepo {
	return &materialMaterialRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (mmr *materialMaterialRepo) GetIsTop(ctx context.Context, productId uint64) (*v1.GetIsTopMaterialsReply, error) {
	list, err := mmr.data.materialuc.GetIsTopMaterials(ctx, &v1.GetIsTopMaterialsRequest{
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (mmr *materialMaterialRepo) List(ctx context.Context, pageNum, pageSize, productId uint64) (*v1.ListMaterialsReply, error) {
	list, err := mmr.data.materialuc.ListMaterials(ctx, &v1.ListMaterialsRequest{
		PageNum:   pageNum,
		PageSize:  pageSize,
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (mmr *materialMaterialRepo) ListAwemesByProductId(ctx context.Context, productId uint64) (*v1.ListAwemesByProductIdReply, error) {
	list, err := mmr.data.materialuc.ListAwemesByProductId(ctx, &v1.ListAwemesByProductIdRequest{
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
