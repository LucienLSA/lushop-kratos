package server

import (
	"cart/internal/conf"
	"context"
	"os"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewRegistrar)

// NewRegistrar 引入 consul
func NewRegistrar(conf *conf.Registry) registry.Registrar {
	if os.Getenv("REGISTRY_DISABLED") == "true" || conf == nil || conf.Consul == nil || conf.Consul.Address == "" {
		return noOpRegistrar{}
	}

	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme

	cli, err := consulAPI.NewClient(c)
	if err != nil {
		log.Errorf("consul client init failed: %v, falling back to no-op registrar", err)
		return noOpRegistrar{}
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

// noOpRegistrar is a safe fallback that disables service registration.
type noOpRegistrar struct{}

func (noOpRegistrar) Register(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}
func (noOpRegistrar) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}
