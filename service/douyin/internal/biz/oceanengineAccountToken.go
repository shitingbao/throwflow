package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/oceanengine/oauth2"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

var (
	DouyinOceanengineAccountTokenNotFound     = errors.InternalServer("DOUYIN_OCEANENGINE_ACCOUNT_TOKEN_NOT_FOUND", "千川授权账户Token不存在")
	DouyinOceanengineAccountTokenListError    = errors.InternalServer("DOUYIN_OCEANENGINE_ACCOUNT_TOKEN_LIST_ERROR", "千川授权账户Token获取失败")
	DouyinOceanengineAccountTokenSyncError    = errors.InternalServer("DOUYIN_OCEANENGINE_ACCOUNT_TOKEN_SYNC_ERROR", "同步千川授权账户Token失败")
	DouyinOceanengineAccountTokenRefreshError = errors.InternalServer("DOUYIN_OCEANENGINE_ACCOUNT_TOKEN_REFRESH_ERROR", "刷新千川授权账户Token失败")
)

type OceanengineAccountTokenRepo interface {
	GetByCompanyIdAndAccountId(context.Context, uint64, uint64) (*domain.OceanengineAccountToken, error)
	List(context.Context) ([]*domain.OceanengineAccountToken, error)
	ListByAppId(context.Context, string) ([]*domain.OceanengineAccountToken, error)
	ListByCompanyId(context.Context, uint64) ([]*domain.OceanengineAccountToken, error)
	Update(context.Context, *domain.OceanengineAccountToken) (*domain.OceanengineAccountToken, error)
	Save(context.Context, *domain.OceanengineAccountToken) (*domain.OceanengineAccountToken, error)
	DeleteByCompanyIdAndAccountId(context.Context, uint64, uint64) error
}

type OceanengineAccountTokenUsecase struct {
	repo   OceanengineAccountTokenRepo
	ocrepo OceanengineConfigRepo
	oarepo OceanengineAccountRepo
	qarepo QianchuanAdvertiserRepo
	tlrepo TaskLogRepo
	conf   *conf.Data
	log    *log.Helper
}

func NewOceanengineAccountTokenUsecase(repo OceanengineAccountTokenRepo, ocrepo OceanengineConfigRepo, oarepo OceanengineAccountRepo, qarepo QianchuanAdvertiserRepo, tlrepo TaskLogRepo, conf *conf.Data, logger log.Logger) *OceanengineAccountTokenUsecase {
	return &OceanengineAccountTokenUsecase{repo: repo, ocrepo: ocrepo, oarepo: oarepo, qarepo: qarepo, tlrepo: tlrepo, conf: conf, log: log.NewHelper(logger)}
}

func (oatuc *OceanengineAccountTokenUsecase) GetOceanengineAccountTokens(ctx context.Context, companyId, advertiserId uint64) (*domain.OceanengineAccountToken, error) {
	qianchuanAdvertiser, err := oatuc.qarepo.GetByCompanyIdAndAdvertiserId(ctx, companyId, advertiserId)

	if err != nil {
		return nil, DouyinOceanengineAccountTokenNotFound
	}

	if qianchuanAdvertiser.Status != 1 {
		return nil, DouyinOceanengineAccountTokenNotFound
	}

	oceanengineAccountToken, err := oatuc.repo.GetByCompanyIdAndAccountId(ctx, qianchuanAdvertiser.CompanyId, qianchuanAdvertiser.AccountId)

	if err != nil {
		return nil, DouyinOceanengineAccountTokenNotFound
	}

	return oceanengineAccountToken, nil
}

func (oatuc *OceanengineAccountTokenUsecase) RefreshOceanengineAccountTokens(ctx context.Context) error {
	list, err := oatuc.repo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Oceanengine RefreshTokenError] Description=%s", "获取账户token表失败"))
		inTaskLog.SetCreateTime(ctx)

		oatuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, inOceanengineAccountToken := range list {
		if time.Now().After(inOceanengineAccountToken.UpdateTime.Add((time.Duration(inOceanengineAccountToken.ExpiresIn) * time.Second) - 60*60*time.Second)) {
			if oceanengineConfig, err := oatuc.ocrepo.GetByAppId(ctx, inOceanengineAccountToken.AppId); err == nil {
				if accessToken, err := oauth2.RefreshToken(inOceanengineAccountToken.AppId, oceanengineConfig.AppSecret, inOceanengineAccountToken.RefreshToken); err == nil {
					inOceanengineAccountToken.SetAccessToken(ctx, accessToken.Data.AccessToken)
					inOceanengineAccountToken.SetRefreshToken(ctx, accessToken.Data.RefreshToken)
					inOceanengineAccountToken.SetExpiresIn(ctx, accessToken.Data.ExpiresIn)
					inOceanengineAccountToken.SetRefreshTokenExpiresIn(ctx, accessToken.Data.RefreshTokenExpiresIn)
					inOceanengineAccountToken.SetUpdateTime(ctx)

					if _, err := oatuc.repo.Update(ctx, inOceanengineAccountToken); err != nil {
						inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Oceanengine RefreshTokenError] CompanyId=%d, AccountId=%d, AppId=%s, AccessToken=%s, RefreshToken=%s, ExpiresIn=%d, RefreshTokenExpiresIn=%d, Description=%s", inOceanengineAccountToken.CompanyId, inOceanengineAccountToken.AccountId, inOceanengineAccountToken.AppId, accessToken.Data.AccessToken, accessToken.Data.RefreshToken, accessToken.Data.ExpiresIn, accessToken.Data.RefreshTokenExpiresIn, "获取新的token插入数据库失败"))
						inTaskLog.SetCreateTime(ctx)

						oatuc.tlrepo.Save(ctx, inTaskLog)
					}
				} else {
					inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Oceanengine RefreshTokenError] CompanyId=%d, AccountId=%d, AppId=%s, Description=%s", inOceanengineAccountToken.CompanyId, inOceanengineAccountToken.AccountId, inOceanengineAccountToken.AppId, "调用巨量引擎刷新token接口失败"))
					inTaskLog.SetCreateTime(ctx)

					oatuc.tlrepo.Save(ctx, inTaskLog)
				}
			} else {
				inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Oceanengine RefreshTokenError] CompanyId=%d, AccountId=%d, AppId=%s, Description=%s", inOceanengineAccountToken.CompanyId, inOceanengineAccountToken.AccountId, inOceanengineAccountToken.AppId, "获取巨量引擎配置文件失败"))
				inTaskLog.SetCreateTime(ctx)

				oatuc.tlrepo.Save(ctx, inTaskLog)
			}
		}
	}

	return nil
}
