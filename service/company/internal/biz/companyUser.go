package biz

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var (
	CompanyCompanyUserNotFound           = errors.NotFound("COMPANY_COMPANY_USER_NOT_FOUND", "企业用户不存在")
	CompanyCompanyUserFound              = errors.NotFound("COMPANY_COMPANY_USER_FOUND", "企业用户已经存在")
	CompanyCompanyMainUserFound          = errors.NotFound("COMPANY_COMPANY_MAIN_USER_FOUND", "企业已经存在主管理员")
	CompanyCompanyUserCreateError        = errors.InternalServer("COMPANY_COMPANY_USER_CREATE_ERROR", "企业用户创建失败")
	CompanyCompanyUserUpdateError        = errors.InternalServer("COMPANY_COMPANY_USER_UPDATE_ERROR", "企业用户更新失败")
	CompanyCompanyUserDeleteError        = errors.InternalServer("COMPANY_COMPANY_USER_DELETE_ERROR", "企业用户删除失败")
	CompanyCompanyMainUserNotDelete      = errors.InternalServer("COMPANY_COMPANY_MAIN_USER_NOT_DELETE", "企业主管理员不能删除")
	CompanyCompanyUserLimitError         = errors.InternalServer("COMPANY_COMPANY_USER_LIMIT_ERROR", "企业账户总数超出限制")
	CompanyLoginFailed                   = errors.InternalServer("COMPANY_LOGIN_FAILED", "登录失败")
	CompanyLoginPhoneNotExist            = errors.InternalServer("COMPANY_LOGIN_PHONE_NOE_EXIST", "登录账号不存在")
	CompanyCompanyDisable                = errors.InternalServer("COMPANY_COMPANY_DISABLE", "公司被禁用")
	CompanyCompanyUserDisable            = errors.InternalServer("COMPANY_COMPANY_USER_DISABLE", "账户被禁用")
	CompanyCompanyExpire                 = errors.InternalServer("COMPANY_COMPANY_EXPIRE", "账户权限已过期")
	CompanyLoginError                    = errors.InternalServer("COMPANY_LOGIN_ERROR", "登录异常错误")
	CompanyChangeCompanyUserCompanyError = errors.InternalServer("COMPANY_CHANGE_COMPANY_USER_COMPANY_ERROR", "切换公司失败")
	CompanyCompanyUserRoleNotFound       = errors.InternalServer("COMPANY_COMPANY_USER_ROLE_NOT_FOUND", "没有被分配广告组")
)

type CompanyUserRepo interface {
	GetById(context.Context, uint64, uint64) (*domain.CompanyUser, error)
	GetByPhone(context.Context, uint64, string) (*domain.CompanyUser, error)
	GetByPhoneAndNotInUserId(context.Context, uint64, string, []uint64) (*domain.CompanyUser, error)
	GetByOrganizationIdAndPhone(context.Context, uint64, uint64, string) (*domain.CompanyUser, error)
	GetByRole(context.Context, uint64, uint8) (*domain.CompanyUser, error)
	ListByCompanyId(context.Context, uint64) ([]*domain.CompanyUser, error)
	ListByCompanyIdAndOrganizationId(context.Context, uint64, uint64) ([]*domain.CompanyUser, error)
	ListByPhone(context.Context, string) ([]*domain.CompanyUser, error)
	List(context.Context, uint64, int, int, string) ([]*domain.CompanyUser, error)
	Count(context.Context, uint64, string) (int64, error)
	Save(context.Context, *domain.CompanyUser) (*domain.CompanyUser, error)
	Update(context.Context, *domain.CompanyUser) (*domain.CompanyUser, error)
	DeleteByCompanyId(context.Context, uint64) error
	DeleteByCompanyIdAndUserId(context.Context, uint64, []uint64) error
	DeleteByCompanyIdAndOrganizationId(context.Context, uint64, uint64) error
	Delete(context.Context, *domain.CompanyUser) error

	SaveCacheHash(context.Context, string, map[string]string, time.Duration) error
	UpdateCacheHash(context.Context, string, map[string]string) error
	GetCacheHash(context.Context, string, string) (string, error)

	SaveCacheString(context.Context, string, string, time.Duration) error
	GetCacheString(context.Context, string) (string, error)

	ExpireCache(context.Context, string, time.Duration) error
	DeleteCache(context.Context, string) error
}

type CompanyUserUsecase struct {
	repo     CompanyUserRepo
	crepo    CompanyRepo
	cucrepo  CompanyUserCompanyRepo
	qarepo   QianchuanAdvertiserRepo
	qahrepo  QianchuanAdvertiserHistoryRepo
	qrarepo  QianchuanReportAdvertiserRepo
	qrawrepo QianchuanReportAwemeRepo
	qrprepo  QianchuanReportProductRepo
	qadrepo  QianchuanAdRepo
	currepo  CompanyUserRoleRepo
	csrepo   CompanySetRepo
	clrpeo   ClueRepo
	mrepo    MenuRepo
	cuwrepo  CompanyUserWhiteRepo
	cqsrepo  CompanyUserQianchuanSearchRepo
	srepo    SmsRepo
	tm       Transaction
	conf     *conf.Data
	log      *log.Helper
}

func NewCompanyUserUsecase(repo CompanyUserRepo, crepo CompanyRepo, cucrepo CompanyUserCompanyRepo, qarepo QianchuanAdvertiserRepo, qahrepo QianchuanAdvertiserHistoryRepo, qrarepo QianchuanReportAdvertiserRepo, qrawrepo QianchuanReportAwemeRepo, qrprepo QianchuanReportProductRepo, qadrepo QianchuanAdRepo, currepo CompanyUserRoleRepo, csrepo CompanySetRepo, clrpeo ClueRepo, mrepo MenuRepo, cuwrepo CompanyUserWhiteRepo, cqsrepo CompanyUserQianchuanSearchRepo, srepo SmsRepo, tm Transaction, conf *conf.Data, logger log.Logger) *CompanyUserUsecase {
	return &CompanyUserUsecase{repo: repo, crepo: crepo, cucrepo: cucrepo, qarepo: qarepo, qahrepo: qahrepo, qrarepo: qrarepo, qrawrepo: qrawrepo, qrprepo: qrprepo, qadrepo: qadrepo, currepo: currepo, csrepo: csrepo, clrpeo: clrpeo, mrepo: mrepo, cuwrepo: cuwrepo, cqsrepo: cqsrepo, srepo: srepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (cuuc *CompanyUserUsecase) ListCompanyUsers(ctx context.Context, companyId, pageNum, pageSize uint64, keyword string) (*domain.CompanyUserList, error) {
	companyUsers, err := cuuc.repo.List(ctx, companyId, int(pageNum), int(pageSize), keyword)

	if err != nil {
		return nil, CompanyDataError
	}

	list := make([]*domain.CompanyUser, 0)

	for _, companyUser := range companyUsers {
		companyUser, _ = cuuc.getCompanyUsersById(ctx, companyUser.Id, companyUser.CompanyId)

		list = append(list, companyUser)
	}

	total, err := cuuc.repo.Count(ctx, companyId, keyword)

	if err != nil {
		return nil, CompanyDataError
	}

	return &domain.CompanyUserList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUsersByPhone(ctx context.Context, phone string) ([]*domain.CompanyUser, error) {
	companyUsers, err := cuuc.repo.ListByPhone(ctx, phone)

	if err != nil {
		return nil, CompanyDataError
	}

	list := make([]*domain.CompanyUser, 0)

	for _, companyUser := range companyUsers {
		companyUser, _ = cuuc.getCompanyUsersById(ctx, companyUser.Id, companyUser.CompanyId)

		list = append(list, companyUser)
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListQianchuanAdvertisersCompanyUsers(ctx context.Context, id, companyId uint64) ([]*domain.QianchuanAdvertiserList, error) {
	companyUser, err := cuuc.repo.GetById(ctx, id, companyId)

	if err != nil {
		return nil, CompanyCompanyUserNotFound
	}

	qianchuanAdvertisers, err := cuuc.qarepo.List(ctx, companyId, 0, 0, "", "", "1")

	if err != nil {
		return nil, CompanyDataError
	}

	list := make([]*domain.QianchuanAdvertiserList, 0)

	day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

	if companyUser.Role == 2 {
		isExistRoles := make([]uint64, 0)

		if companyUserRoles, err := cuuc.currepo.ListByUserIdAndCompanyIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, uint32(day)); err == nil {
			for _, companyUserRole := range companyUserRoles {
				isExistRoles = append(isExistRoles, companyUserRole.AdvertiserId)
			}
		}

		for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
			lqianchuanAdvertiser := &domain.QianchuanAdvertiserList{
				AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
				Status:         qianchuanAdvertiser.Status,
				IsSelect:       0,
				AdvertiserName: qianchuanAdvertiser.AdvertiserName,
				CompanyName:    qianchuanAdvertiser.CompanyName,
			}

			if len(isExistRoles) == 0 {
				lqianchuanAdvertiser.IsSelect = 1
			} else {
				isExit := false

				for _, isExistRole := range isExistRoles {
					if isExistRole == qianchuanAdvertiser.AdvertiserId {
						isExit = true
						break
					}
				}

				if isExit {
					lqianchuanAdvertiser.IsSelect = 1
				}
			}

			list = append(list, lqianchuanAdvertiser)
		}
	} else if companyUser.Role == 3 {
		if exitCompanyUserRoles, err := cuuc.currepo.ListByUserIdAndCompanyIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, uint32(day)); err == nil {
			for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
				lqianchuanAdvertiser := &domain.QianchuanAdvertiserList{
					AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
					Status:         qianchuanAdvertiser.Status,
					IsSelect:       0,
					AdvertiserName: qianchuanAdvertiser.AdvertiserName,
					CompanyName:    qianchuanAdvertiser.CompanyName,
				}

				for _, exitCompanyUserRole := range exitCompanyUserRoles {
					if exitCompanyUserRole.AdvertiserId == qianchuanAdvertiser.AdvertiserId {
						if exitCompanyUserRole.UserId == companyUser.Id {
							lqianchuanAdvertiser.IsSelect = 1
						}

						break
					}
				}

				list = append(list, lqianchuanAdvertiser)
			}
		}
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListSelectCompanyUsers(ctx context.Context) (*domain.SelectCompanyUsers, error) {
	return domain.NewSelectCompanyUsers(), nil
}

func (cuuc *CompanyUserUsecase) StatisticsCompanyUsers(ctx context.Context, companyId uint64) (*domain.StatisticsCompanyUsers, error) {
	company, err := cuuc.crepo.GetById(ctx, companyId)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	totalCompanyUsers := 0
	totalQianchuanAdvertisers := 0

	if companyUsers, err := cuuc.repo.ListByCompanyId(ctx, company.Id); err == nil {
		for _, companyUser := range companyUsers {
			if companyUser.Status == 1 {
				totalCompanyUsers += 1
			}
		}
	}

	if qianchuanAdvertisers, err := cuuc.qarepo.List(ctx, company.Id, 0, 0, "", "", ""); err == nil {
		for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
			if qianchuanAdvertiser.Status == 1 {
				totalQianchuanAdvertisers += 1
			}
		}
	}

	list := make([]*domain.StatisticsCompanyUser, 0)

	list = append(list, &domain.StatisticsCompanyUser{
		Key:   "companyUser",
		Value: strconv.FormatInt(int64(totalCompanyUsers), 10),
	})

	list = append(list, &domain.StatisticsCompanyUser{
		Key:   "qianchuanAdvertiser",
		Value: strconv.FormatInt(int64(totalQianchuanAdvertisers), 10),
	})

	list = append(list, &domain.StatisticsCompanyUser{
		Key:   "totalAccounts",
		Value: strconv.FormatInt(int64(company.Accounts), 10),
	})

	list = append(list, &domain.StatisticsCompanyUser{
		Key:   "totalQianchuanAdvertiser",
		Value: strconv.FormatInt(int64(company.QianchuanAdvertisers), 10),
	})

	return &domain.StatisticsCompanyUsers{
		Statistics: list,
	}, nil
}

func (cuuc *CompanyUserUsecase) GetCompanyUsers(ctx context.Context, id, companyId uint64) (*domain.CompanyUser, error) {
	company, err := cuuc.crepo.GetById(ctx, companyId)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	companyUser, err := cuuc.repo.GetById(ctx, id, company.Id)

	if err != nil {
		return nil, CompanyCompanyUserNotFound
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) CreateCompanyUsers(ctx context.Context, companyId uint64, username, job, phone string, role uint8) (*domain.CompanyUser, error) {
	company, err := cuuc.crepo.GetById(ctx, companyId)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	if err := cuuc.verifyCompanyUserLimit(ctx, company, 0); err != nil {
		return nil, err
	}

	if _, err := cuuc.repo.GetByPhone(ctx, company.Id, phone); err == nil {
		return nil, CompanyCompanyUserFound
	}

	if role == 1 {
		if _, err := cuuc.repo.GetByRole(ctx, company.Id, role); err == nil {
			return nil, CompanyCompanyMainUserFound
		}
	}

	inCompanyUser := domain.NewCompanyUser(ctx, companyId, username, job, phone, role, 1)
	inCompanyUser.SetCreateTime(ctx)
	inCompanyUser.SetUpdateTime(ctx)

	companyUser, err := cuuc.repo.Save(ctx, inCompanyUser)

	if err != nil {
		return nil, CompanyCompanyUserCreateError
	}

	content := domain.AccountOpend{
		Phone:    companyUser.Phone,
		RoleName: companyUser.RoleName,
	}

	if contentByte, err := json.Marshal(content); err == nil {
		cuuc.srepo.Send(ctx, companyUser.Phone, string(contentByte), "accountOpend", tool.GetClientIp(ctx))
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateCompanyUsers(ctx context.Context, id, companyId uint64, username, job, phone string, role uint8) (*domain.CompanyUser, error) {
	inCompanyUser, err := cuuc.repo.GetById(ctx, id, companyId)

	if err != nil {
		return nil, CompanyCompanyUserNotFound
	}

	if inCompanyUser.Phone != phone {
		if _, err := cuuc.repo.GetByPhone(ctx, inCompanyUser.CompanyId, phone); err == nil {
			return nil, CompanyCompanyUserFound
		}
	}

	if role == 1 {
		if inCompanyUser.Role != role {
			if _, err := cuuc.repo.GetByRole(ctx, inCompanyUser.CompanyId, role); err == nil {
				return nil, CompanyCompanyMainUserFound
			}
		}
	}

	oldRole := inCompanyUser.Role

	inCompanyUser.SetUsername(ctx, username)
	inCompanyUser.SetJob(ctx, job)
	inCompanyUser.SetPhone(ctx, phone)
	inCompanyUser.SetRole(ctx, role)
	inCompanyUser.SetUpdateTime(ctx)

	var companyUser *domain.CompanyUser

	err = cuuc.tm.InTx(ctx, func(ctx context.Context) error {
		if companyUser, err = cuuc.repo.Update(ctx, inCompanyUser); err != nil {
			return err
		}

		if role == 1 && oldRole != 1 {
			day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

			if companyUserRoles, err := cuuc.currepo.ListByUserIdAndCompanyIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, uint32(day)); err == nil {
				for _, companyUserRole := range companyUserRoles {
					if companyUserRole.Day == uint32(day) {
						if err := cuuc.currepo.DeleteByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, companyUserRole.AdvertiserId, companyUserRole.Day); err != nil {
							return err
						}
					} else {
						inCompanyUserRole := domain.NewCompanyUserRole(ctx, companyUserRole.UserId, companyUserRole.AdvertiserId, companyUser.CompanyId, uint32(day), 2)
						inCompanyUserRole.SetCreateTime(ctx)
						inCompanyUserRole.SetUpdateTime(ctx)

						if _, err := cuuc.currepo.Save(ctx, inCompanyUserRole); err != nil {
							return err
						}
					}
				}
			} else {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, CompanyCompanyUserUpdateError
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateStatusCompanyUsers(ctx context.Context, id, companyId uint64, status uint8) (*domain.CompanyUser, error) {
	inCompanyUser, err := cuuc.getCompanyUsersById(ctx, id, companyId)

	if err != nil {
		return nil, CompanyCompanyUserNotFound
	}

	if status == 1 {
		company, err := cuuc.crepo.GetById(ctx, companyId)

		if err != nil {
			return nil, CompanyCompanyNotFound
		}

		if err := cuuc.verifyCompanyUserLimit(ctx, company, inCompanyUser.Id); err != nil {
			return nil, err
		}
	}

	inCompanyUser.SetStatus(ctx, status)
	inCompanyUser.SetUpdateTime(ctx)

	companyUser, err := cuuc.repo.Update(ctx, inCompanyUser)

	if err != nil {
		return nil, CompanyCompanyUserUpdateError
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateRoleCompanyUsers(ctx context.Context, id, companyId uint64, roleIds string) (*domain.CompanyUser, error) {
	companyUser, err := cuuc.getCompanyUsersById(ctx, id, companyId)

	if err != nil {
		return nil, CompanyCompanyUserNotFound
	}

	day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

	if len(roleIds) == 0 {
		err = cuuc.tm.InTx(ctx, func(ctx context.Context) error {
			if companyUserRoles, err := cuuc.currepo.ListByUserIdAndCompanyIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, uint32(day)); err == nil {
				for _, inCompanyUserRole := range companyUserRoles {
					if inCompanyUserRole.Day == uint32(day) {
						inCompanyUserRole.SetRoleType(ctx, 2)
						inCompanyUserRole.SetUpdateTime(ctx)

						if _, err := cuuc.currepo.Update(ctx, inCompanyUserRole); err != nil {
							return err
						}
					} else {
						inCompanyUserRole = domain.NewCompanyUserRole(ctx, companyUser.Id, inCompanyUserRole.AdvertiserId, companyUser.CompanyId, uint32(day), 2)
						inCompanyUserRole.SetCreateTime(ctx)
						inCompanyUserRole.SetUpdateTime(ctx)

						if _, err := cuuc.currepo.Save(ctx, inCompanyUserRole); err != nil {
							return err
						}
					}
				}
			} else {
				return err
			}

			return nil
		})

		if err != nil {
			return nil, CompanyCompanyUserUpdateError
		}
	} else {
		qianchuanAdvertisers, err := cuuc.qarepo.List(ctx, companyId, 0, 0, "", "", "1")

		if err != nil {
			return nil, CompanyDataError
		}

		sRoleIds := strings.Split(roleIds, ",")

		if companyUser.Role == 2 || companyUser.Role == 3 {
			err = cuuc.tm.InTx(ctx, func(ctx context.Context) error {
				companyUserRoles, err := cuuc.currepo.ListByUserIdAndCompanyIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, uint32(day))

				if err != nil {
					return err
				}

				for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
					for _, sRoleId := range sRoleIds {
						if sIRoleId, err := strconv.Atoi(sRoleId); err == nil {
							if uint64(sIRoleId) == qianchuanAdvertiser.AdvertiserId {
								isNotExist := true

								for _, inCompanyUserRole := range companyUserRoles {
									if inCompanyUserRole.AdvertiserId == uint64(sIRoleId) {
										isNotExist = false

										break
									}
								}

								if isNotExist {
									if companyUserRole, err := cuuc.currepo.GetByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, uint64(sIRoleId), uint32(day)); err == nil {
										if err := cuuc.currepo.DeleteByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx, companyUserRole.UserId, companyUserRole.CompanyId, companyUserRole.AdvertiserId, companyUserRole.Day); err != nil {
											return err
										}
									} else {
										inCompanyUserRole := domain.NewCompanyUserRole(ctx, companyUser.Id, uint64(sIRoleId), companyUser.CompanyId, uint32(day), 1)
										inCompanyUserRole.SetCreateTime(ctx)
										inCompanyUserRole.SetUpdateTime(ctx)

										if _, err := cuuc.currepo.Save(ctx, inCompanyUserRole); err != nil {
											return err
										}
									}
								}

								break
							}
						} else {
							return err
						}
					}
				}

				for _, companyUserRole := range companyUserRoles {
					isNotExist := true

					for _, sRoleId := range sRoleIds {
						if sIRoleId, err := strconv.Atoi(sRoleId); err == nil {
							if companyUserRole.AdvertiserId == uint64(sIRoleId) {
								isNotExist = false

								break
							}
						} else {
							return err
						}
					}

					if isNotExist {
						if companyUserRole.Day == uint32(day) {
							if err := cuuc.currepo.DeleteByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx, companyUserRole.UserId, companyUserRole.CompanyId, companyUserRole.AdvertiserId, companyUserRole.Day); err != nil {
								return err
							}
						} else {
							inCompanyUserRole := domain.NewCompanyUserRole(ctx, companyUser.Id, companyUserRole.AdvertiserId, companyUser.CompanyId, uint32(day), 2)
							inCompanyUserRole.SetCreateTime(ctx)
							inCompanyUserRole.SetUpdateTime(ctx)

							if _, err := cuuc.currepo.Save(ctx, inCompanyUserRole); err != nil {
								return err
							}
						}
					}
				}

				return nil
			})

			if err != nil {
				return nil, CompanyCompanyUserUpdateError
			}
		}
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateWhiteCompanyUsers(ctx context.Context, id, companyId uint64, isWhite uint8) (*domain.CompanyUser, error) {
	inCompanyUser, err := cuuc.getCompanyUsersById(ctx, id, companyId)

	if err != nil {
		return nil, CompanyCompanyUserNotFound
	}

	if isWhite == 1 {
		inCompanyUserWhite := domain.NewCompanyUserWhite(ctx, inCompanyUser.Phone)
		inCompanyUserWhite.SetCreateTime(ctx)
		inCompanyUserWhite.SetUpdateTime(ctx)

		if _, err := cuuc.cuwrepo.Save(ctx, inCompanyUserWhite); err != nil {
			return nil, CompanyCompanyUserUpdateError
		}
	} else {
		if err := cuuc.cuwrepo.Delete(ctx, inCompanyUser.Phone); err != nil {
			return nil, CompanyCompanyUserUpdateError
		}
	}

	companyUser, err := cuuc.getCompanyUsersById(ctx, id, companyId)

	if err != nil {
		return nil, CompanyCompanyUserUpdateError
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) DeleteCompanyUsers(ctx context.Context, id, companyId uint64) error {
	inCompanyUser, err := cuuc.getCompanyUsersById(ctx, id, companyId)

	if err != nil {
		return CompanyCompanyUserNotFound
	}

	if inCompanyUser.Role == 1 {
		return CompanyCompanyMainUserNotDelete
	}

	err = cuuc.tm.InTx(ctx, func(ctx context.Context) error {
		if err := cuuc.repo.Delete(ctx, inCompanyUser); err != nil {
			return err
		}

		if _, err := cuuc.repo.ListByPhone(ctx, inCompanyUser.Phone); err != nil {
			inCompanyUserCompany, err := cuuc.cucrepo.GetByPhone(ctx, inCompanyUser.Phone)

			if err == nil {
				cuuc.cucrepo.Delete(ctx, inCompanyUserCompany)
			}

			cuuc.cuwrepo.Delete(ctx, inCompanyUser.Phone)
		}

		day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

		if companyUserRoles, err := cuuc.currepo.ListByUserIdAndCompanyIdAndDay(ctx, inCompanyUser.Id, inCompanyUser.CompanyId, uint32(day)); err == nil {
			for _, companyUserRole := range companyUserRoles {
				if companyUserRole.Day == uint32(day) {
					if err := cuuc.currepo.DeleteByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx, companyUserRole.UserId, companyUserRole.CompanyId, companyUserRole.AdvertiserId, companyUserRole.Day); err != nil {
						return err
					}
				} else {
					inCompanyUserRole := domain.NewCompanyUserRole(ctx, companyUserRole.UserId, companyUserRole.AdvertiserId, inCompanyUser.CompanyId, uint32(day), 2)
					inCompanyUserRole.SetCreateTime(ctx)
					inCompanyUserRole.SetUpdateTime(ctx)

					if _, err := cuuc.currepo.Save(ctx, inCompanyUserRole); err != nil {
						return err
					}
				}
			}
		} else {
			return err
		}

		return nil
	})

	if err != nil {
		return CompanyCompanyUserDeleteError
	}

	return nil
}

func (cuuc *CompanyUserUsecase) DeleteRoleCompanyUsers(ctx context.Context, companyId, advertiserId uint64) error {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return CompanyCompanyNotFound
	}

	day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

	if companyUserRoles, err := cuuc.currepo.ListByCompanyIdAndAdvertiserIdAndDay(ctx, companyId, advertiserId, uint32(day)); err == nil {
		for _, companyUserRole := range companyUserRoles {
			if companyUserRole.Day == uint32(day) {
				if err := cuuc.currepo.DeleteByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx, companyUserRole.UserId, companyUserRole.CompanyId, companyUserRole.AdvertiserId, companyUserRole.Day); err != nil {
					return CompanyCompanyUserDeleteError
				}
			} else {
				inCompanyUserRole := domain.NewCompanyUserRole(ctx, companyUserRole.UserId, companyUserRole.AdvertiserId, companyId, uint32(day), 2)
				inCompanyUserRole.SetCreateTime(ctx)
				inCompanyUserRole.SetUpdateTime(ctx)

				if _, err := cuuc.currepo.Save(ctx, inCompanyUserRole); err != nil {
					return CompanyCompanyUserDeleteError
				}
			}
		}
	} else {
		return CompanyCompanyUserDeleteError
	}

	return nil
}

func (cuuc *CompanyUserUsecase) LoginCompanyUser(ctx context.Context, phone string) (*domain.LoginCompanyUser, error) {
	loginCompanyUser := &domain.LoginCompanyUser{}

	if token, err := cuuc.repo.GetCacheString(ctx, "company:user:"+phone); err != nil {
		token = tool.GetToken()

		if companyUserCompany, err := cuuc.cucrepo.GetByPhone(ctx, phone); err != nil {
			if _, err := cuuc.cuwrepo.GetByPhone(ctx, phone); err != nil {
				if companyUsers, err := cuuc.repo.ListByPhone(ctx, phone); err == nil {
					for _, companyUser := range companyUsers {
						if company, err := cuuc.crepo.GetById(ctx, companyUser.CompanyId); err == nil {
							if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
								company.Clue = clue
							}

							if company.Status == 1 && companyUser.Status == 1 {
								loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
									Id:          company.Id,
									CompanyName: company.Clue.CompanyName,
								})
							}
						}
					}
				}
			} else {
				if companys, err := cuuc.crepo.List(ctx, 0, 0, 0, "", "", 0); err == nil {
					for _, company := range companys {
						if company.Status == 1 {
							if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
								company.Clue = clue
							}

							loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
								Id:          company.Id,
								CompanyName: company.Clue.CompanyName,
							})
						}
					}
				}
			}

			loginCompanyUser.Token = token

			cacheData := make(map[string]string)
			cacheData["companyId"] = "0"
			cacheData["phone"] = phone

			if err := cuuc.repo.SaveCacheHash(ctx, "company:user:"+token, cacheData, cuuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
				return nil, CompanyLoginFailed
			}

			if err := cuuc.repo.SaveCacheString(ctx, "company:user:"+phone, token, cuuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
				return nil, CompanyLoginFailed
			}

			return loginCompanyUser, nil
		} else {
			company, err := cuuc.crepo.GetById(ctx, companyUserCompany.CompanyId)

			if err != nil {
				loginCompanyUser.Reason = "公司异常"
			} else {
				if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
					company.Clue = clue

					if company.Status == 0 {
						loginCompanyUser.Reason = "公司被禁用"
					} else if company.Status == 2 {
						loginCompanyUser.Reason = "账户权限已过期"
					} else {
						loginCompanyUser.CurrentCompanyId = company.Id
					}
				} else {
					loginCompanyUser.Reason = "公司异常"
				}
			}

			var companyUser *domain.CompanyUser

			if _, err := cuuc.cuwrepo.GetByPhone(ctx, phone); err != nil {
				if company != nil {
					companyUser, err = cuuc.repo.GetByPhone(ctx, company.Id, phone)

					if err != nil {
						loginCompanyUser.Reason = "公司异常"
					} else {
						if companyUser.Status == 0 {
							loginCompanyUser.Reason = "账户被禁用"
						}
					}
				}

				if companyUsers, err := cuuc.repo.ListByPhone(ctx, phone); err == nil {
					for _, companyUser := range companyUsers {
						if company, err := cuuc.crepo.GetById(ctx, companyUser.CompanyId); err == nil {
							if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
								company.Clue = clue
							}

							if company.Status == 1 && companyUser.Status == 1 {
								loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
									Id:          company.Id,
									CompanyName: company.Clue.CompanyName,
								})
							}
						}
					}
				}

				loginCompanyUser.IsWhite = 0
			} else {
				companyUsers, err := cuuc.repo.ListByPhone(ctx, phone)

				if err != nil || len(companyUsers) == 0 {
					loginCompanyUser.Reason = "公司异常"
				}

				if loginCompanyUser.CurrentCompanyId > 0 {
					isExitCompany := true

					for _, cUser := range companyUsers {
						if cUser.CompanyId == loginCompanyUser.CurrentCompanyId {
							companyUser = cUser

							isExitCompany = false
						}
					}

					if isExitCompany {
						companyUser = companyUsers[0]
					}
				} else {
					companyUser = companyUsers[0]
				}

				if companys, err := cuuc.crepo.List(ctx, 0, 0, 0, "", "", 0); err == nil {
					for _, company := range companys {
						if company.Status == 1 {
							if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
								company.Clue = clue
							}

							loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
								Id:          company.Id,
								CompanyName: company.Clue.CompanyName,
							})
						}
					}
				}

				loginCompanyUser.IsWhite = 1
			}

			if companyUser != nil {
				loginCompanyUser.Id = companyUser.Id
				loginCompanyUser.CompanyId = companyUser.CompanyId
				loginCompanyUser.Username = companyUser.Username
				loginCompanyUser.Job = companyUser.Job
				loginCompanyUser.Phone = companyUser.Phone
				loginCompanyUser.Role = companyUser.Role
				loginCompanyUser.RoleName = companyUser.RoleName
				loginCompanyUser.Role = companyUser.Role
			}

			if company != nil {
				loginCompanyUser.CompanyType = company.CompanyType
				loginCompanyUser.CompanyTypeName = company.CompanyTypeName
				loginCompanyUser.CompanyName = company.Clue.CompanyName
				loginCompanyUser.CompanyStartTime = company.StartTime
				loginCompanyUser.CompanyEndTime = company.EndTime
				loginCompanyUser.Accounts = company.Accounts
				loginCompanyUser.QianchuanAdvertisers = company.QianchuanAdvertisers
				loginCompanyUser.IsTermwork = company.IsTermwork
			}

			loginCompanyUser.Token = token

			cacheData := make(map[string]string)

			cacheData["phone"] = phone

			if loginCompanyUser.CurrentCompanyId > 0 {
				cacheData["companyId"] = strconv.FormatUint(company.Id, 10)
			} else {
				cacheData["companyId"] = "0"
			}

			if err := cuuc.repo.SaveCacheHash(ctx, "company:user:"+token, cacheData, cuuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
				return nil, CompanyLoginFailed
			}

			if err := cuuc.repo.SaveCacheString(ctx, "company:user:"+phone, token, cuuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
				return nil, CompanyLoginFailed
			}

			return loginCompanyUser, nil
		}
	} else {
		cCompanyId, err := cuuc.repo.GetCacheHash(ctx, "company:user:"+token, "companyId")

		if err != nil {
			return nil, CompanyLoginFailed
		}

		companyId, err := strconv.Atoi(cCompanyId)

		if err != nil {
			return nil, CompanyLoginFailed
		}

		if companyId == 0 {
			if _, err := cuuc.cuwrepo.GetByPhone(ctx, phone); err != nil {
				if companyUsers, err := cuuc.repo.ListByPhone(ctx, phone); err == nil {
					for _, companyUser := range companyUsers {
						if company, err := cuuc.crepo.GetById(ctx, companyUser.CompanyId); err == nil {
							if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
								company.Clue = clue
							}

							if company.Status == 1 && companyUser.Status == 1 {
								loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
									Id:          company.Id,
									CompanyName: company.Clue.CompanyName,
								})
							}
						}
					}
				}
			} else {
				if companys, err := cuuc.crepo.List(ctx, 0, 0, 0, "", "", 0); err == nil {
					for _, company := range companys {
						if company.Status == 1 {
							if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
								company.Clue = clue
							}

							loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
								Id:          company.Id,
								CompanyName: company.Clue.CompanyName,
							})
						}
					}
				}
			}

			ltoken := tool.GetToken()

			loginCompanyUser.Token = ltoken

			if err := cuuc.repo.DeleteCache(ctx, "company:user:"+token); err != nil {
				return nil, CompanyLoginFailed
			}

			if err := cuuc.repo.DeleteCache(ctx, "company:user:"+phone); err != nil {
				return nil, CompanyLoginFailed
			}

			cacheData := make(map[string]string)
			cacheData["companyId"] = "0"
			cacheData["phone"] = phone

			if err := cuuc.repo.SaveCacheHash(ctx, "company:user:"+ltoken, cacheData, cuuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
				return nil, CompanyLoginFailed
			}

			if err := cuuc.repo.SaveCacheString(ctx, "company:user:"+phone, ltoken, cuuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
				return nil, CompanyLoginFailed
			}

			return loginCompanyUser, nil
		} else {
			company, err := cuuc.crepo.GetById(ctx, uint64(companyId))

			if err != nil {
				loginCompanyUser.Reason = "公司异常"
			} else {
				if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
					company.Clue = clue

					if company.Status == 0 {
						loginCompanyUser.Reason = "公司被禁用"
					} else if company.Status == 2 {
						loginCompanyUser.Reason = "账户权限已过期"
					} else {
						loginCompanyUser.CurrentCompanyId = company.Id
					}
				} else {
					loginCompanyUser.Reason = "公司异常"
				}
			}

			var companyUser *domain.CompanyUser

			if _, err := cuuc.cuwrepo.GetByPhone(ctx, phone); err != nil {
				if company != nil {
					companyUser, err = cuuc.repo.GetByPhone(ctx, company.Id, phone)

					if err != nil {
						loginCompanyUser.Reason = "公司异常"
					} else {
						if companyUser.Status == 0 {
							loginCompanyUser.Reason = "账户被禁用"
						}
					}
				}

				if companyUsers, err := cuuc.repo.ListByPhone(ctx, phone); err == nil {
					for _, companyUser := range companyUsers {
						if company, err := cuuc.crepo.GetById(ctx, companyUser.CompanyId); err == nil {
							if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
								company.Clue = clue
							}

							if company.Status == 1 && companyUser.Status == 1 {
								loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
									Id:          company.Id,
									CompanyName: company.Clue.CompanyName,
								})
							}
						}
					}
				}

				loginCompanyUser.IsWhite = 0
			} else {
				companyUsers, err := cuuc.repo.ListByPhone(ctx, phone)

				if err != nil || len(companyUsers) == 0 {
					loginCompanyUser.Reason = "公司异常"
				}

				if loginCompanyUser.CurrentCompanyId > 0 {
					isExitCompany := true

					for _, cUser := range companyUsers {
						if cUser.CompanyId == loginCompanyUser.CurrentCompanyId {
							companyUser = cUser

							isExitCompany = false
						}
					}

					if isExitCompany {
						companyUser = companyUsers[0]
					}
				} else {
					companyUser = companyUsers[0]
				}

				if companys, err := cuuc.crepo.List(ctx, 0, 0, 0, "", "", 0); err == nil {
					for _, company := range companys {
						if company.Status == 1 {
							if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
								company.Clue = clue
							}

							loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
								Id:          company.Id,
								CompanyName: company.Clue.CompanyName,
							})
						}
					}
				}

				loginCompanyUser.IsWhite = 1
			}

			if companyUser != nil {
				loginCompanyUser.Id = companyUser.Id
				loginCompanyUser.CompanyId = companyUser.CompanyId
				loginCompanyUser.Username = companyUser.Username
				loginCompanyUser.Job = companyUser.Job
				loginCompanyUser.Phone = companyUser.Phone
				loginCompanyUser.Role = companyUser.Role
				loginCompanyUser.RoleName = companyUser.RoleName
				loginCompanyUser.Role = companyUser.Role
			}

			if company != nil {
				loginCompanyUser.CompanyType = company.CompanyType
				loginCompanyUser.CompanyTypeName = company.CompanyTypeName
				loginCompanyUser.CompanyName = company.Clue.CompanyName
				loginCompanyUser.CompanyStartTime = company.StartTime
				loginCompanyUser.CompanyEndTime = company.EndTime
				loginCompanyUser.Accounts = company.Accounts
				loginCompanyUser.QianchuanAdvertisers = company.QianchuanAdvertisers
				loginCompanyUser.IsTermwork = company.IsTermwork
			}

			ltoken := tool.GetToken()

			loginCompanyUser.Token = ltoken

			if err := cuuc.repo.DeleteCache(ctx, "company:user:"+token); err != nil {
				return nil, CompanyLoginFailed
			}

			if err := cuuc.repo.DeleteCache(ctx, "company:user:"+phone); err != nil {
				return nil, CompanyLoginFailed
			}

			cacheData := make(map[string]string)
			cacheData["companyId"] = strconv.FormatInt(int64(companyId), 10)
			cacheData["phone"] = phone

			if err := cuuc.repo.SaveCacheHash(ctx, "company:user:"+ltoken, cacheData, cuuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
				return nil, CompanyLoginFailed
			}

			if err := cuuc.repo.SaveCacheString(ctx, "company:user:"+phone, ltoken, cuuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
				return nil, CompanyLoginFailed
			}

			return loginCompanyUser, nil
		}
	}
}

func (cuuc *CompanyUserUsecase) GetCompanyUser(ctx context.Context, token string) (*domain.LoginCompanyUser, error) {
	return cuuc.getCompanyUserByToken(ctx, token)
}

func (cuuc *CompanyUserUsecase) ListCompanyUserMenu(ctx context.Context, companyId uint64) ([]*domain.Menu, error) {
	company, err := cuuc.crepo.GetById(ctx, companyId)

	if err != nil {
		return nil, CompanyCompanyNotFound
	}

	if menus, err := cuuc.mrepo.ListByIds(ctx, strings.Split(company.MenuId, ",")); err != nil {
		return nil, CompanyLoginError
	} else {
		return menus, nil
	}
}

func (cuuc *CompanyUserUsecase) ChangeCompanyUserCompany(ctx context.Context, token string, companyId uint64) error {
	phone, err := cuuc.repo.GetCacheHash(ctx, "company:user:"+token, "phone")

	if err != nil {
		return CompanyChangeCompanyUserCompanyError
	}

	cacheData := make(map[string]string)
	cacheData["companyId"] = strconv.FormatUint(companyId, 10)
	cacheData["phone"] = phone

	if err := cuuc.repo.UpdateCacheHash(ctx, "company:user:"+token, cacheData); err != nil {
		return CompanyChangeCompanyUserCompanyError
	}

	if inCompanyUserCompany, err := cuuc.cucrepo.GetByPhone(ctx, phone); err != nil {
		inCompanyUserCompany = domain.NewCompanyUserCompany(ctx, companyId, phone)
		inCompanyUserCompany.SetCreateTime(ctx)
		inCompanyUserCompany.SetUpdateTime(ctx)

		if _, err := cuuc.cucrepo.Save(ctx, inCompanyUserCompany); err != nil {
			return CompanyChangeCompanyUserCompanyError
		}
	} else {
		inCompanyUserCompany.SetCompanyId(ctx, companyId)
		inCompanyUserCompany.SetUpdateTime(ctx)

		if _, err := cuuc.cucrepo.Update(ctx, inCompanyUserCompany); err != nil {
			return CompanyChangeCompanyUserCompanyError
		}
	}

	return nil
}

func (cuuc *CompanyUserUsecase) LogoutCompanyUser(ctx context.Context, token string) error {
	phone, err := cuuc.repo.GetCacheHash(ctx, "company:user:"+token, "phone")

	if err != nil {
		return CompanyLoginAbnormal
	}

	if err := cuuc.repo.DeleteCache(ctx, "company:user:"+token); err != nil {
		return CompanyLoginAbnormal
	}

	if err := cuuc.repo.DeleteCache(ctx, "company:user:"+phone); err != nil {
		return CompanyLoginAbnormal
	}

	return nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserQianchuanAdvertisers(ctx context.Context, userId, companyId uint64, isWhite uint8, day, keyword string) ([]*domain.CompanyUserQianchuanAdvertiser, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	list, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, keyword)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserQianchuanReportAdvertisers(ctx context.Context, userId, companyId uint64, isWhite uint8, day string) ([]*domain.CompanyUserQianchuanReportAdvertiserList, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	list := make([]*domain.CompanyUserQianchuanReportAdvertiserList, 0)
	payOrderAmountDataList := make([]*domain.CompanyUserQianchuanReportAdvertiserData, 0)
	statCostDataList := make([]*domain.CompanyUserQianchuanReportAdvertiserData, 0)

	startTime, _ := tool.StringToTime("2006-01-02 15:04", day+" 00:00")
	endTime, _ := tool.StringToTime("2006-01-02 15:04", day+" 23:59")

	for startTime.Before(endTime) {
		payOrderAmountDataList = append(payOrderAmountDataList, &domain.CompanyUserQianchuanReportAdvertiserData{
			Time:           tool.TimeToString("15:04", startTime),
			Value:          0.00,
			YesterdayValue: 0.00,
		})

		statCostDataList = append(statCostDataList, &domain.CompanyUserQianchuanReportAdvertiserData{
			Time:           tool.TimeToString("15:04", startTime),
			Value:          0.00,
			YesterdayValue: 0.00,
		})

		startTime = startTime.Add(time.Minute * 5)
	}

	qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, "")

	if err != nil {
		return nil, err
	}

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	if len(advertiserIds) > 0 {
		qianchuanReportAdvertisers, err := cuuc.qrarepo.List(ctx, strings.Join(advertiserIds, ","), day)

		if err != nil {
			return nil, CompanyQianchuanReportAdvertiserListError
		}

		for _, qianchuanReportAdvertiser := range qianchuanReportAdvertisers.Data.List {
			for _, l := range payOrderAmountDataList {
				if l.Time == qianchuanReportAdvertiser.Time {
					payOrderAmount, _ := strconv.ParseFloat(qianchuanReportAdvertiser.PayOrderAmount, 64)

					l.Value = payOrderAmount

					break
				}
			}

			for _, l := range statCostDataList {
				if l.Time == qianchuanReportAdvertiser.Time {
					statCost, _ := strconv.ParseFloat(qianchuanReportAdvertiser.StatCost, 64)

					l.Value = statCost

					break
				}
			}
		}

		tday, err := tool.StringToTime("2006-01-02", day)

		if err != nil {
			return nil, CompanyQianchuanAdvertiserListError
		}

		yday := tool.TimeToString("2006-01-02", tday.Add(time.Hour*-24))

		yqianchuanReportAdvertisers, err := cuuc.qrarepo.List(ctx, strings.Join(advertiserIds, ","), yday)

		if err != nil {
			return nil, CompanyQianchuanReportAdvertiserListError
		}

		for _, yqianchuanReportAdvertiser := range yqianchuanReportAdvertisers.Data.List {
			for _, l := range payOrderAmountDataList {
				if l.Time == yqianchuanReportAdvertiser.Time {
					payOrderAmount, _ := strconv.ParseFloat(yqianchuanReportAdvertiser.PayOrderAmount, 64)

					l.YesterdayValue = payOrderAmount

					break
				}
			}

			for _, l := range statCostDataList {
				if l.Time == yqianchuanReportAdvertiser.Time {
					statCost, _ := strconv.ParseFloat(yqianchuanReportAdvertiser.StatCost, 64)

					l.YesterdayValue = statCost

					break
				}
			}
		}
	}

	list = append(list, &domain.CompanyUserQianchuanReportAdvertiserList{
		Key:  "payOrderAmount",
		List: payOrderAmountDataList,
	})

	list = append(list, &domain.CompanyUserQianchuanReportAdvertiserList{
		Key:  "statCost",
		List: statCostDataList,
	})

	return list, nil
}

func (cuuc *CompanyUserUsecase) StatisticsCompanyUserDashboardQianchuanAdvertisers(ctx context.Context, userId, companyId uint64, isWhite uint8, day string) (*v1.StatisticsDashboardQianchuanAdvertisersReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, "")

	if err != nil {
		return nil, err
	}

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	var statistics *v1.StatisticsDashboardQianchuanAdvertisersReply

	if len(advertiserIds) > 0 {
		statistics, err = cuuc.qarepo.StatisticsDashboard(ctx, companyId, day, strings.Join(advertiserIds, ","))

		if err != nil {
			return nil, CompanyQianchuanReportAwemeListError
		}
	}

	return statistics, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserQianchuanReportAwemes(ctx context.Context, userId, companyId, pageNum, pageSize uint64, isDistinction uint32, isWhite uint8, day, keyword string) (*v1.ListQianchuanReportAwemesReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, "")

	if err != nil {
		return nil, err
	}

	var list *v1.ListQianchuanReportAwemesReply

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	if len(advertiserIds) > 0 {
		list, err = cuuc.qrawrepo.List(ctx, pageNum, pageSize, isDistinction, day, keyword, strings.Join(advertiserIds, ","))

		if err != nil {
			return nil, CompanyQianchuanReportAwemeListError
		}
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) StatisticsCompanyUserQianchuanReportAwemes(ctx context.Context, userId, companyId uint64, isWhite uint8, day string) (*v1.StatisticsQianchuanReportAwemesReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, "")

	if err != nil {
		return nil, err
	}

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	var statistics *v1.StatisticsQianchuanReportAwemesReply

	if len(advertiserIds) > 0 {
		statistics, err = cuuc.qrawrepo.Statistics(ctx, day, strings.Join(advertiserIds, ","))

		if err != nil {
			return nil, CompanyQianchuanReportAwemeListError
		}
	}

	return statistics, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserQianchuanReportProducts(ctx context.Context, userId, companyId, advertiserId, pageNum, pageSize uint64, isDistinction uint32, isWhite uint8, day, keyword string) (*v1.ListQianchuanReportProductsReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	kadvertiserId := ""

	if advertiserId > 0 {
		kadvertiserId = strconv.FormatUint(advertiserId, 10)
	}

	qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, kadvertiserId)

	if err != nil {
		return nil, err
	}

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	var list *v1.ListQianchuanReportProductsReply

	if len(advertiserIds) > 0 {
		list, err = cuuc.qrprepo.List(ctx, pageNum, pageSize, isDistinction, day, keyword, strings.Join(advertiserIds, ","))

		if err != nil {
			return nil, CompanyQianchuanReportProductListError
		}
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) StatisticsCompanyUserQianchuanReportProducts(ctx context.Context, userId, companyId uint64, isWhite uint8, day string) (*v1.StatisticsQianchuanReportProductsReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, "")

	if err != nil {
		return nil, err
	}

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	var statistics *v1.StatisticsQianchuanReportProductsReply

	if len(advertiserIds) > 0 {
		statistics, err = cuuc.qrprepo.Statistics(ctx, day, strings.Join(advertiserIds, ","))

		if err != nil {
			return nil, CompanyQianchuanReportProductListError
		}
	}

	return statistics, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserExternalQianchuanAdvertisers(ctx context.Context, userId, companyId, pageNum, pageSize uint64, isWhite uint8, startTime, endTime time.Time, keyword string) (*v1.ListExternalQianchuanAdvertisersReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	var wg sync.WaitGroup
	var historyQianchuanAdvertiserIds sync.Map

	startDay := tool.TimeToString("2006-01-02", startTime)
	endDay := tool.TimeToString("2006-01-02", endTime)

	endTime = endTime.Add(time.Hour * 24)

	for startTime.Before(endTime) {
		wg.Add(1)

		day := tool.TimeToString("2006-01-02", startTime)

		go func(day string) {
			if cuQianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, keyword); err == nil {
				advertiserIds := make([]string, 0)

				for _, cuQianchuanAdvertiser := range cuQianchuanAdvertisers {
					advertiserIds = append(advertiserIds, strconv.FormatUint(cuQianchuanAdvertiser.AdvertiserId, 10))
				}

				if len(advertiserIds) > 0 {
					historyQianchuanAdvertiserIds.Store(day, advertiserIds)
				}
			}

			wg.Done()
		}(day)

		startTime = startTime.Add(time.Hour * 24)
	}

	wg.Wait()

	advertiserIds := make([]string, 0)

	historyQianchuanAdvertiserIds.Range(func(k, v interface{}) bool {
		qianchuanAdvertiserIds := make([]string, 0)
		qianchuanAdvertiserIds = v.([]string)

		for _, qianchuanAdvertiserId := range qianchuanAdvertiserIds {
			isNotExist := true

			for _, advertiserId := range advertiserIds {
				if advertiserId == qianchuanAdvertiserId {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				advertiserIds = append(advertiserIds, qianchuanAdvertiserId)
			}
		}

		return true
	})

	if len(advertiserIds) == 0 {
		return &v1.ListExternalQianchuanAdvertisersReply{
			Code: 200,
			Data: &v1.ListExternalQianchuanAdvertisersReply_Data{
				PageNum:   pageNum,
				PageSize:  pageSize,
				Total:     0,
				TotalPage: 0,
				List:      make([]*v1.ListExternalQianchuanAdvertisersReply_Advertisers, 0),
			},
		}, nil
	}

	list, err := cuuc.qarepo.ListExternal(ctx, pageNum, pageSize, startDay, endDay, strings.Join(advertiserIds, ","))

	if err != nil {
		return nil, CompanyQianchuanAdvertiserListError
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserExternalSelectQianchuanAdvertisers(ctx context.Context, userId, companyId uint64, isWhite uint8, startTime, endTime time.Time) ([]*domain.ExternalSelectQianchuanAdvertiser, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	var wg sync.WaitGroup
	var historyQianchuanAdvertisers sync.Map

	endTime = endTime.Add(time.Hour * 24)

	for startTime.Before(endTime) {
		wg.Add(1)

		day := tool.TimeToString("2006-01-02", startTime)

		go func(day string) {
			if qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, ""); err == nil {
				historyQianchuanAdvertisers.Store(day, qianchuanAdvertisers)
			}

			wg.Done()
		}(day)

		startTime = startTime.Add(time.Hour * 24)
	}

	wg.Wait()

	list := make([]*domain.ExternalSelectQianchuanAdvertiser, 0)

	historyQianchuanAdvertisers.Range(func(k, v interface{}) bool {
		qianchuanAdvertisers := make([]*domain.CompanyUserQianchuanAdvertiser, 0)
		qianchuanAdvertisers = v.([]*domain.CompanyUserQianchuanAdvertiser)

		for _, qianchuanAdvertiser := range qianchuanAdvertisers {
			isNotExist := true

			for _, l := range list {
				if l.AdvertiserId == qianchuanAdvertiser.AdvertiserId {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				if qianchuanAdvertiser.AdvertiserId > 0 {
					if l := utf8.RuneCountInString(qianchuanAdvertiser.AdvertiserName); l > 0 {
						list = append(list, &domain.ExternalSelectQianchuanAdvertiser{
							AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
							AdvertiserName: qianchuanAdvertiser.AdvertiserName,
						})
					}
				}
			}
		}

		return true
	})

	return list, nil
}

func (cuuc *CompanyUserUsecase) StatisticsCompanyUserExternalQianchuanAdvertisers(ctx context.Context, userId, companyId uint64, isWhite uint8, startTime, endTime time.Time, keyword string) (*v1.StatisticsExternalQianchuanAdvertisersReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	var wg sync.WaitGroup
	var historyQianchuanAdvertiserIds sync.Map

	startDay := tool.TimeToString("2006-01-02", startTime)
	endDay := tool.TimeToString("2006-01-02", endTime)

	endTime = endTime.Add(time.Hour * 24)

	for startTime.Before(endTime) {
		wg.Add(1)

		day := tool.TimeToString("2006-01-02", startTime)

		go func(day string) {
			if cuQianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, keyword); err == nil {
				advertiserIds := make([]string, 0)

				for _, cuQianchuanAdvertiser := range cuQianchuanAdvertisers {
					advertiserIds = append(advertiserIds, strconv.FormatUint(cuQianchuanAdvertiser.AdvertiserId, 10))
				}

				if len(advertiserIds) > 0 {
					historyQianchuanAdvertiserIds.Store(day, advertiserIds)
				}
			}

			wg.Done()
		}(day)

		startTime = startTime.Add(time.Hour * 24)
	}

	wg.Wait()

	advertiserIds := make([]string, 0)

	historyQianchuanAdvertiserIds.Range(func(k, v interface{}) bool {
		qianchuanAdvertiserIds := make([]string, 0)
		qianchuanAdvertiserIds = v.([]string)

		for _, qianchuanAdvertiserId := range qianchuanAdvertiserIds {
			isNotExist := true

			for _, advertiserId := range advertiserIds {
				if advertiserId == qianchuanAdvertiserId {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				advertiserIds = append(advertiserIds, qianchuanAdvertiserId)
			}
		}

		return true
	})

	if len(advertiserIds) == 0 {
		return &v1.StatisticsExternalQianchuanAdvertisersReply{
			Code: 200,
			Data: &v1.StatisticsExternalQianchuanAdvertisersReply_Data{
				Statistics: make([]*v1.StatisticsExternalQianchuanAdvertisersReply_Statistics, 0),
			},
		}, nil
	}

	list, err := cuuc.qarepo.StatisticsExternal(ctx, startDay, endDay, strings.Join(advertiserIds, ","))

	if err != nil {
		return nil, CompanyQianchuanAdvertiserListError
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserSelectExternalQianchuanAds(ctx context.Context) (*v1.ListSelectExternalQianchuanAdsReply, error) {
	selects, err := cuuc.qadrepo.ListSelectExternal(ctx)

	if err != nil {
		return nil, CompanyDataError
	}

	return selects, nil
}

func (cuuc *CompanyUserUsecase) GetCompanyUserExternalQianchuanAds(ctx context.Context, userId, companyId, adId uint64, isWhite uint8, startDay, endDay string) ([]*v1.GetExternalQianchuanAdsReply_Ads, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	list := make([]*v1.GetExternalQianchuanAdsReply_Ads, 0)

	if startDay == endDay {
		qianchuanAds, err := cuuc.qadrepo.GetExternal(ctx, adId, startDay)

		if err != nil {
			return nil, CompanyQianchuanAdNotFound
		}

		for _, qianchuanAd := range qianchuanAds.Data.List {
			list = append(list, &v1.GetExternalQianchuanAdsReply_Ads{
				Time:                    qianchuanAd.Time,
				StatCost:                qianchuanAd.StatCost,
				Roi:                     qianchuanAd.Roi,
				PayOrderCount:           qianchuanAd.PayOrderCount,
				PayOrderAmount:          qianchuanAd.PayOrderAmount,
				ClickCnt:                qianchuanAd.ClickCnt,
				ShowCnt:                 qianchuanAd.ShowCnt,
				ConvertCnt:              qianchuanAd.ConvertCnt,
				ClickRate:               qianchuanAd.ClickRate,
				CpmPlatform:             qianchuanAd.CpmPlatform,
				DyFollow:                qianchuanAd.DyFollow,
				PayConvertRate:          qianchuanAd.PayConvertRate,
				ConvertCost:             qianchuanAd.ConvertCost,
				AveragePayOrderStatCost: qianchuanAd.AveragePayOrderStatCost,
				PayOrderAveragePrice:    qianchuanAd.PayOrderAveragePrice,
			})
		}
	} else {
		qianchuanAds, err := cuuc.qadrepo.GetExternalHistory(ctx, adId, startDay, endDay)

		if err != nil {
			return nil, CompanyQianchuanAdNotFound
		}

		for _, qianchuanAd := range qianchuanAds.Data.List {
			list = append(list, &v1.GetExternalQianchuanAdsReply_Ads{
				Time:                    qianchuanAd.Time,
				StatCost:                qianchuanAd.StatCost,
				Roi:                     qianchuanAd.Roi,
				PayOrderCount:           qianchuanAd.PayOrderCount,
				PayOrderAmount:          qianchuanAd.PayOrderAmount,
				ClickCnt:                qianchuanAd.ClickCnt,
				ShowCnt:                 qianchuanAd.ShowCnt,
				ConvertCnt:              qianchuanAd.ConvertCnt,
				ClickRate:               qianchuanAd.ClickRate,
				CpmPlatform:             qianchuanAd.CpmPlatform,
				DyFollow:                qianchuanAd.DyFollow,
				PayConvertRate:          qianchuanAd.PayConvertRate,
				ConvertCost:             qianchuanAd.ConvertCost,
				AveragePayOrderStatCost: qianchuanAd.AveragePayOrderStatCost,
				PayOrderAveragePrice:    qianchuanAd.PayOrderAveragePrice,
			})
		}
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserExternalQianchuanAds(ctx context.Context, userId, companyId, advertiserId, pageNum, pageSize uint64, isWhite uint8, startTime, endTime time.Time, keyword, filter, orderName, orderType string) (*v1.ListExternalQianchuanAdsReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	var wg sync.WaitGroup
	var historyQianchuanAdvertiserIds sync.Map

	startDay := tool.TimeToString("2006-01-02", startTime)
	endDay := tool.TimeToString("2006-01-02", endTime)

	endTime = endTime.Add(time.Hour * 24)

	for startTime.Before(endTime) {
		wg.Add(1)

		day := tool.TimeToString("2006-01-02", startTime)

		go func(day string) {
			if cuQianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, ""); err == nil {
				advertiserIds := make([]string, 0)

				for _, cuQianchuanAdvertiser := range cuQianchuanAdvertisers {
					advertiserIds = append(advertiserIds, strconv.FormatUint(cuQianchuanAdvertiser.AdvertiserId, 10))
				}

				if len(advertiserIds) > 0 {
					historyQianchuanAdvertiserIds.Store(day, advertiserIds)
				}
			}

			wg.Done()
		}(day)

		startTime = startTime.Add(time.Hour * 24)
	}

	wg.Wait()

	advertiserIds := make([]string, 0)

	historyQianchuanAdvertiserIds.Range(func(k, v interface{}) bool {
		qianchuanAdvertiserIds := make([]string, 0)
		qianchuanAdvertiserIds = v.([]string)

		for _, qianchuanAdvertiserId := range qianchuanAdvertiserIds {
			isNotExist := true

			for _, ladvertiserId := range advertiserIds {
				if ladvertiserId == qianchuanAdvertiserId {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				advertiserIds = append(advertiserIds, qianchuanAdvertiserId)
			}
		}

		return true
	})

	if advertiserId > 0 {
		isNotExist := true

		for _, ladvertiserId := range advertiserIds {
			if strconv.FormatUint(advertiserId, 10) == ladvertiserId {
				isNotExist = false

				break
			}
		}

		if !isNotExist {
			advertiserIds = []string{strconv.FormatUint(advertiserId, 10)}
		}

		if isNotExist {
			return &v1.ListExternalQianchuanAdsReply{
				Code: 200,
				Data: &v1.ListExternalQianchuanAdsReply_Data{
					PageNum:  pageNum,
					PageSize: pageSize,
					List:     make([]*v1.ListExternalQianchuanAdsReply_Ads, 0),
				},
			}, nil
		}
	}

	if len(advertiserIds) == 0 {
		return &v1.ListExternalQianchuanAdsReply{
			Code: 200,
			Data: &v1.ListExternalQianchuanAdsReply_Data{
				PageNum:  pageNum,
				PageSize: pageSize,
				List:     make([]*v1.ListExternalQianchuanAdsReply_Ads, 0),
			},
		}, nil
	}

	list, err := cuuc.qadrepo.ListExternal(ctx, pageNum, pageSize, startDay, endDay, keyword, strings.Join(advertiserIds, ","), filter, orderName, orderType)

	if err != nil {
		return nil, CompanyQianchuanAdListError
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) StatisticsCompanyUserExternalQianchuanAds(ctx context.Context, userId, companyId, advertiserId uint64, isWhite uint8, startTime, endTime time.Time, keyword, filter string) (*v1.StatisticsExternalQianchuanAdsReply, error) {
	var wg sync.WaitGroup
	var historyQianchuanAdvertiserIds sync.Map

	startDay := tool.TimeToString("2006-01-02", startTime)
	endDay := tool.TimeToString("2006-01-02", endTime)

	endTime = endTime.Add(time.Hour * 24)

	for startTime.Before(endTime) {
		wg.Add(1)

		day := tool.TimeToString("2006-01-02", startTime)

		go func(day string) {
			if cuQianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, ""); err == nil {
				advertiserIds := make([]string, 0)

				for _, cuQianchuanAdvertiser := range cuQianchuanAdvertisers {
					advertiserIds = append(advertiserIds, strconv.FormatUint(cuQianchuanAdvertiser.AdvertiserId, 10))
				}

				if len(advertiserIds) > 0 {
					historyQianchuanAdvertiserIds.Store(day, advertiserIds)
				}
			}

			wg.Done()
		}(day)

		startTime = startTime.Add(time.Hour * 24)
	}

	wg.Wait()

	advertiserIds := make([]string, 0)

	historyQianchuanAdvertiserIds.Range(func(k, v interface{}) bool {
		qianchuanAdvertiserIds := make([]string, 0)
		qianchuanAdvertiserIds = v.([]string)

		for _, qianchuanAdvertiserId := range qianchuanAdvertiserIds {
			isNotExist := true

			for _, ladvertiserId := range advertiserIds {
				if ladvertiserId == qianchuanAdvertiserId {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				advertiserIds = append(advertiserIds, qianchuanAdvertiserId)
			}
		}

		return true
	})

	if advertiserId > 0 {
		isNotExist := true

		for _, ladvertiserId := range advertiserIds {
			if strconv.FormatUint(advertiserId, 10) == ladvertiserId {
				isNotExist = false

				break
			}
		}

		if !isNotExist {
			advertiserIds = []string{strconv.FormatUint(advertiserId, 10)}
		}

		if isNotExist {
			return &v1.StatisticsExternalQianchuanAdsReply{
				Code: 200,
				Data: &v1.StatisticsExternalQianchuanAdsReply_Data{
					Statistics: make([]*v1.StatisticsExternalQianchuanAdsReply_Statistics, 0),
				},
			}, nil
		}
	}

	if len(advertiserIds) == 0 {
		return &v1.StatisticsExternalQianchuanAdsReply{
			Code: 200,
			Data: &v1.StatisticsExternalQianchuanAdsReply_Data{
				Statistics: make([]*v1.StatisticsExternalQianchuanAdsReply_Statistics, 0),
			},
		}, nil
	}

	list, err := cuuc.qadrepo.StatisticsExternal(ctx, startDay, endDay, keyword, strings.Join(advertiserIds, ","), filter)

	if err != nil {
		return nil, CompanyQianchuanAdListError
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) GetCompanyUserIsCompleteSyncData(ctx context.Context, companyId uint64, day string) (*v1.ListQianchuanAdvertiserHistorysReply, error) {
	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, 0, companyId, 1, day, "")

	if err != nil {
		return &v1.ListQianchuanAdvertiserHistorysReply{
			Code: 200,
			Data: &v1.ListQianchuanAdvertiserHistorysReply_Data{
				List: make([]*v1.ListQianchuanAdvertiserHistorysReply_QianchuanAdvertisers, 0),
			},
		}, nil
	}

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	if len(advertiserIds) == 0 {
		return &v1.ListQianchuanAdvertiserHistorysReply{
			Code: 200,
			Data: &v1.ListQianchuanAdvertiserHistorysReply_Data{
				List: make([]*v1.ListQianchuanAdvertiserHistorysReply_QianchuanAdvertisers, 0),
			},
		}, nil
	}

	list, err := cuuc.qahrepo.List(ctx, day, strings.Join(advertiserIds, ","))

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserQianchuanSearch(ctx context.Context, userId, companyId uint64, isWhite uint8, day, keyword, searchType string) ([]*domain.CompanyUserQianchuanDataSearch, error) {
	list := make([]*domain.CompanyUserQianchuanDataSearch, 0)

	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	var qianchuanAdvertisers []*domain.CompanyUserQianchuanAdvertiser
	var qerr error

	if searchType == "advertiser" {
		qianchuanAdvertisers, qerr = cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, keyword)
	} else {
		qianchuanAdvertisers, qerr = cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, "")
	}

	if qerr != nil {
		return nil, qerr
	}

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	if len(advertiserIds) == 0 {
		return list, nil
	}

	if searchType == "advertiser" {
		if len(qianchuanAdvertisers) > 40 {
			qianchuanAdvertisers = qianchuanAdvertisers[0:40]
		}

		for _, lqianchuanAdvertiser := range qianchuanAdvertisers {
			list = append(list, &domain.CompanyUserQianchuanDataSearch{
				Id:   lqianchuanAdvertiser.AdvertiserId,
				Name: lqianchuanAdvertiser.AdvertiserName,
			})
		}
	} else if searchType == "ad" {
		qianchuanAds, err := cuuc.qadrepo.List(ctx, 1, 40, day, keyword, strings.Join(advertiserIds, ","))

		if err != nil {
			return nil, CompanyUserQianchuanSearchListError
		}

		for _, qianchuanAd := range qianchuanAds.Data.List {
			list = append(list, &domain.CompanyUserQianchuanDataSearch{
				Id:   qianchuanAd.AdId,
				Name: qianchuanAd.Name,
			})
		}
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserExternalHistoryQianchuanSearch(ctx context.Context, userId, companyId uint64, isWhite uint8, startTime, endTime time.Time, keyword, searchType string) ([]*domain.CompanyUserQianchuanDataSearch, error) {
	list := make([]*domain.CompanyUserQianchuanDataSearch, 0)

	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	var wg sync.WaitGroup
	var qianchuanSearches sync.Map

	endTime = endTime.Add(time.Hour * 24)

	for startTime.Before(endTime) {
		wg.Add(1)

		day := tool.TimeToString("2006-01-02", startTime)

		go func(day string) {
			if qianchuanSearch, err := cuuc.ListCompanyUserQianchuanSearch(ctx, userId, companyId, isWhite, day, keyword, searchType); err == nil {
				qianchuanSearches.Store(day, qianchuanSearch)
			}

			wg.Done()
		}(day)

		startTime = startTime.Add(time.Hour * 24)
	}

	wg.Wait()

	qianchuanSearches.Range(func(k, v interface{}) bool {
		cqianchuanSearches := make([]*domain.CompanyUserQianchuanDataSearch, 0)
		cqianchuanSearches = v.([]*domain.CompanyUserQianchuanDataSearch)

		for _, cqianchuanSearch := range cqianchuanSearches {
			isNotExist := true

			for _, l := range list {
				if l.Id == cqianchuanSearch.Id {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				list = append(list, &domain.CompanyUserQianchuanDataSearch{
					Id:   cqianchuanSearch.Id,
					Name: cqianchuanSearch.Name,
				})
			}
		}

		return true
	})

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserExternalQianchuanHistorySearch(ctx context.Context, userId, companyId uint64, isWhite uint8, day, searchType string) ([]*domain.CompanyUserQianchuanDataSearch, error) {
	list := make([]*domain.CompanyUserQianchuanDataSearch, 0)

	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	dayTime, _ := tool.StringToTime("2006-01-02", day)
	uiday, _ := strconv.ParseUint(tool.TimeToString("20060102", dayTime), 10, 64)

	qianchuanAdvertisers, err := cuuc.getCompanyUserQianchuanAdvertisers(ctx, userId, companyId, isWhite, day, "")

	if err != nil {
		return nil, err
	}

	advertiserIds := make([]string, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers {
		advertiserIds = append(advertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
	}

	if len(advertiserIds) == 0 {
		return list, nil
	}

	qianchuanSearches, err := cuuc.cqsrepo.List(ctx, companyId, uint32(uiday), 0, searchType)

	if err != nil {
		return nil, CompanyUserQianchuanSearchListError
	}

	if searchType == "advertiser" {
		if len(qianchuanAdvertisers) > 40 {
			qianchuanAdvertisers = qianchuanAdvertisers[0:40]
		}

		for _, qianchuanSearch := range qianchuanSearches {
			if usearchValue, err := strconv.ParseUint(qianchuanSearch.SearchValue, 10, 64); err == nil {
				for _, lqianchuanAdvertiser := range qianchuanAdvertisers {
					if usearchValue == lqianchuanAdvertiser.AdvertiserId {
						list = append(list, &domain.CompanyUserQianchuanDataSearch{
							Id:   lqianchuanAdvertiser.AdvertiserId,
							Name: lqianchuanAdvertiser.AdvertiserName,
						})
					}
				}
			}
		}
	} else if searchType == "ad" {
		qianchuanAds, err := cuuc.qadrepo.List(ctx, 1, 40, day, "", strings.Join(advertiserIds, ","))

		if err != nil {
			return nil, CompanyUserQianchuanSearchListError
		}

		for _, qianchuanSearch := range qianchuanSearches {
			if usearchValue, err := strconv.ParseUint(qianchuanSearch.SearchValue, 10, 64); err == nil {
				for _, qianchuanAd := range qianchuanAds.Data.List {
					if usearchValue == qianchuanAd.AdId {
						list = append(list, &domain.CompanyUserQianchuanDataSearch{
							Id:   qianchuanAd.AdId,
							Name: qianchuanAd.Name,
						})
					}
				}
			}
		}
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserExternalHistoryQianchuanHistorySearch(ctx context.Context, userId, companyId uint64, isWhite uint8, startTime, endTime time.Time, searchType string) ([]*domain.CompanyUserQianchuanDataSearch, error) {
	list := make([]*domain.CompanyUserQianchuanDataSearch, 0)
	qianchuanList := make([]*domain.CompanyUserQianchuanDataSearch, 0)

	if _, err := cuuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	uistartTime, _ := strconv.ParseUint(tool.TimeToString("20060102", startTime), 10, 64)
	uiendTime, _ := strconv.ParseUint(tool.TimeToString("20060102", endTime), 10, 64)

	var wg sync.WaitGroup
	var qianchuanDatas sync.Map

	endTime = endTime.Add(time.Hour * 24)

	for startTime.Before(endTime) {
		wg.Add(1)

		day := tool.TimeToString("2006-01-02", startTime)

		go func(day string) {
			if qianchuanData, err := cuuc.ListCompanyUserQianchuanSearch(ctx, userId, companyId, isWhite, day, "", searchType); err == nil {
				qianchuanDatas.Store(day, qianchuanData)
			}

			wg.Done()
		}(day)

		startTime = startTime.Add(time.Hour * 24)
	}

	wg.Wait()

	qianchuanDatas.Range(func(k, v interface{}) bool {
		cqianchuanSearches := make([]*domain.CompanyUserQianchuanDataSearch, 0)
		cqianchuanSearches = v.([]*domain.CompanyUserQianchuanDataSearch)

		for _, cqianchuanSearch := range cqianchuanSearches {
			isNotExist := true

			for _, l := range qianchuanList {
				if l.Id == cqianchuanSearch.Id {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				qianchuanList = append(qianchuanList, &domain.CompanyUserQianchuanDataSearch{
					Id:   cqianchuanSearch.Id,
					Name: cqianchuanSearch.Name,
				})
			}
		}

		return true
	})

	qianchuanSearches, err := cuuc.cqsrepo.List(ctx, companyId, uint32(uistartTime), uint32(uiendTime), searchType)

	if err != nil {
		return nil, CompanyUserQianchuanSearchListError
	}

	for _, qianchuanSearch := range qianchuanSearches {
		if usearchValue, err := strconv.ParseUint(qianchuanSearch.SearchValue, 10, 64); err == nil {
			for _, qianchuanData := range qianchuanList {
				if usearchValue == qianchuanData.Id {
					list = append(list, &domain.CompanyUserQianchuanDataSearch{
						Id:   qianchuanData.Id,
						Name: qianchuanData.Name,
					})
				}
			}
		}
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) getCompanyUserQianchuanAdvertisers(ctx context.Context, userId, companyId uint64, isWhite uint8, day, keyword string) ([]*domain.CompanyUserQianchuanAdvertiser, error) {
	cuQianchuanAdvertisers := make([]*domain.CompanyUserQianchuanAdvertiser, 0)

	advertiserIds := make([]uint64, 0)

	tday, err := tool.StringToTime("2006-01-02", day)

	if err != nil {
		return nil, CompanyQianchuanAdvertiserListError
	}

	uiday, err := strconv.ParseUint(tool.TimeToString("20060102", tday), 10, 64)

	if err != nil {
		return nil, CompanyQianchuanAdvertiserListError
	}

	qianchuanAdvertisers, err := cuuc.qarepo.ListByDays(ctx, companyId, day)

	if err != nil {
		return nil, CompanyQianchuanAdvertiserListError
	}

	for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
		if qianchuanAdvertiser.Status == 1 {
			advertiserIds = append(advertiserIds, qianchuanAdvertiser.AdvertiserId)
		}
	}

	if isWhite != 1 {
		companyUser, err := cuuc.repo.GetById(ctx, userId, companyId)

		if err != nil {
			return nil, CompanyCompanyUserNotFound
		}

		if companyUser.Role == 2 {
			companyUserRoles, err := cuuc.currepo.ListByUserIdAndCompanyIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, uint32(uiday))

			if err != nil {
				return nil, CompanyCompanyUserRoleListError
			}

			if len(companyUserRoles) > 0 {
				uadvertiserIds := make([]uint64, 0)

				for _, companyUserRole := range companyUserRoles {
					for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
						if companyUserRole.AdvertiserId == qianchuanAdvertiser.AdvertiserId && qianchuanAdvertiser.Status == 1 {
							uadvertiserIds = append(uadvertiserIds, companyUserRole.AdvertiserId)

							break
						}
					}
				}

				advertiserIds = uadvertiserIds
			}
		} else if companyUser.Role == 3 {
			companyUserRoles, err := cuuc.currepo.ListByUserIdAndCompanyIdAndDay(ctx, companyUser.Id, companyUser.CompanyId, uint32(uiday))

			if err != nil {
				return nil, CompanyCompanyUserRoleListError
			}

			uadvertiserIds := make([]uint64, 0)

			for _, companyUserRole := range companyUserRoles {
				for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
					if companyUserRole.AdvertiserId == qianchuanAdvertiser.AdvertiserId && qianchuanAdvertiser.Status == 1 {
						uadvertiserIds = append(uadvertiserIds, companyUserRole.AdvertiserId)

						break
					}
				}
			}

			advertiserIds = uadvertiserIds
		}
	}

	kqianchuanAdvertisers, err := cuuc.qarepo.List(ctx, companyId, 0, 0, keyword, "", "")

	if err != nil {
		return nil, CompanyQianchuanAdvertiserListError
	}

	for _, advertiserId := range advertiserIds {
		for _, kqianchuanAdvertiser := range kqianchuanAdvertisers.Data.List {
			if advertiserId == kqianchuanAdvertiser.AdvertiserId {
				cuQianchuanAdvertisers = append(cuQianchuanAdvertisers, &domain.CompanyUserQianchuanAdvertiser{
					AdvertiserId:   kqianchuanAdvertiser.AdvertiserId,
					CompanyId:      kqianchuanAdvertiser.CompanyId,
					AccountId:      kqianchuanAdvertiser.AccountId,
					AdvertiserName: kqianchuanAdvertiser.AdvertiserName,
					CompanyName:    kqianchuanAdvertiser.CompanyName,
				})
			}
		}
	}

	return cuQianchuanAdvertisers, nil
}

func (cuuc *CompanyUserUsecase) getCompanyUserByToken(ctx context.Context, token string) (*domain.LoginCompanyUser, error) {
	phone, err := cuuc.repo.GetCacheHash(ctx, "company:user:"+token, "phone")

	if err != nil {
		return nil, CompanyLoginAbnormal
	}

	cCompanyId, err := cuuc.repo.GetCacheHash(ctx, "company:user:"+token, "companyId")

	if err != nil {
		return nil, CompanyLoginAbnormal
	}

	companyId, err := strconv.Atoi(cCompanyId)

	if err != nil {
		return nil, CompanyLoginAbnormal
	}

	loginCompanyUser := &domain.LoginCompanyUser{}

	if companyId == 0 {
		if _, err := cuuc.cuwrepo.GetByPhone(ctx, phone); err != nil {
			if companyUsers, err := cuuc.repo.ListByPhone(ctx, phone); err == nil {
				for _, companyUser := range companyUsers {
					if company, err := cuuc.crepo.GetById(ctx, companyUser.CompanyId); err == nil {
						if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
							company.Clue = clue
						}

						if company.Status == 1 && companyUser.Status == 1 {
							loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
								Id:          company.Id,
								CompanyName: company.Clue.CompanyName,
							})
						}
					}
				}
			}
		} else {
			if companys, err := cuuc.crepo.List(ctx, 0, 0, 0, "", "", 0); err == nil {
				for _, company := range companys {
					if company.Status == 1 {
						if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
							company.Clue = clue
						}

						loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
							Id:          company.Id,
							CompanyName: company.Clue.CompanyName,
						})
					}
				}
			}
		}

		loginCompanyUser.Token = token

		return loginCompanyUser, nil
	} else {
		company, err := cuuc.crepo.GetById(ctx, uint64(companyId))

		if err != nil {
			loginCompanyUser.Reason = "公司异常"
		} else {
			if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
				company.Clue = clue

				if company.Status == 0 {
					loginCompanyUser.Reason = "公司被禁用"
				} else if company.Status == 2 {
					loginCompanyUser.Reason = "账户权限已过期"
				} else {
					loginCompanyUser.CurrentCompanyId = company.Id
				}
			} else {
				loginCompanyUser.Reason = "公司异常"
			}
		}

		var companyUser *domain.CompanyUser

		if _, err := cuuc.cuwrepo.GetByPhone(ctx, phone); err != nil {
			if company != nil {
				companyUser, err = cuuc.repo.GetByPhone(ctx, company.Id, phone)

				if err != nil {
					loginCompanyUser.Reason = "公司异常"
				} else {
					if companyUser.Status == 0 {
						loginCompanyUser.Reason = "账户被禁用"
					}
				}
			}

			if companyUsers, err := cuuc.repo.ListByPhone(ctx, phone); err == nil {
				for _, companyUser := range companyUsers {
					if company, err := cuuc.crepo.GetById(ctx, companyUser.CompanyId); err == nil {
						if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
							company.Clue = clue
						}

						if company.Status == 1 && companyUser.Status == 1 {
							loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
								Id:          company.Id,
								CompanyName: company.Clue.CompanyName,
							})
						}
					}
				}
			}

			loginCompanyUser.IsWhite = 0
		} else {
			companyUsers, err := cuuc.repo.ListByPhone(ctx, phone)

			if err != nil || len(companyUsers) == 0 {
				loginCompanyUser.Reason = "公司异常"
			}

			if loginCompanyUser.CurrentCompanyId > 0 {
				isExitCompany := true

				for _, cUser := range companyUsers {
					if cUser.CompanyId == loginCompanyUser.CurrentCompanyId {
						companyUser = cUser

						isExitCompany = false
					}
				}

				if isExitCompany {
					companyUser = companyUsers[0]
				}
			} else {
				companyUser = companyUsers[0]
			}

			if companys, err := cuuc.crepo.List(ctx, 0, 0, 0, "", "", 0); err == nil {
				for _, company := range companys {
					if company.Status == 1 {
						if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
							company.Clue = clue
						}

						loginCompanyUser.UserCompany = append(loginCompanyUser.UserCompany, &domain.LoginCompanyUserCompany{
							Id:          company.Id,
							CompanyName: company.Clue.CompanyName,
						})
					}
				}
			}

			loginCompanyUser.IsWhite = 1
		}

		if companyUser != nil {
			loginCompanyUser.Id = companyUser.Id
			loginCompanyUser.CompanyId = companyUser.CompanyId
			loginCompanyUser.Username = companyUser.Username
			loginCompanyUser.Job = companyUser.Job
			loginCompanyUser.Phone = companyUser.Phone
			loginCompanyUser.Role = companyUser.Role
			loginCompanyUser.RoleName = companyUser.RoleName
			loginCompanyUser.Role = companyUser.Role
		}

		if company != nil {
			loginCompanyUser.CompanyType = company.CompanyType
			loginCompanyUser.CompanyTypeName = company.CompanyTypeName
			loginCompanyUser.CompanyName = company.Clue.CompanyName
			loginCompanyUser.CompanyStartTime = company.StartTime
			loginCompanyUser.CompanyEndTime = company.EndTime
			loginCompanyUser.Accounts = company.Accounts
			loginCompanyUser.QianchuanAdvertisers = company.QianchuanAdvertisers
			loginCompanyUser.IsTermwork = company.IsTermwork
		}

		loginCompanyUser.Token = token

		return loginCompanyUser, nil
	}
}

func (cuuc *CompanyUserUsecase) listCompanyUsersByPhone(ctx context.Context, phone string) ([]*domain.CompanyUser, error) {
	companyUsers, err := cuuc.repo.ListByPhone(ctx, phone)

	if err != nil {
		return nil, CompanyLoginPhoneNotExist
	}

	list := make([]*domain.CompanyUser, 0)

	for _, companyUser := range companyUsers {
		if company, err := cuuc.crepo.GetById(ctx, companyUser.CompanyId); err == nil {
			companyUser.Company = company

			if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
				companyUser.Company.Clue = clue
			}
		}
	}

	if len(list) == 0 {
		return nil, CompanyLoginPhoneNotExist
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) getCompanyUsersById(ctx context.Context, id, companyId uint64) (*domain.CompanyUser, error) {
	companyUser, err := cuuc.repo.GetById(ctx, id, companyId)

	if err != nil {
		return nil, err
	}

	if company, err := cuuc.crepo.GetById(ctx, companyUser.CompanyId); err == nil {
		companyUser.Company = company

		if clue, err := cuuc.clrpeo.GetById(ctx, company.ClueId); err == nil {
			companyUser.Company.Clue = clue
		}
	}

	if _, err := cuuc.cuwrepo.GetByPhone(ctx, companyUser.Phone); err != nil {
		companyUser.IsWhite = 0
	} else {
		companyUser.IsWhite = 1
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) verifyCompanyUserLimit(ctx context.Context, company *domain.Company, userId uint64) error {
	totalNums := 0

	if companyUsers, err := cuuc.repo.ListByCompanyId(ctx, company.Id); err == nil {
		for _, companyUser := range companyUsers {
			if companyUser.Status == 1 && userId != companyUser.Id {
				totalNums += 1
			}
		}
	}

	if totalNums >= int(company.Accounts) {
		return CompanyCompanyUserLimitError
	}

	return nil
}
