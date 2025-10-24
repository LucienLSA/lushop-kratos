package client

import (
	"order/internal/client/goods"

	"github.com/google/wire"
)

// ProviderSet is client providers.
// 提供所有外部服务的客户端
var ProviderSet = wire.NewSet(
	goods.NewGoodsServiceClient,
	// 未来可以添加其他服务客户端
	// inventory.NewInventoryServiceClient,
	// user.NewUserServiceClient,
)
