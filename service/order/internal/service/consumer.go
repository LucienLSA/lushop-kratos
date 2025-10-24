package service

import (
	"context"
	"encoding/json"
	"order/internal/biz"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
)

// OrderTimeoutConsumer 订单超时消费者
type OrderTimeoutConsumer struct {
	uc  *biz.OrderUsecase
	log *log.Helper
}

// NewOrderTimeoutConsumer 创建订单超时消费者
func NewOrderTimeoutConsumer(uc *biz.OrderUsecase, logger log.Logger) *OrderTimeoutConsumer {
	return &OrderTimeoutConsumer{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// OrderTimeoutMessage 订单超时消息
type OrderTimeoutMessage struct {
	OrderSn string `json:"order_sn"`
	UserId  int32  `json:"user_id"`
}

// OrderTimeout 订单超时处理
// 监听订单超时消息，自动归还库存并关闭订单
// 幂等性：通过订单状态判断，避免重复处理
// 事务性：订单状态更新和发送归还库存消息在同一事务中
func (c *OrderTimeoutConsumer) OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var orderInfo OrderTimeoutMessage
		if err := json.Unmarshal(msgs[i].Body, &orderInfo); err != nil {
			c.log.Errorf("解析订单超时消息失败: %v, body: %s", err, string(msgs[i].Body))
			// 消息格式错误，直接丢弃
			return consumer.ConsumeSuccess, nil
		}

		c.log.Infof("收到订单超时消息: orderSn=%s, time=%v", orderInfo.OrderSn, time.Now())

		// 调用业务层处理订单超时
		if err := c.uc.HandleOrderTimeout(ctx, orderInfo.OrderSn); err != nil {
			c.log.Errorf("处理订单超时失败: orderSn=%s, error=%v", orderInfo.OrderSn, err)
			// 处理失败，返回重试
			return consumer.ConsumeRetryLater, err
		}

		c.log.Infof("订单超时处理成功: orderSn=%s", orderInfo.OrderSn)
	}

	return consumer.ConsumeSuccess, nil
}
