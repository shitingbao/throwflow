package service

import (
	v1 "common/api/common/v1"
	"context"
)

func (cs *CommonService) ListAreas(ctx context.Context, in *v1.ListAreasRequest) (*v1.ListAreasReply, error) {
	areas, err := cs.auc.ListAreas(ctx, in.ParentAreaCode)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListAreasReply_Areas, 0)

	for _, area := range areas {
		list = append(list, &v1.ListAreasReply_Areas{
			AreaCode:       area.AreaCode,
			ParentAreaCode: area.ParentAreaCode,
			AreaName:       area.AreaName,
		})
	}

	return &v1.ListAreasReply{
		Code: 200,
		Data: &v1.ListAreasReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CommonService) GetAreas(ctx context.Context, in *v1.GetAreasRequest) (*v1.GetAreasReply, error) {
	area, err := cs.auc.GetAreas(ctx, in.AreaCode)

	if err != nil {
		return nil, err
	}

	return &v1.GetAreasReply{
		Code: 200,
		Data: &v1.GetAreasReply_Data{
			AreaCode:       area.AreaCode,
			ParentAreaCode: area.ParentAreaCode,
			AreaName:       area.AreaName,
		},
	}, nil
}
