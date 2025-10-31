package ratelimit

import (
	"context"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/go-kratos/kratos/v2/log"
)

// ErrRateLimitExceeded 限流错误
var ErrRateLimitExceeded = errors.New(429, "RATE_LIMIT_EXCEEDED", "请求过于频繁，请稍后再试")

// Option 限流配置选项
type Option func(*options)

type options struct {
	limiter   ratelimit.Limiter
	enabled   bool
	whitelist map[string]bool
	logger    log.Logger
}

// WithLimiter 使用自定义限流器
func WithLimiter(limiter ratelimit.Limiter) Option {
	return func(o *options) {
		o.limiter = limiter
	}
}

// WithEnabled 是否启用限流
func WithEnabled(enabled bool) Option {
	return func(o *options) {
		o.enabled = enabled
	}
}

// WithWhitelist 设置白名单（不限流）
func WithWhitelist(endpoints []string) Option {
	return func(o *options) {
		o.whitelist = make(map[string]bool)
		for _, endpoint := range endpoints {
			o.whitelist[endpoint] = true
		}
	}
}

// WithLogger 设置日志
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// Server 限流中间件
func Server(opts ...Option) middleware.Middleware {
	o := &options{
		enabled:   true,
		whitelist: make(map[string]bool),
		logger:    log.DefaultLogger,
	}

	for _, opt := range opts {
		opt(o)
	}

	// 如果没有自定义限流器，使用默认BBR限流器
	if o.limiter == nil {
		o.limiter = bbr.NewLimiter()
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			// 如果不启用限流，直接通过
			if !o.enabled {
				return handler(ctx, req)
			}

			// 获取请求信息
			info, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			endpoint := info.Operation()

			// 检查白名单
			if o.whitelist[endpoint] {
				return handler(ctx, req)
			}

			// 使用限流器
			done, e := o.limiter.Allow()
			if e != nil {
				// 限流拒绝
				log.Log(log.LevelWarn,
					"msg", "Rate limit exceeded",
					"endpoint", endpoint,
					"error", e,
				)
				return nil, ErrRateLimitExceeded
			}

			// 执行处理
			reply, err = handler(ctx, req)

			// 报告限流结果
			done(ratelimit.DoneInfo{
				Err: err,
			})

			return reply, err
		}
	}
}
