package biz

import (
	"context"
	"douyin/internal/domain"
)

type QianchuanAwemeOrderInfoRepo interface {
	List(context.Context, string, string, string) ([]*domain.AwemeVideoProductQianchuanAwemeOrderInfo, error)
}
