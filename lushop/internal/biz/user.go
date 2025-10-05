package biz

import (
	"context"
	"errors"
	v1 "lushop/api/lushop/v1"
	"lushop/internal/conf"
	"lushop/internal/pkg/captcha"
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
	CreateUser(c context.Context, u *User) (*User, error)
	UserByMobile(ctx context.Context, mobile string) (*User, error)
	UserById(ctx context.Context, Id int64) (*User, error)
	CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error)
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

// 创建用户，用户注册创建后也提供登录状态
func (uc *UserUsecase) CreateUser(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterReply, error) {
	newUser, err := newUser(req.Mobile, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	creatuser, err := uc.uRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)
	claims := auth.CustomClaims{
		ID:          creatuser.ID,
		NickName:    creatuser.NickName,
		AuthorityId: creatuser.Role,
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
		Id:        creatuser.ID,
		Mobile:    creatuser.Mobile,
		Username:  creatuser.NickName,
		Token:     token,
		ExpiredAt: time.Now().Unix() + 60*60*24*30,
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
