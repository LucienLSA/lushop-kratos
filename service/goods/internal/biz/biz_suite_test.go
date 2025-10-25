package biz_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ctl *gomock.Controller
	ctx context.Context
)

func TestBiz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Biz Suite")
}

var _ = BeforeSuite(func() {
	ctx = context.Background()
})

var _ = BeforeEach(func() {
	ctl = gomock.NewController(GinkgoT())
})

var _ = AfterEach(func() {
	ctl.Finish()
})
