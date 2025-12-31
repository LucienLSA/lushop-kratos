# 🚀 Lushop K8s 手动部署指南

## 📋 目录
- [系统概述](#-系统概述)
- [架构图](#-架构图)
- [部署前准备](#-部署前准备)
- [第一步：环境验证](#第一步环境验证)
- [第二步：构建镜像](#第二步构建镜像)
- [第三步：生成密码](#第三步生成密码)
- [第四步：部署基础设施](#第四步部署基础设施)
- [第五步：配置Nacos](#第五步配置nacos)
- [第六步：部署业务服务](#第六步部署业务服务)
- [第七步：配置外部访问](#第七步配置外部访问)
- [第八步：验证部署](#第八步验证部署)
- [故障排除](#-故障排除)
- [常用命令](#-常用命令)
- [清理重置](#-清理重置)

---

## 🎯 系统概述

Lushop 是一个基于 Go-Kratos 框架构建的微服务电商平台，本指南将帮助您在 Kubernetes 集群中完整部署该系统。

### 核心组件

#### 业务服务 (7个)
- **lushop-gateway**: API网关 (HTTP 8001, gRPC 9001)
- **user-service**: 用户服务 (HTTP 8011, gRPC 50051)
- **goods-service**: 商品服务 (HTTP 8012, gRPC 50052)
- **order-service**: 订单服务 (HTTP 8013, gRPC 50053)
- **inventory-service**: 库存服务 (HTTP 8014, gRPC 50054)
- **userauth-service**: 用户认证服务 (HTTP 8015, gRPC 50055)
- **userop-service**: 用户操作服务 (HTTP 8016, gRPC 50056)

#### 基础设施服务 (11个)
- **MySQL**: 关系型数据库 (3306)
- **Redis**: 缓存数据库 (6379)
- **Nacos**: 配置中心 + 服务发现 (8848)
- **Consul**: 服务发现 (8500)
- **RocketMQ**: 消息队列 (9876/10911)
- **Elasticsearch**: 日志存储 (9200)
- **Kibana**: 日志可视化 (5601)
- **Prometheus**: 监控收集 (9090)
- **Grafana**: 监控面板 (3000)
- **Jaeger**: 链路追踪 (14268)
- **NGINX Ingress**: 外部访问代理

---

## 🏗️ 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    🌐 外部访问层                            │
│  • NGINX Ingress Controller (80/443)                       │
│  • Domain: lushop.local                                    │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                    🚪 API 网关                              │
│  • lushop-gateway (HTTP 8001, gRPC 9001)                  │
│  • 限流、认证、路由、熔断                                  │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                 🔄 微服务层 (6个业务服务)                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 🔐 userauth-svc   │ 👤 user-service   │ 📦 userop-svc │   │
│  │ (8015/50055)      │ (8011/50051)      │ (8016/50056)  │   │
│  ├─────────────────┼─────────────────┼─────────────────┤   │
│  │ 🛒 goods-service  │ 📋 order-service  │ 📊 inventory-svc│   │
│  │ (8012/50052)      │ (8013/50053)      │ (8014/50054)  │   │
│  └─────────────────┴─────────────────┴─────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                 🛠️ 基础设施层                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 🔧 Nacos (8848) │ 🔍 Consul (8500) │ 💾 MySQL (3306) │   │
│  ├─────────────────┼─────────────────┼─────────────────┤   │
│  │ 🚀 RocketMQ     │ 📊 Prometheus    │ 📈 Grafana      │   │
│  │ (9876/10911)    │ (9090)          │ (3000)          │   │
│  ├─────────────────┼─────────────────┼─────────────────┤   │
│  │ 🔗 Jaeger       │ 📝 ES            │ 📊 Kibana       │   │
│  │ (14268)         │ (9200)          │ (5601)          │   │
│  └─────────────────┴─────────────────┴─────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                 💽 数据持久化层                             │
│  • PersistentVolumeClaim (PVC)                            │
│  • local-path StorageClass                                │
│  • 数据持久化保证                                         │
└─────────────────────────────────────────────────────────────┘
```

---

## ⚙️ 部署前准备

### 环境要求
- **K8s 版本**: 1.20+
- **CPU**: 4核以上
- **内存**: 8GB以上
- **存储**: 50GB以上可用空间
- **网络**: 可访问 Docker Hub

### 必需工具
```bash
# 检查工具是否安装
kubectl version --client
docker --version

# 检查集群连接
kubectl cluster-info
kubectl get nodes
```

### 网络规划
- **集群内通信**: Service 名称访问
- **外部访问**: Ingress Controller (80/443)
- **域名**: lushop.local (本地 hosts 配置)

---

## 第一步：环境验证

### 1.1 验证 K8s 集群

```bash
# 检查集群状态
kubectl cluster-info

# 检查节点状态
kubectl get nodes -o wide

# 检查 StorageClass (必需)
kubectl get storageclass

# 示例输出:
# NAME                   PROVISIONER               RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
# local-path (default)   rancher.io/local-path      Delete          WaitForFirstConsumer   false                  2d
```

### 1.2 设置 kubeconfig (如需要)

```bash
# 如果 kubectl 无法连接集群
sudo cp /etc/kubernetes/admin.conf ~/.kube/config
sudo chown $(id -u):$(id -g) ~/.kube/config
chmod 600 ~/.kube/config

# 验证连接
kubectl cluster-info
kubectl get nodes
```

### 1.3 创建命名空间

```bash
# 创建 lushop 命名空间
kubectl create namespace lushop --dry-run=client -o yaml | kubectl apply -f -

# 验证创建
kubectl get namespaces
```

---

## 第二步：构建镜像

### 2.1 构建服务镜像

```bash
# 进入项目目录
cd /home/zzx/lucien/projects/lushop-kratos/k8s

# 构建所有服务镜像 (网关 + 6个微服务)
./build-images.sh all

# 查看构建结果
docker images | grep lushop

# 示例输出:
# lushop/gateway      latest    abc123def456    2 minutes ago    156MB
# lushop/user         latest    def456ghi789    2 minutes ago    142MB
# lushop/goods        latest    ghi789jkl012    2 minutes ago    148MB
# lushop/order        latest    jkl012mno345    2 minutes ago    145MB
# lushop/inventory    latest    mno345pqr678    2 minutes ago    141MB
# lushop/userauth     latest    pqr678stu901    2 minutes ago    139MB
# lushop/userop       latest    stu901vwx234    2 minutes ago    143MB
```

### 2.2 导入镜像到 K8s

```bash
# 打包所有镜像
docker save \
  lushop/gateway:latest \
  lushop/user:latest \
  lushop/goods:latest \
  lushop/order:latest \
  lushop/inventory:latest \
  lushop/userauth:latest \
  lushop/userop:latest \
  -o /tmp/lushop-services.tar

# 导入到 containerd (K8s 容器运行时)
sudo ctr -n k8s.io images import /tmp/lushop-services.tar

# 验证导入结果
sudo ctr -n k8s.io images ls | grep lushop

# 清理临时文件
rm /tmp/lushop-services.tar
```

---

## 第三步：生成密码

### 3.1 生成 K8s Secrets

```bash
# 生成包含自定义密码的 Secrets
./gen-secrets-custom.sh

# 查看生成的 Secrets
kubectl get secrets -n lushop

# 示例输出:
# NAME               TYPE     DATA   AGE
# mysql-auth         Opaque   4      10s
# redis-auth         Opaque   1      10s
# nacos-auth         Opaque   2      10s
# rocketmq-credentials Opaque 2      10s
# grafana-auth       Opaque   2      10s
# elasticsearch-auth Opaque   2      10s
```

### 3.2 查看生成的密码 (可选)

```bash
# 运行查看脚本 (如果存在)
chmod +x /tmp/view-secrets.sh 2>/dev/null || true
/tmp/view-secrets.sh 2>/dev/null || echo "查看脚本不存在，请手动查看"

# 或者手动查看特定密码
echo "MySQL Root Password: $(kubectl get secret mysql-auth -n lushop -o jsonpath='{.data.mysql-root-password}' | base64 -d)"
echo "Redis Password: $(kubectl get secret redis-auth -n lushop -o jsonpath='{.data.redis-password}' | base64 -d)"
```

### 3.3 自定义密码 (可选)

```bash
# 设置环境变量自定义密码
export MYSQL_ROOT_PASSWORD="MyRoot@123"
export MYSQL_PASSWORD="lushopDb@123"
export REDIS_PASSWORD="Redis@123"
export NACOS_MYSQL_PASSWORD="Nacos@123"
export GRAFANA_ADMIN_PASSWORD="Grafana@123"

# 重新生成 Secrets
./gen-secrets-custom.sh
```

---

## 第四步：部署基础设施

### 4.1 部署存储服务

```bash
# 部署 MySQL
echo "🚀 部署 MySQL..."
kubectl apply -k base/mysql
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# 检查 MySQL 状态
kubectl get pods -n lushop | grep mysql

# 部署 Redis
echo "🚀 部署 Redis..."
kubectl apply -k base/redis
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=redis -n lushop --timeout=300s

# 检查 Redis 状态
kubectl get pods -n lushop | grep redis

# 部署 RocketMQ
echo "🚀 部署 RocketMQ..."
kubectl apply -k base/rocketmq
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=rocketmq -n lushop --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=rocketmq-namesrv -n lushop --timeout=300s

# 检查 RocketMQ 状态
kubectl get pods -n lushop | grep rocketmq
```

### 4.2 部署配置和服务发现

```bash
# 部署 Nacos
echo "🚀 部署 Nacos..."
kubectl apply -k base/nacos
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=nacos -n lushop --timeout=300s

# 检查 Nacos 状态
kubectl get pods -n lushop | grep nacos

# 部署 Consul (可选)
echo "🚀 部署 Consul..."
kubectl apply -k base/consul
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=consul -n lushop --timeout=300s

# 检查 Consul 状态
kubectl get pods -n lushop | grep consul
```

### 4.3 部署监控和日志

```bash
# 部署 Prometheus
echo "🚀 部署 Prometheus..."
kubectl apply -k base/prometheus

# 部署 Grafana
echo "🚀 部署 Grafana..."
kubectl apply -k base/grafana

# 部署 Jaeger
echo "🚀 部署 Jaeger..."
kubectl apply -k base/jaeger

# 部署 Elasticsearch
echo "🚀 部署 Elasticsearch..."
kubectl apply -k base/elasticsearch
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=elasticsearch -n lushop --timeout=300s

# 部署 Kibana
echo "🚀 部署 Kibana..."
kubectl apply -k base/kibana

# 检查所有基础设施服务
kubectl get pods -n lushop -l 'app.kubernetes.io/part-of=lushop-system' -o wide
```

### 4.4 验证基础设施部署

```bash
# 查看所有 Pods
kubectl get pods -n lushop

# 查看 Services
kubectl get services -n lushop

# 查看 PVC (持久卷声明)
kubectl get pvc -n lushop

# 查看 PV (持久卷)
kubectl get pv
```

---

## 第五步：配置Nacos

### 5.1 导入服务配置

```bash
# 导入所有服务配置到 Nacos
./configure-nacos.sh import

# 验证配置导入
./configure-nacos.sh verify
```

### 5.2 手动验证 Nacos 配置

```bash
# 端口转发访问 Nacos 控制台
kubectl port-forward -n lushop svc/nacos 8848:8848 &

# 在浏览器中访问: http://localhost:8848/nacos
# 登录账号: nacos / nacos
# 命名空间: lushop (ID: de9c6a0e-1fbc-425d-8d3b-09066fea6889)

# 查看配置列表，应该包含:
# - gateway.yaml
# - user.yaml
# - goods.yaml
# - order.yaml
# - inventory.yaml
# - userauth.yaml
# - userop.yaml
```

### 5.3 Nacos 配置说明

Nacos 中的配置包含：
- **数据库连接**: MySQL 连接信息
- **缓存配置**: Redis 连接信息
- **服务发现**: Consul 连接信息
- **消息队列**: RocketMQ 连接信息
- **监控配置**: Jaeger 链路追踪
- **日志配置**: Elasticsearch 连接
- **服务端口**: HTTP 和 gRPC 端口配置

---

## 第六步：部署业务服务

### 6.1 部署业务服务

```bash
# 推荐：一键部署（deploy-lushop.sh）
# 脚本会：
#  - 创建命名空间 lushop（如不存在）
#  - （可选）通过 k8s/build-images.sh 构建服务镜像（使用 --build）
#  - 应用 k8s/ 目录下的 kustomize 清单
#  - 等待核心业务部署（gateway + 6 个微服务）就绪（超时提示）
cd /home/zzx/lucien/projects/lushop-kratos
chmod +x ./deploy-lushop.sh
./deploy-lushop.sh            # 只应用 manifests（不构建镜像）
./deploy-lushop.sh --build    # 构建镜像后再应用 manifests（本地构建）

# 手动方式（等同于一键脚本所做的应用步骤）：
kubectl apply -k k8s

# 或者分别部署
kubectl apply -f deployments.yaml
kubectl apply -f services.yaml
kubectl apply -f common-configmap.yaml
```

#### deploy-lushop.sh 详情
- 用法：`./deploy-lushop.sh [--build]`
- `--build`：调用 `k8s/build-images.sh all` 构建镜像（需要宿主机有构建环境）
- 输出与等待：脚本会等待指定的核心 deployments 可用（默认 5 分钟），如未就绪会输出 Warning
- 建议：生产环境建议直接使用已推到 registry 的镜像，跳过 `--build`；在本地或离线环境使用 `--build` 并按文档将镜像导入到集群运行时（containerd）

#### 常见排错
- 若某个 Deployment 未就绪，先 `kubectl logs -n lushop <pod>` 查看错误日志
- 若镜像拉取失败，请确保 `imagePullSecrets` 或已把镜像导入到节点运行时

### 6.2 等待服务就绪

```bash
# 等待网关就绪
kubectl wait --for=condition=available deployment/lushop-gateway -n lushop --timeout=300s

# 等待各个微服务就绪
kubectl wait --for=condition=available deployment/user-service -n lushop --timeout=300s
kubectl wait --for=condition=available deployment/goods-service -n lushop --timeout=300s
kubectl wait --for=condition=available deployment/order-service -n lushop --timeout=300s
kubectl wait --for=condition=available deployment/inventory-service -n lushop --timeout=300s
kubectl wait --for=condition=available deployment/userauth-service -n lushop --timeout=300s
kubectl wait --for=condition=available deployment/userop-service -n lushop --timeout=300s

# 查看所有服务状态
kubectl get pods -n lushop -o wide
```

### 6.3 检查服务健康状态

```bash
# 检查各个服务的健康状态
for service in lushop-gateway user-service goods-service order-service inventory-service userauth-service userop-service; do
  echo "检查 $service..."
  kubectl get pods -n lushop -l app.kubernetes.io/name=$service
done
```

---

## 第七步：配置外部访问

### 7.1 部署 Ingress Controller

```bash
# 检查是否已安装 NGINX Ingress
kubectl get pods -n ingress-nginx

# 如果未安装，安装 NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# 等待 Ingress Controller 就绪
kubectl wait --for=condition=available deployment/ingress-nginx-controller -n ingress-nginx --timeout=300s
```

### 7.2 配置 Ingress 规则

```bash
# 应用 Ingress 配置
kubectl apply -f ingress.yaml

# 查看 Ingress 状态
kubectl get ingress -n lushop

# 示例输出:
# NAME            CLASS    HOSTS                          ADDRESS   PORTS   AGE
# lushop-ingress  nginx    lushop.local,nacos.lushop.local,grafana.lushop.local,kibana.lushop.local   192.168.1.100   80      10s
```

### 7.3 配置本地 hosts

```bash
# 添加域名解析到 /etc/hosts
echo "127.0.0.1 lushop.local nacos.lushop.local grafana.lushop.local kibana.lushop.local" | sudo tee -a /etc/hosts

# 或者添加到 ~/.hosts 文件 (macOS)
# echo "127.0.0.1 lushop.local nacos.lushop.local grafana.lushop.local kibana.lushop.local" >> ~/.hosts
```

---

## 第八步：验证部署

### 8.1 查看整体状态

```bash
# 查看所有 Pods
kubectl get pods -n lushop -o wide

# 查看所有 Services
kubectl get services -n lushop

# 查看所有 Ingress
kubectl get ingress -n lushop

# 查看所有 PVC
kubectl get pvc -n lushop

# 查看所有 Secrets
kubectl get secrets -n lushop
```

### 8.2 测试服务连通性

```bash
# 测试 API 网关健康检查
curl -v http://lushop.local/api/v1/health

# 测试各个微服务健康检查
for port in 8011 8012 8013 8014 8015 8016; do
  echo "测试端口 $port..."
  curl -v http://lushop.local:$port/health
done
```

### 8.3 访问管理界面

```bash
# Nacos 控制台
kubectl port-forward -n lushop svc/nacos 8848:8848 &
echo "访问 Nacos: http://localhost:8848/nacos (账号: nacos/nacos)"

# Grafana 监控面板
kubectl port-forward -n lushop svc/grafana 3000:3000 &
echo "访问 Grafana: http://localhost:3000 (账号: admin/你的密码)"

# Kibana 日志面板
kubectl port-forward -n lushop svc/kibana 5601:5601 &
echo "访问 Kibana: http://localhost:5601"

# Jaeger 链路追踪
kubectl port-forward -n lushop svc/jaeger 16686:16686 &
echo "访问 Jaeger: http://localhost:16686"
```

### 8.4 验证数据持久化

```bash
# 检查 MySQL 数据持久化
kubectl exec -n lushop mysql-0 -- mysql -u root -p$MYSQL_ROOT_PASSWORD -e "SHOW DATABASES;"

# 检查 Redis 数据持久化
kubectl exec -n lushop redis-0 -- redis-cli -a $REDIS_PASSWORD ping

# 检查 PVC 状态
kubectl get pvc -n lushop -o wide
```

---

## 🔧 故障排除

### 常见问题及解决方法

#### 1. Pod 处于 Pending 状态

```bash
# 检查 PVC 绑定状态
kubectl get pvc -n lushop
kubectl describe pvc <pvc-name> -n lushop

# 检查 StorageClass
kubectl get storageclass

# 检查节点资源
kubectl describe nodes
```

#### 2. 镜像拉取失败 (ImagePullBackOff)

```bash
# 检查镜像是否存在
sudo ctr -n k8s.io images ls | grep lushop

# 重新导入镜像
docker save ... -o /tmp/lushop-services.tar
sudo ctr -n k8s.io images import /tmp/lushop-services.tar

# 重启失败的 Pod
kubectl delete pod <pod-name> -n lushop
```

#### 3. 服务启动失败

```bash
# 查看 Pod 日志
kubectl logs <pod-name> -n lushop -f

# 查看 Pod 详细信息
kubectl describe pod <pod-name> -n lushop

# 检查服务依赖
kubectl get pods -n lushop | grep -E "(mysql|redis|nacos)"
```

#### 4. Nacos 配置问题

```bash
# 检查 Nacos Pod 状态
kubectl get pods -n lushop | grep nacos
kubectl logs <nacos-pod> -n lushop

# 重新导入配置
./configure-nacos.sh import

# 检查配置是否正确导入
kubectl port-forward -n lushop svc/nacos 8848:8848 &
curl http://localhost:8848/nacos/v1/cs/configs?dataId=gateway.yaml&group=lushop_grpc
```

#### 5. 网络连接问题

```bash
# 检查 Service 解析
kubectl exec -n lushop <pod-name> -- nslookup mysql

# 检查网络策略
kubectl get networkpolicies -n lushop

# 测试服务间通信
kubectl exec -n lushop <pod-name> -- curl http://user-service:8011/health
```

#### 6. Ingress 访问问题

```bash
# 检查 Ingress Controller
kubectl get pods -n ingress-nginx

# 检查 Ingress 配置
kubectl describe ingress lushop-ingress -n lushop

# 检查本地 hosts 配置
cat /etc/hosts | grep lushop
```

---

## 📊 常用命令

### 监控和调试

```bash
# 实时查看 Pod 状态
kubectl get pods -n lushop -w

# 查看 Pod 资源使用
kubectl top pods -n lushop

# 查看节点资源使用
kubectl top nodes

# 查看事件
kubectl get events -n lushop --sort-by=.metadata.creationTimestamp

# 查看日志
kubectl logs -f <pod-name> -n lushop

# 进入容器调试
kubectl exec -it <pod-name> -n lushop -- /bin/sh
```

### 服务管理

```bash
# 重启服务
kubectl rollout restart deployment/<service-name> -n lushop

# 扩展副本数
kubectl scale deployment <service-name> -n lushop --replicas=3

# 更新镜像
kubectl set image deployment/<service-name> <container-name>=<new-image> -n lushop

# 查看服务端点
kubectl get endpoints -n lushop
```

### 配置管理

```bash
# 查看 ConfigMap
kubectl get configmap -n lushop
kubectl describe configmap <config-name> -n lushop

# 查看 Secret (注意: 内容是 base64 编码的)
kubectl get secret <secret-name> -n lushop -o yaml

# 解码 Secret 值
kubectl get secret <secret-name> -n lushop -o jsonpath='{.data.<key>}' | base64 -d
```

---

## 🧹 清理重置

### 完全清理

```bash
# 删除整个命名空间 (删除所有资源)
kubectl delete namespace lushop

# 删除 Ingress Controller
kubectl delete namespace ingress-nginx

# 删除所有 PVC 和 PV (谨慎操作!)
kubectl delete pvc,pv --all

# 删除所有镜像 (可选)
docker rmi $(docker images | grep lushop | awk '{print $3}')
sudo ctr -n k8s.io images ls | grep lushop | awk '{print $1}' | xargs sudo ctr -n k8s.io images rm
```

### 一键删除 Lushop 相关资源（推荐，带备份）

仓库提供 `delete-lushop.sh` 脚本用于备份并**选择性**删除 Lushop 命名空间内的业务资源（Deployments/Services/ConfigMaps/Ingress/CronJobs 等）。脚本默认**不会**删除命名空间本身或 PVC（可以避免误删持久数据）。

使用方法：
```bash
cd /home/zzx/lucien/projects/lushop-kratos
chmod +x delete-lushop.sh
./delete-lushop.sh
```

脚本行为（默认安全模式）：
- 将 `lushop` 命名空间内资源导出到 `./lushop-backup-<timestamp>/` 目录作为备份（包含 all、rbac 等）。
- 删除带有 `app.kubernetes.io/name` 标签的 Deployments/StatefulSets/DaemonSets/ReplicaSets/Jobs/CronJobs，删除带标签的 Services/Ingress。
- 删除符合命名规则的 ConfigMaps（以 `-config` 结尾）和 Secrets（以 `-secret` 结尾）以降低误删风险。
- 脚本不会删除 PVC 或 PV；如需删除持久数据，请谨慎手动执行 `kubectl delete pvc <name> -n lushop` 或删除命名空间（见下）。

示例：仅备份，不删除（dry-run 手动步骤）
```bash
TS=$(date +%Y%m%dT%H%M%S)
kubectl get all,cm,secret,ing,sts,daemonset,cronjob,pvc -n lushop -o yaml > "./lushop-backup-${TS}/lushop-all-resources.yaml"
```

如需彻底移除命名空间及所有资源（含 PVC/PV），在确认备份后运行：
```bash
kubectl delete namespace lushop
```

注意事项与恢复
- 备份文件位于运行脚本的当前目录下 `./lushop-backup-<timestamp>/`，包含资源清单，可用于审计或部分恢复。
- 恢复流程通常为解析备份 YAML 并按需 `kubectl apply -f <resource>`（Secret 内容为 base64，注意安全）。

### 部分清理

```bash
# 只删除业务服务
kubectl delete -f deployments.yaml
kubectl delete -f services.yaml

# 只删除基础设施
kubectl delete -k base/

# 删除特定的 PVC
kubectl delete pvc <pvc-name> -n lushop
```

### 重置特定服务

```bash
# 重置单个服务
kubectl delete deployment <service-name> -n lushop
kubectl apply -f deployments.yaml  # 重新部署

# 重置数据库数据
kubectl exec -n lushop mysql-0 -- mysql -u root -p -e "DROP DATABASE lushop; CREATE DATABASE lushop;"
```

---

## 🎯 部署完成检查清单

- [ ] K8s 集群正常运行
- [ ] StorageClass 可用
- [ ] 所有镜像构建完成并导入
- [ ] Secrets 生成成功
- [ ] 基础设施服务全部运行
- [ ] Nacos 配置正确导入
- [ ] 业务服务全部部署
- [ ] Ingress 配置正确
- [ ] 本地 hosts 配置完成
- [ ] 所有服务可通过域名访问
- [ ] 管理界面可正常登录

---

## 📚 学习资源

### 官方文档
- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [Go-Kratos 框架文档](https://go-kratos.dev/)
- [Nacos 配置中心](https://nacos.io/)

### 架构学习
- [微服务架构模式](https://microservices.io/)
- [云原生应用](https://landscape.cncf.io/)
- [容器化最佳实践](https://docs.docker.com/develop/)

---

## 🎉 恭喜完成！

如果您按照此指南成功完成所有步骤，恭喜您已经部署了一个完整的、生产级别的微服务电商平台！

### 访问地址总结

- **API网关**: http://lushop.local/api/v1
- **Nacos控制台**: http://nacos.lushop.local (nacos/nacos)
- **Grafana监控**: http://grafana.lushop.local (admin/你的密码)
- **Kibana日志**: http://kibana.lushop.local
- **Jaeger链路**: http://localhost:16686 (端口转发)

### 下一步建议

1. **测试API接口**: 使用 Postman 或 curl 测试业务接口
2. **配置监控面板**: 在 Grafana 中创建业务监控图表
3. **日志分析**: 配置 Kibana 仪表板进行日志分析
4. **性能调优**: 根据实际负载调整资源限制
5. **备份策略**: 设置定期数据备份
6. **安全加固**: 配置网络策略和访问控制

---

*最后更新时间: 2025年1月1日 | 版本: v1.0 | 作者: AI Assistant*
