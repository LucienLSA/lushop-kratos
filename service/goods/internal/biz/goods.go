package biz

import (
	"context"
	"errors"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

// GoodsRepo is a Goods repo.
type GoodsRepo interface {
	CreateGoods(ctx context.Context, goods *domain.Goods) (*domain.Goods, error)
	GoodsListByIDs(context.Context, ...int64) ([]*domain.Goods, error)
	GoodsByID(ctx context.Context, id int64) (*domain.Goods, error)
	UpdateGoods(ctx context.Context, goods *domain.Goods) error
	DeleteGoods(ctx context.Context, id int64) error
}

// GoodsUsecase is a Goods usecase.
type GoodsUsecase struct {
	repo         GoodsRepo
	tr           Transaction
	brandRepo    BrandRepo
	categoryRepo CategoryRepo
	esGoodsRepo  EsGoodsRepo
	log          *log.Helper
}

// NewGoodsUsecase new a Goods usecase.
func NewGoodsUsecase(repo GoodsRepo, bRepo BrandRepo, cRepo CategoryRepo, tx Transaction, es EsGoodsRepo, logger log.Logger) *GoodsUsecase {
	return &GoodsUsecase{
		repo:         repo,
		log:          log.NewHelper(logger),
		brandRepo:    bRepo,
		categoryRepo: cRepo,
		tr:           tx,
		esGoodsRepo:  es,
	}
}

func (g GoodsUsecase) CreateGoods(ctx context.Context, r *domain.Goods) (*domain.GoodsInfoResponse, error) {
	var (
		err     error
		goods   *domain.Goods
		esGoods *domain.ESGoods
	)
	// 判断商品品牌是否存在
	brand, err := g.brandRepo.GetBrandByID(ctx, r.BrandsID)
	if err != nil {
		return nil, errors.New("品牌不存在")
	}
	// 判断分类是否存在
	category, err := g.categoryRepo.GetCategoryByID(ctx, r.CategoryID)
	if err != nil {
		return nil, errors.New("商品分类不存在")
	}
	// 判断商品是否已存在（根据传入的 ID）
	if r.ID > 0 {
		if _, err := g.repo.GoodsByID(ctx, r.ID); err == nil {
			return nil, errors.New("商品已存在")
		}
	}
	// 开启事务
	err = g.tr.ExecTx(ctx, func(ctx context.Context) error {
		// 更新商品表
		goods, err = g.repo.CreateGoods(ctx, &domain.Goods{
			CategoryID:      r.CategoryID,
			BrandsID:        r.BrandsID,
			Name:            r.Name,
			GoodsSn:         r.GoodsSn,
			MarketPrice:     r.MarketPrice,
			ShopPrice:       r.ShopPrice,
			GoodsBrief:      r.GoodsBrief,
			GoodsFrontImage: r.GoodsFrontImage,
			GoodsImages:     r.GoodsImages,
			OnSale:          r.OnSale,
			IsNew:           r.IsNew,
			IsHot:           r.IsHot,
			ShipFree:        r.ShipFree,
		})
		if err != nil {
			return err
		}
		// 构建并插入 ES 商品文档（最小必要字段）
		esGoods = &domain.ESGoods{
			ID:           goods.ID,
			CategoryID:   category.ID,
			CategoryName: category.Name,
			BrandsID:     brand.ID,
			BrandName:    brand.Name,
			OnSale:       goods.OnSale,
			ShipFree:     goods.ShipFree,
			IsNew:        goods.IsNew,
			IsHot:        goods.IsHot,
			Name:         goods.Name,
			ClickNum:     goods.ClickNum,
			SoldNum:      goods.SoldNum,
			FavNum:       goods.FavNum,
			MarketPrice:  int64(goods.MarketPrice),
			GoodsBrief:   goods.GoodsBrief,
		}
		// 插入 EsGoods
		err = g.esGoodsRepo.InsertEsGoods(ctx, esGoods)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &domain.GoodsInfoResponse{GoodsID: goods.ID}, nil
}

func (g GoodsUsecase) UpdateGoods(ctx context.Context, r *domain.Goods) (*domain.GoodsInfoResponse, error) {
	if r.BrandsID == 0 && r.CategoryID == 0 {
		goods, err := g.repo.GoodsByID(ctx, r.ID)
		if err != nil {
			return nil, errors.New("商品不存在")
		}
		goods.IsNew = r.IsNew
		goods.IsHot = r.IsHot
		goods.OnSale = r.OnSale
		if err := g.repo.UpdateGoods(ctx, goods); err != nil {
			return nil, err
		}
		return &domain.GoodsInfoResponse{GoodsID: goods.ID}, nil
	}
	// 如果要更新商品的分类或品牌，先检查数据库中是否存在对应的分类和品牌。
	brand, err := g.brandRepo.GetBrandByID(ctx, r.BrandsID)
	if err != nil {
		return nil, errors.New("品牌不存在")
	}
	category, err := g.categoryRepo.GetCategoryByID(ctx, r.CategoryID)
	if err != nil {
		return nil, errors.New("商品分类不存在")
	}
	goods, err := g.repo.GoodsByID(ctx, r.ID)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	goods.Brand = *brand
	goods.Category = *category
	goods.BrandsID = r.BrandsID
	goods.CategoryID = r.CategoryID
	goods.Name = r.Name
	goods.GoodsSn = r.GoodsSn
	goods.MarketPrice = r.MarketPrice
	goods.ShipFree = r.ShipFree
	goods.ShopPrice = r.ShopPrice
	goods.GoodsImages = r.GoodsImages
	goods.DescImages = r.DescImages
	goods.GoodsFrontImage = r.GoodsFrontImage
	goods.IsNew = r.IsNew
	goods.IsHot = r.IsHot
	goods.OnSale = r.OnSale
	err = g.repo.UpdateGoods(ctx, goods)
	if err != nil {
		return nil, err
	}
	return &domain.GoodsInfoResponse{GoodsID: goods.ID}, nil
}

func (g GoodsUsecase) BatchGetGoods(ctx context.Context, ids []int64) (*domain.GoodsListResponse, error) {
	goodsList, err := g.repo.GoodsListByIDs(ctx, ids...)
	if err != nil {
		return nil, err
	}
	return &domain.GoodsListResponse{
		Total: int64(len(goodsList)),
		List:  goodsList,
	}, nil
}

func (g GoodsUsecase) DeleteGoods(ctx context.Context, id int64) error {
	return g.repo.DeleteGoods(ctx, id)
}

func (g GoodsUsecase) GetGoodsById(ctx context.Context, id int64) (*domain.Goods, error) {
	return g.repo.GoodsByID(ctx, id)
}
