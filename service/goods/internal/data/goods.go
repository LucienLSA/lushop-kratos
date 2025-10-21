package data

import (
	"context"
	"goods/internal/biz"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// Goods 商品表
type Goods struct {
	BaseFields
	CategoryID      int32 `gorm:"type:int;not null;comment:'商品分类ID'"`
	Category        Category
	BrandsID        int32 `gorm:"type:int;not null"`
	Brand           Brand
	OnSale          bool     `gorm:"default:false;not null;comment:'是否特价'"`
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:'商品编号'"`
	Name            string   `gorm:"type:varchar(100);not null;comment:'商品名称'"`
	ClickNum        int32    `gorm:"type:int;default:0;not null;comment:'商品点击数'"`
	SoldNum         int32    `gorm:"type:int;default:0;not null;comment:'商品销量'"`
	FavNum          int32    `gorm:"type:int;default:0;not null;comment:'商品收藏数'"`
	MarketPrice     float32  `gorm:"not null;comment:'商品市场价'"`
	ShopPrice       float32  `gorm:"not null;comment:'商品实际价'"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null;comment:'商品简介'"`
	ShipFree        bool     `gorm:"default:false;not null;comment:'是否免运费'"`
	Images          GormList `gorm:"type:varchar(1000);not null;comment:'商品图片'"`
	DescImages      GormList `gorm:"type:varchar(5000);not null;comment:'商品详情图片'"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:'商品封面图'"`
	IsNew           bool     `gorm:"default:false;not null;comment:'是否新品'"`
	IsHot           bool     `gorm:"default:false;not null;comment:'是否热卖'"`
}

type GoodsRepo struct {
	data *Data
	log  *log.Helper
}

// NewGoodsRepo .
func NewGoodsRepo(data *Data, logger log.Logger) biz.GoodsRepo {
	return &GoodsRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// data数据层结构体转化为Domain层结构体
func (p *Goods) ToDomain() *domain.Goods {
	return &domain.Goods{
		ID:              p.ID,
		CategoryID:      p.CategoryID,
		BrandsID:        p.BrandsID,
		Name:            p.Name,
		GoodsSn:         p.GoodsSn,
		MarketPrice:     p.MarketPrice,
		GoodsBrief:      p.GoodsBrief,
		GoodsFrontImage: p.GoodsFrontImage,
		GoodsImages:     p.Images,
		DescImages:      p.DescImages,
		OnSale:          p.OnSale,
		ShipFree:        p.ShipFree,
		IsNew:           p.IsNew,
		IsHot:           p.IsHot,
	}
}

func (g GoodsRepo) CreateGoods(c context.Context, goods *domain.Goods) (*domain.Goods, error) {
	product := &Goods{
		CategoryID:      goods.CategoryID,
		BrandsID:        goods.BrandsID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		MarketPrice:     goods.MarketPrice,
		GoodsBrief:      goods.GoodsBrief,
		GoodsFrontImage: goods.GoodsFrontImage,
		Images:          goods.GoodsImages,
		OnSale:          goods.OnSale,
		ShipFree:        goods.ShipFree,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
	}

	result := g.data.DB(c).Save(product)
	if result.Error != nil {
		return nil, errors.InternalServer("GOODS_CREATE_ERROR", "商品创建失败")
	}
	return product.ToDomain(), nil
}

func (g GoodsRepo) GoodsListByIDs(c context.Context, ids ...int64) ([]*domain.Goods, error) {
	var l []*Goods
	if err := g.data.DB(c).Where("id IN (?)", ids).Find(&l).Error; err != nil {
		return nil, errors.NotFound("GOODS_NOT_FOUND", "商品不存在")
	}
	var res []*domain.Goods
	for _, item := range l {
		res = append(res, item.ToDomain())
	}
	return res, nil
}

func (g GoodsRepo) GoodsByID(c context.Context, id int64) (*domain.Goods, error) {
	var m Goods
	if err := g.data.DB(c).First(&m, id).Error; err != nil {
		return nil, errors.NotFound("GOODS_NOT_FOUND", "商品不存在")
	}
	return m.ToDomain(), nil
}

func (g GoodsRepo) UpdateGoods(c context.Context, d *domain.Goods) error {
	updates := map[string]interface{}{
		"category_id":       d.CategoryID,
		"brands_id":         d.BrandsID,
		"name":              d.Name,
		"goods_sn":          d.GoodsSn,
		"market_price":      d.MarketPrice,
		"shop_price":        d.ShopPrice,
		"goods_brief":       d.GoodsBrief,
		"ship_free":         d.ShipFree,
		"images":            d.GoodsImages,
		"desc_images":       d.DescImages,
		"goods_front_image": d.GoodsFrontImage,
		"is_new":            d.IsNew,
		"is_hot":            d.IsHot,
		"on_sale":           d.OnSale,
	}
	if err := g.data.DB(c).Model(&Goods{}).Where("id = ?", d.ID).Updates(updates).Error; err != nil {
		return errors.InternalServer("GOODS_UPDATE_ERROR", "商品更新失败")
	}
	return nil
}

func (g GoodsRepo) DeleteGoods(c context.Context, id int64) error {
	if err := g.data.DB(c).Delete(&Goods{}, id).Error; err != nil {
		return errors.InternalServer("GOODS_DELETE_ERROR", "商品删除失败")
	}
	return nil
}
