package ioutils_test

import (
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("ReadFillBuffer", func() {
	It("should read the whole reader", func() {
		b := make([]byte, 5)
		r := strings.NewReader("12345")
		n, err := ReadFillBuffer(r, b)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(5))
		Expect(string(b[:n])).To(Equal("12345"))
		n, err = ReadFillBuffer(r, b)
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(0))
	})

	It("should read the full buffer", func() {
		b := make([]byte, 5)
		r := strings.NewReader("123456")
		n, err := ReadFillBuffer(r, b)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(5))
		Expect(string(b[:n])).To(Equal("12345"))
		n, err = ReadFillBuffer(r, b)
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(1))
		Expect(string(b[:n])).To(Equal("6"))
	})

	It("should fail if the reader fails", func() {
		b := make([]byte, 5)
		readErr := errors.New("read error")
		r := io.MultiReader(
			strings.NewReader("123"),
			NewErrorReader(readErr),
		)
		n, err := ReadFillBuffer(r, b)
		Expect(err).To(Equal(readErr))
		Expect(n).To(Equal(3))
		Expect(string(b[:n])).To(Equal("123"))
	})
})
