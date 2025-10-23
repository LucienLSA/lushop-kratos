package service

import (
	v1 "order/api/order/v1"
	"order/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewOrderService)

// OrderService is a order service.
type OrderService struct {
	v1.UnimplementedOrderServer

	uc  *biz.OrderUsecase
	log *log.Helper
}

// NewOrderService new a order service.
func NewOrderService(uc *biz.OrderUsecase, logger log.Logger) *OrderService {
	return &OrderService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}
