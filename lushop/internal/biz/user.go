package biz

import (
	"context"
	"errors"
	"fmt"
	v1 "lushop/api/lushop/v1"
	"lushop/internal/conf"
	"lushop/internal/pkg/captcha"
	"lushop/internal/pkg/device"
	"lushop/internal/pkg/middleware/auth"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

// 定义错误
var (
	ErrPasswordInvalid     = errors.New("password invalid")
	ErrUsernameInvalid     = errors.New("username invalid")
	ErrCaptchaInvalid      = errors.New("verification code error")
	ErrMobileInvalid       = errors.New("mobile invalid")
	ErrUserNotFound        = errors.New("user not found")
	ErrLoginFailed         = errors.New("login failed")
	ErrGenerateTokenFailed = errors.New("generate token failed")
	ErrAuthFailed          = errors.New("authentication failed")
	ErrPasswordRepeated    = errors.New("password repeated")
)

// 定义返回的数据的结构体
type User struct {
	ID        int64
	Mobile    string
	Password  string
	NickName  string
	Birthday  int64
	Gender    string
	Role      int
	CreatedAt time.Time
}

type UserRepo interface {
	// mysql
	CreateUser(c context.Context, u *User) (*User, error)
	UserByMobile(ctx context.Context, mobile string) (*User, error)
	UserById(ctx context.Context, Id int64) (*User, error)
	UpdateUser(ctx context.Context, u *User) (*User, error)
	CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error)
	ListUser(ctx context.Context) ([]*User, int, error)
	// redis
	StoreCaptcha(ctx context.Context, CaptchaId, Ans string) error
	StoreToken(ctx context.Context, key, token string, expiration time.Duration) error
	GetToken(ctx context.Context, key string) (string, error)
	DeleteToken(ctx context.Context, key string) error
	StoreRefreshToken(ctx context.Context, userId int64, deviceId, token string, expiration time.Duration) error
	GetRefreshToken(ctx context.Context, userId int64, deviceId string) (string, error)
	DeleteRefreshToken(ctx context.Context, userId int64, deviceId string) error
}

type UserUsecase struct {
	uRepo      UserRepo
	log        *log.Helper
	signingKey string // 这里是为了生存 token 的时候可以直接取配置文件里面的配置
}

func NewUserUsecase(repo UserRepo, logger log.Logger, conf *conf.Auth) *UserUsecase {
	helper := log.NewHelper(log.With(logger, "module", "usecase/lushop"))
	return &UserUsecase{uRepo: repo, log: helper, signingKey: conf.JwtKey}
}

// 获取验证码
func (uc *UserUsecase) GetCaptcha(ctx context.Context) (*v1.CaptchaReply, error) {
	captchaInfo, err := captcha.GetCaptcha(ctx)
	if err != nil {
		return nil, err
	}
	// 将验证码存入Redis（5分钟过期）这里存入captcha_id作为key，ans作为value
	err = uc.uRepo.StoreCaptcha(ctx, captchaInfo.CaptchaId, captchaInfo.Ans)
	if err != nil {
		return nil, err
	}
	return &v1.CaptchaReply{
		CaptchaId: captchaInfo.CaptchaId,
		PicPath:   captchaInfo.PicPath,
		Ans:       captchaInfo.Ans,
	}, nil

}

// 用户ID获取详情
func (uc *UserUsecase) UserDetailByID(ctx context.Context) (*v1.UserDetailResponse, error) {
	// 从上下文取出claims用户权限信息
	var uid int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt5.MapClaims)
		if c["ID"] == nil {
			return nil, ErrAuthFailed
		}
		uid = int64(c["ID"].(float64))
	}
	user, err := uc.uRepo.UserById(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &v1.UserDetailResponse{
		Id:       user.ID,
		NickName: user.NickName,
		Mobile:   user.Mobile,
	}, nil
}

// 用户密码登录
func (uc *UserUsecase) PasswordLogin(ctx context.Context, req *v1.LoginReq) (*v1.RegisterReply, error) {
	// 表单验证
	if len(req.Mobile) <= 0 {
		return nil, ErrMobileInvalid
	}
	if len(req.Password) <= 0 {
		return nil, ErrUsernameInvalid
	}
	// 验证验证码是否正确
	if !captcha.Store.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, ErrCaptchaInvalid
	}
	// 手机号验证
	if user, err := uc.uRepo.UserByMobile(ctx, req.Mobile); err != nil {
		return nil, ErrUserNotFound
	} else {
		// 检查密码
		if passRsp, pasErr := uc.uRepo.CheckPassword(ctx, req.Password, user.Password); pasErr != nil {
			return nil, ErrPasswordInvalid
		} else {
			if passRsp {
				now := time.Now()
				expiresAt := now.Add(30 * 24 * time.Hour)
				claims := auth.CustomClaims{
					ID:          user.ID,
					NickName:    user.NickName,
					AuthorityId: user.Role,
					RegisteredClaims: jwt5.RegisteredClaims{
						NotBefore: jwt5.NewNumericDate(now),
						ExpiresAt: jwt5.NewNumericDate(expiresAt),
						Issuer:    "lucien",
					},
				}
				token, err := auth.CreateToken(claims, uc.signingKey)
				if err != nil {
					return nil, ErrGenerateTokenFailed
				}
				return &v1.RegisterReply{
					Id:        user.ID,
					Mobile:    user.Mobile,
					Username:  user.NickName,
					Token:     token,
					ExpiredAt: time.Now().Unix() + 60*60*24*30,
				}, nil
			} else {
				return nil, ErrLoginFailed
			}
		}
	}
}

// 用户结构体生成
func newUser(mobile, username, password string) (User, error) {
	if len(mobile) <= 0 || len(mobile) > 13 {
		return User{}, ErrMobileInvalid
	}
	if len(username) <= 0 {
		return User{}, ErrUsernameInvalid
	}
	if len(password) <= 0 {
		return User{}, ErrPasswordInvalid
	}
	return User{
		Mobile:   mobile,
		NickName: username,
		Password: password,
	}, nil
}

// 创建用户，用户注册创建后也提供登录状态
func (uc *UserUsecase) CreateUser(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterReply, error) {
	// 验证验证码是否正确
	if !captcha.Store.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, ErrCaptchaInvalid
	}
	newUser, err := newUser(req.Mobile, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	creatuser, err := uc.uRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	aexpiresAt := now.Add(30 * time.Minute)
	aclaims := auth.CustomClaims{
		ID:          creatuser.ID,
		NickName:    creatuser.NickName,
		AuthorityId: creatuser.Role,
		RegisteredClaims: jwt5.RegisteredClaims{
			NotBefore: jwt5.NewNumericDate(now),
			ExpiresAt: jwt5.NewNumericDate(aexpiresAt),
			Issuer:    "lucien",
		},
	}
	// 生成accesstoken
	accessToken, err := auth.CreateToken(aclaims, uc.signingKey)
	if err != nil {
		return nil, ErrGenerateTokenFailed
	}

	// 将access token存入redis
	accessTokenKey := fmt.Sprintf("user_access_token:%d", creatuser.ID)
	err = uc.uRepo.StoreToken(ctx, accessTokenKey, accessToken, 30*time.Minute)
	if err != nil {
		uc.log.Errorf("存储access token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}

	// 生成refresh token
	rexpiresAt := now.Add(24 * 7 * time.Hour)
	deviceId := device.GetDeviceFingerprint(ctx)
	rclaims := auth.CustomClaims{
		ID:       creatuser.ID,
		DeviceId: deviceId,
		RegisteredClaims: jwt5.RegisteredClaims{
			NotBefore: jwt5.NewNumericDate(now),
			ExpiresAt: jwt5.NewNumericDate(rexpiresAt),
			Issuer:    "lucien",
		},
	}
	refreshToken, err := auth.CreateToken(rclaims, uc.signingKey)
	if err != nil {
		return nil, ErrGenerateTokenFailed
	}

	// 将refresh token存入redis
	err = uc.uRepo.StoreRefreshToken(ctx, creatuser.ID, deviceId, refreshToken, 7*24*time.Hour)
	if err != nil {
		uc.log.Errorf("存储refresh token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}

	return &v1.RegisterReply{
		Id:        creatuser.ID,
		Mobile:    creatuser.Mobile,
		Username:  creatuser.NickName,
		Token:     accessToken,
		ExpiredAt: aexpiresAt.Unix(),
	}, nil

}

func (uc *UserUsecase) UpdateUser(ctx context.Context, req *v1.UpdateReq) (*v1.UserDetailResponse, error) {
	// 从JWT token中获取当前登录用户的ID
	var uid int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt5.MapClaims)
		if c["ID"] == nil {
			return nil, ErrAuthFailed
		}
		uid = int64(c["ID"].(float64))
	} else {
		return nil, ErrAuthFailed
	}

	// 根据用户ID获取当前用户信息
	curruser, err := uc.uRepo.UserById(ctx, uid)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	updateUser := &User{
		ID:       curruser.ID,
		Mobile:   curruser.Mobile,   // 手机号通常不允许修改
		Password: curruser.Password, // 密码通过单独的接口修改
		NickName: req.NickName,
		Birthday: req.Birthday,
		Gender:   req.Gender,
		Role:     curruser.Role, // 角色通常不允许用户自己修改
	}

	// 调用repository更新用户信息
	updatedUser, err := uc.uRepo.UpdateUser(ctx, updateUser)
	if err != nil {
		return nil, err
	}

	// 返回更新后的用户信息
	return &v1.UserDetailResponse{
		Id:       updatedUser.ID,
		NickName: updatedUser.NickName,
		Mobile:   updatedUser.Mobile,
		Birthday: updatedUser.Birthday,
		Gender:   updatedUser.Gender,
		Role:     int32(updatedUser.Role),
	}, nil
}

func (uc *UserUsecase) UpdatePassword(ctx context.Context, req *v1.UpdatePwdReq) (*v1.UpdatePwdReply, error) {
	// 从JWT token中获取当前登录用户的ID
	var uid int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt5.MapClaims)
		if c["ID"] == nil {
			return nil, ErrAuthFailed
		}
		uid = int64(c["ID"].(float64))
	} else {
		return nil, ErrAuthFailed
	}

	// 根据用户ID获取当前用户信息
	curruser, err := uc.uRepo.UserById(ctx, uid)
	if err != nil {
		return nil, err
	}

	// 验证旧密码是否正确
	passRsp, pasErr := uc.uRepo.CheckPassword(ctx, req.OldPassword, curruser.Password)
	if pasErr != nil {
		return nil, ErrPasswordInvalid
	}
	if !passRsp {
		return nil, ErrPasswordInvalid
	}

	// 检查新密码是否与旧密码相同
	if req.OldPassword == req.NewPassword {
		return nil, ErrPasswordRepeated
	}

	// 更新用户密码（密码加密由user service层处理）
	updateUser := &User{
		ID:       curruser.ID,
		Mobile:   curruser.Mobile, // 手机号不允许修改
		Password: req.NewPassword, // 传递明文密码，由user service加密
		NickName: curruser.NickName,
		Birthday: curruser.Birthday,
		Gender:   curruser.Gender,
		Role:     curruser.Role, // 角色不允许用户自己修改
	}

	// 调用repository更新用户信息
	_, err = uc.uRepo.UpdateUser(ctx, updateUser)
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePwdReply{
		Success: true,
	}, nil
}

func (uc *UserUsecase) Logout(ctx context.Context) (*v1.LogoutReply, error) {
	// 从JWT token中获取当前登录用户的ID
	var uid int64
	var deviceId string
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt5.MapClaims)
		if c["ID"] == nil {
			return nil, ErrAuthFailed
		}
		uid = int64(c["ID"].(float64))
		if c["DeviceId"] != nil {
			deviceId = c["DeviceId"].(string)
		}
	} else {
		return nil, ErrAuthFailed
	}

	// 删除Redis中的token
	accessTokenKey := fmt.Sprintf("user_access_token:%d", uid)
	err := uc.uRepo.DeleteToken(ctx, accessTokenKey)
	if err != nil {
		uc.log.Errorf("删除access token失败: %v", err)
	}

	// 删除Redis中的refresh token
	if deviceId != "" {
		err = uc.uRepo.DeleteRefreshToken(ctx, uid, deviceId)
		if err != nil {
			uc.log.Errorf("删除refresh token失败: %v", err)
		}
	}

	return &v1.LogoutReply{
		Success: true,
	}, nil
}

func (uc *UserUsecase) List(ctx context.Context) ([]*v1.UserDetailResponse, int, error) {
	list, total, err := uc.uRepo.ListUser(ctx)
	if err != nil {
		return nil, 0, err
	}
	rv := make([]*v1.UserDetailResponse, 0)
	for _, user := range list {
		rv = append(rv, &v1.UserDetailResponse{
			Id:       user.ID,
			NickName: user.NickName,
			Mobile:   user.Mobile,
		})
	}
	return rv, total, nil
}

// RefreshToken 刷新token
func (uc *UserUsecase) RefreshToken(ctx context.Context, refreshToken string) (*v1.RegisterReply, error) {
	// 解析refresh token
	token, err := jwt5.ParseWithClaims(refreshToken, &auth.CustomClaims{}, func(token *jwt5.Token) (interface{}, error) {
		return []byte(uc.signingKey), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrAuthFailed
	}

	claims, ok := token.Claims.(*auth.CustomClaims)
	if !ok {
		return nil, ErrAuthFailed
	}

	// 验证refresh token是否存在于Redis中
	storedToken, err := uc.uRepo.GetRefreshToken(ctx, claims.ID, claims.DeviceId)
	if err != nil || storedToken != refreshToken {
		return nil, ErrAuthFailed
	}

	// 获取用户信息
	user, err := uc.uRepo.UserById(ctx, claims.ID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 生成新的access token
	now := time.Now()
	expiresAt := now.Add(30 * time.Minute)
	newAccessClaims := auth.CustomClaims{
		ID:          user.ID,
		NickName:    user.NickName,
		AuthorityId: user.Role,
		RegisteredClaims: jwt5.RegisteredClaims{
			NotBefore: jwt5.NewNumericDate(now),
			ExpiresAt: jwt5.NewNumericDate(expiresAt),
			Issuer:    "lucien",
		},
	}

	newAccessToken, err := auth.CreateToken(newAccessClaims, uc.signingKey)
	if err != nil {
		return nil, ErrGenerateTokenFailed
	}

	// 存储新的access token
	accessTokenKey := fmt.Sprintf("user_access_token:%d", user.ID)
	err = uc.uRepo.StoreToken(ctx, accessTokenKey, newAccessToken, 30*time.Minute)
	if err != nil {
		uc.log.Errorf("存储新access token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}

	return &v1.RegisterReply{
		Id:        user.ID,
		Mobile:    user.Mobile,
		Username:  user.NickName,
		Token:     newAccessToken,
		ExpiredAt: expiresAt.Unix(),
	}, nil
}
