package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/douyin/v1"
	"interface/internal/biz"
)

type oceanengineConfigRepo struct {
	data *Data
	log  *log.Helper
}

func NewOceanengineConfigRepo(data *Data, logger log.Logger) biz.OceanengineConfigRepo {
	return &oceanengineConfigRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ocr *oceanengineConfigRepo) Rand(ctx context.Context, oceanengineType uint32) (*v1.RandOceanengineConfigsReply, error) {
	list, err := ocr.data.douyinuc.RandOceanengineConfigs(ctx, &v1.RandOceanengineConfigsRequest{
		OceanengineType: oceanengineType,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
