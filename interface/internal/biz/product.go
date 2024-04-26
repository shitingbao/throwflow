package biz

import (
	"context"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type ProductRepo interface {
	List(context.Context, uint64, uint64, uint64, uint64, uint64, string, string) (*v1.ListCompanyProductsReply, error)
	ListExternal(context.Context, uint64, uint64, uint64, uint64, uint64, uint32, string) (*v1.ListExternalCompanyProductsReply, error)
	ListCategory(context.Context) (*v1.ListCompanyProductCategorysReply, error)
	ListCompanyTaskProducts(context.Context, uint64, uint64, string) (*v1.ListCompanyTaskProductsReply, error)
	Statistics(context.Context, uint64, uint64, uint64, string, string) (*v1.StatisticsCompanyProductsReply, error)
	GetUploadId(context.Context, string) (*v1.GetUploadIdCompanyProductsReply, error)
	GetExternal(context.Context, uint64) (*v1.GetExternalCompanyProductsReply, error)
	GetExternalProductShare(context.Context, uint64, uint64) (*v1.GetExternalProductShareCompanyProductsReply, error)
	Create(context.Context, string) (*v1.CreateCompanyProductsReply, error)
	UpdateCommission(context.Context, uint64, string) (*v1.UpdateCommissionCompanyProductsReply, error)
	UpdateStatus(context.Context, uint64, uint32) (*v1.UpdateStatusCompanyProductsReply, error)
	UpdateIsTop(context.Context, uint64, uint32) (*v1.UpdateIsTopCompanyProductsReply, error)
	UpdateMaterial(context.Context, uint64, string) (*v1.UpdateMaterialCompanyProductsReply, error)
	UpdateInvestmentRatio(context.Context, uint64, float64) (*v1.UpdateInvestmentRatioCompanyProductsReply, error)
	Parse(context.Context, string) (*v1.ParseCompanyProductsReply, error)
	Verification(context.Context, uint64) (*v1.VerificationCompanyProductsReply, error)
	UploadPart(context.Context, uint64, uint64, uint64, string, string) (*v1.UploadPartCompanyProductsReply, error)
	CompleteUpload(context.Context, string) (*v1.CompleteUploadCompanyProductsReply, error)
	AbortUpload(context.Context, string) (*v1.AbortUploadCompanyProductsReply, error)
	Delete(context.Context, uint64) (*v1.DeleteCompanyProductsReply, error)
}

type ProductUsecase struct {
	repo  ProductRepo
	urepo UserRepo
	conf  *conf.Data
	log   *log.Helper
}

func NewProductUsecase(repo ProductRepo, urepo UserRepo, conf *conf.Data, logger log.Logger) *ProductUsecase {
	return &ProductUsecase{repo: repo, urepo: urepo, conf: conf, log: log.NewHelper(logger)}
}

func (puc *ProductUsecase) ListProducts(ctx context.Context, pageNum, pageSize, industryId, categoryId, subCategoryId uint64, productStatus, keyword string) (*v1.ListCompanyProductsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(puc.conf.Database.PageSize)
	}

	list, err := puc.repo.List(ctx, industryId, categoryId, subCategoryId, pageNum, pageSize, productStatus, keyword)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (puc *ProductUsecase) ListMiniProducts(ctx context.Context, pageNum, pageSize, industryId, categoryId, subCategoryId uint64, isInvestment uint32, keyword string) (*v1.ListExternalCompanyProductsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(puc.conf.Database.PageSize)
	}

	list, err := puc.repo.ListExternal(ctx, industryId, categoryId, subCategoryId, pageNum, pageSize, isInvestment, keyword)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (puc *ProductUsecase) ListCategorys(ctx context.Context) (*v1.ListCompanyProductCategorysReply, error) {
	list, err := puc.repo.ListCategory(ctx)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PRODUCT_CATEGORY_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (puc *ProductUsecase) ListMiniCategorys(ctx context.Context) (*v1.ListCompanyProductCategorysReply, error) {
	list, err := puc.repo.ListCategory(ctx)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PRODUCT_CATEGORY_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (puc *ProductUsecase) ListCompanyTaskProducts(ctx context.Context, pageNum, pageSize uint64, keyword string) (*v1.ListCompanyTaskProductsReply, error) {
	list, err := puc.repo.ListCompanyTaskProducts(ctx, pageNum, pageSize, keyword)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (puc *ProductUsecase) StatisticsMiniProducts(ctx context.Context, industryId, categoryId, subCategoryId uint64, keyword string) (*v1.StatisticsCompanyProductsReply, error) {
	statistics, err := puc.repo.Statistics(ctx, industryId, categoryId, subCategoryId, "2", keyword)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_STATISTICS_PRODUCT_CATEGORY_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return statistics, nil
}

func (puc *ProductUsecase) GetUploadIdProducts(ctx context.Context, suffix string) (*v1.GetUploadIdCompanyProductsReply, error) {
	uploadId, err := puc.repo.GetUploadId(ctx, suffix)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_UPLOAD_ID_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return uploadId, nil
}

func (puc *ProductUsecase) GetMiniProducts(ctx context.Context, productId uint64) (*v1.GetExternalCompanyProductsReply, error) {
	product, err := puc.repo.GetExternal(ctx, productId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return product, nil
}

func (puc *ProductUsecase) GetMiniProductShareProducts(ctx context.Context, userId, productId uint64) (*v1.GetExternalProductShareCompanyProductsReply, error) {
	productShare, err := puc.repo.GetExternalProductShare(ctx, userId, productId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_PRODUCT_SHARE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return productShare, nil
}

func (puc *ProductUsecase) CreateProducts(ctx context.Context, productCommission string) (*v1.CreateCompanyProductsReply, error) {
	product, err := puc.repo.Create(ctx, productCommission)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return product, nil
}

func (puc *ProductUsecase) UpdateCommissionProducts(ctx context.Context, productId uint64, productCommission string) (*v1.UpdateCommissionCompanyProductsReply, error) {
	product, err := puc.repo.UpdateCommission(ctx, productId, productCommission)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_COMMISSION_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return product, nil
}

func (puc *ProductUsecase) UpdateStatusProducts(ctx context.Context, productId uint64, status uint32) (*v1.UpdateStatusCompanyProductsReply, error) {
	product, err := puc.repo.UpdateStatus(ctx, productId, status)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_STATUS_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return product, nil
}

func (puc *ProductUsecase) UpdateIsTopProducts(ctx context.Context, productId uint64, isTop uint32) (*v1.UpdateIsTopCompanyProductsReply, error) {
	product, err := puc.repo.UpdateIsTop(ctx, productId, isTop)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_IS_TOP_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return product, nil
}

func (puc *ProductUsecase) UpdateMaterialProducts(ctx context.Context, productId uint64, productMaterial string) (*v1.UpdateMaterialCompanyProductsReply, error) {
	list, err := puc.repo.UpdateMaterial(ctx, productId, productMaterial)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_MATERIAL_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (puc *ProductUsecase) UpdateInvestmentRatioProducts(ctx context.Context, productId uint64, investmentRatio float64) (*v1.UpdateInvestmentRatioCompanyProductsReply, error) {
	product, err := puc.repo.UpdateInvestmentRatio(ctx, productId, investmentRatio)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_INVESTMENTRATIO_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return product, nil
}

func (puc *ProductUsecase) ParseMiniProductProducts(ctx context.Context, content string) (*v1.ParseCompanyProductsReply, error) {
	product, err := puc.repo.Parse(ctx, content)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_PARSE_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return product, nil
}

func (puc *ProductUsecase) VerificationMiniProductProducts(ctx context.Context, productId uint64) (*v1.VerificationCompanyProductsReply, error) {
	product, err := puc.repo.Verification(ctx, productId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_VERIFICATION_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return product, nil
}

func (puc *ProductUsecase) UploadPartProducts(ctx context.Context, partNumber, totalPart, contentLength uint64, uploadId, content string) error {
	_, err := puc.repo.UploadPart(ctx, partNumber, totalPart, contentLength, uploadId, content)

	if err != nil {
		return errors.InternalServer("INTERFACE_UPDATE_PART_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return nil
}

func (puc *ProductUsecase) CompleteUploadProducts(ctx context.Context, uploadId string) (*v1.CompleteUploadCompanyProductsReply, error) {
	staticUrl, err := puc.repo.CompleteUpload(ctx, uploadId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPLETE_UPLOAD_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return staticUrl, nil
}

func (puc *ProductUsecase) AbortUploadProducts(ctx context.Context, uploadId string) error {
	_, err := puc.repo.AbortUpload(ctx, uploadId)

	if err != nil {
		return errors.InternalServer("INTERFACE_ABORT_UPLOAD_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return nil
}

func (puc *ProductUsecase) DeleteProducts(ctx context.Context, productId uint64) error {
	_, err := puc.repo.Delete(ctx, productId)

	if err != nil {
		return errors.InternalServer("INTERFACE_DELETE_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return nil
}
