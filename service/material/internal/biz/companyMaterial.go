package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "material/api/service/company/v1"
)

var (
	MaterialCompanyMaterialNotFound    = errors.NotFound("MATERIAL_COMPANY_MATERIAL_NOT_FOUND", "企业素材云库不存在")
	MaterialCompanyMaterialListError   = errors.InternalServer("MATERIAL_COMPANY_MATERIAL_LIST_ERROR", "企业素材云库列表获取失败")
	MaterialCompanyMaterialCreateError = errors.InternalServer("MATERIAL_COMPANY_MATERIAL_CREATE_ERROR", "企业素材云库创建失败")
	MaterialCompanyMaterialUpdateError = errors.InternalServer("MATERIAL_COMPANY_MATERIAL_UPDATE_ERROR", "企业素材云库更新失败")
)

type ListCompanyMaterialsReplyCompanyMaterial struct {
	Id          uint64                                      `json:"id"`
	LibraryName string                                      `json:"libraryName"`
	ChildList   []*ListCompanyMaterialsReplyCompanyMaterial `json:"childList"`
}

type CompanyMaterialRepo interface {
	Get(context.Context, uint64, uint64) (*v1.GetCompanyMaterialsReply, error)
	GetByLibraryName(context.Context, uint64, uint64, string) (*v1.GetCompanyMaterialByLibraryNamesReply, error)
	GetByProductId(context.Context, uint64, uint64, uint64) (*v1.GetCompanyMaterialByProductIdsReply, error)
	List(context.Context, uint64) (*v1.ListCompanyMaterialsReply, error)
	Save(context.Context, uint64, uint64, uint64, uint64, uint64, uint32, string, string, string) (*v1.CreateCompanyMaterialsReply, error)
	UpdateFileSize(context.Context, uint64, uint64, uint64) (*v1.UpdateFileSizeCompanyMaterialsReply, error)
}
