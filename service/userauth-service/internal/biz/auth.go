package biz

import (
	"context"
	"fmt"
	"time"

	v1 "userauth-service/api/userauth/v1"
	"userauth-service/internal/conf"
	"userauth-service/internal/pkg/auth"
	"userauth-service/internal/pkg/captcha"
	"userauth-service/internal/pkg/sms"

	"github.com/go-kratos/kratos/v2/log"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

type AuthRepo interface {
	// 验证码
	StoreCaptcha(ctx context.Context, captchaId, ans string, expiration time.Duration) error
	GetCaptcha(ctx context.Context, captchaId string) (string, error)

	// 短信
	StoreSmsCode(ctx context.Context, mobile, code string, expiration time.Duration) error
	GetSmsCode(ctx context.Context, mobile string) (string, error)
	SetSmsCooldown(ctx context.Context, mobile string, expiration time.Duration) error
	CheckSmsCooldown(ctx context.Context, mobile string) (bool, error)

	// Token
	StoreAccessToken(ctx context.Context, userId int64, token string, expiration time.Duration) error
	StoreRefreshToken(ctx context.Context, userId int64, token string, expiration time.Duration) error
	GetRefreshToken(ctx context.Context, userId int64) (string, error)
	DeleteTokens(ctx context.Context, userId int64) error

	// 黑名单
	AddToBlacklist(ctx context.Context, userId int64, ttl time.Duration) error
	CheckBlacklist(ctx context.Context, userId int64) (bool, error)
}

type AuthUsecase struct {
	repo     AuthRepo
	authConf *conf.Auth
	smsConf  *conf.Sms
	log      *log.Helper
}

func NewAuthUsecase(repo AuthRepo, authConf *conf.Auth, smsConf *conf.Sms, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{
		repo:     repo,
		authConf: authConf,
		smsConf:  smsConf,
		log:      log.NewHelper(log.With(logger, "module", "biz/auth")),
	}
}

// GetCaptcha 生成图形验证码
func (uc *AuthUsecase) GetCaptcha(ctx context.Context) (*v1.CaptchaReply, error) {
	captchaInfo, err := captcha.GetCaptcha(ctx)
	if err != nil {
		return nil, err
	}

	// 存储验证码（5分钟过期）
	err = uc.repo.StoreCaptcha(ctx, captchaInfo.CaptchaId, captchaInfo.Ans, 5*time.Minute)
	if err != nil {
		return nil, err
	}

	return &v1.CaptchaReply{
		CaptchaId: captchaInfo.CaptchaId,
		PicPath:   captchaInfo.PicPath,
	}, nil
}

// VerifyCaptcha 校验图形验证码
func (uc *AuthUsecase) VerifyCaptcha(ctx context.Context, captchaId, captcha string) (bool, error) {
	storedAns, err := uc.repo.GetCaptcha(ctx, captchaId)
	if err != nil || storedAns == "" {
		return false, nil
	}
	return storedAns == captcha, nil
}

// SendSms 发送短信验证码
func (uc *AuthUsecase) SendSms(ctx context.Context, mobile string) (bool, error) {
	// 检查冷却
	cooling, err := uc.repo.CheckSmsCooldown(ctx, mobile)
	if err != nil {
		return false, err
	}
	if cooling {
		uc.log.Infof("手机号 %s 处于冷却期，跳过发送", mobile)
		return true, nil // 冷却中，直接返回成功
	}

	// 生成验证码
	code := sms.GenerateSmsCode(6)
	uc.log.Infof("为手机号 %s 生成验证码: %s", mobile, code)

	// 发送短信（集成阿里云）
	err = sms.SendSms(uc.smsConf, mobile, code)
	if err != nil {
		uc.log.Errorf("发送短信失败: %v", err)
		return false, err
	}

	// 存储验证码（5分钟过期）
	err = uc.repo.StoreSmsCode(ctx, mobile, code, 5*time.Minute)
	if err != nil {
		return false, err
	}

	// 设置冷却（60秒）
	err = uc.repo.SetSmsCooldown(ctx, mobile, 60*time.Second)
	if err != nil {
		return false, err
	}

	return true, nil
}

// VerifySms 校验短信验证码
func (uc *AuthUsecase) VerifySms(ctx context.Context, mobile, smsCode string) (bool, error) {
	storedCode, err := uc.repo.GetSmsCode(ctx, mobile)
	if err != nil || storedCode == "" {
		return false, nil
	}
	return storedCode == smsCode, nil
}

// IssueToken 签发 Token
func (uc *AuthUsecase) IssueToken(ctx context.Context, userId int64, nickName string, role int32) (*v1.TokenReply, error) {
	now := time.Now()

	// 生成 Access Token（30分钟）
	accessExpiry := now.Add(30 * time.Minute)
	accessToken, err := auth.CreateToken(auth.CustomClaims{
		ID:          userId,
		NickName:    nickName,
		AuthorityId: int(role),
		RegisteredClaims: jwt5.RegisteredClaims{
			NotBefore: jwt5.NewNumericDate(now),
			ExpiresAt: jwt5.NewNumericDate(accessExpiry),
			Issuer:    "userauth-service",
		},
	}, uc.authConf.JwtKey)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token（7天）
	refreshExpiry := now.Add(7 * 24 * time.Hour)
	refreshToken, err := auth.CreateToken(auth.CustomClaims{
		ID: userId,
		RegisteredClaims: jwt5.RegisteredClaims{
			NotBefore: jwt5.NewNumericDate(now),
			ExpiresAt: jwt5.NewNumericDate(refreshExpiry),
			Issuer:    "userauth-service",
		},
	}, uc.authConf.JwtKey)
	if err != nil {
		return nil, err
	}

	// 存储到 Redis
	err = uc.repo.StoreAccessToken(ctx, userId, accessToken, 30*time.Minute)
	if err != nil {
		return nil, err
	}
	err = uc.repo.StoreRefreshToken(ctx, userId, refreshToken, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	uc.log.Infof("为用户 %d (%s) 签发 Token 成功", userId, nickName)

	return &v1.TokenReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    accessExpiry.Unix(),
	}, nil
}

// RefreshToken 刷新 Token
func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*v1.TokenReply, error) {
	// 解析 Refresh Token
	token, err := jwt5.ParseWithClaims(refreshToken, &auth.CustomClaims{}, func(token *jwt5.Token) (interface{}, error) {
		return []byte(uc.authConf.JwtKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*auth.CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// 验证 Refresh Token 是否存在于 Redis
	storedToken, err := uc.repo.GetRefreshToken(ctx, claims.ID)
	if err != nil || storedToken != refreshToken {
		return nil, fmt.Errorf("refresh token not found or mismatch")
	}

	// 签发新的 Access Token（保持 Refresh Token 不变）
	now := time.Now()
	accessExpiry := now.Add(30 * time.Minute)
	newAccessToken, err := auth.CreateToken(auth.CustomClaims{
		ID:          claims.ID,
		NickName:    claims.NickName,
		AuthorityId: claims.AuthorityId,
		RegisteredClaims: jwt5.RegisteredClaims{
			NotBefore: jwt5.NewNumericDate(now),
			ExpiresAt: jwt5.NewNumericDate(accessExpiry),
			Issuer:    "userauth-service",
		},
	}, uc.authConf.JwtKey)
	if err != nil {
		return nil, err
	}

	// 更新 Redis 中的 Access Token
	err = uc.repo.StoreAccessToken(ctx, claims.ID, newAccessToken, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	uc.log.Infof("为用户 %d 刷新 Token 成功", claims.ID)

	return &v1.TokenReply{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    accessExpiry.Unix(),
	}, nil
}

// RevokeToken 撤销 Token（登出）
func (uc *AuthUsecase) RevokeToken(ctx context.Context, userId int64) error {
	// 删除 Redis 中的 Token
	err := uc.repo.DeleteTokens(ctx, userId)
	if err != nil {
		uc.log.Errorf("删除用户 %d 的 Token 失败: %v", userId, err)
	}

	// 加入黑名单（30分钟，与 Access Token 过期时间一致）
	err = uc.repo.AddToBlacklist(ctx, userId, 30*time.Minute)
	if err != nil {
		return err
	}

	uc.log.Infof("用户 %d 登出成功，已加入黑名单", userId)
	return nil
}

// CheckBlacklist 检查黑名单
func (uc *AuthUsecase) CheckBlacklist(ctx context.Context, userId int64) (bool, error) {
	return uc.repo.CheckBlacklist(ctx, userId)
}

// AddToBlacklist 添加到黑名单
func (uc *AuthUsecase) AddToBlacklist(ctx context.Context, userId int64, ttl int64) error {
	duration := 30 * time.Minute // 默认 TTL
	if ttl > 0 {
		duration = time.Duration(ttl) * time.Second
	}
	return uc.repo.AddToBlacklist(ctx, userId, duration)
}
