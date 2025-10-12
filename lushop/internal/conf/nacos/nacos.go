package config

import (
	"lushop/internal/conf"

	nacosconfig "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NacosConfigLoader Nacos 配置加载器
type NacosConfigLoader struct {
	logger log.Logger
}

// NewNacosConfigLoader 创建 Nacos 配置加载器
func NewNacosConfigLoader(logger log.Logger) *NacosConfigLoader {
	return &NacosConfigLoader{
		logger: logger,
	}
}

// LoadConfigWithNacos 使用 Nacos 加载配置
// 如果 Nacos 可用，则从 Nacos 加载配置；否则使用本地配置
func (loader *NacosConfigLoader) LoadConfigWithNacos(bc *conf.Bootstrap, flagconf string) (config.Config, error) {
	// 如果配置了 Nacos，则使用 Nacos 作为配置源
	if bc.Nacos != nil {
		loader.logger.Log(log.LevelInfo, "msg", "Using Nacos as config source")

		// 创建 Nacos 客户端
		clientConfig := constant.ClientConfig{
			NamespaceId:         bc.Nacos.NamespaceId,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			LogLevel:            "info",
			Username:            "nacos",
			Password:            "nacos",
		}

		serverConfigs := []constant.ServerConfig{
			{
				IpAddr: bc.Nacos.Addr,
				Port:   uint64(bc.Nacos.Port),
			},
		}

		// 创建配置客户端
		configClient, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			loader.logger.Log(log.LevelWarn, "msg", "Failed to create Nacos client, using local config", "error", err)
			return loader.CreateLocalConfig(flagconf)
		}

		// 创建 Nacos 配置源
		nacosSource := nacosconfig.NewConfigSource(
			configClient,
			nacosconfig.WithDataID(bc.Nacos.DataId),
			nacosconfig.WithGroup(bc.Nacos.GroupId),
		)

		// 创建配置，使用 Nacos 作为主要配置源
		c := config.New(
			config.WithSource(
				nacosSource,
				file.NewSource(flagconf), // 本地配置作为备用
			),
		)

		if err := c.Load(); err != nil {
			loader.logger.Log(log.LevelWarn, "msg", "Failed to load Nacos config, using local config", "error", err)
			return loader.CreateLocalConfig(flagconf)
		}

		loader.logger.Log(log.LevelInfo, "msg", "Successfully loaded config from Nacos")
		return c, nil
	}

	// 如果没有配置 Nacos，使用本地配置
	loader.logger.Log(log.LevelInfo, "msg", "Using local config")
	return loader.CreateLocalConfig(flagconf)
}

// CreateLocalConfig 创建本地配置
func (loader *NacosConfigLoader) CreateLocalConfig(flagconf string) (config.Config, error) {
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)

	if err := c.Load(); err != nil {
		return nil, err
	}

	return c, nil
}

// LoadBootstrapConfig 加载 Bootstrap 配置
func (loader *NacosConfigLoader) LoadBootstrapConfig(c config.Config) (*conf.Bootstrap, error) {
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}
	return &bc, nil
}

// LoadRegistryConfig 加载注册中心配置
// 这里引入consul
func (loader *NacosConfigLoader) LoadRegistryConfig(c config.Config) (*conf.Registry, error) {
	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}
