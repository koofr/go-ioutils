package ioutils_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("EofReader", func() {
	It("should set Eof to true when there it nothing left to read", func() {
		br := bytes.NewReader([]byte("123"))

		r := NewEofReader(br)

		Expect(r.Eof).To(BeFalse())

		n, err := r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Eof).To(BeFalse())

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Eof).To(BeFalse())

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Eof).To(BeFalse())

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(0))
		Expect(err).To(Equal(io.EOF))
		Expect(r.Eof).To(BeTrue())
	})
})
