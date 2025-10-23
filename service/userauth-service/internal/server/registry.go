package server

import (
	"userauth-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/hashicorp/consul/api"
)

// NewRegistrar 创建 Consul 注册中心
func NewRegistrar(conf *conf.Registry, logger log.Logger) registry.Registrar {
	l := log.NewHelper(log.With(logger, "module", "server/registry"))
	
	c := api.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	
	client, err := api.NewClient(c)
	if err != nil {
		l.Fatalf("failed to create consul client: %v", err)
	}
	
	l.Info("consul registrar created successfully")
	
	// 使用 Kratos 的 Consul 注册实现
	reg := consul.New(client, consul.WithHealthCheck(true))
	return reg
}
