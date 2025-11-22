# 🛒 Lushop 微服务电商平台

> 基于 Go-Kratos 框架的生产级微服务电商平台，采用 DDD 领域驱动设计，实现高可用、高并发的分布式系统架构。

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Kratos](https://img.shields.io/badge/Kratos-v2.8.3-blue)](https://go-kratos.dev)
[![gRPC](https://img.shields.io/badge/gRPC-1.69-green)](https://grpc.io)
[![License](https://img.shields.io/badge/License-MIT-yellow)](LICENSE)

## 📖 项目简介

这是一个**生产级微服务电商平台**，采用主流的微服务架构和技术栈，涵盖了用户、商品、订单、库存等核心业务模块。项目严格遵循 DDD 领域驱动设计，实现了服务注册发现、配置中心、链路追踪、分布式事务等企业级特性。

**适用场景**：
- 🎯 学习微服务架构最佳实践
- 🎯 了解 Go 语言在大型项目中的应用
- 🎯 掌握分布式系统设计思想
- 🎯 面试准备和技术积累

**原项目**：https://github.com/LucienLSA/lushop.git

项目架构总览

```
┌─────────────────────────────────────────────────────────────────┐
│                    Lushop API Gateway                            │
│              HTTP: 8001 | gRPC: 9001                            │
│    ┌──────────────────────────────────────────────────┐         │
│    │  User | UserAuth | Cart | Goods | Order |        │         │
│    │  Inventory | UserOp  (7 Services)                │         │
│    └──────────────────────────────────────────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│ User Service  │    │ Goods Service │    │ Order Service │
│  Port: 50051  │    │  Port: 50052  │    │  Port: 50053  │
│  ✅ 100%      │    │  ✅ 100%      │    │  ✅ 100%      │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│UserAuth Svc   │    │Inventory Svc  │    │ UserOp Svc    │
│  Port: 50056  │    │  Port: 50054  │    │  Port: 50055  │
│  ✅ 100%      │    │  ✅ 100%      │    │  ✅ 100%      │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        └─────────────────────┴─────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
                    ▼                   ▼
            ┌──────────────┐    ┌──────────────┐
            │    Consul    │    │    Redis     │
            │ Port: 8500   │    │ Port: 6379   │
            └──────────────┘    └──────────────┘
                    │                   │
                    ▼                   ▼
            ┌──────────────┐    ┌──────────────┐
            │    MySQL     │    │   Jaeger     │
            │ Port: 3306   │    │ Port: 16686  │
            └──────────────┘    └──────────────┘
                    ▼
            ┌──────────────┐
            │  RocketMQ    │
            │ Port: 9876   │
            └──────────────┘
```---

## 🎯 核心技术栈

### 后端框架
- **Go 1.23+** - 高性能编程语言
- **Kratos v2.8.3** - B站开源的微服务框架
- **gRPC 1.69** - 高性能 RPC 框架
- **Protocol Buffers** - 接口定义语言

### 微服务治理
- **Consul** - 服务注册与发现、健康检查
- **Nacos** - 配置中心、动态配置管理
- **Jaeger** - 分布式链路追踪
- **Prometheus** - 服务监控指标

### 数据存储
- **MySQL 8.0** - 关系型数据库
- **Redis 7.0** - 缓存、分布式锁、Session
- **ElasticSearch** - 商品搜索引擎

### 消息队列
- **RocketMQ** - 分布式事务消息、延迟消息
- **Asynq** - 异步任务队列

### 开发工具
- **Wire** - 依赖注入代码生成
- **Protoc** - Proto 文件编译
- **Makefile** - 项目构建自动化
- **Docker Compose** - 本地开发环境

### 安全认证
- **JWT** - 无状态认证
- **BCrypt** - 密码加密
- **图形验证码** - 防机器人
- **短信验证码** - 手机验证

---

## 🏗️ 系统架构设计

### 1. 微服务架构

采用**领域驱动设计（DDD）**，将系统拆分为 6 个独立微服务 + 1 个 API 网关：

| 服务名 | 端口 | 职责 | 技术特点 |
|--------|------|------|----------|
| **User** | 50051 | 用户管理 | JWT认证、密码加密 |
| **Goods** | 50052 | 商品管理 | ES搜索、分类管理 |
| **Order** | 50053 | 订单管理 | 雪花算法、事务消息 |
| **Inventory** | 50054 | 库存管理 | 分布式锁、防超卖 |
| **UserOp** | 50055 | 用户操作 | 地址、收藏、留言 |
| **UserAuth** | 50056 | 认证服务 | Token管理、验证码 |
| **Gateway** | 8001/9001 | API网关 | HTTP→gRPC转换 |

### 2. 四层架构（Clean Architecture）

```
┌─────────────────────────────────────────────┐
│  Transport Layer (HTTP/gRPC Server)         │  ← 协议层
├─────────────────────────────────────────────┤
│  Service Layer (业务编排)                    │  ← 服务层
├─────────────────────────────────────────────┤
│  Biz Layer (业务逻辑)                        │  ← 领域层
├─────────────────────────────────────────────┤
│  Data Layer (数据访问)                       │  ← 数据层
└─────────────────────────────────────────────┘
```

**优势**：
- ✅ 职责清晰，易于维护
- ✅ 依赖倒置，便于测试
- ✅ 业务逻辑与技术实现解耦
- ✅ 支持多种传输协议（HTTP/gRPC）

### 3. 服务间通信

```
┌──────────┐                    ┌──────────┐
│  Order   │  ──── gRPC ────>   │  Goods   │
│ Service  │                    │ Service  │
└──────────┘                    └──────────┘
     │                               │
     │                               │
     v                               v
┌──────────┐                    ┌──────────┐
│Inventory │                    │  MySQL   │
│ Service  │                    │  Redis   │
└──────────┘                    └──────────┘
```

**特点**：
- 🔄 gRPC 高性能通信
- 🔍 Consul 服务发现
- ⚖️ 客户端负载均衡
- 🛡️ 超时控制与重试

---

## 💡 核心业务功能

### 用户体系
- ✅ 用户注册/登录（手机号+验证码）
- ✅ JWT Token 认证与刷新
- ✅ 用户信息管理（CRUD）
- ✅ 密码加密存储（BCrypt）
- ✅ 图形验证码防刷
- ✅ 短信验证码（阿里云）
- ✅ 用户黑名单机制
- ✅ 管理员权限控制

### 商品体系
- ✅ 商品列表/详情查询
- ✅ 商品分类管理
- ✅ 商品搜索（ElasticSearch）
- ✅ 商品管理（CRUD）
- ✅ 批量操作支持

### 订单体系
- ✅ 购物车管理（增删改查）
- ✅ 订单创建（分布式事务）
- ✅ 订单列表/详情
- ✅ 订单状态流转
- ✅ 订单号生成（雪花算法）
- ✅ 订单超时自动取消（延迟消息）
- ✅ 库存扣减与回滚

### 库存体系
- ✅ 库存设置/查询
- ✅ 库存扣减（分布式锁）
- ✅ 库存归还（事务回滚）
- ✅ 防超卖机制

### 用户操作
- ✅ 收货地址管理
- ✅ 商品收藏
- ✅ 用户留言

---

## 🌟 项目亮点（面试重点）

### 1. 分布式事务解决方案 ⭐⭐⭐⭐⭐

**场景**：用户下单时需要同时完成订单创建和库存扣减，如何保证数据一致性？

**解决方案**：RocketMQ 事务消息 + 本地事务表

```go
// 1. 发送半消息（Half Message）
msg := &primitive.Message{
    Topic: "order_inventory_topic",
    Body:  orderData,
}
result, _ := producer.SendMessageInTransaction(ctx, msg)

// 2. 执行本地事务（创建订单）
func ExecuteLocalTransaction(msg *primitive.Message) LocalTransactionState {
    // 创建订单到数据库
    err := createOrder(orderData)
    if err != nil {
        return primitive.RollbackMessageState  // 回滚
    }
    return primitive.CommitMessageState  // 提交
}

// 3. 事务回查（防止网络异常）
func CheckLocalTransaction(msg *primitive.MessageExt) LocalTransactionState {
    // 查询订单是否创建成功
    exists := checkOrderExists(orderId)
    if exists {
        return primitive.CommitMessageState
    }
    return primitive.RollbackMessageState
}

// 4. 消费消息（扣减库存）
func ConsumeMessage(msgs []*primitive.MessageExt) {
    for _, msg := range msgs {
        // 扣减库存
        deductInventory(orderData)
    }
}
```

**优势**：
- ✅ 最终一致性保证
- ✅ 高可用（支持事务回查）
- ✅ 解耦服务（异步处理）

### 2. 分布式锁防超卖 ⭐⭐⭐⭐⭐

**场景**：高并发下如何防止库存超卖？

**解决方案**：Redis 分布式锁 + 乐观锁

```go
// 方案1: Redis 分布式锁
func DeductInventory(goodsId int32, nums int32) error {
    lockKey := fmt.Sprintf("lock:inventory:%d", goodsId)
    
    // 获取分布式锁
    lock := redis.NewLock(lockKey, 10*time.Second)
    if !lock.TryLock() {
        return errors.New("系统繁忙，请稍后重试")
    }
    defer lock.Unlock()
    
    // 查询库存
    inventory := getInventory(goodsId)
    if inventory.Stocks < nums {
        return errors.New("库存不足")
    }
    
    // 扣减库存
    inventory.Stocks -= nums
    updateInventory(inventory)
    
    return nil
}

// 方案2: 数据库乐观锁
UPDATE inventory 
SET stocks = stocks - #{nums}, version = version + 1
WHERE goods_id = #{goodsId} 
  AND stocks >= #{nums}
  AND version = #{version}
```

**对比**：
| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| Redis锁 | 性能高、跨服务 | 需要Redis | 高并发场景 |
| 乐观锁 | 无需额外组件 | 冲突重试 | 中低并发 |

### 3. 服务注册与发现 ⭐⭐⭐⭐

**场景**：微服务如何相互调用？如何实现负载均衡？

**解决方案**：Consul + gRPC 客户端负载均衡

```go
// 1. 服务注册
func RegisterService() {
    consulClient, _ := consulAPI.NewClient(consulAPI.DefaultConfig())
    
    registration := &consulAPI.AgentServiceRegistration{
        ID:      "lushop.order.service-001",
        Name:    "lushop.order.service",
        Address: "127.0.0.1",
        Port:    50053,
        Check: &consulAPI.AgentServiceCheck{
            GRPC:     "127.0.0.1:50053",
            Interval: "10s",
            Timeout:  "5s",
        },
    }
    
    consulClient.Agent().ServiceRegister(registration)
}

// 2. 服务发现
func NewGoodsClient() {
    // 创建 Consul 服务发现
    dis := consul.New(consulClient)
    
    // 创建 gRPC 连接（自动负载均衡）
    conn, _ := grpc.DialInsecure(
        context.Background(),
        grpc.WithEndpoint("discovery:///lushop.goods.service"),
        grpc.WithDiscovery(dis),
    )
    
    client := goodsv1.NewGoodsClient(conn)
}
```

**优势**：
- ✅ 自动服务发现
- ✅ 健康检查（自动剔除故障节点）
- ✅ 客户端负载均衡
- ✅ 动态扩缩容

### 4. 配置中心动态管理 ⭐⭐⭐⭐

**场景**：如何实现配置热更新？避免重启服务？

**解决方案**：Nacos 配置中心

```go
// 1. 从 Nacos 加载配置
func LoadConfig() {
    nacosClient, _ := clients.NewConfigClient(vo.NacosClientParam{
        ServerConfigs: []constant.ServerConfig{{
            IpAddr: "127.0.0.1",
            Port:   8848,
        }},
    })
    
    // 获取配置
    content, _ := nacosClient.GetConfig(vo.ConfigParam{
        DataId: "user.yaml",
        Group:  "lushop_grpc",
    })
    
    // 监听配置变化
    nacosClient.ListenConfig(vo.ConfigParam{
        DataId: "user.yaml",
        Group:  "lushop_grpc",
        OnChange: func(namespace, group, dataId, data string) {
            // 配置变更回调
            reloadConfig(data)
        },
    })
}
```

**优势**：
- ✅ 集中管理配置
- ✅ 配置热更新（无需重启）
- ✅ 多环境配置隔离
- ✅ 配置版本管理

### 5. 链路追踪与监控 ⭐⭐⭐⭐

**场景**：微服务调用链路复杂，如何快速定位问题？

**解决方案**：Jaeger 分布式链路追踪

```go
// 1. 初始化 Tracer
func InitTracer() {
    exporter, _ := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://localhost:14268/api/traces"),
    ))
    
    tp := tracesdk.NewTracerProvider(
        tracesdk.WithBatcher(exporter),
        tracesdk.WithResource(resource.NewSchemaless(
            semconv.ServiceNameKey.String("lushop.order.service"),
        )),
    )
    
    otel.SetTracerProvider(tp)
}

// 2. 自动追踪（Kratos 内置）
// 每个 gRPC 调用自动生成 Trace
```

**效果**：
- ✅ 可视化调用链路
- ✅ 性能瓶颈分析
- ✅ 错误快速定位
- ✅ 服务依赖关系图

### 6. Sentinel 限流降级 ⭐⭐⭐⭐

**场景**：如何防止系统过载？如何保护核心服务？

**解决方案**：Sentinel 限流降级

```go
// 1. 初始化 Sentinel
err := sentinel.Init(logger)

// 2. 加载限流规则
flowRules := []*sentinel.FlowRuleConfig{
    {
        Resource:         "/lushop.lushop.v1.Order/CreateOrder",
        Threshold:        1000,  // QPS 1000
        ControlBehavior:  0,     // Reject
        StatIntervalInMs: 1000,
    },
}
sentinel.LoadFlowRules(flowRules)

// 3. 中间件自动限流检查
entry, blockErr := api.Entry(resource)
if blockErr != nil {
    return ErrRateLimitExceeded  // HTTP 429
}
defer entry.Exit()
```

**功能**：
- ✅ **API 限流**：QPS 限制，防止接口过载
- ✅ **熔断降级**：慢调用、错误率、错误数三种策略
- ✅ **系统保护**：基于 CPU、响应时间、并发数等指标
- ✅ **白名单机制**：公开接口不受限流影响

**优势**：
- ✅ 防止系统过载
- ✅ 快速失败，避免雪崩
- ✅ 保护核心服务
- ✅ 灵活的规则配置

**详细文档**：参考 [Sentinel 调用流程](lushop/docs/SENTINEL_CALL_FLOW.md)

### 7. 依赖注入（Wire） ⭐⭐⭐

**场景**：如何优雅地管理服务依赖？

**解决方案**：Google Wire 自动依赖注入

```go
// wire.go
//go:build wireinject

func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,   // HTTP/gRPC Server
        data.ProviderSet,     // MySQL/Redis Client
        biz.ProviderSet,      // 业务逻辑
        service.ProviderSet,  // 服务层
        newApp,
    ))
}

// 自动生成 wire_gen.go
func wireApp(server *conf.Server, data *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
    dataData, cleanup, err := data.NewData(data, logger)
    if err != nil {
        return nil, nil, err
    }
    userRepo := data.NewUserRepo(dataData, logger)
    userUsecase := biz.NewUserUsecase(userRepo, logger)
    userService := service.NewUserService(userUsecase, logger)
    grpcServer := server.NewGRPCServer(server, userService, logger)
    app := newApp(logger, grpcServer)
    return app, func() {
        cleanup()
    }, nil
}
```

**优势**：
- ✅ 编译期检查（避免运行时错误）
- ✅ 代码自动生成
- ✅ 依赖关系清晰
- ✅ 易于测试（Mock 注入）

---

统计数据

### 服务统计
| 指标 | 数量 | 状态 |
|------|------|------|
| 微服务总数 | 7 个 | ✅ 100% |
| gRPC 微服务 | 6 个 | ✅ 100% |
| API 网关 | 1 个 | ✅ 100% |
| gRPC API 总数 | 48 个 | ✅ 100% |
| HTTP API 总数 | 38 个 | ✅ 100% |
| Wire 依赖注入 | 7/7 | ✅ 100% |
| Consul 注册 | 7/7 | ✅ 100% |

### 技术栈统计
| 技术 | 状态 | 完成度 |
|------|------|--------|
| Kratos v2 | ✅ | 100% |
| gRPC + HTTP | ✅ | 100% |
| Consul | ✅ | 100% |
| Nacos | ✅ | 100% |
| Redis | ✅ | 100% |
| MySQL | ✅ | 100% |
| JWT | ✅ | 100% |
| Jaeger | ✅ | 100% |
| Wire | ✅ | 100% |
| Asynq | ✅ | 100% |
| RocketMQ | ✅ | 100% |

### 代码统计
| 层级 | 文件数 | 状态 |
|------|--------|------|
| Service 层 | 7 个 | ✅ 100% |
| Biz 层 | 8 个 | ✅ 100% |
| Data 层 | 8 个 | ✅ 100% |
| Proto 定义 | 13+ 个 | ✅ 100% |
| Wire 生成 | 7 个 | ✅ 100% |

### 测试统计
| 服务 | 测试文件 | 测试用例 | 覆盖率 | 状态 |
|------|----------|----------|--------|------|
| **Inventory** | 3 个 | 45+ | 70%+ | ✅ 100% |
| **Order** | 3 个 | 57+ | 61%+ | ✅ 100% |
| **User** | - | - | - | 🔄 进行中 |
| **Goods** | - | - | - | 🔄 进行中 |
| **UserOp** | - | - | - | 🔄 进行中 |
| **UserAuth** | - | - | - | 🔄 进行中 |

**测试脚本**: 
- ✅ Inventory: `./run_tests.sh` (自动化测试)
- ✅ Order: `./run_tests.sh` (自动化测试)

---

核心功能完成度

### 1. 用户体系 ✅ 100%
- [x] 用户注册/登录
- [x] 用户信息管理
- [x] 密码管理
- [x] JWT 认证
- [x] 图形验证码
- [x] 短信验证码
- [x] Token 刷新/撤销
- [x] 黑名单机制
- [x] 管理员功能

### 2. 商品体系 ✅ 100%
- [x] 商品列表/详情
- [x] 商品搜索
- [x] 商品分类
- [x] 商品管理 (CRUD)
- [x] 批量操作

### 3. 订单体系 ✅ 100%
- [x] 购物车管理
- [x] 订单创建 (事务消息)
- [x] 订单列表/详情
- [x] 订单状态管理
- [x] 订单号生成 (雪花算法)
- [x] 订单超时处理

### 4. 库存体系 ✅ 100%
- [x] 库存设置/查询
- [x] 库存扣减/归还
- [x] 分布式锁 (防超卖)

### 5. 用户操作 ✅ 100%
- [x] 地址管理
- [x] 留言管理
- [x] 收藏管理

### 6. 基础设施 ✅ 100%
- [x] 服务注册与发现 (Consul)
- [x] 配置中心 (Nacos)
- [x] 链路追踪 (Jaeger)
- [x] 分布式缓存 (Redis)
- [x] 数据持久化 (MySQL)
- [x] 异步任务 (Asynq)
- [x] 消息队列 (RocketMQ)

---


架构特点

### 1. ✅ 微服务架构
- **服务拆分**: 6 个独立微服务 + 1 个 API 网关
- **服务治理**: Consul 服务注册与发现
- **配置管理**: Nacos 集中配置管理
- **负载均衡**: gRPC 客户端负载均衡

### 2. ✅ 四层架构
```
HTTP/gRPC Server → Service → Biz → Data → gRPC Client
```
- **职责清晰**: 每层职责明确
- **易于测试**: 依赖注入便于单元测试
- **易于维护**: 代码结构清晰

### 3. ✅ 统一治理
- **认证**: JWT Token 统一认证
- **追踪**: Jaeger 分布式链路追踪
- **日志**: 结构化日志
- **错误**: 统一错误处理
- **验证**: Proto Validate 参数验证

### 4. ✅ 高可用设计
- **服务降级**: 超时控制
- **熔断机制**: gRPC 重试策略
- **分布式锁**: Redis 防止并发问题
- **健康检查**: Consul 健康检查
- **事务消息**: RocketMQ 保证最终一致性

### 5. ✅ 开发效率
- **Wire**: 自动依赖注入
- **Proto**: 自动代码生成
- **Makefile**: 一键构建
- **脚本**: 快速启动/停止

---

## 📋 HTTP API 路由清单 (38个)

### 用户相关 (10个)
```
POST   /api/user/register          ✅
POST   /api/user/login             ✅
GET    /api/user/detail            ✅
PUT    /api/user/update            ✅
PUT    /api/user/update_pwd        ✅
DELETE /api/user/logout            ✅
GET    /api/user/captcha           ✅
POST   /api/user/send_sms          ✅
POST   /api/user/refresh_token     ✅
GET    /api/admin/users            ✅
```

### 商品相关 (3个)
```
GET    /api/goods/list             ✅
GET    /api/goods/{id}             ✅
GET    /api/goods/search           ✅
```

### 订单相关 (4个)
```
POST   /api/order/create           ✅
GET    /api/order/list             ✅
GET    /api/order/{id}             ✅
DELETE /api/order/{id}             ✅
```

### 购物车相关 (4个)
```
GET    /api/cart/list              ✅
POST   /api/cart/add               ✅
PUT    /api/cart/update            ✅
DELETE /api/cart/{id}              ✅
```

### 库存相关 (2个)
```
POST   /api/inventory/set          ✅
GET    /api/inventory/{goodsId}    ✅
```

### 地址相关 (4个)
```
GET    /api/address/list           ✅
POST   /api/address/create         ✅
PUT    /api/address/{id}           ✅
DELETE /api/address/{id}           ✅
```

### 留言相关 (2个)
```
GET    /api/message/list           ✅
POST   /api/message/create         ✅
```

### 收藏相关 (4个)
```
GET    /api/favorite/list          ✅
POST   /api/favorite/add           ✅
DELETE /api/favorite/{goodsId}     ✅
GET    /api/favorite/check/{goodsId} ✅
```

### 监控相关 (1个)
```
GET    /metrics                    ✅
```

**总计**: 38 个 HTTP API ✅

---

## 🎉 项目亮点

### 1. ✅ 完整的微服务架构
- 7 个服务全部实现
- 服务注册与发现
- 负载均衡
- 健康检查

### 2. ✅ 标准的四层架构
- Service → Biz → Data → gRPC
- 职责清晰
- 易于维护

### 3. ✅ 统一的治理方案
- JWT 认证
- 链路追踪
- 日志记录
- 错误处理
- 参数验证

### 4. ✅ 高可用设计
- 服务降级
- 熔断机制
- 超时控制
- 重试策略
- 分布式锁
- 事务消息

### 5. ✅ 开发效率
- Wire 依赖注入
- Proto 自动生成
- 热重载支持
- 完善的文档

### 6. ✅ 消息队列集成
- RocketMQ 事务消息
- 订单超时延迟消息
- 事务回查机制
- 降级方案

### 7. ✅ 完整的测试体系
- **单元测试**: Data/Biz/Service 三层测试
- **集成测试**: 使用 sqlmock 模拟数据库
- **并发测试**: 验证高并发场景
- **自动化脚本**: 一键运行所有测试
- **覆盖率报告**: 自动生成 HTML 报告

**已完成服务**:
- ✅ **Inventory 服务**: 45+ 测试用例，70%+ 覆盖率
- ✅ **Order 服务**: 57+ 测试用例，61%+ 覆盖率

---

## 🚀 快速开始

### 环境要求

- Go 1.23+
- Docker & Docker Compose
- Make

### 1. 启动基础设施

```bash
# 使用分离式部署（推荐）
cd /home/zzx/GoProject/lushop-kratos-main

# 一键部署所有服务
./deploy.sh
# 或
./deploy/scripts/deploy-all.sh

# 或分步部署
./deploy/scripts/deploy-infrastructure.sh  # 基础设施
./deploy/scripts/deploy-services.sh        # 应用服务
```

包含服务：
- MySQL (3306)
- Redis (6379)
- Consul (8500)
- Nacos (8848)
- Jaeger (16686)
- RocketMQ (9876)

### 2. 初始化数据库

```bash
# 导入数据库脚本
mysql -u root -p < scripts/init_db.sql
```

### 3. 启动微服务

```bash
# 方式1: 使用脚本启动所有服务
./scripts/start_all_services.sh

# 方式2: 单独启动服务
cd service/user && kratos run
cd service/goods && kratos run
cd service/order && kratos run
cd service/inventory && kratos run
cd service/userop && kratos run
cd service/userauth && kratos run
```

### 4. 启动 API 网关

```bash
cd lushop
go run main.go
```

### 5. 验证服务

```bash
# 检查服务注册
curl http://localhost:8500/v1/catalog/services

# 测试 API
curl http://localhost:8001/api/goods/list
```

### 6. 访问监控

- **Consul UI**: http://localhost:8500
- **Nacos UI**: http://localhost:8848/nacos (nacos/nacos)
- **Jaeger UI**: http://localhost:16686

### 7. 运行测试（可选）

```bash
# Inventory 服务测试
cd service/inventory
./run_tests.sh

# Order 服务测试
cd service/order
./run_tests.sh

# 查看覆盖率报告
open coverage.html
```

---

## 💼 面试准备

本项目提供完整的面试准备材料，包括项目介绍模板、常见面试问题详解、技术点总结等。

### 快速导航

- 📖 **[面试指南](docs/INTERVIEW_GUIDE.md)** - 面试速查手册（30秒电梯演讲、核心亮点、问题速查）
- 📝 **[面试问题详解](interview/README.md)** - 9个核心面试问题的详细解答
- 📊 **[项目分析总结](docs/PROJECT_ANALYSIS.md)** - 项目优势、存在的问题、改进方向

### 项目介绍模板（30秒）

> "这是一个**微服务电商平台**，使用 Go 和 Kratos 框架开发。系统拆分为 6 个微服务，通过 gRPC 通信。我主要负责订单和库存模块，使用 **RocketMQ 事务消息**保证分布式一致性，用 **Redis 分布式锁**防止超卖。整个项目采用 Consul 服务发现、Nacos 配置中心、Jaeger 链路追踪，代码结构清晰，易于维护。"

### 核心面试问题

详细的面试问题解答请查看 [interview/](interview/) 文件夹：

- [Q1: 为什么选择微服务？](interview/Q1_为什么选择微服务.md)
- [Q2: 如何拆分服务？](interview/Q2_如何拆分服务.md)
- [Q3: 如何保证数据一致性？](interview/Q3_如何保证数据一致性.md)
- [Q4: 如何防止超卖？](interview/Q4_如何防止超卖.md)
- [Q5: 如何实现负载均衡？](interview/Q5_如何实现负载均衡.md)
- [Q6: 如何排查问题？](interview/Q6_如何排查问题.md)
- [Q7: 遇到的最大挑战？](interview/Q7_遇到的最大挑战.md)
- [Q8: 如果重新设计会怎么做？](interview/Q8_如果重新设计会怎么做.md)
- [Q9: 如何保证代码质量？](interview/Q9_如何保证代码质量.md)

更多面试技巧和常见追问，请查看 [面试指南](docs/INTERVIEW_GUIDE.md)。

---

## 📚 相关文档

### 📖 项目文档
- [📊 项目分析总结](docs/PROJECT_ANALYSIS.md) - **项目优势、存在的问题、改进方向** ⭐
- [🐳 Docker 部署指南](docs/DOCKER_DEPLOY.md) - 完整部署文档
- [🖥️ Ubuntu 24.04 单机部署指南](docs/UBUNTU_SINGLE_NODE_DEPLOY.md) - **Ubuntu 服务器单机部署详细教程** ⭐
- [☸️ Kubernetes 部署指引](k8s/README.md) - Kubernetes 集群部署文档
- [📂 项目结构说明](docs/PROJECT_STRUCTURE.md) - 项目目录结构
- [🧪 测试计划](docs/LUSHOP_TESTING_PLAN.md) - 测试策略和计划
- [🛡️ Sentinel 限流调用流程](lushop/docs/SENTINEL_CALL_FLOW.md) - **Sentinel 限流实现详解** ⭐

### 💼 面试准备
- [📝 面试问题详解](interview/README.md) - 9个核心面试问题详细解答
- [🎯 面试指南](docs/INTERVIEW_GUIDE.md) - 面试速查手册（30秒电梯演讲、核心亮点、问题速查）

---

## 📚 技术深度扩展

### 性能优化

1. **数据库优化**
   - 索引优化（商品ID、用户ID等）
   - 读写分离（主从复制）
   - 分库分表（订单表按用户ID分片）

2. **缓存策略**
   - 商品信息缓存（Redis）
   - 用户Session缓存
   - 热点数据预加载

3. **并发控制**
   - 连接池管理（MySQL/Redis）
   - gRPC 连接复用
   - 协程池限制

### 高可用设计

1. **服务降级**
   - 非核心功能降级（如推荐服务）
   - 超时快速失败
   - 默认值返回

2. **熔断机制**
   - gRPC 重试策略
   - 熔断器模式
   - **Sentinel 限流保护**（✅ 已实现）

3. **容灾方案**
   - 多副本部署
   - 跨机房部署
   - 数据备份策略

---

## 🎓 学习收获

通过这个项目，我深入理解了：

1. **微服务架构**
   - 服务拆分原则
   - 服务间通信
   - 分布式事务
   - 服务治理

2. **Go 语言特性**
   - 并发编程（goroutine/channel）
   - 接口设计
   - 错误处理
   - 性能优化

3. **分布式系统**
   - CAP 理论
   - 最终一致性
   - 分布式锁
   - 消息队列

4. **工程实践**
   - 代码规范
   - 单元测试
   - CI/CD
   - 文档编写

---

## 📞 联系方式

如有问题或建议，欢迎联系：

- **GitHub**: [Your GitHub]
- **Email**: [Your Email]
- **Blog**: [Your Blog]

---

## 📄 License

MIT License

---

**⭐ 如果这个项目对你有帮助，欢迎 Star！**
