package data_test

import (
	"context"

	"userop/internal/biz"
	"userop/internal/data"
	"userop/internal/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FavoriteRepo", func() {
	var (
		repo biz.FavoriteRepo
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = data.NewFavoriteRepo(dataObj, nil)
		ctx = context.Background()
	})

	Context("Favorite CRUD", func() {
		It("should add favorite successfully", func() {
			fav := domain.Favorite{
				UserID:  9001,
				GoodsID: 100,
			}

			err := repo.AddUserFav(ctx, fav)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should get favorite detail successfully", func() {
			// 先添加
			fav := domain.Favorite{
				UserID:  9002,
				GoodsID: 200,
			}
			err := repo.AddUserFav(ctx, fav)
			Ω(err).ShouldNot(HaveOccurred())

			// 查询详情
			result, err := repo.GetUserFavDetail(ctx, fav)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())
			Ω(result.UserID).Should(Equal(int32(9002)))
			Ω(result.GoodsID).Should(Equal(int32(200)))
		})

		It("should get favorite list successfully", func() {
			// 添加多个收藏
			fav1 := domain.Favorite{UserID: 9003, GoodsID: 301}
			fav2 := domain.Favorite{UserID: 9003, GoodsID: 302}

			err := repo.AddUserFav(ctx, fav1)
			Ω(err).ShouldNot(HaveOccurred())

			err = repo.AddUserFav(ctx, fav2)
			Ω(err).ShouldNot(HaveOccurred())

			// 查询列表
			filter := domain.Favorite{UserID: 9003}
			resp, err := repo.GetFavList(ctx, filter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).ShouldNot(BeNil())
			Ω(resp.Total).Should(BeNumerically(">=", 2))
		})

		It("should delete favorite successfully", func() {
			// 先添加
			fav := domain.Favorite{
				UserID:  9004,
				GoodsID: 400,
			}
			err := repo.AddUserFav(ctx, fav)
			Ω(err).ShouldNot(HaveOccurred())

			// 删除
			err = repo.DeleteUserFav(ctx, fav)
			Ω(err).ShouldNot(HaveOccurred())

			// 验证已删除（查询应该返回错误或空）
			_, err = repo.GetUserFavDetail(ctx, fav)
			Ω(err).Should(HaveOccurred())
		})
	})
})
