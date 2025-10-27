package data

import (
	"context"
	slog "log"
	"order/internal/biz"
	"order/internal/conf"
	"order/internal/pkg/rocketmq"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewRedis,
	NewRocketMQProducer,
	NewOrderRepo,
)

// Data .
type Data struct {
	db          *gorm.DB
	rdb         *redis.Client
	producer    *rocketmq.Producer
	txProducer  *rocketmq.TransactionProducer
	goodsClient biz.GoodsService
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, rdb *redis.Client, producer *rocketmq.Producer, goodsClient biz.GoodsService) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if producer != nil {
			producer.Close()
		}
	}
	return &Data{
		db:          db,
		rdb:         rdb,
		producer:    producer,
		txProducer:  nil, // 将在 NewOrderRepo 中设置
		goodsClient: goodsClient,
	}, cleanup, nil
}

// 用来承载事务的上下文
type contextTxKey struct{}

// NewTransaction .
func NewTransaction(d *Data) biz.Transaction {
	return d
}

// ExecTx gorm Transaction
func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// DB 根据此方法来判断当前的 db 是不是使用 事务的 DB
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

// NewDB .
func NewDB(c *conf.Data) *gorm.DB {
	if c == nil || c.Database == nil {
		panic("database configuration is nil")
	}
	// 终端打印输入 sql 执行记录
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢查询 SQL 阈值
			Colorful:      true,        // 禁用彩色打印
			//IgnoreRecordNotFoundError: false,
			LogLevel: logger.Info, // Log lever
		},
	)

	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名是否加 s
		},
	})

	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		panic("failed to connect database")
	}
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
	rdb.AddHook(redisotel.TracingHook{})
	return rdb
}

// NewRocketMQProducer 创建 RocketMQ 生产者
func NewRocketMQProducer(c *conf.Bootstrap, logger log.Logger) *rocketmq.Producer {
	// 如果未启用 RocketMQ，返回 nil
	if c.Rocketmq == nil || !c.Rocketmq.Enable {
		log.NewHelper(logger).Info("RocketMQ is disabled")
		return nil
	}

	producer, err := rocketmq.NewProducer(
		c.Rocketmq.NameServer,
		c.Rocketmq.GroupName,
		c.Rocketmq.Topic,
		logger,
	)
	if err != nil {
		log.NewHelper(logger).Errorf("failed to create RocketMQ producer: %v", err)
		return nil
	}

	log.NewHelper(logger).Info("RocketMQ producer created successfully")
	return producer
}
