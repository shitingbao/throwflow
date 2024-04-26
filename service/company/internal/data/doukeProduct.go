package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type doukeProductRepo struct {
	data *Data
	log  *log.Helper
}

func NewDoukeProductRepo(data *Data, logger log.Logger) biz.DoukeProductRepo {
	return &doukeProductRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (dpsr *doukeProductRepo) List(ctx context.Context, pageNum, pageSize, cosRatioMin, industryId, categoryId, subCategoryId uint64) (*v1.ListDoukeProductsReply, error) {
	product, err := dpsr.data.douyinuc.ListDoukeProducts(ctx, &v1.ListDoukeProductsRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		CosRatioMin:   cosRatioMin,
		IndustryId:    industryId,
		CategoryId:    categoryId,
		SubCategoryId: subCategoryId,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (dpsr *doukeProductRepo) ListByProductId(ctx context.Context, productIds string) (*v1.ListDoukeProductByProductIdsReply, error) {
	product, err := dpsr.data.douyinuc.ListDoukeProductByProductIds(ctx, &v1.ListDoukeProductByProductIdsRequest{
		ProductIds: productIds,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (dpsr *doukeProductRepo) ListCategory(ctx context.Context, parentId uint64) (*v1.ListCategoryDoukeProductsReply, error) {
	product, err := dpsr.data.douyinuc.ListCategoryDoukeProducts(ctx, &v1.ListCategoryDoukeProductsRequest{
		ParentId: parentId,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (dpsr *doukeProductRepo) Save(ctx context.Context, productUrl, externalInfo string) (*v1.CreateShareDoukeProductsReply, error) {
	productShare, err := dpsr.data.douyinuc.CreateShareDoukeProducts(ctx, &v1.CreateShareDoukeProductsRequest{
		ProductUrl:   productUrl,
		ExternalInfo: externalInfo,
	})

	if err != nil {
		return nil, err
	}

	return productShare, err
}

func (dpsr *doukeProductRepo) Parse(ctx context.Context, content string) (*v1.ParseDoukeProductsReply, error) {
	product, err := dpsr.data.douyinuc.ParseDoukeProducts(ctx, &v1.ParseDoukeProductsRequest{
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}
