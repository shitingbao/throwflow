package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/oceanengine"
	account2 "douyin/internal/pkg/oceanengine/account"
	"douyin/internal/pkg/oceanengine/oauth2"
	"douyin/internal/pkg/tool"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"math"
	"strconv"
	"time"
)

var (
	DouyinOceanengineAccountSyncError             = errors.InternalServer("DOUYIN_OCEANENGINE_ACCOUNT_SYNC_ERROR", "同步千川授权账户失败")
	DouyinOceanengineAccountCreateError           = errors.InternalServer("DOUYIN_OCEANENGINE_ACCOUNT_CREATE_ERROR", "巨量引擎获取授权账户信息失败")
	DouyinOceanengineAccountCreateAdvertiserError = errors.InternalServer("DOUYIN_OCEANENGINE_ACCOUNT_CREATE_ADVERTISER_ERROR", "巨量引擎获取授权账户信息失败,授权账号角色:广告账户")
)

type OceanengineAccountRepo interface {
	GetById(context.Context, uint64, uint64) (*domain.OceanengineAccount, error)
	GetBycompanyIdAndAccountId(context.Context, uint64, uint64) (*domain.OceanengineAccount, error)
	Save(context.Context, *domain.OceanengineAccount) (*domain.OceanengineAccount, error)
	Update(context.Context, *domain.OceanengineAccount) (*domain.OceanengineAccount, error)
	DeleteByCompanyIdAndAccountId(context.Context, uint64, uint64) error
}

type OceanengineAccountUsecase struct {
	repo    OceanengineAccountRepo
	ocrepo  OceanengineConfigRepo
	oatrepo OceanengineAccountTokenRepo
	qarepo  QianchuanAdvertiserRepo
	qasrepo QianchuanAdvertiserStatusRepo
	qahrepo QianchuanAdvertiserHistoryRepo
	oalrepo OceanengineApiLogRepo
	crepo   CompanyRepo
	tm      Transaction
	conf    *conf.Data
	log     *log.Helper
}

func NewOceanengineAccountUsecase(repo OceanengineAccountRepo, ocrepo OceanengineConfigRepo, oatrepo OceanengineAccountTokenRepo, qarepo QianchuanAdvertiserRepo, qasrepo QianchuanAdvertiserStatusRepo, qahrepo QianchuanAdvertiserHistoryRepo, oalrepo OceanengineApiLogRepo, crepo CompanyRepo, tm Transaction, conf *conf.Data, logger log.Logger) *OceanengineAccountUsecase {
	return &OceanengineAccountUsecase{repo: repo, ocrepo: ocrepo, oatrepo: oatrepo, qarepo: qarepo, qasrepo: qasrepo, qahrepo: qahrepo, oalrepo: oalrepo, crepo: crepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (oauc *OceanengineAccountUsecase) CreateOceanengineAccounts(ctx context.Context, companyId uint64, appId, authCode string) error {
	if _, err := oauc.crepo.GetById(ctx, companyId); err != nil {
		return err
	}

	oceanengineConfig, err := oauc.ocrepo.GetByAppId(ctx, appId)

	if err != nil {
		return DouyinOceanengineConfigNotFound
	}

	accessToken, err := oauth2.AccessToken(oceanengineConfig.AppId, oceanengineConfig.AppSecret, authCode)

	if err != nil {
		inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, companyId, 0, 0, 0, 0, oceanengineConfig.AppId, "", err.Error())
		inOceanengineApiLog.SetCreateTime(ctx)

		oauc.oalrepo.Save(ctx, inOceanengineApiLog)

		return DouyinOceanengineAccountCreateError
	}

	accounts, err := account2.ListAdvertiser(oceanengineConfig.AppId, oceanengineConfig.AppSecret, accessToken.Data.AccessToken)

	if err != nil {
		inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, companyId, 0, 0, 0, 0, oceanengineConfig.AppId, accessToken.Data.AccessToken, err.Error())
		inOceanengineApiLog.SetCreateTime(ctx)

		oauc.oalrepo.Save(ctx, inOceanengineApiLog)

		return DouyinOceanengineAccountCreateError
	}

	day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

	err = oauc.tm.InTx(ctx, func(ctx context.Context) error {
		for _, account := range accounts.Data.List {
			if !account.IsValid {
				if err := oauc.repo.DeleteByCompanyIdAndAccountId(ctx, companyId, account.AdvertiserId); err != nil {
					return err
				}

				if err := oauc.oatrepo.DeleteByCompanyIdAndAccountId(ctx, companyId, account.AdvertiserId); err != nil {
					return err
				}

				if err := oauc.qarepo.DeleteByCompanyIdAndAccountId(ctx, companyId, account.AdvertiserId); err != nil {
					return err
				}
			} else {
				if account.AccountRole == "PLATFORM_ROLE_QIANCHUAN_AGENT" || account.AccountRole == "PLATFORM_ROLE_SHOP_ACCOUNT" {
					if inOceanengineAccount, err := oauc.repo.GetBycompanyIdAndAccountId(ctx, companyId, account.AdvertiserId); err != nil {
						inOceanengineAccount = domain.NewOceanengineAccount(ctx, oceanengineConfig.AppId, account.AdvertiserName, account.AccountRole, companyId, account.AdvertiserId, 1)
						inOceanengineAccount.SetCreateTime(ctx)
						inOceanengineAccount.SetUpdateTime(ctx)

						if _, err := oauc.repo.Save(ctx, inOceanengineAccount); err != nil {
							return err
						}
					} else {
						inOceanengineAccount.SetAppId(ctx, oceanengineConfig.AppId)
						inOceanengineAccount.SetAccountName(ctx, account.AdvertiserName)
						inOceanengineAccount.SetAccountRole(ctx, account.AccountRole)
						inOceanengineAccount.SetUpdateTime(ctx)

						if _, err := oauc.repo.Update(ctx, inOceanengineAccount); err != nil {
							return err
						}
					}

					if inOceanengineAccountToken, err := oauc.oatrepo.GetByCompanyIdAndAccountId(ctx, companyId, account.AdvertiserId); err != nil {
						inOceanengineAccountToken = domain.NewOceanengineAccountToken(ctx, oceanengineConfig.AppId, accessToken.Data.AccessToken, accessToken.Data.RefreshToken, companyId, account.AdvertiserId, accessToken.Data.ExpiresIn, accessToken.Data.RefreshTokenExpiresIn)
						inOceanengineAccountToken.SetCreateTime(ctx)
						inOceanengineAccountToken.SetUpdateTime(ctx)

						if _, err := oauc.oatrepo.Save(ctx, inOceanengineAccountToken); err != nil {
							return err
						}
					} else {
						inOceanengineAccountToken.SetAppId(ctx, oceanengineConfig.AppId)
						inOceanengineAccountToken.SetAccessToken(ctx, accessToken.Data.AccessToken)
						inOceanengineAccountToken.SetRefreshToken(ctx, accessToken.Data.RefreshToken)
						inOceanengineAccountToken.SetExpiresIn(ctx, accessToken.Data.ExpiresIn)
						inOceanengineAccountToken.SetRefreshTokenExpiresIn(ctx, accessToken.Data.RefreshTokenExpiresIn)
						inOceanengineAccountToken.SetUpdateTime(ctx)

						if _, err := oauc.oatrepo.Update(ctx, inOceanengineAccountToken); err != nil {
							return err
						}
					}

					advertiserIds := make([]string, 0)

					if account.AccountRole == "PLATFORM_ROLE_QIANCHUAN_AGENT" {
						advertisers, err := account2.ListAgentAdvertiser(account.AdvertiserId, accessToken.Data.AccessToken, 1)

						if err != nil {
							inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, companyId, account.AdvertiserId, 0, 0, 0, oceanengineConfig.AppId, accessToken.Data.AccessToken, err.Error())
							inOceanengineApiLog.SetCreateTime(ctx)

							oauc.oalrepo.Save(ctx, inOceanengineApiLog)

							return err
						}

						for _, advertiser := range advertisers.Data.List {
							advertiserIds = append(advertiserIds, strconv.FormatUint(*advertiser, 10))
						}

						if advertisers.Data.PageInfo.TotalPage > 1 {
							var page uint32

							for page = 2; page <= advertisers.Data.PageInfo.TotalPage; page++ {
								if advertisers, err := account2.ListAgentAdvertiser(account.AdvertiserId, accessToken.Data.AccessToken, page); err == nil {
									for _, advertiser := range advertisers.Data.List {
										advertiserIds = append(advertiserIds, strconv.FormatUint(*advertiser, 10))
									}
								} else {
									inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, companyId, account.AdvertiserId, 0, 0, 0, oceanengineConfig.AppId, accessToken.Data.AccessToken, err.Error())
									inOceanengineApiLog.SetCreateTime(ctx)

									oauc.oalrepo.Save(ctx, inOceanengineApiLog)

									return err
								}
							}
						}
					} else if account.AccountRole == "PLATFORM_ROLE_SHOP_ACCOUNT" {
						advertisers, err := account2.ListShopAdvertiser(account.AdvertiserId, accessToken.Data.AccessToken, 1)

						if err != nil {
							inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, companyId, account.AdvertiserId, 0, 0, 0, oceanengineConfig.AppId, accessToken.Data.AccessToken, err.Error())
							inOceanengineApiLog.SetCreateTime(ctx)

							oauc.oalrepo.Save(ctx, inOceanengineApiLog)

							return err
						}

						for _, advertiser := range advertisers.Data.List {
							advertiserIds = append(advertiserIds, strconv.FormatUint(*advertiser, 10))
						}

						if advertisers.Data.PageInfo.TotalPage > 1 {
							var page uint32

							for page = 2; page <= advertisers.Data.PageInfo.TotalPage; page++ {
								if advertisers, err := account2.ListShopAdvertiser(account.AdvertiserId, accessToken.Data.AccessToken, page); err == nil {
									for _, advertiser := range advertisers.Data.List {
										advertiserIds = append(advertiserIds, strconv.FormatUint(*advertiser, 10))
									}
								} else {
									inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, companyId, account.AdvertiserId, 0, 0, 0, oceanengineConfig.AppId, accessToken.Data.AccessToken, err.Error())
									inOceanengineApiLog.SetCreateTime(ctx)

									oauc.oalrepo.Save(ctx, inOceanengineApiLog)

									return err
								}
							}
						}
					}

					totalPage := int(math.Ceil(float64(len(advertiserIds)) / oceanengine.PageSize100))

					for index := 0; index < totalPage; index++ {
						var advertiserInfos *account2.GetAdvertiserPublicInfoResponse

						if len(advertiserIds) < oceanengine.PageSize100*(index+1) {
							advertiserInfos, err = account2.GetAdvertiserPublicInfo(advertiserIds[index*oceanengine.PageSize100:len(advertiserIds)], accessToken.Data.AccessToken)

							if err != nil {
								inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, companyId, account.AdvertiserId, 0, 0, 0, oceanengineConfig.AppId, accessToken.Data.AccessToken, err.Error())
								inOceanengineApiLog.SetCreateTime(ctx)

								oauc.oalrepo.Save(ctx, inOceanengineApiLog)

								return err
							}
						} else {
							advertiserInfos, err = account2.GetAdvertiserPublicInfo(advertiserIds[index*oceanengine.PageSize100:oceanengine.PageSize100*(index+1)], accessToken.Data.AccessToken)

							if err != nil {
								inOceanengineApiLog := domain.NewOceanengineApiLog(ctx, companyId, account.AdvertiserId, 0, 0, 0, oceanengineConfig.AppId, accessToken.Data.AccessToken, err.Error())
								inOceanengineApiLog.SetCreateTime(ctx)

								oauc.oalrepo.Save(ctx, inOceanengineApiLog)

								return err
							}
						}

						for _, advertiserInfo := range advertiserInfos.Data {
							if inQianchuanAdvertiser, err := oauc.qarepo.GetByCompanyIdAndAdvertiserId(ctx, companyId, advertiserInfo.Id); err != nil {
								inQianchuanAdvertiser = domain.NewQianchuanAdvertiser(ctx, companyId, account.AdvertiserId, advertiserInfo.Id, oceanengineConfig.AppId, advertiserInfo.Name, advertiserInfo.Company)
								inQianchuanAdvertiser.SetCreateTime(ctx)
								inQianchuanAdvertiser.SetUpdateTime(ctx)

								if _, err := oauc.qarepo.Save(ctx, inQianchuanAdvertiser); err != nil {
									return err
								}

								inQianchuanAdvertiserStatus := domain.NewQianchuanAdvertiserStatus(ctx, companyId, advertiserInfo.Id, uint32(day), 0)
								inQianchuanAdvertiserStatus.SetCreateTime(ctx)
								inQianchuanAdvertiserStatus.SetUpdateTime(ctx)

								if _, err := oauc.qasrepo.Save(ctx, inQianchuanAdvertiserStatus); err != nil {
									return err
								}
							} else {
								inQianchuanAdvertiser.SetAppId(ctx, oceanengineConfig.AppId)
								inQianchuanAdvertiser.SetAccountId(ctx, account.AdvertiserId)
								inQianchuanAdvertiser.SetCompanyName(ctx, advertiserInfo.Company)
								inQianchuanAdvertiser.SetAdvertiserName(ctx, advertiserInfo.Name)
								inQianchuanAdvertiser.SetUpdateTime(ctx)

								if _, err := oauc.qarepo.Update(ctx, inQianchuanAdvertiser); err != nil {
									return err
								}
							}

							if inQianchuanAdvertiserHistory, err := oauc.qahrepo.GetByAdvertiserId(ctx, advertiserInfo.Id); err != nil {
								inQianchuanAdvertiserHistory = domain.NewQianchuanAdvertiserHistory(ctx, advertiserInfo.Id, uint32(day))
								inQianchuanAdvertiserHistory.SetCreateTime(ctx)
								inQianchuanAdvertiserHistory.SetUpdateTime(ctx)

								if _, err := oauc.qahrepo.Save(ctx, inQianchuanAdvertiserHistory); err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return DouyinOceanengineAccountCreateError
	}

	return nil
}
