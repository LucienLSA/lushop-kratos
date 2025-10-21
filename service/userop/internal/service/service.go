package service

import (
	v1 "userop/api/userop/v1"
	"userop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserOpService)

// UserOpService is a userop service.
type UserOpService struct {
	v1.UnimplementedUserOpServer
	as  *biz.AddressUsecase
	fs  *biz.FavoriteUsecase
	log *log.Helper
}

// NewUserOpService new a userop service.
func NewUserOpService(as *biz.AddressUsecase, fs *biz.FavoriteUsecase, logger log.Logger) *UserOpService {
	return &UserOpService{
		as:  as,
		fs:  fs,
		log: log.NewHelper(logger),
	}
}
