package nacos

import (
	knacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var (
	nacosConfig      config.Config //远程配置中心实例
	nacosIp          string        //nacosip地址
	nacosPort        uint64        //nacos端口
	nacosNameSpaceId string        //当前服务名称空间
	nacosDataId      string        //nacos配置的DataID
	nacosGroup       string        //nacos配置的group分组
)

func readConfig() {
	localConfig := viper.New()                                //新建本地配置中心实例
	localConfig.SetConfigFile("..\\..\\configs\\config.yaml") //指定本地配置文件
	//读取配置文件
	if err := localConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	nacosIp = localConfig.GetString("data.nacos.addr")
	nacosPort = localConfig.GetUint64("data.nacos.port")
	nacosNameSpaceId = localConfig.GetString("data.nacos.namespaceId")
	nacosDataId = localConfig.GetString("data.nacos.dataId")
	nacosGroup = localConfig.GetString("data.nacos.groupId")
}

func InitNacosConfig() (c config.Config) {
	//初始化读取配置文件
	readConfig()
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(nacosIp, nacosPort),
	}
	cc := &constant.ClientConfig{
		NamespaceId:         nacosNameSpaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	c = config.New(
		config.WithSource(
			knacos.NewConfigSource(client,
				knacos.WithGroup(nacosGroup),
				knacos.WithDataID(nacosDataId),
			)))
	if err := c.Load(); err != nil {
		panic(err)
	}
	return c
}

func GetConfig() config.Config {
	if nacosConfig == nil {
		log.Info("准备初始化nacos读取配置")
		nacosConfig = InitNacosConfig()
	}
	return nacosConfig
}
