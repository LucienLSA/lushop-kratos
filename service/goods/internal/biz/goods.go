package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

// GoodsRepo is a Goods repo.
type GoodsRepo interface {
	CreateGoods(ctx context.Context, goods *domain.Goods) (*domain.Goods, error)
	GoodsListByIDs(context.Context, ...int64) ([]*domain.Goods, error)
}

// GoodsUsecase is a Goods usecase.
type GoodsUsecase struct {
	repo              GoodsRepo
	tr                Transaction
	skuRepo           GoodsSkuRepo
	categoryRepo      CategoryRepo
	brandRepo         BrandRepo
	typeRepo          GoodsTypeRepo
	specificationRepo SpecificationRepo
	goodsAttrRepo     GoodsAttrRepo
	inventoryRepo     InventoryRepo
	esGoodsRepo       EsGoodsRepo
	log               *log.Helper
}

// NewGoodsUsecase new a Goods usecase.
func NewGoodsUsecase(repo GoodsRepo, skuRepo GoodsSkuRepo, tx Transaction,
	gRepo GoodsTypeRepo, cRepo CategoryRepo, bRepo BrandRepo,
	sRepo SpecificationRepo, aRepo GoodsAttrRepo, es EsGoodsRepo,
	iRepo InventoryRepo, logger log.Logger) *GoodsUsecase {
	return &GoodsUsecase{
		repo:              repo,
		log:               log.NewHelper(logger),
		skuRepo:           skuRepo,
		tr:                tx,
		typeRepo:          gRepo,
		categoryRepo:      cRepo,
		brandRepo:         bRepo,
		specificationRepo: sRepo,
		goodsAttrRepo:     aRepo,
		esGoodsRepo:       es, // 新增的 es 的repo
		inventoryRepo:     iRepo,
	}
}

func (g GoodsUsecase) CreateGoods(ctx context.Context, r *domain.Goods) (*domain.GoodsInfoResponse, error) {
	var (
		err     error
		goods   *domain.Goods
		esGoods *domain.ESGoods
	)
	// 判断商品品牌是否存在
	brand, err := g.brandRepo.IsBrandByID(ctx, r.BrandsID)
	if err != nil {
		return nil, errors.New("品牌不存在")
	}
	// 判断分类是否存在
	category, err := g.categoryRepo.GetCategoryByID(ctx, r.CategoryID)
	if err != nil {
		return nil, errors.New("商品分类不存在")
	}
	// 判断商品类型是否存在
	goodsType, err := g.typeRepo.IsExistsByID(ctx, r.TypeID)
	if err != nil {
		return nil, errors.New("商品类型不存在")
	}
	// 判断商品规格和属性是否存在
	for _, sku := range r.Sku {
		var sIDs []*int64
		for _, info := range sku.Specification {
			sIDs = append(sIDs, &info.SpecificationID)
		}
		specList, err := g.specificationRepo.ListByIds(ctx, sIDs...)
		if err != nil {
			return nil, err
		}
		for _, sId := range sIDs {
			info := specList.FindById(*sId)
			if info == nil {
				return nil, errors.New("商品规格不存在")
			}
		}
		var attrIDs []int64
		for _, attr := range sku.GroupAttr {
			for _, id := range attr.Attr {
				attrIDs = append(attrIDs, id.AttrID)
			}
		}
		attrList, err := g.goodsAttrRepo.ListByIds(ctx, attrIDs...)
		if err != nil {
			return nil, err
		}
		for _, attr := range sku.GroupAttr {
			for _, id := range attr.Attr {
				attrIDs = append(attrIDs, id.AttrID)
				exist := attrList.IsNotExist(attr.GroupId, id.AttrID)
				if exist {
					return nil, errors.New("商品属性不存在")
				}
			}
		}
	}
	err = g.tr.ExecTx(ctx, func(ctx context.Context) error {
		// 更新商品表
		goods, err = g.repo.CreateGoods(ctx, &domain.Goods{
			CategoryID:      r.CategoryID,
			BrandsID:        r.BrandsID,
			TypeID:          r.TypeID,
			Name:            r.Name,
			NameAlias:       r.NameAlias,
			GoodsSn:         r.GoodsSn,
			GoodsTags:       r.GoodsTags,
			MarketPrice:     r.MarketPrice,
			GoodsBrief:      r.GoodsBrief,
			GoodsFrontImage: r.GoodsFrontImage,
			GoodsImages:     r.GoodsImages,
			OnSale:          r.OnSale,
			IsNew:           r.IsNew,
			IsHot:           r.IsHot,
			ShipFree:        r.ShipFree,
			ShipID:          r.ShipID,
		})
		if err != nil {
			return err
		}
		// 更新商品 SKU 表
		for _, v := range r.Sku {
			res := &domain.GoodsSku{
				GoodsID:        goods.ID,
				GoodsSn:        goods.GoodsSn,
				GoodsName:      goods.Name,
				SkuName:        v.SkuName,
				SkuCode:        v.SkuCode,
				BarCode:        v.BarCode,
				Price:          v.Price,
				PromotionPrice: v.PromotionPrice,
				Points:         v.Points,
				RemarksInfo:    v.RemarksInfo,
				Pic:            v.Pic,
				Inventory:      v.Inventory,
				OnSale:         v.OnSale,
			}
			goodsAttr, err := json.Marshal(v.GroupAttr)
			if err != nil {
				return err
			}
			res.AttrInfo = string(goodsAttr)

			// 插入 sku 表
			skuInfo, err := g.skuRepo.Create(ctx, res)
			if err != nil {
				return err
			}

			// 插入库存表
			_, err = g.inventoryRepo.Create(ctx, &domain.Inventory{
				SkuID:     skuInfo.ID,
				Inventory: skuInfo.Inventory,
			})
			if err != nil {
				return err
			}
			// 插入 sku 规格关联关系表
			var skuRelation []*domain.GoodsSpecificationSku
			for _, spec := range v.Specification {
				skuRelation = append(skuRelation, &domain.GoodsSpecificationSku{
					SkuID:           skuInfo.ID,
					SkuCode:         skuInfo.SkuCode,
					SpecificationId: spec.SpecificationID,
					ValueId:         spec.SpecificationValueID,
				})
			}

			// 插入商品规格关联关系表
			err = g.skuRepo.CreateSkuRelation(ctx, skuRelation)
			if err != nil {
				return err
			}
			// esModel
			{
				esGoods.Sku = append(esGoods.Sku, domain.EsSku{
					SkuID:    skuInfo.ID,
					SkuName:  skuInfo.SkuName,
					SkuPrice: skuInfo.Price,
				})
				esGoods.BrandsID = brand.ID
				esGoods.BrandName = brand.Name
				esGoods.CategoryID = category.ID
				esGoods.CategoryName = category.Name
				esGoods.TypeID = goodsType.ID
				esGoods.TypeName = goodsType.Name
				esGoods.Name = goodsType.Name
				esGoods.ID = goods.ID
				esGoods.OnSale = goods.OnSale
				esGoods.ShipFree = goods.ShipFree
				esGoods.IsNew = goods.IsNew
				esGoods.IsHot = goods.IsHot
				esGoods.Name = goods.Name
				esGoods.GoodsTags = goods.GoodsTags
				esGoods.ClickNum = goods.ClickNum
				esGoods.SoldNum = goods.SoldNum
				esGoods.FavNum = goods.FavNum
				esGoods.MarketPrice = goods.MarketPrice
				esGoods.GoodsBrief = goods.GoodsBrief
			}
			fmt.Println("EsGoods", esGoods)
			// 插入 EsGoods
			err = g.esGoodsRepo.InsertEsGoods(ctx, esGoods)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &domain.GoodsInfoResponse{GoodsID: goods.ID}, nil
}
