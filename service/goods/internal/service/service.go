package service

import (
	v1 "goods/api/goods/v1"
	"goods/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGoodsService)

// GreeterService is a greeter service.
type GoodsService struct {
	v1.UnimplementedGoodsServer
	cac     *biz.CategoryUsecase
	bc      *biz.BrandUsecase
	g       *biz.GoodsUsecase
	esGoods *biz.EsGoodsUsecase
	cacb    *biz.CategoryBrandUsecase
	ban     *biz.BannerUsecase
	log     *log.Helper
}

// NewGoodsService new a goods service.
func NewGoodsService(bc *biz.BrandUsecase, cac *biz.CategoryUsecase, gc *biz.GoodsUsecase,
	esGoods *biz.EsGoodsUsecase, cacb *biz.CategoryBrandUsecase,
	ban *biz.BannerUsecase, logger log.Logger) *GoodsService {
	return &GoodsService{
		bc:      bc,
		cac:     cac,
		g:       gc,
		esGoods: esGoods,
		cacb:    cacb,
		ban:     ban,
		log:     log.NewHelper(logger),
	}
}
