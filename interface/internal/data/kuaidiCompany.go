package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
	"interface/internal/biz"
)

type kuaidiCompanyRepo struct {
	data *Data
	log  *log.Helper
}

func NewKuaidiCompanyRepo(data *Data, logger log.Logger) biz.KuaidiCompanyRepo {
	return &kuaidiCompanyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (kcr *kuaidiCompanyRepo) List(ctx context.Context, pageNum, pageSize uint64, keyword string) (*v1.ListKuaidiCompanysReply, error) {
	list, err := kcr.data.commonuc.ListKuaidiCompanys(ctx, &v1.ListKuaidiCompanysRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		Keyword:  keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
