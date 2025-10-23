package biz

import (
	"context"
	v1 "order/api/order/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// FavoriteRepo is a Favorite repo.
type OrderRepo interface {
	Create(ctx context.Context, in *v1.OrderRequest) (*v1.OrderInfoResponse, error)
	GetCartList(ctx context.Context, userID int32) ([]*v1.ShopCartInfoResponse, error)
}

// OrderUsecase is a Order usecase.
type OrderUsecase struct {
	repo OrderRepo
	log  *log.Helper
}

// NewOrderUsecase new a Order usecase.
func NewOrderUsecase(repo OrderRepo, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, in *v1.OrderRequest) (*v1.OrderInfoResponse, error) {
	return &v1.OrderInfoResponse{}, nil
}

// GetCartList 获取用户购物车列表
func (u *OrderUsecase) GetCartList(ctx context.Context, userID int32) ([]*v1.ShopCartInfoResponse, error) {
	return u.repo.GetCartList(ctx, userID)
}
