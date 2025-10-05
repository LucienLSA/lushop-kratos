package server

import (
	"context"
	v1 "lushop/api/lushop/v1"
	"lushop/internal/conf"
	"lushop/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
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
		),
		http.Filter(handlers.CORS( // 浏览器跨域
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
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
