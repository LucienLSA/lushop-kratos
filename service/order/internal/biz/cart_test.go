package biz_test

import (
	"context"
	"errors"
	"io"

	v1 "order/api/order/v1"
	"order/internal/biz"
	"order/internal/mocks"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CartUsecase", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockOrderRepo
		uc       *biz.OrderUsecase
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockOrderRepo(ctrl)
		
		// 创建一个 discard logger
		logger := log.NewStdLogger(io.Discard)
		uc = biz.NewOrderUsecase(mockRepo, logger)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetCartList", func() {
		It("should get cart list successfully", func() {
			userID := int32(1)
			expectedCarts := []*v1.ShopCartInfoResponse{
				{Id: 1, UserId: 1, GoodsId: 100, Nums: 2},
				{Id: 2, UserId: 1, GoodsId: 200, Nums: 1},
			}

			mockRepo.EXPECT().
				GetCartList(ctx, userID).
				Return(expectedCarts, nil)

			carts, err := uc.GetCartList(ctx, userID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(carts).Should(HaveLen(2))
		})

		It("should return error when repo fails", func() {
			userID := int32(1)

			mockRepo.EXPECT().
				GetCartList(ctx, userID).
				Return(nil, errors.New("database error"))

			carts, err := uc.GetCartList(ctx, userID)
			Ω(err).Should(HaveOccurred())
			Ω(carts).Should(BeNil())
		})
	})

	Context("CreateCartItem", func() {
		It("should create cart item successfully", func() {
			req := &v1.CartItemRequest{
				UserId:  1,
				GoodsId: 100,
				Nums:    2,
			}

			expectedCart := &v1.ShopCartInfoResponse{
				Id:      1,
				UserId:  1,
				GoodsId: 100,
				Nums:    2,
			}

			mockRepo.EXPECT().
				CreateCartItem(ctx, req).
				Return(expectedCart, nil)

			cart, err := uc.CreateCartItem(ctx, req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cart).Should(Equal(expectedCart))
		})

		It("should return error when UserId is missing", func() {
			req := &v1.CartItemRequest{
				GoodsId: 100,
				Nums:    2,
			}

			cart, err := uc.CreateCartItem(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(cart).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("user id is required"))
		})

		It("should return error when GoodsId is missing", func() {
			req := &v1.CartItemRequest{
				UserId: 1,
				Nums:   2,
			}

			cart, err := uc.CreateCartItem(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(cart).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("goods id is required"))
		})

		It("should return error when Nums is invalid", func() {
			req := &v1.CartItemRequest{
				UserId:  1,
				GoodsId: 100,
				Nums:    0,
			}

			cart, err := uc.CreateCartItem(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(cart).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("nums must be greater than 0"))
		})
	})

	Context("UpdateCartItem", func() {
		It("should update cart item successfully", func() {
			req := &v1.CartItemRequest{
				Id:      1,
				UserId:  1,
				GoodsId: 100,
				Nums:    3,
			}

			mockRepo.EXPECT().
				UpdateCartItem(ctx, req).
				Return(nil, nil)

			_, err := uc.UpdateCartItem(ctx, req)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when Id is missing", func() {
			req := &v1.CartItemRequest{
				UserId:  1,
				GoodsId: 100,
				Nums:    3,
			}

			_, err := uc.UpdateCartItem(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("cart item id is required"))
		})
	})

	Context("DeleteCartItem", func() {
		It("should delete cart item successfully", func() {
			req := &v1.CartItemRequest{
				Id:     1,
				UserId: 1,
			}

			mockRepo.EXPECT().
				DeleteCartItem(ctx, req).
				Return(nil, nil)

			_, err := uc.DeleteCartItem(ctx, req)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when Id is missing", func() {
			req := &v1.CartItemRequest{
				UserId: 1,
			}

			_, err := uc.DeleteCartItem(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("cart item id is required"))
		})

		It("should return error when UserId is missing", func() {
			req := &v1.CartItemRequest{
				Id: 1,
			}

			_, err := uc.DeleteCartItem(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("user id is required"))
		})
	})
})
