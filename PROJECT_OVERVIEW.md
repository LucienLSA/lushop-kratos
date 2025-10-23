# Lushop 项目全景概览

**生成时间**: 2025-10-23  
**项目类型**: 基于 Kratos 的微服务电商平台  
**架构模式**: API 网关 + 微服务 + 服务发现

---

## 📊 项目统计

| 指标 | 数量 |
|------|------|
| 微服务数量 | 6 个 |
| API 网关 | 1 个 |
| 配置管理 | Nacos + Consul |
| 通信协议 | gRPC + HTTP |
| 数据存储 | MySQL + Redis |

---

## 🏗️ 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                         客户端                               │
└────────────────────┬────────────────────────────────────────┘
                     │ HTTP
                     ↓
┌─────────────────────────────────────────────────────────────┐
│              API 网关 (lushop)                               │
│              HTTP: 8001  gRPC: 9001                         │
└──┬────────┬────────┬────────┬────────┬────────┬────────────┘
   │        │        │        │        │        │
   │ gRPC   │ gRPC   │ gRPC   │ gRPC   │ gRPC   │ gRPC
   ↓        ↓        ↓        ↓        ↓        ↓
┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌──────────┐
│User │ │Goods│ │Order│ │Inv. │ │UserOp│ │UserAuth  │
│:51  │ │:52  │ │:53  │ │:54  │ │:55  │ │:50056    │
└─────┘ └─────┘ └─────┘ └─────┘ └─────┘ └──────────┘
   │        │        │                      │
   ↓        ↓        ↓                      ↓
┌─────┐ ┌─────┐ ┌─────┐                ┌─────┐
│MySQL│ │MySQL│ │MySQL│                │Redis│
└─────┘ └─────┘ └─────┘                └─────┘
```

---

## 📁 目录结构

```
lushop-kratos-main/
├── lushop/                      # API 网关服务
│   ├── api/                     # Proto 定义
│   │   ├── lushop/v1/          # 网关 HTTP API
│   │   ├── service/            # 后端服务 Proto（客户端）
│   │   │   ├── goods/v1/
│   │   │   ├── order/v1/
│   │   │   └── user/v1/
│   │   └── userauth/v1/        # 认证服务 Proto（客户端）
│   ├── cmd/lushop/             # 主程序入口
│   ├── configs/                # 配置文件
│   │   ├── config.yaml         # Nacos 配置
│   │   └── nacosRemote.yaml    # 实际配置（从 Nacos 拉取）
│   ├── internal/
│   │   ├── biz/                # 业务逻辑层
│   │   │   ├── user.go         # 用户业务
│   │   │   ├── user_auth_adapter.go  # 认证服务适配器
│   │   │   └── http/           # HTTP 业务（商品等）
│   │   ├── data/               # 数据访问层
│   │   │   ├── data.go         # Data 层配置
│   │   │   ├── user.go         # ⚠️ 已废弃（Redis 版本）
│   │   │   └── user_auth_repo_grpc.go  # ✅ 当前使用（gRPC 版本）
│   │   ├── service/            # HTTP 服务层
│   │   │   └── user.go         # 用户 HTTP 接口
│   │   ├── server/             # 服务器配置
│   │   │   ├── http.go         # HTTP 服务器
│   │   │   └── grpc.go         # gRPC 服务器
│   │   ├── pkg/                # 工具包
│   │   │   ├── middleware/     # 中间件（JWT 认证等）
│   │   │   └── util/
│   │   └── task/               # 定时任务
│   ├── bin/                    # 编译产物
│   │   └── lushop              # ✅ 可执行文件
│   └── third_party/            # Proto 依赖
│
├── service/                     # 微服务目录
│   ├── user/                   # 用户服务（gRPC: 50051）
│   │   └── （待实现）
│   ├── goods/                  # 商品服务（gRPC: 50052）
│   │   ├── configs/
│   │   │   ├── config.yaml     # Nacos 配置
│   │   │   └── nacos-config.yaml
│   │   └── （其他文件）
│   ├── order/                  # 订单服务（gRPC: 50053）
│   │   └── （完整实现）
│   ├── inventory/              # 库存服务（gRPC: 50054）
│   │   └── （待实现）
│   ├── userop/                 # 用户操作服务（gRPC: 50055）
│   │   └── （完整实现）
│   └── userauth-service/       # ✨ 用户认证服务（gRPC: 50055）
│       ├── api/userauth/v1/    # Proto 定义
│       ├── cmd/userauth/       # 主程序入口
│       ├── configs/
│       │   └── config.yaml     # 配置文件
│       ├── internal/
│       │   ├── biz/            # 业务逻辑
│       │   │   └── auth.go     # 认证业务
│       │   ├── data/           # 数据访问
│       │   │   └── auth.go     # Redis 操作
│       │   ├── service/        # gRPC 服务实现
│       │   │   └── userauth.go
│       │   ├── server/         # 服务器配置
│       │   │   ├── grpc.go
│       │   │   └── registry.go # Consul 注册
│       │   ├── conf/           # 配置结构
│       │   └── pkg/            # 辅助包
│       │       ├── captcha/    # 验证码生成
│       │       ├── sms/        # 短信发送
│       │       └── auth/       # JWT 工具
│       ├── bin/
│       │   └── userauth        # ✅ 可执行文件
│       ├── Makefile
│       ├── start.sh            # 启动脚本
│       └── README.md
│
├── common/                      # 公共代码
│   └── middleware/             # 公共中间件
│
├── deploy/                      # 部署配置
│   └── nacos/                  # Nacos 配置文件
│
├── README.md                    # 项目说明
├── PROJECT_STRUCTURE.md         # 项目结构详解
├── PROJECT_OVERVIEW.md          # 本文件
└── MIGRATION_SUMMARY.md         # 迁移总结
```

---

## 🎯 服务详情

### 1. API 网关 (lushop)

**端口**: HTTP 8001, gRPC 9001  
**职责**: 
- 对外提供 HTTP REST API
- 路由转发到后端微服务
- JWT 认证与鉴权
- 请求聚合与编排

**关键文件**:
- `internal/biz/user.go` - 用户业务逻辑
- `internal/biz/user_auth_adapter.go` - 认证服务适配器（封装 gRPC 调用）
- `internal/data/user_auth_repo_grpc.go` - gRPC 版本仓库（当前使用）
- `internal/service/user.go` - HTTP 接口实现

**配置管理**: Nacos
- 配置组: `lushop_http`
- DataId: `lushop.yaml`

**依赖服务**:
- User Service (用户资料)
- Goods Service (商品)
- Order Service (订单)
- UserAuth Service (认证域) ✨

---

### 2. UserAuth Service (用户认证服务) ✨

**端口**: gRPC 50056  
**职责**:
- 图形验证码生成与校验
- 短信验证码发送与校验（集成阿里云）
- JWT Token 签发、刷新、撤销
- 用户黑名单管理

**技术栈**:
- gRPC 服务端
- Redis 存储
- base64Captcha 验证码
- 阿里云短信 SDK
- golang-jwt

**服务发现**: Consul
- 服务名: `user-auth-service`
- 健康检查: 启用

**API 接口**:
```protobuf
service UserAuth {
  rpc GetCaptcha(Empty) returns (CaptchaReply);
  rpc VerifyCaptcha(VerifyCaptchaReq) returns (VerifyCaptchaReply);
  rpc SendSms(SendSmsReq) returns (SendSmsReply);
  rpc VerifySms(VerifySmsReq) returns (VerifySmsReply);
  rpc IssueToken(IssueTokenReq) returns (TokenReply);
  rpc RefreshToken(RefreshTokenReq) returns (TokenReply);
  rpc RevokeToken(RevokeTokenReq) returns (Empty);
  rpc AddToBlacklist(AddToBlacklistReq) returns (Empty);
  rpc CheckBlacklist(CheckBlacklistReq) returns (CheckBlacklistReply);
}
```

**Redis 键设计**:
```
captcha:{captcha_id}           # 验证码答案，TTL 5分钟
sms_code:{mobile}              # 短信验证码，TTL 5分钟
sms_cooldown:{mobile}          # 短信冷却，TTL 60秒
user_access_token:{user_id}    # Access Token，TTL 30分钟
user_refresh_token:{user_id}   # Refresh Token，TTL 7天
user_blacklist:{user_id}       # 黑名单，TTL 30分钟
```

---

### 3. User Service (用户服务)

**端口**: gRPC 50051  
**职责**: 用户资料管理（CRUD）  
**状态**: 待实现

---

### 4. Goods Service (商品服务)

**端口**: gRPC 50052  
**职责**: 商品管理、分类、库存查询  
**配置**: Nacos (group: `lushop_grpc`, dataId: `goods.yaml`)

---

### 5. Order Service (订单服务)

**端口**: gRPC 50053  
**职责**: 订单创建、查询、状态管理  
**状态**: 已实现

---

### 6. Inventory Service (库存服务)

**端口**: gRPC 50054  
**职责**: 库存管理、扣减、回滚  
**状态**: 待实现

---

### 7. UserOp Service (用户操作服务)

**端口**: gRPC 50055  
**职责**: 用户行为记录、收藏、浏览历史  
**状态**: 已实现  
**⚠️ 注意**: 与 UserAuth Service 端口冲突，需要调整

---

## 🔄 服务调用链

### 用户注册流程

```
客户端
  ↓ POST /api/user/register
API 网关 (lushop)
  ├─→ UserAuth.GetCaptcha()        # 获取验证码
  ├─→ UserAuth.VerifyCaptcha()     # 校验验证码
  ├─→ User.CreateUser()            # 创建用户
  └─→ UserAuth.IssueToken()        # 签发 Token
```

### 用户登录流程

```
客户端
  ↓ POST /api/user/login
API 网关 (lushop)
  ├─→ UserAuth.VerifyCaptcha()     # 校验验证码
  ├─→ User.GetUserByMobile()       # 获取用户信息
  ├─→ User.CheckPassword()         # 校验密码
  ├─→ UserAuth.CheckBlacklist()    # 检查黑名单
  └─→ UserAuth.IssueToken()        # 签发 Token
```

### 用户登出流程

```
客户端
  ↓ POST /api/user/logout (带 JWT)
API 网关 (lushop)
  ├─→ JWT 中间件解析 Token
  ├─→ UserAuth.CheckBlacklist()    # 检查黑名单
  └─→ UserAuth.RevokeToken()       # 撤销 Token + 加入黑名单
```

---

## 🔧 配置管理

### Nacos 配置

**Nacos 地址**: `127.0.0.1:8848`  
**命名空间**: `de9c6a0e-1fbc-425d-8d3b-09066fea6889`

| 服务 | Group | DataId |
|------|-------|--------|
| API 网关 | lushop_http | lushop.yaml |
| Goods Service | lushop_grpc | goods.yaml |
| Order Service | lushop_grpc | order.yaml |
| UserOp Service | lushop_grpc | userop.yaml |

### Consul 配置

**Consul 地址**: `127.0.0.1:8500`  
**用途**: 服务发现与注册

**已注册服务**:
- `user-auth-service` (UserAuth Service)
- `lushop.user.service` (User Service)
- `lushop.goods.service` (Goods Service)

---

## 🗄️ 数据存储

### MySQL 数据库

| 数据库 | 用途 | 服务 |
|--------|------|------|
| lushop_user | 用户数据 | User Service |
| lushop_goods | 商品数据 | Goods Service |
| lushop_order | 订单数据 | Order Service |

### Redis

| 用途 | 服务 |
|------|------|
| Token 存储 | UserAuth Service |
| 验证码存储 | UserAuth Service |
| 短信验证码 | UserAuth Service |
| 黑名单 | UserAuth Service |
| 缓存 | 各服务 |

---

## 🚀 快速启动

### 1. 启动基础设施

```bash
# MySQL
sudo systemctl start mysql

# Redis
sudo systemctl start redis

# Nacos（可选，如使用本地配置可跳过）
sh nacos/bin/startup.sh -m standalone

# Consul
consul agent -dev
```

### 2. 启动 UserAuth 服务

```bash
cd service/userauth-service
./start.sh
```

### 3. 启动 API 网关

```bash
cd lushop
go run ./cmd/lushop -conf ./configs/config.yaml
```

### 4. 验证部署

```bash
# 检查 Consul 服务注册
curl http://127.0.0.1:8500/v1/catalog/service/user-auth-service

# 测试验证码接口
curl http://127.0.0.1:8001/api/user/captcha
```

---

## 📝 开发指南

### 添加新的 API 接口

1. **定义 Proto** (`lushop/api/lushop/v1/lushop.proto`)
2. **生成代码** (`make api`)
3. **实现 Service** (`lushop/internal/service/`)
4. **实现 Biz** (`lushop/internal/biz/`)
5. **实现 Data** (`lushop/internal/data/`)
6. **Wire 装配** (`lushop/cmd/lushop/wire.go`)

### 添加新的微服务

1. **创建服务目录** (`service/newservice/`)
2. **定义 Proto** (`api/newservice/v1/`)
3. **实现服务端** (`internal/`)
4. **配置 Consul 注册**
5. **在网关添加客户端调用**

---

## ⚠️ 已知问题

### 1. 端口冲突 ✅ 已解决

**原问题**: UserOp Service 和 UserAuth Service 都使用 50055 端口  
**解决方案**: UserAuth Service 已改为 50056 端口  
**状态**: ✅ 已修复，两个服务可以同时启动

### 2. 配置管理混乱

**问题**: 部分服务使用 Nacos，部分使用本地配置  
**建议**: 统一使用 Nacos 或统一使用本地配置

### 3. 服务未实现

**待实现服务**:
- User Service (50051)
- Inventory Service (50054)

---

## 📊 技术栈总结

| 技术 | 用途 |
|------|------|
| Kratos | 微服务框架 |
| gRPC | 服务间通信 |
| HTTP/REST | 对外 API |
| Protobuf | 接口定义 |
| MySQL | 关系型数据库 |
| Redis | 缓存与状态存储 |
| Nacos | 配置管理 |
| Consul | 服务发现 |
| JWT | 认证鉴权 |
| Wire | 依赖注入 |
| OpenTelemetry | 链路追踪 |

---

## 📚 相关文档

- **[README.md](README.md)** - 项目简介与快速开始
- **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - 详细项目结构
- **[MIGRATION_SUMMARY.md](MIGRATION_SUMMARY.md)** - 架构迁移总结
- **[service/userauth-service/README.md](service/userauth-service/README.md)** - UserAuth 服务文档

---

## 🎯 下一步计划

### 短期（1-2周）
- [x] 解决端口冲突问题 ✅ 已完成
- [ ] 实现 User Service 业务逻辑 🔴 高优先级
- [ ] 编译所有服务
- [ ] 端到端测试认证流程
- [ ] 补齐单元测试

### 中期（1个月）
- [ ] 实现 Goods Service 业务逻辑
- [ ] 性能优化与压测
- [ ] 完善监控告警
- [ ] 文档完善

### 长期（3-6个月）
- [ ] 微服务拆分优化
- [ ] 引入消息队列
- [ ] 实现分布式事务
- [ ] 多租户支持

---

**文档维护**: 请在重大架构变更后及时更新本文档  
**最后更新**: 2025-10-23
