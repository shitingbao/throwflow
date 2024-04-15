package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/pkg/csj/command"
	"douyin/internal/pkg/csj/product"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

var (
	DouyinDoukeProductGetError         = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_GET_ERROR", "抖客商品获取失败")
	DouyinDoukeProductShareCreateError = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_SHARE_CREATE_ERROR", "抖客商品分销转链创建失败")
	DouyinDoukeProductParseError       = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_PARSE_ERROR", "抖客抖口令转解析失败")
)

type DoukeProductUsecase struct {
	tm    Transaction
	conf  *conf.Data
	cconf *conf.Csj
	log   *log.Helper
}

func NewDoukeProductUsecase(tm Transaction, conf *conf.Data, cconf *conf.Csj, logger log.Logger) *DoukeProductUsecase {
	return &DoukeProductUsecase{tm: tm, conf: conf, cconf: cconf, log: log.NewHelper(logger)}
}

func (dpuc *DoukeProductUsecase) GetDoukeProducts(ctx context.Context, productId uint64) (*product.DetailResponse, error) {
	product, err := product.Detail(dpuc.cconf.AppId, dpuc.cconf.AppSecret, uuid.NewString(), []uint64{productId})

	if err != nil {
		if len(err.Error()) > 0 {
			return nil, errors.InternalServer("DOUYIN_DOUKE_PRODUCT_GET_ERROR", err.Error())
		} else {
			return nil, DouyinDoukeProductGetError
		}
	}

	if len(product.Data.Products) != 1 {
		return nil, DouyinDoukeProductGetError
	}

	return product, nil
}

func (dpuc *DoukeProductUsecase) CreateShareDoukeProducts(ctx context.Context, productUrl, externalInfo string) (*product.LinkResponse, error) {
	productShare, err := product.Link(dpuc.cconf.AppId, dpuc.cconf.AppSecret, productUrl, externalInfo, uuid.NewString())

	if err != nil {
		if len(err.Error()) > 0 {
			return nil, errors.InternalServer("DOUYIN_DOUKE_PRODUCT_SHARE_CREATE_ERROR", err.Error())
		} else {
			return nil, DouyinDoukeProductShareCreateError
		}
	}

	return productShare, nil
}

func (dpuc *DoukeProductUsecase) ParseDoukeProducts(ctx context.Context, content string) (*command.ParseResponse, error) {
	productParse, err := command.Parse(dpuc.cconf.AppId, dpuc.cconf.AppSecret, content, uuid.NewString())

	if err != nil {
		if len(err.Error()) > 0 {
			return nil, errors.InternalServer("DOUYIN_DOUKE_PRODUCT_PARSE_ERROR", err.Error())
		} else {
			return nil, DouyinDoukeProductParseError
		}
	}

	return productParse, nil
}
