package biz

import (
	"context"
	v1 "interface/api/service/douyin/v1"
)

type QianchuanAdvertiserRepo interface {
	List(context.Context, uint64, uint64, uint64, string) (*v1.ListQianchuanAdvertisersReply, error)
	Statistics(context.Context, uint64) (*v1.StatisticsQianchuanAdvertisersReply, error)
	Update(context.Context, uint64, uint64, uint32) (*v1.UpdateStatusQianchuanAdvertisersReply, error)
}
