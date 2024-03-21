package service

import (
	v1 "company/api/company/v1"
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"encoding/json"
	"time"
)

type ListCompanyMaterialsReplyCompanyMaterial struct {
	Id          uint64                                      `json:"id"`
	LibraryName string                                      `json:"libraryName"`
	ChildList   []*ListCompanyMaterialsReplyCompanyMaterial `json:"childList"`
}

func (cs *CompanyService) ListCompanyMaterials(ctx context.Context, in *v1.ListCompanyMaterialsRequest) (*v1.ListCompanyMaterialsReply, error) {
	companyMaterials, err := cs.cmuc.ListCompanyMaterials(ctx, in.CompanyId)

	if err != nil {
		return nil, err
	}

	listCompanyMaterials := cs.listCompanyMaterials(ctx, companyMaterials)

	list, _ := json.Marshal(listCompanyMaterials)

	return &v1.ListCompanyMaterialsReply{
		Code: 200,
		Data: &v1.ListCompanyMaterialsReply_Data{
			List: string(list),
		},
	}, nil
}

func (cs *CompanyService) GetCompanyMaterials(ctx context.Context, in *v1.GetCompanyMaterialsRequest) (*v1.GetCompanyMaterialsReply, error) {
	companyMaterial, err := cs.cmuc.GetCompanyMaterials(ctx, in.CompanyId, in.CompanyMaterialId)

	if err != nil {
		return nil, err
	}

	return &v1.GetCompanyMaterialsReply{
		Code: 200,
		Data: &v1.GetCompanyMaterialsReply_Data{
			CompanyMaterialId:          companyMaterial.Id,
			CompanyMaterialLibraryName: companyMaterial.LibraryName,
			CompanyMaterialLibraryType: uint32(companyMaterial.LibraryType),
		},
	}, nil
}

func (cs *CompanyService) GetCompanyMaterialByLibraryNames(ctx context.Context, in *v1.GetCompanyMaterialByLibraryNamesRequest) (*v1.GetCompanyMaterialByLibraryNamesReply, error) {
	companyMaterial, err := cs.cmuc.GetCompanyMaterialByLibraryNames(ctx, in.CompanyId, in.ParentId, in.MaterialLibraryName)

	if err != nil {
		return nil, err
	}

	return &v1.GetCompanyMaterialByLibraryNamesReply{
		Code: 200,
		Data: &v1.GetCompanyMaterialByLibraryNamesReply_Data{
			CompanyMaterialId:          companyMaterial.Id,
			CompanyMaterialLibraryName: companyMaterial.LibraryName,
		},
	}, nil
}

func (cs *CompanyService) GetCompanyMaterialByProductIds(ctx context.Context, in *v1.GetCompanyMaterialByProductIdsRequest) (*v1.GetCompanyMaterialByProductIdsReply, error) {
	companyMaterial, err := cs.cmuc.GetCompanyMaterialByProductIds(ctx, in.CompanyId, in.ParentId, in.ProductId)

	if err != nil {
		return nil, err
	}

	return &v1.GetCompanyMaterialByProductIdsReply{
		Code: 200,
		Data: &v1.GetCompanyMaterialByProductIdsReply_Data{
			CompanyMaterialId:          companyMaterial.Id,
			CompanyMaterialLibraryName: companyMaterial.LibraryName,
		},
	}, nil
}

func (cs *CompanyService) GetUploadIdCompanyMaterials(ctx context.Context, in *v1.GetUploadIdCompanyMaterialsRequest) (*v1.GetUploadIdCompanyMaterialsReply, error) {
	uploadId, err := cs.cmuc.GetUploadIdCompanyMaterials(ctx, in.CompanyId, in.Suffix)

	if err != nil {
		return nil, err
	}

	return &v1.GetUploadIdCompanyMaterialsReply{
		Code: 200,
		Data: &v1.GetUploadIdCompanyMaterialsReply_Data{
			UploadId: uploadId,
		},
	}, nil
}

func (cs *CompanyService) CreateCompanyMaterials(ctx context.Context, in *v1.CreateCompanyMaterialsRequest) (*v1.CreateCompanyMaterialsReply, error) {
	if in.MaterialLibraryType == 2 {
		if len(in.MaterialLibraryUrl) == 0 {
			return nil, biz.CompanyValidatorError
		}

		if len(in.MaterialLiberaryFileType) == 0 {
			return nil, biz.CompanyValidatorError
		}

		if in.MaterialLiberaryFileSize == 0 {
			return nil, biz.CompanyValidatorError
		}
	}

	companyMaterial, err := cs.cmuc.CreateCompanyMaterials(ctx, in.CompanyId, in.ParentId, in.ProductId, in.VideoId, in.MaterialLiberaryFileSize, uint8(in.MaterialLibraryType), in.MaterialLibraryName, in.MaterialLibraryUrl, in.MaterialLiberaryFileType)

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyMaterialsReply{
		Code: 200,
		Data: &v1.CreateCompanyMaterialsReply_Data{
			CompanyMaterialId:          companyMaterial.Id,
			CompanyMaterialLibraryName: companyMaterial.LibraryName,
		},
	}, nil
}

func (cs *CompanyService) UpdateFileSizeCompanyMaterials(ctx context.Context, in *v1.UpdateFileSizeCompanyMaterialsRequest) (*v1.UpdateFileSizeCompanyMaterialsReply, error) {
	err := cs.cmuc.UpdateFileSizeCompanyMaterials(ctx, in.CompanyId, in.CompanyMaterialId, in.MaterialLiberaryFileSize)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateFileSizeCompanyMaterialsReply{
		Code: 200,
		Data: &v1.UpdateFileSizeCompanyMaterialsReply_Data{},
	}, nil
}

func (cs *CompanyService) UploadPartCompanyMaterials(ctx context.Context, in *v1.UploadPartCompanyMaterialsRequest) (*v1.UploadPartCompanyMaterialsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := cs.cmuc.UploadPartCompanyMaterials(ctx, in.CompanyId, in.PartNumber, in.TotalPart, in.ContentLength, in.UploadId, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.UploadPartCompanyMaterialsReply{
		Code: 200,
		Data: &v1.UploadPartCompanyMaterialsReply_Data{},
	}, nil
}

func (cs *CompanyService) CompleteUploadCompanyMaterials(ctx context.Context, in *v1.CompleteUploadCompanyMaterialsRequest) (*v1.CompleteUploadCompanyMaterialsReply, error) {
	staticUrl, err := cs.cmuc.CompleteUploadCompanyMaterials(ctx, in.CompanyId, in.UploadId)

	if err != nil {
		return nil, err
	}

	return &v1.CompleteUploadCompanyMaterialsReply{
		Code: 200,
		Data: &v1.CompleteUploadCompanyMaterialsReply_Data{
			StaticUrl: staticUrl,
		},
	}, nil
}

func (cs *CompanyService) AbortUploadCompanyMaterials(ctx context.Context, in *v1.AbortUploadCompanyMaterialsRequest) (*v1.AbortUploadCompanyMaterialsReply, error) {
	err := cs.cmuc.AbortUploadCompanyMaterials(ctx, in.CompanyId, in.UploadId)

	if err != nil {
		return nil, err
	}

	return &v1.AbortUploadCompanyMaterialsReply{
		Code: 200,
		Data: &v1.AbortUploadCompanyMaterialsReply_Data{},
	}, nil
}

func (cs *CompanyService) DeleteCompanyMaterials(ctx context.Context, in *v1.DeleteCompanyMaterialsRequest) (*v1.DeleteCompanyMaterialsReply, error) {
	err := cs.cmuc.DeleteCompanyMaterials(ctx, in.StaticUrl)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteCompanyMaterialsReply{
		Code: 200,
		Data: &v1.DeleteCompanyMaterialsReply_Data{},
	}, nil
}

func (cs *CompanyService) listCompanyMaterials(ctx context.Context, companyMaterials []*domain.CompanyMaterialLibrary) []*ListCompanyMaterialsReplyCompanyMaterial {
	if companyMaterials == nil {
		return nil
	}

	list := make([]*ListCompanyMaterialsReplyCompanyMaterial, 0)

	for _, lcompanyMaterial := range companyMaterials {
		companyMaterial := &ListCompanyMaterialsReplyCompanyMaterial{
			Id:          lcompanyMaterial.Id,
			LibraryName: lcompanyMaterial.LibraryName,
		}

		companyMaterial.ChildList = cs.listCompanyMaterials(ctx, lcompanyMaterial.ChildList)

		list = append(list, companyMaterial)
	}

	return list
}
