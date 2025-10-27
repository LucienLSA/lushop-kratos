package main

import (
	"flag"
	"os"

	"order/internal/conf/metrix"
	nacosconfig "order/internal/conf/nacos"
	"order/internal/pkg/snowflake"

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
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "lushop.order.service"
	// Version is the version of the compiled software.
	Version = "order.v1"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, rr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(id+"order_service"),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(rr), // consul 的引入 服务发现和注册
	)
}

// Set global trace provider 设置链路追逐的方法
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

	// 初始化雪花算法（使用节点ID 1，生产环境应该从配置文件读取）
	if err := snowflake.Init(1); err != nil {
		logger.Log(log.LevelError, "msg", "failed to initialize snowflake", "error", err)
		panic(err)
	}
	logger.Log(log.LevelInfo, "msg", "snowflake initialized successfully")

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
	// 使用 Nacos 优先的配置加载策略
	c, err := configLoader.LoadConfigWithNacosPriority(bc, flagconf)
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
	// 设置链路追踪
	if bc.Trace != nil {
		err = setTracerProvider(bc.Trace.Endpoint)
		if err != nil {
			panic(err)
		}
	} else {
		logger.Log(log.LevelWarn, "msg", "Trace configuration not found, skipping tracer setup")
	}
	app, cleanup, err := wireApp(bc.Server, bc.Data, rc, bc.Service, bc, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	// servicename, err := nacos.GetConfig().Value("service.name").String()
	// if err != nil {
	// 	panic(err)
	// }
	// log.Debugf("servicename:%s", servicename)
	// 初始化 Prometheus metrics
	metrix.Init()
	logger.Log(log.LevelInfo, "msg", "Prometheus metrics initialized")

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
