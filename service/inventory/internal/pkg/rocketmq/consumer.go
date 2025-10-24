package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
)

// MessageHandler 消息处理函数
type MessageHandler func(ctx context.Context, msg *OrderInventoryMessage) error

// Consumer RocketMQ 消费者封装
type Consumer struct {
	consumer rocketmq.PushConsumer
	handler  MessageHandler
	log      *log.Helper
}

// NewConsumer 创建 RocketMQ 消费者
func NewConsumer(nameServers []string, groupName, topic string, handler MessageHandler, logger log.Logger) (*Consumer, error) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(nameServers),
		consumer.WithGroupName(groupName),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	consumerWrapper := &Consumer{
		consumer: c,
		handler:  handler,
		log:      log.NewHelper(logger),
	}

	// 订阅 Topic
	err = c.Subscribe(topic, consumer.MessageSelector{}, consumerWrapper.handleMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe topic: %w", err)
	}

	return consumerWrapper, nil
}

// handleMessage 处理消息
func (c *Consumer) handleMessage(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		c.log.Infof("received message: msgId=%s, topic=%s, tags=%s", msg.MsgId, msg.Topic, msg.GetTags())

		// 解析消息
		var orderMsg OrderInventoryMessage
		if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
			c.log.Errorf("failed to unmarshal message: msgId=%s, error=%v", msg.MsgId, err)
			// 解析失败，直接返回成功，避免重复消费
			continue
		}

		// 调用业务处理函数
		if err := c.handler(ctx, &orderMsg); err != nil {
			c.log.Errorf("failed to handle message: msgId=%s, orderSn=%s, error=%v", msg.MsgId, orderMsg.OrderSn, err)
			// 返回稍后重试
			return consumer.ConsumeRetryLater, err
		}

		c.log.Infof("message handled successfully: msgId=%s, orderSn=%s", msg.MsgId, orderMsg.OrderSn)
	}

	return consumer.ConsumeSuccess, nil
}

// Start 启动消费者
func (c *Consumer) Start() error {
	if err := c.consumer.Start(); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}
	c.log.Info("RocketMQ consumer started")
	return nil
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	return c.consumer.Shutdown()
}
