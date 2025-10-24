package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-kratos/kratos/v2/log"
)

// TransactionProducer RocketMQ 事务消息生产者
type TransactionProducer struct {
	producer rocketmq.TransactionProducer
	topic    string
	log      *log.Helper
}

// NewTransactionProducer 创建事务消息生产者
// listener: 事务监听器，用于执行本地事务和回查
func NewTransactionProducer(nameServers []string, groupName, topic string, listener primitive.TransactionListener, logger log.Logger) (*TransactionProducer, error) {
	p, err := rocketmq.NewTransactionProducer(
		listener,
		producer.WithNameServer(nameServers),
		producer.WithRetry(2),
		producer.WithGroupName(groupName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction producer: %w", err)
	}

	if err := p.Start(); err != nil {
		return nil, fmt.Errorf("failed to start transaction producer: %w", err)
	}

	return &TransactionProducer{
		producer: p,
		topic:    topic,
		log:      log.NewHelper(logger),
	}, nil
}

// SendTransactionMessage 发送事务消息
// key: 消息 key（通常是订单号）
// body: 消息体
// arg: 传递给本地事务执行器的参数
func (p *TransactionProducer) SendTransactionMessage(ctx context.Context, key string, body interface{}, arg interface{}) (*primitive.TransactionSendResult, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &primitive.Message{
		Topic: p.topic,
		Body:  data,
	}
	msg.WithKeys([]string{key})

	// 发送事务消息
	result, err := p.producer.SendMessageInTransaction(ctx, msg)
	if err != nil {
		p.log.Errorf("failed to send transaction message: %v", err)
		return nil, err
	}

	p.log.Infof("send transaction message success: key=%s, sendStatus=%s",
		key, result.State)
	return result, nil
}

// Close 关闭生产者
func (p *TransactionProducer) Close() error {
	return p.producer.Shutdown()
}
