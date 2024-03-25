package auth

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"interface/internal/biz"
)

func Auth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				token := tr.RequestHeader().Get("token")

				if len(token) == 0 {
					return nil, biz.InterfaceLoginTokenError
				}

				ctx = context.WithValue(ctx, "token", token)
			}

			return handler(ctx, req)
		}
	}
}
