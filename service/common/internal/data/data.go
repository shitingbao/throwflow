package data

import (
	"common/internal/conf"
	"context"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	slog "log"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewSmsLogRepo, NewTokenRepo, NewAreaRepo, NewUpdateLogRepo, NewPayAsyncNotificationLogRepo, NewShortCodeLogRepo, NewKuaidiCompanyRepo, NewKuaidiInfoRepo)

// Data .
type Data struct {
	db   *gorm.DB
	rdb  *redis.Client
	conf *conf.Data
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB, rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	data := &Data{conf: c, db: db, rdb: rdb}

	cleanup := func() {
		sqlDB, _ := data.db.DB()
		sqlDB.Close()

		rdb.Close()
	}

	return data, cleanup, nil
}

func NewDB(c *conf.Data, log log.Logger) *gorm.DB {
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢查询 SQL 阈值
			Colorful:      true,        // 禁用彩色打印
			LogLevel:      logger.Info, // Log lever
		},
	)

	db, err := gorm.Open(mysql.Open(c.GetDatabase().GetDsn()), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect database")
	}

	_ = db.AutoMigrate(&SmsLog{})
	_ = db.AutoMigrate(&Area{})
	_ = db.AutoMigrate(&PayAsyncNotificationLog{})
	_ = db.AutoMigrate(&ShortCodeLog{})
	_ = db.AutoMigrate(&KuaidiCompany{})
	_ = db.AutoMigrate(&KuaidiInfo{})

	return db
}

func NewRedis(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), c.Redis.DialTimeout.AsDuration())
	defer cancelFunc()

	err := rdb.Ping(timeout).Err()

	if err != nil {
		panic("failed to connect redis")
	}

	return rdb
}
