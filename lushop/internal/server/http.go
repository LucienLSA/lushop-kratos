package server

import (
	"context"
	"encoding/json"
	v1 "lushop/api/lushop/v1"
	"lushop/internal/conf"
	"lushop/internal/conf/metrix"
	"lushop/internal/service"

	http2 "lushop/internal/biz/http"
	httpNet "net/http"

	"github.com/go-kratos/kratos/v2/encoding"
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
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(), // 接口访问的参数校验
			tracing.Server(),     // 链路追踪
			selector.Server( // jwt 验证
				jwt.Server(func(token *jwt5.Token) (interface{}, error) {
					return []byte(ac.JwtKey), nil
				}, jwt.WithSigningMethod(jwt5.SigningMethodHS256)),
			).Match(NewWhiteListMatcher()).Build(),
			logging.Server(logger),
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
				message := extractMessageFromError(err)
				reply := &http2.BaseResponse{
					Code: 400,
					Msg:  message,
					Data: nil,
				}
				codec := encoding.GetCodec("json")
				data, _ := codec.Marshal(reply)
				writer.Header().Set("Content-Type", "application/json")
				writer.Write(data)
			}),
		http.ResponseEncoder(func(writer httpNet.ResponseWriter, request *httpNet.Request, i interface{}) error {
			reply := &http2.BaseResponse{
				Code: 200,
				Msg:  "请求成功",
				Data: i,
			}
			codec := encoding.GetCodec("json")
			data, _ := codec.Marshal(reply)
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

// NewWhiteListMatcher 白名单不需要token验证的接口
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/lushop.lushop.v1.Lushop/Captcha"] = struct{}{}
	whiteList["/lushop.lushop.v1.Lushop/Login"] = struct{}{}
	whiteList["/lushop.lushop.v1.Lushop/Register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
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
