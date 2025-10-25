# 🐳 Lushop 微服务 Docker 一键部署指南

## 📋 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Lushop 微服务系统                         │
│                  Docker Compose 一键部署                     │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│  基础设施      │    │  微服务        │    │  API 网关     │
│  (7个容器)     │    │  (6个容器)     │    │  (1个容器)    │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
   MySQL/Redis         Goods/Order/etc        HTTP/gRPC
   Consul/Jaeger       Inventory/UserAuth     Gateway
   RocketMQ            UserOp
```

---

## 🚀 快速开始

### 前置条件

1. **Docker** (版本 20.10+)
2. **Docker Compose** (版本 2.0+)
3. **至少 8GB RAM**
4. **至少 20GB 磁盘空间**

### 一键部署

```bash
# 1. 进入项目目录
cd /home/zzx/GoProject/lushop-kratos-main

# 2. 赋予执行权限
chmod +x deploy.sh

# 3. 一键启动所有服务
./deploy.sh start
```

**就这么简单！** 🎉

---

## 📝 部署命令

### 基本命令

```bash
# 启动所有服务
./deploy.sh start

# 停止所有服务
./deploy.sh stop

# 重启所有服务
./deploy.sh restart

# 查看服务日志
./deploy.sh logs

# 查看服务状态
./deploy.sh status

# 清理所有容器和数据
./deploy.sh clean
```

### Docker Compose 命令

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f

# 重启特定服务
docker-compose restart goods-service

# 查看特定服务日志
docker-compose logs -f goods-service
```

---

## 🎯 服务访问地址

### API 网关
- **HTTP API**: http://localhost:8001
- **gRPC API**: localhost:9001

### 管理界面
| 服务 | 地址 | 说明 |
|------|------|------|
| **Consul** | http://localhost:8500 | 服务发现和配置中心 |
| **Jaeger** | http://localhost:16686 | 分布式链路追踪 |
| **RocketMQ Console** | http://localhost:8080 | 消息队列管理 |

### 微服务端口
| 服务 | gRPC 端口 | 说明 |
|------|-----------|------|
| **Goods Service** | 50052 | 商品服务 |
| **Order Service** | 50053 | 订单服务 |
| **Inventory Service** | 50054 | 库存服务 |
| **UserOp Service** | 50055 | 用户操作服务 |
| **UserAuth Service** | 50056 | 用户认证服务 |

### 基础设施
| 服务 | 端口 | 账号 | 密码 |
|------|------|------|------|
| **MySQL** | 3306 | root | root123456 |
| **Redis** | 6379 | - | 123456 |
| **RocketMQ NameServer** | 9876 | - | - |

---

## 📊 部署流程

### 自动化部署流程

```
1. 环境检查
   ├─ 检查 Docker
   ├─ 检查 Docker Compose
   └─ 检查系统资源

2. 清理旧容器
   ├─ 停止运行中的容器
   └─ 删除旧数据卷

3. 构建镜像
   ├─ 构建 Goods Service
   ├─ 构建 Order Service
   ├─ 构建 Inventory Service
   ├─ 构建 UserAuth Service
   ├─ 构建 UserOp Service
   └─ 构建 API Gateway

4. 启动基础设施
   ├─ MySQL (等待健康检查)
   ├─ Redis (等待健康检查)
   ├─ Consul (等待健康检查)
   ├─ Jaeger
   ├─ RocketMQ NameServer
   ├─ RocketMQ Broker
   └─ RocketMQ Console

5. 启动微服务
   ├─ Goods Service
   ├─ Order Service
   ├─ Inventory Service
   ├─ UserAuth Service
   └─ UserOp Service

6. 启动 API 网关
   └─ API Gateway

7. 健康检查
   └─ 验证所有服务状态
```

---

## 🔧 配置说明

### 环境变量

所有服务的配置文件位于 `service/*/configs/` 目录下。

**MySQL 配置**:
```yaml
database:
  driver: mysql
  source: root:root123456@tcp(lushop-mysql:3306)/lushop?charset=utf8mb4&parseTime=True&loc=Local
```

**Redis 配置**:
```yaml
redis:
  addr: lushop-redis:6379
  password: "123456"
  db: 0
```

**Consul 配置**:
```yaml
registry:
  consul:
    address: lushop-consul:8500
    scheme: http
```

### 自定义配置

修改 `docker-compose.yml` 中的环境变量：

```yaml
environment:
  - MYSQL_ROOT_PASSWORD=your_password
  - REDIS_PASSWORD=your_password
```

---

## 📈 监控和日志

### 查看所有服务日志
```bash
./deploy.sh logs
```

### 查看特定服务日志
```bash
# 查看 Goods 服务日志
docker-compose logs -f goods-service

# 查看最近 100 行日志
docker-compose logs --tail=100 goods-service

# 查看实时日志
docker-compose logs -f --tail=50 goods-service
```

### 进入容器调试
```bash
# 进入 Goods 服务容器
docker exec -it lushop-goods-service sh

# 进入 MySQL 容器
docker exec -it lushop-mysql mysql -uroot -proot123456

# 进入 Redis 容器
docker exec -it lushop-redis redis-cli -a 123456
```

---

## 🐛 故障排查

### 问题1: 端口被占用
```bash
# 查找占用端口的进程
sudo lsof -i :8001
sudo lsof -i :3306

# 停止占用端口的进程
kill -9 <PID>
```

### 问题2: 容器启动失败
```bash
# 查看容器日志
docker-compose logs <service-name>

# 查看容器状态
docker-compose ps

# 重启特定服务
docker-compose restart <service-name>
```

### 问题3: 数据库连接失败
```bash
# 检查 MySQL 是否就绪
docker exec lushop-mysql mysqladmin ping -h localhost

# 检查 MySQL 日志
docker-compose logs mysql

# 重启 MySQL
docker-compose restart mysql
```

### 问题4: 服务发现失败
```bash
# 检查 Consul 状态
curl http://localhost:8500/v1/health/state/any

# 查看 Consul 日志
docker-compose logs consul

# 重启 Consul
docker-compose restart consul
```

### 问题5: 内存不足
```bash
# 查看容器资源使用
docker stats

# 限制容器内存
# 在 docker-compose.yml 中添加:
deploy:
  resources:
    limits:
      memory: 512M
```

---

## 🔄 更新和维护

### 更新服务
```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建镜像
docker-compose build

# 3. 重启服务
docker-compose up -d
```

### 备份数据
```bash
# 备份 MySQL 数据
docker exec lushop-mysql mysqldump -uroot -proot123456 lushop > backup.sql

# 备份 Redis 数据
docker exec lushop-redis redis-cli -a 123456 SAVE
docker cp lushop-redis:/data/dump.rdb ./redis-backup.rdb
```

### 恢复数据
```bash
# 恢复 MySQL 数据
docker exec -i lushop-mysql mysql -uroot -proot123456 lushop < backup.sql

# 恢复 Redis 数据
docker cp ./redis-backup.rdb lushop-redis:/data/dump.rdb
docker-compose restart redis
```

---

## 📦 容器列表

| 容器名称 | 镜像 | 端口 | 状态 |
|---------|------|------|------|
| lushop-mysql | mysql:8.0 | 3306 | ✅ |
| lushop-redis | redis:7-alpine | 6379 | ✅ |
| lushop-consul | consul:latest | 8500 | ✅ |
| lushop-jaeger | jaegertracing/all-in-one | 16686 | ✅ |
| lushop-rocketmq-namesrv | apache/rocketmq:4.9.4 | 9876 | ✅ |
| lushop-rocketmq-broker | apache/rocketmq:4.9.4 | 10909-10912 | ✅ |
| lushop-rocketmq-console | apacherocketmq/rocketmq-dashboard | 8080 | ✅ |
| lushop-goods-service | goods:latest | 50052 | ✅ |
| lushop-order-service | order:latest | 50053 | ✅ |
| lushop-inventory-service | inventory:latest | 50054 | ✅ |
| lushop-userop-service | userop:latest | 50055 | ✅ |
| lushop-userauth-service | userauth:latest | 50056 | ✅ |
| lushop-api-gateway | api-gateway:latest | 8001, 9001 | ✅ |

**总计**: 14 个容器

---

## 🎯 性能优化

### 1. 调整容器资源限制
```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
    reservations:
      memory: 256M
```

### 2. 使用多阶段构建
```dockerfile
FROM golang:1.21-alpine AS builder
# 构建阶段

FROM alpine:latest
# 运行阶段
```

### 3. 启用健康检查
```yaml
healthcheck:
  test: ["CMD", "nc", "-z", "localhost", "50052"]
  interval: 30s
  timeout: 3s
  retries: 3
```

---

## 🔐 安全建议

1. **修改默认密码**
   - MySQL root 密码
   - Redis 密码

2. **使用环境变量**
   ```bash
   export MYSQL_ROOT_PASSWORD=your_secure_password
   export REDIS_PASSWORD=your_secure_password
   ```

3. **限制网络访问**
   - 只暴露必要的端口
   - 使用防火墙规则

4. **定期更新镜像**
   ```bash
   docker-compose pull
   docker-compose up -d
   ```

---

## 📚 相关文档

- [README.md](README.md) - 项目总览
- [TESTING_ROADMAP.md](TESTING_ROADMAP.md) - 测试路线图
- [LUSHOP_TESTING_PLAN.md](LUSHOP_TESTING_PLAN.md) - 测试方案
- [test/integration/README.md](test/integration/README.md) - 集成测试

---

## 🎉 总结

使用 Docker Compose 一键部署 Lushop 微服务系统：

✅ **14个容器** - 完整的微服务架构  
✅ **一键启动** - `./deploy.sh start`  
✅ **自动编排** - Docker Compose 管理  
✅ **健康检查** - 自动监控服务状态  
✅ **数据持久化** - Volume 挂载  
✅ **服务发现** - Consul 自动注册  
✅ **链路追踪** - Jaeger 全链路监控  
✅ **消息队列** - RocketMQ 异步处理  

**一条命令，启动整个系统！** 🚀

---

**创建时间**: 2025-10-25  
**版本**: v1.0  
**状态**: ✅ 完成
