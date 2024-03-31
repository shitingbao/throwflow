package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/common/v1"
	"weixin/internal/biz"
)

type kuaidiInfoRepo struct {
	data *Data
	log  *log.Helper
}

func NewKuaidiInfoRepo(data *Data, logger log.Logger) biz.KuaidiInfoRepo {
	return &kuaidiInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (kir *kuaidiInfoRepo) Get(ctx context.Context, code, num, phone string) (*v1.GetKuaidiInfosReply, error) {
	kuaidiInfo, err := kir.data.commonuc.GetKuaidiInfos(ctx, &v1.GetKuaidiInfosRequest{
		Code:  code,
		Num:   num,
		Phone: phone,
	})

	if err != nil {
		return nil, err
	}

	return kuaidiInfo, err
}
