package service

import (
	"context"
	"errors"
	"io"
	"testing"

	v1 "inventory/api/inventory/v1"
	"inventory/internal/biz"
	"inventory/internal/domain"
	"inventory/internal/mocks"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewInventoryService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)

	assert.NotNil(t, svc)
}

func TestSetInv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)
	ctx := context.Background()

	req := &v1.GoodsInvInfo{
		GoodsId: 100,
		Num:     50,
	}

	mockRepo.EXPECT().
		AddInv(ctx, &domain.Inventory{
			Goods:  100,
			Stocks: 50,
		}).
		Return(nil)

	reply, err := svc.SetInv(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestSetInvError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)
	ctx := context.Background()

	req := &v1.GoodsInvInfo{
		GoodsId: 100,
		Num:     50,
	}

	mockRepo.EXPECT().
		AddInv(ctx, gomock.Any()).
		Return(errors.New("database error"))

	reply, err := svc.SetInv(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestInvDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)
	ctx := context.Background()

	req := &v1.GoodsInvInfo{
		GoodsId: 100,
	}

	expectedInv := &domain.Inventory{
		Goods:  100,
		Stocks: 50,
	}

	mockRepo.EXPECT().
		GetInvById(ctx, int32(100)).
		Return(expectedInv, nil)

	reply, err := svc.InvDetail(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(100), reply.GoodsId)
	assert.Equal(t, int32(50), reply.Num)
}

func TestInvDetailNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)
	ctx := context.Background()

	req := &v1.GoodsInvInfo{
		GoodsId: 999,
	}

	mockRepo.EXPECT().
		GetInvById(ctx, int32(999)).
		Return(nil, errors.New("record not found"))

	reply, err := svc.InvDetail(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestSell(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)
	ctx := context.Background()

	req := &v1.SellInfo{
		OrderSn: "ORDER001",
		GoodsInfo: []*v1.GoodsInvInfo{
			{GoodsId: 100, Num: 2},
			{GoodsId: 200, Num: 3},
		},
	}

	mockRepo.EXPECT().
		Sell(ctx, gomock.Any()).
		Return(nil)

	reply, err := svc.Sell(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestSellStockNotEnough(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)
	ctx := context.Background()

	req := &v1.SellInfo{
		OrderSn: "ORDER002",
		GoodsInfo: []*v1.GoodsInvInfo{
			{GoodsId: 100, Num: 1000},
		},
	}

	mockRepo.EXPECT().
		Sell(ctx, gomock.Any()).
		Return(errors.New("not enough stock"))

	reply, err := svc.Sell(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestReback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)
	ctx := context.Background()

	req := &v1.SellInfo{
		OrderSn: "ORDER001",
	}

	mockRepo.EXPECT().
		Reback(ctx, "ORDER001").
		Return(nil)

	reply, err := svc.Reback(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestRebackError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockInventoryRepo(ctrl)
	logger := log.NewStdLogger(io.Discard)
	uc := biz.NewInventoryUsecase(mockRepo, logger)
	svc := NewInventoryService(uc)
	ctx := context.Background()

	req := &v1.SellInfo{
		OrderSn: "ORDER002",
	}

	mockRepo.EXPECT().
		Reback(ctx, "ORDER002").
		Return(errors.New("database error"))

	reply, err := svc.Reback(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}
