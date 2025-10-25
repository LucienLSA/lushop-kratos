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
	User  int32 `gorm:"column:user_id;type:int;index:idx_user_goods,unique;comment:用户id"`
	Goods int32 `gorm:"column:goods_id;type:int;index:idx_user_goods,unique;comment:商品id"`
}

func (UserFav) TableName() string {
	return "user_favs"
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
	dbFav := UserFav{
		User:  userFav.UserID,
		Goods: userFav.GoodsID,
	}
	return f.data.DB(ctx).Create(&dbFav).Error
}

func (f *favoriteRepo) DeleteUserFav(ctx context.Context, userFav domain.Favorite) error {
	return f.data.DB(ctx).Where("user_id = ? AND goods_id = ?", userFav.UserID, userFav.GoodsID).Delete(&UserFav{}).Error
}

func (f *favoriteRepo) GetUserFavDetail(ctx context.Context, userFav domain.Favorite) (*domain.Favorite, error) {
	var dbFav UserFav
	result := f.data.DB(ctx).Where("user_id = ? AND goods_id = ?", userFav.UserID, userFav.GoodsID).First(&dbFav)
	if result.Error != nil {
		return nil, result.Error
	}
	return &domain.Favorite{UserID: dbFav.User, GoodsID: dbFav.Goods}, nil
}

func (f *favoriteRepo) GetFavList(ctx context.Context, filter domain.Favorite) (*domain.UserFavListResponse, error) {
	var (
		list  []UserFav
		total int64
	)
	db := f.data.DB(ctx).Model(&UserFav{})
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
		resp.List = append(resp.List, &domain.UserFavResponse{UserID: item.User, GoodsID: item.Goods})
	}
	return resp, nil
}
