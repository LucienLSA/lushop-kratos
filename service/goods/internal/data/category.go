package data

import (
	"context"
	"fmt"
	"goods/internal/biz"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// Category 商品分类表
type Category struct {
	BaseFields
	Name             string      `gorm:"type:varchar(20);not null;comment:'商品分类名称'" json:"name"`
	ParentCategoryID int32       `json:"parent_category_id"`
	Level            int32       `gorm:"type:int;not null;default:1;comment:'1表示商品分类的等级'" json:"level"`
	IsTab            bool        `gorm:"default:false;not null;comment:'是否Tap栏显示'" json:"is_tab"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
}

// 强制指定表名为 "category"
func (Category) TableName() string {
	return "category"
}

type CategoryRepo struct {
	data *Data
	log  *log.Helper
}

// NewCategoryRepo .
func NewCategoryRepo(data *Data, logger log.Logger) biz.CategoryRepo {
	return &CategoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, id int32) error {
	if res := r.data.DB(ctx).Delete(&Category{}, id); res.RowsAffected == 0 {
		return errors.InternalServer("DELETE_CATGORY_ERROR", res.Error.Error())
	}
	return nil
}

func (r *CategoryRepo) UpdateCategory(ctx context.Context, req *domain.CategoryInfo) error {
	var category Category
	if result := r.data.DB(ctx).First(&category, req.ID); result.RowsAffected == 0 {
		return errors.NotFound("CATEGORY_NOT_FOUND", "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}
	result := r.data.DB(ctx).Save(&category)
	if result.Error != nil {
		return errors.InternalServer("CATEGORY_UPDATE_ERROR", "商品分类创建失败")
	}
	return nil
}

func (r *CategoryRepo) AddCategory(ctx context.Context, req *domain.CategoryInfo) (*domain.Category, error) {
	category := &Category{
		Name:  req.Name,
		Level: req.Level,
		IsTab: req.IsTab,
	}
	if req.Level != 1 {
		var parent Category
		if res := r.data.DB(ctx).First(&parent, req.ParentCategory); res.RowsAffected == 0 {
			return nil, errors.NotFound("CATEGORY_NOT_FOUND", "商品父分类不存在")
		}
		category.ParentCategoryID = req.ParentCategory
	}
	if err := r.data.DB(ctx).Create(category).Error; err != nil {
		return nil, errors.InternalServer("CATEGORY_CREATE_ERROR", err.Error())
	}
	return &domain.Category{
		ID:               int32(category.ID),
		Name:             category.Name,
		ParentCategoryID: category.ParentCategoryID,
		Level:            category.Level,
		IsTab:            category.IsTab,
	}, nil
}

func (r *CategoryRepo) GetCategoryByID(ctx context.Context, id int32) (*domain.Category, error) {
	var categories Category
	if res := r.data.DB(ctx).First(&categories, id); res.RowsAffected == 0 {
		return nil, errors.NotFound("CATEGORY_NOT_FOUND", "分类不存在")
	}

	info := &domain.Category{
		ID:               int32(categories.ID),
		Name:             categories.Name,
		ParentCategoryID: categories.ParentCategoryID,
		Level:            categories.Level,
		IsTab:            categories.IsTab,
	}
	return info, nil
}

func (r *CategoryRepo) CategoryList(ctx context.Context) ([]*domain.Category, error) {
	var cate []*Category
	result := r.data.DB(ctx).Where(&Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&cate)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("CATEGORY_NOT_FOUND", "商品分类不存在")
	}
	if result.Error != nil {
		return nil, errors.InternalServer("CATEGORY_NOT_FOUND", result.Error.Error())
	}

	var res []*domain.Category
	err := copier.Copy(&res, &cate)
	if err != nil {
		return nil, errors.InternalServer("CATEGORY_COPY_ERROR", err.Error())
	}
	return res, nil
}

func (r *CategoryRepo) SubCategory(ctx context.Context, req *domain.Category) ([]*domain.CategoryList, error) {
	var subCategory []Category
	preload := "SubCategory"
	if req.Level == 1 {
		preload = "SubCategory.SubCategory"
	}
	if err := r.data.DB(ctx).Where(&Category{ParentCategoryID: req.ID}).Preload(preload).Find(&subCategory).Error; err != nil {
		return nil, errors.NotFound("CATEGORY_NOT_FOUND", "分类不存在")
	}
	var build func(c *Category) *domain.CategoryList
	build = func(c *Category) *domain.CategoryList {
		dc := &domain.Category{
			ID:               int32(c.ID),
			Name:             c.Name,
			ParentCategoryID: c.ParentCategoryID,
			Level:            c.Level,
			IsTab:            c.IsTab,
		}
		var subs []*domain.CategoryList
		for _, sc := range c.SubCategory {
			subs = append(subs, build(sc))
		}
		return &domain.CategoryList{Category: dc, SubCategory: subs}
	}
	var res []*domain.CategoryList
	for i := range subCategory {
		res = append(res, build(&subCategory[i]))
	}
	return res, nil
}

func (r *CategoryRepo) GetCategoryAll(ctx context.Context, level, id int32) ([]interface{}, error) {
	categoryIds := make([]interface{}, 0)
	var subQuery string
	// 把一级级分类下的所有三级分类都拿到
	if level == 1 {
		subQuery = fmt.Sprintf("SELECT id FROM categories WHERE parent_category_id IN (SELECT id FROM categories WHERE parent_category_id=%d)", id)
	} else if level == 2 {
		subQuery = fmt.Sprintf("SELECT id FROM categories WHERE parent_category_id=%d", id)
	} else if level == 3 {
		subQuery = fmt.Sprintf("SELECT id FROM categories WHERE id=%d", id)
	}

	type Result struct {
		ID int32
	}

	var results []Result
	if err := r.data.DB(ctx).Table("categories").Model(Category{}).Raw(subQuery).Scan(&results).Error; err != nil {
		return nil, errors.InternalServer("CATEGORY_ERROR", err.Error())
	}
	for _, re := range results {
		categoryIds = append(categoryIds, re.ID)
	}
	return categoryIds, nil
}
