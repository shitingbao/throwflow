package service

import (
	v1 "admin/api/admin/v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (as *AdminService) ListIndustries(ctx context.Context, in *emptypb.Empty) (*v1.ListIndustriesReply, error) {
	industries, err := as.iuc.ListIndustries(ctx)

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
