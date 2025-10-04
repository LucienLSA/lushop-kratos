package data_test

import (
	"time"
	"user/internal/biz"
	"user/internal/data"
	"user/internal/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// 用户测试用例
var _ = Describe("User", func() {
	var ro biz.UserRepo
	var uD *biz.User
	BeforeEach(func() {
		//  Db 是 data_suite_test.go 文件里面定义的
		ro = data.NewUserRepo(Db, nil)
		uD = testdata.User()
	})
	// 设置It 块来添加单个规格
	It("CreateUser", func() {
		u, err := ro.CreateUser(ctx, uD)
		Ω(err).ShouldNot(HaveOccurred())
		// 组装数据
		Ω(u.Mobile).Should(Equal("13803881388"))
	})
	It("ListUser", func() {
		user, total, err := ro.ListUser(ctx, 1, 10)
		Ω(err).ShouldNot(HaveOccurred()) // 获取列表不应该出现错误
		Ω(user).ShouldNot(BeEmpty())     // 结果不应该为空
		Ω(total).Should(Equal(1))        // 总数应该为 1，因为上面只创建了一条
		Ω(len(user)).Should(Equal(1))
		Ω(user[0].Mobile).Should(Equal("13803881388"))
	})
	It("UpdateUser", func() {
		birthDay := time.Unix(int64(693646426), 0)
		uD.NickName = "lucien"
		uD.Birthday = &birthDay
		uD.Gender = "male"
		user, err := ro.UpdateUser(ctx, uD)
		Ω(err).ShouldNot(HaveOccurred()) // 更新不应该出现错误
		Ω(user).Should(BeTrue())         // 结果应该为 true
	})
	It("CheckPassword", func() {
		p1 := "123456"
		encryptedPassword := "$2a$12$uhBtaYXOsfgE6l/lUcIarOlvUlbgWUBLWKY0Kx85PddtZgnoyn3Wy"
		password, err := ro.CheckPassword(ctx, p1, encryptedPassword)
		Ω(err).ShouldNot(HaveOccurred()) // 密码验证通过
		Ω(password).Should(BeTrue())     // 结果应该为true
		encryptedPassword1 := "$pbkdf2-sha512$5p7doUNIS9I5mvhA$b18171ff58b04c02ed70ea4f39"
		password1, err := ro.CheckPassword(ctx, p1, encryptedPassword1)
		if err != nil {
			return
		}
		Ω(err).ShouldNot(HaveOccurred())
		Ω(password1).Should(BeFalse()) // 密码验证不通过
	})
})
