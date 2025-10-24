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

// Producer RocketMQ 生产者封装
type Producer struct {
	producer rocketmq.Producer
	topic    string
	log      *log.Helper
}

// NewProducer 创建 RocketMQ 生产者
func NewProducer(nameServers []string, groupName, topic string, logger log.Logger) (*Producer, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(nameServers),
		producer.WithRetry(2),
		producer.WithGroupName(groupName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	if err := p.Start(); err != nil {
		return nil, fmt.Errorf("failed to start producer: %w", err)
	}

	return &Producer{
		producer: p,
		topic:    topic,
		log:      log.NewHelper(logger),
	}, nil
}

// SendMessage 发送消息
func (p *Producer) SendMessage(ctx context.Context, key string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &primitive.Message{
		Topic: p.topic,
		Body:  data,
	}
	msg.WithKeys([]string{key})

	result, err := p.producer.SendSync(ctx, msg)
	if err != nil {
		p.log.Errorf("failed to send message: %v", err)
		return err
	}

	p.log.Infof("send message success: msgId=%s, key=%s", result.MsgID, key)
	return nil
}

// SendDelayMessage 发送延迟消息
func (p *Producer) SendDelayMessage(ctx context.Context, key string, body interface{}, delayLevel int) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &primitive.Message{
		Topic: p.topic,
		Body:  data,
	}
	msg.WithKeys([]string{key})
	msg.WithDelayTimeLevel(delayLevel)

	result, err := p.producer.SendSync(ctx, msg)
	if err != nil {
		p.log.Errorf("failed to send delay message: %v", err)
		return err
	}

	p.log.Infof("send delay message success: msgId=%s, key=%s, delayLevel=%d", result.MsgID, key, delayLevel)
	return nil
}

// Close 关闭生产者
func (p *Producer) Close() error {
	return p.producer.Shutdown()
}
