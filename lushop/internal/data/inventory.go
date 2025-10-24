package data

import (
	"context"
	inventoryV1 "lushop/api/service/inventory/v1"
	"lushop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type inventoryRepo struct {
	data *Data
	log  *log.Helper
}

// NewInventoryRepo 创建库存仓库
func NewInventoryRepo(data *Data, logger log.Logger) biz.InventoryRepo {
	return &inventoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// SetInventory 设置库存
func (r *inventoryRepo) SetInventory(ctx context.Context, goodsId, num int32) error {
	r.log.Infof("调用 Inventory gRPC: SetInv goodsId=%d, num=%d", goodsId, num)

	_, err := r.data.ic.SetInv(ctx, &inventoryV1.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     num,
	})
	if err != nil {
		r.log.Errorf("设置库存失败: goodsId=%d, error=%v", goodsId, err)
		return err
	}

	r.log.Infof("设置库存成功: goodsId=%d, num=%d", goodsId, num)
	return nil
}

// GetInventory 获取库存信息
func (r *inventoryRepo) GetInventory(ctx context.Context, goodsId int32) (int32, error) {
	r.log.Infof("调用 Inventory gRPC: InvDetail goodsId=%d", goodsId)

	resp, err := r.data.ic.InvDetail(ctx, &inventoryV1.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		r.log.Errorf("获取库存失败: goodsId=%d, error=%v", goodsId, err)
		return 0, err
	}

	r.log.Infof("获取库存成功: goodsId=%d, num=%d", goodsId, resp.Num)
	return resp.Num, nil
}

// SellInventory 扣减库存
func (r *inventoryRepo) SellInventory(ctx context.Context, orderSn string, goodsInfo []*biz.GoodsInvInfo) error {
	r.log.Infof("调用 Inventory gRPC: Sell orderSn=%s, goods count=%d", orderSn, len(goodsInfo))

	// 转换为 gRPC 请求
	invList := make([]*inventoryV1.GoodsInvInfo, 0, len(goodsInfo))
	for _, item := range goodsInfo {
		invList = append(invList, &inventoryV1.GoodsInvInfo{
			GoodsId: item.GoodsId,
			Num:     item.Num,
		})
	}

	_, err := r.data.ic.Sell(ctx, &inventoryV1.SellInfo{
		OrderSn:   orderSn,
		GoodsInfo: invList,
	})
	if err != nil {
		r.log.Errorf("扣减库存失败: orderSn=%s, error=%v", orderSn, err)
		return err
	}

	r.log.Infof("扣减库存成功: orderSn=%s", orderSn)
	return nil
}

// RebackInventory 归还库存
func (r *inventoryRepo) RebackInventory(ctx context.Context, orderSn string, goodsInfo []*biz.GoodsInvInfo) error {
	r.log.Infof("调用 Inventory gRPC: Reback orderSn=%s, goods count=%d", orderSn, len(goodsInfo))

	// 转换为 gRPC 请求
	invList := make([]*inventoryV1.GoodsInvInfo, 0, len(goodsInfo))
	for _, item := range goodsInfo {
		invList = append(invList, &inventoryV1.GoodsInvInfo{
			GoodsId: item.GoodsId,
			Num:     item.Num,
		})
	}

	_, err := r.data.ic.Reback(ctx, &inventoryV1.SellInfo{
		OrderSn:   orderSn,
		GoodsInfo: invList,
	})
	if err != nil {
		r.log.Errorf("归还库存失败: orderSn=%s, error=%v", orderSn, err)
		return err
	}

	r.log.Infof("归还库存成功: orderSn=%s", orderSn)
	return nil
}
