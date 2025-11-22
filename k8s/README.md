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
│   ├── elasticsearch/      # Elasticsearch 日志存储
│   ├── kibana/             # Kibana 日志可视化
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
├── generate-secrets.sh     # Secret 生成脚本
├── DEPLOYMENT_ANALYSIS.md  # 部署方案分析与完整指南
└── README.md               # 本文档
```

## 🚀 快速开始

### 前置要求

#### Kubernetes 集群要求

- **Kubernetes 版本**: v1.20+（推荐 v1.24+）
  - 支持 k3s、minikube、kind 或标准 Kubernetes 集群
  - 验证版本: `kubectl version --client --short`
- **kubectl 工具**: v1.20+，已正确配置 kubeconfig
  - 验证配置: `kubectl cluster-info`
  - 验证权限: `kubectl auth can-i create deployments`

#### 集群资源要求

**最小配置**（开发/测试环境）:
- CPU: 4 核
- 内存: 8GB
- 存储: 50GB（推荐使用 SSD）

**推荐配置**（生产环境）:
- CPU: 8+ 核
- 内存: 16GB+
- 存储: 100GB+（推荐使用 SSD）
- 节点数: 3+（高可用）

**资源分配建议**:
- 基础设施服务（MySQL、Redis、Nacos 等）: 约 4GB 内存，2 核 CPU
- 业务服务（7 个微服务 + 网关）: 约 4GB 内存，4 核 CPU
- 监控和日志（Prometheus、Grafana、ES 等）: 约 4GB 内存，2 核 CPU

#### 存储要求

- **StorageClass**: 集群需要配置默认 StorageClass 或手动指定
  - 验证: `kubectl get storageclass`
  - 常见: `local-path`（k3s）、`standard`（minikube）、`hostpath`（kind）
- **持久化存储**: 用于 MySQL、Redis、Elasticsearch 等有状态服务
  - MySQL: 至少 20GB
  - Elasticsearch: 至少 10GB
  - Redis: 至少 5GB

#### 网络要求

- **网络插件**: 集群已安装并配置 CNI 网络插件
  - k3s: 默认使用 Flannel
  - minikube: 默认使用 bridge
  - kind: 默认使用 bridge
- **Service 类型**: 支持 ClusterIP、NodePort、LoadBalancer
- **Ingress**（可选）: 如需使用 Ingress，需安装 Ingress Controller
  - 常见: Nginx Ingress、Traefik

#### 其他要求

- **Docker**: 已安装并运行（用于构建镜像）
  - 验证: `docker --version` 和 `docker ps`
- **镜像仓库**: 
  - 本地构建: Docker 守护进程运行中
  - 远程仓库: 确保集群节点可以拉取镜像（配置镜像拉取密钥）
- **服务镜像**: 已构建或使用预构建镜像
  - 使用 `./build-images.sh` 构建镜像
  - 或使用预构建的镜像（需修改清单中的镜像地址）

#### 快速验证

运行以下命令验证环境是否满足要求：

```bash
# 检查 Kubernetes 版本和连接
kubectl version --short

# 检查集群节点
kubectl get nodes

# 检查存储类
kubectl get storageclass

# 检查命名空间（如果已存在）
kubectl get namespace lushop

# 检查 Docker（用于构建）
docker --version && docker ps
```

### 单机 Kubernetes 部署

如果你还没有 Kubernetes 集群，可以使用以下任一方案在单机上快速部署一个本地集群。推荐使用 **k3s**（最简单轻量）或 **minikube**（功能完整）。

#### 方案一：k3s（推荐 - 最简单轻量）

k3s 是一个轻量级的 Kubernetes 发行版，专为边缘计算和单机部署设计，资源占用小，安装简单。

**安装步骤**：

```bash
# 1. 安装 k3s（使用国内镜像源加速）
curl -sfL https://rancher-mirror.rancher.cn/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn sh -

# 或使用官方源
# curl -sfL https://get.k3s.io | sh -

# 2. 检查安装状态并启动 k3s（如果未运行）
sudo systemctl status k3s

# 如果服务未运行，启动 k3s
sudo systemctl start k3s

# 设置开机自启（可选）
sudo systemctl enable k3s

# 等待几秒让 k3s 完全启动
sleep 5

# 3. 配置 kubectl（k3s 的 kubeconfig 在 /etc/rancher/k3s/k3s.yaml）
# 注意：只有在 k3s 服务运行后，配置文件才会生成
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER:$USER ~/.kube/config

# 4. 验证集群
kubectl cluster-info
kubectl get nodes

# 5. 检查存储类（k3s 默认使用 local-path）
kubectl get storageclass
```

**k3s 特点**：
- ✅ 资源占用小（约 512MB 内存）
- ✅ 安装简单，一条命令完成
- ✅ 内置 Traefik Ingress Controller
- ✅ 内置 local-path-provisioner（自动配置 StorageClass）
- ✅ 适合单机开发和测试

**卸载 k3s**（如需要）：
```bash
# 卸载 k3s
/usr/local/bin/k3s-uninstall.sh
```

---

#### 方案二：minikube（功能完整）

minikube 是官方提供的本地 Kubernetes 开发工具，功能完整，适合学习和开发。

**前置要求**：
- Docker 或虚拟机（VirtualBox/VMware/KVM）

**安装步骤**：

```bash
# 1. 安装 kubectl（如果还没有）
# Linux
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

# 或使用包管理器
# Ubuntu/Debian
sudo apt-get update && sudo apt-get install -y kubectl

# 2. 安装 minikube
# Linux
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# 或使用包管理器
# Ubuntu/Debian
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
sudo dpkg -i minikube_latest_amd64.deb

# 3. 启动 minikube（使用 Docker 驱动，需要 Docker 运行中）
minikube start --driver=docker

# 如果使用其他驱动
# minikube start --driver=virtualbox
# minikube start --driver=kvm2

# 4. 验证集群
kubectl cluster-info
kubectl get nodes

# 5. 检查存储类
kubectl get storageclass
```

**minikube 常用命令**：
```bash
# 启动集群
minikube start

# 停止集群
minikube stop

# 删除集群
minikube delete

# 查看状态
minikube status

# 访问 Dashboard（可选）
minikube dashboard

# 查看服务 URL
minikube service <service-name>
```

**minikube 特点**：
- ✅ 官方维护，功能完整
- ✅ 支持多种驱动（Docker、VirtualBox、KVM 等）
- ✅ 内置 Dashboard
- ✅ 适合学习和开发

---

#### 方案三：kind（基于 Docker）

kind（Kubernetes in Docker）使用 Docker 容器作为 Kubernetes 节点，非常适合 CI/CD 和本地测试。

**前置要求**：
- Docker 已安装并运行

**安装步骤**：

```bash
# 1. 安装 kubectl（如果还没有，参考 minikube 方案）

# 2. 安装 kind
# Linux
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# 3. 创建集群
kind create cluster --name lushop

# 4. 验证集群
kubectl cluster-info --context kind-lushop
kubectl get nodes

# 5. 安装存储类（kind 默认没有 StorageClass）
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.24/deploy/local-path-storage.yaml
kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```

**kind 常用命令**：
```bash
# 创建集群
kind create cluster --name <cluster-name>

# 删除集群
kind delete cluster --name <cluster-name>

# 列出集群
kind get clusters

# 加载本地镜像到集群（用于测试）
kind load docker-image <image-name> --name <cluster-name>
```

**kind 特点**：
- ✅ 基于 Docker，无需虚拟机
- ✅ 启动快速
- ✅ 适合 CI/CD 和自动化测试
- ⚠️ 需要手动配置 StorageClass

---

#### 方案选择建议

| 方案 | 适用场景 | 资源占用 | 安装难度 | 推荐度 |
|------|---------|---------|---------|--------|
| **k3s** | 单机开发、边缘计算 | 低（~512MB） | ⭐ 简单 | ⭐⭐⭐⭐⭐ |
| **minikube** | 学习、开发、功能测试 | 中（~2GB） | ⭐⭐ 中等 | ⭐⭐⭐⭐ |
| **kind** | CI/CD、自动化测试 | 中（~2GB） | ⭐⭐ 中等 | ⭐⭐⭐ |

**推荐**：对于单机部署，优先选择 **k3s**，安装最简单，资源占用最小。

---

#### 安装后验证

无论使用哪种方案，安装完成后请运行以下命令验证：

```bash
# 1. 检查集群连接
kubectl cluster-info

# 2. 检查节点状态
kubectl get nodes

# 3. 检查存储类（重要！）
kubectl get storageclass

# 4. 检查系统 Pod（k3s/minikube/kind 的系统组件）
kubectl get pods -n kube-system

# 5. 测试创建资源
kubectl create namespace test
kubectl delete namespace test
```

如果以上命令都能正常执行，说明 Kubernetes 集群已成功部署，可以继续后续的镜像构建和部署步骤。

---

#### 故障排查

**问题 1：k3s 配置文件不存在**

**错误信息**：
```bash
cp: 对 '/etc/rancher/k3s/k3s.yaml' 调用 stat 失败: 没有那个文件或目录
```

**原因**：k3s 服务未运行，配置文件只有在服务启动后才会生成。

**解决方案**：
```bash
# 1. 检查 k3s 服务状态
sudo systemctl status k3s

# 2. 如果服务未运行，启动 k3s
sudo systemctl start k3s

# 3. 等待几秒让服务完全启动
sleep 5

# 4. 再次检查服务状态，确认运行中
sudo systemctl status k3s

# 5. 现在配置文件应该已经生成，可以继续配置 kubectl
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER:$USER ~/.kube/config
chmod 600 ~/.kube/config
```

**问题 2：kubectl 配置文件权限被拒绝**

**错误信息**：
```bash
WARN[0000] Unable to read /etc/rancher/k3s/k3s.yaml, please start server with --write-kubeconfig-mode or --write-kubeconfig-group to modify kube config permissions 
error: error loading config file "/etc/rancher/k3s/k3s.yaml": open /etc/rancher/k3s/k3s.yaml: permission denied
```

**原因**：k3s 默认创建的配置文件 `/etc/rancher/k3s/k3s.yaml` 只有 root 用户可以读取，普通用户无法直接访问。

**解决方案**（推荐）：
```bash
# 1. 将配置文件复制到用户目录并设置正确的权限
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER:$USER ~/.kube/config
chmod 600 ~/.kube/config

# 2. 设置 KUBECONFIG 环境变量（永久生效）
echo '' >> ~/.bashrc
echo '# Kubernetes kubeconfig' >> ~/.bashrc
echo 'export KUBECONFIG=~/.kube/config' >> ~/.bashrc

# 3. 使环境变量立即生效（或重新打开终端）
source ~/.bashrc

# 4. 验证配置
kubectl cluster-info
kubectl get nodes
```

**注意**：如果配置文件已存在但 kubectl 仍然报错，可能是配置文件过期了。需要重新复制：
```bash
# 更新配置文件
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER:$USER ~/.kube/config
```

**替代方案**（修改 k3s 配置）：
```bash
# 1. 编辑 k3s 服务配置
sudo mkdir -p /etc/rancher/k3s
sudo tee /etc/rancher/k3s/config.yaml <<EOF
write-kubeconfig-mode: "0644"
EOF

# 2. 重启 k3s 使配置生效
sudo systemctl restart k3s

# 3. 等待服务启动
sleep 5

# 4. 现在可以直接使用配置文件
kubectl cluster-info
```

**问题 3：kubectl 无法连接集群**

**错误信息**：
```bash
The connection to the server <server>:6443 was refused
```

**解决方案**：
```bash
# 1. 确认 k3s 服务正在运行
sudo systemctl status k3s

# 2. 如果未运行，启动服务
sudo systemctl start k3s

# 3. 检查 k3s 日志查看错误
sudo journalctl -u k3s -n 50

# 4. 确认配置文件路径正确
ls -la /etc/rancher/k3s/k3s.yaml
```

**问题 4：k3s 启动失败**

**解决方案**：
```bash
# 1. 查看详细错误日志
sudo journalctl -u k3s -n 100 --no-pager

# 2. 检查端口占用（k3s 默认使用 6443 端口）
sudo netstat -tlnp | grep 6443

# 3. 检查防火墙设置
sudo ufw status
# 如果需要，允许相关端口
sudo ufw allow 6443/tcp

# 4. 尝试重启 k3s
sudo systemctl restart k3s

# 5. 如果问题持续，可以尝试重新安装
# 先卸载
/usr/local/bin/k3s-uninstall.sh
# 再重新安装
curl -sfL https://rancher-mirror.rancher.cn/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn sh -
```

**问题 5：minikube 启动失败**

**解决方案**：
```bash
# 1. 查看详细错误信息
minikube start --driver=docker --alsologtostderr

# 2. 检查 Docker 是否运行
docker ps

# 3. 删除旧集群并重新创建
minikube delete
minikube start --driver=docker

# 4. 如果使用其他驱动，确保驱动已安装
# VirtualBox
vboxmanage --version
# KVM
virsh --version
```

**问题 6：存储类（StorageClass）不存在**

**解决方案**：
```bash
# k3s：应该自动配置 local-path，如果没有：
kubectl get storageclass
# 如果为空，检查 local-path-provisioner
kubectl get pods -n local-path-storage

# minikube：应该自动配置，如果没有：
minikube addons enable default-storageclass

# kind：需要手动安装
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.24/deploy/local-path-storage.yaml
kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```

**问题 7：镜像拉取失败**

**错误信息**：
```bash
failed to pull image "rancher/mirrored-pause:3.6": failed to pull and unpack image "docker.io/rancher/mirrored-pause:3.6": 
failed to resolve reference "docker.io/rancher/mirrored-pause:3.6": 
failed to do request: Head "https://registry-1.docker.io/v2/...": dial tcp ...:443: i/o timeout
# 或
read tcp ...:443: read: connection reset by peer
```

**原因**：无法访问 Docker Hub（`docker.io`）拉取镜像，常见原因包括：
- 网络连接问题（防火墙、代理设置）
- 在中国大陆地区访问 Docker Hub 速度慢或被限制
- IPv6 连接问题

**解决方案**：

**方案 1：配置 k3s 使用国内镜像源（推荐，适用于中国大陆用户）**

```bash
# 1. 创建 k3s 镜像源配置文件
sudo mkdir -p /etc/rancher/k3s
sudo tee /etc/rancher/k3s/registries.yaml <<EOF
mirrors:
  "docker.io":
    endpoint:
      - "https://docker.mirrors.ustc.edu.cn"
      - "https://hub-mirror.c.163.com"
      - "https://mirror.baidubce.com"
  "registry.k8s.io":
    endpoint:
      - "https://registry.cn-hangzhou.aliyuncs.com/google_containers"
EOF

# 2. 重启 k3s 使配置生效
sudo systemctl restart k3s

# 3. 等待服务启动（约 10-30 秒）
sleep 10

# 4. 检查 Pod 状态
kubectl get pods -n kube-system

# 5. 如果 Pod 仍然无法启动，可以手动拉取镜像
# 使用 crictl（k3s 自带的容器工具）
sudo crictl pull docker.mirrors.ustc.edu.cn/rancher/mirrored-pause:3.6
sudo crictl tag docker.mirrors.ustc.edu.cn/rancher/mirrored-pause:3.6 rancher/mirrored-pause:3.6
```

**方案 2：手动拉取并导入镜像**

如果镜像源配置后仍然失败，可以手动拉取镜像：

```bash
# 1. 使用 Docker 拉取镜像（如果系统安装了 Docker）
docker pull docker.mirrors.ustc.edu.cn/rancher/mirrored-pause:3.6
docker tag docker.mirrors.ustc.edu.cn/rancher/mirrored-pause:3.6 rancher/mirrored-pause:3.6

# 2. 将镜像导入到 k3s 的 containerd
sudo k3s ctr images import <(docker save rancher/mirrored-pause:3.6)

# 3. 验证镜像已导入
sudo crictl images | grep pause
```

**方案 3：配置代理（如果有可用的代理）**

```bash
# 1. 创建 k3s 环境变量配置文件
sudo mkdir -p /etc/systemd/system/k3s.service.d
sudo tee /etc/systemd/system/k3s.service.d/http-proxy.conf <<EOF
[Service]
Environment="HTTP_PROXY=http://your-proxy:port"
Environment="HTTPS_PROXY=http://your-proxy:port"
Environment="NO_PROXY=localhost,127.0.0.1,0.0.0.0,10.0.0.0/8"
EOF

# 2. 重新加载 systemd 并重启 k3s
sudo systemctl daemon-reload
sudo systemctl restart k3s
```

**方案 4：禁用 IPv6（如果 IPv6 连接有问题）**

```bash
# 1. 编辑 k3s 配置，禁用 IPv6
sudo mkdir -p /etc/rancher/k3s
sudo tee -a /etc/rancher/k3s/config.yaml <<EOF
disable:
  - servicelb
  - traefik
node-ip: "127.0.0.1"
EOF

# 2. 重启 k3s
sudo systemctl restart k3s
```

**验证修复**：

```bash
# 1. 检查 k3s 服务状态
sudo systemctl status k3s

# 2. 检查 Pod 状态（应该看到 Pod 正在运行或已完成）
kubectl get pods -n kube-system

# 3. 查看 Pod 详细信息（如果有问题）
kubectl describe pod <pod-name> -n kube-system

# 4. 查看 Pod 日志
kubectl logs <pod-name> -n kube-system
```

**对于 minikube 用户**：

```bash
# 使用国内镜像源启动 minikube
minikube start --image-mirror-country=cn --driver=docker
```

---

### 构建镜像

在部署之前，需要先构建服务镜像。使用提供的构建脚本可以方便地构建所有或指定的镜像：

```bash
cd k8s

# 构建所有镜像（推荐，包括所有服务和网关）
./build-images.sh all

# 或仅构建服务镜像（不含网关）
./build-images.sh services

# 或仅构建网关镜像
./build-images.sh gateway

# 构建指定服务
./build-images.sh user
./build-images.sh goods
./build-images.sh order

# 查看已构建的镜像
./build-images.sh list

# 使用自定义标签构建
IMAGE_TAG=v1.0.0 ./build-images.sh all

# 查看帮助信息
./build-images.sh help
```

**可用服务**: `user`, `goods`, `order`, `inventory`, `userop`, `userauth`, `gateway`

**环境变量**:
- `IMAGE_TAG`: 镜像标签（默认: `latest`）
- `IMAGE_REGISTRY`: 镜像仓库前缀（默认: `lushop`）

**注意**: 
- 如果遇到 Docker iptables 错误，请参考 [故障排查](#-故障排查) 部分
- 构建过程会从项目根目录作为构建上下文，确保在正确的目录执行脚本

### 一键部署

**方式一：使用自动化部署脚本（推荐）**

使用提供的 `deploy.sh` 脚本可以自动完成整个部署流程：

```bash
cd k8s

# 1. 生成 Secret（首次部署必须）
./generate-secrets.sh

# 2. 构建镜像（如果还未构建）
./build-images.sh all

# 3. 导入镜像到集群（根据使用的集群类型选择）
# k3s: 使用 docker save/load 或配置镜像仓库
# minikube: minikube image load lushop/user:latest
# kind: kind load docker-image lushop/user:latest --name <cluster-name>

# 4. 执行部署
./deploy.sh deploy

# 查看部署状态
./deploy.sh status

# 查看服务日志
./deploy.sh logs gateway
./deploy.sh logs user

# 删除所有资源
./deploy.sh delete
```

**方式二：使用 Kustomize（手动部署）**

如果需要手动控制部署过程：

```bash
cd k8s

# 1. 生成 Secret
./generate-secrets.sh

# 2. 创建命名空间
kubectl apply -f base/namespace.yaml

# 3. 部署基础设施
kubectl apply -k base/redis
kubectl apply -k base/mysql
kubectl apply -k base/rocketmq

# 4. 等待基础设施就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# 5. 初始化数据库（需要创建 init_db.sql 脚本）
# kubectl port-forward -n lushop svc/mysql 3306:3306 &
# mysql -h 127.0.0.1 -uroot -p$(kubectl get secret mysql-auth -n lushop -o jsonpath='{.data.mysql-root-password}' | base64 -d) < ../scripts/init_db.sql

# 6. 部署服务治理组件
kubectl apply -k base/nacos
kubectl apply -k base/consul
kubectl apply -k base/prometheus
kubectl apply -k base/grafana
kubectl apply -k base/jaeger
kubectl apply -k base/elasticsearch
kubectl apply -k base/kibana

# 7. 配置 Nacos（手动）
# kubectl port-forward -n lushop svc/nacos 8848:8848
# 访问 http://localhost:8848/nacos，创建命名空间并导入配置

# 8. 部署业务服务
kubectl apply -k base/services/

# 查看部署状态
kubectl get all -n lushop
```

**详细部署步骤请参考**: [DEPLOYMENT_ANALYSIS.md](./DEPLOYMENT_ANALYSIS.md)

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
kubectl apply -k base/elasticsearch
kubectl apply -k base/kibana
kubectl apply -k base/services/
```

## 🚢 部署服务详细指南

### 部署顺序和依赖关系

为确保服务正常启动，建议按以下顺序部署：

```
1. 基础设施层（存储、消息队列）
   ├── Redis（缓存）
   ├── MySQL（数据库）
   └── RocketMQ（消息队列）

2. 服务治理层（配置中心、注册中心、监控）
   ├── Nacos（配置中心）
   ├── Consul（服务注册发现）
   ├── Prometheus（监控）
   ├── Grafana（可视化）
   ├── Jaeger（链路追踪）
   ├── Elasticsearch（日志存储）
   └── Kibana（日志可视化）

3. 业务服务层（按依赖顺序）
   ├── User（用户服务）
   ├── Goods（商品服务）
   ├── Inventory（库存服务）
   ├── Order（订单服务，依赖 User、Goods、Inventory）
   ├── UserOp（用户操作服务）
   ├── UserAuth（认证服务）
   └── Gateway（API 网关，依赖所有业务服务）
```

### 分步部署详细步骤

#### 步骤 1: 创建命名空间

```bash
# 创建命名空间
kubectl apply -f base/namespace.yaml

# 验证命名空间
kubectl get namespace lushop
```

#### 步骤 2: 部署存储层服务

```bash
# 部署 Redis
kubectl apply -k base/redis
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=redis -n lushop --timeout=300s

# 部署 MySQL
kubectl apply -k base/mysql
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# 部署 RocketMQ
kubectl apply -k base/rocketmq
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=rocketmq -n lushop --timeout=300s
```

#### 步骤 3: 初始化数据库

```bash
# 等待 MySQL 完全就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n lushop --timeout=300s

# 使用 port-forward 连接 MySQL
kubectl port-forward -n lushop svc/mysql 3306:3306 &
PF_PID=$!

# 导入数据库脚本（在另一个终端执行，或等待 port-forward 就绪后执行）
# 注意：需要根据实际情况修改数据库连接信息
mysql -h 127.0.0.1 -uroot -p$(kubectl get secret mysql-auth -n lushop -o jsonpath='{.data.root-password}' | base64 -d) < ../scripts/init_db.sql

# 停止 port-forward
kill $PF_PID
```

#### 步骤 4: 部署服务治理组件

```bash
# 部署 Nacos（配置中心）
kubectl apply -k base/nacos
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=nacos -n lushop --timeout=300s

# 部署 Consul（服务注册发现）
kubectl apply -k base/consul
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=consul -n lushop --timeout=300s

# 部署监控组件
kubectl apply -k base/prometheus
kubectl apply -k base/grafana

# 部署链路追踪
kubectl apply -k base/jaeger

# 部署日志组件
kubectl apply -k base/elasticsearch
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=elasticsearch -n lushop --timeout=300s
kubectl apply -k base/kibana
```

#### 步骤 5: 配置 Nacos

```bash
# 等待 Nacos 就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=nacos -n lushop --timeout=300s

# 使用 port-forward 访问 Nacos 控制台
kubectl port-forward -n lushop svc/nacos 8848:8848

# 在浏览器中访问: http://localhost:8848/nacos
# 默认账号: nacos / nacos
# 
# 需要执行的操作：
# 1. 创建命名空间: de9c6a0e-1fbc-425d-8d3b-09066fea6889
# 2. 为每个服务导入配置（参考各服务的 configs/nacos-config.yaml）
```

#### 步骤 6: 部署业务服务

```bash
# 按依赖顺序部署业务服务

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

#### 步骤 7: 验证部署

```bash
# 查看所有 Pod 状态
kubectl get pods -n lushop

# 查看所有服务
kubectl get svc -n lushop

# 查看所有部署
kubectl get deployments -n lushop

# 检查服务健康状态
kubectl get pods -n lushop -o wide
```

### 部署验证和健康检查

#### 检查服务状态

```bash
# 查看所有资源状态
kubectl get all -n lushop

# 查看 Pod 详细信息
kubectl get pods -n lushop -o wide

# 查看服务端点
kubectl get endpoints -n lushop

# 查看事件（排查问题）
kubectl get events -n lushop --sort-by='.lastTimestamp'
```

#### 服务健康检查

```bash
# 检查基础设施服务
kubectl exec -it -n lushop $(kubectl get pod -l app.kubernetes.io/name=redis -n lushop -o jsonpath='{.items[0].metadata.name}') -- redis-cli ping
kubectl exec -it -n lushop $(kubectl get pod -l app.kubernetes.io/name=mysql -n lushop -o jsonpath='{.items[0].metadata.name}') -- mysqladmin ping -uroot -p$(kubectl get secret mysql-auth -n lushop -o jsonpath='{.data.root-password}' | base64 -d)

# 检查业务服务健康端点
kubectl port-forward -n lushop svc/user-service 8011:8011 &
curl http://localhost:8011/health

kubectl port-forward -n lushop svc/gateway-service 8001:8001 &
curl http://localhost:8001/health
```

#### 查看服务日志

```bash
# 查看所有 Pod 日志
kubectl logs -f -l app.kubernetes.io/name=gateway -n lushop

# 查看指定服务日志
kubectl logs -f deployment/user-service -n lushop
kubectl logs -f deployment/gateway-service -n lushop

# 查看最近 100 行日志
kubectl logs --tail=100 deployment/gateway-service -n lushop

# 查看多个 Pod 的日志（如果有多副本）
kubectl logs -f -l app.kubernetes.io/name=gateway -n lushop --all-containers=true
```

### 更新和回滚服务

```bash
# 更新服务镜像
kubectl set image deployment/user-service user-service=lushop/user:v1.1.0 -n lushop

# 查看更新状态
kubectl rollout status deployment/user-service -n lushop

# 查看更新历史
kubectl rollout history deployment/user-service -n lushop

# 回滚到上一个版本
kubectl rollout undo deployment/user-service -n lushop

# 回滚到指定版本
kubectl rollout undo deployment/user-service --to-revision=2 -n lushop
```

### 扩缩容服务

```bash
# 扩展服务副本数
kubectl scale deployment/user-service --replicas=3 -n lushop

# 查看副本状态
kubectl get deployment user-service -n lushop

# 使用 HPA 自动扩缩容（需要先创建 HPA 资源）
kubectl autoscale deployment/user-service --cpu-percent=70 --min=2 --max=5 -n lushop
```

### 删除服务

```bash
# 删除单个服务
kubectl delete -k base/services/user

# 删除所有业务服务
kubectl delete -k base/services/

# 删除基础设施（谨慎操作，会删除数据）
kubectl delete -k base/redis
kubectl delete -k base/mysql

# 删除所有资源（包括命名空间）
kubectl delete -k base/
kubectl delete namespace lushop
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
- `elasticsearch-auth`（`k8s/base/elasticsearch/secret.yaml`）：Elasticsearch/Kibana 密码

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
- 业务服务请求: 128Mi 内存, 100m CPU
- 业务服务限制: 512Mi 内存, 500m CPU
- Gateway 请求: 256Mi 内存, 200m CPU
- Gateway 限制: 1Gi 内存, 1000m CPU
- 副本数: 1
- 所有服务已配置启动探针、就绪探针和存活探针
- 已添加时区和日志级别环境变量

**生产环境建议**:
- 根据实际负载调整资源限制
- 增加副本数实现高可用
- 配置 HPA/VPA 自动扩缩容

### 存储配置

- MySQL: 20Gi PVC
- Redis: 5Gi PVC
- RocketMQ: 10Gi PVC
- Prometheus: 10Gi PVC
- Elasticsearch: 20Gi PVC (数据) + 5Gi PVC (日志)

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
- **Elasticsearch**: 日志存储和搜索（端口 9200）
- **Kibana**: 日志可视化和分析（端口 5601，NodePort: 30561）

访问方式：
      ```bash
# Prometheus
kubectl port-forward -n lushop svc/prometheus 9090:9090

# Grafana
kubectl port-forward -n lushop svc/grafana 3000:3000

# Jaeger
kubectl port-forward -n lushop svc/jaeger 16686:16686

# Elasticsearch
kubectl port-forward -n lushop svc/elasticsearch 9200:9200
# 或直接访问: http://localhost:9200
# 默认用户名: elastic，密码在 elasticsearch-auth secret 中

# Kibana
# 方式1: NodePort (如果集群支持)
NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
# 访问: http://$NODE_IP:30561
# 方式2: Port Forward
kubectl port-forward -n lushop svc/kibana 5601:5601
# 访问: http://localhost:5601
# 默认用户名: elastic，密码在 elasticsearch-auth secret 中
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
3. **日志**: 已集成 Elasticsearch/Kibana，可配置 Filebeat 或 Logstash 进行日志收集
4. **安全**: 配置 NetworkPolicy，使用 TLS 加密，SealedSecret 管理敏感信息
5. **备份**: MySQL、RocketMQ Store、Prometheus、Elasticsearch 等数据卷需制定备份策略（Velero、快照、CronJob）
6. **CI/CD**: 集成 CI/CD 流水线自动构建和部署
7. **服务配置**: 已完善业务服务配置，包括健康检查、环境变量、资源限制等

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
# 部署所有服务（推荐）
kubectl apply -k base/

# 分步部署基础设施
kubectl apply -k base/redis
kubectl apply -k base/mysql
kubectl apply -k base/nacos
kubectl apply -k base/consul
kubectl apply -k base/rocketmq
kubectl apply -k base/prometheus
kubectl apply -k base/grafana
kubectl apply -k base/jaeger
kubectl apply -k base/elasticsearch
kubectl apply -k base/kibana

# 部署业务服务
kubectl apply -k base/services/

# 查看部署状态
kubectl get all -n lushop
kubectl get pods -n lushop -o wide

# 查看服务日志
kubectl logs -f deployment/gateway-service -n lushop
kubectl logs -f deployment/user-service -n lushop

# 删除服务
kubectl delete -k base/services/        # 删除业务服务
kubectl delete -k base/                 # 删除所有资源

# 扩缩容服务
kubectl scale deployment/user-service --replicas=3 -n lushop

# 更新服务镜像
kubectl set image deployment/user-service user-service=lushop/user:v1.1.0 -n lushop

# 回滚服务
kubectl rollout undo deployment/user-service -n lushop
```

