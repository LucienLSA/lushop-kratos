package biz_test

import (
	"context"
	"errors"

	"userop/internal/biz"
	"userop/internal/domain"
	"userop/internal/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressUsecase", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockAddressRepo
		uc       *biz.AddressUsecase
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockAddressRepo(ctrl)
		uc = biz.NewAddressUsecase(mockRepo, nil)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetAddressList", func() {
		It("should get address list successfully", func() {
			filter := domain.Address{UserID: 1}
			addr1 := &domain.Address{ID: 1, UserID: 1, Province: "广东省"}
			addr2 := &domain.Address{ID: 2, UserID: 1, Province: "北京市"}
			expectedResp := &domain.AddressListResponse{
				Total: 2,
				List:  []*domain.Address{addr1, addr2},
			}

			mockRepo.EXPECT().
				GetAddressList(ctx, filter).
				Return(expectedResp, nil)

			resp, err := uc.GetAddressList(ctx, filter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).Should(Equal(expectedResp))
			Ω(resp.Total).Should(Equal(int64(2)))
		})

		It("should return error when repo fails", func() {
			filter := domain.Address{UserID: 1}

			mockRepo.EXPECT().
				GetAddressList(ctx, filter).
				Return(nil, errors.New("database error"))

			resp, err := uc.GetAddressList(ctx, filter)
			Ω(err).Should(HaveOccurred())
			Ω(resp).Should(BeNil())
		})
	})

	Context("CreateAddress", func() {
		It("should create address successfully", func() {
			address := domain.Address{
				UserID:       1,
				Province:     "广东省",
				City:         "深圳市",
				District:     "南山区",
				Address:      "科技园",
				SignerName:   "张三",
				SignerMobile: "13800138000",
			}

			mockRepo.EXPECT().
				CreateAddress(ctx, address).
				Return(nil)

			err := uc.CreateAddress(ctx, address)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when create fails", func() {
			address := domain.Address{UserID: 1}

			mockRepo.EXPECT().
				CreateAddress(ctx, address).
				Return(errors.New("create failed"))

			err := uc.CreateAddress(ctx, address)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("UpdateAddress", func() {
		It("should update address successfully", func() {
			address := domain.Address{
				ID:       1,
				UserID:   1,
				Province: "广东省",
			}

			mockRepo.EXPECT().
				UpdateAddress(ctx, address).
				Return(nil)

			err := uc.UpdateAddress(ctx, address)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("DeleteAddress", func() {
		It("should delete address successfully", func() {
			address := domain.Address{ID: 1, UserID: 1}

			mockRepo.EXPECT().
				DeleteAddress(ctx, address).
				Return(nil)

			err := uc.DeleteAddress(ctx, address)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("GetMessageList", func() {
		It("should get message list successfully", func() {
			filter := domain.Message{UserID: 1}
			msg1 := &domain.Message{ID: 1, UserID: 1, MessageType: 1}
			expectedResp := &domain.MessageListResponse{
				Total: 1,
				List:  []*domain.Message{msg1},
			}

			mockRepo.EXPECT().
				GetMessageList(ctx, filter).
				Return(expectedResp, nil)

			resp, err := uc.GetMessageList(ctx, filter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).Should(Equal(expectedResp))
		})
	})

	Context("CreateMessage", func() {
		It("should create message successfully", func() {
			msg := domain.Message{
				UserID:      1,
				MessageType: 1,
				Subject:     "测试消息",
				Message:     "这是一条测试消息",
			}

			mockRepo.EXPECT().
				CreateMessage(ctx, msg).
				Return(nil)

			err := uc.CreateMessage(ctx, msg)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
