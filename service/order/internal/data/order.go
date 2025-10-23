package data

import (
	"context"
	v1 "order/api/order/v1"
	"order/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type ShoppingCart struct {
	BaseModel
	User    int32 `gorm:"type:int;index;comment:用户ID"` //在购物车列表中我们需要查询当前用户的购物车记录
	Goods   int32 `gorm:"type:int;index;comment:商品ID"` //加索引：我们需要查询时候， 1. 会影响插入性能 2. 会占用磁盘
	Nums    int32 `gorm:"type:int;comment:数量"`
	Checked bool  `gorm:"comment:是否选中"` //是否选中
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

type OrderInfo struct {
	BaseModel

	User    int32  `gorm:"type:int;index;comment:用户ID"`
	OrderSn string `gorm:"type:varchar(30);index; comment:订单号"` //订单号，我们平台自己生成的订单号
	PayType string `gorm:"type:varchar(20);comment:alipay(支付宝), wechat(微信)"`

	//status大家可以考虑使用iota来做
	Status     string     `gorm:"type:varchar(20); comment:PAYING(待支付), TRADE_SUCCESS(成功), TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)"`
	TradeNo    string     `gorm:"type:varchar(100); comment:交易号"` //交易号就是支付宝的订单号 查账
	OrderMount float32    `gorm:"comment:总金额"`
	PayTime    *time.Time `gorm:"type:datetime; comment:支付时间"`

	Address      string `gorm:"type:varchar(100); comment:收货地址"`
	SignerName   string `gorm:"type:varchar(20); comment:收货名"`
	SingerMobile string `gorm:"type:varchar(11); comment:收获手机"`
	Post         string `gorm:"type:varchar(20); comment:留言备注"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

type OrderGoods struct {
	BaseModel

	Order int32 `gorm:"type:int;index;comment:订单ID"`
	Goods int32 `gorm:"type:int;index;comment:商品ID"`

	//把商品的信息保存下来了 ， 字段冗余， 高并发系统中我们一般都不会遵循三范式  做镜像 记录
	GoodsName  string  `gorm:"type:varchar(100);index;comment:商品名称"`
	GoodsImage string  `gorm:"type:varchar(200);comment:商品图片"`
	GoodsPrice float32 `gorm:"comment:交易时的商品价格(不是最新的价格)"`
	Nums       int32   `gorm:"type:int;comment:数量"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}

type orderRepo struct {
	data *Data
	log  *log.Helper
}

// NewOrderRepo .
func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// GetCartList 获取用户购物车列表
func (r *orderRepo) GetCartList(ctx context.Context, userID int32) ([]*v1.ShopCartInfoResponse, error) {
	var carts []ShoppingCart

	// 查询该用户的所有购物车记录
	if err := r.data.DB(ctx).WithContext(ctx).Where(&ShoppingCart{User: userID}).Find(&carts).Error; err != nil {
		r.log.Errorf("failed to get cart list for user %d: %v", userID, err)
		return nil, err
	}

	// 转换为响应格式
	result := make([]*v1.ShopCartInfoResponse, 0, len(carts))
	for _, cart := range carts {
		result = append(result, &v1.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}

	return result, nil
}
