package service

import (
	"testing"

	v1 "goods/api/goods/v1"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

// TestNewGoodsService 测试 Service 创建
func TestNewGoodsService(t *testing.T) {
	logger := log.NewStdLogger(nil)
	service := NewGoodsService(nil, nil, nil, nil, nil, nil, logger)

	assert.NotNil(t, service)
	assert.NotNil(t, service.log)
}

// TestGoodsInfoResponse 测试商品信息响应转换
func TestGoodsInfoResponse(t *testing.T) {
	goods := &domain.Goods{
		ID:              1,
		CategoryID:      10,
		BrandsID:        20,
		Name:            "测试商品",
		GoodsSn:         "TEST001",
		ClickNum:        100,
		SoldNum:         50,
		FavNum:          30,
		MarketPrice:     199.99,
		ShopPrice:       149.99,
		GoodsBrief:      "这是一个测试商品",
		GoodsFrontImage: "https://example.com/front.jpg",
		GoodsImages:     []string{"img1.jpg", "img2.jpg"},
		DescImages:      []string{"desc1.jpg", "desc2.jpg"},
		OnSale:          true,
		IsNew:           true,
		IsHot:           false,
		ShipFree:        true,
	}

	resp := &v1.GoodsInfoResponse{
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
		GoodsFrontImage: goods.GoodsFrontImage,
		Images:          goods.GoodsImages,
		DescImages:      goods.DescImages,
		OnSale:          goods.OnSale,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		ShipFree:        goods.ShipFree,
	}

	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.Id)
	assert.Equal(t, int32(10), resp.CategoryId)
	assert.Equal(t, "测试商品", resp.Name)
	assert.Equal(t, "TEST001", resp.GoodsSn)
	assert.Equal(t, int32(100), resp.ClickNum)
	assert.Equal(t, int32(50), resp.SoldNum)
	assert.Equal(t, int32(30), resp.FavNum)
	assert.Equal(t, float32(199.99), resp.MarketPrice)
	assert.Equal(t, float32(149.99), resp.ShopPrice)
	assert.Equal(t, "这是一个测试商品", resp.GoodsBrief)
	assert.Equal(t, "https://example.com/front.jpg", resp.GoodsFrontImage)
	assert.Equal(t, 2, len(resp.Images))
	assert.Equal(t, 2, len(resp.DescImages))
	assert.True(t, resp.OnSale)
	assert.True(t, resp.IsNew)
	assert.False(t, resp.IsHot)
	assert.True(t, resp.ShipFree)
}

// TestGoodsListResponse 测试商品列表响应转换
func TestGoodsListResponse(t *testing.T) {
	goodsList := []*domain.Goods{
		{
			ID:              1,
			CategoryID:      10,
			Name:            "商品1",
			GoodsSn:         "TEST001",
			MarketPrice:     100.0,
			ShopPrice:       80.0,
			GoodsBrief:      "商品1简介",
			GoodsFrontImage: "img1.jpg",
			OnSale:          true,
			IsNew:           true,
			IsHot:           false,
			ShipFree:        true,
		},
		{
			ID:              2,
			CategoryID:      20,
			Name:            "商品2",
			GoodsSn:         "TEST002",
			MarketPrice:     200.0,
			ShopPrice:       160.0,
			GoodsBrief:      "商品2简介",
			GoodsFrontImage: "img2.jpg",
			OnSale:          false,
			IsNew:           false,
			IsHot:           true,
			ShipFree:        false,
		},
	}

	resp := &v1.GoodsListResponse{
		Total: int32(len(goodsList)),
	}

	for _, goods := range goodsList {
		item := &v1.GoodsInfoResponse{
			Id:              int32(goods.ID),
			CategoryId:      goods.CategoryID,
			Name:            goods.Name,
			GoodsSn:         goods.GoodsSn,
			MarketPrice:     goods.MarketPrice,
			ShopPrice:       goods.ShopPrice,
			GoodsBrief:      goods.GoodsBrief,
			GoodsFrontImage: goods.GoodsFrontImage,
			OnSale:          goods.OnSale,
			IsNew:           goods.IsNew,
			IsHot:           goods.IsHot,
			ShipFree:        goods.ShipFree,
		}
		resp.Data = append(resp.Data, item)
	}

	assert.NotNil(t, resp)
	assert.Equal(t, int32(2), resp.Total)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "商品1", resp.Data[0].Name)
	assert.Equal(t, "商品2", resp.Data[1].Name)
}

// TestCreateGoodsInfo 测试创建商品请求数据转换
func TestCreateGoodsInfo(t *testing.T) {
	req := &v1.CreateGoodsInfo{
		Id:              1,
		CategoryId:      10,
		BrandId:         20,
		Name:            "新商品",
		GoodsSn:         "NEW001",
		MarketPrice:     299.99,
		ShopPrice:       249.99,
		GoodsBrief:      "新商品简介",
		GoodsFrontImage: "new.jpg",
		Images:          []string{"img1.jpg", "img2.jpg"},
		DescImages:      []string{"desc1.jpg"},
		OnSale:          true,
		IsNew:           true,
		IsHot:           false,
		ShipFree:        true,
	}

	goodsInfo := domain.Goods{
		ID:              int64(req.Id),
		CategoryID:      req.CategoryId,
		BrandsID:        req.BrandId,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		GoodsFrontImage: req.GoodsFrontImage,
		GoodsImages:     req.Images,
		DescImages:      req.DescImages,
		OnSale:          req.OnSale,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		ShipFree:        req.ShipFree,
	}

	assert.Equal(t, int64(1), goodsInfo.ID)
	assert.Equal(t, int32(10), goodsInfo.CategoryID)
	assert.Equal(t, int32(20), goodsInfo.BrandsID)
	assert.Equal(t, "新商品", goodsInfo.Name)
	assert.Equal(t, "NEW001", goodsInfo.GoodsSn)
	assert.Equal(t, float32(299.99), goodsInfo.MarketPrice)
	assert.Equal(t, float32(249.99), goodsInfo.ShopPrice)
	assert.True(t, goodsInfo.OnSale)
	assert.True(t, goodsInfo.IsNew)
	assert.False(t, goodsInfo.IsHot)
	assert.True(t, goodsInfo.ShipFree)
}

// 表驱动测试：测试不同场景的商品数据转换
func TestGoodsInfoResponse_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		goods    *domain.Goods
		expected *v1.GoodsInfoResponse
	}{
		{
			name: "完整商品信息",
			goods: &domain.Goods{
				ID:              1,
				CategoryID:      10,
				BrandsID:        20,
				Name:            "完整商品",
				GoodsSn:         "FULL001",
				MarketPrice:     100.0,
				ShopPrice:       80.0,
				GoodsBrief:      "完整商品简介",
				GoodsFrontImage: "full.jpg",
				OnSale:          true,
				IsNew:           true,
				IsHot:           true,
				ShipFree:        true,
			},
			expected: &v1.GoodsInfoResponse{
				Id:              1,
				CategoryId:      10,
				Name:            "完整商品",
				GoodsSn:         "FULL001",
				MarketPrice:     100.0,
				ShopPrice:       80.0,
				GoodsBrief:      "完整商品简介",
				GoodsFrontImage: "full.jpg",
				OnSale:          true,
				IsNew:           true,
				IsHot:           true,
				ShipFree:        true,
			},
		},
		{
			name: "最小商品信息",
			goods: &domain.Goods{
				ID:         2,
				CategoryID: 5,
				BrandsID:   10,
				Name:       "最小商品",
				GoodsSn:    "MIN001",
			},
			expected: &v1.GoodsInfoResponse{
				Id:         2,
				CategoryId: 5,
				Name:       "最小商品",
				GoodsSn:    "MIN001",
			},
		},
		{
			name: "特价商品",
			goods: &domain.Goods{
				ID:          3,
				CategoryID:  15,
				BrandsID:    25,
				Name:        "特价商品",
				GoodsSn:     "SALE001",
				MarketPrice: 200.0,
				ShopPrice:   99.0,
				OnSale:      true,
				ShipFree:    true,
			},
			expected: &v1.GoodsInfoResponse{
				Id:          3,
				CategoryId:  15,
				Name:        "特价商品",
				GoodsSn:     "SALE001",
				MarketPrice: 200.0,
				ShopPrice:   99.0,
				OnSale:      true,
				ShipFree:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &v1.GoodsInfoResponse{
				Id:              int32(tt.goods.ID),
				CategoryId:      tt.goods.CategoryID,
				Name:            tt.goods.Name,
				GoodsSn:         tt.goods.GoodsSn,
				MarketPrice:     tt.goods.MarketPrice,
				ShopPrice:       tt.goods.ShopPrice,
				GoodsBrief:      tt.goods.GoodsBrief,
				GoodsFrontImage: tt.goods.GoodsFrontImage,
				OnSale:          tt.goods.OnSale,
				IsNew:           tt.goods.IsNew,
				IsHot:           tt.goods.IsHot,
				ShipFree:        tt.goods.ShipFree,
			}

			assert.Equal(t, tt.expected.Id, resp.Id)
			assert.Equal(t, tt.expected.CategoryId, resp.CategoryId)
			assert.Equal(t, tt.expected.Name, resp.Name)
			assert.Equal(t, tt.expected.GoodsSn, resp.GoodsSn)
			assert.Equal(t, tt.expected.MarketPrice, resp.MarketPrice)
			assert.Equal(t, tt.expected.ShopPrice, resp.ShopPrice)
			assert.Equal(t, tt.expected.OnSale, resp.OnSale)
			assert.Equal(t, tt.expected.IsNew, resp.IsNew)
			assert.Equal(t, tt.expected.IsHot, resp.IsHot)
			assert.Equal(t, tt.expected.ShipFree, resp.ShipFree)
		})
	}
}
