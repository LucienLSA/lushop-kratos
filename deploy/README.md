# Docker 部署指南

本目录包含 lushop-kratos 项目的完整 Docker 部署配置，支持一键部署所有基础设施服务。

## 📋 目录结构

```
deploy/
├── consul/consul.sh              # Consul 服务发现
├── elasticsearch/es.sh           # Elasticsearch 搜索引擎
├── grafana/grafana.sh            # Grafana 可视化面板
├── jaeger/jaeger.sh              # Jaeger 链路追踪
├── kibana/kibana.sh              # Kibana 日志分析
├── mysql/mysql.sh                # MySQL 数据库
├── nacos/nacos.sh                # Nacos 配置中心
├── prometheus/prometheus.sh      # Prometheus 监控
├── redis/redis.sh                # Redis 缓存
├── rocketmq/                     # RocketMQ 消息队列
│   ├── docker-compose.yaml
│   └── pre.sh
├── deploy-all.sh                 # 一键部署脚本 ⭐
├── env-template.sh               # 环境变量模板
├── ENVIRONMENT_VARIABLES_REPORT.md # 环境变量配置报告
└── README.md                     # 本文件
```

## 🚀 快速开始

### 1. 环境准备

确保系统已安装 Docker 和 Docker Compose：

```bash
# 检查 Docker
docker --version
docker compose version

# 启动 Docker 服务
sudo systemctl start docker
```

### 2. 配置环境变量

```bash
# 生成环境变量模板
bash deploy/env-template.sh > .env

# 编辑配置文件（修改所有密码）
vim .env

# 如果遇到格式错误，可以使用修复脚本
./deploy/fix-env.sh

# 示例配置
REGISTRY=crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6
MYSQL_ROOT_PASSWORD=your_secure_password
REDIS_PASSWORD=your_redis_password
# ... 其他配置
```

### 3. 一键部署

```bash
# 启动所有服务
./deploy/deploy-all.sh start

# 查看服务状态
./deploy/deploy-all.sh status

# 停止所有服务
./deploy/deploy-all.sh stop
```

### 4. 验证部署

部署完成后，所有服务将自动启动并配置。访问地址：

| 服务 | 访问地址 | 说明 |
|------|----------|------|
| MySQL | `localhost:3306` | 数据库 |
| Redis | `localhost:6379` | 缓存 |
| Consul | `http://localhost:8500` | 服务发现 |
| Nacos | `http://localhost:8848` | 配置中心 |
| Elasticsearch | `http://localhost:9200` | 搜索引擎 |
| Kibana | `http://localhost:5601` | 日志分析 |
| Grafana | `http://localhost:3000` | 可视化面板 |
| Prometheus | `http://localhost:9090` | 监控指标 |
| Jaeger | `http://localhost:16686` | 链路追踪 |
| RocketMQ Dashboard | `http://localhost:8682` | 消息队列管理 |

## 🔧 手动部署

如果需要单独部署某个服务：

```bash
# 部署 MySQL
./deploy/mysql/mysql.sh

# 部署 Redis
./deploy/redis/redis.sh

# 部署 RocketMQ
cd deploy/rocketmq
bash pre.sh
docker-compose up -d
```

## ⚙️ 配置说明

### 环境变量

所有脚本都支持通过环境变量进行配置：

```bash
# 使用环境变量
MYSQL_PORT=3307 ./deploy/mysql/mysql.sh

# 或从 .env 文件加载
./deploy/deploy-all.sh start
```

### 数据持久化

所有服务数据默认存储在 `~/lushop-data/` 目录：

```
~/lushop-data/
├── mysql/           # MySQL 数据
├── redis/           # Redis 数据和配置
├── nacos/           # Nacos 配置和数据
├── elasticsearch/   # ES 数据和配置
├── grafana/         # Grafana 数据
├── prometheus/      # Prometheus 数据和配置
└── rocketmq/        # RocketMQ 数据
```

### 端口配置

默认端口分配（可通过环境变量修改）：

- **应用服务**: 8001-8999
- **基础设施**: 3306, 6379, 8500, 8848, 9200, 5601, 3000, 9090, 16686
- **RocketMQ**: 9876, 10909-10912, 18680-18681, 8682

## 🔍 故障排查

### 环境变量格式错误

**问题**: `-Xmx512m: command not found` 或类似错误

**原因**: `.env` 文件中包含特殊字符的变量值未用引号包围

**解决方法**:
```bash
# 使用修复脚本
./deploy/fix-env.sh

# 或手动修复
vim .env
# 确保包含空格的值用引号包围，如：
# ES_JAVA_OPTS="-Xms512m -Xmx512m"
```

### 查看服务状态

```bash
# 查看所有容器
docker ps -a

# 查看服务日志
docker logs lushop-mysql
docker logs lushop-redis

# 查看 RocketMQ 日志
cd deploy/rocketmq && docker-compose logs
```

### 常见问题

1. **端口冲突**
   ```bash
   # 检查端口占用
   netstat -tulpn | grep :3306

   # 修改端口
   MYSQL_PORT=3307 ./deploy/mysql/mysql.sh
   ```

2. **权限问题**
   ```bash
   # 数据目录权限
   sudo chown -R $USER:$USER ~/lushop-data
   ```

3. **镜像拉取失败**
   ```bash
   # 检查镜像是否存在
   docker images | grep lushop

   # 重新拉取
   docker pull your-registry/mysql:8.0
   ```

## 📊 监控和维护

### 服务监控

```bash
# 查看资源使用
docker stats

# 查看磁盘使用
du -sh ~/lushop-data/*
```

### 备份和恢复

```bash
# MySQL 备份
docker exec lushop-mysql mysqldump -uroot -pPASSWORD lushop > backup.sql

# 数据目录备份
tar -czf lushop-data-backup.tar.gz ~/lushop-data/
```

### 日志管理

```bash
# 查看所有服务日志
for service in mysql redis consul nacos elasticsearch kibana grafana prometheus jaeger; do
  echo "=== $service ==="
  docker logs "lushop-$service" --tail=20
done
```

## 🔒 安全配置

### 生产环境建议

1. **修改默认密码**: 编辑 `.env` 文件，设置强密码
2. **网络隔离**: 使用 Docker 网络或防火墙限制访问
3. **TLS 加密**: 配置服务间 TLS 通信
4. **定期更新**: 保持镜像和系统更新

### 环境变量安全

```bash
# 不要提交 .env 文件到版本控制
echo ".env" >> .gitignore

# 使用加密的密码
MYSQL_ROOT_PASSWORD="$(openssl rand -base64 32)"
```

## 📚 相关文档

- [环境变量配置完整性报告](ENVIRONMENT_VARIABLES_REPORT.md)
- [镜像兼容性报告](../deploy/IMAGE_COMPATIBILITY_REPORT.md)
- [项目主文档](../README.md)

## 🎯 部署检查清单

- [ ] Docker 和 Docker Compose 已安装
- [ ] `.env` 文件已配置（密码已修改）
- [ ] 数据目录权限正确
- [ ] 端口未被占用
- [ ] 网络连接正常
- [ ] 镜像可正常拉取
- [ ] 部署脚本执行成功
- [ ] 所有服务状态正常
- [ ] 可通过浏览器访问管理界面

---

**部署完成后，访问 http://localhost:8848 登录 Nacos 控制台开始配置应用服务！** 🎉
