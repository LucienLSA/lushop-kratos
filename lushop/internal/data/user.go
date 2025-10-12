package data

import (
	"context"
	"fmt"
	"time"

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
func (u *userRepo) StoreRefreshToken(ctx context.Context, userId int64, deviceId, token string, expiration time.Duration) error {
	key := fmt.Sprintf("user_refresh_token:%d:%s", userId, deviceId)
	err := u.data.rdb.Set(ctx, key, token, expiration).Err()
	if err != nil {
		u.log.Errorf("存储refresh token失败: %v", err)
		return err
	}
	return nil
}

// GetRefreshToken 从Redis获取refresh token
func (u *userRepo) GetRefreshToken(ctx context.Context, userId int64, deviceId string) (string, error) {
	key := fmt.Sprintf("user_refresh_token:%d:%s", userId, deviceId)
	token, err := u.data.rdb.Get(ctx, key).Result()
	if err != nil {
		u.log.Errorf("获取refresh token失败: %v", err)
		return "", err
	}
	return token, nil
}

// DeleteRefreshToken 从Redis删除refresh token
func (u *userRepo) DeleteRefreshToken(ctx context.Context, userId int64, deviceId string) error {
	key := fmt.Sprintf("user_refresh_token:%d:%s", userId, deviceId)
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
	}, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	// 调用user service的UpdateUser方法更新用户信息
	_, err := u.data.uc.UpdateUser(ctx, &userService.UpdateUserInfo{
		Id:       user.ID,
		NickName: user.NickName,
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
		Birthday: int64(updatedUser.Birthday),
		Gender:   updatedUser.Gender,
		Role:     int(updatedUser.Role),
	}, nil
}

func (u *userRepo) ListUser(ctx context.Context) ([]*biz.User, int, error) {

}
