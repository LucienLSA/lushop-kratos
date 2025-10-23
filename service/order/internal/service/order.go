package service

import (
	"context"
	v1 "order/api/order/v1"
)

func (s *OrderService) CreateOrder(ctx context.Context, in *v1.OrderRequest) (*v1.OrderInfoResponse, error) {
	return &v1.OrderInfoResponse{}, nil
}
