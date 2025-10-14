package server

import (
	"context"
	"encoding/json"
	v1 "lushop/api/lushop/v1"
	http2 "lushop/internal/biz/http"
	"lushop/internal/conf"
	"lushop/internal/conf/metrix"
	"lushop/internal/pkg/middleware/auth"
	"lushop/internal/pkg/middleware/i18n"
	"lushop/internal/service"
	httpNet "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, ac *conf.Auth, s *service.LushopService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		// 设置请求解码器，处理GET请求的Content-Type问题
		http.RequestDecoder(func(r *httpNet.Request, v interface{}) error {
			// 对于GET请求，不需要检查Content-Type
			if r.Method == "GET" {
				return nil
			}
			// 对于其他请求，使用默认解码器
			return http.DefaultRequestDecoder(r, v)
		}),
		http.Middleware(
			recovery.Recovery(),
			i18n.Translator(),
			validate.Validator(), // 接口访问的参数校验
			tracing.Server(),     // 链路追踪
			// JWT 验证 - 对所有需要认证的接口生效
			selector.Server(
				jwt.Server(func(token *jwt5.Token) (interface{}, error) {
					return []byte(ac.JwtKey), nil
				}, jwt.WithSigningMethod(jwt5.SigningMethodHS256)),
			).Match(NewAuthMatcher()).Build(),
			// 统一的权限中间件 - 自动判断用户/管理员权限
			selector.Server(
				auth.AuthMiddleware(),
			).Match(NewAuthMatcher()).Build(),
			// 日志
			logging.Server(logger),
			// 指标 prometheus
			metrics.Server(
				metrics.WithSeconds(metrix.MetricSeconds),
				metrics.WithRequests(metrix.MetricRequests),
			),
		),
		http.Filter(handlers.CORS( // 浏览器跨域
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
		http.ErrorEncoder(
			func(writer httpNet.ResponseWriter, request *httpNet.Request, err error) {
				log.Infof("拦截到的错误信息是：%s", err.Error())
				log.Infof("请求URL: %s", request.URL.String())
				log.Infof("请求方法: %s", request.Method)
				log.Infof("请求Content-Type: %s", request.Header.Get("Content-Type"))
				log.Infof("请求User-Agent: %s", request.Header.Get("User-Agent"))

				message := extractMessageFromError(err)
				reply := &http2.BaseResponse{
					Code: 400,
					Msg:  message,
					Data: nil,
				}

				// 直接使用标准库JSON编码，避免编码器问题
				data, _ := json.Marshal(reply)
				writer.Header().Set("Content-Type", "application/json")
				writer.Write(data)
			}),
		http.ResponseEncoder(func(writer httpNet.ResponseWriter, request *httpNet.Request, i interface{}) error {
			reply := &http2.BaseResponse{
				Code: 200,
				Msg:  "请求成功",
				Data: i,
			}
			// 直接使用标准库JSON编码，避免编码器问题
			data, _ := json.Marshal(reply)
			writer.Header().Set("Content-Type", "application/json")
			writer.Write(data)
			return nil
		}),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	handler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", handler)
	srv.Handle("/metrics", promhttp.Handler())
	v1.RegisterLushopHTTPServer(srv, s)
	return srv
}

// 公开接口 - 无需认证
var publicEndpoints = map[string]struct{}{
	"/lushop.lushop.v1.Lushop/Captcha":  {},
	"/lushop.lushop.v1.Lushop/Login":    {},
	"/lushop.lushop.v1.Lushop/Register": {},
}

// NewAuthMatcher 需要JWT认证的接口匹配器
func NewAuthMatcher() selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		// 公开接口不需要JWT认证
		if _, ok := publicEndpoints[operation]; ok {
			return false
		}
		// 其他接口都需要JWT认证
		return true
	}
}

func extractMessageFromError(err error) string {
	marshal, err2 := json.Marshal(err)
	if err2 != nil {
		return "系统错误"
	}
	var em ErrorMessage
	e := json.Unmarshal(marshal, &em)
	if e != nil {
		return "系统错误"
	}
	return em.Message
}

type ErrorMessage struct {
	Message string `json:"message"`
}
