package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
)

type DoukeTokenRepo interface {
	Refresh(context.Context) (*v1.RefreshDoukeTokensReply, error)
}

type DoukeTokenUsecase struct {
	repo DoukeTokenRepo
	log  *log.Helper
}

func NewDoukeTokenUsecase(repo DoukeTokenRepo, logger log.Logger) *DoukeTokenUsecase {
	return &DoukeTokenUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (dtuc *DoukeTokenUsecase) RefreshDoukeTokens(ctx context.Context) (*v1.RefreshDoukeTokensReply, error) {
	token, err := dtuc.repo.Refresh(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_REFRESH_DOUKE_TOKEN_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return token, nil
}
