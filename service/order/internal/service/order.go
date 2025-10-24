package service

import (
	"context"
	v1 "order/api/order/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, in *v1.OrderRequest) (*v1.OrderInfoResponse, error) {
	return s.uc.CreateOrder(ctx, in)
}

// OrderList 获取订单列表
func (s *OrderService) OrderList(ctx context.Context, in *v1.OrderFilterRequest) (*v1.OrderListResponse, error) {
	orders, total, err := s.uc.GetOrderList(ctx, in)
	if err != nil {
		return nil, err
	}
	return &v1.OrderListResponse{
		Total: total,
		Data:  orders,
	}, nil
}

// OrderDetail 获取订单详情
func (s *OrderService) OrderDetail(ctx context.Context, in *v1.OrderRequest) (*v1.OrderInfoDetailResponse, error) {
	return s.uc.GetOrderDetail(ctx, in)
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(ctx context.Context, in *v1.OrderStatus) (*emptypb.Empty, error) {
	return s.uc.UpdateOrderStatus(ctx, in)
}
