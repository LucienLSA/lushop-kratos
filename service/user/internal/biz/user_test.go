package biz_test

import (
	"errors"
	"time"
	"user/internal/biz"
	"user/internal/mocks/mrepo"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserUseCase", func() {
	var userCase *biz.UserUsecase
	var mUserRepo *mrepo.MockUserRepo

	BeforeEach(func() {
		// 隔离了数据层，只测试业务逻辑
		mUserRepo = mrepo.NewMockUserRepo(ctl)
		userCase = biz.NewUserUsecase(mUserRepo, nil)
	})
	// 创建用户
	Context("Create", func() {
		// 成功
		It("should create user successfully", func() {
			birthDay := time.Unix(int64(693646426), 0)
			info := &biz.User{
				ID:       1,
				Mobile:   "13803881388",
				Password: "123456",
				NickName: "lucien",
				Role:     1,
				Birthday: &birthDay,
			}
			mUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(info, nil)
			l, err := userCase.Create(ctx, info)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(l.ID).To(Equal(int64(1)))
			Ω(l.Mobile).To(Equal("13803881388"))
		})
		// 失败
		It("should return error when create user fails", func() {
			info := &biz.User{
				Mobile:   "13803881388",
				Password: "123456",
				NickName: "lucien",
			}
			expectedErr := errors.New("create user failed")
			mUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, expectedErr)
			_, err := userCase.Create(ctx, info)
			Ω(err).Should(HaveOccurred())
			Ω(err).To(Equal(expectedErr))
		})
	})
	// 手机号检查用户存在
	Context("UserByMobile", func() {
		// 成功
		It("should get user by mobile successfully", func() {
			birthDay := time.Unix(int64(693646426), 0)
			expectedUser := &biz.User{
				ID:       1,
				Mobile:   "13803881388",
				Password: "123456",
				NickName: "lucien",
				Role:     1,
				Birthday: &birthDay,
			}
			mUserRepo.EXPECT().UserByMobile(ctx, "13803881388").Return(expectedUser, nil)
			user, err := userCase.UserByMobile(ctx, "13803881388")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user).To(Equal(expectedUser))
		})
		// 失败
		It("should return error when user not found", func() {
			expectedErr := errors.New("user not found")
			mUserRepo.EXPECT().UserByMobile(ctx, "13803881388").Return(nil, expectedErr)
			_, err := userCase.UserByMobile(ctx, "13803881388")
			Ω(err).Should(HaveOccurred())
			Ω(err).To(Equal(expectedErr))
		})
	})
	// 获取用户列表
	Context("List", func() {
		It("should list users successfully", func() {
			birthDay := time.Unix(int64(693646426), 0)
			users := []*biz.User{
				{
					ID:       1,
					Mobile:   "13803881388",
					Password: "123456",
					NickName: "lucien",
					Role:     1,
					Birthday: &birthDay,
				},
				{
					ID:       2,
					Mobile:   "13803881389",
					Password: "123456",
					NickName: "alice",
					Role:     1,
					Birthday: &birthDay,
				},
			}
			total := 2
			mUserRepo.EXPECT().ListUser(ctx, 1, 10).Return(users, total, nil)
			result, count, err := userCase.List(ctx, 1, 10)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(result)).To(Equal(2))
			Ω(count).To(Equal(2))
			Ω(result[0].Mobile).To(Equal("13803881388"))
			Ω(result[1].Mobile).To(Equal("13803881389"))
		})

		It("should return error when list users fails", func() {
			expectedErr := errors.New("list users failed")
			mUserRepo.EXPECT().ListUser(ctx, 1, 10).Return(nil, 0, expectedErr)
			_, _, err := userCase.List(ctx, 1, 10)
			Ω(err).Should(HaveOccurred())
			Ω(err).To(Equal(expectedErr))
		})
	})
	// 更新用户
	Context("UpdateUser", func() {
		It("should update user successfully", func() {
			birthDay := time.Unix(int64(693646426), 0)
			user := &biz.User{
				ID:       1,
				Mobile:   "13803881388",
				Password: "123456",
				NickName: "lucien_updated",
				Role:     1,
				Birthday: &birthDay,
			}
			mUserRepo.EXPECT().UpdateUser(ctx, user).Return(true, nil)
			success, err := userCase.UpdateUser(ctx, user)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(success).To(BeTrue())
		})

		It("should return error when update user fails", func() {
			user := &biz.User{
				ID:       1,
				Mobile:   "13803881388",
				NickName: "lucien_updated",
			}
			expectedErr := errors.New("update user failed")
			mUserRepo.EXPECT().UpdateUser(ctx, user).Return(false, expectedErr)
			_, err := userCase.UpdateUser(ctx, user)
			Ω(err).Should(HaveOccurred())
			Ω(err).To(Equal(expectedErr))
		})
	})
	// 检查密码
	Context("CheckPassword", func() {
		It("should check password successfully when password is correct", func() {
			password := "123456"
			encryptedPassword := "encrypted_123456"
			mUserRepo.EXPECT().CheckPassword(ctx, password, encryptedPassword).Return(true, nil)
			valid, err := userCase.CheckPassword(ctx, password, encryptedPassword)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(valid).To(BeTrue())
		})

		It("should return false when password is incorrect", func() {
			password := "wrong_password"
			encryptedPassword := "encrypted_123456"
			mUserRepo.EXPECT().CheckPassword(ctx, password, encryptedPassword).Return(false, nil)
			valid, err := userCase.CheckPassword(ctx, password, encryptedPassword)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(valid).To(BeFalse())
		})

		It("should return error when check password fails", func() {
			password := "123456"
			encryptedPassword := "encrypted_123456"
			expectedErr := errors.New("check password failed")
			mUserRepo.EXPECT().CheckPassword(ctx, password, encryptedPassword).Return(false, expectedErr)
			_, err := userCase.CheckPassword(ctx, password, encryptedPassword)
			Ω(err).Should(HaveOccurred())
			Ω(err).To(Equal(expectedErr))
		})
	})
	// 根据ID查找用户
	Context("UserById", func() {
		It("should get user by id successfully", func() {
			birthDay := time.Unix(int64(693646426), 0)
			expectedUser := &biz.User{
				ID:       1,
				Mobile:   "13803881388",
				Password: "123456",
				NickName: "lucien",
				Role:     1,
				Birthday: &birthDay,
			}
			mUserRepo.EXPECT().GetUserById(ctx, int64(1)).Return(expectedUser, nil)
			user, err := userCase.UserById(ctx, 1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user).To(Equal(expectedUser))
		})

		It("should return error when user not found by id", func() {
			expectedErr := errors.New("user not found")
			mUserRepo.EXPECT().GetUserById(ctx, int64(999)).Return(nil, expectedErr)
			_, err := userCase.UserById(ctx, 999)
			Ω(err).Should(HaveOccurred())
			Ω(err).To(Equal(expectedErr))
		})
	})
})
