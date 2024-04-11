package biz

import (
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"encoding/base64"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	CompanyOrganizationNotFound             = errors.NotFound("COMPANY_ORGANIZATION_NOT_FOUND", "企业机构不存在")
	CompanyOrganizationCreateError          = errors.InternalServer("COMPANY_ORGANIZATION_CREATE_ERROR", "企业机构创建失败")
	CompanyOrganizationCompanyUserError     = errors.InternalServer("COMPANY_ORGANIZATION_COMPANY_USER_ERROR", "该手机号已经是系统用户，暂不能做机构管理员")
	CompanyOrganizationCompanyUserListError = errors.InternalServer("COMPANY_ORGANIZATION_COMPANY_USER_LIST_ERROR", "企业机构关联用户列表获取失败")
	CompanyOrganizationUpdateError          = errors.InternalServer("COMPANY_ORGANIZATION_UPDATE_ERROR", "企业机构更新失败")
	CompanyOrganizationListError            = errors.InternalServer("COMPANY_ORGANIZATION_LIST_ERROR", "企业机构列表获取失败")
	CompanyOrganizationDeleteError          = errors.InternalServer("COMPANY_ORGANIZATION_DELETE_ERROR", "企业机构删除失败")

	Mime = map[string]string{
		"image/jpeg": ".jpeg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}
)

type CompanyOrganizationRepo interface {
	GetById(context.Context, uint64) (*domain.CompanyOrganization, error)
	GetByOrganizationCode(context.Context, string) (*domain.CompanyOrganization, error)
	List(context.Context, int, int) ([]*domain.CompanyOrganization, error)
	Count(context.Context) (int64, error)
	Save(context.Context, *domain.CompanyOrganization) (*domain.CompanyOrganization, error)
	Update(context.Context, *domain.CompanyOrganization) (*domain.CompanyOrganization, error)

	PutContent(context.Context, string, io.Reader) (*ctos.PutObjectV2Output, error)
}

type CompanyOrganizationUsecase struct {
	repo   CompanyOrganizationRepo
	crepo  CompanyRepo
	curepo CompanyUserRepo
	qcrepo QrCodeRepo
	surepo ShortUrlRepo
	screpo ShortCodeRepo
	tm     Transaction
	conf   *conf.Data
	cconf  *conf.Company
	vconf  *conf.Volcengine
	log    *log.Helper
}

func NewCompanyOrganizationUsecase(repo CompanyOrganizationRepo, crepo CompanyRepo, curepo CompanyUserRepo, qcrepo QrCodeRepo, surepo ShortUrlRepo, screpo ShortCodeRepo, tm Transaction, conf *conf.Data, cconf *conf.Company, vconf *conf.Volcengine, logger log.Logger) *CompanyOrganizationUsecase {
	return &CompanyOrganizationUsecase{repo: repo, crepo: crepo, curepo: curepo, qcrepo: qcrepo, surepo: surepo, screpo: screpo, tm: tm, conf: conf, cconf: cconf, vconf: vconf, log: log.NewHelper(logger)}
}

func (couc *CompanyOrganizationUsecase) ListCompanyOrganizations(ctx context.Context, pageNum, pageSize uint64) (*domain.CompanyOrganizationList, error) {
	list := make([]*domain.CompanyOrganization, 0)

	companyOrganizations, err := couc.repo.List(ctx, int(pageNum), int(pageSize))

	if err != nil {
		return nil, CompanyOrganizationListError
	}

	for _, lcompanyOrganization := range companyOrganizations {
		companyOrganization, _ := couc.getCompanyOrganization(ctx, lcompanyOrganization)

		list = append(list, companyOrganization)
	}

	total, err := couc.repo.Count(ctx)

	if err != nil {
		return nil, CompanyOrganizationListError
	}

	return &domain.CompanyOrganizationList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (couc *CompanyOrganizationUsecase) ListSelectCompanyOrganizations(ctx context.Context) (*domain.SelectCompanyOrganizations, error) {
	selectCompanyOrganizations := domain.NewSelectCompanyOrganizations()

	selectCompanyOrganizations.CourseLevel = append(selectCompanyOrganizations.CourseLevel, &domain.CourseLevel{
		Key:   "1",
		Value: "零级",
	})

	selectCompanyOrganizations.CourseLevel = append(selectCompanyOrganizations.CourseLevel, &domain.CourseLevel{
		Key:   "2",
		Value: "初级",
	})

	selectCompanyOrganizations.CourseLevel = append(selectCompanyOrganizations.CourseLevel, &domain.CourseLevel{
		Key:   "3",
		Value: "中级",
	})

	selectCompanyOrganizations.CourseLevel = append(selectCompanyOrganizations.CourseLevel, &domain.CourseLevel{
		Key:   "4",
		Value: "高级",
	})

	selectCompanyOrganizations.OrganizationMcn = append(selectCompanyOrganizations.OrganizationMcn, &domain.OrganizationMcn{
		Key:   "星达小当家",
		Value: "星达小当家",
	})

	selectCompanyOrganizations.OrganizationMcn = append(selectCompanyOrganizations.OrganizationMcn, &domain.OrganizationMcn{
		Key:   "小当家",
		Value: "小当家",
	})

	selectCompanyOrganizations.OrganizationMcn = append(selectCompanyOrganizations.OrganizationMcn, &domain.OrganizationMcn{
		Key:   "莲宝兔旺",
		Value: "莲宝兔旺",
	})

	selectCompanyOrganizations.OrganizationMcn = append(selectCompanyOrganizations.OrganizationMcn, &domain.OrganizationMcn{
		Key:   "壹玖传媒",
		Value: "壹玖传媒",
	})

	return selectCompanyOrganizations, nil
}

func (couc *CompanyOrganizationUsecase) GetCompanyOrganizations(ctx context.Context, organizationId uint64) (*domain.CompanyOrganization, error) {
	inCompanyOrganization, err := couc.repo.GetById(ctx, organizationId)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	companyOrganization, err := couc.getCompanyOrganization(ctx, inCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) GetCompanyOrganizationByOrganizationCodes(ctx context.Context, organizationCode string) (*domain.CompanyOrganization, error) {
	inCompanyOrganization, err := couc.repo.GetByOrganizationCode(ctx, organizationCode)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	companyOrganization, err := couc.getCompanyOrganization(ctx, inCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) CreateCompanyOrganizations(ctx context.Context, organizationName, organizationLogo, organizationMcn, companyName, bankCode, bankDeposit string) (*domain.CompanyOrganization, error) {
	shortCode, err := couc.screpo.Create(ctx)

	if err != nil {
		return nil, CompanyShortCodeCreateError
	}

	sorganizationLogo := strings.Split(organizationLogo, ",")

	if len(sorganizationLogo) != 2 {
		return nil, CompanyOrganizationCreateError
	}

	if _, ok := Mime[sorganizationLogo[0][5:len(sorganizationLogo[0])-7]]; !ok {
		return nil, CompanyOrganizationCreateError
	}

	imagePath := couc.vconf.Tos.Organization.SubFolder + "/" + tool.GetRandCode(time.Now().String()) + Mime[sorganizationLogo[0][5:len(sorganizationLogo[0])-7]]
	imageContent, err := base64.StdEncoding.DecodeString(sorganizationLogo[1])

	if err != nil {
		return nil, CompanyOrganizationCreateError
	}

	if _, err = couc.repo.PutContent(ctx, imagePath, strings.NewReader(string(imageContent))); err != nil {
		return nil, CompanyOrganizationCreateError
	}

	inCompanyOrganization := domain.NewCompanyOrganization(ctx, organizationName, organizationMcn, companyName, bankCode, bankDeposit, couc.vconf.Tos.Organization.Url+"/"+imagePath, shortCode.Data.ShortCode, "", "", "", "", "")
	inCompanyOrganization.SetCreateTime(ctx)
	inCompanyOrganization.SetUpdateTime(ctx)

	if ok := inCompanyOrganization.VerifyOrganizationMcn(ctx); !ok {
		return nil, CompanyOrganizationCompanyUserError
	}

	var tmpCompanyOrganization *domain.CompanyOrganization

	err = couc.tm.InTx(ctx, func(ctx context.Context) error {
		var err error

		tmpCompanyOrganization, err = couc.repo.Save(ctx, inCompanyOrganization)

		if err != nil {
			return err
		}

		/*organizationShortUrl, err := couc.surepo.Create(ctx, "oId="+strconv.FormatUint(tmpCompanyOrganization.Id, 10))

		if err != nil {
			return err
		}

		tmpCompanyOrganization.SetOrganizationShortUrl(ctx, organizationShortUrl.Data.ShortUrl)

		if _, err = couc.repo.Update(ctx, tmpCompanyOrganization); err != nil {
			return err
		}*/

		return nil
	})

	if err != nil {
		return nil, CompanyOrganizationCreateError
	}

	companyOrganization, err := couc.getCompanyOrganization(ctx, tmpCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateCompanyOrganizations(ctx context.Context, organizationId uint64, organizationName, organizationLogo, organizationMcn, companyName, bankCode, bankDeposit string) (*domain.CompanyOrganization, error) {
	inCompanyOrganization, err := couc.repo.GetById(ctx, organizationId)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	if len(organizationLogo) > 0 {
		sorganizationLogo := strings.Split(organizationLogo, ",")

		if len(sorganizationLogo) != 2 {
			return nil, CompanyOrganizationUpdateError
		}

		if _, ok := Mime[sorganizationLogo[0][5:len(sorganizationLogo[0])-7]]; !ok {
			return nil, CompanyOrganizationUpdateError
		}

		imagePath := couc.vconf.Tos.Organization.SubFolder + "/" + tool.GetRandCode(time.Now().String()) + Mime[sorganizationLogo[0][5:len(sorganizationLogo[0])-7]]
		imageContent, err := base64.StdEncoding.DecodeString(sorganizationLogo[1])

		if err != nil {
			return nil, CompanyOrganizationUpdateError
		}

		if _, err = couc.repo.PutContent(ctx, imagePath, strings.NewReader(string(imageContent))); err != nil {
			return nil, CompanyOrganizationUpdateError
		}

		inCompanyOrganization.SetOrganizationLogoUrl(ctx, couc.vconf.Tos.Organization.Url+"/"+imagePath)
	}

	inCompanyOrganization.SetOrganizationMcn(ctx, organizationMcn)
	inCompanyOrganization.SetCompanyName(ctx, companyName)
	inCompanyOrganization.SetBankCode(ctx, bankCode)
	inCompanyOrganization.SetBankDeposit(ctx, bankDeposit)
	inCompanyOrganization.SetOrganizationName(ctx, organizationName)
	inCompanyOrganization.SetUpdateTime(ctx)

	if ok := inCompanyOrganization.VerifyOrganizationMcn(ctx); !ok {
		return nil, CompanyOrganizationUpdateError
	}

	var tmpCompanyOrganization *domain.CompanyOrganization

	err = couc.tm.InTx(ctx, func(ctx context.Context) error {
		var err error

		tmpCompanyOrganization, err = couc.repo.Update(ctx, inCompanyOrganization)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, CompanyOrganizationUpdateError
	}

	companyOrganization, err := couc.getCompanyOrganization(ctx, tmpCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateTeamCompanyOrganizations(ctx context.Context, organizationId uint64, organizationUser string) (*domain.CompanyOrganization, error) {
	inCompanyOrganization, err := couc.repo.GetById(ctx, organizationId)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	inCompanyOrganization.SetOrganizationUser(ctx, organizationUser)
	inCompanyOrganization.SetUpdateTime(ctx)

	if ok := inCompanyOrganization.VerifyOrganizationUser(ctx); !ok {
		return nil, CompanyOrganizationUpdateError
	}

	companyOrganizationUsers, err := couc.curepo.ListByCompanyIdAndOrganizationId(ctx, couc.cconf.DefaultCompanyId, organizationId)

	if err != nil {
		return nil, CompanyOrganizationCompanyUserListError
	}

	deleteCompanyUserIds := make([]uint64, 0)

	for _, companyOrganizationUser := range companyOrganizationUsers {
		isNotExist := true

		for _, lorganizationUser := range inCompanyOrganization.OrganizationUsers {
			if lorganizationUser.UserId > 0 {
				if lorganizationUser.UserId == companyOrganizationUser.Id {
					isNotExist = false

					break
				}
			}
		}

		if isNotExist {
			deleteCompanyUserIds = append(deleteCompanyUserIds, companyOrganizationUser.Id)
		}
	}

	for _, lorganizationUser := range inCompanyOrganization.OrganizationUsers {
		if lorganizationUser.UserId > 0 {
			isNotExist := true

			for _, companyOrganizationUser := range companyOrganizationUsers {
				if companyOrganizationUser.Id == lorganizationUser.UserId {
					if companyOrganizationUser.Phone != lorganizationUser.Phone {
						if _, err := couc.curepo.GetByPhone(ctx, couc.cconf.DefaultCompanyId, lorganizationUser.Phone); err == nil {
							return nil, CompanyOrganizationCompanyUserError
						}
					}

					isNotExist = false

					break
				}
			}

			if isNotExist {
				return nil, CompanyOrganizationCompanyUserError
			}
		} else {
			if _, err := couc.curepo.GetByPhoneAndNotInUserId(ctx, couc.cconf.DefaultCompanyId, lorganizationUser.Phone, deleteCompanyUserIds); err == nil {
				return nil, CompanyOrganizationCompanyUserError
			}
		}
	}

	var tmpCompanyOrganization *domain.CompanyOrganization

	err = couc.tm.InTx(ctx, func(ctx context.Context) error {
		var err error

		tmpCompanyOrganization, err = couc.repo.Update(ctx, inCompanyOrganization)

		if err != nil {
			return err
		}

		if len(deleteCompanyUserIds) > 0 {
			if err = couc.curepo.DeleteByCompanyIdAndUserId(ctx, couc.cconf.DefaultCompanyId, deleteCompanyUserIds); err != nil {
				return err
			}
		}

		for _, lorganizationUser := range inCompanyOrganization.OrganizationUsers {
			if lorganizationUser.UserId > 0 {
				var inCompanyUser *domain.CompanyUser

				for _, companyOrganizationUser := range companyOrganizationUsers {
					if companyOrganizationUser.Id == lorganizationUser.UserId {
						inCompanyUser = companyOrganizationUser

						break
					}
				}

				inCompanyUser.SetPhone(ctx, lorganizationUser.Phone)
				inCompanyUser.SetUsername(ctx, lorganizationUser.Username)
				inCompanyUser.SetOrganizationId(ctx, organizationId)
				inCompanyUser.SetUpdateTime(ctx)

				if _, err = couc.curepo.Update(ctx, inCompanyUser); err != nil {
					return err
				}
			} else {
				inCompanyUser := domain.NewCompanyUser(ctx, couc.cconf.DefaultCompanyId, lorganizationUser.Username, "", lorganizationUser.Phone, 3, 1)
				inCompanyUser.SetOrganizationId(ctx, tmpCompanyOrganization.Id)
				inCompanyUser.SetCreateTime(ctx)
				inCompanyUser.SetUpdateTime(ctx)

				if _, err = couc.curepo.Save(ctx, inCompanyUser); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, CompanyOrganizationUpdateError
	}

	companyOrganization, err := couc.getCompanyOrganization(ctx, tmpCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateCommissionCompanyOrganizations(ctx context.Context, organizationId uint64, organizationCommission string) (*domain.CompanyOrganization, error) {
	inCompanyOrganization, err := couc.repo.GetById(ctx, organizationId)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	inCompanyOrganization.SetOrganizationCommission(ctx, organizationCommission)
	inCompanyOrganization.SetUpdateTime(ctx)

	if ok := inCompanyOrganization.VerifyOrganizationCommission(ctx); !ok {
		return nil, CompanyOrganizationUpdateError
	}

	tmpCompanyOrganization, err := couc.repo.Update(ctx, inCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationUpdateError
	}

	companyOrganization, err := couc.getCompanyOrganization(ctx, tmpCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateColonelCommissionCompanyOrganizations(ctx context.Context, organizationId uint64, organizationColonelCommission string) (*domain.CompanyOrganization, error) {
	inCompanyOrganization, err := couc.repo.GetById(ctx, organizationId)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	inCompanyOrganization.SetOrganizationColonelCommission(ctx, organizationColonelCommission)
	inCompanyOrganization.SetUpdateTime(ctx)

	if ok := inCompanyOrganization.VerifyOrganizationColonelCommission(ctx); !ok {
		return nil, CompanyOrganizationUpdateError
	}

	tmpCompanyOrganization, err := couc.repo.Update(ctx, inCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationUpdateError
	}

	companyOrganization, err := couc.getCompanyOrganization(ctx, tmpCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) UpdateCourseCompanyOrganizations(ctx context.Context, organizationId uint64, organizationCourse string) (*domain.CompanyOrganization, error) {
	inCompanyOrganization, err := couc.repo.GetById(ctx, organizationId)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	inCompanyOrganization.SetOrganizationCourse(ctx, organizationCourse)
	inCompanyOrganization.SetUpdateTime(ctx)

	if ok := inCompanyOrganization.VerifyOrganizationCourse(ctx); !ok {
		return nil, CompanyOrganizationUpdateError
	}

	tmpCompanyOrganization, err := couc.repo.Update(ctx, inCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationUpdateError
	}

	companyOrganization, err := couc.getCompanyOrganization(ctx, tmpCompanyOrganization)

	if err != nil {
		return nil, CompanyOrganizationNotFound
	}

	return companyOrganization, nil
}

func (couc *CompanyOrganizationUsecase) SyncUpdateQrCodeCompanyOrganizations(ctx context.Context) error {
	companyOrganizations, err := couc.repo.List(ctx, 0, 40)

	if err != nil {
		return CompanyOrganizationListError
	}

	for _, inCompanyOrganization := range companyOrganizations {
		if inCompanyOrganization.OrganizationQrCodeUrl == "" {
			qrCodeContent := "oId=" + strconv.FormatUint(inCompanyOrganization.Id, 10)

			qrCode, err := couc.qcrepo.Get(ctx, inCompanyOrganization.Id, qrCodeContent)

			if err != nil {
				return CompanyOrganizationUpdateError
			}

			inCompanyOrganization.SetOrganizationQrCodeUrl(ctx, qrCode.Data.StaticUrl)
			inCompanyOrganization.SetUpdateTime(ctx)

			if _, err := couc.repo.Update(ctx, inCompanyOrganization); err != nil {
				return CompanyOrganizationUpdateError
			}
		}
	}

	return nil
}

func (couc *CompanyOrganizationUsecase) getCompanyOrganization(ctx context.Context, companyOrganization *domain.CompanyOrganization) (*domain.CompanyOrganization, error) {
	companyOrganization.SetOrganizationMcns(ctx)
	companyOrganization.SetOrganizationCommissions(ctx)
	companyOrganization.SetOrganizationCourses(ctx)
	companyOrganization.SetOrganizationColonelCommissions(ctx)

	companyOrganization.OrganizationUsers = make([]*domain.OrganizationUser, 0)

	if companyUsers, err := couc.curepo.ListByCompanyIdAndOrganizationId(ctx, couc.cconf.DefaultCompanyId, companyOrganization.Id); err == nil {
		for _, companyUser := range companyUsers {
			companyOrganization.OrganizationUsers = append(companyOrganization.OrganizationUsers, &domain.OrganizationUser{
				UserId:   companyUser.Id,
				Username: companyUser.Username,
				Phone:    companyUser.Phone,
			})
		}
	}

	return companyOrganization, nil
}
