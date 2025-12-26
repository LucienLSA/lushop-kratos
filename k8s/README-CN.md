# Lushop K8s 完整部署指南

## 📋 项目概述

lushop-kratos 是基于 Go-Kratos 框架的微服务电商平台，包含以下组件：

### 🏗️ 架构组件
- **API Gateway**: HTTP/gRPC 网关服务
- **7个微服务**: User, Goods, Order, Inventory, UserOp, UserAuth
- **基础设施**: MySQL, Redis, Nacos, Consul, Elasticsearch, Jaeger, RocketMQ, Prometheus

### 🎯 部署目标
将完整的微服务架构部署到 Kubernetes 集群，实现生产级别的服务治理、配置管理和可观测性。

## 🚀 完整部署流程

### 阶段 1: 环境准备

#### 1.1 K8s 集群要求
```bash
# 检查集群状态
kubectl cluster-info
kubectl get nodes
kubectl get storageclass

# 要求：
# - K8s 1.20+
# - 默认 StorageClass (如 local-path)
# - 至少 4核8G内存 (推荐8核16G)
```

#### 1.2 本地环境准备
```bash
# 进入项目目录
cd /home/zzx/lucien/projects/lushop-kratos

# 确保 Docker 运行
docker --version && docker ps
```

### 阶段 2: 构建与导入镜像

#### 2.1 构建所有服务镜像
```bash
# 进入 k8s 目录
cd k8s

chmod a+x build-images.sh

# 构建所有镜像 (网关 + 7个微服务)
./build-images.sh all

# 验证构建结果
docker images | grep lushop
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

#### 3.1 生成 K8s Secrets
```bash
# 生成数据库和服务的密码
./gen-secrets-custom.sh

# 查看生成的 secrets
kubectl get secrets -n lushop
```

#### 3.2 创建 K8s 友好的配置文件

由于原始配置文件使用本地地址，需要创建适合 K8s 的配置：

**主要修改内容：**
- MySQL: `127.0.0.1:3306` → `mysql:3306`
- Redis: `127.0.0.1:6379` → `redis:6379`
- Consul: `127.0.0.1:8500` → `consul:8500`
- Elasticsearch: `192.168.x.x:9200` → `elasticsearch:9200`
- Jaeger: `192.168.x.x:14268` → `jaeger:14268`
- RocketMQ: `127.0.0.1:9876` → `rocketmq:9876`

**已创建的配置文件：**
```
service/user/configs/nacos-config-k8s.yaml
service/goods/configs/nacos-config-k8s.yaml
service/order/configs/nacos-config-k8s.yaml
service/inventory/configs/nacos-config-k8s.yaml
service/userop/configs/nacos-config-k8s.yaml
service/userauth/configs/nacosRemote-k8s.yaml
lushop/configs/nacos-config-k8s.yaml
```

### 阶段 4: 部署基础设施

#### 4.1 创建命名空间
```bash
kubectl create namespace lushop
```

#### 4.2 部署存储层
```bash
# MySQL
kubectl apply -k base/mysql
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# Redis
kubectl apply -k base/redis
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=redis -n lushop --timeout=300s

# RocketMQ
kubectl apply -k base/rocketmq
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=rocketmq -n lushop --timeout=300s
```

#### 4.3 部署服务治理
```bash
# Nacos (配置中心)
kubectl apply -k base/nacos
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=nacos -n lushop --timeout=300s

# Consul (服务发现)
kubectl apply -k base/consul
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=consul -n lushop --timeout=300s

# 监控组件
kubectl apply -k base/prometheus
kubectl apply -k base/grafana
kubectl apply -k base/jaeger

# 日志组件
kubectl apply -k base/elasticsearch
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=elasticsearch -n lushop --timeout=300s
kubectl apply -k base/kibana
```

### 阶段 5: 配置 Nacos

#### 5.1 访问 Nacos 控制台
```bash
# 端口转发
kubectl port-forward -n lushop svc/nacos 8848:8848 &

# 访问地址: http://localhost:8848/nacos
# 默认账号: nacos / nacos
```

#### 5.2 创建命名空间
在 Nacos 控制台：
1. 点击 **命名空间** 选项卡
2. 点击 **新建命名空间**
3. 命名空间 ID: `de9c6a0e-1fbc-425d-8d3b-09066fea6889`
4. 命名空间名称: `lushop`

#### 5.3 导入服务配置

为每个服务创建配置（Group: `lushop_grpc`）：

| 服务 | Data ID | 配置文件 |
|------|---------|----------|
| user | `user.yaml` | `service/user/configs/nacos-config-k8s.yaml` |
| goods | `goods.yaml` | `service/goods/configs/nacos-config-k8s.yaml` |
| order | `order.yaml` | `service/order/configs/nacos-config-k8s.yaml` |
| inventory | `inventory.yaml` | `service/inventory/configs/nacos-config-k8s.yaml` |
| userop | `userop.yaml` | `service/userop/configs/nacos-config-k8s.yaml` |
| userauth | `userauth.yaml` | `service/userauth/configs/nacosRemote-k8s.yaml` |
| gateway | `gateway.yaml` | `lushop/configs/nacos-config-k8s.yaml` |

**重要：** 导入配置时，将占位符 `YourDBPasswordHere` 和 `YourRedisPasswordHere` 替换为实际密码。

### 阶段 6: 部署业务服务

#### 6.1 部署微服务 (按依赖顺序)
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

### 阶段 7: 验证部署

#### 7.1 检查服务状态
```bash
# 查看所有 Pod
kubectl get pods -n lushop

# 查看服务
kubectl get svc -n lushop

# 查看部署状态
kubectl get deployments -n lushop
```

#### 7.2 验证服务连通性
```bash
# 测试 API 网关
NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
curl http://$NODE_IP:30080/health

# 端口转发测试
kubectl port-forward -n lushop svc/gateway-service 8001:8001 &
curl http://localhost:8001/health
```

#### 7.3 检查监控组件
```bash
# Prometheus
kubectl port-forward -n lushop svc/prometheus 9090:9090 &

# Grafana (admin/admin)
kubectl port-forward -n lushop svc/grafana 3000:3000 &

# Jaeger
kubectl port-forward -n lushop svc/jaeger 16686:16686 &

# Kibana
kubectl port-forward -n lushop svc/kibana 5601:5601 &
```

## 📊 服务端口映射

| 服务 | HTTP 端口 | gRPC 端口 | 备注 |
|------|----------|----------|------|
| Gateway | 8001 | 9001 | API 网关 (NodePort: 30080/30090) |
| User | 8011 | 50051 | 用户服务 |
| Goods | 8012 | 50052 | 商品服务 |
| Order | 8013 | 50053 | 订单服务 |
| Inventory | 8014 | 50054 | 库存服务 |
| UserOp | 8015 | 50055 | 用户操作服务 |
| UserAuth | 8016 | 50056 | 认证服务 |

## 🔧 自动化部署脚本

项目提供了完整的自动化部署脚本：

```bash
cd k8s

# 一键部署 (包含所有步骤)
./quick-deploy.sh

# 或使用标准部署脚本
./deploy.sh deploy

# 查看状态
./deploy.sh status

# 查看服务日志
./deploy.sh logs gateway
./deploy.sh logs user
```

## ⚠️ 注意事项

### 1. 密码配置
- 确保 `gen-secrets-custom.sh` 中的密码与 Nacos 配置一致
- MySQL 和 Redis 密码需要匹配

### 2. 资源要求
- MySQL: 至少 1GB RAM, 5GB 存储
- Redis: 至少 256MB RAM, 1GB 存储
- Elasticsearch: 至少 2GB RAM, 10GB 存储

### 3. 网络策略
- 确保服务间可以相互通信
- 配置适当的网络策略以提高安全性

### 4. 监控和日志
- 所有服务都集成了 Prometheus 监控
- 日志通过 Elasticsearch + Kibana 收集
- Jaeger 提供分布式链路追踪

## 🔍 故障排查

### 服务无法启动
```bash
# 查看 Pod 日志
kubectl logs -f deployment/<service-name> -n lushop

# 查看 Pod 详情
kubectl describe pod <pod-name> -n lushop
```

### 配置问题
```bash
# 检查 Nacos 配置
kubectl port-forward -n lushop svc/nacos 8848:8848 &
# 访问 http://localhost:8848/nacos 检查配置
```

### 网络连接问题
```bash
# 测试服务发现
kubectl exec -it deployment/user-service -n lushop -- nslookup mysql
kubectl exec -it deployment/user-service -n lushop -- nslookup redis
```

## 🎉 部署完成

部署完成后，你将拥有一个完整的生产级微服务电商平台：

- ✅ 7个微服务 + API 网关
- ✅ 服务注册发现 (Consul)
- ✅ 配置中心 (Nacos)
- ✅ 分布式链路追踪 (Jaeger)
- ✅ 监控告警 (Prometheus + Grafana)
- ✅ 日志收集 (Elasticsearch + Kibana)
- ✅ 消息队列 (RocketMQ)

可以通过 `http://<node-ip>:30080` 访问 API 网关！
