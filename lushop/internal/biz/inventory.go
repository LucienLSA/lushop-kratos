package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// GoodsInvInfo 商品库存信息
type GoodsInvInfo struct {
	GoodsId int32
	Num     int32
}

// InventoryRepo 库存仓库接口
type InventoryRepo interface {
	SetInventory(ctx context.Context, goodsId, num int32) error
	GetInventory(ctx context.Context, goodsId int32) (int32, error)
	SellInventory(ctx context.Context, orderSn string, goodsInfo []*GoodsInvInfo) error
	RebackInventory(ctx context.Context, orderSn string, goodsInfo []*GoodsInvInfo) error
}

// InventoryUsecase 库存用例
type InventoryUsecase struct {
	repo InventoryRepo
	log  *log.Helper
}

// NewInventoryUsecase 创建库存用例
func NewInventoryUsecase(repo InventoryRepo, logger log.Logger) *InventoryUsecase {
	return &InventoryUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// SetInventory 设置库存
func (uc *InventoryUsecase) SetInventory(ctx context.Context, goodsId, num int32) error {
	return uc.repo.SetInventory(ctx, goodsId, num)
}

// GetInventory 获取库存
func (uc *InventoryUsecase) GetInventory(ctx context.Context, goodsId int32) (int32, error) {
	return uc.repo.GetInventory(ctx, goodsId)
}

// SellInventory 扣减库存
func (uc *InventoryUsecase) SellInventory(ctx context.Context, orderSn string, goodsInfo []*GoodsInvInfo) error {
	return uc.repo.SellInventory(ctx, orderSn, goodsInfo)
}

// RebackInventory 归还库存
func (uc *InventoryUsecase) RebackInventory(ctx context.Context, orderSn string, goodsInfo []*GoodsInvInfo) error {
	return uc.repo.RebackInventory(ctx, orderSn, goodsInfo)
}
