package biz

import (
	"context"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

type CategoryRepo interface {
	AddCategory(context.Context, *domain.CategoryInfo) (*domain.CategoryInfo, error)
	UpdateCategory(context.Context, *domain.CategoryInfo) error
	Category(context.Context) ([]*domain.Category, error)
	GetCategoryByID(ctx context.Context, id int32) (*domain.CategoryInfo, error)
	SubCategory(context.Context, *domain.CategoryInfo) ([]*domain.CategoryInfo, error)
	DeleteCategory(context.Context, int32) error
	GetCategoryAll(context.Context, int32, int32) ([]interface{}, error)
}

type CategoryUsecase struct {
	repo CategoryRepo
	log  *log.Helper
}

func NewCategoryUsecase(repo CategoryRepo, logger log.Logger) *CategoryUsecase {
	return &CategoryUsecase{repo: repo, log: log.NewHelper(logger)}
}
func (c *CategoryUsecase) CreateCategory(ctx context.Context, r *domain.CategoryInfo) (*domain.CategoryInfo, error) {
	cateInfo, err := c.repo.AddCategory(ctx, r)
	if err != nil {
		return nil, err
	}
	return cateInfo, nil
}

func (c *CategoryUsecase) DeleteCategory(ctx context.Context, r *domain.CategoryInfo) error {
	// todo 需要验证是否是定级分类,定级分类下面还有没有二级分类
	err := c.repo.DeleteCategory(ctx, r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryUsecase) UpdateCategory(ctx context.Context, r *domain.CategoryInfo) error {
	err := c.repo.UpdateCategory(ctx, r)
	return err
}

func (c *CategoryUsecase) CategoryList(ctx context.Context) ([]*domain.Category, error) {
	return c.repo.Category(ctx)
}

func (c *CategoryUsecase) SubCategoryList(ctx context.Context, cid int32) (*domain.CategoryList, error) {
	cateInfo, err := c.repo.GetCategoryByID(ctx, cid)
	if err != nil {
		return nil, err
	}

	category, err := c.repo.SubCategory(ctx, cateInfo)
	if err != nil {
		return nil, err
	}

	return &domain.CategoryList{
		Category:    cateInfo,
		SubCategory: category,
	}, nil
}
