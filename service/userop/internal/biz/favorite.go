package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

// FavoriteRepo is a Favorite repo.
type FavoriteRepo interface {
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
