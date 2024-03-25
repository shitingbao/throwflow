package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type KuaidiCompanyRepo interface {
	List(context.Context, uint64, uint64, string) (*v1.ListKuaidiCompanysReply, error)
}

type KuaidiCompanyUsecase struct {
	repo KuaidiCompanyRepo
	conf *conf.Data
	log  *log.Helper
}

func NewKuaidiCompanyUsecase(repo KuaidiCompanyRepo, conf *conf.Data, logger log.Logger) *KuaidiCompanyUsecase {
	return &KuaidiCompanyUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (kcuc *KuaidiCompanyUsecase) ListKuaidiCompany(ctx context.Context, pageNum, pageSize uint64, keyword string) (*v1.ListKuaidiCompanysReply, error) {
	if pageSize == 0 {
		pageSize = uint64(kcuc.conf.Database.PageSize)
	}

	list, err := kcuc.repo.List(ctx, pageNum, pageSize, keyword)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_KUAIDI_COMPANY_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}
