package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	v1 "interface/api/interface/v1"
	"interface/internal/conf"
	"interface/internal/pkg/middleware/auth"
	"interface/internal/pkg/middleware/validate"
	"interface/internal/service"
	nhttp "net/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, sinterface *service.InterfaceService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			selector.Server(auth.Auth()).Match(NewWhiteListMatcher()).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	opts = append(opts, http.Filter(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
	)))

	opts = append(opts, http.ErrorEncoder(func(w nhttp.ResponseWriter, r *nhttp.Request, err error) {
		se := errors.FromError(err)

		if se.Reason == "CODEC" {
			se.Code = int32(500)
			se.Reason = "INTERFACE_VALIDATOR_ERROR"
			se.Message = "参数异常"
		}

		codec, _ := http.CodecForRequest(r, "Accept")
		body, err := codec.Marshal(se)

		if err != nil {
			w.WriteHeader(nhttp.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(nhttp.StatusOK)
		_, _ = w.Write(body)
	}))

	srv := http.NewServer(opts...)
	v1.RegisterInterfaceHTTPServer(srv, sinterface)
	return srv
}

// 白名单不需要token验证的接口
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/interface.v1.Interface/ListIndustries"] = struct{}{}
	whiteList["/interface.v1.Interface/ListAreas"] = struct{}{}
	whiteList["/interface.v1.Interface/GetToken"] = struct{}{}
	whiteList["/interface.v1.Interface/SendSms"] = struct{}{}
	whiteList["/interface.v1.Interface/ApplyForm"] = struct{}{}
	whiteList["/interface.v1.Interface/Login"] = struct{}{}
	whiteList["/interface.v1.Interface/ListSelectClues"] = struct{}{}
	whiteList["/interface.v1.Interface/CreateUsers"] = struct{}{}
	whiteList["/interface.v1.Interface/GetConfigUserOpenDouyins"] = struct{}{}
	whiteList["/interface.v1.Interface/GetTicketUserOpenDouyins"] = struct{}{}
	whiteList["/interface.v1.Interface/CreateUserOpenDouyins"] = struct{}{}
	whiteList["/interface.v1.Interface/GetFollows"] = struct{}{}
	whiteList["/interface.v1.Interface/ListMiniProducts"] = struct{}{}
	whiteList["/interface.v1.Interface/ListMiniCategorys"] = struct{}{}
	whiteList["/interface.v1.Interface/GetMiniProducts"] = struct{}{}
	whiteList["/interface.v1.Interface/StatisticsMiniProducts"] = struct{}{}
	whiteList["/interface.v1.Interface/ListMiniMaterials"] = struct{}{}
	whiteList["/interface.v1.Interface/StatisticsMiniMaterials"] = struct{}{}
	whiteList["/interface.v1.Interface/ListCompanyTaskUsable"] = struct{}{}
	whiteList["/interface.v1.Interface/ListMiniMaterialProducts"] = struct{}{}
	whiteList["/interface.v1.Interface/GetVideoUrlUserOrders"] = struct{}{}
	whiteList["/interface.v1.Interface/GetMiniMaterials"] = struct{}{}
	
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
