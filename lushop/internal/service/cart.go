package service

import (
	"context"

	v1 "lushop/api/lushop/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// GetCartList 获取购物车列表
func (s *LushopService) GetCartList(ctx context.Context, req *emptypb.Empty) (*v1.CartListReply, error) {
	return s.cartUc.GetCartList(ctx)
}

// AddToCart 添加商品到购物车
func (s *LushopService) AddToCart(ctx context.Context, req *v1.AddToCartReq) (*v1.CartItemReply, error) {
	return s.cartUc.AddToCart(ctx, req)
}

// UpdateCartItem 更新购物车商品
func (s *LushopService) UpdateCartItem(ctx context.Context, req *v1.UpdateCartItemReq) (*emptypb.Empty, error) {
	err := s.cartUc.UpdateCartItem(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteCartItem 删除购物车商品
func (s *LushopService) DeleteCartItem(ctx context.Context, req *v1.DeleteCartItemReq) (*emptypb.Empty, error) {
	err := s.cartUc.DeleteCartItem(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

