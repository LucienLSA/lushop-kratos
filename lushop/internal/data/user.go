package data

import (
	"context"
	"fmt"
	"time"

	v1 "lushop/api/lushop/v1"
	userService "lushop/api/service/user/v1"
	"lushop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/user")),
	}
}

func (u *userRepo) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	createUser, err := u.data.uc.CreateUser(ctx, &userService.CreateUserInfo{
		NickName: user.NickName,
		Password: user.Password,
		Mobile:   user.Mobile,
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:       createUser.Id,
		Mobile:   createUser.Mobile,
		NickName: createUser.NickName,
	}, nil
}

func (u *userRepo) StoreCaptcha(ctx context.Context, captchaId, ans string) error {
	err := u.data.rdb.SetNX(ctx, "captcha:"+captchaId, ans, 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

// StoreToken 存储access token到Redis
func (u *userRepo) StoreToken(ctx context.Context, key, token string, expiration time.Duration) error {
	err := u.data.rdb.Set(ctx, key, token, expiration).Err()
	if err != nil {
		u.log.Errorf("存储token失败: %v", err)
		return err
	}
	return nil
}

// GetToken 从Redis获取token
func (u *userRepo) GetToken(ctx context.Context, key string) (string, error) {
	token, err := u.data.rdb.Get(ctx, key).Result()
	if err != nil {
		u.log.Errorf("获取token失败: %v", err)
		return "", err
	}
	return token, nil
}

// DeleteToken 从Redis删除token
func (u *userRepo) DeleteToken(ctx context.Context, key string) error {
	err := u.data.rdb.Del(ctx, key).Err()
	if err != nil {
		u.log.Errorf("删除token失败: %v", err)
		return err
	}
	return nil
}

// StoreRefreshToken 存储refresh token到Redis
func (u *userRepo) StoreRefreshToken(ctx context.Context, userId int64, token string, expiration time.Duration) error {
	key := fmt.Sprintf("user_refresh_token:%d", userId)
	err := u.data.rdb.Set(ctx, key, token, expiration).Err()
	if err != nil {
		u.log.Errorf("存储refresh token失败: %v", err)
		return err
	}
	return nil
}

// GetRefreshToken 从Redis获取refresh token
func (u *userRepo) GetRefreshToken(ctx context.Context, userId int64) (string, error) {
	key := fmt.Sprintf("user_refresh_token:%d", userId)
	token, err := u.data.rdb.Get(ctx, key).Result()
	if err != nil {
		u.log.Errorf("获取refresh token失败: %v", err)
		return "", err
	}
	return token, nil
}

// DeleteRefreshToken 从Redis删除refresh token
func (u *userRepo) DeleteRefreshToken(ctx context.Context, userId int64) error {
	key := fmt.Sprintf("user_refresh_token:%d", userId)
	err := u.data.rdb.Del(ctx, key).Err()
	if err != nil {
		u.log.Errorf("删除refresh token失败: %v", err)
		return err
	}
	return nil
}

func (u *userRepo) UserByMobile(ctx context.Context, mobile string) (*biz.User, error) {
	byMobile, err := u.data.uc.GetUserByMobile(ctx, &userService.MobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Mobile:   byMobile.Mobile,
		ID:       byMobile.Id,
		NickName: byMobile.NickName,
		Password: byMobile.Password,
		Role:     int(byMobile.Role),
		Birthday: int64(byMobile.Birthday),
		Gender:   byMobile.Gender,
	}, nil
}

func (u *userRepo) CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error) {
	if byMobile, err := u.data.uc.CheckPassword(ctx, &userService.PasswordCheckInfo{
		Password:          password,
		EncryptedPassword: encryptedPassword,
	}); err != nil {
		return false, err
	} else {
		return byMobile.Success, nil
	}
}

func (u *userRepo) UserById(ctx context.Context, id int64) (*biz.User, error) {
	user, err := u.data.uc.GetUserById(ctx, &userService.IdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:       user.Id,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int(user.Role),
		Password: user.Password,
	}, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	// 调用user service的UpdateUser方法更新用户信息
	_, err := u.data.uc.UpdateUser(ctx, &userService.UpdateUserInfo{
		Id:       user.ID,
		NickName: user.NickName,
		Password: user.Password,
		Gender:   user.Gender,
		Birthday: uint64(user.Birthday),
	})
	if err != nil {
		return nil, err
	}

	// 更新成功后，重新获取用户信息以返回最新的数据
	updatedUser, err := u.data.uc.GetUserById(ctx, &userService.IdRequest{
		Id: user.ID,
	})
	if err != nil {
		return nil, err
	}

	// 返回更新后的用户信息
	return &biz.User{
		ID:       updatedUser.Id,
		Mobile:   updatedUser.Mobile,
		NickName: updatedUser.NickName,
		Password: updatedUser.Password,
		Birthday: int64(updatedUser.Birthday),
		Gender:   updatedUser.Gender,
		Role:     int(updatedUser.Role),
	}, nil
}

func (u *userRepo) StoreLogoutBlacklist(ctx context.Context, userId int64) error {
	key := fmt.Sprintf("logout_blacklist:%d", userId)
	timestamp := time.Now().Unix()
	err := u.data.rdb.Set(ctx, key, timestamp, 24*time.Hour).Err()
	if err != nil {
		u.log.Errorf("存储登出黑名单失败: %v", err)
		return err
	}
	return nil
}

// CheckLogoutBlacklist 检查用户是否在黑名单中
func (u *userRepo) CheckLogoutBlacklist(ctx context.Context, userId int64) (bool, error) {
	key := fmt.Sprintf("logout_blacklist:%d", userId)
	exists, err := u.data.rdb.Exists(ctx, key).Result()
	if err != nil {
		u.log.Errorf("检查登出黑名单失败: %v", err)
		return false, err
	}
	return exists > 0, nil
}

// StoreSmsCode 将短信验证码存入Redis
func (u *userRepo) StoreSmsCode(ctx context.Context, mobile, code string, expiration time.Duration) error {
	key := fmt.Sprintf("sms_code:%s", mobile)
	if err := u.data.rdb.Set(ctx, key, code, expiration).Err(); err != nil {
		u.log.Errorf("存储短信验证码失败: %v", err)
		return err
	}
	return nil
}

// GetSmsCode 获取短信验证码
func (u *userRepo) GetSmsCode(ctx context.Context, mobile string) (string, error) {
	key := fmt.Sprintf("sms_code:%s", mobile)
	code, err := u.data.rdb.Get(ctx, key).Result()
	if err != nil {
		u.log.Errorf("获取短信验证码失败: %v", err)
		return "", err
	}
	return code, nil
}

// SetSmsCooldown 设置短信发送冷却标记
func (u *userRepo) SetSmsCooldown(ctx context.Context, mobile string, expiration time.Duration) error {
	key := fmt.Sprintf("sms_cooldown:%s", mobile)
	if err := u.data.rdb.Set(ctx, key, 1, expiration).Err(); err != nil {
		u.log.Errorf("设置短信冷却失败: %v", err)
		return err
	}
	return nil
}

// CheckSmsCooldown 检查手机号是否处于冷却中
func (u *userRepo) CheckSmsCooldown(ctx context.Context, mobile string) (bool, error) {
	key := fmt.Sprintf("sms_cooldown:%s", mobile)
	exists, err := u.data.rdb.Exists(ctx, key).Result()
	if err != nil {
		u.log.Errorf("检查短信冷却失败: %v", err)
		return false, err
	}
	return exists > 0, nil
}

func (u *userRepo) ListUsers(ctx context.Context, req *v1.ListUsersReq) ([]*biz.User, int, error) {
	// 调用用户服务获取用户列表
	pageInfo := &userService.PageInfo{
		Pn:    uint32(req.Page),
		PSize: uint32(req.PageSize),
	}

	userListResp, err := u.data.uc.GetUserList(ctx, pageInfo)
	if err != nil {
		return nil, 0, err
	}

	// 转换为业务层结构
	result := make([]*biz.User, 0, len(userListResp.Data))
	for _, user := range userListResp.Data {
		result = append(result, &biz.User{
			ID:       user.Id,
			Mobile:   user.Mobile,
			NickName: user.NickName,
			Role:     int(user.Role),
			Password: user.Password,
		})
	}

	return result, int(userListResp.Total), nil
}
