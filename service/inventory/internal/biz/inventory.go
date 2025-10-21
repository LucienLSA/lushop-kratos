package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

// InventoryRepo is a Greater repo.
type InventoryRepo interface {
}

// InventoryUsecase is a Inventory usecase.
type InventoryUsecase struct {
	repo InventoryRepo
	log  *log.Helper
}

// NewInventoryUsecase new a Inventory usecase.
func NewInventoryUsecase(repo InventoryRepo, logger log.Logger) *InventoryUsecase {
	return &InventoryUsecase{repo: repo, log: log.NewHelper(logger)}
}
