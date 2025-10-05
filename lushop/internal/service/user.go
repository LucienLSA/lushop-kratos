package service

import (
	"context"
	v1 "lushop/api/lushop/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *LushopService) Register(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterReply, error) {
	return s.uc.CreateUser(ctx, req)
}

func (s *LushopService) Login(ctx context.Context, req *v1.LoginReq) (*v1.RegisterReply, error) {
	return s.uc.PasswordLogin(ctx, req)
}

func (s *LushopService) Captcha(ctx context.Context, req *emptypb.Empty) (*v1.CaptchaReply, error) {
	return s.uc.GetCaptcha(ctx)
}
func (s *LushopService) Detail(ctx context.Context, req *emptypb.Empty) (*v1.UserDetailResponse, error) {
	return s.uc.UserDetailByID(ctx)
}
