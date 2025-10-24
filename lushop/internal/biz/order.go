package biz

import (
	"context"
	v1 "lushop/api/lushop/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// OrderRepo 订单仓库接口
type OrderRepo interface {
	CreateOrder(ctx context.Context, req *v1.CreateOrderReq) (*v1.CreateOrderReply, error)
	GetOrderList(ctx context.Context, page, pageSize int32, status string) ([]*v1.OrderInfo, int32, error)
	GetOrderDetail(ctx context.Context, id int32) (*v1.GetOrderDetailReply, error)
	CancelOrder(ctx context.Context, id int32) error
}

// OrderUsecase 订单用例
type OrderUsecase struct {
	repo OrderRepo
	log  *log.Helper
}

// NewOrderUsecase 创建订单用例
func NewOrderUsecase(repo OrderRepo, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// CreateOrder 创建订单
func (uc *OrderUsecase) CreateOrder(ctx context.Context, req *v1.CreateOrderReq) (*v1.CreateOrderReply, error) {
	uc.log.Infof("创建订单: address=%s, name=%s", req.Address, req.Name)
	return uc.repo.CreateOrder(ctx, req)
}

// GetOrderList 获取订单列表
func (uc *OrderUsecase) GetOrderList(ctx context.Context, req *v1.GetOrderListReq) (*v1.GetOrderListReply, error) {
	uc.log.Infof("获取订单列表: page=%d, pageSize=%d", req.Page, req.PageSize)
	
	orders, total, err := uc.repo.GetOrderList(ctx, req.Page, req.PageSize, req.Status)
	if err != nil {
		return nil, err
	}
	
	return &v1.GetOrderListReply{
		Total: total,
		Data:  orders,
	}, nil
}

// GetOrderDetail 获取订单详情
func (uc *OrderUsecase) GetOrderDetail(ctx context.Context, req *v1.GetOrderDetailReq) (*v1.GetOrderDetailReply, error) {
	uc.log.Infof("获取订单详情: id=%d", req.Id)
	return uc.repo.GetOrderDetail(ctx, req.Id)
}

// CancelOrder 取消订单
func (uc *OrderUsecase) CancelOrder(ctx context.Context, req *v1.CancelOrderReq) error {
	uc.log.Infof("取消订单: id=%d", req.Id)
	return uc.repo.CancelOrder(ctx, req.Id)
}
