package biz

import (
	"context"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

type CategoryBrandRepo interface {
	// 商品品牌分类
	CreateCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) (*domain.CategoryBrand, error)
	GetCategoryBrandList(ctx context.Context, categoryId int32) (*domain.CategoryBrandList, error)
	CategoryBrandList(ctx context.Context, pg *Pagination) (*domain.CategoryBrandList, int64, error)
	DeleteCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) error
	UpdateCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) error
}

type CategoryBrandUsecase struct {
	repo CategoryBrandRepo
	log  *log.Helper
}

func NewCategoryBrandUsecase(repo CategoryBrandRepo, logger log.Logger) *CategoryBrandUsecase {
	return &CategoryBrandUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (u *CategoryBrandUsecase) CreateCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) (*domain.CategoryBrand, error) {
	categoryBrand, err := u.repo.CreateCategoryBrand(ctx, categoryBrand)
	if err != nil {
		return nil, err
	}
	return categoryBrand, nil
}

func (u *CategoryBrandUsecase) GetCategoryBrandList(ctx context.Context, categoryId int32) (*domain.CategoryBrandList, error) {
	return u.repo.GetCategoryBrandList(ctx, categoryId)
}

func (u *CategoryBrandUsecase) CategoryBrandList(ctx context.Context, pg *Pagination) (*domain.CategoryBrandList, int64, error) {
	return u.repo.CategoryBrandList(ctx, pg)
}

func (u *CategoryBrandUsecase) DeleteCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) error {
	return u.repo.DeleteCategoryBrand(ctx, categoryBrand)
}

func (u *CategoryBrandUsecase) UpdateCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) error {
	return u.repo.UpdateCategoryBrand(ctx, categoryBrand)
}
