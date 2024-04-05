package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "material/api/service/company/v1"
	"material/internal/biz"
)

type companyMaterialRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyMaterialRepo(data *Data, logger log.Logger) biz.CompanyMaterialRepo {
	return &companyMaterialRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cmr *companyMaterialRepo) Get(ctx context.Context, companyId, companyMaterialId uint64) (*v1.GetCompanyMaterialsReply, error) {
	companyMaterial, err := cmr.data.companyuc.GetCompanyMaterials(ctx, &v1.GetCompanyMaterialsRequest{
		CompanyId:         companyId,
		CompanyMaterialId: companyMaterialId,
	})

	if err != nil {
		return nil, err
	}

	return companyMaterial, err
}

func (cmr *companyMaterialRepo) GetByLibraryName(ctx context.Context, companyId, parentId uint64, materialLibraryName string) (*v1.GetCompanyMaterialByLibraryNamesReply, error) {
	list, err := cmr.data.companyuc.GetCompanyMaterialByLibraryNames(ctx, &v1.GetCompanyMaterialByLibraryNamesRequest{
		CompanyId:           companyId,
		ParentId:            parentId,
		MaterialLibraryName: materialLibraryName,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cmr *companyMaterialRepo) GetByProductId(ctx context.Context, companyId, parentId, productId uint64) (*v1.GetCompanyMaterialByProductIdsReply, error) {
	list, err := cmr.data.companyuc.GetCompanyMaterialByProductIds(ctx, &v1.GetCompanyMaterialByProductIdsRequest{
		CompanyId: companyId,
		ParentId:  parentId,
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cmr *companyMaterialRepo) List(ctx context.Context, companyId uint64) (*v1.ListCompanyMaterialsReply, error) {
	list, err := cmr.data.companyuc.ListCompanyMaterials(ctx, &v1.ListCompanyMaterialsRequest{
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cmr *companyMaterialRepo) Save(ctx context.Context, companyId, parentId, productId, videoId, materialLiberaryFileSize uint64, materialLibraryType uint32, materialLibraryName, materialLibraryUrl, materialLiberaryFileType string) (*v1.CreateCompanyMaterialsReply, error) {
	companyMaterial, err := cmr.data.companyuc.CreateCompanyMaterials(ctx, &v1.CreateCompanyMaterialsRequest{
		CompanyId:                companyId,
		ParentId:                 parentId,
		ProductId:                productId,
		VideoId:                  videoId,
		MaterialLibraryName:      materialLibraryName,
		MaterialLibraryType:      materialLibraryType,
		MaterialLibraryUrl:       materialLibraryUrl,
		MaterialLiberaryFileType: materialLiberaryFileType,
		MaterialLiberaryFileSize: materialLiberaryFileSize,
	})

	if err != nil {
		return nil, err
	}

	return companyMaterial, err
}

func (cmr *companyMaterialRepo) UpdateFileSize(ctx context.Context, companyId, companyMaterialId, materialLiberaryFileSize uint64) (*v1.UpdateFileSizeCompanyMaterialsReply, error) {
	companyMaterial, err := cmr.data.companyuc.UpdateFileSizeCompanyMaterials(ctx, &v1.UpdateFileSizeCompanyMaterialsRequest{
		CompanyId:                companyId,
		CompanyMaterialId:        companyMaterialId,
		MaterialLiberaryFileSize: materialLiberaryFileSize,
	})

	if err != nil {
		return nil, err
	}

	return companyMaterial, err
}
