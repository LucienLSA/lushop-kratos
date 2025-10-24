package data

import (
	"context"
	v1 "order/api/order/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type ShoppingCart struct {
	BaseModel
	User    int32 `gorm:"type:int;index;comment:用户ID"` //在购物车列表中我们需要查询当前用户的购物车记录
	Goods   int32 `gorm:"type:int;index;comment:商品ID"` //加索引：我们需要查询时候， 1. 会影响插入性能 2. 会占用磁盘
	Nums    int32 `gorm:"type:int;comment:数量"`
	Checked bool  `gorm:"comment:是否选中"` //是否选中
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

// findCartItemByUserAndGoods 查询用户购物车中是否存在指定商品
// 返回购物车项和错误信息
// 如果返回 gorm.ErrRecordNotFound 表示购物车中没有该商品
func (r *orderRepo) findCartItemByUserAndGoods(ctx context.Context, userId, goodsId int32) (*ShoppingCart, error) {
	var cart ShoppingCart
	err := r.data.DB(ctx).Where(&ShoppingCart{
		User:  userId,
		Goods: goodsId,
	}).First(&cart).Error

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

// GetCartList 获取用户购物车列表
func (r *orderRepo) GetCartList(ctx context.Context, userID int32) ([]*v1.ShopCartInfoResponse, error) {
	var carts []ShoppingCart

	// 查询该用户的所有购物车记录
	if err := r.data.DB(ctx).WithContext(ctx).Where(&ShoppingCart{User: userID}).Find(&carts).Error; err != nil {
		r.log.Errorf("failed to get cart list for user %d: %v", userID, err)
		return nil, err
	}

	// 转换为响应格式
	result := make([]*v1.ShopCartInfoResponse, 0, len(carts))
	for _, cart := range carts {
		result = append(result, &v1.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}

	return result, nil
}

// CreateCartItem 创建或更新购物车项
func (r *orderRepo) CreateCartItem(ctx context.Context, req *v1.CartItemRequest) (*v1.ShopCartInfoResponse, error) {
	// 先查询该用户是否已经有该商品在购物车中
	cartPtr, err := r.findCartItemByUserAndGoods(ctx, req.UserId, req.GoodsId)

	if err != nil && err != gorm.ErrRecordNotFound {
		// 数据库查询错误
		r.log.Errorf("failed to query cart item: %v", err)
		return nil, err
	}

	var cart ShoppingCart
	if cartPtr != nil {
		cart = *cartPtr
	}

	if err == gorm.ErrRecordNotFound {
		// 购物车中没有该商品，创建新记录
		cart = ShoppingCart{
			User:    req.UserId,
			Goods:   req.GoodsId,
			Nums:    req.Nums,
			Checked: req.Checked,
		}

		if err := r.data.DB(ctx).Create(&cart).Error; err != nil {
			r.log.Errorf("failed to create cart item: %v", err)
			return nil, err
		}

		r.log.Infof("created new cart item: user=%d, goods=%d, nums=%d", req.UserId, req.GoodsId, req.Nums)
	} else {
		// 购物车中已有该商品，更新数量
		cart.Nums += req.Nums
		if req.Checked {
			cart.Checked = req.Checked
		}

		if err := r.data.DB(ctx).Save(&cart).Error; err != nil {
			r.log.Errorf("failed to update cart item: %v", err)
			return nil, err
		}

		r.log.Infof("updated cart item: user=%d, goods=%d, nums=%d", req.UserId, req.GoodsId, cart.Nums)
	}

	// 返回购物车项信息
	return &v1.ShopCartInfoResponse{
		Id:      cart.ID,
		UserId:  cart.User,
		GoodsId: cart.Goods,
		Nums:    cart.Nums,
		Checked: cart.Checked,
	}, nil
}

// UpdateCartItem 更新购物车项
func (r *orderRepo) UpdateCartItem(ctx context.Context, req *v1.CartItemRequest) (*emptypb.Empty, error) {
	// 查询购物车项是否存在
	var cart ShoppingCart
	err := r.data.DB(ctx).Where("id = ? AND user = ?", req.Id, req.UserId).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Errorf("cart item not found: id=%d, user=%d", req.Id, req.UserId)
			return nil, errors.NotFound("CART_ITEM_NOT_FOUND", "cart item not found")
		}
		r.log.Errorf("failed to query cart item: %v", err)
		return nil, err
	}

	// 更新购物车项
	cart.Nums = req.Nums
	cart.Checked = req.Checked

	if err := r.data.DB(ctx).Save(&cart).Error; err != nil {
		r.log.Errorf("failed to update cart item: %v", err)
		return nil, err
	}

	r.log.Infof("cart item updated: id=%d, user=%d, nums=%d", req.Id, req.UserId, req.Nums)
	return &emptypb.Empty{}, nil
}

// DeleteCartItem 删除购物车项
func (r *orderRepo) DeleteCartItem(ctx context.Context, req *v1.CartItemRequest) (*emptypb.Empty, error) {
	// 删除购物车项（确保只删除属于该用户的购物车项）
	result := r.data.DB(ctx).Where("id = ? AND user = ?", req.Id, req.UserId).Delete(&ShoppingCart{})

	if result.Error != nil {
		r.log.Errorf("failed to delete cart item: %v", result.Error)
		return nil, result.Error
	}

	// 检查是否真的删除了记录
	if result.RowsAffected == 0 {
		r.log.Errorf("cart item not found or not owned by user: id=%d, user=%d", req.Id, req.UserId)
		return nil, errors.NotFound("CART_ITEM_NOT_FOUND", "cart item not found or not owned by user")
	}

	r.log.Infof("cart item deleted: id=%d, user=%d", req.Id, req.UserId)
	return &emptypb.Empty{}, nil
}
