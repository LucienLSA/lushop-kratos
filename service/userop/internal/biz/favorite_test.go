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

var _ = Describe("FavoriteUsecase", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockFavoriteRepo
		uc       *biz.FavoriteUsecase
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockFavoriteRepo(ctrl)
		uc = biz.NewFavoriteUsecase(mockRepo, nil)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("AddUserFav", func() {
		It("should add favorite successfully", func() {
			fav := domain.Favorite{
				UserID:  1,
				GoodsID: 100,
			}

			mockRepo.EXPECT().
				AddUserFav(ctx, fav).
				Return(nil)

			err := uc.AddUserFav(ctx, fav)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when add fails", func() {
			fav := domain.Favorite{UserID: 1, GoodsID: 100}

			mockRepo.EXPECT().
				AddUserFav(ctx, fav).
				Return(errors.New("already exists"))

			err := uc.AddUserFav(ctx, fav)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("DeleteUserFav", func() {
		It("should delete favorite successfully", func() {
			fav := domain.Favorite{UserID: 1, GoodsID: 100}

			mockRepo.EXPECT().
				DeleteUserFav(ctx, fav).
				Return(nil)

			err := uc.DeleteUserFav(ctx, fav)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("GetUserFavDetail", func() {
		It("should get favorite detail successfully", func() {
			fav := domain.Favorite{UserID: 1, GoodsID: 100}
			expectedFav := &domain.Favorite{
				UserID:  1,
				GoodsID: 100,
			}

			mockRepo.EXPECT().
				GetUserFavDetail(ctx, fav).
				Return(expectedFav, nil)

			result, err := uc.GetUserFavDetail(ctx, fav)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(expectedFav))
		})

		It("should return error when not found", func() {
			fav := domain.Favorite{UserID: 1, GoodsID: 999}

			mockRepo.EXPECT().
				GetUserFavDetail(ctx, fav).
				Return(nil, errors.New("not found"))

			result, err := uc.GetUserFavDetail(ctx, fav)
			Ω(err).Should(HaveOccurred())
			Ω(result).Should(BeNil())
		})
	})

	Context("GetFavList", func() {
		It("should get favorite list successfully", func() {
			filter := domain.Favorite{UserID: 1}
			fav1 := &domain.UserFavResponse{UserID: 1, GoodsID: 100}
			fav2 := &domain.UserFavResponse{UserID: 1, GoodsID: 200}
			expectedResp := &domain.UserFavListResponse{
				Total: 2,
				List:  []*domain.UserFavResponse{fav1, fav2},
			}

			mockRepo.EXPECT().
				GetFavList(ctx, filter).
				Return(expectedResp, nil)

			resp, err := uc.GetFavList(ctx, filter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).Should(Equal(expectedResp))
			Ω(resp.Total).Should(Equal(int64(2)))
		})
	})
})
