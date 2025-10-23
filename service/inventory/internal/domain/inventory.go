package domain

type Inventory struct {
	Goods  int32
	Stocks int32
}
type Delivery struct {
	Goods   int32
	Nums    int32
	OrderSn string
	Status  string // 1.代表等待支付，2.代表支付成功，3.支付失败
}

type GoodsInvInfo struct {
	GoodsID int32
	Nums    int32
}
type SellInfo struct {
	GoodsInvInfo []GoodsInvInfo
	OrderSn      string
}
