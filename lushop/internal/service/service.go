package service

import (
	v1 "lushop/api/lushop/v1"
	"lushop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewLushopService)

// 通过将 NewLushopService 加入 ProviderSet，当需要创建一个 LushopService 实例时，请使用这个 NewLushopService 函数。
// 在项目的 main.go 或初始化部分，所有的 ProviderSet 会被收集起来，Wire 会自动分析依赖关系并生成最终的初始化代码
// LushopService is a lushop service.
type LushopService struct {
	v1.UnimplementedLushopServer
	uc  *biz.UserUsecase
	log *log.Helper
}

// NewLushopService new a shop service.
// 遵循了 Wire 要求的依赖注入规范
// gRPC 服务器启动时，Kratos 的依赖注入系统会调用 NewLushopService
// 自动创建好 LushopService 实例，并把它注册到 gRPC 服务器上，外部就可以通过 gRPC 调用定义的方法
func NewLushopService(uc *biz.UserUsecase, logger log.Logger) *LushopService {
	return &LushopService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/lushop")),
	}
}
