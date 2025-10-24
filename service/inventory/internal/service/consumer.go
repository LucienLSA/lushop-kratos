package service

import (
	"context"
	"inventory/internal/biz"
	"inventory/internal/domain"
	"inventory/internal/pkg/rocketmq"

	"github.com/go-kratos/kratos/v2/log"
)

// InventoryConsumerService 库存消费服务
type InventoryConsumerService struct {
	uc  *biz.InventoryUsecase
	log *log.Helper
}

// NewInventoryConsumerService 创建库存消费服务
func NewInventoryConsumerService(uc *biz.InventoryUsecase, logger log.Logger) *InventoryConsumerService {
	return &InventoryConsumerService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// HandleInventoryMessage 处理库存扣减消息
func (s *InventoryConsumerService) HandleInventoryMessage(ctx context.Context, msg *rocketmq.OrderInventoryMessage) error {
	s.log.Infof("handling inventory message: orderSn=%s, userId=%d, goods_count=%d",
		msg.OrderSn, msg.UserID, len(msg.GoodsInfo))

	// 转换消息为领域模型
	items := make([]domain.GoodsInvInfo, 0, len(msg.GoodsInfo))
	for _, gi := range msg.GoodsInfo {
		items = append(items, domain.GoodsInvInfo{
			GoodsID: gi.GoodsID,
			Nums:    gi.Nums,
		})
	}

	sellInfo := &domain.SellInfo{
		GoodsInvInfo: items,
		OrderSn:      msg.OrderSn,
	}

	// 调用业务层扣减库存
	if err := s.uc.Sell(ctx, sellInfo); err != nil {
		s.log.Errorf("failed to sell inventory: orderSn=%s, error=%v", msg.OrderSn, err)
		return err
	}

	s.log.Infof("inventory sold successfully: orderSn=%s", msg.OrderSn)
	return nil
}
