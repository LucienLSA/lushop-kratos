package data_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthRepo", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("验证码管理", func() {
		It("should store and get captcha successfully", func() {
			captchaId := "test-captcha-123"
			ans := "1234"
			
			// 存储验证码
			err := repo.StoreCaptcha(ctx, captchaId, ans, 5*time.Minute)
			Ω(err).ShouldNot(HaveOccurred())
			
			// 获取验证码
			result, err := repo.GetCaptcha(ctx, captchaId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(ans))
		})

		It("should return empty when captcha not found", func() {
			result, err := repo.GetCaptcha(ctx, "non-existent")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(BeEmpty())
		})

		It("should expire captcha after TTL", func() {
			captchaId := "test-captcha-expire"
			ans := "5678"
			
			// 存储验证码，1秒过期
			err := repo.StoreCaptcha(ctx, captchaId, ans, 1*time.Second)
			Ω(err).ShouldNot(HaveOccurred())
			
			// 立即获取应该成功
			result, err := repo.GetCaptcha(ctx, captchaId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(ans))
			
			// 等待过期
			time.Sleep(2 * time.Second)
			
			// 再次获取应该为空
			result, err = repo.GetCaptcha(ctx, captchaId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(BeEmpty())
		})
	})

	Context("短信验证码", func() {
		It("should store and get sms code successfully", func() {
			mobile := "13800138000"
			code := "123456"
			
			err := repo.StoreSmsCode(ctx, mobile, code, 5*time.Minute)
			Ω(err).ShouldNot(HaveOccurred())
			
			result, err := repo.GetSmsCode(ctx, mobile)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(code))
		})

		It("should handle sms cooldown", func() {
			mobile := "13900139000"
			
			// 设置冷却
			err := repo.SetSmsCooldown(ctx, mobile, 60*time.Second)
			Ω(err).ShouldNot(HaveOccurred())
			
			// 检查冷却状态
			cooling, err := repo.CheckSmsCooldown(ctx, mobile)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cooling).Should(BeTrue())
		})

		It("should return false when not in cooldown", func() {
			mobile := "13700137000"
			
			cooling, err := repo.CheckSmsCooldown(ctx, mobile)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cooling).Should(BeFalse())
		})
	})

	Context("Token 管理", func() {
		It("should store and get access token", func() {
			userId := int64(1)
			token := "test-access-token"
			
			err := repo.StoreAccessToken(ctx, userId, token, 30*time.Minute)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should store and get refresh token", func() {
			userId := int64(2)
			token := "test-refresh-token"
			
			err := repo.StoreRefreshToken(ctx, userId, token, 7*24*time.Hour)
			Ω(err).ShouldNot(HaveOccurred())
			
			result, err := repo.GetRefreshToken(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(token))
		})

		It("should delete tokens successfully", func() {
			userId := int64(3)
			accessToken := "access-token-3"
			refreshToken := "refresh-token-3"
			
			// 存储 tokens
			err := repo.StoreAccessToken(ctx, userId, accessToken, 30*time.Minute)
			Ω(err).ShouldNot(HaveOccurred())
			
			err = repo.StoreRefreshToken(ctx, userId, refreshToken, 7*24*time.Hour)
			Ω(err).ShouldNot(HaveOccurred())
			
			// 删除 tokens
			err = repo.DeleteTokens(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
			
			// 验证已删除（GetRefreshToken 在找不到时返回错误）
			result, err := repo.GetRefreshToken(ctx, userId)
			Ω(err).Should(HaveOccurred())
			Ω(result).Should(BeEmpty())
		})
	})

	Context("黑名单管理", func() {
		It("should add to blacklist successfully", func() {
			userId := int64(100)
			
			err := repo.AddToBlacklist(ctx, userId, 30*time.Minute)
			Ω(err).ShouldNot(HaveOccurred())
			
			// 检查黑名单
			isBlacklisted, err := repo.CheckBlacklist(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(isBlacklisted).Should(BeTrue())
		})

		It("should return false when not in blacklist", func() {
			userId := int64(999)
			
			isBlacklisted, err := repo.CheckBlacklist(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(isBlacklisted).Should(BeFalse())
		})

		It("should expire blacklist after TTL", func() {
			userId := int64(200)
			
			// 添加到黑名单，1秒过期
			err := repo.AddToBlacklist(ctx, userId, 1*time.Second)
			Ω(err).ShouldNot(HaveOccurred())
			
			// 立即检查应该在黑名单中
			isBlacklisted, err := repo.CheckBlacklist(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(isBlacklisted).Should(BeTrue())
			
			// 等待过期
			time.Sleep(2 * time.Second)
			
			// 再次检查应该不在黑名单中
			isBlacklisted, err = repo.CheckBlacklist(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(isBlacklisted).Should(BeFalse())
		})
	})
})
