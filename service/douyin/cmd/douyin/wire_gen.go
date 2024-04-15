// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"douyin/internal/biz"
	"douyin/internal/conf"
	"douyin/internal/data"
	"douyin/internal/server"
	"douyin/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, oceanengine *conf.Oceanengine, event *conf.Event, registry *conf.Registry, confService *conf.Service, developer *conf.Developer, csj *conf.Csj, company *conf.Company, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData, logger)
	v := data.NewRdsDB(confData, logger)
	client := data.NewMongo(confData)
	conn := data.NewClickhouse(confData)
	redisClient := data.NewRedis(confData)
	sender := data.NewKafka(event)
	discovery := data.NewDiscovery(registry)
	companyClient := data.NewCompanyServiceClient(confService, discovery)
	weixinClient := data.NewWeixinServiceClient(confService, discovery)
	dataData, cleanup, err := data.NewData(confData, db, v, client, conn, redisClient, sender, companyClient, weixinClient, logger)
	if err != nil {
		return nil, nil, err
	}
	oceanengineConfigRepo := data.NewOceanengineConfigRepo(dataData, logger)
	oceanengineConfigUsecase := biz.NewOceanengineConfigUsecase(oceanengineConfigRepo, confData, logger)
	oceanengineAccountRepo := data.NewOceanengineAccountRepo(dataData, logger)
	oceanengineAccountTokenRepo := data.NewOceanengineAccountTokenRepo(dataData, logger)
	qianchuanAdvertiserRepo := data.NewQianchuanAdvertiserRepo(dataData, logger)
	qianchuanAdvertiserStatusRepo := data.NewQianchuanAdvertiserStatusRepo(dataData, logger)
	qianchuanAdvertiserHistoryRepo := data.NewQianchuanAdvertiserHistoryRepo(dataData, logger)
	oceanengineApiLogRepo := data.NewOceanengineApiLogRepo(dataData, logger)
	companyRepo := data.NewCompanyRepo(dataData, logger)
	transaction := data.NewTransaction(dataData)
	oceanengineAccountUsecase := biz.NewOceanengineAccountUsecase(oceanengineAccountRepo, oceanengineConfigRepo, oceanengineAccountTokenRepo, qianchuanAdvertiserRepo, qianchuanAdvertiserStatusRepo, qianchuanAdvertiserHistoryRepo, oceanengineApiLogRepo, companyRepo, transaction, confData, logger)
	qianchuanAdvertiserInfoRepo := data.NewQianchuanAdvertiserInfoRepo(dataData, logger)
	companyUserRepo := data.NewCompanyUserRepo(dataData, logger)
	taskLogRepo := data.NewTaskLogRepo(dataData, logger)
	qianchuanCampaignRepo := data.NewQianchuanCampaignRepo(dataData, logger)
	qianchuanProductRepo := data.NewQianchuanProductRepo(dataData, logger)
	qianchuanAwemeRepo := data.NewQianchuanAwemeRepo(dataData, logger)
	qianchuanWalletRepo := data.NewQianchuanWalletRepo(dataData, logger)
	qianchuanAdRepo := data.NewQianchuanAdRepo(dataData, logger)
	qianchuanAdInfoRepo := data.NewQianchuanAdInfoRepo(dataData, logger)
	qianchuanReportAdRepo := data.NewQianchuanReportAdRepo(dataData, logger)
	qianchuanReportAdRealtimeRepo := data.NewQianchuanReportAdRealtimeRepo(dataData, logger)
	qianchuanReportProductRepo := data.NewQianchuanReportProductRepo(dataData, logger)
	qianchuanReportAwemeRepo := data.NewQianchuanReportAwemeRepo(dataData, logger)
	lianshanRealtimeRepo := data.NewLianshanRealtimeRepo(dataData, logger)
	qianchuanAdvertiserUsecase := biz.NewQianchuanAdvertiserUsecase(qianchuanAdvertiserRepo, qianchuanAdvertiserInfoRepo, qianchuanAdvertiserStatusRepo, oceanengineConfigRepo, oceanengineAccountTokenRepo, companyRepo, companyUserRepo, taskLogRepo, qianchuanCampaignRepo, qianchuanProductRepo, qianchuanAwemeRepo, qianchuanWalletRepo, qianchuanAdRepo, qianchuanAdInfoRepo, qianchuanReportAdRepo, qianchuanReportAdRealtimeRepo, qianchuanReportProductRepo, qianchuanReportAwemeRepo, lianshanRealtimeRepo, oceanengineApiLogRepo, transaction, confData, event, oceanengine, logger)
	companySetRepo := data.NewCompanySetRepo(dataData, logger)
	qianchuanCampaignUsecase := biz.NewQianchuanCampaignUsecase(qianchuanCampaignRepo, qianchuanReportAdRealtimeRepo, companySetRepo, qianchuanAdRepo, qianchuanReportAdRepo, qianchuanAdvertiserRepo, confData, logger)
	oceanengineAccountTokenUsecase := biz.NewOceanengineAccountTokenUsecase(oceanengineAccountTokenRepo, oceanengineConfigRepo, oceanengineAccountRepo, qianchuanAdvertiserRepo, taskLogRepo, confData, logger)
	qianchuanReportProductUsecase := biz.NewQianchuanReportProductUsecase(qianchuanReportProductRepo, confData, logger)
	qianchuanReportAwemeUsecase := biz.NewQianchuanReportAwemeUsecase(qianchuanReportAwemeRepo, confData, logger)
	qianchuanAdUsecase := biz.NewQianchuanAdUsecase(qianchuanAdRepo, qianchuanReportAdRepo, qianchuanReportAdRealtimeRepo, companyRepo, companySetRepo, qianchuanAdInfoRepo, qianchuanAdvertiserRepo, qianchuanAdvertiserStatusRepo, qianchuanAdvertiserInfoRepo, qianchuanCampaignRepo, qianchuanWalletRepo, taskLogRepo, confData, event, logger)
	qianchuanAdvertiserHistoryUsecase := biz.NewQianchuanAdvertiserHistoryUsecase(qianchuanAdvertiserHistoryRepo, confData, logger)
	openDouyinTokenRepo := data.NewOpenDouyinTokenRepo(dataData, logger)
	openDouyinUserInfoRepo := data.NewOpenDouyinUserInfoRepo(dataData, logger)
	openDouyinUserInfoCreateLogRepo := data.NewOpenDouyinUserInfoCreateLogRepo(dataData, logger)
	openDouyinApiLogRepo := data.NewOpenDouyinApiLogRepo(dataData, logger)
	jinritemaiApiLogRepo := data.NewJinritemaiApiLogRepo(dataData, logger)
	weixinUserRepo := data.NewWeixinUserRepo(dataData, logger)
	weixinUserOpenDouyinRepo := data.NewWeixinUserOpenDouyinRepo(dataData, logger)
	jinritemaiOrderInfoRepo := data.NewJinritemaiOrderInfoRepo(dataData, logger)
	openDouyinTokenUsecase := biz.NewOpenDouyinTokenUsecase(openDouyinTokenRepo, openDouyinUserInfoRepo, openDouyinUserInfoCreateLogRepo, openDouyinApiLogRepo, jinritemaiApiLogRepo, weixinUserRepo, weixinUserOpenDouyinRepo, jinritemaiOrderInfoRepo, taskLogRepo, transaction, confData, developer, event, logger)
	jinritemaiOrderRepo := data.NewJinritemaiOrderRepo(dataData, logger)
	doukeOrderInfoRepo := data.NewDoukeOrderInfoRepo(dataData, logger)
	weixinUserCommissionRepo := data.NewWeixinUserCommissionRepo(dataData, logger)
	openDouyinVideoRepo := data.NewOpenDouyinVideoRepo(dataData, logger)
	companyProductRepo := data.NewCompanyProductRepo(dataData, logger)
	qianchuanAwemeOrderInfoRepo := data.NewQianchuanAwemeOrderInfoRepo(dataData, logger)
	jinritemaiOrderUsecase := biz.NewJinritemaiOrderUsecase(jinritemaiOrderRepo, jinritemaiOrderInfoRepo, doukeOrderInfoRepo, weixinUserRepo, weixinUserOpenDouyinRepo, weixinUserCommissionRepo, openDouyinTokenRepo, openDouyinUserInfoRepo, openDouyinUserInfoCreateLogRepo, openDouyinVideoRepo, taskLogRepo, jinritemaiApiLogRepo, companyProductRepo, qianchuanAdvertiserStatusRepo, qianchuanAwemeOrderInfoRepo, confData, event, developer, logger)
	jinritemaiStoreRepo := data.NewJinritemaiStoreRepo(dataData, logger)
	jinritemaiStoreInfoRepo := data.NewJinritemaiStoreInfoRepo(dataData, logger)
	jinritemaiStoreUsecase := biz.NewJinritemaiStoreUsecase(jinritemaiStoreRepo, jinritemaiStoreInfoRepo, weixinUserRepo, weixinUserOpenDouyinRepo, openDouyinTokenRepo, companyProductRepo, taskLogRepo, jinritemaiApiLogRepo, confData, company, developer, logger)
	openDouyinVideoUsecase := biz.NewOpenDouyinVideoUsecase(openDouyinVideoRepo, openDouyinTokenRepo, openDouyinUserInfoRepo, taskLogRepo, openDouyinApiLogRepo, jinritemaiOrderInfoRepo, weixinUserOpenDouyinRepo, confData, logger)
	openDouyinUserInfoUsecase := biz.NewOpenDouyinUserInfoUsecase(openDouyinUserInfoRepo, openDouyinVideoRepo, jinritemaiStoreInfoRepo, jinritemaiOrderInfoRepo, weixinUserOpenDouyinRepo, transaction, confData, logger)
	doukeProductUsecase := biz.NewDoukeProductUsecase(transaction, confData, csj, logger)
	doukeOrderRepo := data.NewDoukeOrderRepo(dataData, logger)
	csjApiLogRepo := data.NewCsjApiLogRepo(dataData, logger)
	doukeOrderUsecase := biz.NewDoukeOrderUsecase(doukeOrderRepo, doukeOrderInfoRepo, companyProductRepo, csjApiLogRepo, taskLogRepo, weixinUserCommissionRepo, confData, csj, logger)
	douyinService := service.NewDouyinService(oceanengineConfigUsecase, oceanengineAccountUsecase, qianchuanAdvertiserUsecase, qianchuanCampaignUsecase, oceanengineAccountTokenUsecase, qianchuanReportProductUsecase, qianchuanReportAwemeUsecase, qianchuanAdUsecase, qianchuanAdvertiserHistoryUsecase, openDouyinTokenUsecase, jinritemaiOrderUsecase, jinritemaiStoreUsecase, openDouyinVideoUsecase, openDouyinUserInfoUsecase, doukeProductUsecase, doukeOrderUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, douyinService, logger)
	registrar := server.NewRegistrar(registry)
	app := newApp(logger, grpcServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
