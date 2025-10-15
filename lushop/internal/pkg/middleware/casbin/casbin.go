package casbin

import (
	"context"
	"errors"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	kmjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

type subjectExtractor func(ctx context.Context) (string, error)

// 从 JWT Claims 中提取“主体”，建议用角色名或用户ID；这里用角色名
func RoleFromJWT() subjectExtractor {
	return func(ctx context.Context) (string, error) {
		claims, ok := kmjwt.FromContext(ctx)
		if !ok {
			return "", errors.New("no jwt in context")
		}
		c, ok := claims.(jwt5.MapClaims)
		if !ok {
			return "", errors.New("invalid claims")
		}
		// 你当前 Claims 中是 AuthorityId(int)，这里将其映射成角色字符串
		aid, ok := c["AuthorityId"]
		if !ok {
			return "", errors.New("no AuthorityId")
		}
		role := mapAuthorityToRole(aid) // 自定义：1->admin, 2->user 等
		return role, nil
	}
}

func mapAuthorityToRole(a any) string {
	// 依据你的项目约定：例如 1=admin, 2=user
	switch v := a.(type) {
	case float64:
		if int(v) == 1 {
			return "admin"
		}
		return "user"
	default:
		return "user"
	}
}

func Middleware(e *casbin.Enforcer, extract subjectExtractor) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			sub, err := extract(ctx)
			if err != nil {
				return nil, err
			}
			// 从上下文/传输层获取资源和动作
			// Kratos http 请求可通过 metadata 或者在自定义拦截器中取 path/method
			// 这里做一个通用推断：从错误栈、req类型或 context 中注入，你可按项目现状取 path/method。
			// 假设已在上游注入 path/method 到 context 的 keys 里（示例）：
			path := getPathFromContext(ctx)  // 需实现
			act := getMethodFromContext(ctx) // 需实现

			allowed, err := e.Enforce(sub, path, strings.ToUpper(act))
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, errors.New("forbidden")
			}
			return handler(ctx, req)
		}
	}
}

// 你可以在自定义 http 服务器初始化时，把 path/method 放入 context。
// 这里只是占位，按你的 http.go 实际实现来获取。
func getPathFromContext(ctx context.Context) string {
	if req, ok := khttp.RequestFromServerContext(ctx); ok && req != nil {
		return req.URL.Path
	}
	return ""
}
func getMethodFromContext(ctx context.Context) string {
	if req, ok := khttp.RequestFromServerContext(ctx); ok && req != nil {
		return req.Method
	}
	return ""
}
