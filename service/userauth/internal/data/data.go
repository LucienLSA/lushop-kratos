package data

import (
	"context"
	"time"

	"userauth/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedis, NewAuthRepo)

// Data .
type Data struct {
	log *log.Helper
	rdb *redis.Client
}

// NewData .
func NewData(rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data"))
	cleanup := func() {
		l.Info("closing the data resources")
		if err := rdb.Close(); err != nil {
			l.Errorf("failed to close redis: %v", err)
		}
	}
	return &Data{
		log: l,
		rdb: rdb,
	}, cleanup, nil
}

// NewRedis 创建 Redis 客户端
func NewRedis(conf *conf.Data, logger log.Logger) *redis.Client {
	l := log.NewHelper(log.With(logger, "module", "data/redis"))

	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
	})

	// 启用 OpenTelemetry 追踪
	rdb.AddHook(redisotel.TracingHook{})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		l.Fatalf("failed to connect redis: %v", err)
	}

	l.Info("redis connected successfully")
	return rdb
}
