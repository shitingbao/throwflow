package data

import (
	companyv1 "douyin/api/service/company/v1"
	weixinv1 "douyin/api/service/weixin/v1"
	"douyin/internal/biz"
	"douyin/internal/conf"
	"douyin/internal/pkg/event/event"
	"douyin/internal/pkg/event/kafka"
	"github.com/ClickHouse/clickhouse-go/v2"
	cdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	consulAPI "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	grpcx "google.golang.org/grpc"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	slog "log"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRdsDB, NewMongo, NewClickhouse, NewRedis, NewKafka, NewDiscovery, NewOceanengineConfigRepo, NewOceanengineAccountRepo, NewQianchuanAdvertiserRepo, NewQianchuanAdvertiserStatusRepo, NewOceanengineAccountTokenRepo, NewCompanyUserRepo, NewCompanyRepo, NewCompanyProductRepo, NewOceanengineApiLogRepo, NewQianchuanCampaignRepo, NewQianchuanProductRepo, NewQianchuanAwemeRepo, NewQianchuanWalletRepo, NewLianshanRealtimeRepo, NewQianchuanAdRepo, NewQianchuanReportAdRepo, NewQianchuanReportProductRepo, NewQianchuanReportAwemeRepo, NewCompanySetRepo, NewQianchuanAdInfoRepo, NewQianchuanAdvertiserInfoRepo, NewQianchuanReportAdRealtimeRepo, NewQianchuanAdvertiserHistoryRepo, NewOpenDouyinTokenRepo, NewOpenDouyinUserInfoRepo, NewOpenDouyinUserInfoCreateLogRepo, NewWeixinUserOpenDouyinRepo, NewOpenDouyinApiLogRepo, NewJinritemaiApiLogRepo, NewTaskLogRepo, NewWeixinUserRepo, NewJinritemaiOrderRepo, NewJinritemaiOrderInfoRepo, NewJinritemaiStoreInfoRepo, NewJinritemaiStoreRepo, NewQianchuanAwemeOrderInfoRepo, NewOpenDouyinVideoRepo, NewDoukeOrderRepo, NewDoukeOrderInfoRepo, NewCsjApiLogRepo, NewWeixinUserCommissionRepo, NewTransaction, NewCompanyServiceClient, NewWeixinServiceClient)

type lianshandb struct {
	appId string
	db    *gorm.DB
}

// Data .
type Data struct {
	db          *gorm.DB
	lianshandbs []*lianshandb
	mdb         *mongo.Client
	cdb         cdriver.Conn
	rdb         *redis.Client
	kafka       event.Sender
	companyuc   companyv1.CompanyClient
	weixinuc    weixinv1.WeixinClient
	conf        *conf.Data
}

type contextTxKey struct{}

func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB, lianshandbs []*lianshandb, mdb *mongo.Client, cdb cdriver.Conn, rdb *redis.Client, kafka event.Sender, companyuc companyv1.CompanyClient, weixinuc weixinv1.WeixinClient, logger log.Logger) (*Data, func(), error) {
	data := &Data{conf: c, db: db, lianshandbs: lianshandbs, mdb: mdb, cdb: cdb, rdb: rdb, kafka: kafka, companyuc: companyuc, weixinuc: weixinuc}

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		for _, llianshandb := range lianshandbs {
			sqlDB, _ = llianshandb.db.DB()
			sqlDB.Close()
		}

		mdb.Disconnect(context.TODO())

		cdb.Close()

		rdb.Close()

		kafka.Close()
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

	_ = db.AutoMigrate(&OceanengineConfig{})
	_ = db.AutoMigrate(&OceanengineAccount{})
	_ = db.AutoMigrate(&QianchuanAdvertiser{})
	_ = db.AutoMigrate(&OceanengineAccountToken{})
	_ = db.AutoMigrate(&QianchuanAdvertiserStatus{})
	_ = db.AutoMigrate(&QianchuanAdvertiserHistory{})
	_ = db.AutoMigrate(&OpenDouyinToken{})
	_ = db.AutoMigrate(&OpenDouyinUserInfo{})
	_ = db.AutoMigrate(&OpenDouyinUserInfoCreateLog{})
	_ = db.AutoMigrate(&JinritemaiOrderInfo{})
	_ = db.AutoMigrate(&DoukeOrderInfo{})

	return db
}

func NewRdsDB(c *conf.Data, log log.Logger) []*lianshandb {
	lianshandbs := make([]*lianshandb, 0)

	for _, dsn := range c.Lianshan.Dsns {
		if db, err := gorm.Open(mysql.Open(dsn.Dsn), &gorm.Config{}); err != nil {
			panic("failed to connect rds database")
		} else {
			lianshandbs = append(lianshandbs, &lianshandb{
				appId: dsn.AppId,
				db:    db,
			})
		}
	}

	return lianshandbs
}

func NewMongo(c *conf.Data) *mongo.Client {
	var mdb *mongo.Client
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), c.Mongo.Timeout.AsDuration())
	defer cancel()

	/*var logMonitor = event2.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event2.CommandStartedEvent) {
			fmt.Printf("mongo reqId:%d start on db:%s cmd:%s sql:%+v", startedEvent.RequestID, startedEvent.DatabaseName,
				startedEvent.CommandName, startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event2.CommandSucceededEvent) {
			fmt.Printf("mongo reqId:%d exec cmd:%s success duration %d ns", succeededEvent.RequestID,
				succeededEvent.CommandName, succeededEvent.DurationNanos)
		},
		Failed: func(ctx context.Context, failedEvent *event2.CommandFailedEvent) {
			fmt.Printf("mongo reqId:%d exec cmd:%s failed duration %d ns", failedEvent.RequestID,
				failedEvent.CommandName, failedEvent.DurationNanos)
		},
	}*/

	if utf8.RuneCountInString(c.Mongo.Username) > 0 && utf8.RuneCountInString(c.Mongo.Password) > 0 {
		mdb, err = mongo.Connect(ctx, options.Client().
			ApplyURI(c.Mongo.Dsn).
			SetAuth(options.Credential{
				Username: c.Mongo.Username,
				Password: c.Mongo.Password,
			}))
	} else {
		mdb, err = mongo.Connect(ctx, options.Client().
			ApplyURI(c.Mongo.Dsn))
	}

	if err != nil {
		panic("failed to connect mongo")
	}

	err = mdb.Ping(ctx, readpref.Primary())

	if err != nil {
		panic("failed to connect mongo")
	}

	return mdb
}

func NewClickhouse(c *conf.Data) cdriver.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{c.Clickhouse.Addr},
		Auth: clickhouse.Auth{
			Database: c.Clickhouse.Db,
			Username: c.Clickhouse.Username,
			Password: c.Clickhouse.Password,
		},
		MaxOpenConns: int(c.Clickhouse.MaxOpenConns),
		ReadTimeout:  time.Second * 1,
	})

	if err != nil {
		panic("failed to connect database")
	}

	return conn
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

func NewKafka(e *conf.Event) event.Sender {
	sender, err := kafka.NewKafkaSender(strings.Split(e.Kafka.Addr, ","), e.Kafka.Topic)

	if err != nil {
		panic("failed to connect kafka")
	}

	return sender
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
		grpc.WithTimeout(2*time.Second),
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
		grpc.WithTimeout(2*time.Second),
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
