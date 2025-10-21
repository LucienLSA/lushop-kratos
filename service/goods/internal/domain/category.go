package domain

type Category struct {
	ID               int32
	Name             string
	ParentCategoryID int32
	SubCategory      []*Category
	Level            int32
	IsTab            bool
}

type CategoryList struct {
	Category    *Category
	SubCategory []*CategoryList
}

type CategoryInfo struct {
	ID             int32
	Name           string
	ParentCategory int32
	Level          int32
	IsTab          bool
}
