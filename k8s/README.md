# Lushop Kubernetes 部署指南

本文档详细介绍如何将 Lushop 微服务架构部署到 Kubernetes 集群中。

## 项目概述

Lushop 是一个基于 Go 和 Kratos 框架的微服务电商系统，包含以下服务：

- **lushop-gateway**: API 网关服务
- **goods-service**: 商品服务
- **inventory-service**: 库存服务
- **order-service**: 订单服务
- **user-service**: 用户服务
- **userauth-service**: 用户认证服务
- **userop-service**: 用户操作服务

## 架构图

```
┌─────────────────┐    ┌─────────────────┐
│   Nginx Ingress │    │   Cert Manager  │
│                 │    │                 │
└─────────┬───────┘    └─────────────────┘
          │
          ▼
┌─────────────────┐
│ Lushop Gateway  │
│   (API网关)     │
└─────────┬───────┘
          │
    ┌─────┼─────┐
    │     │     │
    ▼     ▼     ▼
┌─────┐ ┌─────┐ ┌─────┐
│Goods│ │User │ │Order│
└─────┘ └─────┘ └─────┘
    │     │     │
    ▼     ▼     ▼
┌─────┐ ┌─────┐ ┌─────┐
│Inv. │ │Auth │ │User │
│     │ │     │ │ Op. │
└─────┘ └─────┘ └─────┘

Infrastructure Layer:
MySQL, Redis, Nacos, Elasticsearch, RocketMQ
```

## 前置要求

### 1. Kubernetes 集群

- Kubernetes 1.19+
- kubectl 客户端工具
- kubeconfig 配置正确

### 2. 必要的 Kubernetes 组件

```bash
# 安装 NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# 安装 cert-manager (用于 SSL 证书)
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.3/cert-manager.yaml
```

### 3. 容器镜像

确保所有服务镜像已推送到容器注册表：

```bash
# 检查镜像是否存在
docker pull crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/aliyun1123466419/lushop:latest
docker pull crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/aliyun1123466419/goods:latest
# ... 其他服务镜像
```

## 快速开始

### 1. 克隆项目并进入 k8s 目录

```bash
git clone https://github.com/your-org/lushop-kratos.git
cd lushop-kratos/k8s
```

### 2. 配置 Secret

编辑 `common-secret.yaml` 文件，将 base64 编码的敏感信息替换为实际值：

```bash
# 生成 base64 编码的密码
echo -n "your_mysql_password" | base64
echo -n "your_redis_password" | base64
# ... 其他密码
```

### 3. 部署基础设施

```bash
# 部署基础设施服务 (MySQL, Redis, Nacos)
kubectl apply -f infrastructure.yaml
```

### 4. 部署应用服务

```bash
# 使用部署脚本
./deploy.sh deploy
```

或者手动部署：

```bash
# 使用 Kustomize
kubectl apply -k .

# 或直接应用 YAML 文件
kubectl apply -f namespace.yaml
kubectl apply -f common-configmap.yaml
kubectl apply -f common-secret.yaml
kubectl apply -f deployments.yaml
kubectl apply -f services.yaml
kubectl apply -f ingress.yaml
```

### 5. 验证部署

```bash
# 检查所有资源状态
kubectl get all -n lushop

# 查看 Pod 状态
kubectl get pods -n lushop -o wide

# 查看服务状态
kubectl get services -n lushop

# 查看 Ingress
kubectl get ingress -n lushop
```

## 配置说明

### 环境变量配置

所有配置都通过 ConfigMap 和 Secret 管理：

- **common-configmap.yaml**: 应用配置 (数据库连接、端口等)
- **common-secret.yaml**: 敏感信息 (密码、密钥等)

### 域名配置

编辑 `ingress.yaml` 中的域名：

```yaml
spec:
  tls:
  - hosts:
    - api.yourdomain.com  # 替换为您的域名
    - lushop.yourdomain.com
```

### SSL 证书

使用 cert-manager 自动生成 Let's Encrypt 证书：

```yaml
annotations:
  cert-manager.io/cluster-issuer: "letsencrypt-prod"
```

## 部署脚本使用

### 完整部署

```bash
./deploy.sh deploy
```

### 更新镜像

```bash
# 更新所有服务到最新镜像
./update-images.sh

# 更新特定服务
./update-images.sh goods

# 使用特定标签
./update-images.sh -t v1.2.3
```

### 服务扩展

```bash
# 扩展 goods 服务到 5 个副本
./deploy.sh scale goods-service 5
```

### 回滚部署

```bash
# 查看回滚历史
kubectl rollout history deployment/goods-service -n lushop

# 回滚到上一版本
./update-images.sh --rollback 1 goods
```

### 查看状态

```bash
./deploy.sh status
```

## 监控和日志

### 查看日志

```bash
# 查看特定服务日志
kubectl logs -f deployment/goods-service -n lushop

# 查看所有服务日志
kubectl logs -f -l app.kubernetes.io/name=lushop -n lushop
```

### 资源监控

```bash
# 查看资源使用情况
kubectl top pods -n lushop
kubectl top nodes
```

### 健康检查

```bash
# 检查服务健康状态
curl https://api.yourdomain.com/health
```

## 故障排除

### 常见问题

1. **Pod 无法启动**
   ```bash
   kubectl describe pod <pod-name> -n lushop
   kubectl logs <pod-name> -n lushop
   ```

2. **镜像拉取失败**
   ```bash
   # 检查镜像是否存在
   docker pull <image-name>

   # 检查镜像仓库权限
   kubectl create secret docker-registry acr-secret \
     --docker-server=crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com \
     --docker-username=<username> \
     --docker-password=<password> \
     --namespace=lushop
   ```

3. **网络连接问题**
   ```bash
   kubectl exec -it <pod-name> -n lushop -- curl <service-name>:8080/health
   ```

4. **存储问题**
   ```bash
   kubectl get pvc -n lushop
   kubectl describe pvc <pvc-name> -n lushop
   ```

### 调试技巧

```bash
# 进入容器调试
kubectl exec -it <pod-name> -n lushop -- /bin/bash

# 端口转发用于本地调试
kubectl port-forward svc/lushop-gateway-service 8000:8000 -n lushop

# 查看事件
kubectl get events -n lushop --sort-by=.metadata.creationTimestamp
```

## 生产环境考虑

### 安全加固

1. **网络策略**
   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: lushop-network-policy
     namespace: lushop
   spec:
     podSelector: {}
     policyTypes:
     - Ingress
     - Egress
   ```

2. **RBAC**
   ```yaml
   apiVersion: rbac.authorization.k8s.io/v1
   kind: Role
   metadata:
     name: lushop-role
     namespace: lushop
   rules:
   - apiGroups: [""]
     resources: ["pods", "services"]
     verbs: ["get", "list", "watch"]
   ```

### 高可用配置

1. **多副本部署**
   ```yaml
   spec:
     replicas: 3  # 生产环境建议 3+ 副本
   ```

2. **Pod 反亲和性**
   ```yaml
   spec:
     template:
       spec:
         affinity:
           podAntiAffinity:
             requiredDuringSchedulingIgnoredDuringExecution:
             - labelSelector:
                 matchExpressions:
                 - key: app
                   operator: In
                   values:
                   - lushop-gateway
               topologyKey: kubernetes.io/hostname
   ```

### 备份策略

1. **数据库备份**
   ```bash
   # 创建数据库备份 CronJob
   kubectl apply -f backup-job.yaml
   ```

2. **配置备份**
   ```bash
   # 备份所有配置
   kubectl get all -n lushop -o yaml > backup.yaml
   ```

## CI/CD 集成

### GitHub Actions 示例

```yaml
name: Deploy to Kubernetes
on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Configure kubectl
      uses: azure/k8s-set-context@v3
      with:
        method: kubeconfig
        kubeconfig: ${{ secrets.KUBE_CONFIG }}

    - name: Deploy to Kubernetes
      run: |
        cd k8s
        ./deploy.sh deploy
```

## 清理资源

```bash
# 删除所有资源
./deploy.sh cleanup

# 或手动删除
kubectl delete namespace lushop
```

## 支持

如果遇到问题，请：

1. 查看 [故障排除](#故障排除) 部分
2. 检查 Kubernetes 集群状态
3. 查看服务日志
4. 提交 Issue 到项目仓库

## 版本历史

- v1.0.0: 初始 Kubernetes 部署配置
- v1.1.0: 添加基础设施自动化部署
- v1.2.0: 改进监控和日志配置
