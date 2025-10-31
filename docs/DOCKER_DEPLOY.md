# 🐳 Lushop 微服务 Docker 一键部署指南

## 📋 系统架构

项目采用**分离式部署架构**，所有部署文件位于 `deploy/` 目录：

```
deploy/
├── docker-compose.infrastructure.yml  # 基础设施服务（12个服务）
├── docker-compose.services.yml        # 应用服务（7个服务）
└── scripts/                           # 部署脚本
    ├── deploy-all.sh                  # 一键部署所有服务
    ├── deploy-infrastructure.sh       # 部署基础设施
    └── deploy-services.sh            # 部署应用服务
```

### 架构总览

```
┌─────────────────────────────────────────────────────────────┐
│                    Lushop 微服务系统                         │
│           分离式部署 (基础设施 + 应用服务)                    │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│  基础设施      │    │  微服务        │    │  API 网关     │
│  (12个服务)    │    │  (6个服务)     │    │  (1个服务)    │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
   MySQL/Redis         Goods/Order/etc        HTTP/gRPC
   Consul/Nacos       Inventory/UserAuth     Gateway
   Jaeger/Prometheus  UserOp
   RocketMQ/ES
```

---

## 🚀 快速开始

### 前置条件

1. **Docker** (版本 20.10+)
2. **Docker Compose** (版本 2.0+)
3. **至少 8GB RAM**
4. **至少 20GB 磁盘空间**

### 一键部署（推荐方式）

```bash
# 1. 进入项目目录
cd /home/zzx/GoProject/lushop-kratos-main

# 2. 一键部署所有服务（基础设施 + 应用服务）
./deploy.sh

# 或使用分离式部署脚本
./deploy/scripts/deploy-all.sh
```

**就这么简单！** 🎉

### 分离式部署（生产推荐）

项目采用分离式部署，基础设施和应用服务分开管理：

```bash
# 1. 先部署基础设施（MySQL、Redis、Consul、Nacos等）
./deploy/scripts/deploy-infrastructure.sh

# 2. 等待基础设施就绪（约30秒）
sleep 30

# 3. 再部署应用服务（微服务 + API网关）
./deploy/scripts/deploy-services.sh
```

---

## 📝 部署命令

### 使用部署脚本（推荐）

```bash
# 一键部署所有服务
./deploy.sh
# 或
./deploy/scripts/deploy-all.sh

# 分步部署
./deploy/scripts/deploy-infrastructure.sh  # 基础设施
./deploy/scripts/deploy-services.sh        # 应用服务
```

### Docker Compose 命令

#### 基础设施服务

```bash
# 启动基础设施
docker compose -f deploy/docker-compose.infrastructure.yml up -d

# 停止基础设施
docker compose -f deploy/docker-compose.infrastructure.yml down

# 查看基础设施状态
docker compose -f deploy/docker-compose.infrastructure.yml ps

# 查看基础设施日志
docker compose -f deploy/docker-compose.infrastructure.yml logs -f
```

#### 应用服务

```bash
# 构建并启动应用服务
docker compose -f deploy/docker-compose.services.yml up -d --build

# 停止应用服务
docker compose -f deploy/docker-compose.services.yml down

# 查看应用服务状态
docker compose -f deploy/docker-compose.services.yml ps

# 查看应用服务日志
docker compose -f deploy/docker-compose.services.yml logs -f

# 重启特定服务
docker compose -f deploy/docker-compose.services.yml restart goods-service

# 查看特定服务日志
docker compose -f deploy/docker-compose.services.yml logs -f goods-service
```

---

## 🎯 服务访问地址

### API 网关
- **HTTP API**: http://localhost:8001
- **gRPC API**: localhost:9001

### 基础设施服务（12个）

位于 `deploy/docker-compose.infrastructure.yml`：

| 服务 | 端口 | 说明 | 访问地址 |
|------|------|------|----------|
| **MySQL** | 3306 | 关系型数据库 | `localhost:3306` |
| **Redis** | 6379 | 缓存和分布式锁 | `localhost:6379` |
| **Consul** | 8500 | 服务发现和配置 | http://localhost:8500 |
| **Nacos** | 8848 | 配置中心 | http://localhost:8848/nacos |
| **Jaeger** | 16686 | 链路追踪 | http://localhost:16686 |
| **Prometheus** | 9090 | 监控指标收集 | http://localhost:9090 |
| **Grafana** | 3000 | 监控可视化 | http://localhost:3000 |
| **RocketMQ NameServer** | 9876 | 消息队列命名服务 | `localhost:9876` |
| **RocketMQ Broker** | 10909-10912 | 消息队列代理 | `localhost:10909` |
| **RocketMQ Console** | 8080 | 消息队列管理 | http://localhost:8080 |
| **Elasticsearch** | 9200, 9300 | 搜索引擎 | http://localhost:9200 |
| **Kibana** | 5601 | 搜索可视化 | http://localhost:5601 |

**默认账号密码**：
- MySQL: `root` / `root123456`
- Redis: 密码 `123456`
- Nacos: `nacos` / `nacos` (首次登录需修改)
- Grafana: `admin` / `admin`

### 应用服务（7个）

位于 `deploy/docker-compose.services.yml`：

| 服务 | gRPC端口 | HTTP端口 | 说明 |
|------|----------|----------|------|
| **User Service** | 50051 | - | 用户服务 |
| **UserAuth Service** | 50056 | - | 用户认证服务 |
| **Goods Service** | 50052 | 8000 | 商品服务 |
| **Inventory Service** | 50054 | - | 库存服务 |
| **Order Service** | 50053 | - | 订单服务 |
| **UserOp Service** | 50055 | - | 用户操作服务 |
| **API Gateway** | 9001 | 8001 | API网关 |

---

## 📊 部署流程

### 分离式部署流程

#### 第一步：部署基础设施

```bash
./deploy/scripts/deploy-infrastructure.sh
```

该脚本会：
1. 检查系统依赖（iptables）
2. 配置系统参数（Elasticsearch 需要 `vm.max_map_count=262144`）
3. 自动修复配置文件（创建必要的配置文件和目录）
4. 启动所有基础设施服务（12个）
5. 等待服务健康检查通过（约30秒）

#### 第二步：部署应用服务

```bash
./deploy/scripts/deploy-services.sh
```

该脚本会：
1. 检查基础设施服务是否运行（检查网络和服务状态）
2. 构建并启动所有应用服务（7个）
3. 等待服务启动完成（约20秒）

### 网络架构

所有服务通过 Docker 网络 `lushop-network` 通信：

```
lushop-network (bridge)
├── 基础设施服务（infrastructure.yml）
│   ├── mysql (lushop-mysql)
│   ├── redis (lushop-redis)
│   ├── consul (lushop-consul)
│   └── ...
└── 应用服务（services.yml）
    ├── user-service
    ├── goods-service
    └── ...
```

应用服务通过服务名称访问基础设施服务：
- `lushop-mysql:3306`
- `lushop-redis:6379`
- `lushop-consul:8500`
- 其他服务类似

### 自动化部署流程

```
1. 环境检查
   ├─ 检查 Docker
   ├─ 检查 Docker Compose
   └─ 检查系统资源

2. 启动基础设施（infrastructure.yml）
   ├─ MySQL (等待健康检查)
   ├─ Redis (等待健康检查)
   ├─ Consul (等待健康检查)
   ├─ Nacos (等待 MySQL)
   ├─ Jaeger
   ├─ Prometheus
   ├─ Grafana (等待 Prometheus)
   ├─ RocketMQ NameServer
   ├─ RocketMQ Broker (等待 NameServer)
   ├─ RocketMQ Console (等待 Broker)
   ├─ Elasticsearch
   └─ Kibana (等待 Elasticsearch)

3. 构建并启动应用服务（services.yml）
   ├─ 构建 Goods Service
   ├─ 构建 Order Service
   ├─ 构建 Inventory Service
   ├─ 构建 UserAuth Service
   ├─ 构建 UserOp Service
   ├─ 构建 User Service
   └─ 构建 API Gateway

4. 健康检查
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

修改 `deploy/docker-compose.infrastructure.yml` 或 `deploy/docker-compose.services.yml` 中的环境变量：

```yaml
environment:
  - MYSQL_ROOT_PASSWORD=your_password
  - REDIS_PASSWORD=your_password
```

所有配置文件位于 `deploy/` 目录：
- `deploy/docker-compose.infrastructure.yml` - 基础设施配置
- `deploy/docker-compose.services.yml` - 应用服务配置
- `deploy/mysql/init/` - MySQL 初始化脚本
- `deploy/prometheus/prometheus.yml` - Prometheus 配置

**注意**：
- RocketMQ Broker 配置已改为命令行参数，无需 `broker.conf` 文件
- 所有数据目录统一挂载到 `${DATA_DIR:-/home/zzx/GoProject/lushop-data}/` 目录下
- 可通过设置 `DATA_DIR` 环境变量自定义数据存储路径

---

## 📈 监控和日志

### 查看所有服务日志

```bash
# 查看基础设施日志
docker compose -f deploy/docker-compose.infrastructure.yml logs -f

# 查看应用服务日志
docker compose -f deploy/docker-compose.services.yml logs -f
```

### 查看特定服务日志

```bash
# 查看 Goods 服务日志
docker compose -f deploy/docker-compose.services.yml logs -f goods-service

# 查看最近 100 行日志
docker compose -f deploy/docker-compose.services.yml logs --tail=100 goods-service

# 查看实时日志
docker compose -f deploy/docker-compose.services.yml logs -f --tail=50 goods-service
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

### 查看服务状态

```bash
# 查看基础设施状态
docker compose -f deploy/docker-compose.infrastructure.yml ps

# 查看应用服务状态
docker compose -f deploy/docker-compose.services.yml ps

# 查看所有服务（使用过滤器）
docker ps --filter "name=lushop-"

# 快速查看服务健康状态（格式化输出）
docker ps --filter "name=lushop-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
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
docker compose -f deploy/docker-compose.services.yml logs <service-name>
docker compose -f deploy/docker-compose.infrastructure.yml logs <service-name>

# 查看容器状态
docker compose -f deploy/docker-compose.infrastructure.yml ps
docker compose -f deploy/docker-compose.services.yml ps

# 重启特定服务
docker compose -f deploy/docker-compose.services.yml restart <service-name>
```

### 问题3: 数据库连接失败

```bash
# 检查 MySQL 是否就绪
docker exec lushop-mysql mysqladmin ping -h localhost

# 检查 MySQL 日志
docker compose -f deploy/docker-compose.infrastructure.yml logs mysql

# 重启 MySQL
docker compose -f deploy/docker-compose.infrastructure.yml restart mysql
```

### 问题4: 服务发现失败

```bash
# 检查 Consul 状态
curl http://localhost:8500/v1/health/state/any

# 查看 Consul 日志
docker compose -f deploy/docker-compose.infrastructure.yml logs consul

# 重启 Consul
docker compose -f deploy/docker-compose.infrastructure.yml restart consul
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

### 删除并重新部署基础设施

**重要说明**：基础设施容器可以安全删除和重建，数据不会丢失！

#### 为什么可以安全删除？

所有数据使用 **bind mount** 方式挂载到主机目录，`docker compose down` 只会删除容器，不会删除数据目录。

#### 重新部署方法

**方法1：使用 Docker Compose（简单快速）**
```bash
cd /home/zzx/GoProject/lushop-kratos-main/deploy

# 停止并删除容器（数据会保留）
docker compose -f docker-compose.infrastructure.yml down

# 重新启动所有服务
docker compose -f docker-compose.infrastructure.yml up -d
```

**方法2：使用部署脚本（推荐，包含配置检查）**
```bash
cd /home/zzx/GoProject/lushop-kratos-main
./deploy/scripts/deploy-infrastructure.sh
```

**方法3：使用停止重建脚本（一键操作）**
```bash
cd /home/zzx/GoProject/lushop-kratos-main
./deploy/scripts/stop-and-redeploy.sh
```

#### 数据持久化验证

删除容器后，数据目录仍然存在：
- ✅ MySQL 数据：`${DATA_DIR}/mysql/` - 数据库完整保留
- ✅ Redis 数据：`${DATA_DIR}/redis/` - dump.rdb 文件保留
- ✅ Elasticsearch 数据：`${DATA_DIR}/elasticsearch/` - 索引数据保留
- ✅ 其他服务数据目录均保留

**验证数据是否保留**：
```bash
# 检查 MySQL 数据库
docker exec lushop-mysql mysql -uroot -proot123456 -e "SHOW DATABASES;"

# 检查 Redis 数据
ls -lh /home/zzx/GoProject/lushop-data/redis/*.rdb

# 检查 Elasticsearch 数据
du -sh /home/zzx/GoProject/lushop-data/elasticsearch
```

#### 注意事项

1. **数据安全** ✅
   - `docker compose down` 不会删除 bind mount 目录中的数据
   - 只有使用 `docker compose down -v` 才会删除命名卷（本项目不使用命名卷）

2. **服务启动时间** ⏳
   - Nacos：需要约 1-2 分钟完成初始化
   - Elasticsearch：需要约 30-60 秒启动
   - 其他服务：通常 10-30 秒

3. **可能出现的警告** ⚠️
   - Elasticsearch 容器停止时可能出现僵尸进程警告（不影响功能）
   - Jaeger 健康检查可能显示 unhealthy（但服务实际可用）

**详细测试报告**：请参考 [INFRASTRUCTURE_REDEPLOY_TEST.md](INFRASTRUCTURE_REDEPLOY_TEST.md)

---

## 🔄 更新和维护

### 更新服务

```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建并启动应用服务
docker compose -f deploy/docker-compose.services.yml up -d --build

# 3. 重启基础设施（如需要）
docker compose -f deploy/docker-compose.infrastructure.yml up -d
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

#### 基础设施容器（12个）

| 容器名称 | 镜像 | 端口 | 状态 |
|---------|------|------|------|
| lushop-mysql | mysql:8.0 | 3306 | ✅ |
| lushop-redis | redis:7-alpine | 6379 | ✅ |
| lushop-consul | consul:latest | 8500 | ✅ |
| lushop-nacos | nacos/nacos-server:v2.1.1 | 8848, 9848 | ✅ |
| lushop-jaeger | jaegertracing/all-in-one | 16686 | ✅ |
| lushop-prometheus | prom/prometheus:latest | 9090 | ✅ |
| lushop-grafana | grafana/grafana:latest | 3000 | ✅ |
| lushop-rocketmq-namesrv | apache/rocketmq:4.9.4 | 9876 | ✅ |
| lushop-rocketmq-broker | apache/rocketmq:4.9.4 | 10909-10912 | ✅ |
| lushop-rocketmq-console | apacherocketmq/rocketmq-dashboard:latest | 8080 | ✅ |
| lushop-elasticsearch | elasticsearch:8.11.0 | 9200, 9300 | ✅ |
| lushop-kibana | kibana:8.11.0 | 5601 | ✅ |

#### 应用服务容器（7个）

| 容器名称 | 镜像 | 端口 | 状态 |
|---------|------|------|------|
| lushop-goods-service | goods:latest | 50052 | ✅ |
| lushop-order-service | order:latest | 50053 | ✅ |
| lushop-inventory-service | inventory:latest | 50054 | ✅ |
| lushop-userop-service | userop:latest | 50055 | ✅ |
| lushop-userauth-service | userauth:latest | 50056 | ✅ |
| lushop-user-service | user:latest | 50051 | ✅ |
| lushop-api-gateway | api-gateway:latest | 8001, 9001 | ✅ |

**总计**: 19 个容器（12个基础设施 + 7个应用服务）

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

## ⚠️ 重要注意事项

1. **部署顺序**：必须先启动基础设施，再启动应用服务
2. **网络依赖**：应用服务依赖 `lushop-network` 网络，确保基础设施先创建网络
3. **健康检查**：使用健康检查确保服务就绪后再启动依赖服务
4. **资源配置**：确保服务器有足够资源（至少 8GB RAM，20GB 磁盘）
5. **系统参数**：Elasticsearch 需要配置 `vm.max_map_count=262144`（部署脚本会自动配置）
6. **系统依赖**：需要安装 `iptables`（Ubuntu 24.04 可能需要通过 `update-alternatives` 配置）
7. **数据目录**：所有数据统一存储在 `${DATA_DIR:-/home/zzx/GoProject/lushop-data}/` 目录下，可通过环境变量自定义
8. **权限要求**：部分服务（Prometheus、Elasticsearch、Grafana）需要特定的目录权限，部署脚本会自动处理

## 📚 相关文档

- [../README.md](../README.md) - 项目总览
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - 项目目录结构
- [LUSHOP_TESTING_PLAN.md](LUSHOP_TESTING_PLAN.md) - 测试方案
- [INTERVIEW_GUIDE.md](INTERVIEW_GUIDE.md) - 面试指南
- [INFRASTRUCTURE_REDEPLOY_TEST.md](INFRASTRUCTURE_REDEPLOY_TEST.md) - 基础设施删除重建测试报告

---

## 🎉 总结

使用分离式部署架构部署 Lushop 微服务系统：

✅ **19个容器** - 完整的微服务架构（12个基础设施 + 7个应用服务）  
✅ **一键启动** - `./deploy.sh` 或 `./deploy/scripts/deploy-all.sh`  
✅ **分离部署** - 基础设施和应用服务独立管理  
✅ **自动配置修复** - 部署脚本自动检查并修复配置文件问题  
✅ **自动编排** - Docker Compose 管理  
✅ **健康检查** - 自动监控服务状态  
✅ **数据持久化** - Volume 挂载  
✅ **服务发现** - Consul 自动注册  
✅ **配置中心** - Nacos 配置管理  
✅ **链路追踪** - Jaeger 全链路监控  
✅ **消息队列** - RocketMQ 异步处理  
✅ **监控告警** - Prometheus + Grafana  
✅ **搜索引擎** - Elasticsearch + Kibana  

**一条命令，启动整个系统！** 🚀

---

**创建时间**: 2025-10-25  
**版本**: v1.0  
**状态**: ✅ 完成
