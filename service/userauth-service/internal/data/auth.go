package data

import (
	"context"
	"fmt"
	"time"

	"userauth-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type authRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/auth")),
	}
}

// 验证码
func (r *authRepo) StoreCaptcha(ctx context.Context, captchaId, ans string, expiration time.Duration) error {
	key := fmt.Sprintf("captcha:%s", captchaId)
	return r.data.rdb.Set(ctx, key, ans, expiration).Err()
}

func (r *authRepo) GetCaptcha(ctx context.Context, captchaId string) (string, error) {
	key := fmt.Sprintf("captcha:%s", captchaId)
	val, err := r.data.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// 短信
func (r *authRepo) StoreSmsCode(ctx context.Context, mobile, code string, expiration time.Duration) error {
	key := fmt.Sprintf("sms_code:%s", mobile)
	return r.data.rdb.Set(ctx, key, code, expiration).Err()
}

func (r *authRepo) GetSmsCode(ctx context.Context, mobile string) (string, error) {
	key := fmt.Sprintf("sms_code:%s", mobile)
	val, err := r.data.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *authRepo) SetSmsCooldown(ctx context.Context, mobile string, expiration time.Duration) error {
	key := fmt.Sprintf("sms_cooldown:%s", mobile)
	return r.data.rdb.Set(ctx, key, "1", expiration).Err()
}

func (r *authRepo) CheckSmsCooldown(ctx context.Context, mobile string) (bool, error) {
	key := fmt.Sprintf("sms_cooldown:%s", mobile)
	exists, err := r.data.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Token
func (r *authRepo) StoreAccessToken(ctx context.Context, userId int64, token string, expiration time.Duration) error {
	key := fmt.Sprintf("user_access_token:%d", userId)
	return r.data.rdb.Set(ctx, key, token, expiration).Err()
}

func (r *authRepo) StoreRefreshToken(ctx context.Context, userId int64, token string, expiration time.Duration) error {
	key := fmt.Sprintf("user_refresh_token:%d", userId)
	return r.data.rdb.Set(ctx, key, token, expiration).Err()
}

func (r *authRepo) GetRefreshToken(ctx context.Context, userId int64) (string, error) {
	key := fmt.Sprintf("user_refresh_token:%d", userId)
	val, err := r.data.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("refresh token not found")
	}
	return val, err
}

func (r *authRepo) DeleteTokens(ctx context.Context, userId int64) error {
	accessKey := fmt.Sprintf("user_access_token:%d", userId)
	refreshKey := fmt.Sprintf("user_refresh_token:%d", userId)
	pipe := r.data.rdb.Pipeline()
	pipe.Del(ctx, accessKey)
	pipe.Del(ctx, refreshKey)
	_, err := pipe.Exec(ctx)
	return err
}

// 黑名单
func (r *authRepo) AddToBlacklist(ctx context.Context, userId int64, ttl time.Duration) error {
	key := fmt.Sprintf("user_blacklist:%d", userId)
	return r.data.rdb.Set(ctx, key, "1", ttl).Err()
}

func (r *authRepo) CheckBlacklist(ctx context.Context, userId int64) (bool, error) {
	key := fmt.Sprintf("user_blacklist:%d", userId)
	exists, err := r.data.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
