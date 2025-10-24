package biz

import (
	"context"
	v1 "order/api/order/v1"
	"order/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ProviderSet is biz providers.
// Biz 层只提供业务逻辑，不提供客户端
var ProviderSet = wire.NewSet(NewOrderUsecase)

// Transaction 新增事务接口方法
type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

// GoodsService 商品服务接口
// 由 Client 层实现，用于调用商品服务
type GoodsService interface {
	// BatchGetGoods 批量获取商品信息
	BatchGetGoods(ctx context.Context, ids []int32) (map[int32]*domain.GoodsInfo, error)
}

// OrderRepo is a Order repo.
type OrderRepo interface {
	// 订单
	CreateOrder(ctx context.Context, in *v1.OrderRequest) (*v1.OrderInfoResponse, error)
	GetOrderList(ctx context.Context, req *v1.OrderFilterRequest) ([]*v1.OrderInfoResponse, int32, error)
	GetOrderDetail(ctx context.Context, req *v1.OrderRequest) (*v1.OrderInfoDetailResponse, error)
	UpdateOrderStatus(ctx context.Context, req *v1.OrderStatus) (*emptypb.Empty, error)
	HandleOrderTimeout(ctx context.Context, orderSn string) error
	// 购物车
	GetCartList(ctx context.Context, userID int32) ([]*v1.ShopCartInfoResponse, error)
	CreateCartItem(ctx context.Context, req *v1.CartItemRequest) (*v1.ShopCartInfoResponse, error)
	UpdateCartItem(ctx context.Context, req *v1.CartItemRequest) (*emptypb.Empty, error)
	DeleteCartItem(ctx context.Context, req *v1.CartItemRequest) (*emptypb.Empty, error)
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
