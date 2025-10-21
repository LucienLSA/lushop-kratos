package data

import (
	"context"
	"goods/internal/biz"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Banner struct {
	BaseFields
	Image string `gorm:"type:varchar(200);not null;comment:轮播图"`
	Url   string `gorm:"type:varchar(200);not null;comment:'图片链接'"`
	Index int32  `gorm:"type:int;default:1;not null;comment:'轮播图的索引'"`
}

func (Banner) TableName() string {
	return "banner"
}

type BannerRepo struct {
	data *Data
	log  *log.Helper
}

func (b *Banner) ToDomain() *domain.Banner {
	return &domain.Banner{
		ID:    int32(b.ID),
		Image: b.Image,
		Url:   b.Url,
		Index: b.Index,
	}
}

func NewBannerRepo(data *Data, logger log.Logger) biz.BannerRepo {
	return &BannerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (b *BannerRepo) GetBannerByID(ctx context.Context, id int32) (*domain.Banner, error) {
	var banner Banner
	if err := b.data.DB(ctx).Where("id = ?", id).First(&banner).Error; err != nil {
		return nil, errors.InternalServer("BANNER_NOT_FOUND", err.Error())
	}
	return banner.ToDomain(), nil
}

func (b *BannerRepo) Create(ctx context.Context, r *domain.Banner) (*domain.Banner, error) {
	banner := &Banner{
		Image: r.Image,
		Url:   r.Url,
		Index: r.Index,
	}
	if err := b.data.DB(ctx).Create(banner).Error; err != nil {
		return nil, errors.InternalServer("SAVE_BANNER_ERROR", err.Error())
	}
	return banner.ToDomain(), nil
}

func (b *BannerRepo) Delete(ctx context.Context, id int32) error {
	return b.data.DB(ctx).Delete(&Banner{}, id).Error
}

func (b *BannerRepo) Update(ctx context.Context, r *domain.Banner) error {
	banner := &Banner{
		Image: r.Image,
		Url:   r.Url,
		Index: r.Index,
	}
	return b.data.DB(ctx).Save(banner).Error
}

func (b *BannerRepo) BannerList(ctx context.Context) ([]*domain.Banner, int64, error) {
	var banners []*Banner
	result := b.data.DB(ctx).Find(&banners)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, 0, errors.NotFound("BANNER_NOT_FOUND", "轮播图不存在")
	}
	if result.Error != nil {
		return nil, 0, errors.InternalServer("BANNER_NOT_FOUND", result.Error.Error())
	}
	resp := make([]*domain.Banner, 0, len(banners))
	for _, banner := range banners {
		resp = append(resp, banner.ToDomain())
	}
	return resp, result.RowsAffected, nil
}
