package biz

import (
	v1 "company/api/service/material/v1"
	"context"
)

type MaterialMaterialRepo interface {
	GetIsTop(context.Context, uint64) (*v1.GetIsTopMaterialsReply, error)
	List(context.Context, uint64, uint64, uint64) (*v1.ListMaterialsReply, error)
	ListAwemesByProductId(context.Context, uint64) (*v1.ListAwemesByProductIdReply, error)
}
