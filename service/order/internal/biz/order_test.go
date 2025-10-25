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

var _ = Describe("OrderUsecase", func() {
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

	Context("CreateOrder", func() {
		It("should create order successfully", func() {
			req := &v1.OrderRequest{
				UserId:  1,
				Address: "广东省深圳市南山区科技园",
				Name:    "张三",
				Mobile:  "13800138000",
				Post:    "测试订单",
			}

			expectedOrder := &v1.OrderInfoResponse{
				Id:      1,
				OrderSn: "ORDER20250101123456",
				UserId:  1,
				Status:  "PAYING",
			}

			mockRepo.EXPECT().
				CreateOrder(ctx, req).
				Return(expectedOrder, nil)

			order, err := uc.CreateOrder(ctx, req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(order).Should(Equal(expectedOrder))
			Ω(order.OrderSn).Should(ContainSubstring("ORDER"))
		})

		It("should return error when UserId is missing", func() {
			req := &v1.OrderRequest{
				Address: "测试地址",
				Name:    "张三",
				Mobile:  "13800138000",
			}

			order, err := uc.CreateOrder(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(order).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("user id is required"))
		})

		It("should return error when Address is missing", func() {
			req := &v1.OrderRequest{
				UserId: 1,
				Name:   "张三",
				Mobile: "13800138000",
			}

			order, err := uc.CreateOrder(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(order).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("address is required"))
		})

		It("should return error when Name is missing", func() {
			req := &v1.OrderRequest{
				UserId:  1,
				Address: "测试地址",
				Mobile:  "13800138000",
			}

			order, err := uc.CreateOrder(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(order).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("name is required"))
		})

		It("should return error when Mobile is missing", func() {
			req := &v1.OrderRequest{
				UserId:  1,
				Address: "测试地址",
				Name:    "张三",
			}

			order, err := uc.CreateOrder(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(order).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("mobile is required"))
		})

		It("should return error when repo fails", func() {
			req := &v1.OrderRequest{
				UserId:  1,
				Address: "测试地址",
				Name:    "张三",
				Mobile:  "13800138000",
			}

			mockRepo.EXPECT().
				CreateOrder(ctx, req).
				Return(nil, errors.New("database error"))

			order, err := uc.CreateOrder(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(order).Should(BeNil())
		})
	})

	Context("GetOrderList", func() {
		It("should get order list successfully", func() {
			req := &v1.OrderFilterRequest{
				UserId: 1,
				Pages:  1,
			}

			expectedOrders := []*v1.OrderInfoResponse{
				{Id: 1, OrderSn: "ORDER001", UserId: 1},
				{Id: 2, OrderSn: "ORDER002", UserId: 1},
			}

			mockRepo.EXPECT().
				GetOrderList(ctx, req).
				Return(expectedOrders, int32(2), nil)

			orders, total, err := uc.GetOrderList(ctx, req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(orders).Should(HaveLen(2))
			Ω(total).Should(Equal(int32(2)))
		})

		It("should return error when UserId is missing", func() {
			req := &v1.OrderFilterRequest{
				Pages: 1,
			}

			orders, total, err := uc.GetOrderList(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(orders).Should(BeNil())
			Ω(total).Should(Equal(int32(0)))
		})

		It("should set default page to 1 when Pages is 0", func() {
			req := &v1.OrderFilterRequest{
				UserId: 1,
				Pages:  0,
			}

			mockRepo.EXPECT().
				GetOrderList(ctx, gomock.Any()).
				Return([]*v1.OrderInfoResponse{}, int32(0), nil)

			_, _, err := uc.GetOrderList(ctx, req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(req.Pages).Should(Equal(int32(1)))
		})
	})

	Context("GetOrderDetail", func() {
		It("should get order detail successfully", func() {
			req := &v1.OrderRequest{
				Id:     1,
				UserId: 1,
			}

			expectedDetail := &v1.OrderInfoDetailResponse{
				OrderInfo: &v1.OrderInfoResponse{
					Id:      1,
					OrderSn: "ORDER001",
					UserId:  1,
					Status:  "PAID",
				},
			}

			mockRepo.EXPECT().
				GetOrderDetail(ctx, req).
				Return(expectedDetail, nil)

			detail, err := uc.GetOrderDetail(ctx, req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(detail).Should(Equal(expectedDetail))
		})

		It("should return error when Id is missing", func() {
			req := &v1.OrderRequest{
				UserId: 1,
			}

			detail, err := uc.GetOrderDetail(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(detail).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("order id is required"))
		})

		It("should return error when UserId is missing", func() {
			req := &v1.OrderRequest{
				Id: 1,
			}

			detail, err := uc.GetOrderDetail(ctx, req)
			Ω(err).Should(HaveOccurred())
			Ω(detail).Should(BeNil())
			Ω(err.Error()).Should(ContainSubstring("user id is required"))
		})
	})
})
