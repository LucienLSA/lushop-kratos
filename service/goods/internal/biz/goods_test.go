package biz_test

import (
	"context"
	"errors"
	"goods/internal/biz"
	"goods/internal/domain"
	"goods/internal/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoodsUsecase", func() {
	var (
		goodsUsecase *biz.GoodsUsecase
		mockGoods    *mocks.MockGoodsRepo
		mockBrand    *mocks.MockBrandRepo
		mockCategory *mocks.MockCategoryRepo
		mockTx       *mocks.MockTransaction
		mockES       *mocks.MockEsGoodsRepo
		ctx          context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockGoods = mocks.NewMockGoodsRepo(ctl)
		mockBrand = mocks.NewMockBrandRepo(ctl)
		mockCategory = mocks.NewMockCategoryRepo(ctl)
		mockTx = mocks.NewMockTransaction(ctl)
		mockES = mocks.NewMockEsGoodsRepo(ctl)
		goodsUsecase = biz.NewGoodsUsecase(mockGoods, mockBrand, mockCategory, mockTx, mockES, nil)
	})

	// 创建商品测试
	Context("CreateGoods", func() {
		It("should create goods successfully", func() {
			// 准备测试数据
			goodsData := &domain.Goods{
				CategoryID:      1,
				BrandsID:        1,
				Name:            "测试商品",
				GoodsSn:         "TEST001",
				MarketPrice:     100.0,
				ShopPrice:       80.0,
				GoodsBrief:      "测试商品简介",
				GoodsFrontImage: "test.jpg",
				OnSale:          true,
				IsNew:           true,
				IsHot:           false,
				ShipFree:        true,
			}

			brand := &domain.Brand{
				ID:   1,
				Name: "测试品牌",
			}

			category := &domain.Category{
				ID:   1,
				Name: "测试分类",
			}

			createdGoods := &domain.Goods{
				ID:              1,
				CategoryID:      1,
				BrandsID:        1,
				Name:            "测试商品",
				GoodsSn:         "TEST001",
				MarketPrice:     100.0,
				ShopPrice:       80.0,
				GoodsBrief:      "测试商品简介",
				GoodsFrontImage: "test.jpg",
				OnSale:          true,
				IsNew:           true,
				IsHot:           false,
				ShipFree:        true,
			}

			// 设置 Mock 期望
			mockBrand.EXPECT().GetBrandByID(ctx, int32(1)).Return(brand, nil)
			mockCategory.EXPECT().GetCategoryByID(ctx, int32(1)).Return(category, nil)
			mockTx.EXPECT().ExecTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockGoods.EXPECT().CreateGoods(ctx, gomock.Any()).Return(createdGoods, nil)
				mockES.EXPECT().InsertEsGoods(ctx, gomock.Any()).Return(nil)
				return fn(ctx)
			})

			// 执行测试
			result, err := goodsUsecase.CreateGoods(ctx, goodsData)

			// 验证结果
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())
			Ω(result.GoodsID).Should(Equal(int64(1)))
		})

		It("should return error when brand not found", func() {
			goodsData := &domain.Goods{
				CategoryID: 1,
				BrandsID:   999,
				Name:       "测试商品",
			}

			mockBrand.EXPECT().GetBrandByID(ctx, int32(999)).Return(nil, errors.New("brand not found"))

			result, err := goodsUsecase.CreateGoods(ctx, goodsData)

			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("品牌不存在"))
			Ω(result).Should(BeNil())
		})

		It("should return error when category not found", func() {
			goodsData := &domain.Goods{
				CategoryID: 999,
				BrandsID:   1,
				Name:       "测试商品",
			}

			brand := &domain.Brand{ID: 1, Name: "测试品牌"}
			mockBrand.EXPECT().GetBrandByID(ctx, int32(1)).Return(brand, nil)
			mockCategory.EXPECT().GetCategoryByID(ctx, int32(999)).Return(nil, errors.New("category not found"))

			result, err := goodsUsecase.CreateGoods(ctx, goodsData)

			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("商品分类不存在"))
			Ω(result).Should(BeNil())
		})

		It("should return error when goods already exists", func() {
			goodsData := &domain.Goods{
				ID:         1,
				CategoryID: 1,
				BrandsID:   1,
				Name:       "测试商品",
			}

			brand := &domain.Brand{ID: 1, Name: "测试品牌"}
			category := &domain.Category{ID: 1, Name: "测试分类"}
			existingGoods := &domain.Goods{ID: 1}

			mockBrand.EXPECT().GetBrandByID(ctx, int32(1)).Return(brand, nil)
			mockCategory.EXPECT().GetCategoryByID(ctx, int32(1)).Return(category, nil)
			mockGoods.EXPECT().GoodsByID(ctx, int64(1)).Return(existingGoods, nil)

			result, err := goodsUsecase.CreateGoods(ctx, goodsData)

			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("商品已存在"))
			Ω(result).Should(BeNil())
		})
	})

	// 更新商品测试
	Context("UpdateGoods", func() {
		It("should update goods status successfully (without brand and category)", func() {
			goodsData := &domain.Goods{
				ID:     1,
				IsNew:  true,
				IsHot:  true,
				OnSale: true,
			}

			existingGoods := &domain.Goods{
				ID:         1,
				CategoryID: 1,
				BrandsID:   1,
				Name:       "测试商品",
				IsNew:      false,
				IsHot:      false,
				OnSale:     false,
			}

			mockGoods.EXPECT().GoodsByID(ctx, int64(1)).Return(existingGoods, nil)
			mockGoods.EXPECT().UpdateGoods(ctx, gomock.Any()).Return(nil)

			result, err := goodsUsecase.UpdateGoods(ctx, goodsData)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())
			Ω(result.GoodsID).Should(Equal(int64(1)))
		})

		It("should update goods with brand and category successfully", func() {
			goodsData := &domain.Goods{
				ID:              1,
				CategoryID:      2,
				BrandsID:        2,
				Name:            "更新后的商品",
				GoodsSn:         "TEST002",
				MarketPrice:     120.0,
				ShopPrice:       100.0,
				GoodsBrief:      "更新后的简介",
				GoodsFrontImage: "updated.jpg",
				IsNew:           true,
				IsHot:           true,
				OnSale:          true,
				ShipFree:        true,
			}

			brand := &domain.Brand{ID: 2, Name: "新品牌"}
			category := &domain.Category{ID: 2, Name: "新分类"}
			existingGoods := &domain.Goods{
				ID:         1,
				CategoryID: 1,
				BrandsID:   1,
				Name:       "旧商品",
			}

			mockBrand.EXPECT().GetBrandByID(ctx, int32(2)).Return(brand, nil)
			mockCategory.EXPECT().GetCategoryByID(ctx, int32(2)).Return(category, nil)
			mockGoods.EXPECT().GoodsByID(ctx, int64(1)).Return(existingGoods, nil)
			mockGoods.EXPECT().UpdateGoods(ctx, gomock.Any()).Return(nil)

			result, err := goodsUsecase.UpdateGoods(ctx, goodsData)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())
			Ω(result.GoodsID).Should(Equal(int64(1)))
		})

		It("should return error when goods not found", func() {
			goodsData := &domain.Goods{
				ID:     999,
				IsNew:  true,
				IsHot:  true,
				OnSale: true,
			}

			mockGoods.EXPECT().GoodsByID(ctx, int64(999)).Return(nil, errors.New("goods not found"))

			result, err := goodsUsecase.UpdateGoods(ctx, goodsData)

			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("商品不存在"))
			Ω(result).Should(BeNil())
		})

		It("should return error when brand not found for update", func() {
			goodsData := &domain.Goods{
				ID:         1,
				CategoryID: 1,
				BrandsID:   999,
			}

			mockBrand.EXPECT().GetBrandByID(ctx, int32(999)).Return(nil, errors.New("brand not found"))

			result, err := goodsUsecase.UpdateGoods(ctx, goodsData)

			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("品牌不存在"))
			Ω(result).Should(BeNil())
		})
	})

	// 批量获取商品测试
	Context("BatchGetGoods", func() {
		It("should get goods list successfully", func() {
			ids := []int64{1, 2, 3}
			goodsList := []*domain.Goods{
				{ID: 1, Name: "商品1"},
				{ID: 2, Name: "商品2"},
				{ID: 3, Name: "商品3"},
			}

			mockGoods.EXPECT().GoodsListByIDs(ctx, int64(1), int64(2), int64(3)).Return(goodsList, nil)

			result, err := goodsUsecase.BatchGetGoods(ctx, ids)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())
			Ω(result.Total).Should(Equal(int64(3)))
			Ω(len(result.List)).Should(Equal(3))
		})

		It("should return error when get goods list fails", func() {
			ids := []int64{1, 2, 3}

			mockGoods.EXPECT().GoodsListByIDs(ctx, int64(1), int64(2), int64(3)).Return(nil, errors.New("database error"))

			result, err := goodsUsecase.BatchGetGoods(ctx, ids)

			Ω(err).Should(HaveOccurred())
			Ω(result).Should(BeNil())
		})
	})

	// 删除商品测试
	Context("DeleteGoods", func() {
		It("should delete goods successfully", func() {
			mockGoods.EXPECT().DeleteGoods(ctx, int64(1)).Return(nil)

			err := goodsUsecase.DeleteGoods(ctx, 1)

			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when delete fails", func() {
			mockGoods.EXPECT().DeleteGoods(ctx, int64(1)).Return(errors.New("delete error"))

			err := goodsUsecase.DeleteGoods(ctx, 1)

			Ω(err).Should(HaveOccurred())
		})
	})

	// 根据ID获取商品测试
	Context("GetGoodsById", func() {
		It("should get goods by id successfully", func() {
			expectedGoods := &domain.Goods{
				ID:         1,
				Name:       "测试商品",
				CategoryID: 1,
				BrandsID:   1,
			}

			mockGoods.EXPECT().GoodsByID(ctx, int64(1)).Return(expectedGoods, nil)

			result, err := goodsUsecase.GetGoodsById(ctx, 1)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())
			Ω(result.ID).Should(Equal(int64(1)))
			Ω(result.Name).Should(Equal("测试商品"))
		})

		It("should return error when goods not found", func() {
			mockGoods.EXPECT().GoodsByID(ctx, int64(999)).Return(nil, errors.New("goods not found"))

			result, err := goodsUsecase.GetGoodsById(ctx, 999)

			Ω(err).Should(HaveOccurred())
			Ω(result).Should(BeNil())
		})
	})
})
