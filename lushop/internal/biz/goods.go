package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Goods 商品业务对象
type Goods struct {
	ID              int32
	Name            string
	GoodsSn         string
	CategoryID      int32
	ClickNum        int32
	SoldNum         int32
	FavNum          int32
	MarketPrice     float32
	ShopPrice       float32
	GoodsBrief      string
	GoodsDesc       string
	ShipFree        bool
	Images          []string
	DescImages      []string
	GoodsFrontImage string
	IsNew           bool
	IsHot           bool
	OnSale          bool
}

// GoodsRepo 商品仓库接口
type GoodsRepo interface {
	// GetGoodsList 获取商品列表
	GetGoodsList(ctx context.Context, page, pageSize int32, isHot, isNew bool) ([]*Goods, int32, error)

	// GetGoodsDetail 获取商品详情
	GetGoodsDetail(ctx context.Context, id int32) (*Goods, error)

	// SearchGoods 搜索商品
	SearchGoods(ctx context.Context, keyword string, page, pageSize int32) ([]*Goods, int32, error)

	// BatchGetGoods 批量获取商品信息
	BatchGetGoods(ctx context.Context, ids []int32) (map[int32]*Goods, error)
}

// GoodsUsecase 商品业务逻辑
type GoodsUsecase struct {
	repo GoodsRepo
	log  *log.Helper
}

// NewGoodsUsecase 创建商品业务逻辑
func NewGoodsUsecase(repo GoodsRepo, logger log.Logger) *GoodsUsecase {
	return &GoodsUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// GetGoodsList 获取商品列表
func (uc *GoodsUsecase) GetGoodsList(ctx context.Context, page, pageSize int32, isHot, isNew bool) ([]*Goods, int32, error) {
	return uc.repo.GetGoodsList(ctx, page, pageSize, isHot, isNew)
}

// GetGoodsDetail 获取商品详情
func (uc *GoodsUsecase) GetGoodsDetail(ctx context.Context, id int32) (*Goods, error) {
	return uc.repo.GetGoodsDetail(ctx, id)
}

// SearchGoods 搜索商品
func (uc *GoodsUsecase) SearchGoods(ctx context.Context, keyword string, page, pageSize int32) ([]*Goods, int32, error) {
	return uc.repo.SearchGoods(ctx, keyword, page, pageSize)
}

// BatchGetGoods 批量获取商品信息
func (uc *GoodsUsecase) BatchGetGoods(ctx context.Context, ids []int32) (map[int32]*Goods, error) {
	return uc.repo.BatchGetGoods(ctx, ids)
}
