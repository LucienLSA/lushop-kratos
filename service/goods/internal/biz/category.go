package biz

import (
	"context"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type CategoryRepo interface {
	// 商品分类
	AddCategory(context.Context, *domain.CategoryInfo) (*domain.Category, error)
	UpdateCategory(context.Context, *domain.CategoryInfo) error
	CategoryList(context.Context) ([]*domain.Category, error)
	GetCategoryByID(ctx context.Context, id int32) (*domain.Category, error)
	SubCategory(context.Context, *domain.Category) ([]*domain.CategoryList, error)
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

func (c *CategoryUsecase) CreateCategory(ctx context.Context, r *domain.CategoryInfo) (*domain.Category, error) {
	cateInfo, err := c.repo.AddCategory(ctx, r)
	if err != nil {
		return nil, err
	}
	return cateInfo, nil
}

func (c *CategoryUsecase) DeleteCategory(ctx context.Context, r *domain.CategoryInfo) error {
	// 校验分类是否存在
	cate, err := c.repo.GetCategoryByID(ctx, r.ID)
	if err != nil {
		return err
	}
	// 检查是否存在子分类，若存在则不允许删除
	subs, err := c.repo.SubCategory(ctx, cate)
	if err != nil {
		return err
	}
	if len(subs) > 0 {
		return errors.BadRequest("CATEGORY_HAS_CHILDREN", "当前分类下存在子分类，无法删除")
	}
	return c.repo.DeleteCategory(ctx, r.ID)
}

func (c *CategoryUsecase) UpdateCategory(ctx context.Context, r *domain.CategoryInfo) error {
	err := c.repo.UpdateCategory(ctx, r)
	return err
}

func (c *CategoryUsecase) CategoryList(ctx context.Context) ([]*domain.Category, error) {
	return c.repo.CategoryList(ctx)
}

func (c *CategoryUsecase) SubCategoryList(ctx context.Context, cid int32) (*domain.CategoryList, error) {
	// 获取分类ID
	cateInfo, err := c.repo.GetCategoryByID(ctx, cid)
	if err != nil {
		return nil, err
	}
	// 查询子分类
	category, err := c.repo.SubCategory(ctx, cateInfo)
	if err != nil {
		return nil, err
	}

	return &domain.CategoryList{
		Category:    cateInfo,
		SubCategory: category,
	}, nil
}
