package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

// AddressRepo is a Address repo.
type AddressRepo interface {
}

// AddressUsecase is a Address usecase.
type AddressUsecase struct {
	repo AddressRepo
	log  *log.Helper
}

// NewAddressUsecase new a Address usecase.
func NewAddressUsecase(repo AddressRepo, logger log.Logger) *AddressUsecase {
	return &AddressUsecase{repo: repo, log: log.NewHelper(logger)}
}
