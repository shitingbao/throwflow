package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/common/v1"
	"weixin/internal/biz"
)

type areaRepo struct {
	data *Data
	log  *log.Helper
}

func NewAreaRepo(data *Data, logger log.Logger) biz.AreaRepo {
	return &areaRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ar *areaRepo) GetByAreaCode(ctx context.Context, areaCode uint64) (*v1.GetAreasReply, error) {
	area, err := ar.data.commonuc.GetAreas(ctx, &v1.GetAreasRequest{
		AreaCode: areaCode,
	})

	if err != nil {
		return nil, err
	}

	return area, err
}
