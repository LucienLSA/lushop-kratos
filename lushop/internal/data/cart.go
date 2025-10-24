package data

import (
	"context"
	orderv1 "lushop/api/service/order/v1"
	"lushop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// cartRepo 购物车仓储实现
type cartRepo struct {
	data        *Data
	orderClient orderv1.OrderClient
	log         *log.Helper
}

// NewCartRepo 创建购物车仓储
func NewCartRepo(data *Data, orderClient orderv1.OrderClient, logger log.Logger) biz.CartRepo {
	return &cartRepo{
		data:        data,
		orderClient: orderClient,
		log:         log.NewHelper(log.With(logger, "module", "data/cart")),
	}
}

// GetCartList 获取购物车列表
func (r *cartRepo) GetCartList(ctx context.Context, userId int32) ([]*biz.CartItem, int32, error) {
	// 调用 Order Service 的购物车列表接口
	cartList, err := r.orderClient.CartItemList(ctx, &orderv1.UserInfo{
		Id: userId,
	})
	if err != nil {
		r.log.Errorf("failed to get cart list from order service: %v", err)
		return nil, 0, err
	}

	// 转换为领域模型
	items := make([]*biz.CartItem, 0, len(cartList.Data))
	for _, item := range cartList.Data {
		items = append(items, &biz.CartItem{
			Id:      item.Id,
			UserId:  item.UserId,
			GoodsId: item.GoodsId,
			Nums:    item.Nums,
			Checked: item.Checked,
		})
	}

	return items, cartList.Total, nil
}

// AddToCart 添加商品到购物车
func (r *cartRepo) AddToCart(ctx context.Context, userId, goodsId, nums int32, checked bool) (*biz.CartItem, error) {
	// 调用 Order Service 的添加购物车接口
	cartItem, err := r.orderClient.CreateCartItem(ctx, &orderv1.CartItemRequest{
		UserId:  userId,
		GoodsId: goodsId,
		Nums:    nums,
		Checked: checked,
	})
	if err != nil {
		r.log.Errorf("failed to add to cart: %v", err)
		return nil, err
	}

	// 转换为领域模型
	return &biz.CartItem{
		Id:      cartItem.Id,
		UserId:  cartItem.UserId,
		GoodsId: cartItem.GoodsId,
		Nums:    cartItem.Nums,
		Checked: cartItem.Checked,
	}, nil
}

// UpdateCartItem 更新购物车项
func (r *cartRepo) UpdateCartItem(ctx context.Context, id, userId, nums int32, checked bool) error {
	// 调用 Order Service 的更新购物车接口
	_, err := r.orderClient.UpdateCartItem(ctx, &orderv1.CartItemRequest{
		Id:      id,
		UserId:  userId,
		Nums:    nums,
		Checked: checked,
	})
	if err != nil {
		r.log.Errorf("failed to update cart item: %v", err)
		return err
	}

	return nil
}

// DeleteCartItem 删除购物车项
func (r *cartRepo) DeleteCartItem(ctx context.Context, id, userId int32) error {
	// 调用 Order Service 的删除购物车接口
	_, err := r.orderClient.DeleteCartItem(ctx, &orderv1.CartItemRequest{
		Id:     id,
		UserId: userId,
	})
	if err != nil {
		r.log.Errorf("failed to delete cart item: %v", err)
		return err
	}

	return nil
}
