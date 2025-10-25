# 🎯 Lushop 项目面试指南

## 📋 项目概览一句话

**基于 Go-Kratos 的生产级微服务电商平台，采用 DDD 设计，实现了分布式事务、服务治理、链路追踪等企业级特性。**

---

## 🎤 30秒电梯演讲

"这是一个**微服务电商平台**，使用 Go 和 Kratos 框架开发。系统拆分为 6 个微服务，通过 gRPC 通信。我主要负责订单和库存模块，使用 **RocketMQ 事务消息**保证分布式一致性，用 **Redis 分布式锁**防止超卖。整个项目采用 Consul 服务发现、Nacos 配置中心、Jaeger 链路追踪，代码结构清晰，易于维护。"

---

## 💡 核心技术亮点（必须掌握）

### 1. 分布式事务 ⭐⭐⭐⭐⭐
- **技术**：RocketMQ 事务消息
- **场景**：订单创建 + 库存扣减
- **关键词**：半消息、本地事务、事务回查、最终一致性

### 2. 防超卖 ⭐⭐⭐⭐⭐
- **技术**：Redis 分布式锁 + 数据库乐观锁
- **场景**：高并发库存扣减
- **关键词**：分布式锁、乐观锁、version 字段

### 3. 服务治理 ⭐⭐⭐⭐
- **技术**：Consul + gRPC
- **功能**：服务注册、发现、健康检查、负载均衡
- **关键词**：服务发现、客户端负载均衡、健康检查

### 4. 配置管理 ⭐⭐⭐⭐
- **技术**：Nacos
- **功能**：配置中心、热更新
- **关键词**：配置热更新、多环境隔离

### 5. 链路追踪 ⭐⭐⭐⭐
- **技术**：Jaeger
- **功能**：分布式追踪、性能分析
- **关键词**：TraceID、Span、调用链

### 6. 依赖注入 ⭐⭐⭐
- **技术**：Google Wire
- **优势**：编译期检查、代码生成
- **关键词**：依赖注入、ProviderSet

---

## 🔥 高频面试问题速查

### Q1: 分布式事务如何保证一致性？

**关键点**：
- RocketMQ 事务消息（半消息机制）
- 本地事务表
- 事务回查机制
- 最终一致性 vs 强一致性

**回答框架**：
1. 说明场景（订单+库存）
2. 介绍方案（事务消息）
3. 详细流程（4步骤）
4. 优势总结（3点）

### Q2: 如何防止库存超卖？

**关键点**：
- Redis 分布式锁（高并发）
- 数据库乐观锁（中低并发）
- 两种方案对比

**回答框架**：
1. 问题分析（并发扣减）
2. 方案1（Redis锁+代码）
3. 方案2（乐观锁+SQL）
4. 方案选择（场景对比）

### Q3: 微服务如何调用？

**关键点**：
- gRPC 协议
- Consul 服务发现
- 客户端负载均衡
- 健康检查

**回答框架**：
1. 注册流程
2. 发现流程
3. 调用流程
4. 优势说明

### Q4: 如何排查微服务问题？

**关键点**：
- Jaeger 链路追踪
- TraceID 传递
- 可视化分析
- 性能定位

**回答框架**：
1. 工具介绍（Jaeger）
2. 使用方式（TraceID查询）
3. 实际案例（定位慢接口）
4. 效果说明

### Q5: 项目最大挑战？

**推荐回答**：订单超时自动取消

**关键点**：
- 业务场景（30分钟未支付）
- 技术难点（不能轮询）
- 解决方案（延迟消息）
- 实现细节（RocketMQ）

---

## 📊 技术栈速记卡

| 分类 | 技术 | 用途 | 版本 |
|------|------|------|------|
| **语言** | Go | 后端开发 | 1.23+ |
| **框架** | Kratos | 微服务框架 | v2.8.3 |
| **通信** | gRPC | 服务间调用 | 1.69 |
| **注册中心** | Consul | 服务发现 | - |
| **配置中心** | Nacos | 配置管理 | - |
| **数据库** | MySQL | 持久化 | 8.0 |
| **缓存** | Redis | 缓存/锁 | 7.0 |
| **消息队列** | RocketMQ | 事务消息 | - |
| **链路追踪** | Jaeger | 分布式追踪 | - |
| **依赖注入** | Wire | DI | - |

---

## 🏗️ 架构图速记

```
API Gateway (8001)
    ↓
┌───────┬───────┬───────┬───────┬───────┬───────┐
│ User  │Goods  │Order  │Inven  │UserOp │Auth   │
│ 50051 │ 50052 │ 50053 │ 50054 │ 50055 │ 50056 │
└───────┴───────┴───────┴───────┴───────┴───────┘
    ↓       ↓       ↓       ↓       ↓       ↓
┌─────────────────────────────────────────────┐
│  Consul | Nacos | MySQL | Redis | RocketMQ │
└─────────────────────────────────────────────┘
```

---

## 💼 项目介绍模板（背诵版）

### 版本1：简短版（30秒）

"我做过一个**微服务电商平台**，使用 **Go + Kratos** 开发。系统拆分为 6 个服务，通过 **gRPC** 通信。我负责订单和库存模块，用 **RocketMQ 事务消息**保证分布式一致性，用 **Redis 分布式锁**防止超卖。项目采用 **Consul** 服务发现、**Nacos** 配置中心、**Jaeger** 链路追踪。"

### 版本2：详细版（1-2分钟）

"我参与开发了一个基于 Go 语言的**微服务电商平台**。

**技术架构**：采用 **Kratos 框架**，将系统拆分为用户、商品、订单、库存等 6 个独立微服务，通过 **gRPC** 进行高性能通信。

**核心技术**：
- 使用 **Consul** 实现服务注册与发现，支持健康检查和客户端负载均衡
- 使用 **Nacos** 作为配置中心，实现配置热更新
- 使用 **Jaeger** 进行分布式链路追踪，快速定位问题

**我的贡献**：主要负责订单和库存模块，解决了几个关键技术难点：

1. **分布式事务**：使用 **RocketMQ 事务消息**保证订单创建和库存扣减的最终一致性，支持事务回查机制
2. **防超卖**：通过 **Redis 分布式锁**防止高并发下的库存超卖问题
3. **订单超时**：使用 **RocketMQ 延迟消息**实现订单超时自动取消

**项目成果**：系统支持高并发访问，服务可用性达到 99.9%，代码结构清晰，易于维护和扩展。"

---

## 🎯 面试准备 Checklist

### 必须能画出的图

- [ ] 系统架构图
- [ ] 服务调用流程图
- [ ] 分布式事务流程图
- [ ] 服务注册发现流程图

### 必须能解释的概念

- [ ] 微服务架构
- [ ] DDD 领域驱动设计
- [ ] 四层架构（Transport/Service/Biz/Data）
- [ ] gRPC vs HTTP
- [ ] 服务注册与发现
- [ ] 分布式事务
- [ ] 最终一致性
- [ ] CAP 理论

### 必须能回答的问题

- [ ] 为什么选择微服务？
- [ ] 如何拆分服务？
- [ ] 如何保证数据一致性？
- [ ] 如何防止超卖？
- [ ] 如何实现负载均衡？
- [ ] 如何排查问题？
- [ ] 遇到的最大挑战？
- [ ] 如果重新设计会怎么做？

### 必须准备的代码

- [ ] 分布式锁实现
- [ ] 事务消息发送
- [ ] 服务注册代码
- [ ] Wire 依赖注入

---

## 🚀 加分项

### 1. 性能优化经验
- 数据库索引优化
- Redis 缓存策略
- 连接池管理
- 并发控制

### 2. 高可用设计
- 服务降级
- 熔断机制
- 限流保护
- 容灾方案

### 3. 监控告警
- Prometheus 指标
- Grafana 可视化
- 告警规则配置

### 4. DevOps 实践
- Docker 容器化
- Kubernetes 部署
- CI/CD 流程
- 自动化测试

---

## 📝 常见追问及应对

### Q: 为什么不用 2PC/3PC？

**答**：
- 2PC/3PC 是强一致性方案，性能差，不适合高并发场景
- 我们的业务允许短暂不一致，最终一致性即可
- RocketMQ 事务消息性能更好，可用性更高

### Q: Redis 锁如果宕机怎么办？

**答**：
- 使用 Redis 主从+哨兵模式，保证高可用
- 设置合理的锁超时时间，避免死锁
- 也可以用 Redlock 算法（多个 Redis 实例）
- 极端情况可以降级到数据库乐观锁

### Q: 如何保证消息不丢失？

**答**：
- 生产者：同步发送 + 重试机制
- Broker：消息持久化到磁盘
- 消费者：手动 ACK，消费成功后确认
- 还可以配置消息备份（主从同步）

### Q: 服务拆分的原则是什么？

**答**：
- 按业务领域拆分（DDD）
- 单一职责原则
- 高内聚低耦合
- 考虑团队规模和能力

### Q: 运行过程中遇到过 goroutine 泄漏或 panic 吗？

**答**：
是的，在开发和测试过程中确实遇到过这两类问题，这也让我对 Go 的并发编程有了更深的理解。

#### 1. Goroutine 泄漏案例

**问题场景**：
在实现订单超时检查功能时，我发现服务运行一段时间后内存持续增长，通过 `pprof` 分析发现有大量 goroutine 没有正常退出。

**原因分析**：
```go
// ❌ 错误代码
func (s *OrderService) CheckTimeout(ctx context.Context, orderId int64) {
    go func() {
        ticker := time.NewTicker(30 * time.Minute)
        for {
            select {
            case <-ticker.C:
                // 检查订单状态
                s.checkAndCancelOrder(orderId)
            }
        }
    }()
}
```

问题：
- ticker 没有 stop，导致 goroutine 永远不会退出
- 没有监听 context 的 Done 信号
- 每次调用都会创建新的 goroutine

**解决方案**：
```go
// ✅ 正确代码
func (s *OrderService) CheckTimeout(ctx context.Context, orderId int64) error {
    ticker := time.NewTicker(30 * time.Minute)
    defer ticker.Stop()  // 确保 ticker 被停止
    
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Errorf("panic in CheckTimeout: %v", r)
            }
        }()
        
        for {
            select {
            case <-ticker.C:
                s.checkAndCancelOrder(orderId)
            case <-ctx.Done():  // 监听 context 取消
                log.Info("CheckTimeout stopped")
                return
            }
        }
    }()
    
    return nil
}
```

**排查工具**：
```bash
# 1. 使用 pprof 查看 goroutine 数量
curl http://localhost:6060/debug/pprof/goroutine?debug=1

# 2. 生成 goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 3. 查看具体泄漏位置
(pprof) top
(pprof) list CheckTimeout
```

#### 2. Panic 异常案例

**问题场景 1：空指针引用**

在处理 gRPC 请求时，没有判断请求参数是否为 nil：

```go
// ❌ 错误代码
func (s *GoodsService) GetGoodsDetail(ctx context.Context, req *pb.GoodsRequest) (*pb.GoodsResponse, error) {
    // 如果 req 为 nil，这里会 panic
    goods, err := s.uc.GetGoods(ctx, req.Id)
    return &pb.GoodsResponse{Data: goods}, err
}
```

**解决方案**：
```go
// ✅ 正确代码
func (s *GoodsService) GetGoodsDetail(ctx context.Context, req *pb.GoodsRequest) (*pb.GoodsResponse, error) {
    // 1. 参数校验
    if req == nil {
        return nil, errors.New("request is nil")
    }
    
    // 2. 使用 defer recover 捕获 panic
    defer func() {
        if r := recover(); r != nil {
            log.Errorf("panic in GetGoodsDetail: %v, stack: %s", r, debug.Stack())
        }
    }()
    
    goods, err := s.uc.GetGoods(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    
    return &pb.GoodsResponse{Data: goods}, nil
}
```

**问题场景 2：并发写 map**

在缓存实现中，多个 goroutine 同时写 map 导致 panic：

```go
// ❌ 错误代码
type Cache struct {
    data map[string]interface{}
}

func (c *Cache) Set(key string, value interface{}) {
    c.data[key] = value  // concurrent map writes panic
}
```

**解决方案**：
```go
// ✅ 正确代码 - 方案1：使用 sync.Map
type Cache struct {
    data sync.Map
}

func (c *Cache) Set(key string, value interface{}) {
    c.data.Store(key, value)
}

// ✅ 正确代码 - 方案2：使用互斥锁
type Cache struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

#### 3. 预防措施

**开发阶段**：
1. **代码审查**：重点关注 goroutine 的生命周期管理
2. **单元测试**：使用 `-race` 标志检测数据竞争
   ```bash
   go test -race ./...
   ```
3. **静态分析**：使用 `golangci-lint` 检查潜在问题
   ```bash
   golangci-lint run --enable=govet,errcheck,staticcheck
   ```

**运行时监控**：
1. **Prometheus 指标**：
   ```go
   // 监控 goroutine 数量
   prometheus.NewGaugeFunc(prometheus.GaugeOpts{
       Name: "go_goroutines",
       Help: "Number of goroutines",
   }, func() float64 {
       return float64(runtime.NumGoroutine())
   })
   ```

2. **告警规则**：
   ```yaml
   # Prometheus 告警
   - alert: GoroutineLeaking
     expr: go_goroutines > 10000
     for: 5m
     annotations:
       summary: "Goroutine 泄漏告警"
   ```

3. **定期 pprof 分析**：
   ```go
   // 启用 pprof
   import _ "net/http/pprof"
   
   go func() {
       log.Println(http.ListenAndServe("localhost:6060", nil))
   }()
   ```

**最佳实践**：
1. ✅ 总是使用 `context.Context` 控制 goroutine 生命周期
2. ✅ 及时调用 `ticker.Stop()` 和 `cancel()`
3. ✅ 使用 `defer recover()` 捕获 panic
4. ✅ 并发访问共享资源时使用锁或 channel
5. ✅ 使用 `errgroup` 管理多个 goroutine
6. ✅ 定期使用 pprof 分析内存和 goroutine

**经验总结**：
通过这些问题，我学会了：
- 🎯 Go 并发编程的最佳实践
- 🎯 使用 pprof 等工具排查问题
- 🎯 建立完善的监控告警机制
- 🎯 编写更健壮的代码

---

## 🎓 学习建议

### 深入理解的技术点

1. **分布式事务**
   - 2PC/3PC/TCC/Saga
   - 事务消息原理
   - 最终一致性

2. **分布式锁**
   - Redis 实现
   - Redlock 算法
   - ZooKeeper 实现

3. **服务治理**
   - 服务注册发现
   - 负载均衡算法
   - 健康检查机制

4. **微服务模式**
   - API Gateway
   - Service Mesh
   - 熔断降级

---

## 📚 推荐阅读

- 《微服务架构设计模式》
- 《Go 语言高级编程》
- 《分布式系统原理与范型》
- Kratos 官方文档
- RocketMQ 官方文档

---

**💪 祝你面试顺利！记住：自信、清晰、有条理！**
