package domain

// GoodsInfo 商品信息（领域模型）
type GoodsInfo struct {
	ID              int32   // 商品ID
	Name            string  // 商品名称
	GoodsFrontImage string  // 商品图片
	ShopPrice       float32 // 商品价格
}
