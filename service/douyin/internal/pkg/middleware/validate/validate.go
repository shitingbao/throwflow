package validate

import (
	"context"
	"douyin/internal/biz"
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
				if err := v.Validate(); err != nil {
					spew.Dump(err.Error())
					return nil, biz.DouyinValidatorError
				}
			}
			return handler(ctx, req)
		}
	}
}
