package service

import (
	"context"
	v1 "lushop/api/lushop/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// SetInventory 设置库存 HTTP API
func (s *LushopService) SetInventory(ctx context.Context, req *v1.SetInventoryReq) (*emptypb.Empty, error) {
	s.log.Infof("HTTP API: 设置库存 goodsId=%d, num=%d", req.GoodsId, req.Num)

	// 调用业务逻辑（内部会调用 gRPC）
	err := s.inventoryUc.SetInventory(ctx, req.GoodsId, req.Num)
	if err != nil {
		s.log.Errorf("设置库存失败: goodsId=%d, error=%v", req.GoodsId, err)
		return nil, err
	}

	s.log.Infof("设置库存成功: goodsId=%d, num=%d", req.GoodsId, req.Num)
	return &emptypb.Empty{}, nil
}

// GetInventory 获取库存 HTTP API
func (s *LushopService) GetInventory(ctx context.Context, req *v1.GetInventoryReq) (*v1.GetInventoryReply, error) {
	s.log.Infof("HTTP API: 获取库存 goodsId=%d", req.GoodsId)

	// 调用业务逻辑（内部会调用 gRPC）
	num, err := s.inventoryUc.GetInventory(ctx, req.GoodsId)
	if err != nil {
		s.log.Errorf("获取库存失败: goodsId=%d, error=%v", req.GoodsId, err)
		return nil, err
	}

	s.log.Infof("获取库存成功: goodsId=%d, num=%d", req.GoodsId, num)
	return &v1.GetInventoryReply{
		GoodsId: req.GoodsId,
		Num:     num,
	}, nil
}
