package biz_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBiz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Biz Suite")
}
