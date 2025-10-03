package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// 定义返回数据结构体
type User struct {
	ID       int64
	Mobile   string
	Password string
	NickName string
	Birthday int64
	Gender   string
	Role     int
}

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type UserRepo interface {
	CreateUser(context.Context, *User) (*User, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) Create(ctx context.Context, u *User) (*User, error) {
	return uc.repo.CreateUser(ctx, u)
}
