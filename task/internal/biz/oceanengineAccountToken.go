package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
)

type OceanengineAccountTokenRepo interface {
	Refresh(context.Context) (*v1.RefreshOceanengineAccountTokensReply, error)
}

type OceanengineAccountTokenUsecase struct {
	repo OceanengineAccountTokenRepo
	log  *log.Helper
}

func NewOceanengineAccountTokenUsecase(repo OceanengineAccountTokenRepo, logger log.Logger) *OceanengineAccountTokenUsecase {
	return &OceanengineAccountTokenUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (oatuc *OceanengineAccountTokenUsecase) RefreshOceanengineAccountTokens(ctx context.Context) (*v1.RefreshOceanengineAccountTokensReply, error) {
	token, err := oatuc.repo.Refresh(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_REFRESH_OCEANENGINE_ACCOUNT_TOKEN_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return token, nil
}
