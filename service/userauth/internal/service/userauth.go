package service

import (
	"context"

	v1 "userauth/api/userauth/v1"
	"userauth/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserAuthService struct {
	v1.UnimplementedUserAuthServer

	uc *biz.AuthUsecase
}

func NewUserAuthService(uc *biz.AuthUsecase) *UserAuthService {
	return &UserAuthService{uc: uc}
}

func (s *UserAuthService) GetCaptcha(ctx context.Context, req *emptypb.Empty) (*v1.CaptchaReply, error) {
	return s.uc.GetCaptcha(ctx)
}

func (s *UserAuthService) VerifyCaptcha(ctx context.Context, req *v1.VerifyCaptchaReq) (*v1.VerifyCaptchaReply, error) {
	success, err := s.uc.VerifyCaptcha(ctx, req.CaptchaId, req.Captcha)
	if err != nil {
		return nil, err
	}
	return &v1.VerifyCaptchaReply{Success: success}, nil
}

func (s *UserAuthService) SendSms(ctx context.Context, req *v1.SendSmsReq) (*v1.SendSmsReply, error) {
	success, err := s.uc.SendSms(ctx, req.Mobile)
	if err != nil {
		return nil, err
	}
	return &v1.SendSmsReply{Success: success}, nil
}

func (s *UserAuthService) VerifySms(ctx context.Context, req *v1.VerifySmsReq) (*v1.VerifySmsReply, error) {
	success, err := s.uc.VerifySms(ctx, req.Mobile, req.SmsCode)
	if err != nil {
		return nil, err
	}
	return &v1.VerifySmsReply{Success: success}, nil
}

func (s *UserAuthService) IssueToken(ctx context.Context, req *v1.IssueTokenReq) (*v1.TokenReply, error) {
	return s.uc.IssueToken(ctx, req.UserId, req.NickName, req.Role)
}

func (s *UserAuthService) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.TokenReply, error) {
	return s.uc.RefreshToken(ctx, req.RefreshToken)
}

func (s *UserAuthService) RevokeToken(ctx context.Context, req *v1.RevokeTokenReq) (*emptypb.Empty, error) {
	err := s.uc.RevokeToken(ctx, req.UserId)
	return &emptypb.Empty{}, err
}

func (s *UserAuthService) AddToBlacklist(ctx context.Context, req *v1.AddToBlacklistReq) (*emptypb.Empty, error) {
	err := s.uc.AddToBlacklist(ctx, req.UserId, req.TtlSeconds)
	return &emptypb.Empty{}, err
}

func (s *UserAuthService) CheckBlacklist(ctx context.Context, req *v1.CheckBlacklistReq) (*v1.CheckBlacklistReply, error) {
	isBlacklisted, err := s.uc.CheckBlacklist(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.CheckBlacklistReply{IsBlacklisted: isBlacklisted}, nil
}
