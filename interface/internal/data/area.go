package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
	"interface/internal/biz"
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

func (ar *areaRepo) List(ctx context.Context, parentAreaCode uint64) (*v1.ListAreasReply, error) {
	list, err := ar.data.commonuc.ListAreas(ctx, &v1.ListAreasRequest{
		ParentAreaCode: parentAreaCode,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
