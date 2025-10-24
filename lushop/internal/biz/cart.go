package biz

import (
	"context"
	v1 "lushop/api/lushop/v1"
	"lushop/internal/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// CartItem 购物车项领域模型
type CartItem struct {
	Id      int32
	UserId  int32
	GoodsId int32
	Nums    int32
	Checked bool
}

// CartRepo 购物车仓储接口
type CartRepo interface {
	// 获取购物车列表
	GetCartList(ctx context.Context, userId int32) ([]*CartItem, int32, error)
	// 添加商品到购物车
	AddToCart(ctx context.Context, userId, goodsId, nums int32, checked bool) (*CartItem, error)
	// 更新购物车项
	UpdateCartItem(ctx context.Context, id, userId, nums int32, checked bool) error
	// 删除购物车项
	DeleteCartItem(ctx context.Context, id, userId int32) error
}

// CartUsecase 购物车业务用例
type CartUsecase struct {
	repo CartRepo
	log  *log.Helper
}

// NewCartUsecase 创建购物车业务用例
func NewCartUsecase(repo CartRepo, logger log.Logger) *CartUsecase {
	return &CartUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "usecase/cart")),
	}
}

// GetCartList 获取购物车列表
func (uc *CartUsecase) GetCartList(ctx context.Context) (*v1.CartListReply, error) {
	// 从上下文获取用户ID
	userId, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, errors.Unauthorized("UNAUTHORIZED", "please login first")
	}

	// 调用仓储层获取购物车列表
	items, total, err := uc.repo.GetCartList(ctx, int32(userId))
	if err != nil {
		uc.log.Errorf("failed to get cart list: %v", err)
		return nil, err
	}

	// 转换为响应格式
	cartItems := make([]*v1.CartItemReply, 0, len(items))
	for _, item := range items {
		cartItems = append(cartItems, &v1.CartItemReply{
			Id:      item.Id,
			UserId:  item.UserId,
			GoodsId: item.GoodsId,
			Nums:    item.Nums,
			Checked: item.Checked,
		})
	}

	return &v1.CartListReply{
		Total: total,
		Data:  cartItems,
	}, nil
}

// AddToCart 添加商品到购物车
func (uc *CartUsecase) AddToCart(ctx context.Context, req *v1.AddToCartReq) (*v1.CartItemReply, error) {
	// 从上下文获取用户ID
	userId, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, errors.Unauthorized("UNAUTHORIZED", "please login first")
	}

	// 业务规则验证
	if req.GoodsId <= 0 {
		return nil, errors.BadRequest("INVALID_GOODS_ID", "invalid goods id")
	}
	if req.Nums <= 0 {
		return nil, errors.BadRequest("INVALID_NUMS", "nums must be greater than 0")
	}

	// 调用仓储层添加购物车
	item, err := uc.repo.AddToCart(ctx, int32(userId), req.GoodsId, req.Nums, req.Checked)
	if err != nil {
		uc.log.Errorf("failed to add to cart: %v", err)
		return nil, err
	}

	return &v1.CartItemReply{
		Id:      item.Id,
		UserId:  item.UserId,
		GoodsId: item.GoodsId,
		Nums:    item.Nums,
		Checked: item.Checked,
	}, nil
}

// UpdateCartItem 更新购物车项
func (uc *CartUsecase) UpdateCartItem(ctx context.Context, req *v1.UpdateCartItemReq) error {
	// 从上下文获取用户ID
	userId, ok := auth.GetUserID(ctx)
	if !ok {
		return errors.Unauthorized("UNAUTHORIZED", "please login first")
	}

	// 业务规则验证
	if req.Id <= 0 {
		return errors.BadRequest("INVALID_CART_ITEM_ID", "invalid cart item id")
	}
	if req.Nums <= 0 {
		return errors.BadRequest("INVALID_NUMS", "nums must be greater than 0")
	}

	// 调用仓储层更新购物车
	err := uc.repo.UpdateCartItem(ctx, req.Id, int32(userId), req.Nums, req.Checked)
	if err != nil {
		uc.log.Errorf("failed to update cart item: %v", err)
		return err
	}

	return nil
}

// DeleteCartItem 删除购物车项
func (uc *CartUsecase) DeleteCartItem(ctx context.Context, req *v1.DeleteCartItemReq) error {
	// 从上下文获取用户ID
	userId, ok := auth.GetUserID(ctx)
	if !ok {
		return errors.Unauthorized("UNAUTHORIZED", "please login first")
	}

	// 业务规则验证
	if req.Id <= 0 {
		return errors.BadRequest("INVALID_CART_ITEM_ID", "invalid cart item id")
	}

	// 调用仓储层删除购物车
	err := uc.repo.DeleteCartItem(ctx, req.Id, int32(userId))
	if err != nil {
		uc.log.Errorf("failed to delete cart item: %v", err)
		return err
	}

	return nil
}
