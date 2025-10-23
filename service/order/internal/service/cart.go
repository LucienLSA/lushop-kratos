package service

import (
	"context"
	v1 "order/api/order/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *OrderService) CartItemList(ctx context.Context, in *v1.UserInfo) (*v1.CartItemListResponse, error) {
	// 调用业务层获取购物车列表
	cartList, err := s.uc.GetCartList(ctx, in.Id)
	if err != nil {
		s.log.Errorf("failed to get cart list for user %d: %v", in.Id, err)
		return nil, err
	}

	// 构造响应
	return &v1.CartItemListResponse{
		Total: int32(len(cartList)),
		Data:  cartList,
	}, nil
}

func (s *OrderService) CreateCartItem(ctx context.Context, in *v1.CartItemRequest) (*v1.ShopCartInfoResponse, error) {
	// TODO: 实现添加购物车逻辑
	return &v1.ShopCartInfoResponse{
		Id:      in.Id,
		UserId:  in.UserId,
		GoodsId: in.GoodsId,
		Nums:    in.Nums,
		Checked: in.Checked,
	}, nil
}

func (s *OrderService) UpdateCartItem(ctx context.Context, in *v1.CartItemRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *OrderService) DeleteCartItem(ctx context.Context, in *v1.CartItemRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
