package data

import (
	"context"
	"goods/internal/biz"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// GoodsCategoryBrand  商品和分类多对多的表
type GoodsCategoryBrand struct {
	BaseFields
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique;comment:商品和分类联合索引唯一"`
	BrandsID   int32 `gorm:"type:int;index:idx_category_brand,unique:comment:商品和分类联合索引唯一"`
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type CategoryBrandRepo struct {
	data *Data
	log  *log.Helper
}

// NewCategoryBrandRepo .
func NewCategoryBrandRepo(data *Data, logger log.Logger) biz.CategoryBrandRepo {
	return &CategoryBrandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (c *CategoryBrandRepo) CreateCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) (*domain.CategoryBrand, error) {
	m := &GoodsCategoryBrand{
		CategoryID: categoryBrand.CategoryID,
		BrandsID:   categoryBrand.BrandsID,
	}
	if err := c.data.DB(ctx).Create(m).Error; err != nil {
		return nil, errors.InternalServer("CATEGORY_BRAND_CREATE_ERROR", err.Error())
	}
	return &domain.CategoryBrand{
		ID:         int32(m.ID),
		CategoryID: m.CategoryID,
		BrandsID:   m.BrandsID,
	}, nil
}

func (c *CategoryBrandRepo) GetCategoryBrandList(ctx context.Context, categoryId int32) (*domain.CategoryBrandList, error) {
	var rows []GoodsCategoryBrand
	db := c.data.DB(ctx)
	if err := db.Preload("Brand").Where("category_id = ?", categoryId).Find(&rows).Error; err != nil {
		return nil, errors.InternalServer("CATEGORY_BRAND_QUERY_ERROR", err.Error())
	}
	var res domain.CategoryBrandList
	for i := range rows {
		res = append(res, &domain.CategoryBrand{
			ID:         int32(rows[i].ID),
			CategoryID: rows[i].CategoryID,
			BrandsID:   rows[i].BrandsID,
		})
	}
	return &res, nil
}

func (c *CategoryBrandRepo) CategoryBrandList(ctx context.Context, pg *biz.Pagination) (*domain.CategoryBrandList, int64, error) {
	db := c.data.DB(ctx)
	var total int64
	if err := db.Model(&GoodsCategoryBrand{}).Count(&total).Error; err != nil {
		return nil, 0, errors.InternalServer("CATEGORY_BRAND_COUNT_ERROR", err.Error())
	}
	var rows []GoodsCategoryBrand
	if err := db.Model(&GoodsCategoryBrand{}).
		Scopes(Paginate(pg.PageNum, pg.PageSize)).
		Find(&rows).Error; err != nil {
		return nil, 0, errors.InternalServer("CATEGORY_BRAND_LIST_ERROR", err.Error())
	}
	var res domain.CategoryBrandList
	for i := range rows {
		res = append(res, &domain.CategoryBrand{
			ID:         int32(rows[i].ID),
			CategoryID: rows[i].CategoryID,
			BrandsID:   rows[i].BrandsID,
		})
	}
	return &res, total, nil
}

func (c *CategoryBrandRepo) DeleteCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) error {
	return c.data.DB(ctx).Delete(&GoodsCategoryBrand{}, categoryBrand.ID).Error
}

func (c *CategoryBrandRepo) UpdateCategoryBrand(ctx context.Context, categoryBrand *domain.CategoryBrand) error {
	m := &GoodsCategoryBrand{
		CategoryID: categoryBrand.CategoryID,
		BrandsID:   categoryBrand.BrandsID,
	}
	return c.data.DB(ctx).Save(m).Error
}
