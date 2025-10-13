package auth

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	ID          int64
	AuthorityId int
	NickName    string
	jwt5.RegisteredClaims
}

// CreateToken generate token
func CreateToken(c CustomClaims, key string) (string, error) {
	claims := jwt5.NewWithClaims(jwt5.SigningMethodHS256, c)
	signedString, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("generate token failed" + err.Error())
	}
	return signedString, nil
}

// 定义上下文键类型，避免使用字符串作为键
type contextKey string

const (
	userIDKey       contextKey = "user_id"
	userRoleKey     contextKey = "user_role"
	userNicknameKey contextKey = "user_nickname"
)

// AuthMiddleware 统一的权限中间件 - 根据接口路径自动判断权限要求
func AuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 从JWT token中获取用户信息
			claims, ok := jwt.FromContext(ctx)
			if !ok {
				return nil, errors.New("unauthorized: no token found")
			}

			// 检查claims类型
			c, ok := claims.(jwt5.MapClaims)
			if !ok {
				return nil, errors.New("unauthorized: invalid token claims")
			}

			// 获取用户权限ID
			authorityId, ok := c["AuthorityId"]
			if !ok {
				return nil, errors.New("unauthorized: no authority id in token")
			}

			// 验证用户ID存在
			userID, ok := c["ID"]
			if !ok {
				return nil, errors.New("unauthorized: no user id in token")
			}

			// 检查用户是否在黑名单中
			// 注意：这里暂时跳过黑名单检查，因为中间件无法直接访问仓库
			// 在实际项目中，应该通过依赖注入或其他方式访问仓库来实现黑名单检查

			// 将用户信息存储到上下文中
			ctx = context.WithValue(ctx, userIDKey, int64(userID.(float64)))
			ctx = context.WithValue(ctx, userRoleKey, int(authorityId.(float64)))
			if nickname, exists := c["NickName"]; exists {
				ctx = context.WithValue(ctx, userNicknameKey, nickname.(string))
			}

			// 从上下文中获取接口路径
			if operation, exists := ctx.Value("operation").(string); exists {
				// 检查是否为管理员接口
				if isAdminEndpoint(operation) {
					// 管理员接口需要管理员权限
					if int(authorityId.(float64)) != RoleAdmin {
						return nil, errors.New("forbidden: admin access required")
					}
				}
				// 用户接口只需要登录即可，不需要额外检查
			}

			return handler(ctx, req)
		}
	}
}

// isAdminEndpoint 检查是否为管理员接口
func isAdminEndpoint(operation string) bool {
	// 管理员接口列表
	adminEndpoints := map[string]struct{}{
		"/lushop.lushop.v1.Lushop/ListUsers": {}, // 管理员查看用户列表
		"/lushop.lushop.v1.Lushop/KickUser":  {}, // 管理员踢出用户
		// 可以添加更多管理员接口
	}

	_, ok := adminEndpoints[operation]
	return ok
}

// 权限常量定义
const (
	RoleAdmin = 1 // 管理员角色
	RoleUser  = 2 // 普通用户角色
)

// IsAdmin 检查当前用户是否为管理员
func IsAdmin(ctx context.Context) bool {
	role, ok := ctx.Value(userRoleKey).(int)
	return ok && role == RoleAdmin
}

// GetUserRole 获取当前用户角色
func GetUserRole(ctx context.Context) (int, bool) {
	role, ok := ctx.Value(userRoleKey).(int)
	return role, ok
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

// GetUserNickname 从上下文获取用户昵称
func GetUserNickname(ctx context.Context) (string, bool) {
	nickname, ok := ctx.Value(userNicknameKey).(string)
	return nickname, ok
}
