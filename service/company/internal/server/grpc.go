package server

import (
	v1 "company/api/company/v1"
	"company/internal/conf"
	"company/internal/pkg/middleware/validate"
	"company/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, company *service.CompanyService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			validate.Validator(),
		),
		grpc.Options(ggrpc.MaxRecvMsgSize(8 * 1024 * 1024)),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterCompanyServer(srv, company)
	return srv
}
