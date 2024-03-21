package domain

import (
	"context"
	"time"
)

type CompanyMaterialLibrary struct {
	Id                        uint64
	CompanyId                 uint64
	UserId                    uint64
	PrincipalUserId           uint64
	ParentId                  uint64
	LibraryName               string
	LibraryType               uint8
	LibraryUrl                string
	LibraryMengmaAudit        uint8
	LibraryMengmaAuditContent string
	LiberaryFileType          string
	LiberaryFileSize          uint64
	ExamineListType           uint8
	ProductId                 uint64
	VideoId                   uint64
	MaterialId                uint64
	IsDel                     uint8
	CreateTime                time.Time
	UpdateTime                time.Time
	ChildList                 []*CompanyMaterialLibrary
}

func NewCompanyMaterialLibrary(ctx context.Context, companyId, parentId, liberaryFileSize uint64, libraryType uint8, libraryName, libraryUrl, liberaryFileType string) *CompanyMaterialLibrary {
	return &CompanyMaterialLibrary{
		CompanyId:        companyId,
		ParentId:         parentId,
		LibraryName:      libraryName,
		LibraryType:      libraryType,
		LibraryUrl:       libraryUrl,
		LiberaryFileType: liberaryFileType,
		LiberaryFileSize: liberaryFileSize,
	}
}

func (cml *CompanyMaterialLibrary) SetCompanyId(ctx context.Context, companyId uint64) {
	cml.CompanyId = companyId
}

func (cml *CompanyMaterialLibrary) SetUserId(ctx context.Context, userId uint64) {
	cml.UserId = userId
}

func (cml *CompanyMaterialLibrary) SetPrincipalUserId(ctx context.Context, principalUserId uint64) {
	cml.PrincipalUserId = principalUserId
}

func (cml *CompanyMaterialLibrary) SetParentId(ctx context.Context, parentId uint64) {
	cml.ParentId = parentId
}

func (cml *CompanyMaterialLibrary) SetLibraryName(ctx context.Context, libraryName string) {
	cml.LibraryName = libraryName
}

func (cml *CompanyMaterialLibrary) SetLibraryType(ctx context.Context, libraryType uint8) {
	cml.LibraryType = libraryType
}

func (cml *CompanyMaterialLibrary) SetLibraryUrl(ctx context.Context, libraryUrl string) {
	cml.LibraryUrl = libraryUrl
}

func (cml *CompanyMaterialLibrary) SetLibraryMengmaAudit(ctx context.Context, libraryMengmaAudit uint8) {
	cml.LibraryMengmaAudit = libraryMengmaAudit
}

func (cml *CompanyMaterialLibrary) SetLibraryMengmaAuditContent(ctx context.Context, libraryMengmaAuditContent string) {
	cml.LibraryMengmaAuditContent = libraryMengmaAuditContent
}

func (cml *CompanyMaterialLibrary) SetLiberaryFileType(ctx context.Context, liberaryFileType string) {
	cml.LiberaryFileType = liberaryFileType
}

func (cml *CompanyMaterialLibrary) SetLiberaryFileSize(ctx context.Context, liberaryFileSize uint64) {
	cml.LiberaryFileSize = liberaryFileSize
}

func (cml *CompanyMaterialLibrary) SetExamineListType(ctx context.Context, examineListType uint8) {
	cml.ExamineListType = examineListType
}

func (cml *CompanyMaterialLibrary) SetProductId(ctx context.Context, productId uint64) {
	cml.ProductId = productId
}

func (cml *CompanyMaterialLibrary) SetVideoId(ctx context.Context, videoId uint64) {
	cml.VideoId = videoId
}

func (cml *CompanyMaterialLibrary) SetMaterialId(ctx context.Context, materialId uint64) {
	cml.MaterialId = materialId
}

func (cml *CompanyMaterialLibrary) SetIsDel(ctx context.Context, isDel uint8) {
	cml.IsDel = isDel
}

func (cml *CompanyMaterialLibrary) SetUpdateTime(ctx context.Context) {
	cml.UpdateTime = time.Now()
}

func (cml *CompanyMaterialLibrary) SetCreateTime(ctx context.Context) {
	cml.CreateTime = time.Now()
}
