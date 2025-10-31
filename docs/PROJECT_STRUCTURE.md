# 📂 项目目录结构

## 🎯 目录组织原则

项目采用清晰的目录结构，将不同类型的文件分类管理：

- **`deploy/`** - 部署相关配置和脚本
- **`docs/`** - 项目文档
- **`scripts/`** - 开发工具脚本
- **`service/`** - 微服务代码
- **`lushop/`** - API 网关代码
- **`k8s/`** - Kubernetes 配置
- **`interview/`** - 面试准备材料

---

## 📁 详细目录结构

```
lushop-kratos-main/
├── README.md                      # 项目主文档
├── Makefile                       # Make 构建文件
├── deploy.sh                      # 快速部署入口脚本
│
├── deploy/                        # 部署相关
│   ├── docker-compose.infrastructure.yml  # 基础设施配置
│   ├── docker-compose.services.yml        # 应用服务配置
│   ├── scripts/                   # 部署脚本
│   │   ├── deploy-all.sh          # 一键部署
│   │   ├── deploy-infrastructure.sh  # 部署基础设施
│   │   └── deploy-services.sh     # 部署应用服务
│   ├── mysql/                     # MySQL 配置
│   │   └── init/                  # 初始化脚本
│   ├── prometheus/                # Prometheus 配置
│   │   └── prometheus.yml
│   └── nacos/                     # Nacos 配置
│
├── docs/                          # 项目文档
│   ├── DOCKER_DEPLOY.md           # Docker 部署指南（完整）
│   ├── INTERVIEW_GUIDE.md         # 面试指南
│   ├── LUSHOP_TESTING_PLAN.md     # 测试计划
│   └── PROJECT_STRUCTURE.md       # 本文件（项目结构）
│
├── scripts/                       # 开发工具脚本
│   └── regenerate_wire.sh         # Wire 代码生成
│
├── service/                       # 微服务目录
│   ├── user/                      # User Service
│   ├── userauth/                  # UserAuth Service
│   ├── goods/                     # Goods Service
│   ├── inventory/                 # Inventory Service
│   ├── order/                     # Order Service
│   └── userop/                    # UserOp Service
│
├── lushop/                        # API Gateway
│   ├── cmd/                       # 入口文件
│   ├── internal/                 # 内部代码
│   ├── configs/                   # 配置文件
│   └── Dockerfile
│
├── common/                        # 公共代码
│
├── k8s/                           # Kubernetes 配置
│   ├── infrastructure/           # 基础设施
│   └── services/                  # 应用服务
│
├── interview/                     # 面试准备材料
│   └── Q*.md                      # 面试问题
│
└── stress-test/                   # 压力测试
    └── quick_stress_test.sh
```

---

## 📋 主要目录说明

### deploy/ - 部署配置

包含所有部署相关的文件：

- **Docker Compose 配置**
  - `docker-compose.infrastructure.yml` - 基础设施服务（12个服务）
  - `docker-compose.services.yml` - 应用服务（7个服务）

- **部署脚本**
  - `scripts/deploy-all.sh` - 一键部署所有服务
  - `scripts/deploy-infrastructure.sh` - 部署基础设施
  - `scripts/deploy-services.sh` - 部署应用服务

- **服务配置**
  - `mysql/init/` - MySQL 初始化脚本
  - `prometheus/prometheus.yml` - Prometheus 监控配置
  - `nacos/` - Nacos 配置

**注意**：
  - RocketMQ Broker 配置已改为命令行参数，无需 `broker.conf` 文件
  - 所有数据目录统一挂载到 `${DATA_DIR:-/home/zzx/GoProject/lushop-data}/` 目录下

**使用方法**：
```bash
# 从项目根目录执行
./deploy.sh

# 或直接调用脚本
./deploy/scripts/deploy-all.sh

# 使用 docker compose
docker compose -f deploy/docker-compose.infrastructure.yml up -d
docker compose -f deploy/docker-compose.services.yml up -d --build
```

### docs/ - 文档目录

包含所有项目文档：

- `DOCKER_DEPLOY.md` - **Docker 部署详细指南**（推荐用户查看）
- `INFRASTRUCTURE_REDEPLOY_TEST.md` - 基础设施删除重建测试报告
- `INTERVIEW_GUIDE.md` - 面试准备指南
- `LUSHOP_TESTING_PLAN.md` - 测试计划
- `PROJECT_STRUCTURE.md` - 项目结构说明（本文件）

### scripts/ - 开发工具

包含开发过程中使用的工具脚本：

- `regenerate_wire.sh` - 批量重新生成所有服务的 Wire 代码

**使用方法**：
```bash
./scripts/regenerate_wire.sh
```

### service/ - 微服务

每个微服务目录结构：

```
service/[service-name]/
├── cmd/                    # 服务入口
├── internal/              # 内部代码
│   ├── biz/               # 业务逻辑层
│   ├── data/               # 数据访问层
│   └── service/            # 服务层
├── api/                    # API 定义（proto）
├── configs/                # 配置文件
└── Dockerfile             # Docker 镜像构建
```

---

## 🚀 快速导航

### 部署相关
```bash
cd deploy/
ls -la docker-compose*.yml
ls -la scripts/
```

### 文档查阅
```bash
cd docs/
ls -la *.md
```

### 开发工具
```bash
cd scripts/
ls -la *.sh
```

---

## 📝 文件命名规范

### Docker Compose 文件
- `docker-compose.infrastructure.yml` - 基础设施
- `docker-compose.services.yml` - 应用服务

### 脚本文件
- `deploy-*.sh` - 部署脚本
- `regenerate_*.sh` - 代码生成脚本

### 文档文件
- `README.md` - 项目主文档（根目录）
- `*_GUIDE.md` - 指南文档
- `*_PLAN.md` - 计划文档
- `*_STATUS.md` - 状态文档
- `*_REPORT.md` - 报告文档

---

## 🔄 迁移说明

从旧结构迁移到新结构的变化：

| 旧路径 | 新路径 | 说明 |
|--------|--------|------|
| `docker-compose.yml` | `deploy/docker-compose.*.yml` | 分离为两个文件 |
| `docker-compose.infrastructure.yml` | `deploy/docker-compose.infrastructure.yml` | 移动到 deploy/ |
| `docker-compose.services.yml` | `deploy/docker-compose.services.yml` | 移动到 deploy/ |
| `deploy-*.sh` | `deploy/scripts/deploy-*.sh` | 移动到 scripts 子目录 |
| `regenerate_wire.sh` | `scripts/regenerate_wire.sh` | 移动到 scripts/ |
| `*.md` (除 README.md) | `docs/*.md` | 移动到 docs/ |

---

## ✅ 目录组织优势

1. **清晰的职责分离**
   - 部署文件集中管理
   - 文档统一存放
   - 脚本分类明确

2. **易于维护**
   - 相关文件集中在一个目录
   - 减少根目录文件数量
   - 便于查找和更新

3. **标准化结构**
   - 符合常见的项目结构规范
   - 便于新成员理解
   - 便于 CI/CD 集成

---

## 📚 相关文档

- [README.md](../README.md) - 项目主文档
- [DOCKER_DEPLOY.md](DOCKER_DEPLOY.md) - Docker 部署指南（完整用户指南）
- [INFRASTRUCTURE_REDEPLOY_TEST.md](INFRASTRUCTURE_REDEPLOY_TEST.md) - 基础设施删除重建测试报告
- [INTERVIEW_GUIDE.md](INTERVIEW_GUIDE.md) - 面试准备指南

---

**更新时间**: 2025-10-29

