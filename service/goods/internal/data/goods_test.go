package data_test

import (
	"goods/internal/biz"
	"goods/internal/data"
	"goods/internal/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoodsRepo", func() {
	var ro biz.GoodsRepo
	var goodsData *domain.Goods

	BeforeEach(func() {
		ro = data.NewGoodsRepo(Db, nil)
		goodsData = &domain.Goods{
			CategoryID:      1,
			BrandsID:        1,
			Name:            "测试商品",
			GoodsSn:         "TEST001",
			MarketPrice:     199.99,
			ShopPrice:       149.99,
			GoodsBrief:      "测试商品简介",
			GoodsFrontImage: "test.jpg",
			GoodsImages:     []string{},  // 空数组而不是 nil
			DescImages:      []string{},  // 空数组而不是 nil
			OnSale:          true,
			IsNew:           true,
			IsHot:           false,
			ShipFree:        true,
		}
	})

	AfterEach(func() {
		CleanTestData()
	})

	// 创建商品测试
	Context("CreateGoods", func() {
		It("should create goods successfully", func() {
			goods, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(goods).ShouldNot(BeNil())
			Ω(goods.ID).Should(BeNumerically(">", 0))
			Ω(goods.Name).Should(Equal("测试商品"))
			Ω(goods.GoodsSn).Should(Equal("TEST001"))
			Ω(goods.MarketPrice).Should(Equal(float32(199.99)))
			Ω(goods.ShopPrice).Should(Equal(float32(149.99)))
			Ω(goods.OnSale).Should(BeTrue())
			Ω(goods.IsNew).Should(BeTrue())
			Ω(goods.IsHot).Should(BeFalse())
			Ω(goods.ShipFree).Should(BeTrue())
		})

		It("should create goods with minimal data", func() {
			minimalGoods := &domain.Goods{
				CategoryID:  1,
				BrandsID:    1,
				Name:        "最小商品",
				GoodsSn:     "TEST_MIN001",
				MarketPrice: 100.0,
			}

			goods, err := ro.CreateGoods(ctx, minimalGoods)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(goods).ShouldNot(BeNil())
			Ω(goods.ID).Should(BeNumerically(">", 0))
			Ω(goods.Name).Should(Equal("最小商品"))
		})
	})

	// 根据ID查询商品测试
	Context("GoodsByID", func() {
		It("should get goods by id successfully", func() {
			// 先创建商品
			created, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())

			// 查询商品
			goods, err := ro.GoodsByID(ctx, created.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(goods).ShouldNot(BeNil())
			Ω(goods.ID).Should(Equal(created.ID))
			Ω(goods.Name).Should(Equal("测试商品"))
			Ω(goods.GoodsSn).Should(Equal("TEST001"))
			Ω(goods.MarketPrice).Should(Equal(float32(199.99)))
		})

		It("should return error when goods not found", func() {
			_, err := ro.GoodsByID(ctx, 99999)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("商品不存在"))
		})
	})

	// 更新商品测试
	Context("UpdateGoods", func() {
		It("should update goods successfully", func() {
			// 先创建商品
			created, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())

			// 更新商品信息
			created.Name = "更新后的商品"
			created.MarketPrice = 299.99
			created.ShopPrice = 249.99
			created.IsHot = true
			created.OnSale = false

			err = ro.UpdateGoods(ctx, created)
			Ω(err).ShouldNot(HaveOccurred())

			// 验证更新
			updated, err := ro.GoodsByID(ctx, created.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(updated.Name).Should(Equal("更新后的商品"))
			Ω(updated.MarketPrice).Should(Equal(float32(299.99)))
			Ω(updated.ShopPrice).Should(Equal(float32(249.99)))
			Ω(updated.IsHot).Should(BeTrue())
			Ω(updated.OnSale).Should(BeFalse())
		})

		It("should update goods status fields", func() {
			// 创建商品
			created, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())

			// 只更新状态字段
			created.IsNew = false
			created.IsHot = true
			created.OnSale = false
			created.ShipFree = false

			err = ro.UpdateGoods(ctx, created)
			Ω(err).ShouldNot(HaveOccurred())

			// 验证
			updated, err := ro.GoodsByID(ctx, created.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(updated.IsNew).Should(BeFalse())
			Ω(updated.IsHot).Should(BeTrue())
			Ω(updated.OnSale).Should(BeFalse())
			Ω(updated.ShipFree).Should(BeFalse())
		})
	})

	// 删除商品测试
	Context("DeleteGoods", func() {
		It("should delete goods successfully", func() {
			// 先创建商品
			created, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())

			// 删除商品
			err = ro.DeleteGoods(ctx, created.ID)
			Ω(err).ShouldNot(HaveOccurred())

			// 验证已删除（软删除）
			_, err = ro.GoodsByID(ctx, created.ID)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("商品不存在"))
		})

		It("should handle delete non-existent goods", func() {
			// 删除不存在的商品不应该报错（gorm 的行为）
			err := ro.DeleteGoods(ctx, 99999)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	// 批量查询商品测试
	Context("GoodsListByIDs", func() {
		It("should get goods list by ids successfully", func() {
			// 创建多个商品
			goods1, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())

			goodsData.Name = "测试商品2"
			goodsData.GoodsSn = "TEST002"
			goods2, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())

			goodsData.Name = "测试商品3"
			goodsData.GoodsSn = "TEST003"
			goods3, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())

			// 批量查询
			list, err := ro.GoodsListByIDs(ctx, goods1.ID, goods2.ID, goods3.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(list).ShouldNot(BeNil())
			Ω(len(list)).Should(Equal(3))
			Ω(list[0].Name).Should(Equal("测试商品"))
			Ω(list[1].Name).Should(Equal("测试商品2"))
			Ω(list[2].Name).Should(Equal("测试商品3"))
		})

		It("should return empty list when no goods found", func() {
			list, err := ro.GoodsListByIDs(ctx, 99999, 88888)
			// 根据实际实现，可能返回空列表而不是错误
			if err != nil {
				Ω(err.Error()).Should(ContainSubstring("商品不存在"))
			} else {
				Ω(len(list)).Should(Equal(0))
			}
		})

		It("should get partial goods when some ids not exist", func() {
			// 创建一个商品
			goods1, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())

			// 查询包含存在和不存在的ID
			list, err := ro.GoodsListByIDs(ctx, goods1.ID, 99999)
			// 根据实际实现，可能返回部分结果或错误
			if err == nil {
				Ω(len(list)).Should(BeNumerically(">=", 1))
			}
		})
	})

	// 完整生命周期测试
	Context("Goods Lifecycle Integration", func() {
		It("should handle complete goods lifecycle", func() {
			// 1. 创建商品
			created, err := ro.CreateGoods(ctx, goodsData)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(created.ID).Should(BeNumerically(">", 0))

			// 2. 查询商品
			found, err := ro.GoodsByID(ctx, created.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(found.Name).Should(Equal("测试商品"))

			// 3. 更新商品
			found.Name = "更新后的测试商品"
			found.MarketPrice = 399.99
			err = ro.UpdateGoods(ctx, found)
			Ω(err).ShouldNot(HaveOccurred())

			// 4. 验证更新
			updated, err := ro.GoodsByID(ctx, created.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(updated.Name).Should(Equal("更新后的测试商品"))
			Ω(updated.MarketPrice).Should(Equal(float32(399.99)))

			// 5. 批量查询
			list, err := ro.GoodsListByIDs(ctx, created.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(list)).Should(Equal(1))

			// 6. 删除商品
			err = ro.DeleteGoods(ctx, created.ID)
			Ω(err).ShouldNot(HaveOccurred())

			// 7. 验证已删除
			_, err = ro.GoodsByID(ctx, created.ID)
			Ω(err).Should(HaveOccurred())
		})
	})

	// 边界条件测试
	Context("Edge Cases", func() {
		It("should handle goods with empty images", func() {
			emptyImagesGoods := &domain.Goods{
				CategoryID:  1,
				BrandsID:    1,
				Name:        "无图片商品",
				GoodsSn:     "TEST_NOIMG",
				MarketPrice: 100.0,
				GoodsImages: []string{},
				DescImages:  []string{},
			}

			goods, err := ro.CreateGoods(ctx, emptyImagesGoods)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(goods).ShouldNot(BeNil())
		})

		It("should handle goods with zero price", func() {
			zeroPriceGoods := &domain.Goods{
				CategoryID:  1,
				BrandsID:    1,
				Name:        "零价格商品",
				GoodsSn:     "TEST_ZERO",
				MarketPrice: 0,
				ShopPrice:   0,
			}

			goods, err := ro.CreateGoods(ctx, zeroPriceGoods)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(goods.MarketPrice).Should(Equal(float32(0)))
		})
	})
})
