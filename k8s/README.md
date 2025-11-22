# Kubernetes 部署指引

该目录提供了将 `lushop` 项目部署到 Kubernetes 集群的完整清单，支持单机部署和生产环境部署。

## 📁 目录结构

```
k8s/
├── base/                    # 基础清单，可直接 kubectl apply
│   ├── namespace.yaml       # 命名空间
│   ├── redis/              # Redis 配置
│   ├── mysql/              # MySQL 配置
│   ├── nacos/              # Nacos 配置中心
│   ├── consul/             # Consul 服务注册发现
│   ├── jaeger/             # Jaeger 链路追踪
│   ├── rocketmq/           # RocketMQ 消息队列
│   ├── prometheus/         # Prometheus 监控
│   ├── grafana/            # Grafana 可视化
│   └── services/           # 业务服务
│       ├── user/           # 用户服务
│       ├── goods/          # 商品服务
│       ├── order/          # 订单服务
│       ├── inventory/      # 库存服务
│       ├── userop/         # 用户操作服务
│       ├── userauth/       # 认证服务
│       └── gateway/        # API 网关
├── overlays/                # 环境覆盖配置（可选）
│   └── dev/                # 开发环境示例
├── deploy.sh               # 一键部署脚本
└── DEPLOY.md               # 详细部署文档
```

## 🚀 快速开始

### 前置要求

- Kubernetes 集群（k3s、minikube、kind 或标准 k8s）
- kubectl 已配置
- Docker 已安装并运行（用于构建镜像）
- 服务镜像已构建（或使用预构建镜像）

### 构建镜像

在部署之前，需要先构建服务镜像：

```bash
cd k8s

# 构建所有镜像（推荐）
./build-images.sh all

# 或仅构建服务镜像
./build-images.sh services

# 或仅构建网关镜像
./build-images.sh gateway

# 或构建指定服务
./build-images.sh user
./build-images.sh goods

# 查看已构建的镜像
./build-images.sh list
```

**注意**: 如果遇到 Docker iptables 错误，请参考 [故障排查](#-故障排查) 部分。

### 一键部署

```bash
cd k8s

# 部署所有服务
./deploy.sh deploy

# 查看状态
./deploy.sh status

# 查看日志
./deploy.sh logs gateway-service
```

### 手动部署

```bash
# 部署所有资源
kubectl apply -k base/

# 或分步部署
kubectl apply -k base/redis
kubectl apply -k base/mysql
kubectl apply -k base/nacos
kubectl apply -k base/consul
kubectl apply -k base/jaeger
kubectl apply -k base/services/
```

## 📋 部署步骤（详细）

### 1. 准备命名空间与存储

- 所有清单默认部署到 `lushop` 命名空间
- Redis、MySQL、RocketMQ、Prometheus 的 PVC 默认使用集群 `default` StorageClass
- 如需指定 StorageClass（例如 `local-path`、`rook-ceph`），请在各 YAML 中显式设置 `storageClassName`

### 2. 准备 Secret 与配置

**重要**: 仓库中的 Secret 仅为演示值，部署前务必替换！

需要替换的 Secret：
- `redis-auth`（`k8s/base/redis/secret.yaml`）：Redis 密码
- `mysql-auth`（`k8s/base/mysql/secret.yaml`）：MySQL root/业务账号
- `nacos-auth`（`k8s/base/nacos/secret.yaml`）：Nacos MySQL 连接信息
- `rocketmq-credentials`（`k8s/base/rocketmq/secret.yaml`）：RocketMQ ACL 凭据
- `grafana-admin`（`k8s/base/grafana/secret.yaml`）：Grafana 管理员账号

### 3. 初始化数据库

部署 MySQL 后，需要初始化数据库：

```bash
# 等待 MySQL 就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# 导入数据库脚本（使用 port-forward）
kubectl port-forward -n lushop svc/mysql 3306:3306
mysql -h 127.0.0.1 -uroot -proot123456 < scripts/init_db.sql
```

### 4. 配置 Nacos

1. 等待 Nacos 就绪
2. 访问 Nacos 控制台: http://localhost:8848/nacos (nacos/nacos)
3. 创建命名空间: `de9c6a0e-1fbc-425d-8d3b-09066fea6889`
4. 为每个服务创建配置（参考 `service/*/configs/nacos-config.yaml`）

### 5. 应用基础设施组件

```bash
kubectl apply -k k8s/base
```

使用 `kubectl get pods -n lushop` 观察所有组件进入 `Ready` 状态。

### 6. 部署业务服务

业务服务配置在 `k8s/base/services/` 目录下，通过 ConfigMap 注入配置。

### 7. 使用 overlays（可选）

`k8s/overlays/dev` 为示例，可复制成 `stage`、`prod`。利用 overlay patches 覆盖镜像 tag、副本数、资源限制等。

```bash
kustomize build k8s/overlays/dev | kubectl apply -f -
```

## 🔧 配置说明

### 服务端口映射

| 服务 | HTTP 端口 | gRPC 端口 | 说明 |
|------|-----------|-----------|------|
| Gateway | 8001 | 9001 | API 网关（NodePort: 30080/30090） |
| User | 8011 | 50051 | 用户服务 |
| Goods | 8012 | 50052 | 商品服务 |
| Order | 8013 | 50053 | 订单服务 |
| Inventory | 8014 | 50054 | 库存服务 |
| UserOp | 8015 | 50055 | 用户操作服务 |
| UserAuth | 8016 | 50056 | 认证服务 |

### 资源配置

**单机测试环境**（默认）:
- 请求: 128Mi 内存, 100m CPU
- 限制: 512Mi 内存, 500m CPU
- 副本数: 1

**生产环境建议**:
- 根据实际负载调整资源限制
- 增加副本数实现高可用
- 配置 HPA/VPA 自动扩缩容

### 存储配置

- MySQL: 20Gi PVC
- Redis: 5Gi PVC
- RocketMQ: 10Gi PVC
- Prometheus: 10Gi PVC

## 🌐 访问服务

### NodePort（已配置）

Gateway 服务已配置为 NodePort：

```bash
# 获取节点 IP
NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')

# 访问 API
curl http://$NODE_IP:30080/api/goods/list
```

### Port Forward

```bash
# 转发 Gateway 服务
kubectl port-forward -n lushop svc/gateway-service 8001:8001

# 访问 API
curl http://localhost:8001/api/goods/list
```

### Ingress（需要安装 Ingress Controller）

参考 `DEPLOY.md` 中的 Ingress 配置示例。

## 📊 监控与可观测

- **Prometheus**: 服务监控指标收集
- **Grafana**: 可视化仪表盘
- **Jaeger**: 分布式链路追踪（端口 16686）

访问方式：
```bash
# Prometheus
kubectl port-forward -n lushop svc/prometheus 9090:9090

# Grafana
kubectl port-forward -n lushop svc/grafana 3000:3000

# Jaeger
kubectl port-forward -n lushop svc/jaeger 16686:16686
```

## 🐛 故障排查

### Docker 构建问题

**问题**: `iptables: no such file or directory` 或网络相关错误

**解决方案**:
```bash
# 1. 检查 iptables 是否安装
sudo apt-get install iptables -y  # Ubuntu/Debian
sudo yum install iptables -y      # CentOS/RHEL

# 2. 重启 Docker 服务
sudo systemctl restart docker

# 3. 检查 Docker 网络
docker network ls
docker network inspect bridge

# 4. 如果使用 rootless Docker，可能需要配置
# 参考: https://docs.docker.com/engine/security/rootless/
```

**问题**: 构建时找不到目录

**解决方案**:
- 确保在项目根目录执行构建脚本
- 使用提供的 `build-images.sh` 脚本，它会自动处理路径问题

### Pod 无法启动

```bash
# 查看 Pod 状态
kubectl describe pod <pod-name> -n lushop

# 查看日志
kubectl logs <pod-name> -n lushop

# 查看事件
kubectl get events -n lushop --sort-by='.lastTimestamp'
```

### 服务无法连接

```bash
# 检查服务端点
kubectl get endpoints -n lushop

# 测试 DNS 解析
kubectl run -it --rm debug --image=busybox --restart=Never -- nslookup user-service.lushop.svc.cluster.local

# 测试服务连通性
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- curl http://user-service.lushop.svc.cluster.local:8011/health
```

### 配置问题

```bash
# 查看 ConfigMap
kubectl get configmap -n lushop
kubectl describe configmap <configmap-name> -n lushop

# 查看 Secret
kubectl get secret -n lushop
kubectl describe secret <secret-name> -n lushop
```

### 镜像拉取失败

```bash
# 检查镜像是否存在
docker images | grep lushop

# 如果使用私有仓库，检查镜像拉取 Secret
kubectl get secret -n lushop | grep docker

# 如果镜像在本地，需要导入到集群
# 对于 kind: kind load docker-image lushop/user:latest
# 对于 minikube: minikube image load lushop/user:latest
# 对于 k3s: 使用 k3d 或直接使用本地镜像
```

## 📝 后续优化建议

1. **高可用**: 增加副本数，配置 HPA/VPA、PDB
2. **监控**: 配置 ServiceMonitor，完善 Prometheus 监控
3. **日志**: 集成 ELK 或 Loki 进行日志收集
4. **安全**: 配置 NetworkPolicy，使用 TLS 加密，SealedSecret 管理敏感信息
5. **备份**: MySQL、RocketMQ Store、Prometheus 等数据卷需制定备份策略（Velero、快照、CronJob）
6. **CI/CD**: 集成 CI/CD 流水线自动构建和部署
7. **持续扩展**: 补充 Elasticsearch/Kibana、完善业务服务配置

## 📚 相关文档

- [详细部署文档](DEPLOY.md) - 完整的部署步骤和配置说明
- [项目主 README](../README.md) - 项目架构和功能介绍

## 🔗 快速命令参考

### 镜像构建

```bash
# 构建所有镜像
./build-images.sh all

# 构建指定服务
./build-images.sh user
./build-images.sh gateway

# 查看镜像列表
./build-images.sh list
```

### 服务部署

```bash
# 部署所有服务
./deploy.sh deploy

# 删除所有服务
./deploy.sh delete

# 查看状态
./deploy.sh status

# 查看日志
./deploy.sh logs [service-name]

# 仅部署基础设施
./deploy.sh infrastructure

# 仅部署业务服务
./deploy.sh services
```

