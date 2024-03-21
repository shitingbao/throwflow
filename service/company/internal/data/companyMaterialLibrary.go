package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业素材云库表
type CompanyMaterialLibrary struct {
	Id                        uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyId                 uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;index:company_id_parent_id_product_id;not null;comment:企业库ID"`
	UserId                    uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:企业用户ID"`
	PrincipalUserId           uint64    `gorm:"column:principal_user_id;type:bigint(20) UNSIGNED;not null;comment:素材负责人ID"`
	ParentId                  uint64    `gorm:"column:parent_id;type:bigint(20) UNSIGNED;index:company_id_parent_id_product_id;not null;default:0;comment:父级ID"`
	LibraryName               string    `gorm:"column:library_name;type:text;not null;comment:文件夹名/文件名"`
	LibraryType               uint8     `gorm:"column:library_type;type:tinyint(3) UNSIGNED;not null;default:1;comment:类型，1：文件夹，2：文件"`
	LibraryUrl                string    `gorm:"column:library_url;type:varchar(250);not null;default:'';comment:文件链接地址"`
	LibraryMengmaAudit        uint8     `gorm:"column:library_mengma_audit;type:tinyint(3) UNSIGNED;not null;default:0;comment:内审状态：0：默认值，1：已提交，2：审核通过，3：审核不通过"`
	LibraryMengmaAuditContent string    `gorm:"column:library_mengma_audit_content;type:text;comment:内审审核内容"`
	LiberaryFileType          string    `gorm:"column:liberary_file_type;type:varchar(50);not null;default:'';comment:文件类型"`
	LiberaryFileSize          uint64    `gorm:"column:liberary_file_size;type:bigint(20) UNSIGNED;not null;default:0;comment:文件大小，单位:Byte"`
	ExamineListType           uint8     `gorm:"column:examine_list_type;type:tinyint(3) UNSIGNED;not null;default:0;comment:审批列表状态"`
	ProductId                 uint64    `gorm:"column:product_id;type:bigint(20) UNSIGNED;index:company_id_parent_id_product_id;not null;default:0;comment:商品ID"`
	VideoId                   uint64    `gorm:"column:video_id;type:bigint(20) UNSIGNED;not null;default:0;comment:抖音账号下的视频ID"`
	MaterialId                uint64    `gorm:"column:material_id;type:bigint(20) UNSIGNED;not null;default:0;comment:素材ID"`
	IsDel                     uint8     `gorm:"column:is_del;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否删除，1：已删除，0：未删除"`
	CreateTime                time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime                time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyMaterialLibrary) TableName() string {
	return "company_material_library"
}

type companyMaterialLibraryRepo struct {
	data *Data
	log  *log.Helper
}

func (cml *CompanyMaterialLibrary) ToDomain() *domain.CompanyMaterialLibrary {
	companyMaterialLibrary := &domain.CompanyMaterialLibrary{
		Id:                        cml.Id,
		CompanyId:                 cml.CompanyId,
		UserId:                    cml.UserId,
		PrincipalUserId:           cml.PrincipalUserId,
		ParentId:                  cml.ParentId,
		LibraryName:               cml.LibraryName,
		LibraryType:               cml.LibraryType,
		LibraryUrl:                cml.LibraryUrl,
		LibraryMengmaAudit:        cml.LibraryMengmaAudit,
		LibraryMengmaAuditContent: cml.LibraryMengmaAuditContent,
		LiberaryFileType:          cml.LiberaryFileType,
		LiberaryFileSize:          cml.LiberaryFileSize,
		ExamineListType:           cml.ExamineListType,
		ProductId:                 cml.ProductId,
		VideoId:                   cml.VideoId,
		MaterialId:                cml.MaterialId,
		IsDel:                     cml.IsDel,
		CreateTime:                cml.CreateTime,
		UpdateTime:                cml.UpdateTime,
	}

	return companyMaterialLibrary
}

func NewCompanyMaterialLibraryRepo(data *Data, logger log.Logger) biz.CompanyMaterialLibraryRepo {
	return &companyMaterialLibraryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cmlr *companyMaterialLibraryRepo) GetByCompanyId(ctx context.Context, companyId, materialLibraryId uint64) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary := &CompanyMaterialLibrary{}

	if result := cmlr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("id = ?", materialLibraryId).
		Where("is_del = 0").
		First(companyMaterialLibrary); result.Error != nil {
		return nil, result.Error
	}

	return companyMaterialLibrary.ToDomain(), nil
}

func (cmlr *companyMaterialLibraryRepo) GetByCompanyIdAndLibraryName(ctx context.Context, companyId uint64, libraryName string) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary := &CompanyMaterialLibrary{}

	if result := cmlr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("library_name = ?", libraryName).
		Where("library_type = 1").
		Where("is_del = 0").
		First(companyMaterialLibrary); result.Error != nil {
		return nil, result.Error
	}

	return companyMaterialLibrary.ToDomain(), nil
}

func (cmlr *companyMaterialLibraryRepo) GetByCompanyIdAndParentIdAndLibraryName(ctx context.Context, companyId, parentId uint64, libraryName string) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary := &CompanyMaterialLibrary{}

	if result := cmlr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("parent_id = ?", parentId).
		Where("library_name = ?", libraryName).
		Where("library_type = 1").
		Where("is_del = 0").
		Order("create_time asc").
		First(companyMaterialLibrary); result.Error != nil {
		return nil, result.Error
	}

	return companyMaterialLibrary.ToDomain(), nil
}

func (cmlr *companyMaterialLibraryRepo) GetByCompanyIdAndParentIdAndProductId(ctx context.Context, companyId, parentId, productId uint64) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary := &CompanyMaterialLibrary{}

	if result := cmlr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("parent_id = ?", parentId).
		Where("product_id = ?", productId).
		Where("library_type = 1").
		Where("is_del = 0").
		First(companyMaterialLibrary); result.Error != nil {
		return nil, result.Error
	}

	return companyMaterialLibrary.ToDomain(), nil
}

func (cmlr *companyMaterialLibraryRepo) ListByParentIdAndLibraryType(ctx context.Context, companyId, id uint64, libraryType uint8) ([]*domain.CompanyMaterialLibrary, error) {
	var companyMaterialLibraries []CompanyMaterialLibrary
	list := make([]*domain.CompanyMaterialLibrary, 0)

	if result := cmlr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("parent_id = ?", id).
		Where("library_type = ?", libraryType).
		Where("is_del = 0").
		Order("create_time DESC,id DESC").
		Find(&companyMaterialLibraries); result.Error != nil {
		return nil, result.Error
	}

	for _, companyMaterialLibrary := range companyMaterialLibraries {
		list = append(list, companyMaterialLibrary.ToDomain())
	}

	return list, nil
}

func (cmlr *companyMaterialLibraryRepo) ListByParentId(ctx context.Context, companyId, id uint64) ([]*domain.CompanyMaterialLibrary, error) {
	var companyMaterialLibraries []CompanyMaterialLibrary
	list := make([]*domain.CompanyMaterialLibrary, 0)

	if result := cmlr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("parent_id = ?", id).
		Where("is_del = 0").
		Order("create_time ASC,id ASC").
		Find(&companyMaterialLibraries); result.Error != nil {
		return nil, result.Error
	}

	for _, companyMaterialLibrary := range companyMaterialLibraries {
		list = append(list, companyMaterialLibrary.ToDomain())
	}

	return list, nil
}

func (cmlr *companyMaterialLibraryRepo) CountByParentIdAndLibraryType(ctx context.Context, id uint64, libraryType uint8) (int64, error) {
	var count int64

	if result := cmlr.data.db.WithContext(ctx).
		Where("parent_id = ?", id).
		Where("library_type = ?", libraryType).
		Where("is_del = 0").
		Model(&CompanyMaterialLibrary{}).
		Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cmlr *companyMaterialLibraryRepo) Save(ctx context.Context, in *domain.CompanyMaterialLibrary) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary := &CompanyMaterialLibrary{
		CompanyId:                 in.CompanyId,
		UserId:                    in.UserId,
		PrincipalUserId:           in.PrincipalUserId,
		ParentId:                  in.ParentId,
		LibraryName:               in.LibraryName,
		LibraryType:               in.LibraryType,
		LibraryUrl:                in.LibraryUrl,
		LibraryMengmaAudit:        in.LibraryMengmaAudit,
		LibraryMengmaAuditContent: in.LibraryMengmaAuditContent,
		LiberaryFileType:          in.LiberaryFileType,
		LiberaryFileSize:          in.LiberaryFileSize,
		ExamineListType:           in.ExamineListType,
		ProductId:                 in.ProductId,
		VideoId:                   in.VideoId,
		MaterialId:                in.MaterialId,
		IsDel:                     in.IsDel,
		CreateTime:                in.CreateTime,
		UpdateTime:                in.UpdateTime,
	}

	if result := cmlr.data.DB(ctx).Create(companyMaterialLibrary); result.Error != nil {
		return nil, result.Error
	}

	return companyMaterialLibrary.ToDomain(), nil
}

func (cmlr *companyMaterialLibraryRepo) Update(ctx context.Context, in *domain.CompanyMaterialLibrary) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary := &CompanyMaterialLibrary{
		Id:                        in.Id,
		CompanyId:                 in.CompanyId,
		UserId:                    in.UserId,
		PrincipalUserId:           in.PrincipalUserId,
		ParentId:                  in.ParentId,
		LibraryName:               in.LibraryName,
		LibraryType:               in.LibraryType,
		LibraryUrl:                in.LibraryUrl,
		LibraryMengmaAudit:        in.LibraryMengmaAudit,
		LibraryMengmaAuditContent: in.LibraryMengmaAuditContent,
		LiberaryFileType:          in.LiberaryFileType,
		LiberaryFileSize:          in.LiberaryFileSize,
		ExamineListType:           in.ExamineListType,
		ProductId:                 in.ProductId,
		VideoId:                   in.VideoId,
		MaterialId:                in.MaterialId,
		IsDel:                     in.IsDel,
		CreateTime:                in.CreateTime,
		UpdateTime:                in.UpdateTime,
	}

	if result := cmlr.data.DB(ctx).Save(companyMaterialLibrary); result.Error != nil {
		return nil, result.Error
	}

	return companyMaterialLibrary.ToDomain(), nil
}
