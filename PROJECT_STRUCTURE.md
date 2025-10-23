# 项目结构说明

## 目录结构

```
lushop-kratos-main/
├── lushop/                          # API 网关服务
│   ├── api/                         # Proto 定义
│   │   ├── lushop/v1/              # 网关 API
│   │   ├── service/user/v1/        # 用户服务 API
│   │   └── userauth/v1/            # 认证服务 API（客户端）
│   ├── cmd/lushop/                 # 主程序入口
│   ├── configs/                    # 配置文件
│   ├── internal/
│   │   ├── biz/                    # 业务逻辑层
│   │   ├── data/                   # 数据访问层
│   │   │   ├── user.go            # Redis 版本（已废弃）
│   │   │   └── user_auth_repo_grpc.go  # gRPC 版本（当前使用）
│   │   ├── service/                # HTTP 服务层
│   │   └── server/                 # 服务器配置
│   └── third_party/                # Proto 依赖
│
├── service/                         # 微服务目录
│   └── userauth-service/           # 用户认证服务（gRPC）
│       ├── api/userauth/v1/        # Proto 定义与生成代码
│       ├── cmd/userauth/           # 主程序入口
│       ├── configs/                # 配置文件
│       ├── internal/
│       │   ├── biz/                # 业务逻辑
│       │   │   └── auth.go        # 认证业务
│       │   ├── data/               # 数据访问
│       │   │   └── auth.go        # Redis 操作
│       │   ├── service/            # gRPC 服务实现
│       │   │   └── userauth.go
│       │   ├── server/             # 服务器配置
│       │   │   ├── grpc.go
│       │   │   └── registry.go    # Consul 注册
│       │   ├── conf/               # 配置结构
│       │   └── pkg/                # 辅助包
│       │       ├── captcha/        # 验证码生成
│       │       ├── sms/            # 短信发送
│       │       └── auth/           # JWT 工具
│       ├── third_party/            # Proto 依赖
│       ├── Makefile
│       ├── start.sh                # 启动脚本
│       └── README.md
│
├── MIGRATION_GUIDE.md              # 完整迁移指南
├── MIGRATION_SUMMARY.md            # 迁移总结
├── USER_AUTH_SERVICE_IMPLEMENTATION.md  # 服务端实现指南
└── PROJECT_STRUCTURE.md            # 本文件
```

## 服务说明

### 1. lushop（API 网关）
- **端口**：HTTP 8001, gRPC 9001
- **职责**：
  - 对外提供 HTTP API
  - 路由转发到后端微服务
  - 通过 gRPC 调用 UserAuth 服务处理认证
  - 通过 gRPC 调用 User 服务处理用户资料
- **服务发现**：Consul
- **配置**：`lushop/configs/config.yaml`

### 2. userauth-service（认证服务）
- **端口**：gRPC 50055
- **职责**：
  - 图形验证码生成与校验
  - 短信验证码发送与校验
  - JWT Token 签发、刷新、撤销
  - 用户黑名单管理
- **服务发现**：Consul（服务名：`user-auth-service`）
- **配置**：`service/userauth-service/configs/config.yaml`
- **存储**：Redis

## 服务调用链

```
客户端
  ↓ HTTP
API 网关 (lushop)
  ↓ gRPC
  ├─→ UserAuth Service (认证域)
  │   └─→ Redis (Token/验证码/黑名单)
  │
  └─→ User Service (用户资料域)
      └─→ MySQL (用户数据)
```

## 快速启动

### 1. 启动基础设施

```bash
# Redis
sudo systemctl start redis

# Consul
consul agent -dev
```

### 2. 启动 UserAuth 服务

```bash
cd service/userauth-service
./start.sh
# 或
./bin/userauth -conf ./configs/config.yaml
```

### 3. 启动 API 网关

```bash
cd lushop
go run ./cmd/lushop -conf ./configs/config.yaml
```

## 配置说明

### API 网关配置（lushop/configs/config.yaml）

```yaml
service:
  user:
    endpoint: "discovery:///lushop.user.service"
  user_auth:
    endpoint: "discovery:///user-auth-service"  # 指向 UserAuth 服务
```

### UserAuth 服务配置（service/userauth-service/configs/config.yaml）

```yaml
server:
  grpc:
    addr: 0.0.0.0:50055  # gRPC 监听端口

data:
  redis:
    addr: 127.0.0.1:6379

auth:
  jwt_key: "lushop-secret-key-2024"  # 与网关保持一致

registry:
  consul:
    address: 127.0.0.1:8500
```

## 开发工作流

### 修改 Proto 定义

1. 修改 `lushop/api/userauth/v1/user_auth.proto`
2. 在网关生成客户端代码：
   ```bash
   cd lushop && make api
   ```
3. 复制到服务端并生成服务端代码：
   ```bash
   cp -r lushop/api/userauth/v1/user_auth.proto service/userauth-service/api/userauth/v1/
   cd service/userauth-service && make api
   ```

### 添加新的认证功能

1. 在 `service/userauth-service/internal/biz/auth.go` 添加业务逻辑
2. 在 `service/userauth-service/internal/data/auth.go` 添加数据访问
3. 在 `service/userauth-service/internal/service/userauth.go` 实现 gRPC 接口
4. 在网关 `lushop/internal/biz/user_auth_adapter.go` 添加客户端调用

## 部署架构

```
                    ┌─────────────┐
                    │   Consul    │
                    │  (服务发现)  │
                    └─────────────┘
                           ↑
                           │ 注册/发现
              ┌────────────┴────────────┐
              ↓                         ↓
    ┌──────────────────┐      ┌──────────────────┐
    │   API Gateway    │      │  UserAuth Service│
    │   (lushop)       │◄────►│  (认证服务)       │
    │   :8001/:9001    │ gRPC │   :50055         │
    └──────────────────┘      └──────────────────┘
              ↑                         ↓
              │ HTTP                    │
              │                    ┌─────────┐
          客户端                    │  Redis  │
                                   └─────────┘
```

## 监控与日志

### 日志位置
- API 网关：标准输出
- UserAuth 服务：标准输出

### 健康检查
```bash
# 检查 Consul 服务注册
curl http://127.0.0.1:8500/v1/catalog/service/user-auth-service

# 检查 API 网关
curl http://127.0.0.1:8001/api/user/captcha
```

## 故障排查

### UserAuth 服务无法启动
1. 检查 Redis：`redis-cli ping`
2. 检查 Consul：`curl http://127.0.0.1:8500/v1/status/leader`
3. 检查端口：`lsof -i:50055`

### 网关无法调用 UserAuth
1. 检查 Consul 服务发现：`curl http://127.0.0.1:8500/v1/catalog/service/user-auth-service`
2. 检查网关配置：`service.user_auth.endpoint`
3. 查看网关日志中的 gRPC 连接错误

## 相关文档

- **完整迁移指南**：`MIGRATION_GUIDE.md`
- **迁移总结**：`MIGRATION_SUMMARY.md`
- **服务端实现**：`USER_AUTH_SERVICE_IMPLEMENTATION.md`
- **UserAuth README**：`service/userauth-service/README.md`
