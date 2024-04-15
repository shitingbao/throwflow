//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"douyin/internal/biz"
	"douyin/internal/conf"
	"douyin/internal/data"
	"douyin/internal/server"
	"douyin/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Oceanengine, *conf.Event, *conf.Registry, *conf.Service, *conf.Developer, *conf.Csj, *conf.Company, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
