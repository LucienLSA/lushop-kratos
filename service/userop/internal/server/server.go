package server

import (
	"context"
	"os"
	"time"
	"userop/internal/conf"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewRegistrar)

// NewRegistrar 创建 Consul 注册器，支持健康检查
func NewRegistrar(conf *conf.Registry) registry.Registrar {
	// 检查是否禁用注册或配置为空
	if os.Getenv("REGISTRY_DISABLED") == "true" || conf == nil {
		log.Warnf("❌ Service registration DISABLED - using noOpRegistrar")
		return &noOpRegistrar{}
	}

	// 检查 Consul 配置
	if conf.Consul == nil || conf.Consul.Address == "" {
		log.Warnf("⚠️  Consul configuration is missing - using noOpRegistrar")
		return &noOpRegistrar{}
	}

	// 创建 Consul 客户端
	config := consulAPI.DefaultConfig()
	config.Address = conf.Consul.Address
	if conf.Consul.Scheme != "" {
		config.Scheme = conf.Consul.Scheme
	}
	if config.HttpClient != nil {
		config.HttpClient.Timeout = 10 * time.Second
	}

	client, err := consulAPI.NewClient(config)
	if err != nil {
		log.Errorf("consul client init failed: %v, falling back to no-op registrar", err)
		return &noOpRegistrar{}
	}

	// 测试连接
	if _, err := client.Agent().Self(); err != nil {
		log.Errorf("consul connection failed: %v", err)
		return &noOpRegistrar{}
	}

	log.Infof("✅ Consul registrar initialized successfully! Address: %s", conf.Consul.Address)

	// 创建注册器，启用健康检查
	return consul.New(client, consul.WithHealthCheck(true))
}

// noOpRegistrar 禁用服务注册的回退实现
type noOpRegistrar struct{}

func (noOpRegistrar) Register(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}

func (noOpRegistrar) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}

// NewDiscovery 创建 Consul 服务发现客户端
func NewDiscovery(conf *conf.Registry) registry.Discovery {
	if conf == nil || conf.Consul == nil || conf.Consul.Address == "" {
		return nil
	}

	config := consulAPI.DefaultConfig()
	config.Address = conf.Consul.Address
	config.Scheme = conf.Consul.Scheme

	client, err := consulAPI.NewClient(config)
	if err != nil {
		return nil
	}

	return consul.New(client)
}
