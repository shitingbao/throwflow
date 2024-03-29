package biz

import (
	"context"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	DouyinQianchuanWalletListError = errors.InternalServer("DOUYIN_QIANCHUAN_WALLET_LIST_ERROR", "账户钱包获取失败")
)

type QianchuanWalletRepo interface {
	List(context.Context, string) ([]*domain.QianchuanWallet, error)
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.QianchuanWallet) error
}
