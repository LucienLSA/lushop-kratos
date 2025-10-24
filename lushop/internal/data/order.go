package data

import (
	"context"
	v1 "lushop/api/lushop/v1"
	orderV1 "lushop/api/service/order/v1"
	"lushop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

// NewOrderRepo 创建订单仓库
func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateOrder 创建订单
func (r *orderRepo) CreateOrder(ctx context.Context, req *v1.CreateOrderReq) (*v1.CreateOrderReply, error) {
	r.log.Infof("调用 Order gRPC: CreateOrder address=%s", req.Address)

	// 调用订单服务 gRPC API
	resp, err := r.data.oc.CreateOrder(ctx, &orderV1.OrderRequest{
		Address: req.Address,
		Name:    req.Name,
		Mobile:  req.Mobile,
		Post:    req.Post,
	})
	if err != nil {
		r.log.Errorf("创建订单失败: error=%v", err)
		return nil, err
	}

	r.log.Infof("创建订单成功: orderSn=%s", resp.OrderSn)
	return &v1.CreateOrderReply{
		Id:          resp.Id,
		OrderSn:     resp.OrderSn,
		TotalAmount: resp.Total,
		Status:      "PAYING", // 待支付状态
	}, nil
}

// GetOrderList 获取订单列表
func (r *orderRepo) GetOrderList(ctx context.Context, page, pageSize int32, status string) ([]*v1.OrderInfo, int32, error) {
	r.log.Infof("调用 Order gRPC: OrderList page=%d, pageSize=%d", page, pageSize)

	// 调用订单服务 gRPC API
	resp, err := r.data.oc.OrderList(ctx, &orderV1.OrderFilterRequest{
		Pages:       page,
		PagePerNums: pageSize,
	})
	if err != nil {
		r.log.Errorf("获取订单列表失败: error=%v", err)
		return nil, 0, err
	}

	// 转换为业务对象
	orders := make([]*v1.OrderInfo, 0, len(resp.Data))
	for _, item := range resp.Data {
		orders = append(orders, &v1.OrderInfo{
			Id:          item.Id,
			OrderSn:     item.OrderSn,
			TotalAmount: item.Total,
			Status:      item.Status,
			Address:     item.Address,
			Name:        item.Name,
			Mobile:      item.Mobile,
			CreatedAt:   0, // addTime 是 string，需要转换
		})
	}

	r.log.Infof("获取订单列表成功: total=%d", resp.Total)
	return orders, resp.Total, nil
}

// GetOrderDetail 获取订单详情
func (r *orderRepo) GetOrderDetail(ctx context.Context, id int32) (*v1.GetOrderDetailReply, error) {
	r.log.Infof("调用 Order gRPC: OrderDetail id=%d", id)

	// 调用订单服务 gRPC API
	resp, err := r.data.oc.OrderDetail(ctx, &orderV1.OrderRequest{
		Id: id,
	})
	if err != nil {
		r.log.Errorf("获取订单详情失败: id=%d, error=%v", id, err)
		return nil, err
	}

	// 转换商品信息
	goods := make([]*v1.OrderGoodsInfo, 0, len(resp.Goods))
	for _, item := range resp.Goods {
		goods = append(goods, &v1.OrderGoodsInfo{
			GoodsId:    item.GoodsId,
			GoodsName:  item.GoodsName,
			GoodsImage: item.GoodsImage,
			GoodsPrice: item.GoodsPrice,
			Nums:       item.Nums,
		})
	}

	r.log.Infof("获取订单详情成功: orderSn=%s", resp.OrderInfo.OrderSn)
	return &v1.GetOrderDetailReply{
		Id:          resp.OrderInfo.Id,
		OrderSn:     resp.OrderInfo.OrderSn,
		TotalAmount: resp.OrderInfo.Total,
		Status:      resp.OrderInfo.Status,
		Address:     resp.OrderInfo.Address,
		Name:        resp.OrderInfo.Name,
		Mobile:      resp.OrderInfo.Mobile,
		Post:        resp.OrderInfo.Post,
		CreatedAt:   0, // addTime 是 string，需要转换
		Goods:       goods,
	}, nil
}

// CancelOrder 取消订单
func (r *orderRepo) CancelOrder(ctx context.Context, id int32) error {
	r.log.Infof("调用 Order gRPC: UpdateOrderStatus id=%d, status=TRADE_CLOSED", id)

	// 调用订单服务 gRPC API 更新订单状态为已关闭
	_, err := r.data.oc.UpdateOrderStatus(ctx, &orderV1.OrderStatus{
		OrderSn: "", // 通过 ID 查询
		Status:  "TRADE_CLOSED",
	})
	if err != nil {
		r.log.Errorf("取消订单失败: id=%d, error=%v", id, err)
		return err
	}

	r.log.Infof("取消订单成功: id=%d", id)
	return nil
}
