# Prometheus 监控配置说明

## 概述

已为 lushop 项目配置完整的 Prometheus + Grafana 监控系统。

## 架构

- **Prometheus**: 指标收集和存储 (端口: 9090)
- **Grafana**: 数据可视化 (端口: 3000)
- **Goods Service**: HTTP metrics 端点 (端口: 8000)

## 配置内容

### 1. Goods Service 配置

- ✅ 初始化 Prometheus exporter (`metrix.Init()`)
- ✅ 注册 metrics 中间件 (请求计数和耗时)
- ✅ 创建 HTTP 服务器暴露 `/metrics` 端点
- ✅ 配置端口 8000 用于 metrics

### 2. Prometheus 配置

- ✅ 配置文件: `deploy/prometheus/prometheus.yml`
- ✅ 监控所有微服务 (goods, user, order, inventory, userop, userauth)
- ✅ 抓取间隔: 15秒

### 3. Docker Compose 配置

- ✅ Prometheus 服务 (端口 9090)
- ✅ Grafana 服务 (端口 3000)
- ✅ Goods Service 暴露 8000 端口

## 使用方法

### 启动服务

```bash
# 启动所有服务（包括 Prometheus 和 Grafana）
docker-compose up -d

# 仅启动监控相关服务
docker-compose up -d prometheus grafana

# 启动 goods 服务
docker-compose up -d goods-service
```

### 访问监控界面

1. **Prometheus UI**
   - URL: http://localhost:9090
   - 查看指标: http://localhost:9090/targets
   - 查询示例:
     ```promql
     # 请求总数
     http_server_requests_total{service="goods-service"}
     
     # 请求耗时
     http_server_request_duration_seconds{service="goods-service"}
     ```

2. **Grafana**
   - URL: http://localhost:3000
   - 默认账号: admin
   - 默认密码: admin

3. **Goods Service Metrics**
   - URL: http://localhost:8000/metrics
   - 直接查看原始指标数据

### 配置 Grafana

1. 登录 Grafana (http://localhost:3000)
2. 添加数据源:
   - Configuration → Data Sources → Add data source
   - 选择 Prometheus
   - URL: http://prometheus:9090
   - 点击 "Save & Test"

3. 导入仪表板:
   - Create → Import
   - 输入仪表板 ID: 
     - `6417` (Kratos 服务监控)
     - `1860` (Node Exporter)
   - 选择 Prometheus 数据源
   - 点击 Import

## 监控指标

### Goods Service 指标

- `http_server_requests_total`: HTTP 请求总数
- `http_server_request_duration_seconds`: HTTP 请求耗时直方图
- `grpc_server_requests_total`: gRPC 请求总数
- `grpc_server_request_duration_seconds`: gRPC 请求耗时直方图

### 标签

- `service`: 服务名称 (如: goods-service)
- `method`: 请求方法
- `code`: 响应状态码

## 本地开发验证

如果在本地开发环境运行（非 Docker）:

```bash
# 1. 启动 goods 服务
cd service/goods
go run cmd/goods/main.go

# 2. 验证 metrics 端点
curl http://localhost:8000/metrics

# 3. 查看 Prometheus 是否能抓取数据
# 访问 http://localhost:9090/targets
```

## 故障排查

### Goods Service 无法启动

检查配置文件 `configs/nacos-config.yaml`:
```yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:50052
    timeout: 1s
```

### Prometheus 无法抓取数据

1. 检查 Prometheus targets 状态: http://localhost:9090/targets
2. 确认 goods-service 的 8000 端口可访问
3. 查看 Prometheus 日志:
   ```bash
   docker logs lushop-prometheus
   ```

### Metrics 端点返回 404

确认以下内容:
1. `metrix.Init()` 已在 main.go 中调用
2. HTTP 服务器已正确创建和注册
3. Wire 依赖注入已重新生成: `wire gen`

## 扩展配置

### 添加新的监控指标

在 `internal/conf/metrix/metrix.go` 中添加:

```go
var CustomMetric metric.Int64Counter

func Init() {
    // ... 现有代码 ...
    
    CustomMetric, err = meter.Int64Counter("custom_metric_name")
    if err != nil {
        panic(err)
    }
}
```

### 添加告警规则

在 `prometheus.yml` 中添加:

```yaml
rule_files:
  - 'alerts.yml'
```

创建 `alerts.yml`:

```yaml
groups:
  - name: goods_service
    rules:
      - alert: HighErrorRate
        expr: rate(http_server_requests_total{code=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate on {{ $labels.service }}"
```

## 参考文档

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)
- [Kratos Metrics 中间件](https://go-kratos.dev/docs/component/middleware/metrics)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
