package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/pkg/douke/kolProduct"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	DouyinDoukeProductShareCreateError = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_SHARE_CREATE_ERROR", "抖客商品分销转链创建失败")
)

type DoukeProductShareUsecase struct {
	dtrepo DoukeTokenRepo
	tm     Transaction
	conf   *conf.Data
	dconf  *conf.Douke
	log    *log.Helper
}

func NewDoukeProductShareUsecase(dtrepo DoukeTokenRepo, tm Transaction, conf *conf.Data, dconf *conf.Douke, logger log.Logger) *DoukeProductShareUsecase {
	return &DoukeProductShareUsecase{dtrepo: dtrepo, tm: tm, conf: conf, dconf: dconf, log: log.NewHelper(logger)}
}

func (dpsuc *DoukeProductShareUsecase) CreateDoukeProductShares(ctx context.Context, productUrl, externalInfo string) (*kolProduct.ShareKolProductResponse, error) {
	doukeToken, err := dpsuc.dtrepo.Get(ctx, dpsuc.dconf.AuthorityId, dpsuc.dconf.AuthSubjectType)

	if err != nil {
		return nil, DouyinDoukeTokenNotFound
	}

	productShare, err := kolProduct.ShareKolProduct(dpsuc.dconf.AppKey, dpsuc.dconf.AppSecret, doukeToken.AccessToken, productUrl, dpsuc.dconf.Pid, externalInfo)

	if err != nil {
		if len(err.Error()) > 0 {
			return nil, errors.InternalServer("DOUYIN_DOUKE_PRODUCT_SHARE_CREATE_ERROR", err.Error())
		} else {
			return nil, DouyinDoukeProductShareCreateError
		}
	}

	return productShare, nil
}
