package service

import (
	"context"
	v1 "lushop/api/lushop/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateOrder 创建订单 HTTP API
func (s *LushopService) CreateOrder(ctx context.Context, req *v1.CreateOrderReq) (*v1.CreateOrderReply, error) {
	s.log.Infof("HTTP API: 创建订单 address=%s, name=%s", req.Address, req.Name)
	return s.orderUc.CreateOrder(ctx, req)
}

// GetOrderList 获取订单列表 HTTP API
func (s *LushopService) GetOrderList(ctx context.Context, req *v1.GetOrderListReq) (*v1.GetOrderListReply, error) {
	s.log.Infof("HTTP API: 获取订单列表 page=%d, pageSize=%d", req.Page, req.PageSize)
	return s.orderUc.GetOrderList(ctx, req)
}

// GetOrderDetail 获取订单详情 HTTP API
func (s *LushopService) GetOrderDetail(ctx context.Context, req *v1.GetOrderDetailReq) (*v1.GetOrderDetailReply, error) {
	s.log.Infof("HTTP API: 获取订单详情 id=%d", req.Id)
	return s.orderUc.GetOrderDetail(ctx, req)
}

// CancelOrder 取消订单 HTTP API
func (s *LushopService) CancelOrder(ctx context.Context, req *v1.CancelOrderReq) (*emptypb.Empty, error) {
	s.log.Infof("HTTP API: 取消订单 id=%d", req.Id)
	err := s.orderUc.CancelOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
