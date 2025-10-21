package data

import (
	"context"
	"goods/internal/biz"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// Brand 商品品牌表
type Brand struct {
	BaseFields
	Name string `gorm:"type:varchar(50);not null;comment:'品牌名称'"`
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:'品牌Logo图片'"`
}

func (p *Brand) ToDomain() *domain.Brand {
	return &domain.Brand{
		ID:   int32(p.ID),
		Name: p.Name,
		Logo: p.Logo,
	}
}

type BrandRepo struct {
	data *Data
	log  *log.Helper
}

// NewBrandRepo .
func NewBrandRepo(data *Data, logger log.Logger) biz.BrandRepo {
	return &BrandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *BrandRepo) Create(ctx context.Context, b *domain.Brand) (*domain.Brand, error) {
	brand := &Brand{
		Name: b.Name,
		Logo: b.Logo,
	}
	if err := r.data.DB(ctx).Create(brand).Error; err != nil {
		return nil, errors.InternalServer("SAVE_BRAND_ERROR", err.Error())
	}
	return brand.ToDomain(), nil
}

func (r *BrandRepo) GetBrandByName(ctx context.Context, name string) (*domain.Brand, error) {
	var brand Brand
	result := r.data.DB(ctx).Where("name = ?", name).First(&brand)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("BRAND_NOT_FOUND", "brand not found")
		}
		return nil, errors.InternalServer("BRAND_QUERY_ERROR", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.NotFound("BRAND_NOT_FOUND", "brand not found")
	}
	return brand.ToDomain(), nil
}

func (r *BrandRepo) Update(ctx context.Context, b *domain.Brand) error {
	brands := Brand{}
	if result := r.data.DB(ctx).Where("id=?", b.ID).First(&brands); result.RowsAffected == 0 {
		return errors.NotFound("BRAND_NOT_FOUND", "brand not found")
	}

	if b.Name != "" {
		brands.Name = b.Name
	}
	if b.Logo != "" {
		brands.Logo = b.Logo
	}
	if err := r.data.DB(ctx).Save(&brands).Error; err != nil {
		return errors.InternalServer("UPDATE_BRAND_ERROR", err.Error())
	}
	return nil
}

func (r *BrandRepo) List(ctx context.Context, b *biz.Pagination) ([]*domain.Brand, int64, error) {
	var brands []Brand
	result := r.data.DB(ctx).Scopes(Paginate(b.PageNum, b.PageSize)).Find(&brands)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, 0, errors.NotFound("BRAND_NOT_FOUND", "brand not found")
	}
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var rsp []*domain.Brand
	var total int64
	result = r.data.DB(ctx).Table("brands").Model(&Brand{}).Count(&total)
	if result.Error != nil {
		return nil, 0, errors.NotFound("BRAND_NOT_FOUND", "brand not found")
	}
	for _, v := range brands {
		br := &domain.Brand{
			ID:   int32(v.ID),
			Name: v.Name,
			Logo: v.Logo,
		}
		rsp = append(rsp, br)
	}
	return rsp, total, nil
}

func (r *BrandRepo) IsBrand(ctx context.Context, ids []int32) error {
	idCount := len(ids)
	if idCount == 0 {
		return errors.InternalServer("BRAND_NOT_FOUND", "brand not found")
	}
	var count int64
	result := r.data.DB(ctx).Table("brands").Where("id IN (?)", ids).Count(&count)
	if result.Error != nil {
		return errors.InternalServer("BRAND_NOT_FOUND", result.Error.Error())
	}
	if int64(idCount) != count {
		return errors.InternalServer("BRAND_NOT_FOUND", "品牌不存在")
	}
	return nil
}

func (r *BrandRepo) GetBrandByID(ctx context.Context, id int32) (*domain.Brand, error) {
	var b Brand
	if err := r.data.DB(ctx).Table("brands").Where("id = ?", id).First(&b).Error; err != nil {
		return nil, errors.InternalServer("BRAND_NOT_FOUND", err.Error())
	}
	return b.ToDomain(), nil
}

func (r *BrandRepo) ListByIds(ctx context.Context, ids ...int32) (domain.BrandList, error) {
	if len(ids) == 0 {
		return nil, errors.InternalServer("BRAND_NOT_FOUND", "请选择品牌")
	}

	var l []*Brand
	if err := r.data.DB(ctx).Where("id IN (?)", ids).Find(&l).Error; err != nil {
		return nil, errors.InternalServer("BRAND_NOT_FOUND", err.Error())
	}

	var res domain.BrandList
	for _, item := range l {
		res = append(res, item.ToDomain())
	}
	return res, nil
}

func (r *BrandRepo) Delete(ctx context.Context, id int32) error {
	return r.data.DB(ctx).Delete(&Brand{}, id).Error
}
