package data_test

import (
	"context"

	"userop/internal/biz"
	"userop/internal/data"
	"userop/internal/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressRepo", func() {
	var (
		repo biz.AddressRepo
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = data.NewAddressRepo(dataObj, nil)
		ctx = context.Background()
	})

	Context("Address CRUD", func() {
		It("should create address successfully", func() {
			address := domain.Address{
				UserID:       9001,
				Province:     "广东省",
				City:         "深圳市",
				District:     "南山区",
				Address:      "科技园",
				SignerName:   "TEST张三",
				SignerMobile: "13800138000",
			}

			err := repo.CreateAddress(ctx, address)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should get address list successfully", func() {
			// 先创建测试数据
			address := domain.Address{
				UserID:       9002,
				Province:     "北京市",
				City:         "北京市",
				District:     "朝阳区",
				Address:      "望京",
				SignerName:   "TEST李四",
				SignerMobile: "13900139000",
			}
			err := repo.CreateAddress(ctx, address)
			Ω(err).ShouldNot(HaveOccurred())

			// 查询列表
			filter := domain.Address{UserID: 9002}
			resp, err := repo.GetAddressList(ctx, filter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).ShouldNot(BeNil())
			Ω(resp.Total).Should(BeNumerically(">", 0))
		})

		It("should update address successfully", func() {
			// 先创建
			address := domain.Address{
				UserID:       9003,
				Province:     "上海市",
				SignerName:   "TEST王五",
				SignerMobile: "13700137000",
			}
			err := repo.CreateAddress(ctx, address)
			Ω(err).ShouldNot(HaveOccurred())

			// 获取创建的地址
			filter := domain.Address{UserID: 9003}
			resp, err := repo.GetAddressList(ctx, filter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(resp.List)).Should(BeNumerically(">", 0))

			// 更新
			if len(resp.List) > 0 {
				updateAddr := *resp.List[0]
				updateAddr.Province = "浙江省"
				err = repo.UpdateAddress(ctx, updateAddr)
				Ω(err).ShouldNot(HaveOccurred())
			}
		})

		It("should delete address successfully", func() {
			// 先创建
			address := domain.Address{
				UserID:       9004,
				Province:     "广东省",
				SignerName:   "TEST赵六",
				SignerMobile: "13600136000",
			}
			err := repo.CreateAddress(ctx, address)
			Ω(err).ShouldNot(HaveOccurred())

			// 获取创建的地址
			filter := domain.Address{UserID: 9004}
			resp, err := repo.GetAddressList(ctx, filter)
			Ω(err).ShouldNot(HaveOccurred())

			// 删除
			if len(resp.List) > 0 {
				err = repo.DeleteAddress(ctx, *resp.List[0])
				Ω(err).ShouldNot(HaveOccurred())
			}
		})
	})

	Context("Message CRUD", func() {
		It("should create message successfully", func() {
			msg := domain.Message{
				UserID:      9001,
				MessageType: 1,
				Subject:     "TEST消息",
				Message:     "这是一条测试消息",
			}

			err := repo.CreateMessage(ctx, msg)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should get message list successfully", func() {
			// 先创建
			msg := domain.Message{
				UserID:      9002,
				MessageType: 1,
				Subject:     "TEST消息2",
				Message:     "测试消息内容",
			}
			err := repo.CreateMessage(ctx, msg)
			Ω(err).ShouldNot(HaveOccurred())

			// 查询
			filter := domain.Message{UserID: 9002}
			resp, err := repo.GetMessageList(ctx, filter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).ShouldNot(BeNil())
			Ω(resp.Total).Should(BeNumerically(">", 0))
		})
	})
})
