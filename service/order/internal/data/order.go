package data

import (
	"context"
	"encoding/json"
	v1 "order/api/order/v1"
	"order/internal/biz"
	"order/internal/conf"
	"order/internal/pkg/rocketmq"
	"order/internal/pkg/snowflake"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

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
func NewOrderRepo(data *Data, c *conf.Bootstrap, logger log.Logger) biz.OrderRepo {
	repo := &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
	
	// 创建事务生产者（如果启用了 RocketMQ）
	if c.Rocketmq != nil && c.Rocketmq.Enable && data.txProducer == nil {
		txProducer, err := rocketmq.NewTransactionProducer(
			c.Rocketmq.NameServer,
			c.Rocketmq.GroupName+"-tx",
			c.Rocketmq.Topic,
			repo, // repo 实现了 primitive.TransactionListener 接口
			logger,
		)
		if err != nil {
			log.NewHelper(logger).Errorf("failed to create transaction producer: %v", err)
		} else {
			data.txProducer = txProducer
			log.NewHelper(logger).Info("RocketMQ transaction producer created successfully")
		}
	}
	
	return repo
}

// GenerateOrderSn 使用雪花算法生成订单号
// 雪花算法保证全局唯一性和高性能
func GenerateOrderSn(userId int32) string {
	// 使用项目内置的 snowflake 包生成唯一ID
	return snowflake.GenerateIDString()
}

// CreateOrder 使用 RocketMQ 事务消息创建订单
// 保证订单创建和库存扣减的最终一致性
func (r *orderRepo) CreateOrder(ctx context.Context, req *v1.OrderRequest) (*v1.OrderInfoResponse, error) {
	// 如果事务生产者可用，使用事务消息方案
	if r.data.txProducer != nil {
		return r.CreateOrderWithTransactionMessage(ctx, req)
	}
	
	// 否则使用原来的方案（普通消息，存在一致性风险）
	r.log.Warn("transaction producer is nil, using legacy order creation method")
	return r.createOrderLegacy(ctx, req)
}

// createOrderLegacy 原有的订单创建方法（保留作为降级方案）
// 注意：这个方法存在数据一致性风险，仅在事务消息不可用时使用
func (r *orderRepo) createOrderLegacy(ctx context.Context, req *v1.OrderRequest) (*v1.OrderInfoResponse, error) {
	var orderInfo OrderInfo
	var totalAmount float32

	// 开启事务
	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查询用户购物车中选中的商品
		var cartItems []ShoppingCart
		if err := tx.Where(&ShoppingCart{User: req.UserId, Checked: true}).Find(&cartItems).Error; err != nil {
			r.log.Errorf("failed to query cart items: %v", err)
			return err
		}

		if len(cartItems) == 0 {
			return errors.BadRequest("EMPTY_CART", "no items in cart")
		}

		// 2. 生成订单号
		orderSn := GenerateOrderSn(req.UserId)

		// 3. 创建订单基本信息
		orderInfo = OrderInfo{
			User:         req.UserId,
			OrderSn:      orderSn,
			Address:      req.Address,
			SignerName:   req.Name,
			SingerMobile: req.Mobile,
			Post:         req.Post,
			Status:       "WAIT_BUYER_PAY", // 待支付状态
		}

		if err := tx.Create(&orderInfo).Error; err != nil {
			r.log.Errorf("failed to create order: %v", err)
			return err
		}

		// 4. 批量查询商品信息
		goodsIds := make([]int32, 0, len(cartItems))
		for _, item := range cartItems {
			goodsIds = append(goodsIds, item.Goods)
		}

		// 调用商品服务批量获取商品信息
		goodsMap, err := r.data.goodsClient.BatchGetGoods(ctx, goodsIds)
		if err != nil {
			r.log.Errorf("failed to get goods info: %v", err)
			return err
		}

		// 5. 创建订单商品明细并计算总金额
		totalAmount = 0
		for _, item := range cartItems {
			goods, ok := goodsMap[item.Goods]
			if !ok {
				r.log.Errorf("goods not found: goods_id=%d", item.Goods)
				return errors.NotFound("GOODS_NOT_FOUND", "goods not found")
			}

			orderGoods := OrderGoods{
				Order:      orderInfo.ID,
				Goods:      item.Goods,
				GoodsName:  goods.Name,
				GoodsImage: goods.GoodsFrontImage,
				GoodsPrice: goods.ShopPrice,
				Nums:       item.Nums,
			}

			if err := tx.Create(&orderGoods).Error; err != nil {
				r.log.Errorf("failed to create order goods: %v", err)
				return err
			}

			totalAmount += goods.ShopPrice * float32(item.Nums)
		}

		// 6. 更新订单总金额
		orderInfo.OrderMount = totalAmount
		if err := tx.Save(&orderInfo).Error; err != nil {
			r.log.Errorf("failed to update order amount: %v", err)
			return err
		}

		// 7. 清空购物车选中的商品
		if err := tx.Where(&ShoppingCart{User: req.UserId, Checked: true}).Delete(&ShoppingCart{}).Error; err != nil {
			r.log.Errorf("failed to delete cart items: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	r.log.Infof("order created: order_sn=%s, user=%d, amount=%.2f", orderInfo.OrderSn, req.UserId, totalAmount)

	// 事务成功后，发送 MQ 消息到库存服务进行库存扣减
	if r.data.producer != nil {
		// 查询订单商品明细
		var orderGoods []OrderGoods
		if err := r.data.DB(ctx).Where("order = ?", orderInfo.ID).Find(&orderGoods).Error; err != nil {
			r.log.Errorf("failed to query order goods: %v", err)
		} else {
			// 构建消息
			goodsInfo := make([]rocketmq.GoodsInvInfo, 0, len(orderGoods))
			for _, goods := range orderGoods {
				goodsInfo = append(goodsInfo, rocketmq.GoodsInvInfo{
					GoodsID: goods.Goods,
					Nums:    goods.Nums,
				})
			}

			msg := rocketmq.OrderInventoryMessage{
				OrderSn:   orderInfo.OrderSn,
				UserID:    orderInfo.User,
				GoodsInfo: goodsInfo,
			}

			// 异步发送消息（不阻塞订单创建）
			go func() {
				if err := r.data.producer.SendMessage(context.Background(), orderInfo.OrderSn, msg); err != nil {
					r.log.Errorf("failed to send inventory message: order_sn=%s, error=%v", orderInfo.OrderSn, err)
					// TODO: 可以将失败的消息存入数据库，通过定时任务补偿
				} else {
					r.log.Infof("inventory message sent: order_sn=%s", orderInfo.OrderSn)
				}
			}()
		}
	}

	// 返回订单信息
	return &v1.OrderInfoResponse{
		Id:      orderInfo.ID,
		UserId:  orderInfo.User,
		OrderSn: orderInfo.OrderSn,
		Status:  orderInfo.Status,
		Total:   orderInfo.OrderMount,
		Address: orderInfo.Address,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SingerMobile,
		Post:    orderInfo.Post,
		AddTime: orderInfo.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetOrderList 获取订单列表
func (r *orderRepo) GetOrderList(ctx context.Context, req *v1.OrderFilterRequest) ([]*v1.OrderInfoResponse, int32, error) {
	var orders []OrderInfo
	var total int64

	// 查询总数
	if err := r.data.DB(ctx).Model(&OrderInfo{}).Where("user = ?", req.UserId).Count(&total).Error; err != nil {
		r.log.Errorf("failed to count orders: %v", err)
		return nil, 0, err
	}

	// 查询订单列表（使用 Paginate 分页）
	if err := r.data.DB(ctx).
		Where("user = ?", req.UserId).
		Order("id desc").
		Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).
		Find(&orders).Error; err != nil {
		r.log.Errorf("failed to get order list: %v", err)
		return nil, 0, err
	}

	// 转换为响应格式
	result := make([]*v1.OrderInfoResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, &v1.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, int32(total), nil
}

// GetOrderDetail 获取订单详情
func (r *orderRepo) GetOrderDetail(ctx context.Context, req *v1.OrderRequest) (*v1.OrderInfoDetailResponse, error) {
	var order OrderInfo
	var orderGoods []OrderGoods

	// 查询订单基本信息（验证用户所有权）
	if err := r.data.DB(ctx).Where("id = ? AND user = ?", req.Id, req.UserId).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Errorf("order not found: id=%d, user=%d", req.Id, req.UserId)
			return nil, errors.NotFound("ORDER_NOT_FOUND", "order not found")
		}
		r.log.Errorf("failed to get order: %v", err)
		return nil, err
	}

	// 查询订单商品明细
	if err := r.data.DB(ctx).Where("order = ?", order.ID).Find(&orderGoods).Error; err != nil {
		r.log.Errorf("failed to get order goods: %v", err)
		return nil, err
	}

	// 转换订单信息
	orderInfo := &v1.OrderInfoResponse{
		Id:      order.ID,
		UserId:  order.User,
		OrderSn: order.OrderSn,
		PayType: order.PayType,
		Status:  order.Status,
		Post:    order.Post,
		Total:   order.OrderMount,
		Address: order.Address,
		Name:    order.SignerName,
		Mobile:  order.SingerMobile,
		AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 转换商品明细
	goods := make([]*v1.OrderItemResponse, 0, len(orderGoods))
	for _, item := range orderGoods {
		goods = append(goods, &v1.OrderItemResponse{
			Id:         item.ID,
			OrderId:    item.Order,
			GoodsId:    item.Goods,
			GoodsName:  item.GoodsName,
			GoodsImage: item.GoodsImage,
			GoodsPrice: item.GoodsPrice,
			Nums:       item.Nums,
		})
	}

	return &v1.OrderInfoDetailResponse{
		OrderInfo: orderInfo,
		Goods:     goods,
	}, nil
}

// UpdateOrderStatus 更新订单状态
func (r *orderRepo) UpdateOrderStatus(ctx context.Context, req *v1.OrderStatus) (*emptypb.Empty, error) {
	var order OrderInfo

	// 查询订单
	if err := r.data.DB(ctx).Where("order_sn = ?", req.OrderSn).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Errorf("order not found: order_sn=%s", req.OrderSn)
			return nil, errors.NotFound("ORDER_NOT_FOUND", "order not found")
		}
		r.log.Errorf("failed to get order: %v", err)
		return nil, err
	}

	// 更新订单状态
	order.Status = req.Status
	if err := r.data.DB(ctx).Save(&order).Error; err != nil {
		r.log.Errorf("failed to update order status: %v", err)
		return nil, err
	}

	r.log.Infof("order status updated: order_sn=%s, status=%s", req.OrderSn, req.Status)
	return &emptypb.Empty{}, nil
}

// HandleOrderTimeout 处理订单超时
// 1. 查询订单状态
// 2. 如果订单未支付，发送归还库存消息到 MQ
// 3. 更新订单状态为 TRADE_CLOSED
// 幂等性：通过订单状态判断，已支付或已关闭的订单不处理
func (r *orderRepo) HandleOrderTimeout(ctx context.Context, orderSn string) error {
	return r.data.ExecTx(ctx, func(ctx context.Context) error {
		db := r.data.DB(ctx)
		
		// 查询订单
		var order OrderInfo
		if err := db.Where("order_sn = ?", orderSn).First(&order).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 订单不存在，幂等处理
				r.log.Infof("订单不存在，无需处理: order_sn=%s", orderSn)
				return nil
			}
			r.log.Errorf("查询订单失败: order_sn=%s, error=%v", orderSn, err)
			return errors.InternalServer("QUERY_ORDER_ERROR", err.Error())
		}
		
		r.log.Infof("订单状态: order_sn=%s, status=%s", orderSn, order.Status)
		
		// 如果订单已支付或已关闭，不处理
		if order.Status == "TRADE_SUCCESS" || order.Status == "TRADE_CLOSED" || order.Status == "TRADE_FINISHED" {
			r.log.Infof("订单已支付或已关闭，无需处理: order_sn=%s, status=%s", orderSn, order.Status)
			return nil
		}
		
		// 发送归还库存消息到 MQ
		type RebackMessage struct {
			OrderSn string `json:"order_sn"`
		}
		msg := RebackMessage{OrderSn: orderSn}
		
		if err := r.data.producer.SendMessage(ctx, orderSn, msg); err != nil {
			r.log.Errorf("发送归还库存消息失败: order_sn=%s, error=%v", orderSn, err)
			return errors.InternalServer("SEND_REBACK_MESSAGE_ERROR", err.Error())
		}
		
		// 更新订单状态为 TRADE_CLOSED
		order.Status = "TRADE_CLOSED"
		if err := db.Save(&order).Error; err != nil {
			r.log.Errorf("更新订单状态失败: order_sn=%s, error=%v", orderSn, err)
			return errors.InternalServer("UPDATE_ORDER_STATUS_ERROR", err.Error())
		}
		
		r.log.Infof("订单超时处理成功: order_sn=%s, 已发送归还库存消息并关闭订单", orderSn)
		return nil
	})
}

// ==================== RocketMQ 事务消息监听器 ====================

// ExecuteLocalTransaction 执行本地事务
// 当发送半消息（Half Message）成功后，Broker 会回调此方法执行本地事务
// 返回值：
//   - CommitMessageState: 提交消息，Consumer 可以消费
//   - RollbackMessageState: 回滚消息，Consumer 不会消费
//   - UnknowState: 未知状态，Broker 会定时回查
func (r *orderRepo) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	r.log.Infof("executing local transaction: keys=%v", msg.GetKeys())

	// 从消息中解析订单库存扣减消息
	var orderMsg rocketmq.OrderInventoryMessage
	if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
		r.log.Errorf("failed to unmarshal message: %v", err)
		return primitive.RollbackMessageState
	}

	orderSn := orderMsg.OrderSn
	ctx := context.Background()

	// 执行本地事务：创建订单
	err := r.createOrderInTransaction(ctx, orderSn, orderMsg.UserID)
	if err != nil {
		r.log.Errorf("local transaction failed: orderSn=%s, error=%v", orderSn, err)
		return primitive.RollbackMessageState
	}

	r.log.Infof("local transaction success: orderSn=%s", orderSn)
	return primitive.CommitMessageState
}

// CheckLocalTransaction 回查本地事务状态
// 当 ExecuteLocalTransaction 返回 UnknowState 或网络异常时，Broker 会定时回调此方法
// 通过查询数据库确认本地事务是否执行成功
func (r *orderRepo) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	r.log.Infof("checking local transaction: msgId=%s, keys=%v", msg.MsgId, msg.GetKeys())

	// 从消息 key 中获取订单号
	keys := msg.GetKeys()
	if len(keys) == 0 {
		r.log.Error("message keys is empty")
		return primitive.RollbackMessageState
	}
	orderSn := keys[0]

	// 查询订单是否存在
	var order OrderInfo
	ctx := context.Background()
	if err := r.data.DB(ctx).Where("order_sn = ?", orderSn).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Infof("order not found, rollback message: orderSn=%s", orderSn)
			return primitive.RollbackMessageState
		}
		// 数据库查询异常，返回未知状态，稍后再查
		r.log.Errorf("failed to query order: orderSn=%s, error=%v", orderSn, err)
		return primitive.UnknowState
	}

	r.log.Infof("order exists, commit message: orderSn=%s, status=%s", orderSn, order.Status)
	return primitive.CommitMessageState
}

// createOrderInTransaction 在事务中创建订单
// 这是本地事务的核心逻辑，只负责创建订单，不发送消息（消息由事务消息机制保证）
func (r *orderRepo) createOrderInTransaction(ctx context.Context, orderSn string, userId int32) error {
	return r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查询用户购物车中选中的商品
		var cartItems []ShoppingCart
		if err := tx.Where(&ShoppingCart{User: userId, Checked: true}).Find(&cartItems).Error; err != nil {
			r.log.Errorf("failed to query cart items: %v", err)
			return err
		}

		if len(cartItems) == 0 {
			return errors.BadRequest("EMPTY_CART", "no items in cart")
		}

		// 2. 创建订单基本信息（使用传入的订单号）
		orderInfo := OrderInfo{
			User:    userId,
			OrderSn: orderSn,
			Status:  "WAIT_BUYER_PAY", // 待支付状态
		}

		if err := tx.Create(&orderInfo).Error; err != nil {
			r.log.Errorf("failed to create order: %v", err)
			return err
		}

		// 3. 批量查询商品信息
		goodsIds := make([]int32, 0, len(cartItems))
		for _, item := range cartItems {
			goodsIds = append(goodsIds, item.Goods)
		}

		// 调用商品服务批量获取商品信息
		goodsMap, err := r.data.goodsClient.BatchGetGoods(ctx, goodsIds)
		if err != nil {
			r.log.Errorf("failed to get goods info: %v", err)
			return err
		}

		// 4. 创建订单商品明细并计算总金额
		var totalAmount float32
		for _, item := range cartItems {
			goods, ok := goodsMap[item.Goods]
			if !ok {
				r.log.Errorf("goods not found: goods_id=%d", item.Goods)
				return errors.NotFound("GOODS_NOT_FOUND", "goods not found")
			}

			orderGoods := OrderGoods{
				Order:      orderInfo.ID,
				Goods:      item.Goods,
				GoodsName:  goods.Name,
				GoodsImage: goods.GoodsFrontImage,
				GoodsPrice: goods.ShopPrice,
				Nums:       item.Nums,
			}

			if err := tx.Create(&orderGoods).Error; err != nil {
				r.log.Errorf("failed to create order goods: %v", err)
				return err
			}

			totalAmount += goods.ShopPrice * float32(item.Nums)
		}

		// 5. 更新订单总金额
		orderInfo.OrderMount = totalAmount
		if err := tx.Save(&orderInfo).Error; err != nil {
			r.log.Errorf("failed to update order amount: %v", err)
			return err
		}

		// 6. 清空购物车选中的商品
		if err := tx.Where(&ShoppingCart{User: userId, Checked: true}).Delete(&ShoppingCart{}).Error; err != nil {
			r.log.Errorf("failed to delete cart items: %v", err)
			return err
		}

		r.log.Infof("order created in transaction: order_sn=%s, user=%d, amount=%.2f", orderSn, userId, totalAmount)
		return nil
	})
}

// CreateOrderWithTransactionMessage 使用事务消息创建订单
// 这是对外暴露的方法，保证订单创建和库存扣减的最终一致性
func (r *orderRepo) CreateOrderWithTransactionMessage(ctx context.Context, req *v1.OrderRequest) (*v1.OrderInfoResponse, error) {
	// 1. 预先生成订单号
	orderSn := snowflake.GenerateIDString()

	// 2. 查询购物车商品信息（用于构建消息）
	var cartItems []ShoppingCart
	if err := r.data.DB(ctx).Where(&ShoppingCart{User: req.UserId, Checked: true}).Find(&cartItems).Error; err != nil {
		r.log.Errorf("failed to query cart items: %v", err)
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.BadRequest("EMPTY_CART", "no items in cart")
	}

	// 3. 构建库存扣减消息
	goodsInfo := make([]rocketmq.GoodsInvInfo, 0, len(cartItems))
	for _, item := range cartItems {
		goodsInfo = append(goodsInfo, rocketmq.GoodsInvInfo{
			GoodsID: item.Goods,
			Nums:    item.Nums,
		})
	}

	msg := rocketmq.OrderInventoryMessage{
		OrderSn:   orderSn,
		UserID:    req.UserId,
		GoodsInfo: goodsInfo,
	}

	// 4. 发送事务消息
	// 事务消息会先发送半消息，然后执行本地事务（创建订单）
	// 如果本地事务成功，提交消息；失败则回滚消息
	result, err := r.data.txProducer.SendTransactionMessage(ctx, orderSn, msg, nil)
	if err != nil {
		r.log.Errorf("failed to send transaction message: orderSn=%s, error=%v", orderSn, err)
		return nil, errors.InternalServer("SEND_TRANSACTION_MESSAGE_ERROR", err.Error())
	}

	r.log.Infof("transaction message sent: orderSn=%s, msgId=%s, state=%s",
		orderSn, result.MsgID, result.State)

	// 5. 查询订单信息返回
	var orderInfo OrderInfo
	if err := r.data.DB(ctx).Where("order_sn = ?", orderSn).First(&orderInfo).Error; err != nil {
		r.log.Errorf("failed to query order: orderSn=%s, error=%v", orderSn, err)
		return nil, err
	}

	// 6. 发送订单超时延迟消息（30分钟后触发）
	if r.data.producer != nil {
		timeoutMsg := struct {
			OrderSn string `json:"order_sn"`
			UserId  int32  `json:"user_id"`
		}{
			OrderSn: orderSn,
			UserId:  req.UserId,
		}

		// delayLevel 16 = 30分钟
		go func() {
			if err := r.data.producer.SendDelayMessage(context.Background(), orderSn, timeoutMsg, 16); err != nil {
				r.log.Errorf("failed to send timeout message: order_sn=%s, error=%v", orderSn, err)
			} else {
				r.log.Infof("timeout message sent: order_sn=%s, delay=30min", orderSn)
			}
		}()
	}

	// 返回订单信息
	return &v1.OrderInfoResponse{
		Id:      orderInfo.ID,
		UserId:  orderInfo.User,
		OrderSn: orderInfo.OrderSn,
		Status:  orderInfo.Status,
		Total:   orderInfo.OrderMount,
		Address: orderInfo.Address,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SingerMobile,
		Post:    orderInfo.Post,
		AddTime: orderInfo.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
