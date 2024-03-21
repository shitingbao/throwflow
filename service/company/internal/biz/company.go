package biz

import (
	v1 "company/api/service/common/v1"
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"time"
)

var (
	CompanyCompanyNotFound      = errors.NotFound("COMPANY_COMAPNY_NOT_FOUND", "企业不存在")
	CompanyClueFound            = errors.InternalServer("COMPANY_CLUE_FOUND", "线索已经开通过企业")
	CompanyCompanyCreateError   = errors.InternalServer("COMPANY_COMAPNY_CREATE_ERROR", "企业创建失败")
	CompanyCompanyUpdateError   = errors.InternalServer("COMPANY_COMAPNY_UPDATE_ERROR", "企业更新失败")
	CompanyCompanyDeleteError   = errors.InternalServer("COMPANY_COMAPNY_DELETE_ERROR", "企业删除失败")
	CompanyCompanyExitMainAdmin = errors.InternalServer("COMPANY_COMAPNY_EXIT_MAIN_ADMIN", "企业已存在主管理员")
	CompanyRoleMenuNotFound     = errors.InternalServer("COMPANY_ROLE_MENU_NOT_FOUND", "企业关联菜单不存在")
	CompanyCompanyNotUpdate     = errors.InternalServer("COMPANY_COMAPNY_NOT_UPDATE", "企业已过期不能更新状态")
	CompanyCompanyNotDelete     = errors.InternalServer("COMPANY_COMAPNY_DELETE_ERROR", "企业只有过期才能被移到线索")
)

type CompanyRepo interface {
	GetById(context.Context, uint64) (*domain.Company, error)
	GetByClueId(context.Context, uint64) (*domain.Company, error)
	List(context.Context, int, int, uint64, string, string, uint8) ([]*domain.Company, error)
	Count(context.Context, uint64, string, string, uint8) (int64, error)
	Statistics(context.Context, uint8) (int64, error)
	Save(context.Context, *domain.Company) (*domain.Company, error)
	Update(context.Context, *domain.Company) (*domain.Company, error)
	Delete(context.Context, *domain.Company) error
}

type CompanyUsecase struct {
	repo    CompanyRepo
	crepo   ClueRepo
	irepo   IndustryRepo
	curepo  CompanyUserRepo
	arepo   AreaRepo
	qarepo  QianchuanAdvertiserRepo
	cprrepo CompanyPerformanceRuleRepo
	csrepo  CompanySetRepo
	tm      Transaction
	conf    *conf.Data
	econf   *conf.Event
	log     *log.Helper
}

func NewCompanyUsecase(repo CompanyRepo, crepo ClueRepo, irepo IndustryRepo, curepo CompanyUserRepo, arepo AreaRepo, qarepo QianchuanAdvertiserRepo, cprrepo CompanyPerformanceRuleRepo, csrepo CompanySetRepo, tm Transaction, conf *conf.Data, econf *conf.Event, logger log.Logger) *CompanyUsecase {
	return &CompanyUsecase{repo: repo, crepo: crepo, irepo: irepo, curepo: curepo, arepo: arepo, qarepo: qarepo, cprrepo: cprrepo, csrepo: csrepo, tm: tm, conf: conf, econf: econf, log: log.NewHelper(logger)}
}

func (cuc *CompanyUsecase) ListCompanys(ctx context.Context, pageNum, pageSize, industryId uint64, keyword, status string, companyType uint8) (*domain.CompanyList, error) {
	companys, err := cuc.repo.List(ctx, int(pageNum), int(pageSize), industryId, keyword, status, companyType)

	if err != nil {
		return nil, CompanyDataError
	}

	list := make([]*domain.Company, 0)

	for _, company := range companys {
		company, _ = cuc.getCompanyById(ctx, company.Id)

		list = append(list, company)
	}

	total, err := cuc.repo.Count(ctx, industryId, keyword, status, companyType)

	if err != nil {
		return nil, CompanyDataError
	}

	return &domain.CompanyList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (cuc *CompanyUsecase) ListSelectCompanys(ctx context.Context) (*domain.SelectCompanys, error) {
	return domain.NewSelectCompanys(), nil
}

func (cuc *CompanyUsecase) StatisticsCompanys(ctx context.Context) (*domain.StatisticsCompanys, error) {
	selects := domain.NewSelectCompanys()

	statistics := make([]*domain.StatisticsCompany, 0)

	for _, companyType := range selects.CompanyType {
		iCompanyType, _ := strconv.Atoi(companyType.Key)

		count, _ := cuc.repo.Statistics(ctx, uint8(iCompanyType))

		statistics = append(statistics, &domain.StatisticsCompany{
			Key:   companyType.Value,
			Value: strconv.FormatInt(count, 10),
		})
	}

	return &domain.StatisticsCompanys{
		Statistics: statistics,
	}, nil
}

func (cuc *CompanyUsecase) GetCompanys(ctx context.Context, id uint64) (*domain.Company, error) {
	company, err := cuc.getCompanyById(ctx, id)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	return company, nil
}

func (cuc *CompanyUsecase) CreateCompanys(ctx context.Context, companyName, contactInformation, source, seller, facilitator, adminName, adminPhone, address, industryId string, userId, clueId, areaCode uint64, companyType, qianchuanUse, status uint8) (*domain.Company, error) {
	var inClue *domain.Clue

	if clueId > 0 {
		var err error

		inClue, err = cuc.crepo.GetById(ctx, clueId)

		if err != nil {
			return nil, CompanyClueNotFound
		}

		if inCompany, err := cuc.repo.GetByClueId(ctx, clueId); err == nil {
			if inCompany.IsDel == 0 {
				return nil, CompanyClueFound
			} else {
				inCompany.SetUpdateTime(ctx)
				inCompany.SetIsDel(ctx, 0)

				if company, err := cuc.repo.Update(ctx, inCompany); err != nil {
					return nil, CompanyClueCreateError
				} else {
					company, err = cuc.getCompanyById(ctx, company.Id)

					if err != nil {
						return nil, CompanyCompanyNotFound
					}

					return company, nil
				}
			}
		}
	}

	industryIds := strings.Split(industryId, ",")
	industryIds = tool.RemoveEmptyString(industryIds)

	sIndustryIds := make([]string, 0)
	sIndustryNames := make([]string, 0)

	for _, sIndustryId := range industryIds {
		if iIndustryId, err := strconv.ParseUint(sIndustryId, 10, 64); err == nil {
			if industry, err := cuc.irepo.GetById(ctx, iIndustryId); err == nil {
				sIndustryIds = append(sIndustryIds, strconv.FormatUint(industry.Id, 10))
				sIndustryNames = append(sIndustryNames, industry.IndustryName)
			}
		}
	}

	/*if len(sIndustryIds) == 0 {
		return nil, CompanyIndustryNotFound
	}*/

	var area *v1.GetAreasReply
	var err error

	if areaCode > 0 {
		area, err = cuc.arepo.GetByAreaCode(ctx, areaCode)

		if err != nil {
			return nil, err
		}
	}

	if clueId == 0 {
		inClue = domain.NewClue(ctx, companyName, contactInformation, source, seller, facilitator, address, industryId, area.Data.AreaCode, companyType, qianchuanUse, status)
		inClue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
		inClue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))
		inClue.SetCreateTime(ctx)
		inClue.SetUpdateTime(ctx)

		inClue.SetOperationLog(ctx, userId, "", "创建", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
	} else {
		inClue.SetCompanyName(ctx, companyName)
		inClue.SetCompanyType(ctx, companyType)
		inClue.SetQianchuanUse(ctx, qianchuanUse)
		inClue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
		inClue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))
		inClue.SetContactInformation(ctx, contactInformation)
		inClue.SetSeller(ctx, seller)
		inClue.SetFacilitator(ctx, facilitator)
		inClue.SetStatus(ctx, status)
		inClue.SetAreaCode(ctx, area.Data.AreaCode)
		inClue.SetAddress(ctx, address)

		if area != nil {
			inClue.SetAreaName(ctx, area.Data.AreaName)
		}

		inClue.SetUpdateTime(ctx)
	}

	if ok := inClue.VerifyContactInformation(ctx); !ok {
		return nil, CompanyDataError
	}

	var company *domain.Company

	err = cuc.tm.InTx(ctx, func(ctx context.Context) error {
		var clue *domain.Clue
		var err error

		if inClue.Id > 0 {
			clue, err = cuc.crepo.Update(ctx, inClue)

			if err != nil {
				return CompanyClueUpdateError
			}
		} else {
			clue, err = cuc.crepo.Save(ctx, inClue)

			if err != nil {
				return CompanyClueCreateError
			}
		}

		inCompany := domain.NewCompany(ctx, clue.Id)
		inCompany.SetCreateTime(ctx)
		inCompany.SetUpdateTime(ctx)
		inCompany.SetStartTime(ctx, time.Now())
		inCompany.SetEndTime(ctx, time.Now())

		company, err = cuc.repo.Save(ctx, inCompany)

		if err != nil {
			return CompanyCompanyCreateError
		}

		inCompanyUser := domain.NewCompanyUser(ctx, company.Id, adminName, "", adminPhone, 1, 1)
		inCompanyUser.SetCreateTime(ctx)
		inCompanyUser.SetUpdateTime(ctx)

		if _, err := cuc.curepo.Save(ctx, inCompanyUser); err != nil {
			return CompanyCompanyUserCreateError
		}

		inCompanyPerformanceRuleA := domain.NewCompanyPerformanceRule(ctx, company.Id, "提升绩效方案A", "")
		inCompanyPerformanceRuleA.SetCreateTime(ctx)
		inCompanyPerformanceRuleA.SetUpdateTime(ctx)

		if _, err := cuc.cprrepo.Save(ctx, inCompanyPerformanceRuleA); err != nil {
			return CompanyCompanyPerformanceRuleCreateError
		}

		inCompanyPerformanceRuleB := domain.NewCompanyPerformanceRule(ctx, company.Id, "提升绩效方案B", "")
		inCompanyPerformanceRuleB.SetCreateTime(ctx)
		inCompanyPerformanceRuleB.SetUpdateTime(ctx)

		if _, err := cuc.cprrepo.Save(ctx, inCompanyPerformanceRuleB); err != nil {
			return CompanyCompanyPerformanceRuleCreateError
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	company, err = cuc.getCompanyById(ctx, company.Id)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	return company, nil
}

func (cuc *CompanyUsecase) UpdateCompanys(ctx context.Context, id uint64, companyName, contactInformation, seller, facilitator, adminName, adminPhone, address, industryId string, clueId, areaCode uint64, companyType, qianchuanUse, status uint8) (*domain.Company, error) {
	inClue, err := cuc.crepo.GetById(ctx, clueId)

	if err != nil {
		return nil, CompanyClueNotFound
	}

	industryIds := strings.Split(industryId, ",")
	industryIds = tool.RemoveEmptyString(industryIds)

	sIndustryIds := make([]string, 0)
	sIndustryNames := make([]string, 0)

	for _, sIndustryId := range industryIds {
		if iIndustryId, err := strconv.ParseUint(sIndustryId, 10, 64); err == nil {
			if industry, err := cuc.irepo.GetById(ctx, iIndustryId); err == nil {
				sIndustryIds = append(sIndustryIds, strconv.FormatUint(industry.Id, 10))
				sIndustryNames = append(sIndustryNames, industry.IndustryName)
			}
		}
	}

	/*if len(sIndustryIds) == 0 {
		return nil, CompanyIndustryNotFound
	}*/

	var area *v1.GetAreasReply

	if areaCode > 0 {
		area, err = cuc.arepo.GetByAreaCode(ctx, areaCode)

		if err != nil {
			return nil, err
		}
	}

	inClue.SetCompanyName(ctx, companyName)
	inClue.SetCompanyType(ctx, companyType)
	inClue.SetQianchuanUse(ctx, qianchuanUse)
	inClue.SetAddress(ctx, address)
	inClue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
	inClue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))
	inClue.SetContactInformation(ctx, contactInformation)
	inClue.SetSeller(ctx, seller)
	inClue.SetFacilitator(ctx, facilitator)
	inClue.SetStatus(ctx, status)
	inClue.SetAreaCode(ctx, area.Data.AreaCode)

	if area != nil {
		inClue.SetAreaName(ctx, area.Data.AreaName)
	}

	inClue.SetUpdateTime(ctx)

	if ok := inClue.VerifyContactInformation(ctx); !ok {
		return nil, CompanyDataError
	}

	inCompany, err := cuc.getCompanyById(ctx, id)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	isExitMainAdmin := false
	var inCompanyUser *domain.CompanyUser

	for _, companyUser := range inCompany.CompanyUser {
		if companyUser.Role == 1 {
			isExitMainAdmin = true
			inCompanyUser = companyUser
			break
		}
	}

	if isExitMainAdmin {
		if inCompanyUser.Phone != adminPhone {
			return nil, CompanyCompanyExitMainAdmin
		}
	}

	err = cuc.tm.InTx(ctx, func(ctx context.Context) error {
		if _, err := cuc.crepo.Update(ctx, inClue); err != nil {
			return CompanyClueUpdateError
		}

		inCompany.SetUpdateTime(ctx)

		if _, err := cuc.repo.Update(ctx, inCompany); err != nil {
			return CompanyCompanyCreateError
		}

		inCompanyUser.SetUsername(ctx, adminName)
		inCompanyUser.SetUpdateTime(ctx)

		if _, err = cuc.curepo.Update(ctx, inCompanyUser); err != nil {
			return CompanyCompanyUserCreateError
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	company, err := cuc.getCompanyById(ctx, inCompany.Id)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	return company, nil
}

func (cuc *CompanyUsecase) UpdateStatusCompanys(ctx context.Context, id uint64, status uint8) (*domain.Company, error) {
	inCompany, err := cuc.getCompanyById(ctx, id)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	if inCompany.Status == 2 {
		return nil, CompanyCompanyNotUpdate
	}

	inCompany.SetStatus(ctx, status)
	inCompany.SetUpdateTime(ctx)

	if _, err := cuc.repo.Update(ctx, inCompany); err != nil {
		return nil, CompanyCompanyUpdateError
	}

	company, err := cuc.getCompanyById(ctx, inCompany.Id)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	if status == 0 {
		cuc.qarepo.UpdateStatusByCompanyId(ctx, company.Id, 0)
	}

	return company, nil
}

func (cuc *CompanyUsecase) UpdateRoleCompanys(ctx context.Context, id, userId uint64, ids []string, startTime, endTime time.Time, accounts, qianchuanAdvertisers uint32, companyType, isTermwork uint8) (*domain.Company, error) {
	inCompany, err := cuc.getCompanyById(ctx, id)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	oldAccounts := inCompany.Accounts
	oldQianchuanAdvertisers := inCompany.QianchuanAdvertisers

	inCompany.SetAccounts(ctx, accounts)
	inCompany.SetQianchuanAdvertisers(ctx, qianchuanAdvertisers)
	inCompany.SetCompanyType(ctx, companyType)
	inCompany.SetIsTermwork(ctx, isTermwork)
	inCompany.SetMenuId(ctx, strings.Join(ids, ","))
	inCompany.SetStartTime(ctx, startTime)
	inCompany.SetEndTime(ctx, endTime)
	inCompany.SetUpdateTime(ctx)

	if inCompany.Status == 2 {
		inCompany.SetStatus(ctx, 0)
	}

	if _, err := cuc.repo.Update(ctx, inCompany); err != nil {
		return nil, CompanyCompanyUpdateError
	}

	company, err := cuc.getCompanyById(ctx, inCompany.Id)

	if company.Status == 2 || company.Status == 0 {
		if inClue, err := cuc.crepo.GetById(ctx, inCompany.ClueId); err == nil {
			if companyType == 1 {
				inClue.SetOperationLog(ctx, userId, "", "试用版开通", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
			} else if companyType == 2 {
				inClue.SetOperationLog(ctx, userId, "", "基础版开通", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
			} else if companyType == 3 {
				inClue.SetOperationLog(ctx, userId, "", "专业版开通", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
			} else if companyType == 4 {
				inClue.SetOperationLog(ctx, userId, "", "旗舰版开通", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
			} else if companyType == 5 {
				inClue.SetOperationLog(ctx, userId, "", "尊享版开通", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
			}

			cuc.crepo.Update(ctx, inClue)
		}
	}

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	if accounts < oldAccounts {
		if companyUsers, err := cuc.curepo.ListByCompanyId(ctx, company.Id); err == nil {
			if len(companyUsers) > int(accounts) {
				outAccounts := len(companyUsers) - int(accounts)

				for i := 1; i <= outAccounts; i++ {
					key := len(companyUsers) - i

					if key >= 0 {
						companyUsers[key].SetStatus(ctx, 0)

						cuc.curepo.Update(ctx, companyUsers[key])
					}
				}
			}
		}
	}

	if qianchuanAdvertisers < oldQianchuanAdvertisers {
		if dqianchuanAdvertisers, err := cuc.qarepo.List(ctx, company.Id, 0, 0, "", "", ""); err == nil {
			if len(dqianchuanAdvertisers.Data.List) > int(qianchuanAdvertisers) {
				outQianchuanAdvertisers := len(dqianchuanAdvertisers.Data.List) - int(qianchuanAdvertisers)

				for i := 1; i <= outQianchuanAdvertisers; i++ {
					key := len(dqianchuanAdvertisers.Data.List) - i

					if key >= 0 {
						cuc.qarepo.Update(ctx, dqianchuanAdvertisers.Data.List[key].CompanyId, dqianchuanAdvertisers.Data.List[key].AdvertiserId, 0)
					}
				}
			}
		}
	}

	return company, nil
}

func (cuc *CompanyUsecase) DeleteCompanys(ctx context.Context, id uint64) error {
	inCompany, err := cuc.getCompanyById(ctx, id)

	if err != nil {
		return CompanyCompanyNotFound
	}

	if inCompany.Status != 2 {
		return CompanyCompanyNotDelete
	}

	inCompany.SetIsDel(ctx, 1)
	inCompany.SetUpdateTime(ctx)

	if _, err := cuc.repo.Update(ctx, inCompany); err != nil {
		return CompanyCompanyUpdateError
	}

	cuc.qarepo.UpdateStatusByCompanyId(ctx, inCompany.Id, 0)

	return nil
}

func (cuc *CompanyUsecase) SyncUpdateStatusCompanys(ctx context.Context) error {
	companys, err := cuc.repo.List(ctx, 0, 0, 0, "", "", 0)

	if err != nil {
		return err
	}

	for _, company := range companys {
		if company.Status != 2 {
			if company.EndTime.AddDate(0, 0, 1).Before(time.Now()) {
				err = cuc.tm.InTx(ctx, func(ctx context.Context) error {
					company.SetStatus(ctx, 2)
					company.SetUpdateTime(ctx)

					if _, err := cuc.repo.Update(ctx, company); err != nil {
						return err
					}

					cuc.qarepo.UpdateStatusByCompanyId(ctx, company.Id, 0)

					if inClue, err := cuc.crepo.GetById(ctx, company.ClueId); err != nil {
						return err
					} else {
						if company.CompanyType == 1 {
							inClue.SetOperationLog(ctx, 0, "", "试用版到期", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
						} else if company.CompanyType == 2 {
							inClue.SetOperationLog(ctx, 0, "", "基础版到期", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
						} else if company.CompanyType == 3 {
							inClue.SetOperationLog(ctx, 0, "", "专业版到期", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
						} else if company.CompanyType == 4 {
							inClue.SetOperationLog(ctx, 0, "", "旗舰版到期", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
						} else if company.CompanyType == 5 {
							inClue.SetOperationLog(ctx, 0, "", "尊享版到期", tool.TimeToString("2006-01-02 15:04:05", time.Now()))
						}

						if _, err := cuc.crepo.Update(ctx, inClue); err != nil {
							return err
						}
					}

					return nil
				})

				if err != nil {
					return CompanyCompanyUpdateError
				}
			}
		}
	}

	return nil
}

func (cuc *CompanyUsecase) getCompanyById(ctx context.Context, id uint64) (*domain.Company, error) {
	company, err := cuc.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	company.Clue, _ = cuc.crepo.GetById(ctx, company.ClueId)
	company.CompanyUser, _ = cuc.curepo.ListByCompanyId(ctx, company.Id)

	if company.Clue != nil {
		industryIds := strings.Split(company.Clue.IndustryId, ",")

		sIndustryIds := make([]string, 0)
		sIndustryNames := make([]string, 0)

		for _, sIndustryId := range industryIds {
			if iIndustryId, err := strconv.ParseUint(sIndustryId, 10, 64); err == nil {
				if industry, err := cuc.irepo.GetById(ctx, iIndustryId); err == nil {
					sIndustryIds = append(sIndustryIds, strconv.FormatUint(industry.Id, 10))
					sIndustryNames = append(sIndustryNames, industry.IndustryName)
				}
			}
		}

		company.Clue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
		company.Clue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))

		if company.Clue.AreaCode > 0 {
			if area, err := cuc.arepo.GetByAreaCode(ctx, company.Clue.AreaCode); err == nil {
				company.Clue.SetAreaName(ctx, area.Data.AreaName)
			}
		}
	}

	return company, nil
}
