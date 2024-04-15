package data

import (
	commonv1 "company/api/service/common/v1"
	douyinv1 "company/api/service/douyin/v1"
	materialv1 "company/api/service/material/v1"
	weixinv1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"company/internal/conf"
	"company/internal/pkg/event/event"
	"company/internal/pkg/event/kafka"
	ctos "company/internal/pkg/volcengine/tos"
	"context"
	slog "log"
	"os"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	grpcx "google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewClickhouse, NewKafka, NewTos, NewDiscovery, NewTransaction, NewMenuRepo, NewIndustryRepo, NewClueRepo, NewCompanyRepo, NewCompanyUserRepo, NewCompanyUserCompanyRepo, NewCompanyUserRoleRepo, NewQianchuanAdvertiserRepo, NewCompanyUserWhiteRepo, NewCompanyPerformanceRuleRepo, NewCompanyPerformanceDailyRepo, NewCompanyPerformanceRebalanceRepo, NewCompanyPerformanceMonthlyRepo, NewQianchuanAdAdvertiserRepo, NewCompanySetRepo, NewQianchuanReportAdvertiserRepo, NewQianchuanReportProductRepo, NewQianchuanReportAwemeRepo, NewQianchuanAdRepo, NewCompanyUserQianchuanSearchRepo, NewQianchuanAdvertiserHistoryRepo, NewSmsRepo, NewAreaRepo, NewCompanyProductRepo, NewCompanyProductCategoryRepo, NewQrCodeRepo, NewCompanyMaterialRepo, NewCompanyMaterialLibraryRepo, NewOpenDouyinUserInfoRepo, NewJinritemaiStoreRepo, NewWeixinUserScanRecordRepo, NewJinritemaiOrderRepo, NewMaterialMaterialRepo, NewCompanyOrganizationRepo, NewShortUrlRepo, NewShortCodeRepo, NewDoukeProductRepo, NewDoukeOrderRepo, NewCompanyTaskRepo, NewCompanyTaskAccountRelationRepo, NewCompanyTaskDetailRepo, NewWeixinUserRepo, NewWeixinUserOpenDouyinRepo, NewWeixinUserCommissionRepo, NewDouyinServiceClient, NewCommonServiceClient, NewWeixinServiceClient, NewMaterialServiceClient)

type tos struct {
	name string
	tos  *ctos.Tos
}

// Data .
type Data struct {
	db         *gorm.DB
	rdb        *redis.Client
	cdb        driver.Conn
	kafka      event.Sender
	toses      []*tos
	douyinuc   douyinv1.DouyinClient
	commonuc   commonv1.CommonClient
	weixinuc   weixinv1.WeixinClient
	materialuc materialv1.MaterialClient
	conf       *conf.Data
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
func NewData(c *conf.Data, db *gorm.DB, rdb *redis.Client, cdb driver.Conn, kafka event.Sender, toses []*tos, douyinuc douyinv1.DouyinClient, commonuc commonv1.CommonClient, weixinuc weixinv1.WeixinClient, materialuc materialv1.MaterialClient, logger log.Logger) (*Data, func(), error) {
	data := &Data{conf: c, db: db, rdb: rdb, cdb: cdb, kafka: kafka, toses: toses, douyinuc: douyinuc, commonuc: commonuc, weixinuc: weixinuc, materialuc: materialuc}

	cleanup := func() {
		sqlDB, _ := data.db.DB()
		sqlDB.Close()

		rdb.Close()

		cdb.Close()

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

	_ = db.AutoMigrate(&Menu{})
	_ = db.AutoMigrate(&Clue{})
	_ = db.AutoMigrate(&Company{})
	_ = db.AutoMigrate(&CompanyUser{})
	_ = db.AutoMigrate(&CompanyUserCompany{})
	_ = db.AutoMigrate(&CompanyUserWhite{})
	_ = db.AutoMigrate(&CompanyUserRole{})
	_ = db.AutoMigrate(&CompanyPerformanceRule{})
	_ = db.AutoMigrate(&CompanyPerformanceRebalance{})
	_ = db.AutoMigrate(&CompanyPerformanceDaily{})
	_ = db.AutoMigrate(&CompanyPerformanceMonthly{})
	_ = db.AutoMigrate(&CompanySet{})
	_ = db.AutoMigrate(&CompanyUserQianchuanSearch{})
	_ = db.AutoMigrate(&CompanyProduct{})
	_ = db.AutoMigrate(&CompanyMaterialLibrary{})
	_ = db.AutoMigrate(&CompanyOrganization{})
	_ = db.AutoMigrate(&CompanyTask{})
	_ = db.AutoMigrate(&CompanyTaskAccountRelation{})
	_ = db.AutoMigrate(&CompanyTaskDetail{})

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

func NewClickhouse(c *conf.Data) driver.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{c.Clickhouse.Addr},
		Auth: clickhouse.Auth{
			Database: c.Clickhouse.Db,
			Username: c.Clickhouse.Username,
			Password: c.Clickhouse.Password,
		},
	})

	if err != nil {
		panic("failed to connect database")
	}

	return conn
}

func NewKafka(e *conf.Event) event.Sender {
	sender, err := kafka.NewKafkaSender(strings.Split(e.Kafka.Addr, ","), e.Kafka.Topic)

	if err != nil {
		panic("failed to connect kafka")
	}

	return sender
}

func NewTos(c *conf.Volcengine) []*tos {
	toses := make([]*tos, 0)

	productClient := &ctos.Tos{
		AccessKey:  c.Tos.AccessKey,
		SecretKey:  c.Tos.SecretKey,
		Endpoint:   c.Tos.Product.Endpoint,
		Region:     c.Tos.Product.Region,
		BucketName: c.Tos.Product.BucketName,
	}

	err := productClient.NewClient()

	if err != nil {
		panic("failed to connect tos")
	}

	toses = append(toses, &tos{
		name: "product",
		tos:  productClient,
	})

	materialClient := &ctos.Tos{
		AccessKey:  c.Tos.AccessKey,
		SecretKey:  c.Tos.SecretKey,
		Endpoint:   c.Tos.Material.Endpoint,
		Region:     c.Tos.Material.Region,
		BucketName: c.Tos.Material.BucketName,
	}

	err = materialClient.NewClient()

	if err != nil {
		panic("failed to connect tos")
	}

	toses = append(toses, &tos{
		name: "material",
		tos:  materialClient,
	})

	organizationClient := &ctos.Tos{
		AccessKey:  c.Tos.AccessKey,
		SecretKey:  c.Tos.SecretKey,
		Endpoint:   c.Tos.Organization.Endpoint,
		Region:     c.Tos.Organization.Region,
		BucketName: c.Tos.Organization.BucketName,
	}

	err = organizationClient.NewClient()

	if err != nil {
		panic("failed to connect tos")
	}

	toses = append(toses, &tos{
		name: "organization",
		tos:  organizationClient,
	})

	taskClient := &ctos.Tos{
		AccessKey:  c.Tos.AccessKey,
		SecretKey:  c.Tos.SecretKey,
		Endpoint:   c.Tos.Task.Endpoint,
		Region:     c.Tos.Task.Region,
		BucketName: c.Tos.Task.BucketName,
	}

	err = taskClient.NewClient()

	if err != nil {
		panic("failed to connect tos")
	}

	toses = append(toses, &tos{
		name: "task",
		tos:  taskClient,
	})

	return toses
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
		grpc.WithTimeout(5*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := douyinv1.NewDouyinClient(conn)
	return c
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
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := commonv1.NewCommonClient(conn)
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
		grpc.WithTimeout(5*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := weixinv1.NewWeixinClient(conn)
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
		grpc.WithTimeout(5*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := materialv1.NewMaterialClient(conn)
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
