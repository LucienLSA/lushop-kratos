package goods

import (
	"context"
	"fmt"
	goodsV1 "goods/api/goods/v1"
	"order/internal/biz"
	"order/internal/conf"
	"order/internal/domain"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	consulAPI "github.com/hashicorp/consul/api"
)

// client 商品服务客户端（适配器 Adapter）
// 实现 biz.GoodsService 接口
type client struct {
	grpcClient goodsV1.GoodsClient
	log        *log.Helper
}

// NewGoodsServiceClient 创建商品服务客户端
// 返回 biz.GoodsService 接口，而不是具体类型
func NewGoodsServiceClient(c *conf.Service, r *conf.Registry, logger log.Logger) (biz.GoodsService, error) {
	l := log.NewHelper(logger)

	// 创建 Consul 注册中心客户端
	consulConfig := consulAPI.DefaultConfig()
	consulConfig.Address = r.Consul.Address
	consulConfig.Scheme = r.Consul.Scheme

	consulClient, err := consulAPI.NewClient(consulConfig)
	if err != nil {
		l.Errorf("failed to create consul client: %v", err)
		return nil, err
	}

	// 创建服务发现
	dis := consul.New(consulClient)

	// 创建 gRPC 连接
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.Goods.Endpoint),
		grpc.WithDiscovery(dis),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		l.Errorf("failed to dial goods service: %v", err)
		return nil, err
	}

	// 创建商品服务客户端
	grpcClient := goodsV1.NewGoodsClient(conn)
	l.Info("goods service client created successfully")

	return &client{
		grpcClient: grpcClient,
		log:        l,
	}, nil
}

// BatchGetGoods 实现 biz.GoodsService 接口
func (c *client) BatchGetGoods(ctx context.Context, ids []int32) (map[int32]*domain.GoodsInfo, error) {
	// 调用商品服务的 BatchGetGoods 接口
	resp, err := c.grpcClient.BatchGetGoods(ctx, &goodsV1.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		c.log.Errorf("failed to batch get goods: ids=%v, error=%v", ids, err)
		return nil, fmt.Errorf("failed to get goods info: %w", err)
	}

	// 转换为领域模型（biz.GoodsInfo）
	result := make(map[int32]*domain.GoodsInfo)
	for _, goods := range resp.Data {
		result[goods.Id] = &domain.GoodsInfo{
			ID:              goods.Id,
			Name:            goods.Name,
			GoodsFrontImage: goods.GoodsFrontImage,
			ShopPrice:       goods.ShopPrice,
		}
	}

	c.log.Infof("batch get goods success: ids=%v, count=%d", ids, len(result))
	return result, nil
}
