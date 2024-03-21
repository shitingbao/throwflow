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
	CompanyClueNotFound         = errors.NotFound("COMPANY_CLUE_NOT_FOUND", "线索不存在")
	CompanyClueCreateError      = errors.InternalServer("COMPANY_CLUE_CREATE_ERROR", "线索创建失败")
	CompanyClueUpdateError      = errors.InternalServer("COMPANY_CLUE_UPDATE_ERROR", "线索更新失败")
	CompanyClueClaimError       = errors.InternalServer("COMPANY_CLUE_CLAIM_ERROR", "线索认领失败")
	CompanyClueCancelClaimError = errors.InternalServer("COMPANY_CLUE_CANCEL_CLAIM_ERROR", "线索取消认领失败")
	CompanyClueDeleteError      = errors.InternalServer("COMPANY_CLUE_DELETE_ERROR", "线索删除失败")
	CompanyClueNotDelete        = errors.InternalServer("COMPANY_CLUE_NOT_DELETE", "当前线索有关联企业不能删除")
	CompanyClueNotClaim         = errors.InternalServer("COMPANY_CLUE_NOT_CLAIM", "当前线索已被其他人认领")
)

type ClueRepo interface {
	GetById(context.Context, uint64) (*domain.Clue, error)
	List(context.Context, int, int, uint64, string, uint8) ([]*domain.Clue, error)
	Count(context.Context, uint64, string, uint8) (int64, error)
	Statistics(context.Context, uint8) (int64, error)
	Save(context.Context, *domain.Clue) (*domain.Clue, error)
	Update(context.Context, *domain.Clue) (*domain.Clue, error)
	Delete(context.Context, *domain.Clue) error
}

type ClueUsecase struct {
	repo     ClueRepo
	irepo    IndustryRepo
	crepo    CompanyRepo
	curepo   CompanyUserRepo
	currepo  CompanyUserRoleRepo
	cucrepo  CompanyUserCompanyRepo
	cpmrepo  CompanyPerformanceMonthlyRepo
	cpdrepo  CompanyPerformanceDailyRepo
	cprrepo  CompanyPerformanceRebalanceRepo
	cprurepo CompanyPerformanceRuleRepo
	arepo    AreaRepo
	tm       Transaction
	conf     *conf.Data
	log      *log.Helper
}

func NewClueUsecase(repo ClueRepo, irepo IndustryRepo, crepo CompanyRepo, curepo CompanyUserRepo, currepo CompanyUserRoleRepo, cucrepo CompanyUserCompanyRepo, cpmrepo CompanyPerformanceMonthlyRepo, cpdrepo CompanyPerformanceDailyRepo, cprrepo CompanyPerformanceRebalanceRepo, cprurepo CompanyPerformanceRuleRepo, arepo AreaRepo, tm Transaction, conf *conf.Data, logger log.Logger) *ClueUsecase {
	return &ClueUsecase{repo: repo, irepo: irepo, crepo: crepo, curepo: curepo, currepo: currepo, cucrepo: cucrepo, cpmrepo: cpmrepo, cpdrepo: cpdrepo, cprrepo: cprrepo, cprurepo: cprurepo, arepo: arepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (cuc *ClueUsecase) GetClueById(ctx context.Context, id uint64) (*domain.Clue, error) {
	clue, err := cuc.getClueById(ctx, id)

	if err != nil {
		return nil, CompanyClueNotFound
	}

	return clue, nil
}

func (cuc *ClueUsecase) ListClues(ctx context.Context, pageNum, pageSize, industryId uint64, keyword string, status uint8) (*domain.ClueList, error) {
	list, err := cuc.repo.List(ctx, int(pageNum), int(pageSize), industryId, keyword, status)

	if err != nil {
		return nil, CompanyDataError
	}

	for _, clue := range list {
		industryIds := strings.Split(clue.IndustryId, ",")

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

		clue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
		clue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))

		if clue.AreaCode > 0 {
			if area, err := cuc.arepo.GetByAreaCode(ctx, clue.AreaCode); err == nil {
				clue.SetAreaName(ctx, area.Data.AreaName)
			}
		}

		if company, err := cuc.crepo.GetByClueId(ctx, clue.Id); err == nil {
			clue.IsAffiliates = 1

			if companyUsers, err := cuc.curepo.ListByCompanyId(ctx, company.Id); err == nil {
				for _, companyUser := range companyUsers {
					if companyUser.Role == 1 {
						clue.AdminName = companyUser.Username
						clue.AdminPhone = companyUser.Phone

						break
					}
				}
			}
		}
	}

	total, err := cuc.repo.Count(ctx, industryId, keyword, status)

	if err != nil {
		return nil, CompanyDataError
	}

	return &domain.ClueList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (cuc *ClueUsecase) ListSelectClues(ctx context.Context) (*domain.SelectClues, error) {
	return domain.NewSelectClues(), nil
}

func (cuc *ClueUsecase) StatisticsClues(ctx context.Context) (*domain.StatisticsClues, error) {
	selects := domain.NewSelectClues()

	statistics := make([]*domain.StatisticsClue, 0)

	for _, status := range selects.Status {
		iStatus, _ := strconv.Atoi(status.Key)

		count, _ := cuc.repo.Statistics(ctx, uint8(iStatus))

		statistics = append(statistics, &domain.StatisticsClue{
			Key:   status.Value,
			Value: strconv.FormatInt(count, 10),
		})

	}

	return &domain.StatisticsClues{
		Statistics: statistics,
	}, nil
}

func (cuc *ClueUsecase) CreateClues(ctx context.Context, companyName, contactInformation, source, seller, facilitator, address, industryId string, userId, areaCode uint64, companyType, qianchuanUse, status uint8) (*domain.Clue, error) {
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

	inClue := domain.NewClue(ctx, companyName, contactInformation, source, seller, facilitator, address, industryId, areaCode, companyType, qianchuanUse, status)
	inClue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
	inClue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))
	inClue.SetCreateTime(ctx)
	inClue.SetUpdateTime(ctx)

	if ok := inClue.VerifyContactInformation(ctx); !ok {
		return nil, CompanyDataError
	}

	inClue.SetOperationLog(ctx, userId, "", "创建", tool.TimeToString("2006-01-02 15:04:05", time.Now()))

	var clue *domain.Clue

	err = cuc.tm.InTx(ctx, func(ctx context.Context) error {
		clue, err = cuc.repo.Save(ctx, inClue)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, CompanyClueCreateError
	}

	clue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
	clue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))

	if area != nil {
		clue.SetAreaName(ctx, area.Data.AreaName)
	}

	return clue, nil
}

func (cuc *ClueUsecase) UpdateClues(ctx context.Context, id, userId uint64, companyName, contactInformation, seller, facilitator, address, industryId string, areaCode uint64, companyType, qianchuanUse, status uint8) (*domain.Clue, error) {
	inClue, err := cuc.getClueById(ctx, id)

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
	inClue.SetAddress(ctx, address)
	inClue.SetQianchuanUse(ctx, qianchuanUse)
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
	inClue.SetOperationLog(ctx, userId, "", "编辑", tool.TimeToString("2006-01-02 15:04:05", time.Now()))

	if ok := inClue.VerifyContactInformation(ctx); !ok {
		return nil, CompanyDataError
	}

	clue, err := cuc.repo.Update(ctx, inClue)

	if err != nil {
		return nil, CompanyClueUpdateError
	}

	clue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
	clue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))

	if area != nil {
		clue.SetAreaName(ctx, area.Data.AreaName)
	}

	return clue, nil
}

func (cuc *ClueUsecase) UpdateCompanyNameClues(ctx context.Context, companyId uint64, companyName string) (*domain.Clue, error) {
	company, err := cuc.crepo.GetById(ctx, companyId)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	inClue, err := cuc.getClueById(ctx, company.ClueId)

	if err != nil {
		return nil, CompanyClueNotFound
	}

	inClue.SetCompanyName(ctx, companyName)
	inClue.SetUpdateTime(ctx)

	clue, err := cuc.repo.Update(ctx, inClue)

	if err != nil {
		return nil, CompanyClueUpdateError
	}

	return clue, nil
}

func (cuc *ClueUsecase) UpdateOperationLogClues(ctx context.Context, id, userId uint64, content string, operationTime time.Time) (*domain.Clue, error) {
	inClue, err := cuc.getClueById(ctx, id)

	if err != nil {
		return nil, CompanyClueNotFound
	}

	inClue.SetUpdateTime(ctx)
	inClue.SetOperationLog(ctx, userId, "", content, tool.TimeToString("2006-01-02 15:04:05", operationTime))

	if ok := inClue.VerifyContactInformation(ctx); !ok {
		return nil, CompanyDataError
	}

	clue, err := cuc.repo.Update(ctx, inClue)

	if err != nil {
		return nil, CompanyClueUpdateError
	}

	return clue, nil
}

func (cuc *ClueUsecase) DeleteClues(ctx context.Context, id uint64) error {
	inClue, err := cuc.getClueById(ctx, id)

	if err != nil {
		return CompanyClueNotFound
	}

	if inCompany, err := cuc.crepo.GetByClueId(ctx, inClue.Id); err == nil {
		if inCompany.IsDel == 1 {
			cerr := cuc.tm.InTx(ctx, func(ctx context.Context) error {
				inClue.SetIsDel(ctx, 1)

				if _, err := cuc.repo.Update(ctx, inClue); err != nil {
					return err
				}

				if err := cuc.cpdrepo.DeleteByCompanyId(ctx, inCompany.Id); err != nil {
					return err
				}

				if err := cuc.cpmrepo.DeleteByCompanyId(ctx, inCompany.Id); err != nil {
					return err
				}

				if err := cuc.cprrepo.DeleteByCompanyId(ctx, inCompany.Id); err != nil {
					return err
				}

				if err := cuc.cprurepo.DeleteByCompanyId(ctx, inCompany.Id); err != nil {
					return err
				}

				if err := cuc.curepo.DeleteByCompanyId(ctx, inCompany.Id); err != nil {
					return err
				}

				day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

				if companyUserRoles, err := cuc.currepo.ListByCompanyIdAndDay(ctx, inCompany.Id, uint32(day)); err == nil {
					for _, companyUserRole := range companyUserRoles {
						if companyUserRole.Day == uint32(day) {
							if err := cuc.currepo.DeleteByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx, companyUserRole.UserId, companyUserRole.CompanyId, companyUserRole.AdvertiserId, companyUserRole.Day); err != nil {
								return err
							}
						} else {
							inCompanyUserRole := domain.NewCompanyUserRole(ctx, companyUserRole.UserId, companyUserRole.AdvertiserId, inCompany.Id, uint32(day), 2)
							inCompanyUserRole.SetCreateTime(ctx)
							inCompanyUserRole.SetUpdateTime(ctx)

							if _, err := cuc.currepo.Save(ctx, inCompanyUserRole); err != nil {
								return err
							}
						}
					}
				} else {
					return err
				}

				if err := cuc.cucrepo.DeleteByCompanyId(ctx, inCompany.Id); err != nil {
					return err
				}

				return nil
			})

			if cerr != nil {
				return CompanyClueDeleteError
			}
		} else {
			return CompanyClueNotDelete
		}
	} else {
		if derr := cuc.repo.Delete(ctx, inClue); derr != nil {
			return CompanyClueDeleteError
		}
	}

	return nil
}

func (cuc *ClueUsecase) getClueById(ctx context.Context, id uint64) (*domain.Clue, error) {
	clue, err := cuc.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	industryIds := strings.Split(clue.IndustryId, ",")

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

	clue.SetIndustryId(ctx, strings.Join(sIndustryIds, ","))
	clue.SetIndustryName(ctx, strings.Join(sIndustryNames, ","))

	if clue.AreaCode > 0 {
		if area, err := cuc.arepo.GetByAreaCode(ctx, clue.AreaCode); err == nil {
			clue.SetAreaName(ctx, area.Data.AreaName)
		}
	}

	return clue, nil
}
