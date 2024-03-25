package data

import (
	"context"
	commonv1 "interface/api/service/common/v1"
	companyv1 "interface/api/service/company/v1"
	douyinv1 "interface/api/service/douyin/v1"
	materialv1 "interface/api/service/material/v1"
	weixinv1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
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
var ProviderSet = wire.NewSet(NewData, NewDiscovery, NewIndustryRepo, NewAreaRepo, NewSmsRepo, NewTokenRepo, NewKuaidiCompanyRepo, NewLoginRepo, NewCompanyUserRepo, NewCompanySetRepo, NewQianchuanAdvertiserRepo, NewClueRepo, NewOceanengineConfigRepo, NewMaterialRepo, NewPerformanceRuleRepo, NewPerformanceRepo, NewPerformanceRebalanceRepo, NewUpdateLogRepo, NewUserRepo, NewProductRepo, NewCompanyRepo, NewCompanyMaterialRepo, NewCompanyOrganizationRepo, NewUserAddressRepo, NewUserStoreRepo, NewUserOpenDouyinRepo, NewUserSampleOrderRepo, NewUserScanRecordRepo, NewJinritemailOrderRepo, NewDoukeOrderRepo, NewCommonServiceClient, NewDouyinServiceClient, NewUserOrderRepo, NewUserOrganizationRepo, NewUserCouponRepo, NewUserBalanceRepo, NewUserContractRepo, NewUserBankRepo, NewCompanyTaskRepo, NewCompanyServiceClient, NewMaterialServiceClient, NewWeixinServiceClient)

// Data .
type Data struct {
	commonuc   commonv1.CommonClient
	douyinuc   douyinv1.DouyinClient
	companyuc  companyv1.CompanyClient
	materialuc materialv1.MaterialClient
	weixinuc   weixinv1.WeixinClient
}

// NewData
func NewData(commonuc commonv1.CommonClient, douyinuc douyinv1.DouyinClient, companyuc companyv1.CompanyClient, materialuc materialv1.MaterialClient, weixinuc weixinv1.WeixinClient, logger log.Logger) (*Data, func(), error) {
	data := &Data{commonuc: commonuc, douyinuc: douyinuc, companyuc: companyuc, materialuc: materialuc, weixinuc: weixinuc}

	cleanup := func() {
	}

	return data, cleanup, nil
}

func NewCommonServiceClient(sr *conf.Service, rr registry.Discovery) commonv1.CommonClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Common.GetEndpoint()),
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
	c := commonv1.NewCommonClient(conn)
	return c
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
		grpc.WithTimeout(50*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := companyv1.NewCompanyClient(conn)
	return c
}

func NewMaterialServiceClient(sr *conf.Service, rr registry.Discovery) materialv1.MaterialClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Material.GetEndpoint()),
		grpc.WithDiscovery(rr),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithTimeout(10*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := materialv1.NewMaterialClient(conn)
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
		grpc.WithTimeout(10*time.Second),
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
