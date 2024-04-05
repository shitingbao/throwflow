package service

import (
	"github.com/google/wire"
	v1 "material/api/material/v1"
	"material/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewMaterialService)

type MaterialService struct {
	v1.UnimplementedMaterialServer

	muc *biz.MaterialUsecase
	cuc *biz.CollectUsecase
}

func NewMaterialService(muc *biz.MaterialUsecase, cuc *biz.CollectUsecase) *MaterialService {
	return &MaterialService{muc: muc, cuc: cuc}
}
