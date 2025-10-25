package service

import (
	"context"
	"testing"

	v1 "userop/api/userop/v1"
	"userop/internal/biz"
	"userop/internal/domain"
	"userop/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewFavoriteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAddrRepo := mocks.NewMockAddressRepo(ctrl)
	mockRepo := mocks.NewMockFavoriteRepo(ctrl)
	addrUc := biz.NewAddressUsecase(mockAddrRepo, nil)
	uc := biz.NewFavoriteUsecase(mockRepo, nil)
	svc := NewUserOpService(addrUc, uc, nil)

	assert.NotNil(t, svc)
}

func TestAddUserFav(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAddrRepo := mocks.NewMockAddressRepo(ctrl)
	mockRepo := mocks.NewMockFavoriteRepo(ctrl)
	addrUc := biz.NewAddressUsecase(mockAddrRepo, nil)
	uc := biz.NewFavoriteUsecase(mockRepo, nil)
	svc := NewUserOpService(addrUc, uc, nil)
	ctx := context.Background()

	mockRepo.EXPECT().
		AddUserFav(ctx, gomock.Any()).
		Return(nil)

	req := &v1.UserFavRequest{
		UserId:  1,
		GoodsId: 100,
	}

	reply, err := svc.AddUserFav(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestDeleteUserFav(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAddrRepo := mocks.NewMockAddressRepo(ctrl)
	mockRepo := mocks.NewMockFavoriteRepo(ctrl)
	addrUc := biz.NewAddressUsecase(mockAddrRepo, nil)
	uc := biz.NewFavoriteUsecase(mockRepo, nil)
	svc := NewUserOpService(addrUc, uc, nil)
	ctx := context.Background()

	mockRepo.EXPECT().
		DeleteUserFav(ctx, gomock.Any()).
		Return(nil)

	req := &v1.UserFavRequest{
		UserId:  1,
		GoodsId: 100,
	}

	reply, err := svc.DeleteUserFav(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestGetUserFavDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAddrRepo := mocks.NewMockAddressRepo(ctrl)
	mockRepo := mocks.NewMockFavoriteRepo(ctrl)
	addrUc := biz.NewAddressUsecase(mockAddrRepo, nil)
	uc := biz.NewFavoriteUsecase(mockRepo, nil)
	svc := NewUserOpService(addrUc, uc, nil)
	ctx := context.Background()

	expectedFav := &domain.Favorite{
		UserID:  1,
		GoodsID: 100,
	}

	mockRepo.EXPECT().
		GetUserFavDetail(ctx, gomock.Any()).
		Return(expectedFav, nil)

	req := &v1.UserFavRequest{
		UserId:  1,
		GoodsId: 100,
	}

	reply, err := svc.GetUserFavDetail(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestGetFavList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAddrRepo := mocks.NewMockAddressRepo(ctrl)
	mockRepo := mocks.NewMockFavoriteRepo(ctrl)
	addrUc := biz.NewAddressUsecase(mockAddrRepo, nil)
	uc := biz.NewFavoriteUsecase(mockRepo, nil)
	svc := NewUserOpService(addrUc, uc, nil)
	ctx := context.Background()

	fav1 := &domain.UserFavResponse{UserID: 1, GoodsID: 100}
	fav2 := &domain.UserFavResponse{UserID: 1, GoodsID: 200}
	expectedResp := &domain.UserFavListResponse{
		Total: 2,
		List:  []*domain.UserFavResponse{fav1, fav2},
	}

	mockRepo.EXPECT().
		GetFavList(ctx, gomock.Any()).
		Return(expectedResp, nil)

	req := &v1.UserFavRequest{UserId: 1}
	reply, err := svc.GetFavList(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(2), reply.Total)
}
