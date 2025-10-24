package service

import (
	"context"
	v1 "lushop/api/lushop/v1"
)

// GetGoodsList 获取商品列表 HTTP API
func (s *LushopService) GetGoodsList(ctx context.Context, req *v1.GoodsListReq) (*v1.GoodsListReply, error) {
	s.log.Infof("HTTP API: 获取商品列表 page=%d, pageSize=%d", req.Page, req.PageSize)

	// 调用业务逻辑（内部会调用 gRPC）
	goodsList, total, err := s.goodsUc.GetGoodsList(ctx, req.Page, req.PageSize, req.IsHot, req.IsNew)
	if err != nil {
		return nil, err
	}

	// 转换为 HTTP 响应
	items := make([]*v1.GoodsInfo, 0, len(goodsList))
	for _, goods := range goodsList {
		items = append(items, &v1.GoodsInfo{
			Id:              goods.ID,
			Name:            goods.Name,
			GoodsSn:         goods.GoodsSn,
			ShopPrice:       goods.ShopPrice,
			GoodsFrontImage: goods.GoodsFrontImage,
			IsNew:           goods.IsNew,
			IsHot:           goods.IsHot,
		})
	}

	return &v1.GoodsListReply{
		Total: total,
		Data:  items,
	}, nil
}

// GetGoodsDetail 获取商品详情 HTTP API
func (s *LushopService) GetGoodsDetail(ctx context.Context, req *v1.GoodsDetailReq) (*v1.GoodsDetailReply, error) {
	s.log.Infof("HTTP API: 获取商品详情 id=%d", req.Id)

	// 调用业务逻辑（内部会调用 gRPC）
	goods, err := s.goodsUc.GetGoodsDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 转换为 HTTP 响应
	return &v1.GoodsDetailReply{
		Id:              goods.ID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		CategoryId:      goods.CategoryID,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		GoodsDesc:       goods.GoodsDesc,
		Images:          goods.Images,
		DescImages:      goods.DescImages,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
	}, nil
}

// SearchGoods 搜索商品 HTTP API
func (s *LushopService) SearchGoods(ctx context.Context, req *v1.SearchGoodsReq) (*v1.GoodsListReply, error) {
	s.log.Infof("HTTP API: 搜索商品 keyword=%s", req.Keyword)

	// 调用业务逻辑（内部会调用 gRPC）
	goodsList, total, err := s.goodsUc.SearchGoods(ctx, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 转换为 HTTP 响应
	items := make([]*v1.GoodsInfo, 0, len(goodsList))
	for _, goods := range goodsList {
		items = append(items, &v1.GoodsInfo{
			Id:              goods.ID,
			Name:            goods.Name,
			ShopPrice:       goods.ShopPrice,
			GoodsFrontImage: goods.GoodsFrontImage,
		})
	}

	return &v1.GoodsListReply{
		Total: total,
		Data:  items,
	}, nil
}
