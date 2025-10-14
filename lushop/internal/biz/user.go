package biz

import (
	"context"
	"errors"
	"fmt"
	v1 "lushop/api/lushop/v1"
	"lushop/internal/conf"
	"lushop/internal/pkg/captcha"
	"lushop/internal/pkg/middleware/auth"
	"lushop/internal/pkg/sms"
	"os"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-kratos/kratos/v2/log"
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
	// 管理员可查看用户列表
	ListUsers(ctx context.Context, req *v1.ListUsersReq) ([]*User, int, error)

	// redis
	StoreCaptcha(ctx context.Context, CaptchaId, Ans string) error
	StoreToken(ctx context.Context, key, token string, expiration time.Duration) error
	GetToken(ctx context.Context, key string) (string, error)
	DeleteToken(ctx context.Context, key string) error
	StoreRefreshToken(ctx context.Context, userId int64, token string, expiration time.Duration) error
	GetRefreshToken(ctx context.Context, userId int64) (string, error)
	DeleteRefreshToken(ctx context.Context, userId int64) error
	StoreLogoutBlacklist(ctx context.Context, userId int64) error
	CheckLogoutBlacklist(ctx context.Context, userId int64) (bool, error)
	GetTokenTTL(ctx context.Context, key string) (time.Duration, error)
	StoreLogoutBlacklistWithTTL(ctx context.Context, userId int64, ttl time.Duration) error

	// sms
	StoreSmsCode(ctx context.Context, mobile, code string, expiration time.Duration) error
	GetSmsCode(ctx context.Context, mobile string) (string, error)
	SetSmsCooldown(ctx context.Context, mobile string, expiration time.Duration) error
	CheckSmsCooldown(ctx context.Context, mobile string) (bool, error)
}

type UserUsecase struct {
	uRepo      UserRepo
	log        *log.Helper
	signingKey string // 这里是为了生存 token 的时候可以直接取配置文件里面的配置
	smsConf    *conf.Sms
}

func NewUserUsecase(repo UserRepo, logger log.Logger, conf *conf.Auth, smsConf *conf.Sms) *UserUsecase {
	helper := log.NewHelper(log.With(logger, "module", "usecase/lushop"))
	return &UserUsecase{uRepo: repo, log: helper, signingKey: conf.JwtKey, smsConf: smsConf}
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
	// 从上下文获取用户ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.uRepo.CheckLogoutBlacklist(ctx, uid)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, ErrAuthFailed
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
				aexpiresAt := now.Add(30 * time.Minute)
				aclaims := auth.CustomClaims{
					ID:          user.ID,
					NickName:    user.NickName,
					AuthorityId: user.Role,
					RegisteredClaims: jwt5.RegisteredClaims{
						NotBefore: jwt5.NewNumericDate(now),
						ExpiresAt: jwt5.NewNumericDate(aexpiresAt),
						Issuer:    "lucien",
					},
				}
				accessToken, err := auth.CreateToken(aclaims, uc.signingKey)
				if err != nil {
					return nil, ErrGenerateTokenFailed
				}
				// 将access token存入redis
				accessTokenKey := fmt.Sprintf("user_access_token:%d", user.ID)
				err = uc.uRepo.StoreToken(ctx, accessTokenKey, accessToken, 30*time.Minute)
				if err != nil {
					uc.log.Errorf("存储access token失败: %v", err)
					return nil, ErrGenerateTokenFailed
				}
				// 生成refresh token
				rexpiresAt := now.Add(24 * 7 * time.Hour)
				rclaims := auth.CustomClaims{
					ID: user.ID,
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
				err = uc.uRepo.StoreRefreshToken(ctx, user.ID, refreshToken, 7*24*time.Hour)
				if err != nil {
					uc.log.Errorf("存储refresh token失败: %v", err)
					return nil, ErrGenerateTokenFailed
				}
				// 删除redis的验证码
				return &v1.RegisterReply{
					Id:           user.ID,
					Mobile:       user.Mobile,
					Username:     user.NickName,
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
					ExpiredAt:    aexpiresAt.Unix(),
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
	rclaims := auth.CustomClaims{
		ID: creatuser.ID,
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
	err = uc.uRepo.StoreRefreshToken(ctx, creatuser.ID, refreshToken, 7*24*time.Hour)
	if err != nil {
		uc.log.Errorf("存储refresh token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}
	// 删除redis验证码
	return &v1.RegisterReply{
		Id:           creatuser.ID,
		Mobile:       creatuser.Mobile,
		Username:     creatuser.NickName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    aexpiresAt.Unix(),
	}, nil

}

// 发送手机验证码
func (uc *UserUsecase) SendSms(ctx context.Context, req *v1.SendSmsReq) (*v1.SendSmsReply, error) {
	// 检验手机号是否合法
	if len(req.Mobile) != 11 {
		return nil, ErrMobileInvalid
	}
	// 检验手机冷却时间（例如 60s 内不可重复发送）
	cooling, err := uc.uRepo.CheckSmsCooldown(ctx, req.Mobile)
	if err != nil {
		return nil, err
	}
	if cooling {
		return &v1.SendSmsReply{Success: true}, nil
	}

	// 生成验证码
	code := sms.GenerateSmsCode(6)

	// TODO: 集成实际短信发送通道
	// 从环境变量读取，优先于配置
	var apiKey, apiSecret string
	if uc.smsConf != nil {
		apiKey = uc.smsConf.ApiKey
		apiSecret = uc.smsConf.ApiSecret
	}
	if v := os.Getenv("SMS_API_KEY"); v != "" {
		apiKey = v
	}
	if v := os.Getenv("SMS_API_SECRET"); v != "" {
		apiSecret = v
	}

	if apiKey != "" && apiSecret != "" {
		config := &openapi.Config{
			AccessKeyId:     tea.String(apiKey),
			AccessKeySecret: tea.String(apiSecret),
			RegionId:        tea.String(uc.smsConf.RegionId),
		}
		config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
		client, _ := dysmsapi.NewClient(config)
		request := &dysmsapi.SendSmsRequest{}
		request.SetTemplateCode(uc.smsConf.TemplateCode)
		request.SetTemplateParam("{\"code\":" + code + "}")
		request.SetPhoneNumbers(req.Mobile)
		request.SetSignName(uc.smsConf.SignName)
		response, err := client.SendSms(request)
		if err != nil {
			return nil, err
		}
		uc.log.Infof("发送短信响应: %v", response)
		// 保存验证码到Redis，过期时间5分钟
		if err := uc.uRepo.StoreSmsCode(ctx, req.Mobile, code, 5*time.Minute); err != nil {
			return nil, err
		}
		// 设置发送冷却60秒
		if err := uc.uRepo.SetSmsCooldown(ctx, req.Mobile, 60*time.Second); err != nil {
			return nil, err
		}
	}
	return &v1.SendSmsReply{Success: true}, nil
}

// 验证手机验证码
func (uc *UserUsecase) VerifySms(ctx context.Context, req *v1.VerifySmsReq) (*v1.RegisterReply, error) {
	// 校验参数
	if len(req.Mobile) != 11 {
		return nil, ErrMobileInvalid
	}
	if len(req.SmsCode) != 6 {
		return nil, ErrCaptchaInvalid
	}

	// 获取并校验验证码
	code, err := uc.uRepo.GetSmsCode(ctx, req.Mobile)
	if err != nil || code == "" {
		return nil, ErrCaptchaInvalid
	}
	if code != req.SmsCode {
		return nil, ErrCaptchaInvalid
	}

	// 获取用户信息（必须已注册）
	user, err := uc.uRepo.UserByMobile(ctx, req.Mobile)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 生成token，与密码登录一致
	now := time.Now()
	aexpiresAt := now.Add(30 * time.Minute)
	aclaims := auth.CustomClaims{
		ID:          user.ID,
		NickName:    user.NickName,
		AuthorityId: user.Role,
		RegisteredClaims: jwt5.RegisteredClaims{
			NotBefore: jwt5.NewNumericDate(now),
			ExpiresAt: jwt5.NewNumericDate(aexpiresAt),
			Issuer:    "lucien",
		},
	}
	accessToken, err := auth.CreateToken(aclaims, uc.signingKey)
	if err != nil {
		return nil, ErrGenerateTokenFailed
	}
	accessTokenKey := fmt.Sprintf("user_access_token:%d", user.ID)
	if err := uc.uRepo.StoreToken(ctx, accessTokenKey, accessToken, 30*time.Minute); err != nil {
		uc.log.Errorf("存储access token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}

	// 生成refresh token
	rexpiresAt := now.Add(24 * 7 * time.Hour)
	rclaims := auth.CustomClaims{
		ID: user.ID,
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
	if err := uc.uRepo.StoreRefreshToken(ctx, user.ID, refreshToken, 7*24*time.Hour); err != nil {
		uc.log.Errorf("存储refresh token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}
	// 验证码使用后删除验证码
	return &v1.RegisterReply{
		Id:           user.ID,
		Mobile:       user.Mobile,
		Username:     user.NickName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    aexpiresAt.Unix(),
	}, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, req *v1.UpdateReq) (*v1.UserDetailResponse, error) {
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.uRepo.CheckLogoutBlacklist(ctx, uid)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, ErrAuthFailed
	}
	// 根据用户ID获取当前用户信息
	curruser, err := uc.uRepo.UserById(ctx, uid)
	if err != nil {
		return nil, err
	}

	// 解析生日字符串为时间戳
	var birthdayTimestamp int64
	if req.Birthday != "" {
		// 解析日期字符串 "2025-02-16" 格式
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			uc.log.Errorf("解析生日格式失败: %v", err)
			return nil, errors.New("invalid birthday format, expected YYYY-MM-DD")
		}
		birthdayTimestamp = birthday.Unix()
	} else {
		birthdayTimestamp = curruser.Birthday // 保持原有值
	}

	// 更新用户信息
	updateUser := &User{
		ID:       curruser.ID,
		Mobile:   curruser.Mobile,   // 手机号通常不允许修改
		Password: curruser.Password, // 密码通过单独的接口修改
		NickName: req.NickName,
		Birthday: birthdayTimestamp,
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
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.uRepo.CheckLogoutBlacklist(ctx, uid)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
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
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}

	// 删除Redis中的access token
	accessTokenKey := fmt.Sprintf("user_access_token:%d", uid)
	err := uc.uRepo.DeleteToken(ctx, accessTokenKey)
	if err != nil {
		uc.log.Errorf("删除access token失败: %v", err)
	}

	// 删除Redis中的refresh token
	err = uc.uRepo.DeleteRefreshToken(ctx, uid)
	if err != nil {
		uc.log.Errorf("删除refresh token失败: %v", err)
	}

	// 将用户登出信息存入redis黑名单
	ttl, err := uc.uRepo.GetTokenTTL(ctx, accessTokenKey)
	if err != nil {
		uc.log.Errorf("获取用户%d token TTL失败: %v", uid, err)
	}
	if ttl <= 0 {
		// 若获取失败或无效，使用一个安全的最小TTL（例如30分钟）
		ttl = 30 * time.Minute
	}
	if err := uc.uRepo.StoreLogoutBlacklistWithTTL(ctx, uid, ttl); err != nil {
		uc.log.Errorf("加入用户%d到登出黑名单失败: %v", uid, err)
		return nil, ErrAuthFailed
	}

	return &v1.LogoutReply{
		Success: true,
	}, nil
}

// RefreshToken 刷新token
func (uc *UserUsecase) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenReply, error) {
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.uRepo.CheckLogoutBlacklist(ctx, uid)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, ErrAuthFailed
	}
	// 解析refresh token
	token, err := jwt5.ParseWithClaims(req.RefreshToken, &auth.CustomClaims{}, func(token *jwt5.Token) (interface{}, error) {
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
	storedToken, err := uc.uRepo.GetRefreshToken(ctx, claims.ID)
	if err != nil || storedToken != req.RefreshToken {
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
	// 生成新的access token
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

	return &v1.RefreshTokenReply{
		AccessToken:  newAccessToken,
		RefreshToken: storedToken,
	}, nil
}

func (uc *UserUsecase) ListUsers(ctx context.Context, req *v1.ListUsersReq) (*v1.ListUsersReply, error) {
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.uRepo.CheckLogoutBlacklist(ctx, uid)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, ErrAuthFailed
	}
	// 检查管理员权限
	if !auth.IsAdmin(ctx) {
		return nil, errors.New("forbidden: admin access required")
	}

	list, total, err := uc.uRepo.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	users := make([]*v1.UserDetailResponse, 0)
	for _, user := range list {
		users = append(users, &v1.UserDetailResponse{
			Id:       user.ID,
			NickName: user.NickName,
			Mobile:   user.Mobile,
			Role:     int32(user.Role),
		})
	}

	return &v1.ListUsersReply{
		Users: users,
		Total: int32(total),
	}, nil
}

// 管理员注销删除用户
func (uc *UserUsecase) DeleteUser(ctx context.Context, req *v1.KickUserReq) (*v1.KickUserReply, error) {
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.uRepo.CheckLogoutBlacklist(ctx, uid)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, ErrAuthFailed
	}
	// 检查管理员权限
	if !auth.IsAdmin(ctx) {
		return nil, errors.New("forbidden: admin access required")
	}

	// 踢出目标用户：清理其令牌并加入黑名单
	// 1) 删除access token
	accessTokenKey := fmt.Sprintf("user_access_token:%d", req.GetId())
	if err := uc.uRepo.DeleteToken(ctx, accessTokenKey); err != nil {
		uc.log.Errorf("删除用户%d access token失败: %v", req.GetId(), err)
	}

	// 2) 删除refresh token
	if err := uc.uRepo.DeleteRefreshToken(ctx, req.GetId()); err != nil {
		uc.log.Errorf("删除用户%d refresh token失败: %v", req.GetId(), err)
	}

	// 3) 将用户加入登出黑名单，过期时间为 access token 剩余时间
	ttl, err := uc.uRepo.GetTokenTTL(ctx, accessTokenKey)
	if err != nil {
		uc.log.Errorf("获取用户%d token TTL失败: %v", req.GetId(), err)
	}
	if ttl <= 0 {
		// 若获取失败或无效，使用一个安全的最小TTL（例如30分钟）
		ttl = 30 * time.Minute
	}
	if err := uc.uRepo.StoreLogoutBlacklistWithTTL(ctx, req.GetId(), ttl); err != nil {
		uc.log.Errorf("加入用户%d到登出黑名单失败: %v", req.GetId(), err)
		return nil, ErrAuthFailed
	}

	return &v1.KickUserReply{Success: true}, nil
}
