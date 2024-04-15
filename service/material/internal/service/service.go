package service

import (
	v1 "material/api/material/v1"
	"material/internal/biz"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewMaterialService)

type MaterialService struct {
	v1.UnimplementedMaterialServer

	muc  *biz.MaterialUsecase
	mcuc *biz.MaterialContentUsecase
	cuc  *biz.CollectUsecase
}

func NewMaterialService(muc *biz.MaterialUsecase, mcuc *biz.MaterialContentUsecase, cuc *biz.CollectUsecase) *MaterialService {
	return &MaterialService{muc: muc, mcuc: mcuc, cuc: cuc}
}
