package biz

import (
	"context"
	userauthV1 "lushop/api/userauth/v1"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserAuthAdapter 封装对 UserAuth 服务的调用
// 用于在 biz 层统一调用认证服务
type UserAuthAdapter struct {
	client userauthV1.UserAuthClient
	log    *log.Helper
}

// NewUserAuthAdapter 创建认证服务适配器
func NewUserAuthAdapter(client userauthV1.UserAuthClient, logger log.Logger) *UserAuthAdapter {
	return &UserAuthAdapter{
		client: client,
		log:    log.NewHelper(log.With(logger, "module", "biz/user-auth-adapter")),
	}
}

// GetCaptcha 获取图形验证码
func (a *UserAuthAdapter) GetCaptcha(ctx context.Context) (*userauthV1.CaptchaReply, error) {
	reply, err := a.client.GetCaptcha(ctx, &emptypb.Empty{})
	if err != nil {
		a.log.Errorf("获取验证码失败: %v", err)
		return nil, err
	}
	return reply, nil
}

// VerifyCaptcha 校验图形验证码
func (a *UserAuthAdapter) VerifyCaptcha(ctx context.Context, captchaId, captcha string) (bool, error) {
	reply, err := a.client.VerifyCaptcha(ctx, &userauthV1.VerifyCaptchaReq{
		CaptchaId: captchaId,
		Captcha:   captcha,
	})
	if err != nil {
		a.log.Errorf("校验验证码失败: %v", err)
		return false, err
	}
	return reply.Success, nil
}

// SendSms 发送短信验证码
func (a *UserAuthAdapter) SendSms(ctx context.Context, mobile string) (bool, error) {
	reply, err := a.client.SendSms(ctx, &userauthV1.SendSmsReq{
		Mobile: mobile,
	})
	if err != nil {
		a.log.Errorf("发送短信失败: %v", err)
		return false, err
	}
	return reply.Success, nil
}

// VerifySms 校验短信验证码
func (a *UserAuthAdapter) VerifySms(ctx context.Context, mobile, smsCode string) (bool, error) {
	reply, err := a.client.VerifySms(ctx, &userauthV1.VerifySmsReq{
		Mobile:  mobile,
		SmsCode: smsCode,
	})
	if err != nil {
		a.log.Errorf("校验短信验证码失败: %v", err)
		return false, err
	}
	return reply.Success, nil
}

// IssueToken 签发 Token（用于登录成功后）
func (a *UserAuthAdapter) IssueToken(ctx context.Context, userId int64, nickName string, role int) (*userauthV1.TokenReply, error) {
	reply, err := a.client.IssueToken(ctx, &userauthV1.IssueTokenReq{
		UserId:   userId,
		NickName: nickName,
		Role:     int32(role),
	})
	if err != nil {
		a.log.Errorf("签发Token失败: %v", err)
		return nil, err
	}
	return reply, nil
}

// RefreshToken 刷新 Token
func (a *UserAuthAdapter) RefreshToken(ctx context.Context, refreshToken string) (*userauthV1.TokenReply, error) {
	reply, err := a.client.RefreshToken(ctx, &userauthV1.RefreshTokenReq{
		RefreshToken: refreshToken,
	})
	if err != nil {
		a.log.Errorf("刷新Token失败: %v", err)
		return nil, err
	}
	return reply, nil
}

// RevokeToken 撤销 Token（登出）
func (a *UserAuthAdapter) RevokeToken(ctx context.Context, userId int64) error {
	_, err := a.client.RevokeToken(ctx, &userauthV1.RevokeTokenReq{
		UserId: userId,
	})
	if err != nil {
		a.log.Errorf("撤销Token失败: %v", err)
		return err
	}
	return nil
}

// CheckBlacklist 检查用户是否在黑名单
func (a *UserAuthAdapter) CheckBlacklist(ctx context.Context, userId int64) (bool, error) {
	reply, err := a.client.CheckBlacklist(ctx, &userauthV1.CheckBlacklistReq{
		UserId: userId,
	})
	if err != nil {
		a.log.Errorf("检查黑名单失败: %v", err)
		return false, err
	}
	return reply.IsBlacklisted, nil
}
