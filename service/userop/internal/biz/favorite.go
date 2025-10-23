package biz

import (
	"context"
	"userop/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

// FavoriteRepo is a Favorite repo.
type FavoriteRepo interface {
	AddUserFav(ctx context.Context, userFav domain.Favorite) error
	DeleteUserFav(ctx context.Context, userFav domain.Favorite) error
	GetUserFavDetail(ctx context.Context, userFav domain.Favorite) (*domain.Favorite, error)
	GetFavList(ctx context.Context, filter domain.Favorite) (*domain.UserFavListResponse, error)
}

// FavoriteUsecase is a Favorite usecase.
type FavoriteUsecase struct {
	repo FavoriteRepo
	log  *log.Helper
}

// NewFavoriteUsecase new a Favorite usecase.
func NewFavoriteUsecase(repo FavoriteRepo, logger log.Logger) *FavoriteUsecase {
	return &FavoriteUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (f *FavoriteUsecase) AddUserFav(ctx context.Context, userFav domain.Favorite) error {
	return f.repo.AddUserFav(ctx, userFav)
}

func (f *FavoriteUsecase) DeleteUserFav(ctx context.Context, userFav domain.Favorite) error {
	return f.repo.DeleteUserFav(ctx, userFav)
}

func (f *FavoriteUsecase) GetUserFavDetail(ctx context.Context, userFav domain.Favorite) (*domain.Favorite, error) {
	return f.repo.GetUserFavDetail(ctx, userFav)
}

func (f *FavoriteUsecase) GetFavList(ctx context.Context, filter domain.Favorite) (*domain.UserFavListResponse, error) {
	return f.repo.GetFavList(ctx, filter)
}
