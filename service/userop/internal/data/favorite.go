package data

import (
	"userop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type favoriteRepo struct {
	data *Data
	log  *log.Helper
}

// NewFavoriteRepo .
func NewFavoriteRepo(data *Data, logger log.Logger) biz.FavoriteRepo {
	return &favoriteRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
