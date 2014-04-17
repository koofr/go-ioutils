package ioutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIoutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ioutils Suite")
}
