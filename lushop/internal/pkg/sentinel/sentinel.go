package sentinel

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/go-kratos/kratos/v2/log"
)

// Init 初始化 Sentinel
func Init(logger log.Logger) error {
	err := api.InitDefault()
	if err != nil {
		return err
	}
	log.NewHelper(logger).Info("Sentinel initialized successfully")
	return nil
}

// LoadDefaultRules 加载默认限流规则
// 为常用接口加载默认限流规则
func LoadDefaultRules() error {
	// 为常用接口配置默认限流规则
	rules := []*FlowRuleConfig{
		// 订单相关接口
		{
			Resource:         "/lushop.lushop.v1.Order/CreateOrder",
			Threshold:        1000,
			ControlBehavior:  0, // Reject
			StatIntervalInMs: 1000,
		},
		// 购物车相关接口
		{
			Resource:         "/lushop.lushop.v1.Cart/AddCart",
			Threshold:        2000,
			ControlBehavior:  0,
			StatIntervalInMs: 1000,
		},
		{
			Resource:         "/lushop.lushop.v1.Cart/UpdateCart",
			Threshold:        2000,
			ControlBehavior:  0,
			StatIntervalInMs: 1000,
		},
		// 库存相关接口
		{
			Resource:         "/lushop.lushop.v1.Inventory/SetInv",
			Threshold:        500,
			ControlBehavior:  0,
			StatIntervalInMs: 1000,
		},
		// 用户操作相关接口
		{
			Resource:         "/lushop.lushop.v1.UserOp/CreateAddress",
			Threshold:        1000,
			ControlBehavior:  0,
			StatIntervalInMs: 1000,
		},
	}

	return LoadFlowRules(rules)
}

// LoadFlowRules 加载限流规则
func LoadFlowRules(rules []*FlowRuleConfig) error {
	if len(rules) == 0 {
		return nil
	}

	flowRules := make([]*flow.Rule, 0, len(rules))
	for _, rule := range rules {
		flowRules = append(flowRules, &flow.Rule{
			Resource:               rule.Resource,
			Threshold:              rule.Threshold,
			RelationStrategy:       flow.RelationStrategy(rule.RelationStrategy),
			TokenCalculateStrategy: flow.TokenCalculateStrategy(rule.TokenCalculateStrategy),
			ControlBehavior:        flow.ControlBehavior(rule.ControlBehavior),
			StatIntervalInMs:       rule.StatIntervalInMs,
		})
	}

	_, err := flow.LoadRules(flowRules)
	return err
}

// LoadCircuitBreakerRules 加载熔断规则
func LoadCircuitBreakerRules(rules []*CircuitBreakerRuleConfig) error {
	if len(rules) == 0 {
		return nil
	}

	cbRules := make([]*circuitbreaker.Rule, 0, len(rules))
	for _, rule := range rules {
		var strategy circuitbreaker.Strategy
		switch rule.Strategy {
		case "SlowRequestRatio":
			strategy = circuitbreaker.SlowRequestRatio
		case "ErrorRatio":
			strategy = circuitbreaker.ErrorRatio
		case "ErrorCount":
			strategy = circuitbreaker.ErrorCount
		default:
			strategy = circuitbreaker.SlowRequestRatio
		}

		cbRules = append(cbRules, &circuitbreaker.Rule{
			Resource:         rule.Resource,
			Strategy:         strategy,
			RetryTimeoutMs:   rule.RetryTimeoutMs,
			MinRequestAmount: rule.MinRequestAmount,
			StatIntervalMs:   rule.StatIntervalMs,
			MaxAllowedRtMs:   rule.MaxAllowedRtMs,
			Threshold:        rule.Threshold,
		})
	}

	_, err := circuitbreaker.LoadRules(cbRules)
	return err
}

// LoadSystemRules 加载系统规则
func LoadSystemRules(rules []*SystemRuleConfig) error {
	if len(rules) == 0 {
		return nil
	}

	systemRules := make([]*system.Rule, 0, len(rules))
	for _, rule := range rules {
		systemRules = append(systemRules, &system.Rule{
			MetricType:   system.MetricType(rule.MetricType),
			TriggerCount: rule.TriggerCount,
		})
	}

	_, err := system.LoadRules(systemRules)
	return err
}

// FlowRuleConfig 限流规则配置
type FlowRuleConfig struct {
	Resource               string  // 资源名（接口路径）
	Threshold              float64 // 阈值（QPS）
	RelationStrategy       int32   // 关系策略（0:Current, 1:Associated）
	TokenCalculateStrategy int32   // Token计算策略（0:Direct, 1:WarmUp）
	ControlBehavior        int32   // 控制行为（0:Reject, 1:Throttling）
	StatIntervalInMs       uint32  // 统计间隔（毫秒）
}

// CircuitBreakerRuleConfig 熔断规则配置
type CircuitBreakerRuleConfig struct {
	Resource         string  // 资源名（接口路径）
	Strategy         string  // 策略（SlowRequestRatio/ErrorRatio/ErrorCount）
	RetryTimeoutMs   uint32  // 重试超时时间（毫秒）
	MinRequestAmount uint64  // 最小请求数
	StatIntervalMs   uint32  // 统计间隔（毫秒）
	MaxAllowedRtMs   uint64  // 最大允许响应时间（毫秒）
	Threshold        float64 // 阈值
}

// SystemRuleConfig 系统规则配置
type SystemRuleConfig struct {
	MetricType   int32   // 指标类型（0:Load, 1:AvgRT, 2:Concurrency, 3:InboundQPS, 4:CpuUsage）
	TriggerCount float64 // 触发阈值
}
