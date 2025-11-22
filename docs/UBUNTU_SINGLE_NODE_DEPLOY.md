# Ubuntu 24.04 单机部署指南

本文档详细说明如何在 Ubuntu 24.04 服务器上单机部署 lushop-kratos 微服务电商平台。

## 📋 目录

- [系统要求](#系统要求)
- [前置准备](#前置准备)
- [安装依赖](#安装依赖)
- [获取项目代码](#获取项目代码)
- [构建镜像](#构建镜像)
- [配置 Secret](#配置-secret)
- [部署基础设施](#部署基础设施)
- [初始化数据库](#初始化数据库)
- [配置 Nacos](#配置-nacos)
- [部署业务服务](#部署业务服务)
- [验证部署](#验证部署)
- [访问服务](#访问服务)
- [故障排查](#故障排查)
- [常用命令](#常用命令)

---

## 系统要求

### 硬件要求

- **CPU**: 至少 4 核（推荐 8 核）
- **内存**: 至少 8GB（推荐 16GB，标准 Kubernetes 需要更多资源）
- **磁盘**: 至少 100GB 可用空间（推荐 SSD）
- **网络**: 可访问互联网（用于拉取镜像和组件）

> **注意**: 标准 Kubernetes 相比轻量级发行版（如 k3s）需要更多资源。如果资源受限，可以考虑使用 k3s 或 minikube。

### 软件要求

- **操作系统**: Ubuntu 24.04 LTS
- **内核版本**: 5.15+（推荐 6.8+）
- **用户权限**: sudo 权限

---

## 前置准备

### 1. 更新系统

```bash
# 更新系统包
sudo apt update && sudo apt upgrade -y

# 安装基础工具
sudo apt install -y curl wget git vim net-tools
```

### 2. 配置主机名（可选）

```bash
# 设置主机名
sudo hostnamectl set-hostname lushop-server

# 验证
hostname
```

### 3. 配置防火墙（如果启用）

```bash
# 检查防火墙状态
sudo ufw status

# 如果需要开放端口（根据实际需求）
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 30080/tcp   # Gateway HTTP
sudo ufw allow 30090/tcp   # Gateway gRPC
sudo ufw allow 30561/tcp   # Kibana
```

---

## 安装依赖

### 1. 安装 Docker

```bash
# 卸载旧版本（如果有）
sudo apt remove -y docker docker-engine docker.io containerd runc

# 安装依赖
sudo apt install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# 添加 Docker 官方 GPG 密钥
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# 添加 Docker 仓库
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安装 Docker Engine
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 启动 Docker 服务
sudo systemctl enable docker
sudo systemctl start docker

# 验证安装
docker --version
docker compose version

# 将当前用户添加到 docker 组（可选，避免每次使用 sudo）
sudo usermod -aG docker $USER
# 注意：需要重新登录才能生效，或使用 newgrp docker
```

### 2. 安装 Kubernetes 组件

```bash
# 安装必要的工具
sudo apt install -y apt-transport-https ca-certificates curl gpg

# 设置 Kubernetes 版本（可根据需要修改，推荐使用最新稳定版）
K8S_VERSION="v1.30"  # 或使用 v1.29, v1.28 等

# 添加 Kubernetes 官方 GPG 密钥
curl -fsSL https://pkgs.k8s.io/core:/stable:/${K8S_VERSION}/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

# 添加 Kubernetes 仓库
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/${K8S_VERSION}/deb/ /" | sudo tee /etc/apt/sources.list.d/kubernetes.list

# 更新并安装 kubectl、kubelet、kubeadm
sudo apt update
sudo apt install -y kubelet kubeadm kubectl

# 锁定版本，防止自动更新
sudo apt-mark hold kubelet kubeadm kubectl

# 验证安装
kubectl version --client
kubeadm version
```

> **注意**: 确保 kubelet、kubeadm 和 kubectl 版本一致。如果需要安装特定版本，可以在安装时指定：`sudo apt install -y kubelet=1.30.0-1.1 kubeadm=1.30.0-1.1 kubectl=1.30.0-1.1`

### 3. 配置容器运行时

Kubernetes 需要容器运行时。由于我们已经安装了 Docker，需要配置 containerd（Docker 使用 containerd 作为运行时）：

```bash
# 配置 containerd
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml

# 修改配置以使用 systemd cgroup 驱动
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml

# 重启 containerd
sudo systemctl restart containerd
sudo systemctl enable containerd

# 验证 containerd 状态
sudo systemctl status containerd
```

### 4. 初始化 Kubernetes 集群

```bash
# 关闭 swap（Kubernetes 要求）
sudo swapoff -a
# 永久禁用 swap（编辑 /etc/fstab，注释掉 swap 行）
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

# 加载必要的内核模块
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter

# 配置网络参数
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sudo sysctl --system

# 初始化单节点 Kubernetes 集群
sudo kubeadm init \
  --pod-network-cidr=10.244.0.0/16 \
  --apiserver-advertise-address=$(hostname -I | awk '{print $1}') \
  --control-plane-endpoint=$(hostname -I | awk '{print $1}')

# 按照输出提示配置 kubectl（非 root 用户）
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# 如果是 root 用户，执行：
# export KUBECONFIG=/etc/kubernetes/admin.conf

# 验证集群状态
kubectl get nodes
# 注意：节点会显示 NotReady，因为还没有安装网络插件
```

### 5. 安装网络插件（Flannel）

```bash
# 安装 Flannel 网络插件
kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml

# 等待网络插件就绪
kubectl wait --for=condition=ready pod -l app=flannel -n kube-flannel --timeout=300s

# 验证节点状态（应该显示 Ready）
kubectl get nodes

# 允许在 master 节点上调度 Pod（单节点集群需要）
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
```

### 6. 安装 kustomize（可选）

```bash
# 安装 kustomize
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
sudo mv kustomize /usr/local/bin/

# 验证
kustomize version
```

### 7. 安装 MySQL 客户端（用于初始化数据库）

```bash
sudo apt install -y mysql-client
```

---

## 获取项目代码

```bash
# 克隆项目（如果还没有）
cd ~
git clone <your-repo-url> lushop-kratos
# 或使用已有项目
cd /home/lucien/goproject/lushop-kratos

# 进入项目目录
cd lushop-kratos
```

---

## 构建镜像

### 1. 检查 Docker 服务

```bash
# 确保 Docker 正在运行
sudo systemctl status docker

# 如果未运行，启动它
sudo systemctl start docker
```

### 2. 构建服务镜像

项目需要构建以下镜像：

- `lushop/user:latest`
- `lushop/goods:latest`
- `lushop/order:latest`
- `lushop/inventory:latest`
- `lushop/userop:latest`
- `lushop/userauth:latest`
- `lushop/gateway:latest`

**方式一：使用构建脚本（如果有）**

```bash
cd k8s
# 如果存在 build-images.sh
chmod +x build-images.sh
./build-images.sh all
```

**方式二：手动构建**

```bash
# 构建用户服务
cd service/user
docker build -t lushop/user:latest .

# 构建商品服务
cd ../goods
docker build -t lushop/goods:latest .

# 构建订单服务
cd ../order
docker build -t lushop/order:latest .

# 构建库存服务
cd ../inventory
docker build -t lushop/inventory:latest .

# 构建用户操作服务
cd ../userop
docker build -t lushop/userop:latest .

# 构建认证服务
cd ../userauth
docker build -t lushop/userauth:latest .

# 构建网关服务
cd ../../lushop
docker build -t lushop/gateway:latest .
```

**方式三：使用 Makefile（如果支持）**

```bash
# 在项目根目录
make build
```

### 3. 验证镜像

```bash
# 查看所有镜像
docker images | grep lushop

# 应该看到以下镜像：
# lushop/user:latest
# lushop/goods:latest
# lushop/order:latest
# lushop/inventory:latest
# lushop/userop:latest
# lushop/userauth:latest
# lushop/gateway:latest
```

### 4. 配置镜像访问（重要）

标准 Kubernetes 使用 containerd 作为容器运行时（Docker 也使用 containerd）。有几种方式让 Kubernetes 使用本地构建的镜像：

**方法1：使用本地镜像仓库（推荐）**

```bash
# 为镜像打标签，使其看起来像是从本地仓库拉取的
# 或者直接使用本地镜像（如果 imagePullPolicy 设置为 IfNotPresent 或 Never）

# 查看镜像
docker images | grep lushop
```

**方法2：配置 containerd 使用 Docker 镜像**

由于 Docker 和 Kubernetes 都使用 containerd，镜像通常是共享的。确保 Deployment 中的 `imagePullPolicy` 设置为 `IfNotPresent` 或 `Never`：

```yaml
# 在 Deployment 中设置
spec:
  containers:
  - name: user
    image: lushop/user:latest
    imagePullPolicy: IfNotPresent  # 或 Never
```

**方法3：导入镜像到 containerd（如果需要）**

```bash
# 将 Docker 镜像导入到 containerd
for img in lushop/user lushop/goods lushop/order lushop/inventory lushop/userop lushop/userauth lushop/gateway; do
    docker save $img:latest | sudo ctr -n k8s.io images import -
done

# 验证镜像已导入
sudo ctr -n k8s.io images ls | grep lushop
```

**方法4：配置私有镜像仓库（生产环境推荐）**

在生产环境中，建议搭建私有镜像仓库（如 Harbor）或使用 Docker Hub 等公共仓库。

---

## 配置 Secret

### 1. 修改 MySQL Secret

```bash
cd k8s/base/mysql
vim secret.yaml
```

修改密码（建议使用强密码）：

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysql-auth
type: Opaque
stringData:
  mysql-root-password: "your-root-password"      # 修改这里
  mysql-user: "lushop"
  mysql-password: "your-lushop-password"         # 修改这里
  mysql-database: "lushop"
```

### 2. 修改 Redis Secret

```bash
cd ../redis
vim secret.yaml
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: redis-auth
type: Opaque
stringData:
  redis-password: "your-redis-password"          # 修改这里
```

### 3. 修改 Elasticsearch Secret

```bash
cd ../elasticsearch
vim secret.yaml
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: elasticsearch-auth
type: Opaque
stringData:
  elasticsearch-password: "your-elastic-password" # 修改这里
  kibana-password: "your-kibana-password"        # 修改这里
```

### 4. 修改其他 Secret（可选）

根据需要修改：
- `k8s/base/nacos/secret.yaml` - Nacos 配置
- `k8s/base/rocketmq/secret.yaml` - RocketMQ 配置
- `k8s/base/grafana/secret.yaml` - Grafana 配置

---

## 部署基础设施

### 1. 部署所有资源

```bash
cd k8s/base

# 一键部署所有服务
kubectl apply -k .

# 或分步部署（推荐，便于排查问题）
kubectl apply -f namespace.yaml
kubectl apply -k redis
kubectl apply -k mysql
kubectl apply -k consul
kubectl apply -k nacos
kubectl apply -k jaeger
kubectl apply -k rocketmq
kubectl apply -k prometheus
kubectl apply -k grafana
kubectl apply -k elasticsearch
kubectl apply -k kibana
```

### 2. 检查部署状态

```bash
# 查看所有 Pod 状态
kubectl get pods -n lushop

# 查看 Pod 详细信息
kubectl get pods -n lushop -o wide

# 实时监控 Pod 状态
watch -n 2 kubectl get pods -n lushop

# 查看服务状态
kubectl get svc -n lushop

# 查看 PVC 状态
kubectl get pvc -n lushop
```

### 3. 等待所有 Pod 就绪

```bash
# 等待所有 Pod 就绪（最多等待 10 分钟）
kubectl wait --for=condition=ready pod --all -n lushop --timeout=600s

# 如果某些 Pod 未就绪，查看日志
kubectl logs <pod-name> -n lushop
kubectl describe pod <pod-name> -n lushop
```

### 4. 常见问题处理

**问题1：Pod 一直处于 Pending 状态**

```bash
# 查看 Pod 详情
kubectl describe pod <pod-name> -n lushop

# 常见原因：
# 1. 资源不足 - 检查节点资源
kubectl describe node

# 2. PVC 未绑定 - 检查 StorageClass
kubectl get storageclass
# 如果不存在 StorageClass，需要安装一个（如 local-path-provisioner 或 nfs-client）
# 安装 local-path-provisioner（推荐用于单节点）
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.24/deploy/local-path-storage.yaml
kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```

**问题2：Pod 一直处于 ImagePullBackOff 状态**

```bash
# 检查镜像是否存在
docker images | grep lushop
# 或检查 containerd
sudo ctr -n k8s.io images ls | grep lushop

# 如果不存在，重新构建或导入镜像
# 方法1：重新构建
docker build -t lushop/user:latest ./service/user

# 方法2：导入到 containerd
docker save lushop/user:latest | sudo ctr -n k8s.io images import -

# 方法3：检查 Deployment 的 imagePullPolicy
kubectl get deployment <deployment-name> -n lushop -o yaml | grep imagePullPolicy
# 确保设置为 IfNotPresent 或 Never
```

**问题3：MySQL Pod 无法启动**

```bash
# 查看 MySQL 日志
kubectl logs -f <mysql-pod-name> -n lushop

# 检查 PVC 是否正确创建
kubectl get pvc -n lushop | grep mysql
```

---

## 初始化数据库

### 1. 等待 MySQL 就绪

```bash
# 等待 MySQL Pod 就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# 验证 MySQL 服务
kubectl get svc mysql -n lushop
```

### 2. 端口转发

```bash
# 在另一个终端执行（保持运行）
kubectl port-forward -n lushop svc/mysql 3306:3306
```

### 3. 导入数据库脚本

```bash
# 如果项目有数据库初始化脚本
# 方式1：使用 mysql 客户端
mysql -h 127.0.0.1 -uroot -p<your-root-password> < scripts/init_db.sql

# 方式2：使用 kubectl exec
# 先找到 MySQL Pod 名称
MYSQL_POD=$(kubectl get pod -n lushop -l app.kubernetes.io/name=mysql -o jsonpath='{.items[0].metadata.name}')
# 复制脚本到 Pod
kubectl cp scripts/init_db.sql lushop/$MYSQL_POD:/tmp/init_db.sql
# 执行脚本
kubectl exec -it $MYSQL_POD -n lushop -- mysql -uroot -p<your-root-password> < /tmp/init_db.sql

# 方式3：如果没有初始化脚本，手动创建数据库和表
mysql -h 127.0.0.1 -uroot -p<your-root-password>
```

在 MySQL 客户端中执行：

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS lushop CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE lushop;

-- 创建表（根据项目实际情况）
-- 这里需要根据项目的数据库脚本创建表结构
```

### 4. 验证数据库

```bash
# 连接数据库验证
mysql -h 127.0.0.1 -uroot -p<your-root-password> -e "SHOW DATABASES;"
mysql -h 127.0.0.1 -uroot -p<your-root-password> -e "USE lushop; SHOW TABLES;"
```

---

## 配置 Nacos

### 1. 等待 Nacos 就绪

```bash
# 等待 Nacos Pod 就绪
kubectl wait --for=condition=ready pod -l app=nacos -n lushop --timeout=300s

# 查看 Nacos 服务
kubectl get svc nacos -n lushop
```

### 2. 访问 Nacos 控制台

```bash
# 端口转发
kubectl port-forward -n lushop svc/nacos 8848:8848
```

在浏览器中访问：http://localhost:8848/nacos

- **用户名**: nacos
- **密码**: nacos

### 3. 创建命名空间

1. 登录 Nacos 控制台
2. 进入「命名空间」菜单
3. 创建新命名空间：
   - **命名空间 ID**: `de9c6a0e-1fbc-425d-8d3b-09066fea6889`
   - **命名空间名**: `lushop_grpc`
   - **描述**: Lushop 微服务配置

### 4. 导入服务配置

根据项目中的配置文件（通常在 `service/*/configs/nacos-config.yaml`），为每个服务创建配置：

1. 进入「配置管理」→「配置列表」
2. 选择命名空间：`lushop_grpc`
3. 点击「+」创建配置
4. 为每个服务创建配置（DataId 格式：`<service-name>.yaml`）

**示例配置（user.yaml）**：

```yaml
server:
  http:
    addr: 0.0.0.0:8011
    timeout: 1s
  grpc:
    addr: 0.0.0.0:50051
    timeout: 1s

data:
  database:
    driver: mysql
    source: lushop:lushop123456@tcp(mysql.lushop.svc.cluster.local:3306)/lushop?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    addr: redis.lushop.svc.cluster.local:6379
    password: your-redis-password
    db: 0
    dial_timeout: 1s
    read_timeout: 0.2s
    write_timeout: 0.2s

registry:
  consul:
    address: consul.lushop.svc.cluster.local:8500
    scheme: http
```

为以下服务创建配置：
- `user.yaml`
- `goods.yaml`
- `order.yaml`
- `inventory.yaml`
- `userop.yaml`
- `userauth.yaml`
- `gateway.yaml`

---

## 部署业务服务

### 1. 部署所有业务服务

```bash
cd k8s/base

# 部署所有业务服务
kubectl apply -k services/

# 或单独部署
kubectl apply -k services/user
kubectl apply -k services/goods
kubectl apply -k services/order
kubectl apply -k services/inventory
kubectl apply -k services/userop
kubectl apply -k services/userauth
kubectl apply -k services/gateway
```

### 2. 检查服务状态

```bash
# 查看所有业务服务 Pod
kubectl get pods -n lushop -l app.kubernetes.io/part-of=lushop

# 查看服务
kubectl get svc -n lushop

# 查看 Gateway 服务（应该显示 NodePort）
kubectl get svc gateway-service -n lushop
```

### 3. 查看服务日志

```bash
# 查看 Gateway 日志
kubectl logs -f -l app.kubernetes.io/name=gateway -n lushop

# 查看 User 服务日志
kubectl logs -f -l app.kubernetes.io/name=user -n lushop

# 查看所有服务日志
for svc in user goods order inventory userop userauth gateway; do
    echo "=== $svc ==="
    kubectl logs -l app.kubernetes.io/name=$svc -n lushop --tail=20
done
```

---

## 验证部署

### 1. 检查所有 Pod 状态

```bash
# 所有 Pod 应该处于 Running 状态
kubectl get pods -n lushop

# 预期输出示例：
# NAME                              READY   STATUS    RESTARTS   AGE
# consul-xxx                        1/1     Running   0          5m
# elasticsearch-0                   1/1     Running   0          5m
# elasticsearch-1                   1/1     Running   0          5m
# gateway-xxx                       1/1     Running   0          2m
# goods-xxx                         1/1     Running   0          2m
# grafana-xxx                       1/1     Running   0          5m
# inventory-xxx                     1/1     Running   0          2m
# jaeger-xxx                        1/1     Running   0          5m
# kibana-xxx                        1/1     Running   0          5m
# mysql-xxx                         1/1     Running   0          5m
# nacos-xxx                         1/1     Running   0          5m
# order-xxx                         1/1     Running   0          2m
# prometheus-xxx                    1/1     Running   0          5m
# redis-xxx                         1/1     Running   0          5m
# rocketmq-xxx                      1/1     Running   0          5m
# user-xxx                          1/1     Running   0          2m
# userauth-xxx                      1/1     Running   0          2m
# userop-xxx                        1/1     Running   0          2m
```

### 2. 检查服务注册

```bash
# 端口转发 Consul
kubectl port-forward -n lushop svc/consul 8500:8500

# 在另一个终端查看服务
curl http://localhost:8500/v1/catalog/services

# 或在浏览器访问
# http://localhost:8500
```

### 3. 测试 API

```bash
# 获取节点 IP
NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')

# 测试 Gateway API（使用 NodePort）
curl http://$NODE_IP:30080/api/goods/list

# 或使用端口转发
kubectl port-forward -n lushop svc/gateway-service 8001:8001
curl http://localhost:8001/api/goods/list
```

### 4. 检查健康状态

```bash
# 检查 Gateway 健康状态
curl http://localhost:8001/health

# 检查各个服务的健康状态
for svc in user goods order inventory userop userauth; do
    echo "=== $svc ==="
    kubectl exec -n lushop $(kubectl get pod -n lushop -l app.kubernetes.io/name=$svc -o jsonpath='{.items[0].metadata.name}') -- wget -qO- http://localhost:8011/health || echo "Failed"
done
```

---

## 访问服务

### 1. Gateway API

**方式1：NodePort（推荐）**

```bash
# 获取节点 IP
NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')

# HTTP API
curl http://$NODE_IP:30080/api/goods/list

# gRPC API（需要使用 grpcurl 或客户端）
```

**方式2：Port Forward**

```bash
# HTTP
kubectl port-forward -n lushop svc/gateway-service 8001:8001
curl http://localhost:8001/api/goods/list

# gRPC
kubectl port-forward -n lushop svc/gateway-service 9001:9001
```

### 2. 监控服务

**Consul UI**
```bash
kubectl port-forward -n lushop svc/consul 8500:8500
# 访问: http://localhost:8500
```

**Nacos UI**
```bash
kubectl port-forward -n lushop svc/nacos 8848:8848
# 访问: http://localhost:8848/nacos (nacos/nacos)
```

**Jaeger UI**
```bash
kubectl port-forward -n lushop svc/jaeger 16686:16686
# 访问: http://localhost:16686
```

**Prometheus**
```bash
kubectl port-forward -n lushop svc/prometheus 9090:9090
# 访问: http://localhost:9090
```

**Grafana**
```bash
kubectl port-forward -n lushop svc/grafana 3000:3000
# 访问: http://localhost:3000
# 默认用户名/密码在 grafana-admin secret 中
```

**Kibana**
```bash
# 方式1：NodePort
NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
# 访问: http://$NODE_IP:30561

# 方式2：Port Forward
kubectl port-forward -n lushop svc/kibana 5601:5601
# 访问: http://localhost:5601
# 默认用户名: elastic，密码在 elasticsearch-auth secret 中
```

### 3. 服务端口映射

| 服务 | HTTP 端口 | gRPC 端口 | NodePort | 访问方式 |
|------|-----------|-----------|----------|----------|
| Gateway | 8001 | 9001 | 30080/30090 | NodePort 或 Port Forward |
| User | 8011 | 50051 | - | Port Forward |
| Goods | 8012 | 50052 | - | Port Forward |
| Order | 8013 | 50053 | - | Port Forward |
| Inventory | 8014 | 50054 | - | Port Forward |
| UserOp | 8015 | 50055 | - | Port Forward |
| UserAuth | 8016 | 50056 | - | Port Forward |

---

## 故障排查

### 1. Pod 无法启动

```bash
# 查看 Pod 详情
kubectl describe pod <pod-name> -n lushop

# 查看 Pod 日志
kubectl logs <pod-name> -n lushop

# 查看前一个容器的日志（如果容器重启）
kubectl logs <pod-name> -n lushop --previous
```

### 2. 服务无法连接

```bash
# 检查服务端点
kubectl get endpoints -n lushop

# 测试 DNS 解析
kubectl run -it --rm debug --image=busybox --restart=Never -n lushop -- nslookup mysql.lushop.svc.cluster.local

# 测试服务连通性
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -n lushop -- curl http://mysql.lushop.svc.cluster.local:3306
```

### 3. 镜像拉取失败

```bash
# 检查镜像是否存在
docker images | grep lushop
# 或检查 containerd
sudo ctr -n k8s.io images ls | grep lushop

# 如果不存在，重新构建或导入
docker build -t lushop/user:latest ./service/user
# 导入到 containerd
docker save lushop/user:latest | sudo ctr -n k8s.io images import -

# 检查 Pod 的镜像拉取策略
kubectl get pod <pod-name> -n lushop -o yaml | grep imagePullPolicy

# 如果策略是 Always，修改为 IfNotPresent
kubectl patch deployment <deployment-name> -n lushop -p '{"spec":{"template":{"spec":{"containers":[{"name":"<container-name>","imagePullPolicy":"IfNotPresent"}]}}}}'
```

### 4. 存储问题

```bash
# 检查 PVC 状态
kubectl get pvc -n lushop

# 检查 StorageClass
kubectl get storageclass

# 如果不存在 StorageClass，安装 local-path-provisioner（单节点推荐）
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.24/deploy/local-path-storage.yaml
kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

# 或使用其他存储方案（如 NFS、Ceph 等）
```

### 5. 资源不足

```bash
# 查看节点资源
kubectl describe node

# 查看 Pod 资源使用
kubectl top pods -n lushop

# 如果资源不足，可以：
# 1. 减少副本数
# 2. 降低资源限制
# 3. 增加服务器资源
```

### 6. 网络问题

```bash
# 检查 Service
kubectl get svc -n lushop

# 检查 Service 详情
kubectl describe svc <service-name> -n lushop

# 测试服务连通性
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -n lushop -- curl http://<service-name>.lushop.svc.cluster.local:<port>
```

### 7. 配置问题

```bash
# 查看 ConfigMap
kubectl get configmap -n lushop
kubectl describe configmap <configmap-name> -n lushop

# 查看 Secret
kubectl get secret -n lushop
kubectl describe secret <secret-name> -n lushop

# 查看环境变量
kubectl exec <pod-name> -n lushop -- env | grep -i <key>
```

---

## 常用命令

### 查看状态

```bash
# 查看所有资源
kubectl get all -n lushop

# 查看 Pod 状态
kubectl get pods -n lushop

# 查看服务状态
kubectl get svc -n lushop

# 查看 PVC 状态
kubectl get pvc -n lushop

# 查看事件
kubectl get events -n lushop --sort-by='.lastTimestamp'
```

### 日志查看

```bash
# 查看所有服务日志
for svc in user goods order inventory userop userauth gateway; do
    echo "=== $svc ==="
    kubectl logs -l app.kubernetes.io/name=$svc -n lushop --tail=50
done

# 实时查看日志
kubectl logs -f -l app.kubernetes.io/name=gateway -n lushop
```

### 重启服务

```bash
# 重启指定服务
kubectl rollout restart deployment <deployment-name> -n lushop

# 重启所有业务服务
for svc in user goods order inventory userop userauth gateway; do
    kubectl rollout restart deployment $svc -n lushop
done
```

### 删除和清理

```bash
# 删除所有资源
kubectl delete -k k8s/base

# 删除命名空间（会删除所有资源）
kubectl delete namespace lushop

# 清理 PVC（谨慎操作，会删除数据）
kubectl delete pvc --all -n lushop
```

### 调试命令

```bash
# 进入 Pod 执行命令
kubectl exec -it <pod-name> -n lushop -- /bin/sh

# 查看 Pod 资源使用
kubectl top pod <pod-name> -n lushop

# 查看节点资源
kubectl top node
```

---

## 后续优化

1. **高可用**: 增加副本数，配置 HPA/VPA
2. **监控**: 配置 ServiceMonitor，完善 Prometheus 监控
3. **日志**: 配置 Filebeat 或 Logstash 收集日志到 Elasticsearch
4. **安全**: 配置 NetworkPolicy，使用 TLS 加密
5. **备份**: 制定 MySQL、Elasticsearch 等数据备份策略
6. **CI/CD**: 集成 CI/CD 流水线自动构建和部署

---

## 参考文档

- [Kubernetes 部署指引](../k8s/README.md)
- [项目主 README](../README.md)
- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [kubeadm 安装指南](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/)
- [Docker 官方文档](https://docs.docker.com/)
- [containerd 官方文档](https://containerd.io/docs/)

---

## 常见问题 FAQ

### Q1: 单节点 Kubernetes 集群有什么限制？

A: 单节点集群适合开发和测试环境。主要限制包括：
- 没有高可用性（单点故障）
- 资源受限（所有组件运行在同一节点）
- 不适合生产环境

生产环境建议使用多节点集群，至少 3 个 master 节点和多个 worker 节点。

### Q2: 如何备份数据？

A: 可以使用 Velero 或手动备份 PVC。对于 MySQL，可以使用 `mysqldump` 或 `kubectl exec` 执行备份命令。

### Q3: 如何更新服务？

A: 构建新镜像后，更新 Deployment 的镜像标签，或使用 `kubectl set image` 命令。

### Q4: 如何查看服务日志？

A: 使用 `kubectl logs` 命令，或配置日志收集工具（如 Filebeat）将日志发送到 Elasticsearch。

### Q5: 如何扩容服务？

A: 使用 `kubectl scale` 命令增加副本数，或配置 HPA 自动扩缩容。

---

**部署完成后，恭喜您成功部署了 lushop-kratos 微服务电商平台！** 🎉

如有问题，请参考故障排查部分或查看项目文档。

