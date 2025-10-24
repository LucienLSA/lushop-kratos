package biz

import (
	"context"
	"inventory/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

// InventoryRepo is a Inventory repo.
type InventoryRepo interface {
	AddInv(ctx context.Context, inv *domain.Inventory) error
	GetInvById(ctx context.Context, goodsId int32) (*domain.Inventory, error)
	Sell(ctx context.Context, sell *domain.SellInfo) error
	Reback(ctx context.Context, orderSn string) error
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

func (uc *InventoryUsecase) SetInv(ctx context.Context, inv *domain.Inventory) error {
	return uc.repo.AddInv(ctx, inv)
}

func (uc *InventoryUsecase) GetInvById(ctx context.Context, goodsId int32) (*domain.Inventory, error) {
	return uc.repo.GetInvById(ctx, goodsId)
}

func (uc *InventoryUsecase) Sell(ctx context.Context, sell *domain.SellInfo) error {
	return uc.repo.Sell(ctx, sell)
}

// Reback 库存归还
// 通过订单号归还库存，确保幂等性
func (uc *InventoryUsecase) Reback(ctx context.Context, orderSn string) error {
	return uc.repo.Reback(ctx, orderSn)
}
