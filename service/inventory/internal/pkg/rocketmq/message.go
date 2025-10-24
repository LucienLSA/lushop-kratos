package rocketmq

// OrderInventoryMessage 订单库存扣减消息
type OrderInventoryMessage struct {
	OrderSn   string          `json:"order_sn"`   // 订单号
	UserID    int32           `json:"user_id"`    // 用户ID
	GoodsInfo []GoodsInvInfo  `json:"goods_info"` // 商品信息列表
}

// GoodsInvInfo 商品库存信息
type GoodsInvInfo struct {
	GoodsID int32 `json:"goods_id"` // 商品ID
	Nums    int32 `json:"nums"`     // 数量
}
