package data_test

import (
	"fmt"
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
	var testCounter int

	BeforeEach(func() {
		//  Db 是 data_suite_test.go 文件里面定义的
		ro = data.NewUserRepo(Db, nil)
		// 为每个测试生成唯一的手机号，避免重复
		testCounter++
		uD = testdata.User()
		uD.Mobile = fmt.Sprintf("1380388%04d", testCounter)
	})

	AfterEach(func() {
		// 清理测试数据，确保测试之间隔离
		if Db != nil {
			Db.CleanTestData()
		}
	})

	// 创建用户
	Context("CreateUser", func() {
		It("should create user successfully", func() {
			u, err := ro.CreateUser(ctx, uD)
			Ω(err).ShouldNot(HaveOccurred())
			// 组装数据
			Ω(u.Mobile).Should(Equal(uD.Mobile))
			Ω(u.NickName).Should(Equal("user1"))
			Ω(u.Gender).Should(Equal("male"))
			Ω(u.Role).Should(Equal(1))
		})

		It("should return error when user already exists", func() {
			// 先创建一个用户
			_, err := ro.CreateUser(ctx, uD)
			Ω(err).ShouldNot(HaveOccurred())

			// 尝试创建相同手机号的用户
			duplicateUser := testdata.User()
			duplicateUser.Mobile = uD.Mobile // 使用相同的手机号
			duplicateUser.NickName = "duplicate_user"
			_, err = ro.CreateUser(ctx, duplicateUser)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("用户已存在"))
		})
	})

	// 获取用户列表
	Context("ListUser", func() {
		BeforeEach(func() {
			// 创建测试数据
			_, err := ro.CreateUser(ctx, uD)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should list users successfully", func() {
			user, total, err := ro.ListUser(ctx, 1, 10)
			Ω(err).ShouldNot(HaveOccurred()) // 获取列表不应该出现错误
			Ω(user).ShouldNot(BeEmpty())     // 结果不应该为空
			Ω(total).Should(Equal(1))        // 总数应该为 1，因为上面只创建了一条
			Ω(len(user)).Should(Equal(1))
			Ω(user[0].Mobile).Should(Equal(uD.Mobile))
		})

		It("should handle pagination correctly", func() {
			// 创建第二个用户
			user2 := testdata.User()
			user2.Mobile = fmt.Sprintf("1380388%04d", testCounter+1)
			user2.NickName = "user2"
			_, err := ro.CreateUser(ctx, user2)
			Ω(err).ShouldNot(HaveOccurred())

			// 测试分页
			users, total, err := ro.ListUser(ctx, 1, 1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(users)).Should(Equal(1))
			Ω(total).Should(Equal(2))

			users, total, err = ro.ListUser(ctx, 2, 1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(users)).Should(Equal(1))
			Ω(total).Should(Equal(2))
		})

		It("should handle invalid page parameters", func() {
			// 测试无效的分页参数
			users, total, err := ro.ListUser(ctx, 0, 0)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(users)).Should(Equal(1)) // 默认页面大小应该是 10
			Ω(total).Should(Equal(1))

			users, total, err = ro.ListUser(ctx, 1, 101)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(users)).Should(Equal(1)) // 最大页面大小应该是 100
			Ω(total).Should(Equal(1))
		})
	})

	// 根据手机号查找用户
	Context("UserByMobile", func() {
		BeforeEach(func() {
			// 创建测试数据
			_, err := ro.CreateUser(ctx, uD)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should get user by mobile successfully", func() {
			user, err := ro.UserByMobile(ctx, uD.Mobile)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user.Mobile).Should(Equal(uD.Mobile))
			Ω(user.NickName).Should(Equal("user1"))
		})

		It("should return error when user not found", func() {
			_, err := ro.UserByMobile(ctx, "13803881399")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("user not found"))
		})
	})

	// 根据ID查找用户
	Context("GetUserById", func() {
		var createdUser *biz.User

		BeforeEach(func() {
			// 创建测试数据
			var err error
			createdUser, err = ro.CreateUser(ctx, uD)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should get user by id successfully", func() {
			user, err := ro.GetUserById(ctx, createdUser.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user.ID).Should(Equal(createdUser.ID))
			Ω(user.Mobile).Should(Equal(createdUser.Mobile))
			Ω(user.NickName).Should(Equal("user1"))
		})

		It("should return error when user not found by id", func() {
			_, err := ro.GetUserById(ctx, 99999)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("user not found"))
		})
	})

	// 更新用户
	Context("UpdateUser", func() {
		var createdUser *biz.User

		BeforeEach(func() {
			// 创建测试数据
			var err error
			createdUser, err = ro.CreateUser(ctx, uD)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should update user successfully", func() {
			birthDay := time.Unix(int64(693646426), 0)
			updateUser := &biz.User{
				ID:       createdUser.ID,
				NickName: "lucien_updated",
				Birthday: &birthDay,
				Gender:   "female",
			}
			success, err := ro.UpdateUser(ctx, updateUser)
			Ω(err).ShouldNot(HaveOccurred()) // 更新不应该出现错误
			Ω(success).Should(BeTrue())      // 结果应该为 true

			// 验证更新结果
			updatedUser, err := ro.GetUserById(ctx, createdUser.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(updatedUser.NickName).Should(Equal("lucien_updated"))
			Ω(updatedUser.Gender).Should(Equal("female"))
		})

		It("should return error when user not found for update", func() {
			updateUser := &biz.User{
				ID:       99999,
				NickName: "not_found_user",
			}
			success, err := ro.UpdateUser(ctx, updateUser)
			Ω(err).Should(HaveOccurred())
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("user not found"))
		})
	})

	// 检查密码
	Context("CheckPassword", func() {
		It("should check password successfully when password is correct", func() {
			p1 := "123456"
			encryptedPassword := "$2a$12$uhBtaYXOsfgE6l/lUcIarOlvUlbgWUBLWKY0Kx85PddtZgnoyn3Wy"
			password, err := ro.CheckPassword(ctx, p1, encryptedPassword)
			Ω(err).ShouldNot(HaveOccurred()) // 密码验证通过
			Ω(password).Should(BeTrue())     // 结果应该为true
		})

		It("should return false when password is incorrect", func() {
			p1 := "wrong_password"
			encryptedPassword := "$2a$12$uhBtaYXOsfgE6l/lUcIarOlvUlbgWUBLWKY0Kx85PddtZgnoyn3Wy"
			password, err := ro.CheckPassword(ctx, p1, encryptedPassword)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(password).Should(BeFalse()) // 密码验证不通过
		})

		It("should handle invalid encrypted password format", func() {
			p1 := "123456"
			encryptedPassword1 := "$pbkdf2-sha512$5p7doUNIS9I5mvhA$b18171ff58b04c02ed70ea4f39"
			password1, err := ro.CheckPassword(ctx, p1, encryptedPassword1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(password1).Should(BeFalse()) // 密码验证不通过
		})

		It("should handle empty password", func() {
			password, err := ro.CheckPassword(ctx, "", "some_encrypted_password")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(password).Should(BeFalse())
		})

		It("should handle empty encrypted password", func() {
			password, err := ro.CheckPassword(ctx, "123456", "")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(password).Should(BeFalse())
		})
	})

	// 集成测试：完整的用户生命周期
	Context("User Lifecycle Integration", func() {
		It("should handle complete user lifecycle", func() {
			// 1. 创建用户
			user, err := ro.CreateUser(ctx, uD)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user.Mobile).Should(Equal(uD.Mobile))

			// 2. 根据手机号查找用户
			foundUser, err := ro.UserByMobile(ctx, uD.Mobile)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(foundUser.ID).Should(Equal(user.ID))

			// 3. 根据ID查找用户
			foundUserById, err := ro.GetUserById(ctx, user.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(foundUserById.ID).Should(Equal(user.ID))

			// 4. 更新用户信息
			updateUser := &biz.User{
				ID:       user.ID,
				NickName: "updated_nickname",
				Gender:   "female",
			}
			success, err := ro.UpdateUser(ctx, updateUser)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(success).Should(BeTrue())

			// 5. 验证更新结果
			updatedUser, err := ro.GetUserById(ctx, user.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(updatedUser.NickName).Should(Equal("updated_nickname"))
			Ω(updatedUser.Gender).Should(Equal("female"))

			// 6. 验证密码
			valid, err := ro.CheckPassword(ctx, "123456", user.Password)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(valid).Should(BeTrue())

			// 7. 获取用户列表
			users, total, err := ro.ListUser(ctx, 1, 10)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(total).Should(Equal(1))
			Ω(len(users)).Should(Equal(1))
			Ω(users[0].ID).Should(Equal(user.ID))
		})
	})
})
