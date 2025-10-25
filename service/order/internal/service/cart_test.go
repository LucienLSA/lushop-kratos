package service

import (
	"context"
	"testing"

	v1 "order/api/order/v1"
	"order/internal/biz"
	"order/internal/mocks"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

func TestCartList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewOrderUsecase(mockRepo, logger)
	svc := NewOrderService(uc, logger)
	ctx := context.Background()

	req := &v1.UserInfo{Id: 1}

	expectedCarts := []*v1.ShopCartInfoResponse{
		{Id: 1, UserId: 1, GoodsId: 100, Nums: 2},
		{Id: 2, UserId: 1, GoodsId: 200, Nums: 1},
	}

	mockRepo.EXPECT().
		GetCartList(ctx, int32(1)).
		Return(expectedCarts, nil)

	reply, err := svc.CartItemList(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(2), reply.Total)
	assert.Equal(t, 2, len(reply.Data))
}

func TestCreateCartItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewOrderUsecase(mockRepo, logger)
	svc := NewOrderService(uc, logger)
	ctx := context.Background()

	req := &v1.CartItemRequest{
		UserId:  1,
		GoodsId: 100,
		Nums:    2,
	}

	expectedCart := &v1.ShopCartInfoResponse{
		Id:      1,
		UserId:  1,
		GoodsId: 100,
		Nums:    2,
	}

	mockRepo.EXPECT().
		CreateCartItem(ctx, req).
		Return(expectedCart, nil)

	reply, err := svc.CreateCartItem(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(1), reply.Id)
}

func TestUpdateCartItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewOrderUsecase(mockRepo, logger)
	svc := NewOrderService(uc, logger)
	ctx := context.Background()

	req := &v1.CartItemRequest{
		Id:     1,
		UserId: 1,
		Nums:   3,
	}

	mockRepo.EXPECT().
		UpdateCartItem(ctx, req).
		Return(&emptypb.Empty{}, nil)

	reply, err := svc.UpdateCartItem(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestDeleteCartItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewOrderUsecase(mockRepo, logger)
	svc := NewOrderService(uc, logger)
	ctx := context.Background()

	req := &v1.CartItemRequest{
		Id:     1,
		UserId: 1,
	}

	mockRepo.EXPECT().
		DeleteCartItem(ctx, req).
		Return(&emptypb.Empty{}, nil)

	reply, err := svc.DeleteCartItem(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}
