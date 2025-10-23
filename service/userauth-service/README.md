# UserAuth Service

用户认证服务 - 统一管理验证码、短信、Token 签发与黑名单

## 功能

- **图形验证码**：生成与校验
- **短信验证码**：发送与校验（集成阿里云短信）
- **Token 管理**：签发、刷新、撤销
- **黑名单管理**：登出黑名单、风控黑名单

## 架构

```
UserAuth Service (gRPC)
├── API Layer (gRPC Server)
├── Business Layer (biz)
│   └── AuthUsecase
├── Data Layer (data)
│   └── AuthRepo (Redis)
└── Infrastructure
    ├── Captcha (base64Captcha)
    ├── SMS (阿里云)
    └── JWT (golang-jwt)
```

## 快速开始

### 1. 配置

编辑 `configs/config.yaml`：

```yaml
server:
  grpc:
    addr: 0.0.0.0:50055  # gRPC 端口

data:
  redis:
    addr: 127.0.0.1:6379  # Redis 地址

auth:
  jwt_key: "your-secret-key"  # JWT 密钥（与网关保持一致）

sms:
  api_key: "your-aliyun-key"      # 阿里云 AccessKey
  api_secret: "your-aliyun-secret"
  sign_name: "签名"
  template_code: "SMS_123456"

registry:
  consul:
    address: 127.0.0.1:8500  # Consul 地址
```

### 2. 编译

```bash
make build
```

### 3. 运行

```bash
./bin/userauth -conf ./configs/config.yaml
```

或使用环境变量覆盖短信配置：

```bash
export SMS_API_KEY="your-key"
export SMS_API_SECRET="your-secret"
./bin/userauth -conf ./configs/config.yaml
```

### 4. 验证

服务启动后：
- gRPC 监听：`0.0.0.0:50055`
- Consul 注册：服务名 `user-auth-service`

检查 Consul：
```bash
curl http://127.0.0.1:8500/v1/catalog/service/user-auth-service
```

## 开发

### 生成代码

```bash
# 生成 proto 代码
make api

# 生成配置代码
make config

# 生成 Wire
make wire

# 一键生成所有
make all
```

### 目录结构

```
userauth-service/
├── api/userauth/v1/       # Proto 定义与生成代码
├── cmd/userauth/          # 主程序入口
├── configs/               # 配置文件
├── internal/
│   ├── biz/               # 业务逻辑
│   │   └── auth.go        # 认证业务
│   ├── data/              # 数据访问
│   │   └── auth.go        # Redis 操作
│   ├── service/           # gRPC 服务实现
│   │   └── userauth.go
│   ├── server/            # 服务器配置
│   │   ├── grpc.go
│   │   └── registry.go    # Consul 注册
│   ├── conf/              # 配置结构
│   └── pkg/               # 辅助包
│       ├── captcha/       # 验证码生成
│       ├── sms/           # 短信发送
│       └── auth/          # JWT 工具
└── third_party/           # Proto 依赖
```

## API 接口

### gRPC 接口

```protobuf
service UserAuth {
  // 图形验证码
  rpc GetCaptcha(google.protobuf.Empty) returns (CaptchaReply);
  rpc VerifyCaptcha(VerifyCaptchaReq) returns (VerifyCaptchaReply);
  
  // 短信验证码
  rpc SendSms(SendSmsReq) returns (SendSmsReply);
  rpc VerifySms(VerifySmsReq) returns (VerifySmsReply);
  
  // Token 管理
  rpc IssueToken(IssueTokenReq) returns (TokenReply);
  rpc RefreshToken(RefreshTokenReq) returns (TokenReply);
  rpc RevokeToken(RevokeTokenReq) returns (google.protobuf.Empty);
  
  // 黑名单
  rpc AddToBlacklist(AddToBlacklistReq) returns (google.protobuf.Empty);
  rpc CheckBlacklist(CheckBlacklistReq) returns (CheckBlacklistReply);
}
```

## Redis 键设计

```
captcha:{captcha_id}           # 验证码答案，TTL 5分钟
sms_code:{mobile}              # 短信验证码，TTL 5分钟
sms_cooldown:{mobile}          # 短信冷却，TTL 60秒
user_access_token:{user_id}    # Access Token，TTL 30分钟
user_refresh_token:{user_id}   # Refresh Token，TTL 7天
user_blacklist:{user_id}       # 黑名单，TTL 30分钟
```

## 监控与日志

### 日志

服务使用结构化日志，包含：
- `ts`：时间戳
- `caller`：调用位置
- `service.id`：服务实例 ID
- `service.name`：服务名
- `trace.id`：链路追踪 ID

### 指标

- gRPC 请求计数
- gRPC 请求延迟
- Redis 连接状态
- 短信发送成功率

## 故障排查

### 服务无法启动

1. 检查 Redis 连接：`redis-cli ping`
2. 检查 Consul 连接：`curl http://127.0.0.1:8500/v1/status/leader`
3. 检查端口占用：`lsof -i:50055`

### 短信发送失败

1. 检查阿里云 API 凭证
2. 检查短信模板是否审核通过
3. 检查手机号格式（仅支持中国大陆）
4. 查看日志中的详细错误信息

### Token 签发失败

1. 检查 JWT 密钥配置
2. 检查 Redis 连接
3. 确认用户 ID 有效

## 性能优化

- Redis 连接池：默认配置已优化
- gRPC 连接复用：客户端自动处理
- 短信冷却：60秒防刷
- Token 缓存：Redis 存储，避免重复签发

## 安全建议

1. **JWT 密钥**：使用强随机密钥，定期轮换
2. **短信限流**：已实现 60秒冷却，可根据需要调整
3. **Redis 密码**：生产环境务必设置 Redis 密码
4. **Consul ACL**：生产环境启用 Consul ACL
5. **TLS**：生产环境启用 gRPC TLS

## 部署

### Docker

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/userauth .
COPY --from=builder /app/configs ./configs
CMD ["./userauth", "-conf", "./configs/config.yaml"]
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: userauth-service
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: userauth
        image: userauth-service:latest
        ports:
        - containerPort: 50055
        env:
        - name: SMS_API_KEY
          valueFrom:
            secretKeyRef:
              name: sms-credentials
              key: api-key
```

## 许可证

MIT
