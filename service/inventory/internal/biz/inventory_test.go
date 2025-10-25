package biz_test

import (
	"context"
	"errors"
	"io"

	"inventory/internal/biz"
	"inventory/internal/domain"
	"inventory/internal/mocks"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InventoryUsecase", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockInventoryRepo
		uc       *biz.InventoryUsecase
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockInventoryRepo(ctrl)

		// 创建一个 discard logger
		logger := log.NewStdLogger(io.Discard)
		uc = biz.NewInventoryUsecase(mockRepo, logger)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("SetInv", func() {
		It("should set inventory successfully", func() {
			inv := &domain.Inventory{
				Goods:  100,
				Stocks: 50,
			}

			mockRepo.EXPECT().
				AddInv(ctx, inv).
				Return(nil)

			err := uc.SetInv(ctx, inv)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when repo fails", func() {
			inv := &domain.Inventory{
				Goods:  100,
				Stocks: 50,
			}

			mockRepo.EXPECT().
				AddInv(ctx, inv).
				Return(errors.New("database error"))

			err := uc.SetInv(ctx, inv)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("database error"))
		})
	})

	Context("GetInvById", func() {
		It("should get inventory by id successfully", func() {
			goodsId := int32(100)
			expectedInv := &domain.Inventory{
				Goods:  100,
				Stocks: 50,
			}

			mockRepo.EXPECT().
				GetInvById(ctx, goodsId).
				Return(expectedInv, nil)

			inv, err := uc.GetInvById(ctx, goodsId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inv).ShouldNot(BeNil())
			Ω(inv.Goods).Should(Equal(int32(100)))
			Ω(inv.Stocks).Should(Equal(int32(50)))
		})

		It("should return error when goods not found", func() {
			goodsId := int32(999)

			mockRepo.EXPECT().
				GetInvById(ctx, goodsId).
				Return(nil, errors.New("record not found"))

			inv, err := uc.GetInvById(ctx, goodsId)
			Ω(err).Should(HaveOccurred())
			Ω(inv).Should(BeNil())
		})
	})

	Context("Sell", func() {
		It("should sell inventory successfully", func() {
			sell := &domain.SellInfo{
				GoodsInvInfo: []domain.GoodsInvInfo{
					{GoodsID: 100, Nums: 2},
					{GoodsID: 200, Nums: 3},
				},
				OrderSn: "ORDER001",
			}

			mockRepo.EXPECT().
				Sell(ctx, sell).
				Return(nil)

			err := uc.Sell(ctx, sell)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when stock not enough", func() {
			sell := &domain.SellInfo{
				GoodsInvInfo: []domain.GoodsInvInfo{
					{GoodsID: 100, Nums: 1000},
				},
				OrderSn: "ORDER002",
			}

			mockRepo.EXPECT().
				Sell(ctx, sell).
				Return(errors.New("not enough stock"))

			err := uc.Sell(ctx, sell)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("not enough stock"))
		})

		It("should return error when repo fails", func() {
			sell := &domain.SellInfo{
				GoodsInvInfo: []domain.GoodsInvInfo{
					{GoodsID: 100, Nums: 2},
				},
				OrderSn: "ORDER003",
			}

			mockRepo.EXPECT().
				Sell(ctx, sell).
				Return(errors.New("database error"))

			err := uc.Sell(ctx, sell)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("Reback", func() {
		It("should reback inventory successfully", func() {
			orderSn := "ORDER001"

			mockRepo.EXPECT().
				Reback(ctx, orderSn).
				Return(nil)

			err := uc.Reback(ctx, orderSn)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should handle idempotent reback", func() {
			orderSn := "ORDER002"

			mockRepo.EXPECT().
				Reback(ctx, orderSn).
				Return(nil)

			err := uc.Reback(ctx, orderSn)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when repo fails", func() {
			orderSn := "ORDER003"

			mockRepo.EXPECT().
				Reback(ctx, orderSn).
				Return(errors.New("database error"))

			err := uc.Reback(ctx, orderSn)
			Ω(err).Should(HaveOccurred())
		})
	})
})
