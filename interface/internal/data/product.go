package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (pr *productRepo) List(ctx context.Context, industryId, categoryId, subCategoryId, pageNum, pageSize uint64, productStatus, keyword string) (*v1.ListCompanyProductsReply, error) {
	list, err := pr.data.companyuc.ListCompanyProducts(ctx, &v1.ListCompanyProductsRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		ProductStatus: productStatus,
		IndustryId:    industryId,
		CategoryId:    categoryId,
		SubCategoryId: subCategoryId,
		Keyword:       keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (pr *productRepo) ListExternal(ctx context.Context, industryId, categoryId, subCategoryId, pageNum, pageSize uint64, isInvestment uint32, keyword string) (*v1.ListExternalCompanyProductsReply, error) {
	list, err := pr.data.companyuc.ListExternalCompanyProducts(ctx, &v1.ListExternalCompanyProductsRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		IndustryId:    industryId,
		CategoryId:    categoryId,
		SubCategoryId: subCategoryId,
		IsInvestment:  isInvestment,
		Keyword:       keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (pr *productRepo) ListCategory(ctx context.Context) (*v1.ListCompanyProductCategorysReply, error) {
	list, err := pr.data.companyuc.ListCompanyProductCategorys(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (pr *productRepo) Statistics(ctx context.Context, industryId, categoryId, subCategoryId uint64, productStatus, keyword string) (*v1.StatisticsCompanyProductsReply, error) {
	list, err := pr.data.companyuc.StatisticsCompanyProducts(ctx, &v1.StatisticsCompanyProductsRequest{
		ProductStatus: productStatus,
		IndustryId:    industryId,
		CategoryId:    categoryId,
		SubCategoryId: subCategoryId,
		Keyword:       keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (pr *productRepo) GetUploadId(ctx context.Context, suffix string) (*v1.GetUploadIdCompanyProductsReply, error) {
	uploadId, err := pr.data.companyuc.GetUploadIdCompanyProducts(ctx, &v1.GetUploadIdCompanyProductsRequest{
		Suffix: suffix,
	})

	if err != nil {
		return nil, err
	}

	return uploadId, err
}

func (pr *productRepo) GetExternal(ctx context.Context, productId uint64) (*v1.GetExternalCompanyProductsReply, error) {
	product, err := pr.data.companyuc.GetExternalCompanyProducts(ctx, &v1.GetExternalCompanyProductsRequest{
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) GetExternalProductShare(ctx context.Context, userId, productId uint64) (*v1.GetExternalProductShareCompanyProductsReply, error) {
	productShare, err := pr.data.companyuc.GetExternalProductShareCompanyProducts(ctx, &v1.GetExternalProductShareCompanyProductsRequest{
		ProductId: productId,
		UserId:    userId,
	})

	if err != nil {
		return nil, err
	}

	return productShare, err
}

func (pr *productRepo) Create(ctx context.Context, commission string) (*v1.CreateCompanyProductsReply, error) {
	product, err := pr.data.companyuc.CreateCompanyProducts(ctx, &v1.CreateCompanyProductsRequest{
		Commission: commission,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) UpdateCommission(ctx context.Context, productId uint64, commission string) (*v1.UpdateCommissionCompanyProductsReply, error) {
	product, err := pr.data.companyuc.UpdateCommissionCompanyProducts(ctx, &v1.UpdateCommissionCompanyProductsRequest{
		ProductId:  productId,
		Commission: commission,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) UpdateStatus(ctx context.Context, productId uint64, status uint32) (*v1.UpdateStatusCompanyProductsReply, error) {
	product, err := pr.data.companyuc.UpdateStatusCompanyProducts(ctx, &v1.UpdateStatusCompanyProductsRequest{
		ProductId: productId,
		Status:    status,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) UpdateIsTop(ctx context.Context, productId uint64, isTop uint32) (*v1.UpdateIsTopCompanyProductsReply, error) {
	product, err := pr.data.companyuc.UpdateIsTopCompanyProducts(ctx, &v1.UpdateIsTopCompanyProductsRequest{
		ProductId: productId,
		IsTop:     isTop,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) UpdateSampleThreshold(ctx context.Context, productId, sampleThresholdValue uint64, sampleThresholdType uint32) (*v1.UpdateSampleThresholdCompanyProductsReply, error) {
	product, err := pr.data.companyuc.UpdateSampleThresholdCompanyProducts(ctx, &v1.UpdateSampleThresholdCompanyProductsRequest{
		ProductId:            productId,
		SampleThresholdType:  sampleThresholdType,
		SampleThresholdValue: sampleThresholdValue,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) UpdateMaterial(ctx context.Context, productId uint64, productMaterial string) (*v1.UpdateMaterialCompanyProductsReply, error) {
	list, err := pr.data.companyuc.UpdateMaterialCompanyProducts(ctx, &v1.UpdateMaterialCompanyProductsRequest{
		ProductId:       productId,
		ProductMaterial: productMaterial,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (pr *productRepo) UpdateInvestmentRatio(ctx context.Context, productId uint64, investmentRatio float64) (*v1.UpdateInvestmentRatioCompanyProductsReply, error) {
	list, err := pr.data.companyuc.UpdateInvestmentRatioCompanyProducts(ctx, &v1.UpdateInvestmentRatioCompanyProductsRequest{
		ProductId:       productId,
		InvestmentRatio: investmentRatio,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (pr *productRepo) UploadPart(ctx context.Context, partNumber, totalPart, contentLength uint64, uploadId, content string) (*v1.UploadPartCompanyProductsReply, error) {
	product, err := pr.data.companyuc.UploadPartCompanyProducts(ctx, &v1.UploadPartCompanyProductsRequest{
		UploadId:      uploadId,
		PartNumber:    partNumber,
		TotalPart:     totalPart,
		ContentLength: contentLength,
		Content:       content,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) CompleteUpload(ctx context.Context, uploadId string) (*v1.CompleteUploadCompanyProductsReply, error) {
	product, err := pr.data.companyuc.CompleteUploadCompanyProducts(ctx, &v1.CompleteUploadCompanyProductsRequest{
		UploadId: uploadId,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) AbortUpload(ctx context.Context, uploadId string) (*v1.AbortUploadCompanyProductsReply, error) {
	product, err := pr.data.companyuc.AbortUploadCompanyProducts(ctx, &v1.AbortUploadCompanyProductsRequest{
		UploadId: uploadId,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}

func (pr *productRepo) Delete(ctx context.Context, productId uint64) (*v1.DeleteCompanyProductsReply, error) {
	product, err := pr.data.companyuc.DeleteCompanyProducts(ctx, &v1.DeleteCompanyProductsRequest{
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return product, err
}
