package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
)

func (ds *DouyinService) CreateDoukeProductShares(ctx context.Context, in *v1.CreateDoukeProductSharesRequest) (*v1.CreateDoukeProductSharesReply, error) {
	productShare, err := ds.dpsuc.CreateDoukeProductShares(ctx, in.ProductUrl, in.ExternalInfo)

	if err != nil {
		return nil, err
	}
	
	return &v1.CreateDoukeProductSharesReply{
		Code: 200,
		Data: &v1.CreateDoukeProductSharesReply_Data{
			DyPassword: productShare.Data.DyPassword,
		},
	}, nil
}
