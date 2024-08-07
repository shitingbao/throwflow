package validate

import (
	"company/internal/biz"
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-kratos/kratos/v2/middleware"
)

type validator interface {
	Validate() error
}

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(validator); ok {
				spew.Dump(v.Validate())
				if err := v.Validate(); err != nil {
					return nil, biz.CompanyValidatorError
				}
			}
			return handler(ctx, req)
		}
	}
}
