package data

import (
	v1 "admin/api/service/common/v1"
	"admin/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
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
