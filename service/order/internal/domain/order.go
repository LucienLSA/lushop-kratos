package domain

import "time"

// ShoppingCart is the domain model for a user's cart item.
type ShoppingCart struct {
    ID      int32
    UserID  int32
    GoodsID int32
    Nums    int32
    Checked bool
}

// OrderInfo is the domain model for an order header.
type OrderInfo struct {
    ID          int32
    UserID      int32
    OrderSn     string
    PayType     string
    Status      string
    TradeNo     string
    OrderMount  float32
    PayTime     *time.Time
    Address     string
    SignerName  string
    SingerMobile string
    Post        string
}

// OrderGoods is the domain model for items under an order.
type OrderGoods struct {
    ID         int32
    OrderID    int32
    GoodsID    int32
    GoodsName  string
    GoodsImage string
    GoodsPrice float32
    Nums       int32
}
