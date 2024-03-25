package service

import (
	"context"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) ListAreas(ctx context.Context, in *v1.ListAreasRequest) (*v1.ListAreasReply, error) {
	areas, err := is.auc.ListAreas(ctx, in.ParentAreaCode)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListAreasReply_Areas, 0)

	for _, area := range areas.Data.List {
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
