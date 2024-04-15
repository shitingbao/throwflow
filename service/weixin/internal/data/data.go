package data

import (
	"context"
	slog "log"
	"os"
	"time"
	"unicode/utf8"
	commonv1 "weixin/api/service/common/v1"
	companyv1 "weixin/api/service/company/v1"
	douyinv1 "weixin/api/service/douyin/v1"
	"weixin/internal/biz"
	"weixin/internal/conf"
	ctos "weixin/internal/pkg/volcengine/tos"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/sony/sonyflake"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	grpcx "google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewMongo, NewSonyflake, NewTos, NewDiscovery, NewTransaction, NewUserRepo, NewUserOpenIdRepo, NewCompanyRepo, NewCompanyProductRepo, NewUserOpenDouyinRepo, NewOpenDouyinUserInfoRepo, NewUserAddressRepo, NewUserSampleOrderRepo, NewUserCommissionRepo, NewJinritemaiOrderRepo, NewTaskLogRepo, NewQrCodeRepo, NewCompanyOrganizationRepo, NewUserOrderRepo, NewUserOrganizationRelationRepo, NewUserScanRecordRepo, NewUserCouponRepo, NewUserCouponCreateLogRepo, NewAreaRepo, NewPayRepo, NewShortUrlRepo, NewShortCodeRepo, NewDjAwemeRepo, NewTuUserRepo, NewCouponUserRepo, NewUserIntegralRelationRepo, NewUserContractRepo, NewUserBankRepo, NewUserBalanceLogRepo, NewKuaidiInfoRepo, NewDoukeOrderRepo, NewAwemesAdvertiserWeixinAuthRepo, NewCourseRepo, NewCourseUserRepo, NewCompanyTaskRepo, NewCompanyServiceClient, NewDouyinServiceClient, NewCommonServiceClient)

type tos struct {
	name string
	tos  *ctos.Tos
}

// Data .
type Data struct {
	db        *gorm.DB
	rdb       *redis.Client
	mdb       *mongo.Client
	sonyflake *sonyflake.Sonyflake
	toses     []*tos
	companyuc companyv1.CompanyClient
	douyinuc  douyinv1.DouyinClient
	commonuc  commonv1.CommonClient
	conf      *conf.Data
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
func NewData(c *conf.Data, db *gorm.DB, rdb *redis.Client, mdb *mongo.Client, sonyflake *sonyflake.Sonyflake, toses []*tos, companyuc companyv1.CompanyClient, douyinuc douyinv1.DouyinClient, commonuc commonv1.CommonClient, logger log.Logger) (*Data, func(), error) {
	data := &Data{conf: c, db: db, rdb: rdb, mdb: mdb, sonyflake: sonyflake, toses: toses, companyuc: companyuc, douyinuc: douyinuc, commonuc: commonuc}

	cleanup := func() {
		sqlDB, _ := data.db.DB()
		sqlDB.Close()

		mdb.Disconnect(context.TODO())

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

	_ = db.AutoMigrate(&User{})
	_ = db.AutoMigrate(&UserOpenId{})
	_ = db.AutoMigrate(&UserOpenDouyin{})
	_ = db.AutoMigrate(&UserAddress{})
	_ = db.AutoMigrate(&UserSampleOrder{})
	_ = db.AutoMigrate(&UserOrganizationRelation{})
	_ = db.AutoMigrate(&UserIntegralRelation{})
	_ = db.AutoMigrate(&UserOrder{})
	_ = db.AutoMigrate(&UserScanRecord{})
	_ = db.AutoMigrate(&UserCommission{})
	_ = db.AutoMigrate(&UserCouponCreateLog{})
	_ = db.AutoMigrate(&UserCoupon{})
	_ = db.AutoMigrate(&UserContract{})
	_ = db.AutoMigrate(&UserBank{})
	_ = db.AutoMigrate(&UserBalanceLog{})
	_ = db.AutoMigrate(&Course{})
	_ = db.AutoMigrate(&CourseUser{})

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

func NewMongo(c *conf.Data) *mongo.Client {
	var mdb *mongo.Client
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), c.Mongo.Timeout.AsDuration())
	defer cancel()

	/*var logMonitor = event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Printf("mongo reqId:%d start on db:%s cmd:%s sql:%+v", startedEvent.RequestID, startedEvent.DatabaseName,
				startedEvent.CommandName, startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
			fmt.Printf("mongo reqId:%d exec cmd:%s success duration %d ns", succeededEvent.RequestID,
				succeededEvent.CommandName, succeededEvent.DurationNanos)
		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
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

func NewSonyflake() *sonyflake.Sonyflake {
	sonyflake, err := sonyflake.New(sonyflake.Settings{
		StartTime: time.Now(),
	})

	if err != nil {
		panic("failed to initialization sonyflake")
	}

	return sonyflake
}

func NewTos(c *conf.Volcengine) []*tos {
	toses := make([]*tos, 0)

	companyClient := &ctos.Tos{
		AccessKey:  c.Tos.AccessKey,
		SecretKey:  c.Tos.SecretKey,
		Endpoint:   c.Tos.Company.Endpoint,
		Region:     c.Tos.Company.Region,
		BucketName: c.Tos.Company.BucketName,
	}

	err := companyClient.NewClient()

	if err != nil {
		panic("failed to connect tos")
	}

	toses = append(toses, &tos{
		name: "company",
		tos:  companyClient,
	})

	avatarClient := &ctos.Tos{
		AccessKey:  c.Tos.AccessKey,
		SecretKey:  c.Tos.SecretKey,
		Endpoint:   c.Tos.Avatar.Endpoint,
		Region:     c.Tos.Avatar.Region,
		BucketName: c.Tos.Avatar.BucketName,
	}

	err = avatarClient.NewClient()

	if err != nil {
		panic("failed to connect tos")
	}

	toses = append(toses, &tos{
		name: "avatar",
		tos:  avatarClient,
	})

	return toses
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

func NewDouyinServiceClient(sr *conf.Service, rr registry.Discovery) douyinv1.DouyinClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Douyin.GetEndpoint()),
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
