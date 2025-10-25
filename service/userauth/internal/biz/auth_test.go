package biz_test

import (
	"context"
	"errors"
	"io"
	"time"

	"userauth/internal/biz"
	"userauth/internal/conf"
	"userauth/internal/mocks"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthUsecase", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockAuthRepo
		uc       *biz.AuthUsecase
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockAuthRepo(ctrl)
		
		authConf := &conf.Auth{
			JwtKey: "test-jwt-key-for-testing-only-must-be-long-enough",
		}
		smsConf := &conf.Sms{
			ApiKey:       "test-api-key",
			ApiSecret:    "test-api-secret",
			SignName:     "测试签名",
			TemplateCode: "SMS_123456",
		}
		
		// 创建一个测试用的 logger
		logger := log.NewStdLogger(io.Discard)
		uc = biz.NewAuthUsecase(mockRepo, authConf, smsConf, logger)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("VerifyCaptcha", func() {
		It("should verify captcha successfully", func() {
			captchaId := "test-captcha-id"
			captcha := "1234"
			
			mockRepo.EXPECT().
				GetCaptcha(ctx, captchaId).
				Return(captcha, nil)
			
			valid, err := uc.VerifyCaptcha(ctx, captchaId, captcha)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(valid).Should(BeTrue())
		})

		It("should return false for invalid captcha", func() {
			captchaId := "test-captcha-id"
			
			mockRepo.EXPECT().
				GetCaptcha(ctx, captchaId).
				Return("1234", nil)
			
			valid, err := uc.VerifyCaptcha(ctx, captchaId, "wrong")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(valid).Should(BeFalse())
		})

		It("should return false when captcha not found", func() {
			captchaId := "test-captcha-id"
			
			mockRepo.EXPECT().
				GetCaptcha(ctx, captchaId).
				Return("", errors.New("not found"))
			
			valid, _ := uc.VerifyCaptcha(ctx, captchaId, "1234")
			Ω(valid).Should(BeFalse())
		})
	})

	Context("VerifySms", func() {
		It("should verify sms code successfully", func() {
			mobile := "13800138000"
			code := "123456"
			
			mockRepo.EXPECT().
				GetSmsCode(ctx, mobile).
				Return(code, nil)
			
			valid, err := uc.VerifySms(ctx, mobile, code)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(valid).Should(BeTrue())
		})

		It("should return false for invalid sms code", func() {
			mobile := "13800138000"
			
			mockRepo.EXPECT().
				GetSmsCode(ctx, mobile).
				Return("123456", nil)
			
			valid, err := uc.VerifySms(ctx, mobile, "wrong")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(valid).Should(BeFalse())
		})

		It("should return false when sms code not found", func() {
			mobile := "13800138000"
			
			mockRepo.EXPECT().
				GetSmsCode(ctx, mobile).
				Return("", errors.New("not found"))
			
			valid, _ := uc.VerifySms(ctx, mobile, "123456")
			Ω(valid).Should(BeFalse())
		})
	})

	Context("IssueToken", func() {
		It("should issue token successfully", func() {
			userId := int64(1)
			nickName := "testuser"
			role := int32(1)
			
			mockRepo.EXPECT().
				StoreAccessToken(ctx, userId, gomock.Any(), 30*time.Minute).
				Return(nil)
			
			mockRepo.EXPECT().
				StoreRefreshToken(ctx, userId, gomock.Any(), 7*24*time.Hour).
				Return(nil)
			
			reply, err := uc.IssueToken(ctx, userId, nickName, role)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(reply).ShouldNot(BeNil())
			Ω(reply.AccessToken).ShouldNot(BeEmpty())
			Ω(reply.RefreshToken).ShouldNot(BeEmpty())
			Ω(reply.ExpiredAt).Should(BeNumerically(">", time.Now().Unix()))
		})

		It("should return error when store access token fails", func() {
			userId := int64(1)
			nickName := "testuser"
			role := int32(1)
			
			mockRepo.EXPECT().
				StoreAccessToken(ctx, userId, gomock.Any(), 30*time.Minute).
				Return(errors.New("redis error"))
			
			reply, err := uc.IssueToken(ctx, userId, nickName, role)
			Ω(err).Should(HaveOccurred())
			Ω(reply).Should(BeNil())
		})
	})

	Context("RevokeToken", func() {
		It("should revoke token successfully", func() {
			userId := int64(1)
			
			mockRepo.EXPECT().
				DeleteTokens(ctx, userId).
				Return(nil)
			
			mockRepo.EXPECT().
				AddToBlacklist(ctx, userId, 30*time.Minute).
				Return(nil)
			
			err := uc.RevokeToken(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return error when add to blacklist fails", func() {
			userId := int64(1)
			
			mockRepo.EXPECT().
				DeleteTokens(ctx, userId).
				Return(nil)
			
			mockRepo.EXPECT().
				AddToBlacklist(ctx, userId, 30*time.Minute).
				Return(errors.New("redis error"))
			
			err := uc.RevokeToken(ctx, userId)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("CheckBlacklist", func() {
		It("should check blacklist successfully", func() {
			userId := int64(1)
			
			mockRepo.EXPECT().
				CheckBlacklist(ctx, userId).
				Return(true, nil)
			
			inBlacklist, err := uc.CheckBlacklist(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inBlacklist).Should(BeTrue())
		})

		It("should return false when user not in blacklist", func() {
			userId := int64(1)
			
			mockRepo.EXPECT().
				CheckBlacklist(ctx, userId).
				Return(false, nil)
			
			inBlacklist, err := uc.CheckBlacklist(ctx, userId)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(inBlacklist).Should(BeFalse())
		})
	})

	Context("AddToBlacklist", func() {
		It("should add to blacklist with custom TTL", func() {
			userId := int64(1)
			ttl := int64(3600)
			
			mockRepo.EXPECT().
				AddToBlacklist(ctx, userId, time.Duration(ttl)*time.Second).
				Return(nil)
			
			err := uc.AddToBlacklist(ctx, userId, ttl)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should add to blacklist with default TTL", func() {
			userId := int64(1)
			
			mockRepo.EXPECT().
				AddToBlacklist(ctx, userId, 30*time.Minute).
				Return(nil)
			
			err := uc.AddToBlacklist(ctx, userId, 0)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
