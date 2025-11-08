package sentinel

import (
	"context"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// ErrRateLimitExceeded 限流错误
var ErrRateLimitExceeded = errors.New(429, "RATE_LIMIT_EXCEEDED", "请求过于频繁，请稍后再试")

// Option 配置选项
type Option func(*options)

type options struct {
	enabled   bool
	whitelist map[string]bool
	logger    log.Logger
}

// WithEnabled 是否启用 Sentinel
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

// Server Sentinel 限流中间件
func Server(opts ...Option) middleware.Middleware {
	o := &options{
		enabled:   true,
		whitelist: make(map[string]bool),
		logger:    log.DefaultLogger,
	}

	for _, opt := range opts {
		opt(o)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 如果不启用，直接通过
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

			// 使用接口路径作为资源名
			resource := endpoint

			// Sentinel 限流检查
			entry, blockErr := api.Entry(resource)
			if blockErr != nil {
				// 限流被触发
				log.NewHelper(o.logger).Warnw(
					"msg", "Rate limit exceeded",
					"endpoint", endpoint,
					"resource", resource,
					"blockType", blockErr.BlockType().String(),
				)
				return nil, ErrRateLimitExceeded
			}

			// 执行处理
			reply, err = handler(ctx, req)

			// 记录错误（用于后续的熔断等功能）
			if err != nil {
				entry.SetError(err)
			}

			// 退出 Sentinel 资源
			entry.Exit()

			return reply, err
		}
	}
}
