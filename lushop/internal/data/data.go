package data

import (
	"context"
	"lushop/internal/conf"
	"time"

	goodsV1 "lushop/api/service/goods/v1"
	inventoryV1 "lushop/api/service/inventory/v1"
	orderV1 "lushop/api/service/order/v1"
	userV1 "lushop/api/service/user/v1"
	userauthV1 "lushop/api/service/userauth/v1"
	useropV1 "lushop/api/service/userop/v1"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	grpcx "google.golang.org/grpc"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewUserAuthRepoGRPC, // 使用 gRPC 版本的用户仓库（统一治理方案）
	// NewUserRepo,       // 旧的 Redis 版本（如需回滚可切换）
	NewCartRepo,      // 购物车仓库
	NewGoodsRepo,     // 商品仓库
	NewInventoryRepo, // 库存仓库
	NewOrderRepo,     // 订单仓库
	NewUserOpRepo,    // 用户操作仓库
	NewUserServiceClient,
	NewUserAuthClient,
	NewGoodsServiceClient,     // Goods Service 客户端
	NewInventoryServiceClient, // Inventory Service 客户端
	NewOrderServiceClient,     // Order Service 客户端
	NewUserOpServiceClient,    // UserOp Service 客户端
	NewRegister,
	NewDiscovery,
	NewRedis,
)

// Data .
type Data struct {
	log *log.Helper
	uc  userV1.UserClient           // 用户服务的客户端
	uac userauthV1.UserAuthClient   // 用户认证服务的客户端
	gc  goodsV1.GoodsClient         // 商品服务的客户端
	ic  inventoryV1.InventoryClient // 库存服务的客户端
	oc  orderV1.OrderClient         // 订单服务的客户端
	uoc useropV1.UserOpClient       // 用户操作服务的客户端
	rdb *redis.Client
}

// Rdb 返回共享的 Redis 客户端，供其他组件（如任务服务）复用
func (d *Data) Rdb() *redis.Client {
	return d.rdb
}

// NewData .
func NewData(c *conf.Data, uc userV1.UserClient, uac userauthV1.UserAuthClient, gc goodsV1.GoodsClient, ic inventoryV1.InventoryClient, oc orderV1.OrderClient, uoc useropV1.UserOpClient, logger log.Logger, rdb *redis.Client) (*Data, error) {
	l := log.NewHelper(log.With(logger, "module", "data"))
	return &Data{log: l, uc: uc, uac: uac, gc: gc, ic: ic, oc: oc, uoc: uoc, rdb: rdb}, nil
}

// NewUserServiceClient 链接用户grpc服务
func NewUserServiceClient(ac *conf.Auth, sr *conf.Service, rr registry.Discovery) userV1.UserClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.User.Endpoint), // consul
		grpc.WithDiscovery(rr),              // consul
		grpc.WithMiddleware(
			tracing.Client(), // 链路追踪
			recovery.Recovery(),
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := userV1.NewUserClient(conn)
	return c
}

// NewUserAuthClient 链接用户认证grpc服务
func NewUserAuthClient(ac *conf.Auth, sr *conf.Service, rr registry.Discovery) userauthV1.UserAuthClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.UserAuth.Endpoint), // consul
		grpc.WithDiscovery(rr),                  // consul
		grpc.WithMiddleware(
			tracing.Client(), // 链路追踪
			recovery.Recovery(),
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := userauthV1.NewUserAuthClient(conn)
	return c
}

// NewGoodsServiceClient 链接商品服务grpc服务
func NewGoodsServiceClient(ac *conf.Auth, sr *conf.Service, rr registry.Discovery) goodsV1.GoodsClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Goods.Endpoint), // consul
		grpc.WithDiscovery(rr),               // consul
		grpc.WithMiddleware(
			tracing.Client(), // 链路追踪
			recovery.Recovery(),
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := goodsV1.NewGoodsClient(conn)
	return c
}

// NewOrderServiceClient 链接订单服务grpc服务
func NewOrderServiceClient(ac *conf.Auth, sr *conf.Service, rr registry.Discovery) orderV1.OrderClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Order.Endpoint), // consul
		grpc.WithDiscovery(rr),               // consul
		grpc.WithMiddleware(
			tracing.Client(), // 链路追踪
			recovery.Recovery(),
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := orderV1.NewOrderClient(conn)
	return c
}

// NewInventoryServiceClient 创建库存服务客户端
func NewInventoryServiceClient(sr *conf.Service, rr registry.Discovery) inventoryV1.InventoryClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Inventory.Endpoint), // consul
		grpc.WithDiscovery(rr),                   // consul
		grpc.WithMiddleware(
			tracing.Client(), // 链路追踪
			recovery.Recovery(),
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := inventoryV1.NewInventoryClient(conn)
	return c
}

// NewUserOpServiceClient 创建用户操作服务客户端
func NewUserOpServiceClient(sr *conf.Service, rr registry.Discovery) useropV1.UserOpClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Userop.Endpoint), // consul
		grpc.WithDiscovery(rr),                // consul
		grpc.WithMiddleware(
			tracing.Client(), // 链路追踪
			recovery.Recovery(),
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := useropV1.NewUserOpClient(conn)
	return c
}

// NewRegister add consul
func NewRegister(conf *conf.Registry) registry.Registrar {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

// NewDiscovery
func NewDiscovery(conf *conf.Registry) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
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
	// 测试连接是否正常
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Errorf("Failed to connect to Redis: %v", err)
		// 不要关闭连接，让调用者处理错误
	}
	return rdb
}
