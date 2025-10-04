package biz_test

import (
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
	It("Create", func() {
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
		Ω(err).ToNot(HaveOccurred())
		Ω(l.ID).To(Equal(int64(1)))
		Ω(l.Mobile).To(Equal("13803881388"))
	})
})
