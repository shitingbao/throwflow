package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/pkg/csj/command"
	"douyin/internal/pkg/csj/product"
	"douyin/internal/pkg/tool"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

var (
	DouyinDoukeProductGetError          = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_GET_ERROR", "抖客商品获取失败")
	DouyinDoukeProductListError         = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_LIST_ERROR", "抖客商品列表获取失败")
	DouyinDoukeProductCategoryListError = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_CATEGORY_LIST_ERROR", "抖客商品分类列表获取失败")
	DouyinDoukeProductShareCreateError  = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_SHARE_CREATE_ERROR", "抖客商品分销转链创建失败")
	DouyinDoukeProductParseError        = errors.InternalServer("DOUYIN_DOUKE_PRODUCT_PARSE_ERROR", "抖客抖口令转解析失败")
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

func (dpuc *DoukeProductUsecase) ListDoukeProducts(ctx context.Context, pageNum, pageSize, cosRatioMin, firstCid, secondCid, thirdCid uint64) (*product.ListResponse, error) {
	products, err := product.List(dpuc.cconf.AppId, dpuc.cconf.AppSecret, uuid.NewString(), pageNum, pageSize, cosRatioMin, []uint64{firstCid}, []uint64{secondCid}, []uint64{thirdCid})

	if err != nil {
		if len(err.Error()) > 0 {
			return nil, errors.InternalServer("DOUYIN_DOUKE_PRODUCT_LIST_ERROR", err.Error())
		} else {
			return nil, DouyinDoukeProductListError
		}
	}

	return products, nil
}

func (dpuc *DoukeProductUsecase) ListDoukeProductByProductIds(ctx context.Context, productIds string) (*product.DetailResponse, error) {
	sproductIds := tool.RemoveEmptyString(strings.Split(productIds, ","))

	if len(sproductIds) == 0 {
		return nil, DouyinValidatorError
	}

	iproductIds := make([]uint64, 0)

	for _, sproductId := range sproductIds {
		if iproductId, err := strconv.ParseUint(sproductId, 10, 64); err == nil {
			iproductIds = append(iproductIds, iproductId)
		}
	}

	if len(iproductIds) == 0 {
		return nil, DouyinValidatorError
	}

	products, err := product.Detail(dpuc.cconf.AppId, dpuc.cconf.AppSecret, uuid.NewString(), iproductIds)

	if err != nil {
		if len(err.Error()) > 0 {
			return nil, errors.InternalServer("DOUYIN_DOUKE_PRODUCT_GET_ERROR", err.Error())
		} else {
			return nil, DouyinDoukeProductGetError
		}
	}

	return products, nil
}

func (dpuc *DoukeProductUsecase) ListCategoryDoukeProducts(ctx context.Context, parentId uint64) (*product.CategoryResponse, error) {
	categories, err := product.Category(dpuc.cconf.AppId, dpuc.cconf.AppSecret, uuid.NewString(), parentId)

	if err != nil {
		if len(err.Error()) > 0 {
			return nil, errors.InternalServer("DOUYIN_DOUKE_PRODUCT_CATEGORY_LIST_ERROR", err.Error())
		} else {
			return nil, DouyinDoukeProductCategoryListError
		}
	}

	return categories, nil
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
