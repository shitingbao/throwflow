package service

import (
	v1 "douyin/api/douyin/v1"
	"douyin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewDouyinService)

type DouyinService struct {
	v1.UnimplementedDouyinServer

	ocuc   *biz.OceanengineConfigUsecase
	oauc   *biz.OceanengineAccountUsecase
	qauc   *biz.QianchuanAdvertiserUsecase
	qcuc   *biz.QianchuanCampaignUsecase
	oatuc  *biz.OceanengineAccountTokenUsecase
	qrpuc  *biz.QianchuanReportProductUsecase
	qrauc  *biz.QianchuanReportAwemeUsecase
	qaduc  *biz.QianchuanAdUsecase
	qahuc  *biz.QianchuanAdvertiserHistoryUsecase
	odtuc  *biz.OpenDouyinTokenUsecase
	jouc   *biz.JinritemaiOrderUsecase
	jsuc   *biz.JinritemaiStoreUsecase
	odvuc  *biz.OpenDouyinVideoUsecase
	oduiuc *biz.OpenDouyinUserInfoUsecase
	dpuc   *biz.DoukeProductUsecase
	douc   *biz.DoukeOrderUsecase

	log *log.Helper
}

func NewDouyinService(ocuc *biz.OceanengineConfigUsecase, oauc *biz.OceanengineAccountUsecase, qauc *biz.QianchuanAdvertiserUsecase, qcuc *biz.QianchuanCampaignUsecase, oatuc *biz.OceanengineAccountTokenUsecase, qrpuc *biz.QianchuanReportProductUsecase, qrauc *biz.QianchuanReportAwemeUsecase, qaduc *biz.QianchuanAdUsecase, qahuc *biz.QianchuanAdvertiserHistoryUsecase, odtuc *biz.OpenDouyinTokenUsecase, jouc *biz.JinritemaiOrderUsecase, jsuc *biz.JinritemaiStoreUsecase, odvuc *biz.OpenDouyinVideoUsecase, oduiuc *biz.OpenDouyinUserInfoUsecase, dpuc *biz.DoukeProductUsecase, douc *biz.DoukeOrderUsecase, logger log.Logger) *DouyinService {
	log := log.NewHelper(log.With(logger, "module", "service/douyin"))

	return &DouyinService{ocuc: ocuc, oauc: oauc, qauc: qauc, qcuc: qcuc, oatuc: oatuc, qrpuc: qrpuc, qrauc: qrauc, qaduc: qaduc, qahuc: qahuc, odtuc: odtuc, jouc: jouc, jsuc: jsuc, odvuc: odvuc, oduiuc: oduiuc, dpuc: dpuc, douc: douc, log: log}
}
