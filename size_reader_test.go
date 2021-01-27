package ioutils_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("SizeReader", func() {
	It("should calculate reader size", func() {
		br := bytes.NewReader([]byte("123"))

		r := NewSizeReader(br)

		Expect(r.Size()).To(Equal(int64(0)))

		n, err := r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Size()).To(Equal(int64(1)))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Size()).To(Equal(int64(2)))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Size()).To(Equal(int64(3)))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(0))
		Expect(err).To(Equal(io.EOF))
		Expect(r.Size()).To(Equal(int64(3)))
	})
})
