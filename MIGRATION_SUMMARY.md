# 用户服务统一治理（gRPC）迁移总结

## ✅ 已完成工作

### 阶段 1：Proto 定义（已完成）
- **文件**：`lushop/api/userauth/v1/user_auth.proto`
- **内容**：
  - GetCaptcha/VerifyCaptcha：图形验证码
  - SendSms/VerifySms：短信验证码
  - IssueToken/RefreshToken/RevokeToken：Token 管理
  - AddToBlacklist/CheckBlacklist：黑名单管理
- **生成代码**：已通过 `make api` 生成客户端代码

### 阶段 2：网关侧客户端注入（已完成）
- **修改文件**：
  - `internal/conf/conf.proto`：增加 `Service.UserAuth.Endpoint` 配置
  - `internal/data/data.go`：注入 `UserAuthClient`，扩展 `Data` 结构体
  - `configs/nacosRemote.yaml`：增加 `service.user_auth.endpoint` 配置
- **状态**：Wire 已重新生成，依赖已清理

### 阶段 3：gRPC 版仓库实现（已完成）
- **新建文件**：`internal/data/user_auth_repo_grpc.go`
- **实现**：
  - 实现 `biz.UserRepo` 接口
  - 认证域方法通过 `UserAuth` gRPC 调用
  - 用户资料方法继续通过 `User` 服务
  - 旧的 Redis 方法标记为废弃（保留接口兼容）

### 阶段 4：业务层改造（已完成）
- **新建文件**：`internal/biz/user_auth_adapter.go`
  - 封装对 UserAuth 服务的调用
  - 提供 GetCaptcha、VerifyCaptcha、SendSms、VerifySms、IssueToken、RefreshToken、RevokeToken、CheckBlacklist 方法

- **修改文件**：`internal/biz/user.go`
  - ✅ `GetCaptcha()`：调用 UserAuth.GetCaptcha
  - ✅ `PasswordLogin()`：验证码校验 + Token 签发改为 gRPC
  - ✅ `CreateUser()`：验证码校验 + Token 签发改为 gRPC
  - ✅ `SendSms()`：调用 UserAuth.SendSms
  - ✅ `VerifySms()`：调用 UserAuth.VerifySms + IssueToken
  - ✅ `RefreshToken()`：调用 UserAuth.RefreshToken
  - ✅ `Logout()`：调用 UserAuth.RevokeToken
  - ✅ `UserDetailByID/UpdateUser/UpdatePassword/DeleteUser/ListUsers`：黑名单检查改为 authAdapter.CheckBlacklist

- **清理**：移除未使用的导入（captcha、sms、os、阿里云 SDK 等）

### 阶段 5：Wire 装配切换（已完成）
- **修改文件**：
  - `internal/biz/biz.go`：ProviderSet 增加 `NewUserAuthAdapter`
  - `internal/data/data.go`：ProviderSet 切换到 `NewUserAuthRepoGRPC`（注释掉 `NewUserRepo`）
- **状态**：Wire 已重新生成 `cmd/lushop/wire_gen.go`

## 📋 待完成工作

### 阶段 6：服务侧实现（待实现）
- **任务**：创建独立的 UserAuth 服务项目
- **参考文档**：`USER_AUTH_SERVICE_IMPLEMENTATION.md`
- **核心文件**：
  - `internal/biz/auth.go`：认证业务逻辑
  - `internal/data/auth.go`：Redis 数据访问
  - `internal/service/userauth.go`：gRPC 服务实现
  - `internal/server/grpc.go`：gRPC 服务器配置
  - `cmd/userauth/main.go`：服务入口

### 阶段 7：联调测试（待执行）
- [ ] 启动 UserAuth 服务（端口 50055）
- [ ] 启动 API 网关（端口 8001/9001）
- [ ] 验证 Consul 服务发现
- [ ] 端到端测试：
  - 用户注册流程
  - 密码登录流程
  - 短信登录流程
  - Token 刷新流程
  - 登出流程
  - 黑名单校验

### 阶段 8：观测与治理（待补齐）
- [ ] 补齐审计日志（登录/注册/短信/Token 操作）
- [ ] 限流策略（短信 60s 冷却、登录失败限制、验证码 IP 限流）
- [ ] Tracing 链路追踪
- [ ] Prometheus 指标监控
- [ ] 告警配置

### 阶段 9：清理（可选）
- [ ] 移除网关侧冗余的 Redis 认证逻辑（`internal/data/user.go` 中的 Token/验证码/短信方法）
- [ ] 移除未使用的依赖（阿里云 SMS SDK 可保留在 UserAuth 服务中）
- [ ] 更新文档与注释

## 🎯 架构变化对比

### 迁移前（网关自管 Redis）
```
客户端 → API 网关
           ├─ JWT 验证（中间件）
           ├─ 本地生成验证码 → Redis
           ├─ 本地发送短信 → Redis
           ├─ 本地签发 Token → Redis
           └─ 本地管理黑名单 → Redis
```

### 迁移后（用户服务统一治理）
```
客户端 → API 网关
           ├─ JWT 验证（中间件，轻量）
           └─ 转调 UserAuth 服务（gRPC）
                 ├─ 生成验证码 → Redis
                 ├─ 发送短信 → Redis
                 ├─ 签发 Token → Redis
                 ├─ 刷新 Token → Redis
                 ├─ 撤销 Token → Redis
                 └─ 管理黑名单 → Redis
```

## 📊 收益分析

### 优势
1. **领域边界清晰**：认证逻辑集中在 UserAuth 服务，网关保持轻量
2. **易于演进**：
   - 接入 SSO/OIDC/SAML
   - 多端策略（Web/App/小程序不同 Token 策略）
   - 设备指纹、IP 风控、异地登录检测
3. **审计与合规**：集中埋点，便于审计与合规（GDPR、等保）
4. **可扩展性**：UserAuth 服务可独立扩容，不影响网关
5. **故障隔离**：认证服务故障不影响其他业务服务

### 注意事项
1. **多一跳延迟**：网关 → UserAuth 服务（可通过连接池、缓存优化）
2. **复杂度增加**：需要维护独立服务、服务发现、可用性建设
3. **Redis 依赖**：UserAuth 服务依赖 Redis，需要高可用配置

## 🔄 回滚方案

如需回滚到 Redis 方案：

### 1. 修改 data ProviderSet
```go
// internal/data/data.go
var ProviderSet = wire.NewSet(
    NewData,
    NewUserRepo,         // 恢复 Redis 版本
    // NewUserAuthRepoGRPC,  // 注释掉 gRPC 版本
    NewUserServiceClient,
    NewUserAuthClient,
    NewRegister,
    NewDiscovery,
    NewRedis,
)
```

### 2. 重新生成 Wire
```bash
cd cmd/lushop && wire
```

### 3. 重启服务
```bash
go build ./cmd/lushop
./lushop -conf ./configs/config.yaml
```

## 📚 相关文档

- **完整迁移指南**：`MIGRATION_GUIDE.md`
- **服务端实现指南**：`USER_AUTH_SERVICE_IMPLEMENTATION.md`
- **阶段 2 命令**：`MIGRATION_STEP2_COMMANDS.md`

## 🚀 下一步行动

1. **立即执行**：实现 UserAuth 服务端（参考 `USER_AUTH_SERVICE_IMPLEMENTATION.md`）
2. **联调测试**：启动服务并验证所有流程
3. **性能测试**：压测关键接口（登录/注册/刷新）
4. **上线准备**：补齐监控、告警、文档

## 📞 技术支持

如遇问题，请检查：
1. Consul 服务发现是否正常
2. gRPC 连接是否建立（检查日志）
3. Redis 连接是否正常
4. 配置文件中的 `service.user_auth.endpoint` 是否正确
5. UserAuth 服务是否启动并注册到 Consul

---

**迁移完成度**：网关侧 100%，服务侧 0%（待实现）
**预计剩余工作量**：2-3 天（服务端实现 + 联调测试）
