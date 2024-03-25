package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) ListIndustries(ctx context.Context, in *emptypb.Empty) (*v1.ListIndustriesReply, error) {
	industries, err := is.iuc.ListIndustries(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListIndustriesReply_Industries, 0)

	for _, industry := range industries.Data.List {
		list = append(list, &v1.ListIndustriesReply_Industries{
			Id:           industry.Id,
			IndustryName: industry.IndustryName,
		})
	}

	return &v1.ListIndustriesReply{
		Code: 200,
		Data: &v1.ListIndustriesReply_Data{
			List: list,
		},
	}, nil
}
