package biz

import (
	"context"
	"errors"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

type BannerRepo interface {
	Create(ctx context.Context, banner *domain.Banner) (*domain.Banner, error)
	GetBannerByID(ctx context.Context, id int32) (*domain.Banner, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, banner *domain.Banner) error
	BannerList(ctx context.Context) ([]*domain.Banner, int64, error)
}
type BannerUsecase struct {
	repo BannerRepo
	log  *log.Helper
}

func NewBannerUsecase(repo BannerRepo, logger log.Logger) *BannerUsecase {
	return &BannerUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *BannerUsecase) CreateBanner(ctx context.Context, req *domain.Banner) (*domain.Banner, error) {
	if existing, err := uc.repo.GetBannerByID(ctx, req.ID); err == nil && existing != nil {
		return nil, errors.New("当前轮播图已经存在")
	} else if err != nil {
		return uc.repo.Create(ctx, req)
	}
	return nil, errors.New("轮播图创建失败")
}

func (uc *BannerUsecase) BannerList(ctx context.Context) ([]*domain.Banner, int64, error) {
	return uc.repo.BannerList(ctx)
}

func (uc *BannerUsecase) UpdateBanner(ctx context.Context, req *domain.Banner) error {
	return uc.repo.Update(ctx, req)
}

func (uc *BannerUsecase) DeleteBanner(ctx context.Context, id int32) error {
	return uc.repo.Delete(ctx, id)
}
