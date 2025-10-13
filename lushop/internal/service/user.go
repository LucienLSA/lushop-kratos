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

func (s *LushopService) SendSms(ctx context.Context, req *v1.SendSmsReq) (*v1.SendSmsReply, error) {
	return s.uc.SendSms(ctx, req)
}

func (s *LushopService) VerifySms(ctx context.Context, req *v1.VerifySmsReq) (*v1.RegisterReply, error) {
	return s.uc.VerifySms(ctx, req)
}

func (s *LushopService) Captcha(ctx context.Context, req *emptypb.Empty) (*v1.CaptchaReply, error) {
	return s.uc.GetCaptcha(ctx)
}
func (s *LushopService) Detail(ctx context.Context, req *emptypb.Empty) (*v1.UserDetailResponse, error) {
	return s.uc.UserDetailByID(ctx)
}

func (s *LushopService) Update(ctx context.Context, req *v1.UpdateReq) (*v1.UserDetailResponse, error) {
	return s.uc.UpdateUser(ctx, req)
}

func (s *LushopService) UpdatePwd(ctx context.Context, req *v1.UpdatePwdReq) (*v1.UpdatePwdReply, error) {
	return s.uc.UpdatePassword(ctx, req)
}

func (s *LushopService) Logout(ctx context.Context, req *emptypb.Empty) (*v1.LogoutReply, error) {
	return s.uc.Logout(ctx)
}

func (s *LushopService) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenReply, error) {
	return s.uc.RefreshToken(ctx, req)
}

// 管理员可查看用户列表
func (s *LushopService) ListUsers(ctx context.Context, req *v1.ListUsersReq) (*v1.ListUsersReply, error) {
	return s.uc.ListUsers(ctx, req)
}

// 管理员删除用户
// 管理员踢出用户（注销登录）
func (s *LushopService) KickUser(ctx context.Context, req *v1.KickUserReq) (*v1.KickUserReply, error) {
	return s.uc.DeleteUser(ctx, req)
}
