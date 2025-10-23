package biz

import (
	"context"
	"errors"
	v1 "lushop/api/lushop/v1"
	"lushop/internal/conf"
	"lushop/internal/pkg/middleware/auth"
	"time"

	"github.com/go-kratos/kratos/v2/log"
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
	// 用户资料相关（通过 User 服务）
	CreateUser(c context.Context, u *User) (*User, error)
	UserByMobile(ctx context.Context, mobile string) (*User, error)
	UserById(ctx context.Context, Id int64) (*User, error)
	UpdateUser(ctx context.Context, u *User) (*User, error)
	CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error)
	ListUsers(ctx context.Context, req *v1.ListUsersReq) ([]*User, int, error)

	// 黑名单相关（通过 UserAuth 服务）
	StoreLogoutBlacklist(ctx context.Context, userId int64) error
	CheckLogoutBlacklist(ctx context.Context, userId int64) (bool, error)
	StoreLogoutBlacklistWithTTL(ctx context.Context, userId int64, ttl time.Duration) error
}

type UserUsecase struct {
	uRepo       UserRepo
	authAdapter *UserAuthAdapter // 用户认证服务适配器
	log         *log.Helper
	signingKey  string // 保留用于向后兼容，实际 Token 签发由 UserAuth 服务处理
	smsConf     *conf.Sms
}

func NewUserUsecase(repo UserRepo, authAdapter *UserAuthAdapter, logger log.Logger, conf *conf.Auth, smsConf *conf.Sms) *UserUsecase {
	helper := log.NewHelper(log.With(logger, "module", "usecase/lushop"))
	return &UserUsecase{
		uRepo:       repo,
		authAdapter: authAdapter,
		log:         helper,
		signingKey:  conf.JwtKey,
		smsConf:     smsConf,
	}
}

// 获取验证码 - 通过 UserAuth 服务
func (uc *UserUsecase) GetCaptcha(ctx context.Context) (*v1.CaptchaReply, error) {
	reply, err := uc.authAdapter.GetCaptcha(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CaptchaReply{
		CaptchaId: reply.CaptchaId,
		PicPath:   reply.PicPath,
		// 注意：不返回 ans，安全考虑
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
	isBlacklisted, err := uc.authAdapter.CheckBlacklist(ctx, uid)
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

// 用户密码登录 - 通过 UserAuth 服务签发 Token
func (uc *UserUsecase) PasswordLogin(ctx context.Context, req *v1.LoginReq) (*v1.RegisterReply, error) {
	// 表单验证
	if len(req.Mobile) <= 0 {
		return nil, ErrMobileInvalid
	}
	if len(req.Password) <= 0 {
		return nil, ErrUsernameInvalid
	}
	
	// 验证验证码是否正确 - 通过 UserAuth 服务
	success, err := uc.authAdapter.VerifyCaptcha(ctx, req.CaptchaId, req.Captcha)
	if err != nil || !success {
		return nil, ErrCaptchaInvalid
	}
	
	// 手机号验证
	user, err := uc.uRepo.UserByMobile(ctx, req.Mobile)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	// 检查密码
	passRsp, pasErr := uc.uRepo.CheckPassword(ctx, req.Password, user.Password)
	if pasErr != nil {
		return nil, ErrPasswordInvalid
	}
	if !passRsp {
		return nil, ErrLoginFailed
	}
	
	// 通过 UserAuth 服务签发 Token
	tokenReply, err := uc.authAdapter.IssueToken(ctx, user.ID, user.NickName, user.Role)
	if err != nil {
		uc.log.Errorf("签发Token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}
	
	return &v1.RegisterReply{
		Id:           user.ID,
		Mobile:       user.Mobile,
		Username:     user.NickName,
		AccessToken:  tokenReply.AccessToken,
		RefreshToken: tokenReply.RefreshToken,
		ExpiredAt:    tokenReply.ExpiredAt,
	}, nil
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

// 创建用户，用户注册创建后也提供登录状态 - 通过 UserAuth 服务签发 Token
func (uc *UserUsecase) CreateUser(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterReply, error) {
	// 验证验证码是否正确 - 通过 UserAuth 服务
	success, err := uc.authAdapter.VerifyCaptcha(ctx, req.CaptchaId, req.Captcha)
	if err != nil || !success {
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
	
	// 通过 UserAuth 服务签发 Token
	tokenReply, err := uc.authAdapter.IssueToken(ctx, creatuser.ID, creatuser.NickName, creatuser.Role)
	if err != nil {
		uc.log.Errorf("签发Token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}
	
	return &v1.RegisterReply{
		Id:           creatuser.ID,
		Mobile:       creatuser.Mobile,
		Username:     creatuser.NickName,
		AccessToken:  tokenReply.AccessToken,
		RefreshToken: tokenReply.RefreshToken,
		ExpiredAt:    tokenReply.ExpiredAt,
	}, nil
}

// 发送手机验证码 - 通过 UserAuth 服务
func (uc *UserUsecase) SendSms(ctx context.Context, req *v1.SendSmsReq) (*v1.SendSmsReply, error) {
	// 检验手机号是否合法
	if len(req.Mobile) != 11 {
		return nil, ErrMobileInvalid
	}
	
	// 通过 UserAuth 服务发送短信（包含验证码生成、冷却控制、短信发送）
	success, err := uc.authAdapter.SendSms(ctx, req.Mobile)
	if err != nil {
		uc.log.Errorf("发送短信失败: %v", err)
		return nil, err
	}
	
	return &v1.SendSmsReply{Success: success}, nil
}

// 验证手机验证码 - 通过 UserAuth 服务
func (uc *UserUsecase) VerifySms(ctx context.Context, req *v1.VerifySmsReq) (*v1.RegisterReply, error) {
	// 校验参数
	if len(req.Mobile) != 11 {
		return nil, ErrMobileInvalid
	}
	if len(req.SmsCode) != 6 {
		return nil, ErrCaptchaInvalid
	}

	// 通过 UserAuth 服务校验短信验证码
	success, err := uc.authAdapter.VerifySms(ctx, req.Mobile, req.SmsCode)
	if err != nil || !success {
		return nil, ErrCaptchaInvalid
	}

	// 获取用户信息（必须已注册）
	user, err := uc.uRepo.UserByMobile(ctx, req.Mobile)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 通过 UserAuth 服务签发 Token
	tokenReply, err := uc.authAdapter.IssueToken(ctx, user.ID, user.NickName, user.Role)
	if err != nil {
		uc.log.Errorf("签发Token失败: %v", err)
		return nil, ErrGenerateTokenFailed
	}

	return &v1.RegisterReply{
		Id:           user.ID,
		Mobile:       user.Mobile,
		Username:     user.NickName,
		AccessToken:  tokenReply.AccessToken,
		RefreshToken: tokenReply.RefreshToken,
		ExpiredAt:    tokenReply.ExpiredAt,
	}, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, req *v1.UpdateReq) (*v1.UserDetailResponse, error) {
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.authAdapter.CheckBlacklist(ctx, uid)
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
	isBlacklisted, err := uc.authAdapter.CheckBlacklist(ctx, uid)
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

	// 通过 UserAuth 服务撤销 Token（包含删除 Token 和加入黑名单）
	err := uc.authAdapter.RevokeToken(ctx, uid)
	if err != nil {
		uc.log.Errorf("撤销Token失败: %v", err)
		return nil, ErrAuthFailed
	}

	return &v1.LogoutReply{
		Success: true,
	}, nil
}

// RefreshToken 刷新token - 通过 UserAuth 服务
func (uc *UserUsecase) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenReply, error) {
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.authAdapter.CheckBlacklist(ctx, uid)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, ErrAuthFailed
	}
	
	// 通过 UserAuth 服务刷新 Token
	tokenReply, err := uc.authAdapter.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		uc.log.Errorf("刷新Token失败: %v", err)
		return nil, ErrAuthFailed
	}

	return &v1.RefreshTokenReply{
		AccessToken:  tokenReply.AccessToken,
		RefreshToken: tokenReply.RefreshToken,
	}, nil
}

func (uc *UserUsecase) ListUsers(ctx context.Context, req *v1.ListUsersReq) (*v1.ListUsersReply, error) {
	// 从上下文获取当前登录用户的ID
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, ErrAuthFailed
	}
	// 检查用户是否在黑名单中
	isBlacklisted, err := uc.authAdapter.CheckBlacklist(ctx, uid)
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
	isBlacklisted, err := uc.authAdapter.CheckBlacklist(ctx, uid)
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

	// 踢出目标用户：撤销其令牌并加入黑名单（通过 UserAuth 服务）
	err = uc.authAdapter.RevokeToken(ctx, req.GetId())
	if err != nil {
		uc.log.Errorf("撤销用户%d Token失败: %v", req.GetId(), err)
		return nil, ErrAuthFailed
	}

	return &v1.KickUserReply{Success: true}, nil
}
