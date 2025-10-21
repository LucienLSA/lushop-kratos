package domain

// CategoryBrand 表示分类与品牌的关联关系
type CategoryBrand struct {
    ID         int32
    CategoryID int32
    BrandsID   int32
}

// CategoryBrandList 用于批量返回分类-品牌关联
type CategoryBrandList []*CategoryBrand
