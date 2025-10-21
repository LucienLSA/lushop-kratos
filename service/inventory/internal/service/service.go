package service

import (
	v1 "inventory/api/inventory/v1"
	"inventory/internal/biz"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewInventoryService)

// InventoryService is a inventory service.
type InventoryService struct {
	v1.UnimplementedInventoryServer
	uc *biz.InventoryUsecase
}

// NewInventoryService new a inventory service.
func NewInventoryService(uc *biz.InventoryUsecase) *InventoryService {
	return &InventoryService{uc: uc}
}
