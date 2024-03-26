package service

import (
	v1 "common/api/common/v1"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"time"
)

func (cs *CommonService) CreateShortCode(ctx context.Context, in *empty.Empty) (*v1.CreateShortCodeReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shortCode, err := cs.scuc.CreateShortCode(ctx)

	if err != nil {
		return nil, err
	}

	return &v1.CreateShortCodeReply{
		Code: 200,
		Data: &v1.CreateShortCodeReply_Data{
			ShortCode: shortCode,
		},
	}, nil
}
