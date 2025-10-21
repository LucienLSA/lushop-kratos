package config

import (
	"userop/internal/conf"

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
// 优先使用 Nacos 配置中心，如果 Nacos 不可用则使用本地配置作为备用
func (loader *NacosConfigLoader) LoadConfigWithNacos(bc *conf.Bootstrap, flagconf string) (config.Config, error) {
	// 如果配置了 Nacos，则优先使用 Nacos 作为配置源
	if bc.Nacos != nil {
		loader.logger.Log(log.LevelInfo, "msg", "Attempting to load config from Nacos")

		// 创建 Nacos 客户端
		clientConfig := constant.ClientConfig{
			NamespaceId:         bc.Nacos.NamespaceId,
			TimeoutMs:           10000, // 增加超时时间
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			LogLevel:            "info",
			Username:            "nacos", // 添加用户名
			Password:            "nacos", // 添加密码
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
			loader.logger.Log(log.LevelWarn, "msg", "Failed to create Nacos client, falling back to local config", "error", err)
			return loader.CreateLocalConfig(flagconf)
		}

		// 创建 Nacos 配置源
		nacosSource := nacosconfig.NewConfigSource(
			configClient,
			nacosconfig.WithDataID(bc.Nacos.DataId),
			nacosconfig.WithGroup(bc.Nacos.GroupId),
		)

		// 首先尝试只使用 Nacos 配置源
		c := config.New(
			config.WithSource(nacosSource),
		)

		if err := c.Load(); err != nil {
			loader.logger.Log(log.LevelWarn, "msg", "Failed to load config from Nacos, falling back to local config", "error", err)
			return loader.CreateLocalConfig(flagconf)
		}

		// 验证 Nacos 配置是否包含必要的数据
		var testConfig conf.Bootstrap
		if err := c.Scan(&testConfig); err != nil {
			loader.logger.Log(log.LevelWarn, "msg", "Nacos config is invalid, falling back to local config", "error", err)
			return loader.CreateLocalConfig(flagconf)
		}

		// 检查关键配置是否存在
		if testConfig.Data == nil || testConfig.Server == nil {
			loader.logger.Log(log.LevelWarn, "msg", "Nacos config missing required fields, falling back to local config")
			return loader.CreateLocalConfig(flagconf)
		}

		loader.logger.Log(log.LevelInfo, "msg", "Successfully loaded config from Nacos")
		return c, nil
	}

	// 如果没有配置 Nacos，使用本地配置
	loader.logger.Log(log.LevelInfo, "msg", "No Nacos config found, using local config")
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

// LoadConfigWithNacosPriority 使用 Nacos 优先的配置加载策略
// 优先使用 Nacos 配置，如果 Nacos 配置不完整，则用本地配置补充
func (loader *NacosConfigLoader) LoadConfigWithNacosPriority(bc *conf.Bootstrap, flagconf string) (config.Config, error) {
	// 如果配置了 Nacos，则优先使用 Nacos 作为配置源
	if bc.Nacos != nil {
		loader.logger.Log(log.LevelInfo, "msg", "Attempting to load config from Nacos with priority")

		// 创建 Nacos 客户端
		clientConfig := constant.ClientConfig{
			NamespaceId:         bc.Nacos.NamespaceId,
			TimeoutMs:           10000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			LogLevel:            "info",
			Username:            "nacos", // 添加用户名
			Password:            "nacos", // 添加密码
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

		// 创建配置，Nacos 优先，本地配置作为补充
		c := config.New(
			config.WithSource(
				nacosSource,              // Nacos 配置优先
				file.NewSource(flagconf), // 本地配置作为补充
			),
		)

		if err := c.Load(); err != nil {
			loader.logger.Log(log.LevelWarn, "msg", "Failed to load config, using local config only", "error", err)
			return loader.CreateLocalConfig(flagconf)
		}

		loader.logger.Log(log.LevelInfo, "msg", "Successfully loaded config with Nacos priority")
		return c, nil
	}

	// 如果没有配置 Nacos，使用本地配置
	loader.logger.Log(log.LevelInfo, "msg", "No Nacos config found, using local config")
	return loader.CreateLocalConfig(flagconf)
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
