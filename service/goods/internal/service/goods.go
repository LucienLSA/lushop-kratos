package service

import (
	"context"
	v1 "goods/api/goods/v1"
	"goods/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateGoods 创建商品
func (g *GoodsService) CreateGoods(ctx context.Context, r *v1.CreateGoodsInfo) (*v1.GoodsInfoResponse, error) {
	goodsInfo := domain.Goods{
		ID:              int64(r.Id),
		CategoryID:      r.CategoryId,
		BrandsID:        r.BrandId,
		Name:            r.Name,
		GoodsSn:         r.GoodsSn,
		OnSale:          r.OnSale,
		MarketPrice:     r.MarketPrice,
		ShopPrice:       r.ShopPrice,
		GoodsBrief:      r.GoodsBrief,
		GoodsFrontImage: r.GoodsFrontImage,
		GoodsImages:     r.Images,
		ShipFree:        r.ShipFree,
		IsNew:           r.IsNew,
		IsHot:           r.IsHot,
	}

	result, err := g.g.CreateGoods(ctx, &goodsInfo)
	if err != nil {
		return nil, err
	}
	return &v1.GoodsInfoResponse{Id: int32(result.GoodsID)}, nil
}

func (g *GoodsService) GoodsList(ctx context.Context, r *v1.GoodsFilterRequest) (*v1.GoodsListResponse, error) {
	return nil, nil
}

func (g *GoodsService) GoodsListES(ctx context.Context, r *v1.GoodsFilterRequest) (*v1.GoodsListResponse, error) {
	// 将 proto 过滤条件映射为 domain 的 ES 过滤条件
	goodsFilter := &domain.ESGoodsFilter{
		CategoryID:  r.TopCategory,
		BrandsID:    r.Brand,
		Keywords:    r.KeyWords,
		IsNew:       r.IsNew,
		IsHot:       r.IsHot,
		MaxPrice:    int64(r.PriceMax),
		MinPrice:    int64(r.PriceMin),
		Pages:       int64(r.Pages),
		PagePerNums: int64(r.PagePerNums),
	}

	result, err := g.esGoods.GoodsListES(ctx, goodsFilter)
	if err != nil {
		return nil, err
	}
	response := v1.GoodsListResponse{
		Total: int32(result.Total),
	}
	for _, goods := range result.List {
		res := v1.GoodsInfoResponse{
			Id:              int32(goods.ID),
			CategoryId:      goods.CategoryID,
			Name:            goods.Name,
			GoodsSn:         goods.GoodsSn,
			ClickNum:        int32(goods.ClickNum),
			SoldNum:         int32(goods.SoldNum),
			FavNum:          int32(goods.FavNum),
			MarketPrice:     goods.MarketPrice,
			ShopPrice:       goods.ShopPrice,
			GoodsBrief:      goods.GoodsBrief,
			GoodsDesc:       goods.GoodsBrief,
			ShipFree:        goods.ShipFree,
			Images:          goods.GoodsImages,
			DescImages:      goods.DescImages,
			GoodsFrontImage: goods.GoodsFrontImage,
			IsNew:           goods.IsNew,
			IsHot:           goods.IsHot,
			OnSale:          goods.OnSale,
		}
		response.Data = append(response.Data, &res)
	}
	return &response, nil
}

func (g *GoodsService) BatchGetGoods(ctx context.Context, r *v1.BatchGoodsIdInfo) (*v1.GoodsListResponse, error) {
	var ids []int64
	for _, id := range r.Id {
		ids = append(ids, int64(id))
	}
	result, err := g.g.BatchGetGoods(ctx, ids)
	if err != nil {
		return nil, err
	}
	rsp := &v1.GoodsListResponse{
		Total: int32(result.Total),
	}
	for _, goods := range result.List {
		item := &v1.GoodsInfoResponse{
			Id:              int32(goods.ID),
			CategoryId:      int32(goods.CategoryID),
			Name:            goods.Name,
			GoodsSn:         goods.GoodsSn,
			ClickNum:        int32(goods.ClickNum),
			SoldNum:         int32(goods.SoldNum),
			FavNum:          int32(goods.FavNum),
			MarketPrice:     goods.MarketPrice,
			ShopPrice:       goods.ShopPrice,
			GoodsBrief:      goods.GoodsBrief,
			GoodsDesc:       goods.GoodsBrief,
			ShipFree:        goods.ShipFree,
			Images:          goods.GoodsImages,
			DescImages:      goods.DescImages,
			GoodsFrontImage: goods.GoodsFrontImage,
			IsNew:           goods.IsNew,
			IsHot:           goods.IsHot,
			OnSale:          goods.OnSale,
		}
		rsp.Data = append(rsp.Data, item)
	}
	return rsp, nil
}

func (g *GoodsService) UpdateGoods(ctx context.Context, r *v1.CreateGoodsInfo) (*emptypb.Empty, error) {
	goods := &domain.Goods{
		ID:              int64(r.Id),
		CategoryID:      r.CategoryId,
		BrandsID:        r.BrandId,
		Name:            r.Name,
		GoodsSn:         r.GoodsSn,
		OnSale:          r.OnSale,
		MarketPrice:     r.MarketPrice,
		ShopPrice:       r.ShopPrice,
		GoodsBrief:      r.GoodsBrief,
		GoodsFrontImage: r.GoodsFrontImage,
		DescImages:      r.DescImages,
		GoodsImages:     r.Images,
		ShipFree:        r.ShipFree,
		IsNew:           r.IsNew,
		IsHot:           r.IsHot,
	}
	if _, err := g.g.UpdateGoods(ctx, goods); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsService) DeleteGoods(ctx context.Context, r *v1.DeleteGoodsInfo) (*emptypb.Empty, error) {
	err := g.g.DeleteGoods(ctx, int64(r.Id))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsService) GetGoodsDetail(ctx context.Context, r *v1.GoodInfoRequest) (*v1.GoodsInfoResponse, error) {
	goods, err := g.g.GetGoodsById(ctx, int64(r.Id))
	if err != nil {
		return nil, err
	}
	resp := &v1.GoodsInfoResponse{
		Id:              int32(goods.ID),
		CategoryId:      int32(goods.CategoryID),
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        int32(goods.ClickNum),
		SoldNum:         int32(goods.SoldNum),
		FavNum:          int32(goods.FavNum),
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		GoodsDesc:       goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		Images:          goods.GoodsImages,
		DescImages:      goods.DescImages,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
	}
	return resp, nil
}
