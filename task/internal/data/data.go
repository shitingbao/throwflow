package data

import (
	"context"
	companyv1 "task/api/service/company/v1"
	douyinv1 "task/api/service/douyin/v1"
	weixinv1 "task/api/service/weixin/v1"
	"task/internal/conf"
	"time"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	grpcx "google.golang.org/grpc"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewOceanengineAccountTokenRepo, NewOpenDouyinTokenRepo, NewQianchuanAdvertiserRepo, NewQianchuanAdRepo, NewCompanyRepo, NewCompanyOrganizationRepo, NewJinritemaiOrderRepo, NewJinritemaiStoreRepo, NewOpenDouyinVideoRepo, NewWeixinUserRepo, NewWeixinUserOrganizationRepo, NewWeixinUserCommissionRepo, NewWeixinUserCouponRepo, NewWeixinUserBalanceRepo, NewDoukeTokenRepo, NewCompanyTaskRepo, NewDiscovery, NewDouyinServiceClient, NewCompanyServiceClient, NewWeixinServiceClient)

// Data .
type Data struct {
	douyinuc  douyinv1.DouyinClient
	companyuc companyv1.CompanyClient
	weixinuc  weixinv1.WeixinClient
}

// NewData .
func NewData(douyinuc douyinv1.DouyinClient, companyuc companyv1.CompanyClient, weixinuc weixinv1.WeixinClient, logger log.Logger) (*Data, func(), error) {
	data := &Data{douyinuc: douyinuc, companyuc: companyuc, weixinuc: weixinuc}

	cleanup := func() {
	}

	return data, cleanup, nil
}

func NewDouyinServiceClient(sr *conf.Service, rr registry.Discovery) douyinv1.DouyinClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Douyin.GetEndpoint()),
		grpc.WithDiscovery(rr),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithTimeout(20*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := douyinv1.NewDouyinClient(conn)
	return c
}

func NewCompanyServiceClient(sr *conf.Service, rr registry.Discovery) companyv1.CompanyClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Company.GetEndpoint()),
		grpc.WithDiscovery(rr),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithTimeout(20*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := companyv1.NewCompanyClient(conn)
	return c
}

func NewWeixinServiceClient(sr *conf.Service, rr registry.Discovery) weixinv1.WeixinClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Weixin.GetEndpoint()),
		grpc.WithDiscovery(rr),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithTimeout(20*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := weixinv1.NewWeixinClient(conn)
	return c
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	c := consulAPI.DefaultConfig()

	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}
