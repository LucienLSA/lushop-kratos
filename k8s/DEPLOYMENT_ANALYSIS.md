# K8s 部署方案分析与完整部署指南

## 📊 方案合理性分析

### ✅ 优点

1. **目录结构清晰**
   - 使用 Kustomize 管理配置，支持多环境
   - 基础设施和业务服务分离
   - 每个服务独立目录，便于管理

2. **资源配置合理**
   - 资源请求和限制设置合理
   - 健康检查配置完整（startup/readiness/liveness）
   - 使用 StatefulSet 管理有状态服务（MySQL、Redis 等）

3. **服务发现配置正确**
   - 使用 Kubernetes Service DNS 名称（如 `nacos`、`consul`）
   - ConfigMap 配置清晰

4. **镜像构建脚本完善**
   - 支持单独构建和批量构建
   - 包含错误检查和提示

### ⚠️ 存在的问题

1. **缺少数据库初始化脚本**
   - README 中提到了 `scripts/init_db.sql`，但实际不存在
   - 需要创建或明确数据库初始化方式

2. **部署顺序依赖不明确**
   - 虽然 README 有部署顺序，但缺少自动化脚本
   - 没有依赖检查和等待机制

3. **镜像构建和部署分离**
   - 需要先构建镜像再部署，但流程不够自动化
   - 本地镜像需要导入到集群（k3s/minikube/kind）

4. **Nacos 配置导入步骤不清晰**
   - 需要手动在 Nacos 控制台导入配置
   - 缺少自动化脚本或 Job

5. **缺少部署验证脚本**
   - 没有自动化的健康检查脚本
   - 缺少端到端测试

6. **Secret 管理**
   - Secret 使用明文示例值，需要提醒用户修改
   - 缺少 Secret 生成脚本

7. **存储类依赖**
   - 没有检查 StorageClass 是否存在
   - 某些环境（如 kind）需要手动安装

## 🔧 完整部署方案

### 阶段 1: 环境准备

#### 1.1 检查 Kubernetes 集群

```bash
# 检查集群状态
kubectl cluster-info
kubectl get nodes

# 检查存储类
kubectl get storageclass
# 如果没有默认存储类，需要安装（见 README 中的故障排查部分）
```

#### 1.2 检查 Docker 和镜像构建

```bash
# 检查 Docker
docker --version
docker ps

# 构建所有镜像
cd /home/zzx/GoProject/lushop-kratos/k8s
./build-images.sh all
```

#### 1.3 导入镜像到集群（如需要）

**对于 k3s:**
```bash
# k3s 使用 containerd，需要导入镜像
# 方法1: 使用 k3d 或直接使用本地镜像（如果配置了镜像仓库）
# 方法2: 使用 docker save/load
docker save lushop/user:latest | sudo k3s ctr images import -
docker save lushop/goods:latest | sudo k3s ctr images import -
# ... 其他镜像
```

**对于 minikube:**
```bash
minikube image load lushop/user:latest
minikube image load lushop/goods:latest
# ... 其他镜像
```

**对于 kind:**
```bash
kind load docker-image lushop/user:latest --name <cluster-name>
kind load docker-image lushop/goods:latest --name <cluster-name>
# ... 其他镜像
```

### 阶段 2: 准备 Secret 和配置

#### 2.1 生成 Secret（重要！）

**创建 Secret 生成脚本**（见下方完整脚本）

```bash
# 生成所有 Secret
./generate-secrets.sh
```

#### 2.2 验证 Secret

```bash
kubectl get secrets -n lushop
```

### 阶段 3: 分阶段部署

#### 3.1 创建命名空间

```bash
kubectl apply -f base/namespace.yaml
```

#### 3.2 部署基础设施层（存储）

```bash
# 1. Redis
kubectl apply -k base/redis
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=redis -n lushop --timeout=300s

# 2. MySQL
kubectl apply -k base/mysql
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# 3. RocketMQ
kubectl apply -k base/rocketmq
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=rocketmq -n lushop --timeout=300s
```

#### 3.3 初始化数据库

```bash
# 等待 MySQL 完全就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# 获取 MySQL root 密码
MYSQL_ROOT_PASSWORD=$(kubectl get secret mysql-auth -n lushop -o jsonpath='{.data.mysql-root-password}' | base64 -d)

# 使用 port-forward 连接 MySQL
kubectl port-forward -n lushop svc/mysql 3306:3306 &
PF_PID=$!
sleep 5

# 初始化数据库（需要创建 init_db.sql 脚本）
# 如果脚本不存在，可以手动连接数据库执行初始化
mysql -h 127.0.0.1 -uroot -p"$MYSQL_ROOT_PASSWORD" < ../scripts/init_db.sql || echo "请手动初始化数据库"

# 停止 port-forward
kill $PF_PID
```

#### 3.4 部署服务治理组件

```bash
# 1. Nacos
kubectl apply -k base/nacos
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=nacos -n lushop --timeout=300s

# 2. Consul
kubectl apply -k base/consul
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=consul -n lushop --timeout=300s

# 3. 监控组件
kubectl apply -k base/prometheus
kubectl apply -k base/grafana

# 4. 链路追踪
kubectl apply -k base/jaeger

# 5. 日志组件
kubectl apply -k base/elasticsearch
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=elasticsearch -n lushop --timeout=300s
kubectl apply -k base/kibana
```

#### 3.5 配置 Nacos

```bash
# 等待 Nacos 就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=nacos -n lushop --timeout=300s

# 使用 port-forward 访问 Nacos 控制台
kubectl port-forward -n lushop svc/nacos 8848:8848 &
NACOS_PF_PID=$!

echo "Nacos 控制台: http://localhost:8848/nacos"
echo "默认账号: nacos / nacos"
echo ""
echo "需要执行的操作："
echo "1. 创建命名空间: de9c6a0e-1fbc-425d-8d3b-09066fea6889"
echo "2. 为每个服务导入配置（参考各服务的 configs/nacos-config.yaml）"
echo ""
echo "按 Enter 继续（保持 port-forward 运行）..."
read

# 停止 port-forward（可选，如果需要继续使用可以保持运行）
# kill $NACOS_PF_PID
```

#### 3.6 部署业务服务

```bash
# 1. 基础服务（无依赖或依赖较少）
kubectl apply -k base/services/user
kubectl apply -k base/services/goods
kubectl apply -k base/services/inventory
kubectl apply -k base/services/userop
kubectl apply -k base/services/userauth

# 2. 等待基础服务就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=user -n lushop --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=goods -n lushop --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=inventory -n lushop --timeout=300s

# 3. 部署依赖服务
kubectl apply -k base/services/order

# 4. 部署 API 网关（最后部署，依赖所有业务服务）
kubectl apply -k base/services/gateway
```

### 阶段 4: 验证部署

#### 4.1 检查所有 Pod 状态

```bash
kubectl get pods -n lushop
kubectl get pods -n lushop -o wide
```

#### 4.2 检查服务端点

```bash
kubectl get svc -n lushop
kubectl get endpoints -n lushop
```

#### 4.3 健康检查

```bash
# 检查基础设施服务
kubectl exec -it -n lushop $(kubectl get pod -l app.kubernetes.io/name=redis -n lushop -o jsonpath='{.items[0].metadata.name}') -- redis-cli ping

# 检查业务服务
kubectl port-forward -n lushop svc/gateway-service 8001:8001 &
curl http://localhost:8001/health
```

## 🚀 一键部署脚本（推荐）

创建自动化部署脚本可以大大简化部署流程，见 `deploy.sh` 脚本。

## 📝 关键注意事项

1. **镜像必须提前构建并导入到集群**
2. **Secret 必须使用真实值，不要使用示例值**
3. **数据库初始化必须在部署业务服务前完成**
4. **Nacos 配置必须手动导入（或使用自动化脚本）**
5. **确保 StorageClass 存在，否则 PVC 无法创建**
6. **按照依赖顺序部署，不要一次性部署所有服务**

## 🔍 故障排查

参考 README.md 中的故障排查部分，常见问题包括：
- 镜像拉取失败
- Pod 无法启动
- 服务无法连接
- 存储类不存在
- 配置问题

