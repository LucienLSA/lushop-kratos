package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
// Wire 依赖注入 的核心部分,注册了如何创建某个对象的 “工具”（即构造函数）
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)
