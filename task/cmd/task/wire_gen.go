// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"task/internal/biz"
	"task/internal/conf"
	"task/internal/data"
	"task/internal/server"
	"task/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(task *conf.Task, registry *conf.Registry, confService *conf.Service, logger log.Logger) (*kratos.App, func(), error) {
	discovery := data.NewDiscovery(registry)
	douyinClient := data.NewDouyinServiceClient(confService, discovery)
	companyClient := data.NewCompanyServiceClient(confService, discovery)
	weixinClient := data.NewWeixinServiceClient(confService, discovery)
	dataData, cleanup, err := data.NewData(douyinClient, companyClient, weixinClient, logger)
	if err != nil {
		return nil, nil, err
	}
	oceanengineAccountTokenRepo := data.NewOceanengineAccountTokenRepo(dataData, logger)
	oceanengineAccountTokenUsecase := biz.NewOceanengineAccountTokenUsecase(oceanengineAccountTokenRepo, logger)
	openDouyinTokenRepo := data.NewOpenDouyinTokenRepo(dataData, logger)
	openDouyinTokenUsecase := biz.NewOpenDouyinTokenUsecase(openDouyinTokenRepo, logger)
	openDouyinVideoRepo := data.NewOpenDouyinVideoRepo(dataData, logger)
	openDouyinVideoUsecase := biz.NewOpenDouyinVideoUsecase(openDouyinVideoRepo, logger)
	qianchuanAdvertiserRepo := data.NewQianchuanAdvertiserRepo(dataData, logger)
	qianchuanAdvertiserUsecase := biz.NewQianchuanAdvertiserUsecase(qianchuanAdvertiserRepo, logger)
	qianchuanAdRepo := data.NewQianchuanAdRepo(dataData, logger)
	qianchuanAdUsecase := biz.NewQianchuanAdUsecase(qianchuanAdRepo, logger)
	companyRepo := data.NewCompanyRepo(dataData, logger)
	companyUsecase := biz.NewCompanyUsecase(companyRepo, logger)
	companyOrganizationRepo := data.NewCompanyOrganizationRepo(dataData, logger)
	companyOrganizationUsecase := biz.NewCompanyOrganizationUsecase(companyOrganizationRepo, logger)
	companyTaskRepo := data.NewCompanyTaskRepo(dataData, logger)
	companyTaskUsecase := biz.NewCompanyTaskUsecase(companyTaskRepo, logger)
	companyProductRepo := data.NewCompanyProductRepo(dataData, logger)
	companyProductUsecase := biz.NewCompanyProductUsecase(companyProductRepo, logger)
	jinritemaiStoreRepo := data.NewJinritemaiStoreRepo(dataData, logger)
	jinritemaiStoreUsecase := biz.NewJinritemaiStoreUsecase(jinritemaiStoreRepo, logger)
	jinritemaiOrderRepo := data.NewJinritemaiOrderRepo(dataData, logger)
	jinritemaiOrderUsecase := biz.NewJinritemaiOrderUsecase(jinritemaiOrderRepo, logger)
	weixinUserRepo := data.NewWeixinUserRepo(dataData, logger)
	weixinUserUsecase := biz.NewWeixinUserUsecase(weixinUserRepo, logger)
	weixinUserCommissionRepo := data.NewWeixinUserCommissionRepo(dataData, logger)
	weixinUserCommissionUsecase := biz.NewWeixinUserCommissionUsecase(weixinUserCommissionRepo, logger)
	weixinUserCouponRepo := data.NewWeixinUserCouponRepo(dataData, logger)
	weixinUserCouponUsecase := biz.NewWeixinUserCouponUsecase(weixinUserCouponRepo, logger)
	weixinUserBalanceRepo := data.NewWeixinUserBalanceRepo(dataData, logger)
	weixinUserBalanceUsecase := biz.NewWeixinUserBalanceUsecase(weixinUserBalanceRepo, logger)
	doukeOrderRepo := data.NewDoukeOrderRepo(dataData, logger)
	doukeOrderUsecase := biz.NewDoukeOrderUsecase(doukeOrderRepo, logger)
	taskService := service.NewTaskService(logger, oceanengineAccountTokenUsecase, openDouyinTokenUsecase, openDouyinVideoUsecase, qianchuanAdvertiserUsecase, qianchuanAdUsecase, companyUsecase, companyOrganizationUsecase, companyTaskUsecase, companyProductUsecase, jinritemaiStoreUsecase, jinritemaiOrderUsecase, weixinUserUsecase, weixinUserCommissionUsecase, weixinUserCouponUsecase, weixinUserBalanceUsecase, doukeOrderUsecase)
	serverTask := server.NewTaskServer(task, taskService)
	app := newApp(logger, serverTask)
	return app, func() {
		cleanup()
	}, nil
}
