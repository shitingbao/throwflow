package biz

import (
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	CommonKuaidiCompanyListError = errors.InternalServer("COMMON_KUAIDI_COMPANY_LIST_ERROR", "快递公司编码列表获取失败")
)

type KuaidiCompanyRepo interface {
	Get(context.Context, string) (*domain.KuaidiCompany, error)
	List(context.Context, int, int, string) ([]*domain.KuaidiCompany, error)
	Count(context.Context, string) (int64, error)
}

type KuaidiCompanyUsecase struct {
	repo KuaidiCompanyRepo
	log  *log.Helper
}

func NewKuaidiCompanyUsecase(repo KuaidiCompanyRepo, logger log.Logger) *KuaidiCompanyUsecase {
	return &KuaidiCompanyUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (kcuc *KuaidiCompanyUsecase) ListKuaidiCompanys(ctx context.Context, pageNum, pageSize uint64, keyword string) (*domain.KuaidiCompanyList, error) {
	list, err := kcuc.repo.List(ctx, int(pageNum), int(pageSize), keyword)

	if err != nil {
		return nil, CommonKuaidiCompanyListError
	}

	total, err := kcuc.repo.Count(ctx, keyword)

	if err != nil {
		return nil, CommonKuaidiCompanyListError
	}

	return &domain.KuaidiCompanyList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}
