package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	consulAPI "github.com/hashicorp/consul/api"
	"material/internal/conf"
	"material/internal/pkg/event/event"
	"material/internal/pkg/event/kafka"
	ctos "material/internal/pkg/volcengine/tos"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	grpcx "google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	slog "log"
	companyv1 "material/api/service/company/v1"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewKafka, NewTos, NewDiscovery, NewMaterialCategoryRepo, NewMaterialRepo, NewMaterialProductRepo, NewCollectRepo, NewCompanyMaterialRepo, NewCompanyProductRepo, NewCompanyServiceClient)

// Data .
type Data struct {
	db        *gorm.DB
	rdb       *redis.Client
	kafka     event.Sender
	tos       *ctos.Tos
	companyuc companyv1.CompanyClient
	conf      *conf.Data
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB, rdb *redis.Client, kafka event.Sender, tos *ctos.Tos, companyuc companyv1.CompanyClient, logger log.Logger) (*Data, func(), error) {
	data := &Data{conf: c, db: db, rdb: rdb, kafka: kafka, tos: tos, companyuc: companyuc}

	cleanup := func() {
		sqlDB, _ := data.db.DB()
		sqlDB.Close()

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

	_ = db.AutoMigrate(&Collect{})

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

func NewKafka(e *conf.Event) event.Sender {
	sender, err := kafka.NewKafkaSender(strings.Split(e.Kafka.Addr, ","), e.Kafka.Topic)

	if err != nil {
		panic("failed to connect kafka")
	}

	return sender
}

func NewTos(c *conf.Volcengine) *ctos.Tos {
	client := &ctos.Tos{
		AccessKey:  c.Tos.AccessKey,
		SecretKey:  c.Tos.SecretKey,
		Endpoint:   c.Tos.Material.Endpoint,
		Region:     c.Tos.Material.Region,
		BucketName: c.Tos.Material.BucketName,
	}

	err := client.NewClient()

	if err != nil {
		panic("failed to connect tos")
	}

	return client
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
