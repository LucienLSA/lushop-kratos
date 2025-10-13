package main

import (
	"flag"
	nacosconfig "lushop/internal/conf/nacos"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// go build -ldflags "-X main.Version=x.y.z" 定义版本服务名称等
var (
	// Name is the name of the compiled software.
	Name = "lushop.api"
	// Version is the version of the compiled software.
	Version = "lushop.api.v1"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

// 启动指定路径配置文件
func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// wire依赖注入实现的方法，整合了应用所需的依赖，应用实例创建
func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, rr registry.Registrar) *kratos.App {
	// 设置应用的ID、名称、版本、元数据、日志、服务器和注册中心
	return kratos.New(
		kratos.ID(id+"lushop.api"),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(rr),
	)
}

func setTracerProvider(url string) error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	// 创建配置加载器
	configLoader := nacosconfig.NewNacosConfigLoader(logger)
	// 首先加载本地配置获取 Nacos 连接信息
	localConfig, err := configLoader.CreateLocalConfig(flagconf)
	if err != nil {
		panic(err)
	}
	defer localConfig.Close()
	// 获取初始配置
	bc, err := configLoader.LoadBootstrapConfig(localConfig)
	if err != nil {
		panic(err)
	}
	// 使用 Nacos 或本地配置加载完整配置
	c, err := configLoader.LoadConfigWithNacos(bc, flagconf)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	// 重新加载配置（可能来自 Nacos）
	bc, err = configLoader.LoadBootstrapConfig(c)
	if err != nil {
		panic(err)
	}
	// 加载注册中心配置
	rc, err := configLoader.LoadRegistryConfig(c)
	if err != nil {
		panic(err)
	}
	// 初始化链路追踪
	err = setTracerProvider(bc.Trace.Endpoint)
	if err != nil {
		panic(err)
	}

	// 通过 wireApp 函数（由 Wire 生成）构建应用实例，并获取清理函数
	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Auth, bc.Service, bc.Sms, rc, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
