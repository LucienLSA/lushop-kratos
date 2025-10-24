package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUsecase, NewUserAuthAdapter, NewCartUsecase, NewGoodsUsecase, NewInventoryUsecase, NewOrderUsecase, NewUserOpUsecase)
