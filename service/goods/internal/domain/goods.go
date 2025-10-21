package domain

type Goods struct {
	ID              int64
	CategoryID      int32
	BrandsID        int32
	Category        Category
	Brand           Brand
	OnSale          bool
	GoodsSn         string
	Name            string
	ClickNum        int64
	SoldNum         int64
	FavNum          int64
	MarketPrice     float32
	ShopPrice       float32
	GoodsBrief      string
	ShipFree        bool
	GoodsFrontImage string
	DescImages      []string
	GoodsImages     []string
	IsNew           bool
	IsHot           bool
}

type GoodsInfoResponse struct {
	GoodsID int64
}

type GoodsListResponse struct {
	Total int64
	List  []*Goods
}
