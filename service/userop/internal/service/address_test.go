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

func TestNewUserOpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAddrRepo := mocks.NewMockAddressRepo(ctrl)
	mockFavRepo := mocks.NewMockFavoriteRepo(ctrl)
	addrUc := biz.NewAddressUsecase(mockAddrRepo, nil)
	favUc := biz.NewFavoriteUsecase(mockFavRepo, nil)
	svc := NewUserOpService(addrUc, favUc, nil)

	assert.NotNil(t, svc)
}

func TestGetAddressList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAddressRepo(ctrl)
	mockFavRepo := mocks.NewMockFavoriteRepo(ctrl)
	uc := biz.NewAddressUsecase(mockRepo, nil)
	favUc := biz.NewFavoriteUsecase(mockFavRepo, nil)
	svc := NewUserOpService(uc, favUc, nil)
	ctx := context.Background()

	addr1 := &domain.Address{ID: 1, UserID: 1, Province: "广东省"}
	expectedResp := &domain.AddressListResponse{
		Total: 1,
		List:  []*domain.Address{addr1},
	}

	mockRepo.EXPECT().
		GetAddressList(ctx, gomock.Any()).
		Return(expectedResp, nil)

	req := &v1.AddressRequest{UserId: 1}
	reply, err := svc.GetAddressList(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(1), reply.Total)
}

func TestCreateAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAddressRepo(ctrl)
	mockFavRepo := mocks.NewMockFavoriteRepo(ctrl)
	uc := biz.NewAddressUsecase(mockRepo, nil)
	favUc := biz.NewFavoriteUsecase(mockFavRepo, nil)
	svc := NewUserOpService(uc, favUc, nil)
	ctx := context.Background()

	mockRepo.EXPECT().
		CreateAddress(ctx, gomock.Any()).
		Return(nil)

	req := &v1.AddressRequest{
		UserId:       1,
		Province:     "广东省",
		City:         "深圳市",
		District:     "南山区",
		Address:      "科技园",
		SignerName:   "张三",
		SignerMobile: "13800138000",
	}

	reply, err := svc.CreateAddress(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestUpdateAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAddressRepo(ctrl)
	mockFavRepo := mocks.NewMockFavoriteRepo(ctrl)
	uc := biz.NewAddressUsecase(mockRepo, nil)
	favUc := biz.NewFavoriteUsecase(mockFavRepo, nil)
	svc := NewUserOpService(uc, favUc, nil)
	ctx := context.Background()

	mockRepo.EXPECT().
		UpdateAddress(ctx, gomock.Any()).
		Return(nil)

	req := &v1.AddressRequest{
		Id:       1,
		UserId:   1,
		Province: "广东省",
	}

	reply, err := svc.UpdateAddress(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestDeleteAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAddressRepo(ctrl)
	mockFavRepo := mocks.NewMockFavoriteRepo(ctrl)
	uc := biz.NewAddressUsecase(mockRepo, nil)
	favUc := biz.NewFavoriteUsecase(mockFavRepo, nil)
	svc := NewUserOpService(uc, favUc, nil)
	ctx := context.Background()

	mockRepo.EXPECT().
		DeleteAddress(ctx, gomock.Any()).
		Return(nil)

	req := &v1.AddressRequest{
		Id:     1,
		UserId: 1,
	}

	reply, err := svc.DeleteAddress(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestMessageList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAddressRepo(ctrl)
	mockFavRepo := mocks.NewMockFavoriteRepo(ctrl)
	uc := biz.NewAddressUsecase(mockRepo, nil)
	favUc := biz.NewFavoriteUsecase(mockFavRepo, nil)
	svc := NewUserOpService(uc, favUc, nil)
	ctx := context.Background()

	msg1 := &domain.Message{ID: 1, UserID: 1, MessageType: 1}
	expectedResp := &domain.MessageListResponse{
		Total: 1,
		List:  []*domain.Message{msg1},
	}

	mockRepo.EXPECT().
		GetMessageList(ctx, gomock.Any()).
		Return(expectedResp, nil)

	req := &v1.MessageRequest{UserId: 1}
	reply, err := svc.MessageList(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, int32(1), reply.Total)
}

func TestCreateMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAddressRepo(ctrl)
	mockFavRepo := mocks.NewMockFavoriteRepo(ctrl)
	uc := biz.NewAddressUsecase(mockRepo, nil)
	favUc := biz.NewFavoriteUsecase(mockFavRepo, nil)
	svc := NewUserOpService(uc, favUc, nil)
	ctx := context.Background()

	mockRepo.EXPECT().
		CreateMessage(ctx, gomock.Any()).
		Return(nil)

	req := &v1.MessageRequest{
		UserId:      1,
		MessageType: 1,
		Subject:     "测试消息",
		Message:     "这是一条测试消息",
	}

	reply, err := svc.CreateMessage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}
