package biz

import (
	"context"
	v1 "order/api/order/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateOrder 创建订单
func (u *OrderUsecase) CreateOrder(ctx context.Context, req *v1.OrderRequest) (*v1.OrderInfoResponse, error) {
	// 业务层参数验证
	if req.UserId == 0 {
		u.log.Error("user id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "user id is required")
	}
	if req.Address == "" {
		u.log.Error("address is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "address is required")
	}
	if req.Name == "" {
		u.log.Error("name is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "name is required")
	}
	if req.Mobile == "" {
		u.log.Error("mobile is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "mobile is required")
	}

	// 调用数据层创建订单
	order, err := u.repo.CreateOrder(ctx, req)
	if err != nil {
		u.log.Errorf("failed to create order: user=%d, error=%v", req.UserId, err)
		return nil, err
	}

	u.log.Infof("order created successfully: order_id=%d, order_sn=%s, user=%d", order.Id, order.OrderSn, req.UserId)
	return order, nil
}

// GetOrderList 获取订单列表
func (u *OrderUsecase) GetOrderList(ctx context.Context, req *v1.OrderFilterRequest) ([]*v1.OrderInfoResponse, int32, error) {
	// 业务层参数验证
	if req.UserId == 0 {
		u.log.Error("user id is required")
		return nil, 0, errors.BadRequest("INVALID_PARAMETER", "user id is required")
	}
	if req.Pages <= 0 {
		req.Pages = 1
	}
	if req.PagePerNums <= 0 {
		req.PagePerNums = 10
	}

	// 调用数据层获取订单列表
	orders, total, err := u.repo.GetOrderList(ctx, req)
	if err != nil {
		u.log.Errorf("failed to get order list: user=%d, error=%v", req.UserId, err)
		return nil, 0, err
	}

	u.log.Infof("order list retrieved: user=%d, total=%d, page=%d", req.UserId, total, req.Pages)
	return orders, total, nil
}

// GetOrderDetail 获取订单详情
func (u *OrderUsecase) GetOrderDetail(ctx context.Context, req *v1.OrderRequest) (*v1.OrderInfoDetailResponse, error) {
	// 业务层参数验证
	if req.Id == 0 {
		u.log.Error("order id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "order id is required")
	}
	if req.UserId == 0 {
		u.log.Error("user id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "user id is required")
	}

	// 调用数据层获取订单详情
	detail, err := u.repo.GetOrderDetail(ctx, req)
	if err != nil {
		u.log.Errorf("failed to get order detail: order_id=%d, user=%d, error=%v", req.Id, req.UserId, err)
		return nil, err
	}

	u.log.Infof("order detail retrieved: order_id=%d, user=%d", req.Id, req.UserId)
	return detail, nil
}

// UpdateOrderStatus 更新订单状态
func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, req *v1.OrderStatus) (*emptypb.Empty, error) {
	// 业务层参数验证
	if req.OrderSn == "" {
		u.log.Error("order sn is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "order sn is required")
	}
	if req.Status == "" {
		u.log.Error("status is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "status is required")
	}

	// 验证状态值是否合法
	validStatuses := map[string]bool{
		"WAIT_BUYER_PAY": true,
		"PAYING":         true,
		"TRADE_SUCCESS":  true,
		"TRADE_CLOSED":   true,
		"TRADE_FINISHED": true,
	}
	if !validStatuses[req.Status] {
		u.log.Errorf("invalid status: %s", req.Status)
		return nil, errors.BadRequest("INVALID_STATUS", "invalid order status")
	}

	// 调用数据层更新订单状态
	_, err := u.repo.UpdateOrderStatus(ctx, req)
	if err != nil {
		u.log.Errorf("failed to update order status: order_sn=%s, status=%s, error=%v", req.OrderSn, req.Status, err)
		return nil, err
	}

	u.log.Infof("order status updated: order_sn=%s, status=%s", req.OrderSn, req.Status)
	return &emptypb.Empty{}, nil
}

// HandleOrderTimeout 处理订单超时
// 1. 查询订单状态
// 2. 如果未支付，发送归还库存消息
// 3. 更新订单状态为 TRADE_CLOSED
// 幂等性：通过订单状态判断，已支付或已关闭的订单不处理
func (u *OrderUsecase) HandleOrderTimeout(ctx context.Context, orderSn string) error {
	if orderSn == "" {
		u.log.Error("order sn is required")
		return errors.BadRequest("INVALID_PARAMETER", "order sn is required")
	}

	// 调用数据层处理订单超时
	if err := u.repo.HandleOrderTimeout(ctx, orderSn); err != nil {
		u.log.Errorf("failed to handle order timeout: order_sn=%s, error=%v", orderSn, err)
		return err
	}

	u.log.Infof("order timeout handled: order_sn=%s", orderSn)
	return nil
}
