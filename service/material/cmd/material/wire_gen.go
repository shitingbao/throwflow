// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"material/internal/biz"
	"material/internal/conf"
	"material/internal/data"
	"material/internal/server"
	"material/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, company *conf.Company, volcengine *conf.Volcengine, confService *conf.Service, registry *conf.Registry, event *conf.Event, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData, logger)
	client := data.NewRedis(confData)
	sender := data.NewKafka(event)
	tos := data.NewTos(volcengine)
	discovery := data.NewDiscovery(registry)
	companyClient := data.NewCompanyServiceClient(confService, discovery)
	dataData, cleanup, err := data.NewData(confData, db, client, sender, tos, companyClient, logger)
	if err != nil {
		return nil, nil, err
	}
	materialRepo := data.NewMaterialRepo(dataData, logger)
	materialCategoryRepo := data.NewMaterialCategoryRepo(dataData, logger)
	materialProductRepo := data.NewMaterialProductRepo(dataData, logger)
	companyProductRepo := data.NewCompanyProductRepo(dataData, logger)
	companyMaterialRepo := data.NewCompanyMaterialRepo(dataData, logger)
	collectRepo := data.NewCollectRepo(dataData, logger)
	materialUsecase := biz.NewMaterialUsecase(materialRepo, materialCategoryRepo, materialProductRepo, companyProductRepo, companyMaterialRepo, collectRepo, confData, volcengine, company, logger)
	materialContentRepo := data.NewMaterialContentRepo(dataData, logger)
	materialContentUsecase := biz.NewMaterialContentUsecase(materialContentRepo, materialRepo, confData, logger)
	collectUsecase := biz.NewCollectUsecase(collectRepo, materialRepo, logger)
	materialService := service.NewMaterialService(materialUsecase, materialContentUsecase, collectUsecase)
	grpcServer := server.NewGRPCServer(confServer, materialService, logger)
	registrar := server.NewRegistrar(registry)
	app := newApp(logger, grpcServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}