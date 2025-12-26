# 🚀 Lushop K8s 一键部署完整指南

> 将 lushop 微服务电商平台完整部署到 Kubernetes 集群

## 📋 部署概览

**架构**: 7个微服务 + API网关 + 基础设施服务
**配置中心**: Nacos + Consul
**监控**: Prometheus + Grafana + Jaeger
**消息队列**: RocketMQ
**存储**: MySQL + Redis + Elasticsearch

---

## 🎯 快速开始 (推荐)

### 前置要求
```bash
# 1. K8s 集群 (1.20+)
kubectl version --short

# 2. 默认 StorageClass
kubectl get storageclass

# 3. Docker 环境
docker --version
```

### 一键部署
```bash
cd /home/zzx/lucien/projects/lushop-kratos/k8s

# 🚀 一键部署所有组件
./quick-deploy.sh

# 📊 查看部署状态
./quick-deploy.sh status

# 📝 查看服务日志
./quick-deploy.sh logs gateway
```

---

## 📝 详细部署步骤 (手动执行)

### 阶段 1: 环境准备

#### 1.1 验证集群环境
```bash
# 检查集群状态
kubectl cluster-info
kubectl get nodes
kubectl get storageclass

# 要求:
# - K8s 版本: 1.20+
# - 可用 CPU: 4核以上
# - 可用内存: 8GB以上
# - StorageClass: 默认存储类存在
```

#### 1.2 设置 K8s 配置 (如需要)
```bash
# 如果 kubectl 无法连接集群
sudo cp /etc/kubernetes/admin.conf ~/.kube/config
sudo chown $(id -u):$(id -g) ~/.kube/config
chmod 600 ~/.kube/config
```

### 阶段 2: 构建和准备镜像

#### 2.1 构建服务镜像
```bash
cd /home/zzx/lucien/projects/lushop-kratos/k8s

# 构建所有镜像 (网关 + 7个微服务)
./build-images.sh all

# 验证构建结果
docker images | grep lushop

# 注意：如果看到重复镜像，请清理只保留 lushop/* 系列
docker rmi $(docker images | grep "localhost:5000/lushop" | awk '{print $1":"$2}') 2>/dev/null || true
```

#### 2.2 导入镜像到 K8s 集群
```bash
# 打包镜像
IMAGES=(lushop/user:latest lushop/goods:latest lushop/order:latest \
        lushop/inventory:latest lushop/userop:latest lushop/userauth:latest \
        lushop/gateway:latest)
docker save "${IMAGES[@]}" -o /tmp/lushop-images.tar

# 导入到 containerd
sudo ctr -n k8s.io images import /tmp/lushop-images.tar

# 验证导入
sudo ctr -n k8s.io images ls | grep lushop
```

### 阶段 3: 配置管理

#### 3.1 生成密码和 Secrets
```bash
# 生成自定义密码的 K8s Secrets
./gen-secrets-custom.sh

# 验证 Secrets 创建
kubectl get secrets -n lushop
```

**脚本功能：**
- 生成 MySQL、Redis、Nacos、RocketMQ、Grafana、Elasticsearch 的密码
- 支持通过环境变量自定义密码
- 自动创建 K8s Secrets
- 显示生成的密码摘要

#### 3.2 创建命名空间
```bash
kubectl create namespace lushop --dry-run=client -o yaml | kubectl apply -f -
```

### 阶段 4: 部署基础设施

#### 4.1 部署存储服务
```bash
# MySQL 数据库
kubectl apply -k base/mysql
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# Redis 缓存
kubectl apply -k base/redis
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=redis -n lushop --timeout=300s

# RocketMQ 消息队列
kubectl apply -k base/rocketmq
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=rocketmq -n lushop --timeout=300s
```

#### 4.2 部署配置和注册中心
```bash
# Nacos 配置中心
kubectl apply -k base/nacos
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=nacos -n lushop --timeout=300s

# Consul 服务发现
kubectl apply -k base/consul
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=consul -n lushop --timeout=300s
```

#### 4.3 部署监控和日志
```bash
# Prometheus 监控
kubectl apply -k base/prometheus

# Grafana 可视化
kubectl apply -k base/grafana

# Jaeger 链路追踪
kubectl apply -k base/jaeger

# Elasticsearch 日志存储
kubectl apply -k base/elasticsearch
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=elasticsearch -n lushop --timeout=300s

# Kibana 日志可视化
kubectl apply -k base/kibana
```

### 阶段 5: 配置 Nacos

#### 5.1 导入服务配置
```bash
# 自动导入所有服务配置到 Nacos
./configure-nacos.sh import

# 验证配置导入
./configure-nacos.sh verify
```

#### 5.2 手动验证 (可选)
```bash
# 端口转发访问 Nacos 控制台
kubectl port-forward -n lushop svc/nacos 8848:8848 &

# 浏览器访问: http://localhost:8848/nacos
# 账号: nacos / nacos
# 命名空间: lushop (ID: de9c6a0e-1fbc-425d-8d3b-09066fea6889)
```

### 阶段 6: 部署业务服务

#### 6.1 按依赖顺序部署微服务
```bash
# 基础服务 (无依赖)
kubectl apply -k base/services/user
kubectl apply -k base/services/goods
kubectl apply -k base/services/inventory
kubectl apply -k base/services/userop
kubectl apply -k base/services/userauth

# 等待基础服务就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=user -n lushop --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=goods -n lushop --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=inventory -n lushop --timeout=300s

# 依赖服务
kubectl apply -k base/services/order

# API 网关 (最后部署)
kubectl apply -k base/services/gateway
```

#### 6.2 等待所有服务就绪
```bash
# 等待网关就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=gateway -n lushop --timeout=300s

# 查看所有服务状态
kubectl get pods -n lushop
kubectl get svc -n lushop
```

### 阶段 7: 验证部署

#### 7.1 检查服务健康状态
```bash
# 查看所有 Pod 状态
kubectl get pods -n lushop -o wide

# 查看服务端点
kubectl get svc -n lushop

# 检查部署状态
kubectl get deployments -n lushop
```

#### 7.2 测试 API 访问
```bash
# 获取节点 IP
NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')

# 测试网关健康检查
curl http://$NODE_IP:30080/health

# 测试商品服务
curl http://$NODE_IP:30080/api/goods/list
```

#### 7.3 验证服务间通信
```bash
# 检查服务日志
kubectl logs -f deployment/gateway -n lushop

# 验证数据库连接
kubectl exec -it deployment/user-service -n lushop -- mysql -h mysql -u lushop -p -e "SHOW DATABASES;"

# 验证 Redis 连接
kubectl exec -it deployment/user-service -n lushop -- redis-cli -h redis -a <password> ping
```

### 阶段 8: 访问监控组件

#### 8.1 设置端口转发
```bash
# Prometheus (监控指标)
kubectl port-forward -n lushop svc/prometheus 9090:9090 &

# Grafana (可视化面板) - admin/admin
kubectl port-forward -n lushop svc/grafana 3000:3000 &

# Jaeger (链路追踪)
kubectl port-forward -n lushop svc/jaeger 16686:16686 &

# Kibana (日志分析)
kubectl port-forward -n lushop svc/kibana 5601:5601 &

# Nacos (配置管理) - nacos/nacos
kubectl port-forward -n lushop svc/nacos 8848:8848 &

# Consul (服务发现)
kubectl port-forward -n lushop svc/consul 8500:8500 &
```

#### 8.2 访问地址
- **API 网关**: `http://<node-ip>:30080`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)
- **Jaeger**: `http://localhost:16686`
- **Kibana**: `http://localhost:5601`
- **Nacos**: `http://localhost:8848` (nacos/nacos)
- **Consul**: `http://localhost:8500`

---

## 🛠️ 故障排查

### 常见问题解决方案

#### 1. 镜像拉取失败
```bash
# 检查镜像是否存在
sudo ctr -n k8s.io images ls | grep lushop

# 重新导入镜像
docker save lushop/user:latest -o /tmp/user.tar
sudo ctr -n k8s.io images import /tmp/user.tar
```

#### 2. Pod 启动失败
```bash
# 查看 Pod 日志
kubectl logs -f pod/<pod-name> -n lushop

# 查看 Pod 详细信息
kubectl describe pod/<pod-name> -n lushop

# 检查事件
kubectl get events -n lushop --sort-by='.lastTimestamp'
```

#### 3. 配置问题
```bash
# 验证 Nacos 配置
kubectl port-forward -n lushop svc/nacos 8848:8848 &
curl "http://localhost:8848/nacos/v1/cs/configs?dataId=user.yaml&group=lushop_grpc"

# 重新导入配置
./configure-nacos.sh import
```

#### 4. 服务间通信失败
```bash
# 测试服务发现
kubectl exec -it deployment/user-service -n lushop -- nslookup goods-service

# 测试端口连通性
kubectl exec -it deployment/user-service -n lushop -- nc -zv goods-service 8012
```

#### 5. 数据库连接失败
```bash
# 检查 MySQL Pod
kubectl logs -f deployment/mysql -n lushop

# 测试数据库连接
kubectl exec -it deployment/mysql -n lushop -- mysql -u root -p -e "SHOW DATABASES;"

# 验证密码
kubectl get secret mysql-auth -n lushop -o yaml
```

---

## 📊 服务架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                    Lushop API Gateway                            │
│              HTTP: 30080 | gRPC: 30090                           │
│    ┌──────────────────────────────────────────────────┐         │
│    │  User | UserAuth | Goods | Order | Inventory |    │         │
│    │  UserOp (7 Services)                              │         │
│    └──────────────────────────────────────────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│    MySQL      │    │    Redis      │    │   RocketMQ    │
│  Port: 3306   │    │  Port: 6379   │    │  Port: 9876   │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│    Nacos      │    │   Consul      │    │  Prometheus   │
│ Port: 8848    │    │ Port: 8500    │    │ Port: 9090    │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        └─────────────────────┴─────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
                    ▼                   ▼
            ┌──────────────┐    ┌──────────────┐
            │  Elasticsearch │    │   Jaeger      │
            │  Port: 9200    │    │ Port: 16686   │
            └───────────────┘    └───────────────┘
```

---

## 🎯 服务端口映射

| 服务 | HTTP端口 | gRPC端口 | NodePort | 备注 |
|------|----------|----------|----------|------|
| Gateway | 8001 | 9001 | 30080/30090 | API网关 |
| User | 8011 | 50051 | - | 用户服务 |
| Goods | 8012 | 50052 | - | 商品服务 |
| Order | 8013 | 50053 | - | 订单服务 |
| Inventory | 8014 | 50054 | - | 库存服务 |
| UserOp | 8015 | 50055 | - | 用户操作服务 |
| UserAuth | 8016 | 50056 | - | 认证服务 |

---

## 📚 相关文档

- [详细部署指南](k8s/README-CN.md)
- [配置协调机制](k8s/config-coordination.md)
- [配置使用说明](k8s/config-usage.md)
- [Nacos配置导入](k8s/configure-nacos.sh)
- [一键部署脚本](k8s/quick-deploy.sh)

---

## 🎉 部署完成！

部署完成后，你将拥有一个完整的生产级微服务电商平台，支持：

- ✅ **7个微服务** + API网关
- ✅ **服务注册发现** (Consul)
- ✅ **配置中心** (Nacos)
- ✅ **分布式链路追踪** (Jaeger)
- ✅ **监控告警** (Prometheus + Grafana)
- ✅ **日志收集** (Elasticsearch + Kibana)
- ✅ **消息队列** (RocketMQ)

**访问你的微服务平台**: `http://<node-ip>:30080` 🚀
