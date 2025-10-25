package data_test

import (
	"context"
	"inventory/internal/biz"
	"inventory/internal/data"
	"inventory/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("InventoryRepo", func() {
	var (
		repo biz.InventoryRepo
		ctx  context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		
		// 创建 Data 实例
		d := data.NewTestData(testDB, testRedis)
		
		logger := log.DefaultLogger
		repo = data.NewInventoryRepo(d, logger)

		// 清理测试数据
		testDB.Exec("TRUNCATE TABLE inventory")
		testDB.Exec("TRUNCATE TABLE stockselldetail")
	})

	Context("AddInv", func() {
		It("should create new inventory", func() {
			inv := &domain.Inventory{
				Goods:  100,
				Stocks: 50,
			}

			err := repo.AddInv(ctx, inv)
			Ω(err).ShouldNot(HaveOccurred())

			// 验证数据已插入
			var dbInv data.Inventory
			err = testDB.Where("goods = ?", 100).First(&dbInv).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(dbInv.Goods).Should(Equal(int32(100)))
			Ω(dbInv.Stocks).Should(Equal(int32(50)))
		})

		It("should update existing inventory", func() {
			// 先插入一条记录
			dbInv := data.Inventory{Goods: 100, Stocks: 30}
			err := testDB.Create(&dbInv).Error
			Ω(err).ShouldNot(HaveOccurred())

			// 更新库存
			inv := &domain.Inventory{
				Goods:  100,
				Stocks: 80,
			}
			err = repo.AddInv(ctx, inv)
			Ω(err).ShouldNot(HaveOccurred())

			// 验证已更新
			var updated data.Inventory
			err = testDB.Where("goods = ?", 100).First(&updated).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(updated.Stocks).Should(Equal(int32(80)))
		})
	})

	Context("GetInvById", func() {
		It("should get inventory by goods id", func() {
			// 插入测试数据
			dbInv := data.Inventory{Goods: 200, Stocks: 100}
			err := testDB.Create(&dbInv).Error
			Ω(err).ShouldNot(HaveOccurred())

			// 查询
			inv, err := repo.GetInvById(ctx, 200)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inv).ShouldNot(BeNil())
			Ω(inv.Goods).Should(Equal(int32(200)))
			Ω(inv.Stocks).Should(Equal(int32(100)))
		})

		It("should return error when goods not found", func() {
			inv, err := repo.GetInvById(ctx, 999)
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(Equal(gorm.ErrRecordNotFound))
			Ω(inv).Should(BeNil())
		})
	})

	Context("Sell", func() {
		BeforeEach(func() {
			// 准备库存数据
			invs := []data.Inventory{
				{Goods: 100, Stocks: 50},
				{Goods: 200, Stocks: 30},
			}
			for _, inv := range invs {
				err := testDB.Create(&inv).Error
				Ω(err).ShouldNot(HaveOccurred())
			}
		})

		It("should sell inventory successfully", func() {
			sell := &domain.SellInfo{
				OrderSn: "ORDER001",
				GoodsInvInfo: []domain.GoodsInvInfo{
					{GoodsID: 100, Nums: 5},
					{GoodsID: 200, Nums: 3},
				},
			}

			err := repo.Sell(ctx, sell)
			Ω(err).ShouldNot(HaveOccurred())

			// 验证库存已扣减
			var inv1 data.Inventory
			err = testDB.Where("goods = ?", 100).First(&inv1).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inv1.Stocks).Should(Equal(int32(45))) // 50 - 5

			var inv2 data.Inventory
			err = testDB.Where("goods = ?", 200).First(&inv2).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inv2.Stocks).Should(Equal(int32(27))) // 30 - 3

			// 验证扣减明细已记录
			var detail data.StockSellDetail
			err = testDB.Where("order_sn = ?", "ORDER001").First(&detail).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(detail.Status).Should(Equal(int32(1)))
			Ω(len(detail.Detail)).Should(Equal(2))
		})

		It("should return error when stock not enough", func() {
			sell := &domain.SellInfo{
				OrderSn: "ORDER002",
				GoodsInvInfo: []domain.GoodsInvInfo{
					{GoodsID: 100, Nums: 100}, // 超过库存
				},
			}

			err := repo.Sell(ctx, sell)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("not enough stock"))
		})

		It("should return error when goods not found", func() {
			sell := &domain.SellInfo{
				OrderSn: "ORDER003",
				GoodsInvInfo: []domain.GoodsInvInfo{
					{GoodsID: 999, Nums: 1},
				},
			}

			err := repo.Sell(ctx, sell)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("Reback", func() {
		BeforeEach(func() {
			// 准备库存数据
			invs := []data.Inventory{
				{Goods: 100, Stocks: 40},
				{Goods: 200, Stocks: 20},
			}
			for _, inv := range invs {
				err := testDB.Create(&inv).Error
				Ω(err).ShouldNot(HaveOccurred())
			}

			// 先执行一次扣减
			sell := &domain.SellInfo{
				OrderSn: "ORDER_REBACK_001",
				GoodsInvInfo: []domain.GoodsInvInfo{
					{GoodsID: 100, Nums: 10},
					{GoodsID: 200, Nums: 5},
				},
			}
			err := repo.Sell(ctx, sell)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should reback inventory successfully", func() {
			err := repo.Reback(ctx, "ORDER_REBACK_001")
			Ω(err).ShouldNot(HaveOccurred())

			// 验证库存已归还
			var inv1 data.Inventory
			err = testDB.Where("goods = ?", 100).First(&inv1).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inv1.Stocks).Should(Equal(int32(40))) // 40 - 10 + 10 = 40

			var inv2 data.Inventory
			err = testDB.Where("goods = ?", 200).First(&inv2).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inv2.Stocks).Should(Equal(int32(20))) // 20 - 5 + 5 = 20

			// 验证扣减明细状态已更新
			var detail data.StockSellDetail
			err = testDB.Where("order_sn = ?", "ORDER_REBACK_001").First(&detail).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(detail.Status).Should(Equal(int32(2))) // 已归还
		})

		It("should be idempotent - reback twice should not error", func() {
			// 第一次归还
			err := repo.Reback(ctx, "ORDER_REBACK_001")
			Ω(err).ShouldNot(HaveOccurred())

			// 记录第一次归还后的库存
			var inv1 data.Inventory
			err = testDB.Where("goods = ?", 100).First(&inv1).Error
			Ω(err).ShouldNot(HaveOccurred())
			firstStock := inv1.Stocks

			// 第二次归还（幂等性测试）
			err = repo.Reback(ctx, "ORDER_REBACK_001")
			Ω(err).ShouldNot(HaveOccurred())

			// 验证库存没有重复归还
			var inv2 data.Inventory
			err = testDB.Where("goods = ?", 100).First(&inv2).Error
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inv2.Stocks).Should(Equal(firstStock)) // 库存应该保持不变
		})

		It("should handle non-existent order gracefully", func() {
			err := repo.Reback(ctx, "NON_EXISTENT_ORDER")
			Ω(err).ShouldNot(HaveOccurred()) // 幂等性：不存在的订单不报错
		})
	})
})
