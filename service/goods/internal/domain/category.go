package domain

type Category struct {
	ID               int32
	Name             string
	ParentCategoryID int32
	SubCategory      []*Category
	Level            int32
	IsTab            bool
	Sort             int32
}

type CategoryList struct {
	Category    *CategoryInfo
	SubCategory []*CategoryInfo
}

type CategoryInfo struct {
	ID             int32
	Name           string
	ParentCategory int32
	Level          int32
	IsTab          bool
	Sort           int32
}
