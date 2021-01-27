package ioutils_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("TempFile", func() {
	It("create temp file and remove it after close", func() {
		file, err := NewTempFile("tempfile-test-")
		Expect(err).NotTo(HaveOccurred())

		_, err = os.Stat(file.Name())
		Expect(err).NotTo(HaveOccurred())

		n, err := file.Write([]byte("123"))
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(3))

		_ = file.Read
		_ = file.Write
		_ = file.Seek

		err = file.Close()
		Expect(err).NotTo(HaveOccurred())

		_, err = os.Stat(file.Name())
		Expect(os.IsNotExist(err)).To(BeTrue())
	})
})
