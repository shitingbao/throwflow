package service

import (
	v1 "common/api/common/v1"
	"common/internal/biz"
	"common/internal/conf"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewCommonService)

type CommonService struct {
	v1.UnimplementedCommonServer

	suc  *biz.SmsUsecase
	puc  *biz.PayUsecase
	tuc  *biz.TokenUsecase
	auc  *biz.AreaUsecase
	uluc *biz.UpdateLogUsecase
	suuc *biz.ShortUrlUsecase
	scuc *biz.ShortCodeUsecase
	kcuc *biz.KuaidiCompanyUsecase
	kiuc *biz.KuaidiInfoUsecase

	conf *conf.Data
}

func NewCommonService(suc *biz.SmsUsecase, puc *biz.PayUsecase, tuc *biz.TokenUsecase, auc *biz.AreaUsecase, uluc *biz.UpdateLogUsecase, suuc *biz.ShortUrlUsecase, scuc *biz.ShortCodeUsecase, kcuc *biz.KuaidiCompanyUsecase, kiuc *biz.KuaidiInfoUsecase, conf *conf.Data) *CommonService {
	return &CommonService{suc: suc, puc: puc, tuc: tuc, auc: auc, uluc: uluc, suuc: suuc, scuc: scuc, kcuc: kcuc, kiuc: kiuc, conf: conf}
}
