package data

import (
	"common/internal/biz"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type tokenRepo struct {
	data *Data
	log  *log.Helper
}

func NewTokenRepo(data *Data, logger log.Logger) biz.TokenRepo {
	return &tokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (tr *tokenRepo) Save(ctx context.Context, key string, val string, timeout time.Duration) error {
	if _, err := tr.data.rdb.Set(ctx, key, val, timeout).Result(); err != nil {
		return err
	}

	return nil
}

func (tr *tokenRepo) Verify(ctx context.Context, key string) error {
	luaScript := "local keyExists = redis.call('EXISTS', KEYS[1]) \n if keyExists == 1 then \n redis.call('DEL', KEYS[1]) \n return 1 \n else \n return 0 \n end"

	var luaKeys []string
	luaKeys = append(luaKeys, key)

	result, err := tr.data.rdb.Eval(ctx, luaScript, luaKeys).Result()

	if err != nil {
		return err
	}

	tr.data.rdb.Del(ctx, key).Result()

	if result == int64(0) {
		return errors.New("TOKEN验证失败")
	}

	return nil
}
