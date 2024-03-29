package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/douke/token"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

var (
	DouyinDoukeTokenNotFound    = errors.InternalServer("DOUYIN_DOUKE_TOKEN_NOT_FOUND", "抖客授权账户Token不存在")
	DouyinDoukeTokenCreateError = errors.InternalServer("DOUYIN_DOUKE_TOKEN_CREATE_ERROR", "抖客授权账户信息创建失败")
)

type DoukeTokenRepo interface {
	Get(context.Context, string, string) (*domain.DoukeToken, error)
	List(context.Context) ([]*domain.DoukeToken, error)
	Save(context.Context, *domain.DoukeToken) (*domain.DoukeToken, error)
	Update(context.Context, *domain.DoukeToken) (*domain.DoukeToken, error)
}

type DoukeTokenUsecase struct {
	repo   DoukeTokenRepo
	tlrepo TaskLogRepo
	tm     Transaction
	conf   *conf.Data
	dconf  *conf.Douke
	log    *log.Helper
}

func NewDoukeTokenUsecase(repo DoukeTokenRepo, tlrepo TaskLogRepo, tm Transaction, conf *conf.Data, dconf *conf.Douke, logger log.Logger) *DoukeTokenUsecase {
	return &DoukeTokenUsecase{repo: repo, tlrepo: tlrepo, tm: tm, conf: conf, dconf: dconf, log: log.NewHelper(logger)}
}

func (dtuc *DoukeTokenUsecase) CreateDoukeTokens(ctx context.Context, code string) error {
	token, err := token.CreateToken(dtuc.dconf.AppKey, dtuc.dconf.AppSecret, code)

	if err != nil {
		return DouyinDoukeTokenCreateError
	}

	if inDoukeToken, err := dtuc.repo.Get(ctx, token.Data.AuthorityId, token.Data.AuthSubjectType); err != nil {
		inDoukeToken = domain.NewDoukeToken(ctx, token.Data.AuthorityId, token.Data.AuthSubjectType, token.Data.AccessToken, token.Data.RefreshToken, token.Data.ExpiresIn)
		inDoukeToken.SetCreateTime(ctx)
		inDoukeToken.SetUpdateTime(ctx)

		if _, err := dtuc.repo.Save(ctx, inDoukeToken); err != nil {
			return DouyinDoukeTokenCreateError
		}
	} else {
		inDoukeToken.SetAccessToken(ctx, token.Data.AccessToken)
		inDoukeToken.SetExpiresIn(ctx, token.Data.ExpiresIn)
		inDoukeToken.SetRefreshToken(ctx, token.Data.RefreshToken)
		inDoukeToken.SetUpdateTime(ctx)

		if _, err := dtuc.repo.Update(ctx, inDoukeToken); err != nil {
			return DouyinDoukeTokenCreateError
		}
	}

	return nil
}

func (dtuc *DoukeTokenUsecase) RefreshDoukeTokens(ctx context.Context) error {
	list, err := dtuc.repo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Douke RefreshTokenError] Description=%s", "抖客账户token表获取失败"))
		inTaskLog.SetCreateTime(ctx)

		dtuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, inDoukeToken := range list {
		if time.Now().After(inDoukeToken.UpdateTime.Add((time.Duration(inDoukeToken.ExpiresIn) * time.Second) - 60*60*time.Second)) {
			if token, err := token.RefreshToken(dtuc.dconf.AppKey, dtuc.dconf.AppSecret, inDoukeToken.RefreshToken); err == nil {
				inDoukeToken.SetAccessToken(ctx, token.Data.AccessToken)
				inDoukeToken.SetExpiresIn(ctx, token.Data.ExpiresIn)
				inDoukeToken.SetRefreshToken(ctx, token.Data.RefreshToken)
				inDoukeToken.SetUpdateTime(ctx)

				if _, err := dtuc.repo.Update(ctx, inDoukeToken); err != nil {
					inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Douke RefreshTokenError] AuthorityId=%s, AuthSubjectType=%s, AccessToken=%s, RefreshToken=%s, ExpiresIn=%d, Description=%s", inDoukeToken.AuthorityId, inDoukeToken.AuthSubjectType, inDoukeToken.AccessToken, inDoukeToken.RefreshToken, inDoukeToken.ExpiresIn, "新的token插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					dtuc.tlrepo.Save(ctx, inTaskLog)
				}
			} else {
				fmt.Println(err)
				inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Douke RefreshTokenError] AuthorityId=%s, AuthSubjectType=%s, Description=%s", inDoukeToken.AuthorityId, inDoukeToken.AuthSubjectType, "调用抖客平台刷新token接口失败"))
				inTaskLog.SetCreateTime(ctx)

				dtuc.tlrepo.Save(ctx, inTaskLog)
			}
		}
	}

	return nil
}
