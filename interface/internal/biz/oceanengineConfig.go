package biz

import (
	"context"
	v1 "interface/api/service/douyin/v1"
)

type OceanengineConfigRepo interface {
	Rand(context.Context, uint32) (*v1.RandOceanengineConfigsReply, error)
}
