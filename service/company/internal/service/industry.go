package service

import (
	v1 "company/api/company/v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (cs *CompanyService) ListIndustries(ctx context.Context, in *emptypb.Empty) (*v1.ListIndustriesReply, error) {
	industries, err := cs.iuc.ListIndustries(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListIndustriesReply_Industry, 0)

	for _, industry := range industries {
		list = append(list, &v1.ListIndustriesReply_Industry{
			Id:           industry.Id,
			IndustryName: industry.IndustryName,
			Status:       uint32(industry.Status),
		})
	}

	return &v1.ListIndustriesReply{
		Code: 200,
		Data: &v1.ListIndustriesReply_Data{
			List: list,
		},
	}, nil
}
