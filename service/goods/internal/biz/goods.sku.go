package biz

import (
	"context"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

type GoodsSkuRepo interface {
	Create(context.Context, *domain.GoodsSku) (*domain.GoodsSku, error)
	CreateSkuRelation(context.Context, []*domain.GoodsSpecificationSku) error
}

type GoodsSkuUsecase struct {
	repo GoodsSkuRepo
	log  *log.Helper
}

func NewGoodsSkuUsecase(repo GoodsSkuRepo, logger log.Logger) *GoodsSkuUsecase {
	return &GoodsSkuUsecase{repo: repo, log: log.NewHelper(logger)}
}
