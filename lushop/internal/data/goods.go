package data

import (
	"context"
	goodsV1 "lushop/api/service/goods/v1"
	"lushop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type goodsRepo struct {
	data *Data
	log  *log.Helper
}

// NewGoodsRepo 创建商品仓库
func NewGoodsRepo(data *Data, logger log.Logger) biz.GoodsRepo {
	return &goodsRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// GetGoodsList 获取商品列表
func (r *goodsRepo) GetGoodsList(ctx context.Context, page, pageSize int32, isHot, isNew bool) ([]*biz.Goods, int32, error) {
	r.log.Infof("获取商品列表: page=%d, pageSize=%d, isHot=%v, isNew=%v", page, pageSize, isHot, isNew)

	// 调用商品服务 gRPC API
	resp, err := r.data.gc.GoodsList(ctx, &goodsV1.GoodsFilterRequest{
		PriceMin:    0,
		PriceMax:    0,
		IsHot:       isHot,
		IsNew:       isNew,
		IsTab:       false,
		TopCategory: 0,
		Pages:       page,
		PagePerNums: pageSize,
	})
	if err != nil {
		r.log.Errorf("获取商品列表失败: %v", err)
		return nil, 0, err
	}

	// 转换为业务对象
	goodsList := make([]*biz.Goods, 0, len(resp.Data))
	for _, item := range resp.Data {
		goodsList = append(goodsList, &biz.Goods{
			ID:              item.Id,
			Name:            item.Name,
			GoodsSn:         item.GoodsSn,
			CategoryID:      item.CategoryId,
			ClickNum:        item.ClickNum,
			SoldNum:         item.SoldNum,
			FavNum:          item.FavNum,
			MarketPrice:     item.MarketPrice,
			ShopPrice:       item.ShopPrice,
			GoodsBrief:      item.GoodsBrief,
			ShipFree:        item.ShipFree,
			Images:          item.Images,
			GoodsFrontImage: item.GoodsFrontImage,
			IsNew:           item.IsNew,
			IsHot:           item.IsHot,
			OnSale:          item.OnSale,
		})
	}

	return goodsList, resp.Total, nil
}

// GetGoodsDetail 获取商品详情
func (r *goodsRepo) GetGoodsDetail(ctx context.Context, id int32) (*biz.Goods, error) {
	r.log.Infof("获取商品详情: id=%d", id)

	// 调用商品服务 gRPC API
	resp, err := r.data.gc.GetGoodsDetail(ctx, &goodsV1.GoodInfoRequest{
		Id: id,
	})
	if err != nil {
		r.log.Errorf("获取商品详情失败: id=%d, error=%v", id, err)
		return nil, err
	}

	// 转换为业务对象
	return &biz.Goods{
		ID:              resp.Id,
		Name:            resp.Name,
		GoodsSn:         resp.GoodsSn,
		CategoryID:      resp.CategoryId,
		ClickNum:        resp.ClickNum,
		SoldNum:         resp.SoldNum,
		FavNum:          resp.FavNum,
		MarketPrice:     resp.MarketPrice,
		ShopPrice:       resp.ShopPrice,
		GoodsBrief:      resp.GoodsBrief,
		GoodsDesc:       resp.GoodsDesc,
		ShipFree:        resp.ShipFree,
		Images:          resp.Images,
		DescImages:      resp.DescImages,
		GoodsFrontImage: resp.GoodsFrontImage,
		IsNew:           resp.IsNew,
		IsHot:           resp.IsHot,
		OnSale:          resp.OnSale,
	}, nil
}

// SearchGoods 搜索商品
func (r *goodsRepo) SearchGoods(ctx context.Context, keyword string, page, pageSize int32) ([]*biz.Goods, int32, error) {
	r.log.Infof("搜索商品: keyword=%s, page=%d, pageSize=%d", keyword, page, pageSize)

	// 调用商品服务 gRPC API
	resp, err := r.data.gc.GoodsList(ctx, &goodsV1.GoodsFilterRequest{
		KeyWords:    keyword,
		Pages:       page,
		PagePerNums: pageSize,
	})
	if err != nil {
		r.log.Errorf("搜索商品失败: keyword=%s, error=%v", keyword, err)
		return nil, 0, err
	}

	// 转换为业务对象
	goodsList := make([]*biz.Goods, 0, len(resp.Data))
	for _, item := range resp.Data {
		goodsList = append(goodsList, &biz.Goods{
			ID:              item.Id,
			Name:            item.Name,
			GoodsSn:         item.GoodsSn,
			ShopPrice:       item.ShopPrice,
			GoodsFrontImage: item.GoodsFrontImage,
			IsNew:           item.IsNew,
			IsHot:           item.IsHot,
		})
	}

	return goodsList, resp.Total, nil
}

// BatchGetGoods 批量获取商品信息
func (r *goodsRepo) BatchGetGoods(ctx context.Context, ids []int32) (map[int32]*biz.Goods, error) {
	r.log.Infof("批量获取商品: ids=%v", ids)

	// 调用商品服务 gRPC API
	resp, err := r.data.gc.BatchGetGoods(ctx, &goodsV1.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		r.log.Errorf("批量获取商品失败: ids=%v, error=%v", ids, err)
		return nil, err
	}

	// 转换为 map
	goodsMap := make(map[int32]*biz.Goods)
	for _, item := range resp.Data {
		goodsMap[item.Id] = &biz.Goods{
			ID:              item.Id,
			Name:            item.Name,
			ShopPrice:       item.ShopPrice,
			GoodsFrontImage: item.GoodsFrontImage,
		}
	}

	return goodsMap, nil
}
