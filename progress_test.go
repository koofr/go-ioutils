package ioutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("ProgressReader", func() {
	It("should show progress for reader", func() {
		f, err := NewTempFile("progress")
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()

		_, err = f.Write([]byte("0123456789"))
		Expect(err).NotTo(HaveOccurred())
		_, err = f.Seek(0, 0)
		Expect(err).NotTo(HaveOccurred())

		r, err := NewProgressReader(f)
		Expect(err).NotTo(HaveOccurred())

		_, err = r.Read([]byte{0})
		Expect(err).NotTo(HaveOccurred())

		_, err = r.Read([]byte{0, 0, 0, 0, 0})
		Expect(err).NotTo(HaveOccurred())

		_, err = r.Read([]byte{0, 0, 0, 0, 0})
		Expect(err).NotTo(HaveOccurred())
	})
})
