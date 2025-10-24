package biz

import (
	"context"
	v1 "order/api/order/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetCartList 获取用户购物车列表
func (u *OrderUsecase) GetCartList(ctx context.Context, userID int32) ([]*v1.ShopCartInfoResponse, error) {
	return u.repo.GetCartList(ctx, userID)
}

// CreateCartItem 创建购物车项
func (u *OrderUsecase) CreateCartItem(ctx context.Context, req *v1.CartItemRequest) (*v1.ShopCartInfoResponse, error) {
	// 业务层参数验证
	if req.UserId == 0 {
		u.log.Error("user id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "user id is required")
	}
	if req.GoodsId == 0 {
		u.log.Error("goods id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "goods id is required")
	}
	if req.Nums <= 0 {
		u.log.Error("nums must be greater than 0")
		return nil, errors.BadRequest("INVALID_PARAMETER", "nums must be greater than 0")
	}

	// 调用数据层创建购物车项
	cart, err := u.repo.CreateCartItem(ctx, req)
	if err != nil {
		u.log.Errorf("failed to create cart item: user=%d, goods=%d, error=%v", req.UserId, req.GoodsId, err)
		return nil, err
	}

	u.log.Infof("cart item created successfully: user=%d, goods=%d, cart_id=%d", req.UserId, req.GoodsId, cart.Id)
	return cart, nil
}

// UpdateCartItem 更新购物车项
func (u *OrderUsecase) UpdateCartItem(ctx context.Context, req *v1.CartItemRequest) (*emptypb.Empty, error) {
	// 业务层参数验证
	if req.Id == 0 {
		u.log.Error("cart item id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "cart item id is required")
	}
	if req.UserId == 0 {
		u.log.Error("user id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "user id is required")
	}
	if req.Nums <= 0 {
		u.log.Error("nums must be greater than 0")
		return nil, errors.BadRequest("INVALID_PARAMETER", "nums must be greater than 0")
	}

	// 调用数据层更新购物车项
	result, err := u.repo.UpdateCartItem(ctx, req)
	if err != nil {
		u.log.Errorf("failed to update cart item: id=%d, user=%d, error=%v", req.Id, req.UserId, err)
		return nil, err
	}

	u.log.Infof("cart item updated successfully: id=%d, user=%d, nums=%d", req.Id, req.UserId, req.Nums)
	return result, nil
}

// DeleteCartItem 删除购物车项
func (u *OrderUsecase) DeleteCartItem(ctx context.Context, req *v1.CartItemRequest) (*emptypb.Empty, error) {
	// 业务层参数验证
	if req.Id == 0 {
		u.log.Error("cart item id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "cart item id is required")
	}
	if req.UserId == 0 {
		u.log.Error("user id is required")
		return nil, errors.BadRequest("INVALID_PARAMETER", "user id is required")
	}

	// 调用数据层删除购物车项
	result, err := u.repo.DeleteCartItem(ctx, req)
	if err != nil {
		u.log.Errorf("failed to delete cart item: id=%d, user=%d, error=%v", req.Id, req.UserId, err)
		return nil, err
	}

	u.log.Infof("cart item deleted successfully: id=%d, user=%d", req.Id, req.UserId)
	return result, nil
}
