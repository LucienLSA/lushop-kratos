# 🚀 Lushop 快速启动指南

## 一键部署（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/your-repo/lushop-kratos-main.git
cd lushop-kratos-main

# 2. 一键启动
chmod +x deploy.sh
./deploy.sh start
```

**就这么简单！** 🎉

---

## 使用 Makefile

```bash
# 查看所有命令
make help

# 启动所有服务
make start

# 停止所有服务
make stop

# 查看服务状态
make status

# 查看日志
make logs
```

---

## 访问服务

启动成功后，访问以下地址：

| 服务 | 地址 |
|------|------|
| **API 网关** | http://localhost:8001 |
| **Consul** | http://localhost:8500 |
| **Jaeger** | http://localhost:16686 |
| **RocketMQ** | http://localhost:8080 |

---

## 常用命令

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f goods-service

# 重启服务
docker-compose restart goods-service

# 停止所有服务
./deploy.sh stop

# 清理所有数据
./deploy.sh clean
```

---

## 测试

```bash
# 运行集成测试
cd test/integration
./run_integration_tests.sh
```

---

## 详细文档

- [Docker 部署文档](DOCKER_DEPLOY.md)
- [测试方案](LUSHOP_TESTING_PLAN.md)
- [集成测试](test/integration/README.md)

---

**就是这么简单！开始使用 Lushop 吧！** 🚀
