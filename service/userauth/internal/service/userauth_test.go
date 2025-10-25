package service

import (
	"context"
	"errors"
	"testing"
	"time"

	v1 "userauth/api/userauth/v1"
	"userauth/internal/biz"
	"userauth/internal/conf"
	"userauth/internal/mocks"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

func TestNewUserAuthService(t *testing.T) {
	logger := log.NewStdLogger(io.Discard)
	authConf := &conf.Auth{JwtKey: "test-key-for-testing-only-must-be-long"}
	smsConf := &conf.Sms{}
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := mocks.NewMockAuthRepo(ctrl)
	uc := biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
	svc := NewUserAuthService(uc)
	
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.uc)
}

func TestVerifyCaptcha(t *testing.T) {
	logger := log.NewStdLogger(io.Discard)
	authConf := &conf.Auth{JwtKey: "test-key"}
	smsConf := &conf.Sms{}
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := mocks.NewMockAuthRepo(ctrl)
	uc := biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
	svc := NewUserAuthService(uc)
	ctx := context.Background()
	
	// Mock 期望
	mockRepo.EXPECT().
		GetCaptcha(ctx, "test-id").
		Return("1234", nil)
	
	req := &v1.VerifyCaptchaReq{
		CaptchaId: "test-id",
		Captcha:   "1234",
	}
	
	reply, err := svc.VerifyCaptcha(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.Success)
}

func TestVerifySms(t *testing.T) {
	logger := log.NewStdLogger(io.Discard)
	authConf := &conf.Auth{JwtKey: "test-key"}
	smsConf := &conf.Sms{}
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := mocks.NewMockAuthRepo(ctrl)
	uc := biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
	svc := NewUserAuthService(uc)
	ctx := context.Background()
	
	// Mock 期望
	mockRepo.EXPECT().
		GetSmsCode(ctx, "13800138000").
		Return("123456", nil)
	
	req := &v1.VerifySmsReq{
		Mobile:  "13800138000",
		SmsCode: "123456",
	}
	
	reply, err := svc.VerifySms(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.Success)
}

func TestIssueToken(t *testing.T) {
	logger := log.NewStdLogger(io.Discard)
	authConf := &conf.Auth{JwtKey: "test-jwt-key-for-testing-only-must-be-long-enough"}
	smsConf := &conf.Sms{}
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := mocks.NewMockAuthRepo(ctrl)
	uc := biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
	svc := NewUserAuthService(uc)
	ctx := context.Background()
	
	// Mock 期望
	mockRepo.EXPECT().
		StoreAccessToken(ctx, int64(1), gomock.Any(), 30*time.Minute).
		Return(nil)
	
	mockRepo.EXPECT().
		StoreRefreshToken(ctx, int64(1), gomock.Any(), 7*24*time.Hour).
		Return(nil)
	
	req := &v1.IssueTokenReq{
		UserId:   1,
		NickName: "testuser",
		Role:     1,
	}
	
	reply, err := svc.IssueToken(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.NotEmpty(t, reply.AccessToken)
	assert.NotEmpty(t, reply.RefreshToken)
	assert.Greater(t, reply.ExpiredAt, time.Now().Unix())
}

func TestRevokeToken(t *testing.T) {
	logger := log.NewStdLogger(io.Discard)
	authConf := &conf.Auth{JwtKey: "test-key"}
	smsConf := &conf.Sms{}
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := mocks.NewMockAuthRepo(ctrl)
	uc := biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
	svc := NewUserAuthService(uc)
	ctx := context.Background()
	
	// Mock 期望
	mockRepo.EXPECT().
		DeleteTokens(ctx, int64(1)).
		Return(nil)
	
	mockRepo.EXPECT().
		AddToBlacklist(ctx, int64(1), 30*time.Minute).
		Return(nil)
	
	req := &v1.RevokeTokenReq{
		UserId: 1,
	}
	
	reply, err := svc.RevokeToken(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.IsType(t, &emptypb.Empty{}, reply)
}

func TestCheckBlacklist(t *testing.T) {
	logger := log.NewStdLogger(io.Discard)
	authConf := &conf.Auth{JwtKey: "test-key"}
	smsConf := &conf.Sms{}
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := mocks.NewMockAuthRepo(ctrl)
	uc := biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
	svc := NewUserAuthService(uc)
	ctx := context.Background()
	
	// Mock 期望
	mockRepo.EXPECT().
		CheckBlacklist(ctx, int64(1)).
		Return(true, nil)
	
	req := &v1.CheckBlacklistReq{
		UserId: 1,
	}
	
	reply, err := svc.CheckBlacklist(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.IsBlacklisted)
}

func TestAddToBlacklist(t *testing.T) {
	logger := log.NewStdLogger(io.Discard)
	authConf := &conf.Auth{JwtKey: "test-key"}
	smsConf := &conf.Sms{}
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := mocks.NewMockAuthRepo(ctrl)
	uc := biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
	svc := NewUserAuthService(uc)
	ctx := context.Background()
	
	// Mock 期望
	mockRepo.EXPECT().
		AddToBlacklist(ctx, int64(1), time.Duration(3600)*time.Second).
		Return(nil)
	
	req := &v1.AddToBlacklistReq{
		UserId:     1,
		TtlSeconds: 3600,
	}
	
	reply, err := svc.AddToBlacklist(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.IsType(t, &emptypb.Empty{}, reply)
}

// 测试错误场景
func TestVerifyCaptchaError(t *testing.T) {
	logger := log.NewStdLogger(io.Discard)
	authConf := &conf.Auth{JwtKey: "test-key"}
	smsConf := &conf.Sms{}
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := mocks.NewMockAuthRepo(ctrl)
	uc := biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
	svc := NewUserAuthService(uc)
	ctx := context.Background()
	
	// Mock 返回错误
	mockRepo.EXPECT().
		GetCaptcha(ctx, "test-id").
		Return("", errors.New("redis error"))
	
	req := &v1.VerifyCaptchaReq{
		CaptchaId: "test-id",
		Captcha:   "1234",
	}
	
	reply, err := svc.VerifyCaptcha(ctx, req)
	
	assert.NoError(t, err) // VerifyCaptcha 不返回错误，只返回 false
	assert.NotNil(t, reply)
	assert.False(t, reply.Success)
}
