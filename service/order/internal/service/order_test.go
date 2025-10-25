package service

import (
	"context"
	"testing"

	v1 "order/api/order/v1"
	"order/internal/biz"
	"order/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"

	"github.com/go-kratos/kratos/v2/log"
)

func TestNewOrderService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewOrderUsecase(mockRepo, logger)
	svc := NewOrderService(uc, logger)

	assert.NotNil(t, svc)
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewOrderUsecase(mockRepo, logger)
	svc := NewOrderService(uc, logger)
	ctx := context.Background()

	req := &v1.OrderRequest{
		UserId:  1,
		Address: "广东省深圳市南山区",
		Name:    "张三",
		Mobile:  "13800138000",
	}

	expectedOrder := &v1.OrderInfoResponse{
		Id:      1,
		OrderSn: "ORDER001",
		UserId:  1,
	}

	mockRepo.EXPECT().
		CreateOrder(ctx, req).
		Return(expectedOrder, nil)

	reply, err := svc.CreateOrder(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(1), reply.Id)
}

func TestGetOrderList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewOrderUsecase(mockRepo, logger)
	svc := NewOrderService(uc, logger)
	ctx := context.Background()

	req := &v1.OrderFilterRequest{
		UserId: 1,
		Pages:  1,
	}

	expectedOrders := []*v1.OrderInfoResponse{
		{Id: 1, OrderSn: "ORDER001", UserId: 1},
	}

	mockRepo.EXPECT().
		GetOrderList(ctx, req).
		Return(expectedOrders, int32(1), nil)

	reply, err := svc.OrderList(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(1), reply.Total)
}

func TestGetOrderDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewOrderUsecase(mockRepo, logger)
	svc := NewOrderService(uc, logger)
	ctx := context.Background()

	req := &v1.OrderRequest{
		Id:     1,
		UserId: 1,
	}

	expectedDetail := &v1.OrderInfoDetailResponse{
		OrderInfo: &v1.OrderInfoResponse{
			Id:      1,
			OrderSn: "ORDER001",
			UserId:  1,
		},
	}

	mockRepo.EXPECT().
		GetOrderDetail(ctx, req).
		Return(expectedDetail, nil)

	reply, err := svc.OrderDetail(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(1), reply.OrderInfo.Id)
}
