package data

import (
	"context"
	"userop/internal/biz"
	"userop/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

type UserFav struct {
	BaseFields
	//联合索引
	User  int32 `gorm:"type:int;index:idx_user_goods,unique;comment:用户id"`
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique;comment:商品id"`
}

func (UserFav) TableName() string {
	return "userfav"
}

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

func (f *favoriteRepo) AddUserFav(ctx context.Context, userFav domain.Favorite) error {
	return f.data.DB(ctx).Create(&userFav).Error
}

func (f *favoriteRepo) DeleteUserFav(ctx context.Context, userFav domain.Favorite) error {
	return f.data.DB(ctx).Delete(&userFav).Error
}

func (f *favoriteRepo) GetUserFavDetail(ctx context.Context, userFav domain.Favorite) (*domain.Favorite, error) {
	result := f.data.DB(ctx).Where(&userFav).First(&userFav)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userFav, nil
}

func (f *favoriteRepo) GetFavList(ctx context.Context, filter domain.Favorite) (*domain.UserFavListResponse, error) {
	var (
		list  []domain.Favorite
		total int64
	)
	db := f.data.DB(ctx).Model(&domain.Favorite{})
	if filter.UserID != 0 {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.GoodsID != 0 {
		db = db.Where("goods_id = ?", filter.GoodsID)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	resp := &domain.UserFavListResponse{Total: total}
	for _, item := range list {
		resp.List = append(resp.List, &domain.UserFavResponse{UserID: item.UserID, GoodsID: item.GoodsID})
	}
	return resp, nil
}
