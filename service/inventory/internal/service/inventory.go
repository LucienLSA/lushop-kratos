package service

import (
	"context"
	v1 "inventory/api/inventory/v1"
	"inventory/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *InventoryService) SetInv(ctx context.Context, req *v1.GoodsInvInfo) (*emptypb.Empty, error) {
	goodsInv := &domain.Inventory{
		Goods:  req.GoodsId,
		Stocks: req.Num,
	}
	err := s.uc.SetInv(ctx, goodsInv)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *InventoryService) InvDetail(ctx context.Context, req *v1.GoodsInvInfo) (*v1.GoodsInvInfo, error) {
	inv, err := s.uc.GetInvById(ctx, req.GoodsId)
	if err != nil {
		return nil, err
	}
	return &v1.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (s *InventoryService) Sell(ctx context.Context, req *v1.SellInfo) (*emptypb.Empty, error) {
	items := make([]domain.GoodsInvInfo, 0, len(req.GoodsInfo))
	for _, gi := range req.GoodsInfo {
		items = append(items, domain.GoodsInvInfo{GoodsID: gi.GoodsId, Nums: gi.Num})
	}
	sell := &domain.SellInfo{
		GoodsInvInfo: items,
		OrderSn:      req.OrderSn,
	}
	if err := s.uc.Sell(ctx, sell); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Reback 库存归还
// 通过订单号归还库存，支持幂等性
// 适用场景：订单取消、超时未支付等需要归还库存的场景
func (s *InventoryService) Reback(ctx context.Context, req *v1.SellInfo) (*emptypb.Empty, error) {
	if err := s.uc.Reback(ctx, req.OrderSn); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
