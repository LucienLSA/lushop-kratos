# Sentinel 限流调用流程详解

本文档说明 Sentinel 限流在代码中的每一层关系和完整调用过程。

---

## 📋 目录结构

```
lushop/
├── cmd/lushop/main.go                    # 应用入口，初始化 Sentinel
├── internal/
│   ├── conf/conf.proto                   # 配置定义（Protobuf）
│   ├── server/http.go                    # HTTP 服务器，注册中间件
│   └── pkg/
│       ├── sentinel/sentinel.go          # Sentinel 核心工具包（初始化、规则加载）
│       └── middleware/sentinel/sentinel.go  # Sentinel HTTP 中间件（请求拦截）
└── configs/
    └── nacos-config.yaml                 # 配置文件（限流规则）
```

---

## 🔄 完整调用流程图

```
┌─────────────────────────────────────────────────────────┐
│              应用启动阶段                                  │
└─────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│  Step 1: main() 启动                                      │
│  - 加载配置文件 (nacos-config.yaml)                      │
│  - 解析配置到 conf.Server                                │
└─────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│  Step 2: initSentinel()                                  │
│  位置: cmd/lushop/main.go:224                            │
│  职责: 初始化 Sentinel 并加载规则                         │
│    ├─ sentinel.Init() - 初始化 Sentinel 核心             │
│    ├─ LoadFlowRules() - 加载限流规则                     │
│    └─ LoadCircuitBreakerRules() - 加载熔断规则（可选）   │
└─────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│  Step 3: NewHTTPServer()                                 │
│  位置: internal/server/http.go:35                        │
│  职责: 创建 HTTP 服务器并注册中间件                       │
│    └─ initSentinelMiddleware() - 创建 Sentinel 中间件    │
└─────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│              运行时阶段（HTTP 请求处理）                   │
└─────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│  Step 4: HTTP 请求到达                                   │
│  - 客户端发送请求                                         │
│  - Kratos HTTP Server 接收请求                           │
└─────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│  Step 5: 中间件链执行                                     │
│  1. recovery.Recovery()                                  │
│  2. Sentinel 中间件 ← 限流检查                           │
│  3. validate.Validator()                                 │
│  4. ... 其他中间件                                        │
└─────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│  Step 6: Sentinel 中间件执行                             │
│  位置: pkg/middleware/sentinel/sentinel.go:62           │
│    ├─ 检查是否启用                                        │
│    ├─ 检查白名单                                          │
│    ├─ api.Entry(resource) - 限流检查                     │
│    ├─ 执行业务逻辑 handler(ctx, req)                     │
│    ├─ entry.SetError(err) - 记录错误（熔断统计）         │
│    └─ entry.Exit() - 释放资源，更新统计                    │
└─────────────────────────────────────────────────────────┘
                    │
            ┌───────┴───────┐
            │               │
            ▼               ▼
┌──────────────────┐  ┌──────────────────┐
│ 限流触发          │  │ 限流通过          │
│ 返回 429 错误     │  │ 返回业务响应      │
└──────────────────┘  └──────────────────┘
```

---

## 📚 各层详细说明

### 第一层：应用入口层 (`cmd/lushop/main.go`)

**职责**：应用启动、配置加载、Sentinel 初始化

**关键函数**：`initSentinel()`

**执行流程**：
1. 调用 `sentinel.Init()` 初始化 Sentinel 核心
2. 从配置文件读取限流规则（`FlowRules`）
3. 如果有配置，调用 `sentinel.LoadFlowRules()` 加载规则
4. 如果没有配置，调用 `sentinel.LoadDefaultRules()` 加载默认规则
5. 可选：加载熔断规则和系统规则

**调用时机**：在 `main()` 函数中，应用启动时调用（加载配置之后，创建 HTTP Server 之前）

---

### 第二层：Sentinel 核心工具包 (`internal/pkg/sentinel/sentinel.go`)

**职责**：Sentinel 初始化、规则管理、配置转换

**关键函数**：

- **`Init()`**：调用 `api.InitDefault()` 初始化 Sentinel 核心组件
- **`LoadFlowRules()`**：将配置结构体转换为 Sentinel SDK 的 `flow.Rule`，并加载到 Sentinel 核心
- **`LoadDefaultRules()`**：当配置文件中没有指定限流规则时，使用代码中定义的默认规则

---

### 第三层：HTTP 服务器层 (`internal/server/http.go`)

**职责**：创建 HTTP 服务器、注册中间件链

**关键函数**：

- **`NewHTTPServer()`**：创建 HTTP 服务器，在中间件链中注册 Sentinel 中间件（第二位，在 `recovery` 之后）
- **`initSentinelMiddleware()`**：从配置中读取 `enabled` 和 `whitelist`，创建 Sentinel 中间件实例

---

### 第四层：Sentinel 中间件层 (`internal/pkg/middleware/sentinel/sentinel.go`)

**职责**：请求拦截、限流检查、错误处理

**执行流程**：

1. **检查是否启用**：如果 `enabled = false`，直接放行
2. **获取请求信息**：从 `context` 中提取 `endpoint`（接口路径）
3. **检查白名单**：如果接口在白名单中，直接放行
4. **Sentinel 限流检查**：调用 `api.Entry(resource)` 检查是否触发限流
5. **限流触发处理**：如果触发限流，返回 `ErrRateLimitExceeded`（HTTP 429）
6. **执行业务逻辑**：如果通过限流检查，调用 `handler(ctx, req)` 执行实际业务逻辑
7. **记录错误**：如果业务逻辑返回错误，调用 `entry.SetError(err)` 记录（用于熔断统计）
8. **退出资源**：调用 `entry.Exit()` 释放 Sentinel 资源，更新统计信息

---

## 🔍 关键调用点

### 1. `api.Entry(resource)` - Sentinel 限流检查

**位置**：`pkg/middleware/sentinel/sentinel.go:85`

**作用**：
- Sentinel Go SDK 的核心 API
- 根据资源名（接口路径）检查是否触发限流
- 内部检查：限流规则（FlowRule）、熔断规则（CircuitBreakerRule）、系统规则（SystemRule）

**返回值**：
- `entry`：Sentinel 资源入口，用于后续的 `SetError()` 和 `Exit()`
- `blockErr`：如果触发限流，返回错误；否则为 `nil`

### 2. `entry.SetError(err)` - 记录错误

**位置**：`pkg/middleware/sentinel/sentinel.go:102`

**作用**：记录业务逻辑返回的错误，用于熔断规则的错误率统计

### 3. `entry.Exit()` - 退出资源

**位置**：`pkg/middleware/sentinel/sentinel.go:106`

**作用**：释放 Sentinel 资源，更新统计信息（QPS、响应时间、并发数）。**必须调用**，否则会导致资源泄漏。

---

## 📊 数据流向图

```
配置文件 (nacos-config.yaml)
    │ 解析
    ▼
conf.Server (Protobuf 结构体)
    │ 读取
    ▼
initSentinel() (main.go)
    │ 转换
    ▼
sentinel.FlowRuleConfig (Go 结构体)
    │ 加载
    ▼
flow.Rule (Sentinel SDK)
    │ 存储
    ▼
Sentinel 核心 (内存中的规则表)
    │ 查询
    ▼
api.Entry(resource) (运行时检查)
    │ 判断
    ▼
限流结果 (通过/拒绝)
```

---

## 🎯 关键设计点

### 1. 分层设计

- **配置层**：`conf.proto` + `nacos-config.yaml`
- **初始化层**：`main.go` + `pkg/sentinel/sentinel.go`
- **中间件层**：`server/http.go` + `pkg/middleware/sentinel/sentinel.go`
- **运行时层**：Sentinel Go SDK

### 2. 职责分离

- **`pkg/sentinel`**：负责配置管理和规则加载（启动时）
- **`pkg/middleware/sentinel`**：负责请求拦截和限流检查（运行时）

### 3. 配置优先级

1. **配置文件规则**（`nacos-config.yaml` 中的 `flow_rules`）
2. **默认规则**（`pkg/sentinel/sentinel.go` 中的 `LoadDefaultRules()`）
3. **无规则**（如果都没有，则不限流）

### 4. 白名单机制

- 白名单接口直接放行，不经过 Sentinel 检查
- 适用于公开接口（登录、注册、商品列表等）

---

## 🚀 总结

Sentinel 限流的完整调用链：

1. **启动阶段**：`main.go` → `initSentinel()` → `pkg/sentinel` → 加载规则到 Sentinel 核心
2. **注册阶段**：`server/http.go` → `initSentinelMiddleware()` → `pkg/middleware/sentinel` → 注册中间件
3. **运行时**：HTTP 请求 → 中间件链 → Sentinel 中间件 → `api.Entry()` → 限流检查 → 业务逻辑

每一层都有明确的职责，通过清晰的接口和配置实现了解耦和可维护性。
