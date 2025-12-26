# K8s 中 Nacos 与各服务配置协调机制

## 📋 配置协调概述

在 K8s 环境中，lushop 项目的配置协调采用**分层配置管理**策略：

```
本地配置 (config.yaml) → Nacos 配置中心 → 业务服务
```

## 🔄 配置协调流程

### 阶段 1: 部署顺序协调
```
1. 基础设施层 (MySQL, Redis, RocketMQ)
2. 配置中心 (Nacos)
3. 配置导入 (自动/手动)
4. 业务服务 (按依赖顺序)
```

### 阶段 2: 配置同步机制
```
业务服务启动时:
1. 读取本地 config.yaml 获取 Nacos 连接信息
2. 连接 Nacos 拉取完整业务配置
3. 应用配置并启动服务
```

## 📁 配置层次结构

### 本地配置层 (`config.yaml`)
**位置**: `service/*/configs/config.yaml`
**用途**: 提供 Nacos 连接信息
**内容**:
```yaml
nacos:
  addr: nacos
  port: 8848
  namespaceId: de9c6a0e-1fbc-425d-8d3b-09066fea6889
  dataId: {service}.yaml
  groupId: lushop_grpc
```

### Nacos 配置层 (`nacos-config-k8s.yaml`)
**位置**: `service/*/configs/nacos-config-k8s.yaml`
**用途**: 完整的业务配置
**特点**:
- 使用 K8s 服务名 (mysql, redis, consul)
- 包含完整的业务逻辑配置
- 支持环境变量占位符

## 🔧 协调机制详解

### 1. 服务发现协调
```yaml
# Nacos 配置中的服务发现
registry:
  consul:
    address: consul:8500
    scheme: http

service:
  user:
    endpoint: "discovery:///lushop.user.service"
  goods:
    endpoint: "discovery:///lushop.goods.service"
```

### 2. 数据库连接协调
```yaml
# Nacos 配置中的数据库连接
data:
  database:
    driver: mysql
    source: lushop:{password}@tcp(mysql:3306)/lushop_user
  redis:
    addr: redis:6379
    password: "{password}"
```

### 3. 密码管理协调
**K8s Secrets** ↔ **Nacos 配置** ↔ **环境变量**

```bash
# 1. K8s Secrets 中的密码
kubectl get secret mysql-auth -n lushop -o yaml

# 2. Nacos 配置中的占位符
source: lushop:YourDBPasswordHere@tcp(mysql:3306)

# 3. 配置导入时自动替换
./configure-nacos.sh import
```

## 🚀 自动化协调脚本

### 配置导入脚本
```bash
# 自动导入所有配置到 Nacos
./configure-nacos.sh import

# 验证配置导入
./configure-nacos.sh verify
```

### 一键部署脚本
```bash
# 包含完整协调流程的部署
./quick-deploy.sh

# 内部执行顺序:
# 1. deploy_infrastructure
# 2. configure_nacos      ← 新增的配置导入步骤
# 3. deploy_application
# 4. wait_for_services
```

## ⚠️ 常见协调问题及解决方案

### 问题 1: 配置导入失败
**现象**: Nacos 中缺少服务配置
**原因**: 配置导入脚本执行失败或密码不匹配
**解决**:
```bash
# 手动导入配置
kubectl port-forward -n lushop svc/nacos 8848:8848 &
./configure-nacos.sh import
```

### 问题 2: 服务启动失败
**现象**: Pod 状态为 CrashLoopBackOff
**原因**: Nacos 连接失败或配置错误
**解决**:
```bash
# 检查服务日志
kubectl logs -f deployment/user-service -n lushop

# 验证 Nacos 连接
kubectl exec -it deployment/user-service -n lushop -- curl http://nacos:8848/nacos/v1/console/health/readiness
```

### 问题 3: 密码不一致
**现象**: 数据库连接失败
**原因**: K8s Secrets 与 Nacos 配置中的密码不匹配
**解决**:
```bash
# 重新生成 secrets
./gen-secrets-custom.sh

# 重新导入配置
./configure-nacos.sh import

# 重启服务
kubectl rollout restart deployment -n lushop
```

## 📊 配置协调验证

### 验证 Nacos 配置
```bash
# 检查配置是否存在
kubectl port-forward -n lushop svc/nacos 8848:8848 &
curl "http://localhost:8848/nacos/v1/cs/configs?dataId=user.yaml&group=lushop_grpc&tenant=de9c6a0e-1fbc-425d-8d3b-09066fea6889"
```

### 验证服务配置加载
```bash
# 检查服务日志中的配置加载信息
kubectl logs deployment/user-service -n lushop | grep -i "nacos\|config"

# 验证服务间通信
kubectl exec -it deployment/user-service -n lushop -- curl http://goods-service:8012/health
```

## 🔄 配置更新流程

### 生产环境配置更新
```bash
# 1. 更新 Nacos 配置
kubectl port-forward -n lushop svc/nacos 8848:8848 &
# 在 Nacos 控制台更新配置

# 2. 触发服务重启
kubectl rollout restart deployment/user-service -n lushop

# 3. 验证配置生效
kubectl logs -f deployment/user-service -n lushop
```

### 批量配置更新
```bash
# 更新所有服务的配置版本
for service in user goods order inventory userop userauth; do
  kubectl rollout restart deployment/$service-service -n lushop
done
```

## 🎯 配置协调最佳实践

### 1. 配置分层管理
- **本地配置**: 连接信息，不包含敏感数据
- **Nacos 配置**: 完整业务配置，支持动态更新
- **K8s Secrets**: 敏感信息管理

### 2. 环境隔离
- 不同环境使用不同的 Nacos 命名空间
- 配置中的服务名保持一致
- 密码通过环境变量管理

### 3. 监控和告警
- 监控配置加载状态
- 监控服务间通信
- 配置变更审计

### 4. 备份和恢复
- 定期备份 Nacos 配置
- 配置变更走审批流程
- 准备配置回滚方案

## 📚 相关文档

- [配置使用指南](./config-usage.md)
- [部署指南](./README-CN.md)
- [Nacos 配置导入脚本](./configure-nacos.sh)
- [一键部署脚本](./quick-deploy.sh)
