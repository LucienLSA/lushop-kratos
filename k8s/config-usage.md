# Lushop K8s 配置文件使用指南

## 📋 问题说明

原始的服务配置文件使用本地地址和固定IP（如 `127.0.0.1`、`192.168.185.128`），这在 K8s 环境中无法工作。

## ✅ 已创建的 K8s 友好配置文件

### 各服务的配置文件对应关系：

| 服务 | 原始配置文件 | K8s 配置文件 |
|------|-------------|-------------|
| user | `service/user/configs/nacos-config.yaml` | `service/user/configs/nacos-config-k8s.yaml` |
| goods | `service/goods/configs/nacos-config.yaml` | `service/goods/configs/nacos-config-k8s.yaml` |
| order | `service/order/configs/nacos-config.yaml` | `service/order/configs/nacos-config-k8s.yaml` |
| inventory | `service/inventory/configs/nacos-config.yaml` | `service/inventory/configs/nacos-config-k8s.yaml` |
| userop | `service/userop/configs/nacos-config.yaml` | `service/userop/configs/nacos-config-k8s.yaml` |
| userauth | `service/userauth/configs/nacosRemote.yaml` | `service/userauth/configs/nacosRemote-k8s.yaml` |
| gateway | `lushop/configs/nacos-config.yaml` | `lushop/configs/nacos-config-k8s.yaml` |

## 🔧 K8s 配置的主要修改

### 1. 服务地址替换
- **MySQL**: `127.0.0.1:3306` → `mysql:3306`
- **Redis**: `127.0.0.1:6379` → `redis:6379`
- **Consul**: `127.0.0.1:8500` → `consul:8500`
- **Jaeger**: `127.0.0.1:14268` → `jaeger:14268`
- **Elasticsearch**: `192.168.185.128:9200` → `elasticsearch:9200`
- **RocketMQ**: `127.0.0.1:9876` → `rocketmq:9876`

### 2. 密码占位符
将硬编码密码替换为占位符：
- `YourDBPasswordHere` - MySQL 密码
- `YourRedisPasswordHere` - Redis 密码

## 🚀 使用 K8s 配置文件的步骤

### 1. 更新 Nacos 配置
在 Nacos 控制台中为每个服务创建新配置，使用对应的 `-k8s.yaml` 文件内容。

### 2. 更新 Secrets 中的密码
确保 `gen-secrets-custom.sh` 中设置的密码与配置文件中的占位符匹配：

```bash
# 编辑 gen-secrets-custom.sh 中的密码
MYSQL_PASSWORD="${MYSQL_PASSWORD:-lushopDb@123}"
REDIS_PASSWORD="${REDIS_PASSWORD:-Redis@123}"

# 运行生成
./gen-secrets-custom.sh
```

### 3. 部署验证
```bash
# 部署后检查服务是否能连接到依赖
kubectl logs -f deployment/user-service -n lushop
kubectl logs -f deployment/goods-service -n lushop
```

## 📊 服务端口映射

| 服务 | HTTP 端口 | gRPC 端口 | 服务名 |
|------|----------|----------|--------|
| Gateway | 8001 | 9001 | lushop-gateway |
| User | 8011 | 50051 | lushop.user.service |
| Goods | 8012 | 50052 | lushop.goods.service |
| Order | 8013 | 50053 | lushop.order.service |
| Inventory | 8014 | 50054 | lushop.inventory.service |
| UserOp | 8015 | 50055 | lushop.userop.service |
| UserAuth | 8016 | 50056 | lushop.userauth.service |

## ⚠️ 注意事项

1. **密码一致性**: 确保 Nacos 配置中的密码与 K8s Secrets 中的值完全匹配
2. **服务发现**: 所有服务都使用 `discovery://` 前缀进行服务发现
3. **数据库初始化**: 部署前需要初始化各服务的数据库表
4. **网络策略**: 确保 K8s 网络策略允许服务间的通信

## 🔍 故障排查

### 连接失败
```bash
# 检查服务是否能解析其他服务名
kubectl exec -it deployment/user-service -n lushop -- nslookup mysql
kubectl exec -it deployment/user-service -n lushop -- nslookup redis
```

### 配置问题
```bash
# 检查 Nacos 配置是否正确加载
kubectl logs -f deployment/user-service -n lushop | grep -i nacos
```

### 数据库连接
```bash
# 测试数据库连接
kubectl exec -it deployment/mysql -n lushop -- mysql -u lushop -p -e "SHOW DATABASES;"
```
