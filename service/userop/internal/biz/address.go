package biz

import (
	"context"
	"userop/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

// AddressRepo is a Address repo.
type AddressRepo interface {
	GetAddressList(ctx context.Context, filter domain.Address) (*domain.AddressListResponse, error)
	CreateAddress(ctx context.Context, address domain.Address) error
	DeleteAddress(ctx context.Context, address domain.Address) error
	UpdateAddress(ctx context.Context, address domain.Address) error
	GetMessageList(ctx context.Context, filter domain.Message) (*domain.MessageListResponse, error)
	CreateMessage(ctx context.Context, msg domain.Message) error
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

func (f *AddressUsecase) GetAddressList(ctx context.Context, filter domain.Address) (*domain.AddressListResponse, error) {
	return f.repo.GetAddressList(ctx, filter)
}

func (f *AddressUsecase) CreateAddress(ctx context.Context, address domain.Address) error {
	return f.repo.CreateAddress(ctx, address)
}

func (f *AddressUsecase) DeleteAddress(ctx context.Context, address domain.Address) error {
	return f.repo.DeleteAddress(ctx, address)
}

func (f *AddressUsecase) UpdateAddress(ctx context.Context, address domain.Address) error {
	return f.repo.UpdateAddress(ctx, address)
}

func (f *AddressUsecase) GetMessageList(ctx context.Context, filter domain.Message) (*domain.MessageListResponse, error) {
	return f.repo.GetMessageList(ctx, filter)
}

func (f *AddressUsecase) CreateMessage(ctx context.Context, msg domain.Message) error {
	return f.repo.CreateMessage(ctx, msg)
}
