package biz

import (
	"context"
	"errors"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

type Pagination struct {
	PageNum  int
	PageSize int
}
type BrandRepo interface {
	Create(context.Context, *domain.Brand) (*domain.Brand, error)
	GetBrandByName(context.Context, string) (*domain.Brand, error)
	Update(context.Context, *domain.Brand) error
	List(context.Context, *Pagination) ([]*domain.Brand, int64, error)
	GetBrandByID(context.Context, int32) (*domain.Brand, error)
	IsBrand(context.Context, []int32) error
	ListByIds(context.Context, ...int32) (domain.BrandList, error)
	Delete(context.Context, int32) error
}
type BrandUsecase struct {
	repo BrandRepo
	log  *log.Helper
}

func NewBrandUsecase(repo BrandRepo, logger log.Logger) *BrandUsecase {
	return &BrandUsecase{repo: repo, log: log.NewHelper(logger)}
}
func (uc *BrandUsecase) CreateBrand(ctx context.Context, b *domain.Brand) (*domain.Brand, error) {
	if existing, err := uc.repo.GetBrandByName(ctx, b.Name); err == nil && existing != nil {
		return nil, errors.New("当前品牌已经存在")
	} else if err != nil {
		// 如果是未找到则允许创建，其它错误直接返回
		// 由于仓储层未定义特定错误类型，这里直接尝试创建，若底层约束冲突会返回错误
		return uc.repo.Create(ctx, b)
	}
	return nil, errors.New("品牌创建失败")
}
func (uc *BrandUsecase) UpdateBrand(ctx context.Context, b *domain.Brand) error {
	if existing, err := uc.repo.GetBrandByName(ctx, b.Name); err != nil || existing == nil {
		return errors.New("当前品牌不存在")
	}
	return uc.repo.Update(ctx, b)
}

func (uc *BrandUsecase) BrandList(ctx context.Context, b *Pagination) ([]*domain.Brand, int64, error) {
	list, total, err := uc.repo.List(ctx, b)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil

}
func (uc *BrandUsecase) ListByIds(ctx context.Context, ids ...int32) (domain.BrandList, error) {
	return uc.repo.ListByIds(ctx, ids...)
}
func (uc *BrandUsecase) DeleteBrand(ctx context.Context, id int32) error {
	return uc.repo.Delete(ctx, id)
}
