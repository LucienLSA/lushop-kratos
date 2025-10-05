package data

import (
	"context"

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
